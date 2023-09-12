package cbreak

import (
	"context"
	"errors"
	"testing"
	"time"
)

type M struct {
	b     bool
	count int
}

func (m *M) Func() string {
	m.count++
	if m.count%2 == 0 {
		if m.count > 21 {
			if m.b {
				return "not cool"
			}
			return "cool"
		}
		return "not cool"
	}
	return "cool"
}
func ChangeState(s int) {
}
func TestExecuteOpenToHalf(t *testing.T) {
	t.Parallel()
	cb := New(ChangeState)
	m := M{}
	for i := 0; i < 21; i++ {
		time.Sleep(1 * time.Millisecond)
		cb.Execute(context.Background(), func() (interface{}, error) {
			l := m.Func()

			if l != "cool" {
				return "", errors.New("WTF")
			}
			return "", nil
		})
	}
	want := Open
	got := cb.ReturnState()
	if want != got {
		t.Errorf("Want %d, but got %d", want, got)
	}
	time.Sleep(1 * time.Second)
	// The state should have changed to halfState
	want = Half
	got = cb.ReturnState()
	if want != got {
		t.Errorf("Want %d, but got %d", want, got)
	}

}

func TestExecuteHalfToClosed(t *testing.T) {
	t.Parallel()
	cb := New(ChangeState)
	m := M{}
	for i := 0; i < 21; i++ {
		time.Sleep(1 * time.Millisecond)
		cb.Execute(context.Background(), func() (interface{}, error) {
			l := m.Func()

			if l != "cool" {
				return "", errors.New("WTF")
			}
			return "", nil
		})
	}
	// changes state to half state
	time.Sleep(2 * time.Second)
	for i := 0; i < 5; i++ {
		cb.Execute(context.Background(), func() (interface{}, error) {
			l := m.Func()

			if l != "cool" {
				return "", errors.New("WTF")
			}
			return "", nil
		})
	}
	want := Closed
	got := cb.ReturnState()
	if want != got {
		t.Errorf("Want %d, but got %d", want, got)
	}

}

func TestExecuteHalfToOpen(t *testing.T) {
	t.Parallel()
	cb := New(ChangeState)
	m := M{true, 0}
	for i := 0; i < 21; i++ {
		time.Sleep(10 * time.Millisecond)
		cb.Execute(context.Background(), func() (interface{}, error) {
			l := m.Func()

			if l != "cool" {
				return "", errors.New("WTF")
			}
			return "", nil
		})
	}
	// changes state to half state
	time.Sleep(2 * time.Second)

	for i := 0; i < 5; i++ {
		cb.Execute(context.Background(), func() (interface{}, error) {

			l := m.Func()
			if i == 3 {
				return "", errors.New("Some error")
			}
			if l != "cool" {
				return "", errors.New("WTF")
			}
			return "", nil
		})
	}
	want := Open
	got := cb.ReturnState()
	if want != got {
		t.Errorf("Want %v, but got %v", want, got)
	}

}
