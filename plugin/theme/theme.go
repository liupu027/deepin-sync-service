package theme

import (
	themeManagerAPI "dbus/com/deepin/daemon/thememanager"
	"fmt"
	"pkg.deepin.io/sync"
	"time"
)

type themeEntry struct {
	monitor  *time.Timer
	oldValue string
}

var _theme themeEntry

func init() {
	sync.Register(&_theme)
	go _theme.Mon()
}

const (
	MANAGER_DEST = "com.deepin.daemon.ThemeManager"
	MANAGER_PATH = "/com/deepin/daemon/ThemeManager"
)

func (e *themeEntry) Mon() {
	_themeManger, _ := themeManagerAPI.NewThemeManager(MANAGER_DEST, MANAGER_PATH)
	e.monitor = time.NewTimer(10 * time.Second)
	for {
		select {
		case <-e.monitor.C:
			fmt.Print(e)
			if "" == e.oldValue {
				e.oldValue = _themeManger.CurrentTheme.Get()
			}

			if e.oldValue != _themeManger.CurrentTheme.Get() {
				e.oldValue = _themeManger.CurrentTheme.Get()
				sync.CommitEntry(e.ID())
			}
			e.monitor.Reset(10 * time.Second)
		}
	}
}

func (e *themeEntry) ID() string {
	return "theme"
}

func (e *themeEntry) Snapshot() ([]byte, error) {
	return []byte(e.oldValue), nil
}

func (e *themeEntry) Restore(data []byte) error {
	_themeManger, _ := themeManagerAPI.NewThemeManager(MANAGER_DEST, MANAGER_PATH)
	_themeManger.CurrentTheme.Set(string(data))
	return nil
}
