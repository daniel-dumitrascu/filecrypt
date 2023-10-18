package workerpool

import (
	"server/request"
)

func Init(nrRoutines int,
	tasks chan request.RequestData,
	handlers *[3]func(req *request.RequestData)) {

	for id := 1; id <= nrRoutines; id++ {
		go worker(id, tasks, handlers)
	}
}

func worker(id int,
	tasks chan request.RequestData,
	handlers *[3]func(req *request.RequestData)) {

	for task := range tasks {
		handlers[task.ActionType](&task)
	}
}
