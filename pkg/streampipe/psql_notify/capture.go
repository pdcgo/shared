package psql_notify

import (
	"context"
	"fmt"
	"log"
	"net"
	"regexp"
	"time"

	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
	"github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/proxy"
	"github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBType string

const (
	DBTypePostgres         DBType = "postgres"
	DBTypeGcpCloudPostgres DBType = "gcpcloudpostgres"
)

type CapturePostgres struct {
	driverType DBType
	dsn        string
	err        error
}

func NewCapturePostgres(driverType DBType, dsn string) *CapturePostgres {
	return &CapturePostgres{
		driverType: driverType,
		dsn:        dsn,
	}
}

func (cdc *CapturePostgres) Listen(ctx context.Context, channel string, routineMode bool, handle func(raw string)) *CapturePostgres {
	if cdc.err != nil {
		return cdc
	}

	var listener *pq.Listener
	switch cdc.driverType {
	case DBTypePostgres:
		listener = pq.NewListener(cdc.dsn, time.Second*2, time.Hour, nil)
	case DBTypeGcpCloudPostgres:
		listener = pq.NewDialListener(dialer{}, cdc.dsn, time.Second*2, time.Hour, nil)
	}

	err := listener.Listen(channel)
	if err != nil {
		return cdc.setErr(err)
	}

	runner := func() {
		fmt.Printf("Listening for %s changes...\n", channel)
		defer listener.Close()

		for {
			select {
			case <-ctx.Done():
				return
			case notification := <-listener.Notify:
				if notification != nil {
					handle(notification.Extra)
				}
			default:
				if cdc.err != nil {
					return
				}
			}
		}

	}

	if routineMode {
		go runner()
	} else {
		runner()
	}

	return cdc
}

func (cdc *CapturePostgres) PrepareChannel(tableName, channelName string) *CapturePostgres {
	if cdc.err != nil {
		return cdc
	}

	return cdc.dbops(func(tx *gorm.DB) error {
		var err error
		err = tx.Transaction(func(tx *gorm.DB) error {
			err = tx.Exec(fmt.Sprintf(`LOCK TABLE %s IN ACCESS EXCLUSIVE MODE`, tableName)).Error
			if err != nil {
				return err
			}
			return cdc.createChannel(tx, channelName, tableName)
		})
		return err

	})
}

// its not safe query ada kemungkinan injection jadi hati hati
func (cdc *CapturePostgres) BackfillWithQuery(channel string, raw string) *CapturePostgres {
	if cdc.err != nil {
		return cdc
	}

	backfillScript := `
DO $$
DECLARE
    emp RECORD;
BEGIN
    FOR emp IN %s LOOP
        PERFORM pg_notify('%s',
                      json_build_object(
                          'operation', 'BACKFILL',
                          'id', emp.id,
                          'data', emp

                      )::text);
    END LOOP;
END $$;
	`

	return cdc.dbops(func(tx *gorm.DB) error {
		script := fmt.Sprintf(backfillScript, raw, channel)
		return tx.Exec(script).Error
	})
}

func (cdc *CapturePostgres) Err() error {
	return cdc.err
}

func (cdc *CapturePostgres) dbops(handler func(tx *gorm.DB) error) *CapturePostgres {
	var dialect gorm.Dialector
	switch cdc.driverType {
	case DBTypePostgres:
		dialect = postgres.Open(cdc.dsn)
	case DBTypeGcpCloudPostgres:
		dialect = postgres.New(postgres.Config{
			DriverName: "cloudsqlpostgres",
			DSN:        cdc.dsn,
		})
		// postgres.New(
	}
	db, err := gorm.Open(dialect, &gorm.Config{})
	if err != nil {
		return cdc.setErr(err)
	}

	defer func() {}()
	sql, err := db.DB()
	if err != nil {
		return cdc.setErr(err)
	}

	defer sql.Close()

	err = db.Transaction(handler)

	if err != nil {
		return cdc.setErr(err)
	}

	return cdc
}

func (cdc *CapturePostgres) setErr(err error) *CapturePostgres {
	if cdc.err != nil {
		return cdc
	}
	if err != nil {
		cdc.err = err

	}
	return cdc
}

type FuncItem struct {
	RoutineName string
	RoutineType string
	DataType    string
}

func (cdc *CapturePostgres) ShowFunction(db *gorm.DB) *CapturePostgres {
	if cdc.err != nil {
		return cdc
	}
	hasil := []*FuncItem{}
	err := db.Raw(
		`
SELECT routine_name, routine_type, data_type 
FROM information_schema.routines 
WHERE routine_schema NOT IN ('pg_catalog', 'information_schema');
		`,
	).Scan(&hasil).Error

	if err != nil {
		return cdc.setErr(err)
	}

	for _, item := range hasil {
		log.Println(item.RoutineName, item.RoutineType, item.DataType)
	}

	return cdc
}

func (cdc *CapturePostgres) createChannel(db *gorm.DB, channel string, tablename string) error {
	log.Printf("creating function %s", channel)
	cfunc := `
CREATE OR REPLACE FUNCTION notify_%s()
RETURNS TRIGGER AS $$
BEGIN
    PERFORM pg_notify('%s', 
                      json_build_object(
                          'operation', TG_OP,
                          'id', NEW.id,
                          'data', NEW
                          
                      )::text);
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;
	`
	cfunc = fmt.Sprintf(cfunc, channel, channel)
	err := db.Exec(cfunc).Error
	if err != nil {
		return err
	}
	log.Printf("creating trigger %s", channel)
	dropfirst := fmt.Sprintf(`DROP TRIGGER IF EXISTS %s_trigger ON %s;`, channel, tablename)
	err = db.Exec(dropfirst).Error
	if err != nil {
		return err
	}

	ctrig := `
CREATE TRIGGER %s_trigger
AFTER INSERT OR UPDATE OR DELETE ON %s
FOR EACH ROW EXECUTE FUNCTION notify_%s();
	`
	ctrig = fmt.Sprintf(ctrig, channel, tablename, channel)
	err = db.Exec(ctrig).Error
	if err != nil {
		return err
	}

	return err
}

type dialer struct{}

var instanceRegexp = regexp.MustCompile(`^\[(.+)\]:[0-9]+$`)

func (d dialer) Dial(ntw, addr string) (net.Conn, error) {
	matches := instanceRegexp.FindStringSubmatch(addr)
	if len(matches) != 2 {
		return nil, fmt.Errorf("failed to parse addr: %q. It should conform to the regular expression %q", addr, instanceRegexp)
	}
	instance := matches[1]
	return proxy.Dial(instance)
}
func (d dialer) DialTimeout(ntw, addr string, timeout time.Duration) (net.Conn, error) {
	return nil, fmt.Errorf("timeout is not currently supported for cloudsqlpostgres dialer")
}

// performing backfilling
// DO $$
// DECLARE
//     emp RECORD;
// BEGIN
//     FOR emp IN SELECT * FROM employees LOOP
//         PERFORM pg_notify('testchan',
//                       json_build_object(
//                           'operation', 'BACKFILL',
//                           'id', emp.id,
//                           'data', emp

//                       )::text);
//     END LOOP;
// END $$;
