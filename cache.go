package sync

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"pkg.deepin.io/lib/utils"
)

const (
	cacheDir = "/.cache/deepin/sync/entry/"
)

type cache struct {
	entry Entry
	Buf   []byte
	Sign  signature
}

func checkCreateDir(userCacheDir string) error {
	fi, err := os.Stat(userCacheDir)
	if nil == err && !fi.IsDir() {
		os.RemoveAll(userCacheDir)
	}
	return os.MkdirAll(userCacheDir, 0775)
}

func newCache(e Entry) (*cache, error) {
	c := &cache{entry: e}
	checkCreateDir(c.locleSpace())
	err := c.load()
	return c, err
}

func (c *cache) locleSpace() string {
	return utils.GetHomeDir() + cacheDir + c.entry.ID()
}

func (c *cache) clone(cl *cache) error {
	c.Buf = cl.Buf
	c.Sign = cl.Sign
	return nil
}

func (c *cache) save() error {
	head := c.locleSpace() + "/head"
	data, err := json.Marshal(c)
	if nil != err {
		return err
	}
	return ioutil.WriteFile(head, data, 0644)
}

func (c *cache) load() error {
	head := c.locleSpace() + "/head"
	data, err := ioutil.ReadFile(head)
	var headCache cache
	if err = json.Unmarshal(data, &headCache); nil != err {
		Log.Info(err)
	}
	return c.clone(&headCache)
}
