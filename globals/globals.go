package globals

import (
	"sync"
)

var WaitGroup sync.WaitGroup

func GetWaitGroup() *sync.WaitGroup {
	return &WaitGroup
}
