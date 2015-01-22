package sync

import (
	"time"
)

type signature struct {
	ModifyTime time.Time
	Sha1       string
}

func (s *signature) equal(s1 *signature) bool {
	return s.ModifyTime.Equal(s1.ModifyTime) && s.Sha1 == s1.Sha1
}

func (s *signature) before(s1 *signature) bool {
	return s.ModifyTime.Before(s1.ModifyTime)
}

func (s *signature) after(s1 *signature) bool {
	return s.ModifyTime.After(s1.ModifyTime)
}
