package ramstore

import (
	"sync"
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
}

type saveObj struct {
	Key string
	Obj Obj
}
