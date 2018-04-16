package etcd

import (
	"sync"

	"github.com/pkg/errors"
	"github.com/utrack/petrichor/settings/setetcd"
)

type Valuer struct {
	cli setetcd.Client

	valMtx sync.RWMutex
	valMap map[string]string

	subMtx sync.RWMutex
	subs   map[string][]chan string
}

func NewValuer(cli etcd.Client) (*Valuer, error) {
	ret := &Valuer{
		cli:    cli,
		valMap: map[string]string{},
		subs:   map[string][]chan string{},
	}
	err := ret.initValues()
	if err != nil {
		return nil, errors.Wrap(err, "couldn't get seed values from the client")
	}

	// TODO run refresher pump
	// TODO close all sub chans if seed chan from client closes
	return ret, nil
}

// createSubChan creates and registers a subscription chan for single value
// should be called after initValues
func (v *Valuer) createSubChan(vname string) <-chan string {
	ret := make(chan string, 10)
	v.valMtx.RLock()
	defer v.valMtx.RUnlock()
	// TODO retrieve default from the registerer if not found
	ret <- v.valMap[vname]

	v.subMtx.Lock()
	defer v.subMtx.Unlock()
	v.subs[vname] = append(v.subs[vname], ret)
	return ret
}

// get seed values from the client
// prevent zero-values from being propagated to the listeners first
func (v *Valuer) initValues() error {
	v.valMtx.Lock()
	defer v.valMtx.Unlock()

	var err error
	v.valMap, err = v.cli.Values()
	return errors.Wrap(err, "got an error from etcd client")
}

func (v *Valuer) setValue(name, value string) {
	v.valMtx.Lock()
	defer v.valMtx.Unlock()
	v.valMap[name] = value
}

func (v *Valuer) outputValue(name, value string) {
	v.subMtx.Lock()
	defer v.subMtx.Unlock()

	for _, s := range v.subs[name] {
		select {
		case s <- value:
		default:
		}
	}
}
