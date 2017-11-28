package worker

import (
	"time"

	"github.com/golang/glog"
	"github.com/qetuantuan/jengo_recap/action"
	"github.com/qetuantuan/jengo_recap/constant"
	"github.com/qetuantuan/jengo_recap/queue"
	"github.com/qetuantuan/jengo_recap/task"
)

type Executor struct {
	Base
}

const (
	INPUT_DATA_ERROR = "input data error"
	CHECK_INTERVAL   = 10
	CHAN_CAP         = 1000
)

type executeInfo struct {
	Result  action.ActionState
	Reason  string
	Message queue.Msg
}

func processResult(resChan chan executeInfo, executor *Executor) {
	for {
		res := <-resChan
		//todo: add retry logic
		taskInfo := res.Message.Body().(task.General)
		taskInfo.Run.State = action.ActionStateToString(res.Result)
		taskInfo.Run.Status = constant.RUN_FINISHED
		for _, outTopic := range executor.OutputTopic {
			err := executor.Queue.WriteMsgs(outTopic, []interface{}{res.Message.Body()})
			if err != nil {
				glog.Errorf("Write msg to queue:%s failed! err:%v", outTopic, err)
			}
		}
		executor.Queue.AckMsg(executor.ListeningTopic, res.Message, true)
	}
}

func (executor *Executor) Run() {
	resultChan := make(chan executeInfo, CHAN_CAP)
	go processResult(resultChan, executor)
	for {
		// Get task.Execute from listening topic
		messages, err := executor.Queue.GetMsgExclusively(executor.ListeningTopic, 1)
		if err != nil {
			glog.Errorf("get message from queue error:%v", err)
			continue
		}
		for _, message := range messages {
			taskInfo := message.Body().(task.General)
			if taskInfo.Version < 1 {
				glog.Errorf("")
				continue
			}
		}
		// Get manifest and template from task.Schedule
		// Spawn a runner process
		//   runner process
		//     create context
		//     put manifest to context
		//     create template-step-actions object tree from template string
		//     run actions
		//     report status to executeWorker
		// if child.wait(), output task.Finalize with status
		// if parallel == max, wait for one child finish
		time.Sleep(time.Second * 1)
	}
}
