package confc

import (
	"errors"
	"sync"
)

// Registerer publishes settings' info.
type Registerer interface {
	Register(TypedSettingDesc) error
	MustRegister(TypedSettingDesc)

	Desc(name string) *TypedSettingDesc

	// TODO add module desc mapping id -> name etc, see design
}

// defaultRegisterer is a default registry point for the created values.
var defaultRegisterer Registerer = NewDelayedRegisterer()

// SetRegisterer sets the default registerer for values.
// All previously registered values are relayed to the new one.
// This function should only be called once.
func SetRegisterer(r Registerer) {
	defaultRegisterer.(*delayedRegisterer).Relay(r)
}

// delayedRegisterer is required to support values registering on init(), but
// the registerer itself appearing at runtimedelayedRegisterer
type delayedRegisterer struct {
	sync.RWMutex

	realR    Registerer
	settings map[string]TypedSettingDesc
}

func NewDelayedRegisterer() *delayedRegisterer {
	return &delayedRegisterer{
		settings: map[string]TypedSettingDesc{},
	}
}

func (r *delayedRegisterer) Register(s TypedSettingDesc) error {
	r.Lock()
	defer r.Unlock()

	if r.realR != nil {
		return r.realR.Register(s)
	}

	if _, ok := r.settings[s.Name]; ok {
		return errors.New("setting of this name already exists")
	}
	r.settings[s.Name] = s

	return nil
}

func (r *delayedRegisterer) MustRegister(s TypedSettingDesc) {
	if err := r.Register(s); err != nil {
		panic(err)
	}
}

func (r *delayedRegisterer) Relay(n Registerer) {
	r.Lock()
	if r.settings == nil || r.realR != nil {
		panic("Relay called twice!")
	}
	r.Unlock()

	for _, s := range r.settings {
		n.MustRegister(s)
	}

	r.settings = nil
	r.realR = n
}

func (r *delayedRegisterer) Desc(name string) *TypedSettingDesc {
	r.RLock()
	defer r.RUnlock()

	if r.realR != nil {
		return r.realR.Desc(name)
	}

	ret, ok := r.settings[name]
	if !ok {
		return nil
	}

	return &ret
}
