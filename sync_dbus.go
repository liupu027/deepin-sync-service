package sync

import (
	"fmt"
	"pkg.deepin.io/lib"
	"pkg.deepin.io/lib/dbus"
)

const (
	DBusName = "com.deepin.sync.service"
	DBusPath = "/com/deepin/sync/service"
	DBusIfc  = "com.deepin.sync.service"
)

type syncDBus struct {
	Sync
}

func (sy *Sync) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		DBusName,
		DBusPath,
		DBusIfc,
	}
}

func LoadDBus() error {
	Log.Info("Start Sync Service")
	if !lib.UniqueOnSession(DBusName) {
		return fmt.Errorf("There already has an Sync daemon running.")
	}

	err := dbus.InstallOnSession(getSync())
	if err != nil {
		return err
	}

	return nil
}
