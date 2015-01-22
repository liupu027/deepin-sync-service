package sync

import (
	"time"
)

type syncEntry struct {
	entry  Entry
	remote *remoteEntry
	cache  *cache
}

func (se *syncEntry) commit() error {
	//when you commit , you always push locate setting to the server!!!!
	buf, err := se.entry.Snapshot()
	if nil != err {
		Log.Error(err)
		return err
	}
	se.cache.Buf = buf
	sign, err := se.remote.sign()
	if nil != err {
		// That mean you are maybe offine
		Log.Warning("You are maybe offine", err)
	}
	if nil != sign {
		sign.ModifyTime = sign.ModifyTime.Add(10 * time.Second)
		se.cache.Sign = *sign
	} else {
		//TODO: hope the time is not so diff
		se.cache.Sign.ModifyTime = time.Now().AddDate(9999, 99, 99)
	}
	se.cache.save()
	return nil

}

func (se *syncEntry) sync() error {
	Log.Warning("syncEntry sync")
	if err := se.cache.load(); nil != err {
		Log.Error(err)
		return err
	}

	Log.Warning("syncEntry sync")
	sign, err := se.remote.sign()
	if nil != err {
		Log.Error(err)
		return err
	}

	if se.cache.Sign.equal(sign) {
		Log.Warning("Sign Equal")
		return nil
	}

	Log.Warning(se.cache.Sign, sign)
	if se.cache.Sign.before(sign) {
		Log.Warning("syncEntry sync")
		//cache is old, featch data frome server
		c, err := se.remote.pull()
		if nil != err {
			Log.Error(err)
			return err
		}
		// if not modify, do not Restore
		if err := se.entry.Restore(c.Buf); nil != err {
			Log.Error(err)
			return err
		}
		se.cache.clone(c)
		se.cache.Sign = *sign
		se.cache.save()
		// TODO write cache
		return nil
	}

	// cache is new, push to server
	if err := se.remote.push(se.cache); nil != err {
		Log.Error(err)
		return err
	}
	sign, err = se.remote.sign()
	if nil != err {
		Log.Error(err)
		return err
	}

	se.cache.Sign = *sign
	se.cache.save()
	return nil
}
