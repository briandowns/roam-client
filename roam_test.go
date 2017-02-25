package roam

import (
	"reflect"
	"testing"
)

// TestNew
func TestNew(t *testing.T) {
	c := New("asdf1234", nil)
	if reflect.TypeOf(c).String() != "*goroamclient.Client" {
		t.Error("incorrect data type returned")
	}
}

// TestSetPoint
func TestSetPoint(t *testing.T) {
	c := New("asdf1234", nil)
	if _, err := c.SetPoint(&Point{}); err != nil {
		t.Error(err)
	}
}

// TestSetHook
func TestSetHook(t *testing.T) {
	c := New("asdf1234", nil)
	if _, err := c.SetHook(&Hook{}); err != nil {
		t.Error(err)
	}
}

// TestDeleteHook
func TestDeleteHook(t *testing.T) {
	c := New("asdf1234", nil)
	if _, err := c.DeleteHook("hook"); err != nil {
		t.Error(err)
	}
}

// TestHooks
func TestHooks(t *testing.T) {
	c := New("asdf1234", nil)
	_, err := c.Hooks()
	if err != nil {
		t.Error(err)
	}
}

type testRoam struct{}

// call
func (t *testRoam) call(m, url string, result interface{}) error {
	return nil
}
