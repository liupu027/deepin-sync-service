package sync

import (
	"os"
)

func getRunDir() string {
	run := os.Getenv("XDG_RUNTIME_DIR")
	if "" == run {
		run = os.TempDir()
	}
	return run
}
