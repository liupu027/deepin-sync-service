package sync

import (
	"bytes"
	"encoding/json"
	"os"
	"pkg.deepin.io/storage"
	"pkg.deepin.io/storage/backend/qihoo"
	"strings"
	"time"
)

type remoteEntry struct {
	entry Entry
	s     storage.Storage
	l     *lock
}

func newRemoteEntry(e Entry) (*remoteEntry, error) {
	qs, err := qihoo.NewQihooStorage()
	if nil != err {
		return nil, err
	}
	re := &remoteEntry{
		entry: e,
		s:     qs,
		l:     newLock(),
	}
	os.Mkdir(re.sapce(), 0755)
	//todo make namespace
	return re, nil
}

func (re *remoteEntry) unlock() error {
	tmpRemoteLockPath := re.sapce() + re.l.Name()
	newRemoteLockName := "lock.head." + re.l.Id
	newRemoteLockPath := re.sapce() + newRemoteLockName
	Log.Infof("Remove Lock: %v, %v", tmpRemoteLockPath, newRemoteLockPath)
	return re.s.Del([]string{tmpRemoteLockPath, newRemoteLockPath})
}

func (re *remoteEntry) lock() error {
	//create tmplock
	tmpLockPath := getRunDir() + re.sapce() + re.l.Name()
	tmpLockFile, err := os.Create(tmpLockPath)
	if nil != err {
		return err
	}
	re.l.CreateAt = time.Now()
	tmpLockFile.Write(re.l.json())
	tmpLockFile.Close()

	tmpRemoteLockPath := re.sapce() + re.l.Name()
	err = re.s.PutFile(tmpLockPath, tmpRemoteLockPath)
	if nil != err {
		Log.Error(err)
		return err
	}

	list, _ := re.s.List(re.sapce())
	var tmpLockNode storage.Status
	for _, v := range list {
		if v.FullPath == tmpRemoteLockPath {
			tmpLockNode = v
		}
	}

	LockHead := "lock.head."
	for _, v := range list {
		if strings.Contains(v.Name, LockHead) {
			//check id
			if err := re.l.check(strings.Replace(v.Name, LockHead, "", -1)); nil != err {
				//check time expired
				Log.Info(tmpLockNode, v)
				if tmpLockNode.RemoteModfiyTime.Sub(v.RemoteModfiyTime) < (10 * time.Minute) {
					Log.Error(LockVerfiyFailed)
					return LockVerfiyFailed
				}
				Log.Infof("Remote Lock %v Expire, LastModifyAt %v", v.Name, v.RemoteModfiyTime)
			}
		}
	}

	//move lock
	newRemoteLockName := "lock.head." + re.l.Id
	newRemoteLockPath := re.sapce() + newRemoteLockName
	if err := re.s.Rename(tmpRemoteLockPath, newRemoteLockName); nil != err {
		Log.Error(err)
		return err
	}

	//recheck locker
	var buf []byte
	rd := bytes.NewBuffer(buf)
	if err := re.s.Get(rd, newRemoteLockPath); nil != err {
		Log.Error(err)
		return err
	}

	return re.l.verfiy(rd.Bytes())
}

func (re *remoteEntry) sapce() string {
	return BaseNameSpace + re.entry.ID() + "/"
}

func (re *remoteEntry) modifyTime() (time.Time, error) {
	head := re.sapce() + "head"
	st, err := re.s.Stat(head)
	return st.RemoteModfiyTime, err
}

func (re *remoteEntry) push(cache *cache) error {
	Log.Warning(cache.locleSpace()+"/head", re.sapce()+"head")
	return re.s.PutFile(cache.locleSpace()+"/head", re.sapce()+"head")
}

func (re *remoteEntry) pull() (*cache, error) {
	var buf []byte
	rd := bytes.NewBuffer(buf)
	head := re.sapce() + "head"
	if err := re.s.Get(rd, head); nil != err {
		return nil, err
	}
	var c cache
	Log.Info(string(rd.Bytes()))
	if err := json.Unmarshal(rd.Bytes(), &c); nil != err {
		return nil, err
	}
	return &c, nil
}

func (re *remoteEntry) sign() (*signature, error) {
	head := re.sapce() + "head"
	st, err := re.s.Stat(head)
	if nil != err {
		return nil, err
	}
	sign := signature{
		ModifyTime: st.RemoteModfiyTime,
		Sha1:       st.Sha1,
	}
	return &sign, nil
}
