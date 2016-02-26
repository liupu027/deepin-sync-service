package main

import (
	"fmt"
	"os"
	"pkg.deepin.io/sync"
	_ "pkg.deepin.io/sync/plugin/theme"
	"pkg.deepin.io/lib/dbus"
)

func main() {
	if err := sync.LoadDBus(); nil != err {
		os.Stderr.Write([]byte(fmt.Sprint(err)))
		os.Exit(1)
	}

	dbus.DealWithUnhandledMessage()

	if err := dbus.Wait(); nil != err {
		os.Stderr.Write([]byte(fmt.Sprint(err)))
		os.Exit(1)
	}
	os.Stdout.Write([]byte("daemon Exit"))
	os.Exit(0)
}
