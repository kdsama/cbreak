package cbreak

import (
	"context"
	"errors"
	"testing"
	"time"
)

func Func() string {
	return "Cool"
}
func TestExecute(t *testing.T) {
	cb := New()
	for i := 0; i < 100; i++ {
		time.Sleep(1 * time.Second)
		_, err := cb.Execute(context.Background(), func() (interface{}, error) {
			l := Func()
			if l != "cool" {
				return "", errors.New("WTF")
			}
			return "", nil
		})
		if err != nil {
			t.Error(err)
		}
	}

}
