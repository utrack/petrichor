package confc

import (
	"strconv"
)

// NewBoolean creates new Boolean value using default Registerer and Valuer.
func NewBoolean(name string, desc SettingDesc, def bool) *Boolean {
	return NewBooleanV(name, desc, def, defaultRegisterer, defaultValuer)
}

// NewBooleanV creates new Boolean value using default Registerer and custom Valuer.
func NewBooleanV(name string, desc SettingDesc, def bool, r Registerer, v Valuer) *Boolean {
	info := newTypedDesc(name, def, TypeBoolean, desc)
	r.MustRegister(info)
	return initBooleanValue(def, v.ChanForValue(name))
}

func initBooleanValue(def bool, v chan string) *Boolean {
	ret := &Boolean{
		v:          def,
		updateChan: v,
		outUpdates: []chan bool{},
	}
	go boolPump(ret, boolUpdateChan(v))
	return ret
}

// Boolean provides a dynamic Boolean value.
type Boolean struct {
	v bool
	updateChan chan string

	// TODO make it thread safe
	outUpdates []chan bool
}

// Value returns current value.
func (b *Boolean) Value() bool {
	return b.v
}

func (b *Boolean) UpdateChan() chan<- string {
	return b.updateChan
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
