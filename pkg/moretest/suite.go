package moretest

import (
	"errors"
	"os/exec"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

type SetupFunc func(t *testing.T) func() error

type SetupListFunc []SetupFunc

func Suite(t *testing.T, name string, setups SetupListFunc, handler func(t *testing.T)) {
	tearDown := []func() error{}
	for _, setup := range setups {
		tearfunc := setup(t)
		tearDown = append(tearDown, tearfunc)
	}
	t.Run(name, handler)

	var err error

	for tlen := len(tearDown) - 1; tlen >= 0; tlen-- {
		tear := tearDown[tlen]
		if tear == nil {
			continue
		}
		err = tear()
		assert.Nil(t, err)
	}
}

func CheckGCPAuth() error {
	out, err := exec.Command("gcloud", "auth", "list", "--format=value(account)").Output()
	if err != nil {
		return err
	}
	if string(out) == "" {
		return errors.New("account gcloud not setup")
	}
	return nil
}

var isGcpLogin error
var isGcpLoginOnce sync.Once

func SkipGcpNotLogin(t *testing.T, handler func(t *testing.T)) func(t *testing.T) {

	isGcpLoginOnce.Do(func() {
		isGcpLogin = CheckGCPAuth()

	})

	return func(t *testing.T) {
		if isGcpLogin != nil {
			t.Skip("gcp error " + isGcpLogin.Error())
		}

		handler(t)
	}

}
