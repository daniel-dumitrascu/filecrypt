package workerpool

import (
	"server/request"
)

// TODO - adauga array-ul la init si mai apoi foloseste ca param la processor
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
