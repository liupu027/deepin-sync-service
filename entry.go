package sync

var (
	BaseNameSpace = "/data/deepin/sync/entrys/"
)

type Entry interface {
	ID() string
	Snapshot() ([]byte, error)
	Restore([]byte) error
}
