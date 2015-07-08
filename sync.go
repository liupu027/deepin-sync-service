package sync

import (
	"errors"
	"pkg.deepin.io/lib/log"
	"time"
)

var Log = log.NewLogger("deepin-sync")

var (
	ErrNullSyncEntry = errors.New("Invaild Sync Entry Name")
)

type Sync struct {
	entryMap map[string](*syncEntry)

	// however, a push server is more useful than query
	syncTimer *time.Timer

	Failed func(id string)
	Finish func(id string)
}

var _sync *Sync

var syncInterval = 20 * time.Second

func getSync() *Sync {
	if nil == _sync {
		_sync = &Sync{
			entryMap:  map[string](*syncEntry){},
			syncTimer: time.NewTimer(syncInterval),
		}
		go _sync.startSync()
	}
	return _sync
}

func (sy *Sync) startSync() {
	for {
		select {
		case <-sy.syncTimer.C:
			sy.SyncAll()
			sy.syncTimer.Reset(30 * time.Second)
		}
	}
}

func Register(e Entry) error {
	//create cache namespace for entry
	c, err := newCache(e)
	if nil != err {
		Log.Error(err)
		return err
	}

	//create remote Entry
	re, err := newRemoteEntry(e)
	if nil != err {
		Log.Error(err)
		return err
	}
	se := &syncEntry{
		remote: re,
		entry:  e,
		cache:  c,
	}
	Log.Info(e.ID(), se)
	getSync().entryMap[e.ID()] = se
	Log.Info(getSync().entryMap)
	return nil
}

func CommitEntry(id string) error {
	return getSync().CommitEntry(id)
}

func (sy *Sync) CommitEntry(id string) error {
	entry := sy.entryMap[id]
	Log.Info(sy.entryMap, id, entry)
	if nil == entry {
		Log.Error(entry, ErrNullSyncEntry)
		return ErrNullSyncEntry
	}
	return entry.commit()
}

func (sy *Sync) SyncEntry(id string) error {
	entry := sy.entryMap[id]
	Log.Info(sy.entryMap, id, entry)
	if nil == entry {
		Log.Error(entry, ErrNullSyncEntry)
		return ErrNullSyncEntry
	}
	return entry.sync()
}

func (sy *Sync) SyncAll() error {
	for k, v := range sy.entryMap {
		if err := v.sync(); nil != err {
			Log.Warningf("Sync %v Failed: %v", k, err)
		}
	}
	return nil
}
