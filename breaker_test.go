package cbreak

import (
	"context"
	"testing"
)

func Func() string {
	return "Cool"
}
func TestExecute(t *testing.T) {
	cb := New()
	_, err := cb.Execute(context.Background(), func() (interface{}, error) {

		return "", nil
	})
	if err != nil {
		t.Error("Error")
	}
}
