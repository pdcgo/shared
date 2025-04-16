package gorm_commenter

import (
	"fmt"
	"strings"

	gorm "gorm.io/gorm"
	gormclause "gorm.io/gorm/clause"
)

type Comment struct {
	Content string
}

func (c Comment) Name() string {
	return "COMMENT"
}

func (c Comment) Build(builder gormclause.Builder) {
	builder.WriteString("/* ")
	builder.WriteString(c.Content)
	builder.WriteString(" */")
}

func (c Comment) MergeClause(mergeClause *gormclause.Clause) {
}

func (c Comment) ModifyStatement(stmt *gorm.Statement) {
	clause := stmt.Clauses[c.Name()]
	// Perhatikan bahwa ini harus berupa Ekspresi, karena jika Ekspresi nihil, eksekusi metode Build tidak akan terpicu.
	// Di sini kami mengacu pada BeforeExpression yang didaftarkan oleh petunjuk di awal, yang menyebabkan Build tidak dapat dieksekusi, baru setelah kami menyelesaikan seluruh proses gorm kami menemukan masalahnya.
	clause.Expression = c
	stmt.Clauses[c.Name()] = clause
}

var extraClause = []string{"COMMENT"}

type CommentClausePlugin struct{}

// NewCommentClausePlugin create a new ExtraPlugin
func NewCommentClausePlugin() *CommentClausePlugin {
	return &CommentClausePlugin{}
}

// Name plugin name
func (ep *CommentClausePlugin) Name() string {
	return "CommentClausePlugin"
}

// Initialize register BuildClauses
func (ep *CommentClausePlugin) Initialize(db *gorm.DB) (err error) {
	initClauses(db)
	db.Callback().Create().After("gorm:create").Register("CommentClausePlugin", AddAnnotation)
	db.Callback().Delete().After("gorm:delete").Register("CommentClausePlugin", AddAnnotation)
	db.Callback().Query().After("gorm:query").Register("CommentClausePlugin", AddAnnotation)
	db.Callback().Update().After("gorm:update").Register("CommentClausePlugin", AddAnnotation)
	db.Callback().Raw().After("gorm:raw").Register("CommentClausePlugin", AddAnnotation)
	db.Callback().Row().After("gorm:row").Register("CommentClausePlugin", AddAnnotation)

	return
}

type TraceData map[string]string

func GetTraceData(db *gorm.DB) string {
	traceData := TraceData{
		"action":      "",
		"controller":  "",
		"framework":   "",
		"route":       "",
		"application": "",
		"db driver":   "",
	}

	for key := range traceData {
		if v, ok := db.Statement.Context.Value(key).(string); ok {
			traceData[key] = v
		}
	}

	content := []string{}
	for key, value := range traceData {
		if value == "" {
			continue
		}
		content = append(content, fmt.Sprintf("%s=%s", key, value))
	}

	return strings.Join(content, ",\n")
}

func AddAnnotation(db *gorm.DB) {
	if db.Error != nil {
		return
	}

	content := GetTraceData(db)
	if db.Statement.SQL.Len() > 0 {
		oldSQL := db.Statement.SQL.String()
		db.Statement.SQL.Reset()
		db.Statement.SQL.WriteString(fmt.Sprintf("%s /* %s */", oldSQL, content))
		return
	}

	db.Statement.AddClause(Comment{Content: content})
}

// initClauses init SQL clause
func initClauses(db *gorm.DB) {
	if db.Error != nil {
		return
	}
	createClause := append(extraClause, db.Callback().Create().Clauses...)
	deleteClause := append(extraClause, db.Callback().Delete().Clauses...)
	queryClause := append(extraClause, db.Callback().Query().Clauses...)
	updateClause := append(extraClause, db.Callback().Update().Clauses...)
	rawClause := append(extraClause, db.Callback().Raw().Clauses...)
	rowClause := append(extraClause, db.Callback().Row().Clauses...)
	db.Callback().Create().Clauses = createClause
	db.Callback().Delete().Clauses = deleteClause
	db.Statement.Callback().Query().Clauses = queryClause
	db.Callback().Update().Clauses = updateClause
	db.Callback().Raw().Clauses = rawClause
	db.Callback().Row().Clauses = rowClause
}

// main.go
// package main

// import (
//     "context"
//     "plugins"

//     gorm "gorm.io/gorm"
//     "github.com/google/uuid"
// )

// type Product struct {
//   gorm.Model
//   Code  string
//   Price uint
// }

// func main() {
//     dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
//     db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

//     db.Use(plugins.NewCommentClausePlugin())

//     db.Create(&Product{Code: "D42", Price: 100})

//     // 传入context，指定rid
//     ctx := context.WithValue(context.Background(), "rid", uuid.New().String())
//     db.WithContext(ctx).Create(&Product{Code: "D42", Price: 100})
// }
