package cbreak

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"
)

type M struct {
	count int
}

func (m *M) Func() string {
	m.count++
	if m.count%2 == 0 {
		return "not cool"
	}
	return "cool"
}
func ChangeState(s int) {
}
func TestExecute(t *testing.T) {
	cb := New(ChangeState)
	m := M{}
	for i := 0; i < 21; i++ {
		time.Sleep(300 * time.Millisecond)
		cb.Execute(context.Background(), func() (interface{}, error) {
			l := m.Func()

			if l != "cool" {
				return "", errors.New("WTF")
			}
			return "", nil
		})
	}
	fmt.Println("We here ?")
	want := 2
	got := cb.ReturnState()
	if want != got {
		t.Errorf("Want %d, but got %d", want, got)
	}
	time.Sleep(5 * time.Second)
	// The state should have changed to halfState
	want = 1
	got = cb.ReturnState()
	if want != got {
		t.Errorf("Want %d, but got %d", want, got)
	}

}
