package setetcd

import (
	"github.com/pkg/errors"
	"github.com/utrack/petrichor/client/confc"
)

type Storage interface {
	Client
}

type Client interface {
	// ExportDefinitions pushes the TypedSettingDescs to etcd.
	// It tries to grab the update mtx in process, returning ErrClientHasNoMutex
	// if there's any other instance holding the mutex.
	//
	// Nevertheless, an update process starts that tries to re-grab the mtx
	// and write the definition just in case if the mtx holder goes down.
	// If mtx grab succeeds - it rewrites the definition.
	ExportDefinitions(map[string]confc.TypedSettingDesc) error

	// RegisterVersionAndMigrate sets current version as latest and migrates
	// settings of old version to this one (if it wasn't done before).
	//
	// Can return ErrClientHasNoMutex.
	RegisterVersionAndMigrate() error

	Values() (map[string]string, error)

	// Changes returns a settings' change stream.
	Changes() (<-chan ChangeNotification, error)
}

// ChangeNotification is a notification about a change of a setting's value.
// Does not include any info about type/tag/module changes.
type ChangeNotification struct {
	Name  string
	Value string
}

var (
	ErrClientHasNoMutex = errors.New("this Client has no exclusive mutex")
)
