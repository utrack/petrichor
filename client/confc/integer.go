package confc

import (
	"strconv"
)

// NewInteger creates new Integer value using default Registerer and Valuer.
func NewInteger(name string, desc SettingDesc, def int64) *Integer {
	return NewIntegerV(name, desc, def, defaultRegisterer, defaultValuer)
}

// NewIntegerV creates new Integer value using default Registerer and custom Valuer.
func NewIntegerV(name string, desc SettingDesc, def int64, r Registerer, v Valuer) *Integer {
	info := newTypedDesc(name, def, TypeInteger, desc)
	r.MustRegister(info)
	return initIntegerValue(def, v.ChanForValue(name))
}

func initIntegerValue(def int64, v chan string) *Integer {
	ret := &Integer{
		v:          def,
		updateChan: v,
		outUpdates: []chan int64{},
	}
	go intPump(ret, intUpdateChan(v))
	return ret
}

// Integer provides a dynamic Integer value.
type Integer struct {
	v int64
	updateChan chan string

	// TODO make it thread safe
	outUpdates []chan int64
}

// Value returns current value.
func (b *Integer) Value() int64 {
	return b.v
}

func (b *Integer) UpdateChan() chan<- string {
	return b.updateChan
}

// Updates creates a new channel pumping out this setting's updates.
func (b *Integer) Updates() <-chan int64 {
	ret := make(chan int64, updateBufferCap)
	b.outUpdates = append(b.outUpdates, ret)
	return ret
}

func intPump(b *Integer, updates <-chan int64) {
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

func intUpdateChan(v <-chan string) <-chan int64 {
	ret := make(chan int64, updateBufferCap)
	go func() {
		defer close(ret)
		for val := range v {
			newVal, err := parseIntValue(val)
			if err != nil {
				// TODO log debug?
				continue
			}
			ret <- newVal
		}
	}()
	return ret
}

func parseIntValue(v string) (int64, error) {
	return strconv.ParseInt(v, 10, 64)
}
