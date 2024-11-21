package round_robin

import "sync"

type RoundRobin struct {
	id []int64
	name []string
	connType []int64
	service []string
	index int
	mutex sync.Mutex
}

var globalScheduler *RoundRobin

func InitRoundRobin(id []int64, name []string, connType []int64, service []string) {
	globalScheduler = &RoundRobin{
		id: id,
		name: name,
		connType: connType,
		service: service,
		index: 0,
	}
}

func NextTask() (int64, string, int64, string) {
	globalScheduler.mutex.Lock()
	defer globalScheduler.mutex.Unlock()
	id := globalScheduler.id[globalScheduler.index]
	name := globalScheduler.name[globalScheduler.index]
	connType := globalScheduler.connType[globalScheduler.index]
	service := globalScheduler.service[globalScheduler.index]
	globalScheduler.index = (globalScheduler.index + 1) % len(globalScheduler.service)
	return id, name, connType, service
}