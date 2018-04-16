package tests

import
(
	"testing"

	p "github.com/utrack/petrichor/client/confc"
)

func TestIntegerValue(t *testing.T) {
	A := p.NewIntegerV("A", p.SettingDesc{}, 42, p.NewDelayedRegisterer(), p.NewProxyValuer())
	value := A.Value()
	if value != 42 {
		t.Error("Expected 42, got ", value)
	}
}

func TestIntegerUpdateViaChan(t *testing.T) {
	pv := p.NewProxyValuer()
	A := p.NewIntegerV("A", p.SettingDesc{}, 42, p.NewDelayedRegisterer(), pv)

	uc := A.Updates()
	A.UpdateChan() <- "24"
	v := <-uc
	if v != 24 {
		t.Error("Expected 24 from chan, got ", v)
	}

	if  A.Value() != 24 {
		t.Error("Expected 24 from Value, got ", v)
	}
}