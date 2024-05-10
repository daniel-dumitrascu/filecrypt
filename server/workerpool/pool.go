package workerpool

import (
	"server/config"
	"server/request"
	"server/utils"
)

type Pool struct {
	tasks chan *request.RequestData
}

func (p *Pool) Init(nrRoutines int, handlers *[4]func(req *request.RequestData)) {
	p.tasks = make(chan *request.RequestData, config.Tasks_channel_size)

	for id := 1; id <= nrRoutines; id++ {
		go worker(id, p.tasks, handlers)
	}
}

func (p *Pool) AddTask(req *request.RequestData) {
	select {
	case p.tasks <- req:
	default:
		log := utils.GetLogger()
		log.Error("Server is full at the moment. Please try later")
	}
}

func worker(id int,
	tasks chan *request.RequestData,
	handlers *[4]func(req *request.RequestData)) {

	for req := range tasks {
		handlers[req.ActionType](req)
	}
}
