package worker

import (
	"time"

	"github.com/golang/glog"
	"github.com/qetuantuan/jengo_recap/registry"
	"github.com/qetuantuan/jengo_recap/task"
)

type Scheduler struct {
	Base
	workerRegistry *registry.WorkerRegistry
}

func (scheduler *Scheduler) Run() {
	for {
		// Get task.Schedule from listening topic
		messages, err := scheduler.Queue.GetMsgExclusively(scheduler.ListeningTopic, 1)
		if err != nil {
			glog.Errorf("get message from queue error:%v", err)
			continue
		}
		for _, message := range messages {
			msg := message.Body().(task.General)

			for _, topic := range scheduler.OutputTopic {
				if err := scheduler.Queue.WriteMsgs(topic, []interface{}{msg}); err != nil {
					glog.Errorf("write message into topic:%s failed! err=%v", topic, err)
					continue
				}
			}
			if err = scheduler.Queue.AckMsg(scheduler.ListeningTopic, message, true); err != nil {
				glog.Errorf("ack message failed! err=%v", err)
				continue
			}
		}
		time.Sleep(time.Second * 1)
	}
}
