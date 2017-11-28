package worker

import (
	"time"

	"github.com/golang/glog"
	"github.com/qetuantuan/jengo_recap/api"
	"github.com/qetuantuan/jengo_recap/client"
	"github.com/qetuantuan/jengo_recap/constant"
	"github.com/qetuantuan/jengo_recap/dao"
	"github.com/qetuantuan/jengo_recap/queue"
	"github.com/qetuantuan/jengo_recap/service"
	"github.com/qetuantuan/jengo_recap/task"
)

type Finalize struct {
	Base
	ProjectClient client.ProjectStoreClientInterface
	RunDao        dao.RunDaoInterface
}

func (w *Finalize) Work() {
	glog.Info("(TODO: $WorkerId) started")
	for {
		msgs, err := w.Queue.GetMsgExclusively(queue.TopicFinalizeGroup1Name, 1)
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
			updateData := make(map[string]interface{})
			updateData["run.status"] = constant.RUN_FINISHED
			updateData["run.state"] = constant.RUN_STATE_SUCCESS
			updateData["run.updatedat"] = time.Now().UTC()

			if err := w.RunDao.UpdateRunProperties(t.Run.Id, updateData); err != nil {
				glog.Errorf("Failed to update run status as finished. Error is %v", err)
				continue
			}

			now := time.Now().UTC()
			var patchRun = &api.PatchRun{
				Id:        t.Run.Id,
				ProjectId: t.Run.ProjectId,
				UserId:    t.Run.UserId,
				Status:    constant.RUN_FINISHED,
				State:     constant.RUN_STATE_SUCCESS,
				UpdatedAt: &now,
			}
			if _, _, err := w.ProjectClient.UpdatePatchRun(patchRun, t.Run.BuildId); err != nil {
				glog.Errorf("Failed to update run status as finished. Error is %v", err)
				continue
			}

			//TODO: refactor later for PutLog
			if r, bytes, err := w.ProjectClient.PutLog(
				patchRun.UserId,
				patchRun.ProjectId,
				t.Run.BuildId,
				patchRun.Id,
				patchRun.Id+service.GetRandomString(100-len(patchRun.Id))); err != nil {

				glog.Errorf("Project Service Response: %v %v %v", r.Status, string(bytes), err)
				continue
			}

			// TODO: handle ack error with retry
			w.Queue.AckMsg(queue.TopicFinalizeGroup1Name, msgs[0], true)
			glog.Infof("%s is finshed", t.Run.Id)
		}
	}
}
