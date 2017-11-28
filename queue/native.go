package queue

import (
	"errors"
	"time"

	"github.com/golang/glog"
	"github.com/qetuantuan/jengo_recap/task"
	"github.com/satori/go.uuid"
)

const (
	TopicStartGroup1Name    = "TopicStartGroup1"
	TopicParseGroup1Name    = "TopicParseGroup1"
	TopicScheduleGroup1Name = "TopicScheduleGroup1"
	TopicExecuteGroup1Name  = "TopicExecuteGroup1"
	TopicFinalizeGroup1Name = "TopicFinalizeGroup1"

	AckTimeout = 20 * time.Minute
	bufferSize = 20
)

var (
	errGeneralTaskType = errors.New("Wrong type of msg sent to queue")
)

type NativeMsg struct {
	id   string
	body interface{}

	T     time.Time
	Topic string
}

func (m *NativeMsg) Id() string {
	return m.id
}

func (m *NativeMsg) Body() interface{} {
	return m.body
}

type NativeTaskQueue struct {
	Chans map[string]chan task.General

	HoldingMsgs map[string]*NativeMsg
}

func NewNativeTaskQueue() *NativeTaskQueue {
	return &NativeTaskQueue{
		Chans: map[string]chan task.General{
			TopicStartGroup1Name:    make(chan task.General, bufferSize),
			TopicParseGroup1Name:    make(chan task.General, bufferSize),
			TopicScheduleGroup1Name: make(chan task.General, bufferSize),
			TopicExecuteGroup1Name:  make(chan task.General, bufferSize),
			TopicFinalizeGroup1Name: make(chan task.General, bufferSize),
		},
		HoldingMsgs: make(map[string]*NativeMsg),
	}
}

func (q *NativeTaskQueue) GetMsgExclusively(topic string, size int) (msgs []Msg, err error) {
	// TODO: This version is blocking, may provide non-blocking method later
	//       anyway return 1 msg for now
	glog.Infof("GetMsgExclusively: topic: %s", topic)
	ch, ok := q.Chans[topic]
	if !ok {
		err = errors.New("topic not found: " + topic)
		return
	}

	task := <-ch
	m := &NativeMsg{
		id:    uuid.NewV4().String(),
		body:  task,
		T:     time.Now().UTC(),
		Topic: topic,
	}
	q.HoldingMsgs[m.Id()] = m
	glog.Info(topic+", m: %v", m)
	msgs = append(msgs, m)

	return
}

func (q *NativeTaskQueue) AckMsg(topic string, msg Msg, done bool) (err error) {
	_, ok := q.Chans[topic]
	if !ok {
		err = errors.New("topic not found: " + topic)
		return
	}

	if _, ok := q.HoldingMsgs[msg.Id()]; ok {
		if !done { // write back
			if err = q.WriteMsgs(topic, []interface{}{msg.Body()}); err == nil {
				delete(q.HoldingMsgs, msg.Id())
			} else {
				err = errors.New("Ack failed msg error: " + err.Error())
			}
			return
		} else {
			delete(q.HoldingMsgs, msg.Id())
		}
	} else {
		err = errors.New("Acking non-existence msg")
	}
	return
}

func (q *NativeTaskQueue) WriteMsgs(topic string, objs []interface{} /* Msg data */) (err error) {
	// bruce force error handling
	ch, ok := q.Chans[topic]
	if !ok {
		err = errors.New("topic not found: " + topic)
		return
	}
	for i := 0; i < len(objs); i++ {
		t, ok := objs[i].(task.General)
		if ok {
			glog.Infof("write msg %d: %v", i, t)
			ch <- t
		} else {
			err = errGeneralTaskType
			break
		}
	}
	return
}

func (q *NativeTaskQueue) Start() {
	go func() {
		for {
			now := time.Now().UTC().Add(-10 * time.Second) // 10s buffer
			for k, v := range q.HoldingMsgs {
				if now.Sub(v.T) > AckTimeout {
					if err := q.WriteMsgs(v.Topic, []interface{}{v.Body()}); err == nil {
						delete(q.HoldingMsgs, k)
					}
				}
			}
			time.Sleep(1 * time.Minute)
		}
	}()
}
