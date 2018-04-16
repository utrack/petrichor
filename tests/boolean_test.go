package tests

import (
	"testing"

	p "github.com/utrack/petrichor/client/confc"
)

func TestBooleanValue(t *testing.T) {
	A := p.NewBooleanV("A", p.SettingDesc{}, true, p.NewDelayedRegisterer(), p.NewProxyValuer())
	if !A.Value() {
		t.Error("Expected true, got false")
	}
}

func TestBooleanUpdateViaChan(t *testing.T) {
	pv := p.NewProxyValuer()
	A := p.NewBooleanV("A", p.SettingDesc{}, true, p.NewDelayedRegisterer(), pv)

	uc := A.Updates()
	A.UpdateChan() <- "false"
	v := <-uc
	if v != false {
		t.Error("Expected false from chan, got true")
	}

	if  A.Value() != false {
		t.Error("Expected false from Value, got true")
	}
}