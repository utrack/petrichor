package etcd

import (
	"sync"

	"github.com/pkg/errors"
	"github.com/utrack/petrichor/client/confc"
	"github.com/utrack/petrichor/settings/setetcd"
)

type Registerer struct {
	cli setetcd.Client

	settMu sync.Mutex
	setts  map[string]confc.TypedSettingDesc
}

func NewRegisterer(cli setetcd.Client) *Registerer {
	return &Registerer{
		cli:   cli,
		setts: map[string]confc.TypedSettingDesc{},
	}
}

func (r *Registerer) Register(d confc.TypedSettingDesc) error {
	r.settMu.Lock()
	defer r.settMu.Unlock()
	if _, ok := r.setts[d.Name]; ok {
		return errors.Errorf("setting %v is already registered", d.Name)
	}
	r.setts[d.Name] = d
	err := r.cli.ExportDefinitions(r.setts)
	if err == setetcd.ErrClientHasNoMutex {
		// TODO log debug?
		err = nil
	}
	return err
}

func (r *Registerer) MustRegister(d confc.TypedSettingDesc) {
	if err := r.Register(d); err != nil {
		panic(err)
	}
}

func (r *Registerer) Desc(name string) *confc.TypedSettingDesc {
	r.settMu.Lock()
	defer r.settMu.Unlock()

	ret, ok := r.setts[name]
	if !ok {
		return nil
	}
	return &ret
}
