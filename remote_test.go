package sync

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestAll(t *testing.T) {
	re, err := newRemoteEntry("theme")
	Convey("Test Remote Entry", t, func() {
		So(err, ShouldBeNil)
		ConveyRemoteLock(re)
		ConveyRemoteUnlock(re)
	})

}

func ConveyRemoteLock(re *remoteEntry) {
	Convey("Lock", func() {
		err := re.lock()
		So(err, ShouldBeNil)
	})
}

func ConveyRemoteUnlock(re *remoteEntry) {
	Convey("Unlock", func() {
		err := re.unlock()
		So(err, ShouldBeNil)
	})
}
