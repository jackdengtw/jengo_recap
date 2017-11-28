package worker

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang/glog"
	"github.com/qetuantuan/jengo_recap/api"
	"github.com/qetuantuan/jengo_recap/client"
	"github.com/qetuantuan/jengo_recap/dao"
	"github.com/qetuantuan/jengo_recap/queue"
	"github.com/qetuantuan/jengo_recap/task"
)

type Start struct {
	Base
	ProjectClient client.ProjectStoreClientInterface
	RunDao        dao.RunDaoInterface
}

func (w *Start) Work() {
	glog.Info("(TODO: $WorkerId) started.")
	for {
		msgs, err := w.Queue.GetMsgExclusively(queue.TopicStartGroup1Name, 1)
		if err != nil {
			glog.Errorf("(TODO: $WorkerId): GetMsgExclusively error: %v", err)
			continue
		}
		if len(msgs) < 1 {
			glog.Info("(TODO: $WorkerId): There's no tasks right now..")
			time.Sleep(1 * time.Second)
			continue
		}

		t, ok := msgs[0].Body().(task.General)
		glog.Info("(TODO: $WorkerId): Get Msg: %v", t)
		if !ok {
			glog.Errorf("cast last msg error: task: %v", t)
		} else {
			run := t.Run
			// parse worker update task to status=running
			resp, raw, err := w.ProjectClient.InsertRun(&run.Run)
			if err != nil {
				glog.Errorf("Failed to insert run to ps: %v", err)
				continue
			}
			if resp.StatusCode != http.StatusOK {
				glog.Errorf("Failed to get 200 from s.ProjectClient.InsertRun")
				continue
			}

			buildResponse := &api.ProjectStoreCreateRunResponse{}
			err = json.Unmarshal(raw, buildResponse)
			if err != nil {
				glog.Errorf("Failed to unmarshal project store insert run's response: %v, %s", err, string(raw))
				continue
			}

			run.BuildId = buildResponse.BuildId
			glog.Infof("task: %v", t)

			updateData := make(map[string]interface{})
			updateData["run.build_id"] = t.Run.BuildId

			if err := w.RunDao.UpdateRunProperties(t.Run.Id, updateData); err != nil {
				glog.Errorf("Failed to update build id. Error is %v", err)
				continue
			}

			if err := w.Queue.WriteMsgs(
				// change to parse topic for debugging, should be w.OutputTopic
				queue.TopicParseGroup1Name,
				[]interface{}{
					t,
				}); err != nil {
				glog.Errorf("Write finalize msg failed: %v", err)
			}

			// TODO: handle ack error with retry
			w.Queue.AckMsg(queue.TopicStartGroup1Name, msgs[0], true)
			glog.Infof("%s is finshed", t.Run.Id)
		}
	}
}
