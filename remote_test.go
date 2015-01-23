package sync

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

type testEntry struct {
}

func (*testEntry) ID() string {
	return "test"
}
func (e *testEntry) Snapshot() ([]byte, error) {
	return []byte("testValue"), nil
}

func (e *testEntry) Restore(data []byte) error {
	return nil
}

func TestAll(t *testing.T) {
	var te testEntry
	re, err := newRemoteEntry(&te)
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
