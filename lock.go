package sync

import (
	"encoding/json"
	"errors"
	"pkg.deepin.io/lib/utils"
	"time"
)

var (
	LockVerfiyFailed = errors.New("Verfiy Lock Failed, Maybe lock by Other Client")
)

type lock struct {
	Id       string
	CreateAt time.Time
}

func newLock() *lock {
	return &lock{
		Id: utils.GenUuid(),
	}
}

func (l *lock) Name() string {
	return "lock." + l.Id
}

func (l *lock) check(id string) error {
	if id != l.Id {
		return LockVerfiyFailed
	}
	return nil
}

func (l *lock) json() []byte {
	j, err := json.Marshal(l)
	if nil != err {
		Log.Error(err)
	}
	return j
}

func (l *lock) verfiy(lockjson []byte) error {
	var verfiyLock lock
	err := json.Unmarshal(lockjson, &verfiyLock)
	if nil != err {
		return err
	}

	if l.Id != verfiyLock.Id {
		return LockVerfiyFailed
	}

	return nil
}
