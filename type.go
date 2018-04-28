package ramstore

import (
	"sync"
	"time"
)

type stor struct {
	Data map[string]Obj
	sync.RWMutex
}

// Obj - Объект зранения
type Obj struct {
	Time    int64
	Data    []byte
	Deleted bool
	Expire  int
}

func (o *Obj) checkOK() bool {
	if o.Deleted || (o.Expire > 0 && o.Expire < int(time.Now().Unix())) {
		return false
	}

	return true
}

type saveObj struct {
	Key string
	Obj Obj
}
