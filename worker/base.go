package worker

import (
	"github.com/qetuantuan/jengo_recap/queue"
)

type Base struct {
	ListeningTopic string
	OutputTopic    []string

	Id   string
	Name string
	Ip   string
	Type string

	Queue queue.TaskQueue
}
