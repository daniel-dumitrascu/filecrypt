package workerpool

import (
	"fmt"
	"server/config"
	"server/request"
)

type Pool struct {
	tasks chan *request.RequestData
}

func (p *Pool) Init(nrRoutines int, handlers *[3]func(req *request.RequestData)) {
	p.tasks = make(chan *request.RequestData, config.Tasks_channel_size)

	for id := 1; id <= nrRoutines; id++ {
		go worker(id, p.tasks, handlers)
	}
}

func (p *Pool) AddTask(req *request.RequestData) {
	select {
	case p.tasks <- req:
	default:
		fmt.Println("Server is full at the moment. Please try later")
	}
}

func worker(id int,
	tasks chan *request.RequestData,
	handlers *[3]func(req *request.RequestData)) {

	for req := range tasks {
		handlers[req.ActionType](req)
	}
}
