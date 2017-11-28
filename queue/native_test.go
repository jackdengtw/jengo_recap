package queue

import (
	"testing"

	"github.com/qetuantuan/jengo_recap/task"
)

func TestWaterfall(t *testing.T) {
	var q TaskQueue
	q = NewNativeTaskQueue()
	tsk := task.General{Status: "1"}

	go func() {
		if err := q.WriteMsgs(
			TopicParseGroup1Name,
			[]interface{}{
				tsk,
			}); err != nil {
			t.Fatalf("Write parse msg failed: %v", err)
		}
	}()

	go func() {
		if msgs, err := q.GetMsgExclusively(TopicParseGroup1Name, 1); err != nil {
			t.Fatalf("Get parse msg failed: %v", err)
		} else {
			if len(msgs) != 1 {
				t.Fatalf("Get %v parse msg instead of 1", len(msgs))
			}
			if err := q.AckMsg(TopicParseGroup1Name, msgs[0], true); err != nil {
				t.Fatalf("ack parse msg failed " + err.Error())
			}
			if err := q.WriteMsgs(
				TopicScheduleGroup1Name,
				[]interface{}{
					msgs[0].Body().(task.General),
				}); err != nil {
				t.Fatalf("Write schedule msg failed: %v", err)
			}
		}
	}()

	if msgs, err := q.GetMsgExclusively(TopicScheduleGroup1Name, 1); err != nil {
		t.Fatalf("Get schedule msg failed: %v", err)
	} else {
		if len(msgs) != 1 {
			t.Fatalf("Get %v schedule msg instead of 1", len(msgs))
		}
		q.AckMsg(TopicScheduleGroup1Name, msgs[0], true)
	}

	if len(q.(*NativeTaskQueue).HoldingMsgs) > 0 {
		t.Fatalf("Holding msg not clean up")
	}
}

func TestAckNotDone(t *testing.T) {
	var q TaskQueue
	q = NewNativeTaskQueue()
	tsk := task.General{Status: "1"}

	go func() {
		if err := q.WriteMsgs(
			TopicParseGroup1Name,
			[]interface{}{
				tsk,
			}); err != nil {
			t.Fatalf("Write parse msg failed: %v", err)
		}
	}()

	go func() {
		if msgs, err := q.GetMsgExclusively(TopicParseGroup1Name, 1); err != nil {
			t.Fatalf("Get parse msg failed: %v", err)
		} else {
			if len(msgs) != 1 {
				t.Fatalf("Get %v parse msg instead of 1", len(msgs))
			}
			if err := q.AckMsg(TopicParseGroup1Name, msgs[0], false); err != nil {
				t.Fatalf("ack parse msg failed " + err.Error())
			}
		}
	}()

	if msgs, err := q.GetMsgExclusively(TopicParseGroup1Name, 1); err != nil {
		t.Fatalf("Get Parse msg again failed: %v", err)
	} else {
		if len(msgs) != 1 {
			t.Fatalf("Get %v schedule msg instead of 1", len(msgs))
		}
		q.AckMsg(TopicParseGroup1Name, msgs[0], true)
	}

	if len(q.(*NativeTaskQueue).HoldingMsgs) > 0 {
		t.Fatalf("Holding msg not clean up")
	}
}
