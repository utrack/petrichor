package confc

import (
	"strconv"
)

const (
	updateBufferCap = 5
)

// Valuer provides value updates from a provider.
type Valuer interface {
	ChanForValue(name string) <-chan string
	Value(string) string
}

// NewBoolean creates new Boolean value using default Registerer and Valuer.
func NewBoolean(name string, desc SettingDesc, def bool) *Boolean {
	return NewBooleanV(name, desc, def, DefaultValuer)
}

// NewBooleanV creates new Boolean value using default Registerer and custom Valuer.
func NewBooleanV(name string, desc SettingDesc, def bool, v Valuer) *Boolean {
	info := newTypedDesc(name, def, TypeBoolean, desc)
	DefaultRegisterer.MustRegister(info)
	return initBooleanValue(def, v.ChanForValue(name))
}

// NewBooleanG returns a getter func that retrieves a boolean value from a given
// Valuer.
//
// Returns last successfully retrieved value if Valuer returns garbage.
func NewBooleanG(name string, desc SettingDesc, def bool) func(Valuer) bool {
	info := newTypedDesc(name, def, TypeBoolean, desc)
	DefaultRegisterer.MustRegister(info)

	value := def
	return func(v Valuer) bool {
		newVal, err := parseBoolValue(v.Value(name))
		if err == nil {
			value = newVal
		}
		return value
	}
}

func initBooleanValue(def bool, v <-chan string) *Boolean {
	ret := &Boolean{
		v:          def,
		outUpdates: []chan bool{},
	}
	go boolPump(ret, boolUpdateChan(v))
	return ret
}

// Boolean provides a dynamic Boolean value.
type Boolean struct {
	v bool

	// TODO make it thread safe
	outUpdates []chan bool
}

// Value returns current value.
func (b *Boolean) Value() bool {
	return b.v
}

// Updates creates a new channel pumping out this setting's updates.
func (b *Boolean) Updates() <-chan bool {
	ret := make(chan bool, updateBufferCap)
	b.outUpdates = append(b.outUpdates, ret)
	return ret
}

func boolPump(b *Boolean, updates <-chan bool) {
	for tempVal := range updates {
		v := tempVal
		// TODO thread safe mu
		b.v = v
		for _, ch := range b.outUpdates {
			ch <- v
		}
	}

	for _, ch := range b.outUpdates {
		close(ch)
	}
}

func boolUpdateChan(v <-chan string) <-chan bool {
	ret := make(chan bool, updateBufferCap)
	go func() {
		defer close(ret)
		for val := range v {
			newVal, err := parseBoolValue(val)
			if err != nil {
				// TODO log debug?
				continue
			}
			ret <- newVal
		}
	}()
	return ret
}

func parseBoolValue(v string) (bool, error) {
	return strconv.ParseBool(v)
}
