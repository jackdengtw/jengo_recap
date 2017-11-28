package worker

import (
	"errors"
	"time"

	_ "gopkg.in/yaml.v2"

	"github.com/golang/glog"
	"github.com/qetuantuan/jengo_recap/api"
	"github.com/qetuantuan/jengo_recap/client"
	"github.com/qetuantuan/jengo_recap/definition"
	"github.com/qetuantuan/jengo_recap/model"
	"github.com/qetuantuan/jengo_recap/queue"
	"github.com/qetuantuan/jengo_recap/scm"
	"github.com/qetuantuan/jengo_recap/task"
)

type Parse struct {
	Base
	Uc client.UserStoreClientInterface
	Pc client.ProjectStoreClientInterface
	Gc scm.Scm
}

func (w *Parse) Work() {
	for {
		// Get task.Parse from listening topic
		// http get yml file
		// parse yml into manifest
		// TODO later: parse yml into template
		// select template with manifest
		// clone template and serialize template to string
		// output task.Schedule with manifest and template

		// test yaml
		// yaml.Marshal()

		glog.Info("try to get a task from queue")
		msgs, err := w.Queue.GetMsgExclusively(w.ListeningTopic, 1)
		if err != nil {
			glog.Errorf("call GetMsgExclusively error: err(%v)", err)
			continue
		}
		if len(msgs) < 1 {
			glog.Info("there's no message right now...")
			continue
		}

		for _, msg := range msgs {

			// parse each message received from queue to task.General
			t, ok := msg.Body().(task.General)
			if !ok {
				glog.Errorf("not expected task type: m(%v)", msg)
				// TODO: handle ack error with retry
				if err = w.Queue.AckMsg(w.ListeningTopic, msg, true); err != nil {
					err = errors.New("not expected task type and call AckMsgs error: " + err.Error())
					continue
				}
				continue
			}

			// process each task.General
			if err = w.Process(t); err != nil {
				glog.Errorf("process error: err(%v), msg(%v)", err, msg)
				continue
			}

			// ack each well-processed message
			// TODO: handle ack error with retry
			if err = w.Queue.AckMsg(w.ListeningTopic, msg, true); err != nil {
				glog.Errorf("call AckMsgs error: err(%v)", err)
				continue
			}

			glog.Infof("finished to process task: t(%s)", msg.Body())
		}
	}
}

func (w *Parse) Process(t task.General) (err error) {
	glog.Infof("started to Process")

	var content []byte
	var template definition.Template
	var manifest definition.Manifest

	if t, content, err = w.GetYml(t); err != nil {
		glog.Errorf("call GetYml error: err(%v)", err)
		return
	}

	if template, manifest, err = w.ParseYml(content); err != nil {
		glog.Errorf("call ParseYml error: err(%v)", err)
		// TODO: error log to log service
		return
	}

	if err = w.MoveOn(t, template, manifest); err != nil {
		glog.Errorf("call MoveOn error: err(%v)", err)
		return
	}

	glog.Infof("succeeded to Process: t(%v)", t)
	return
}

// download yaml file from scm
func (w *Parse) GetYml(t task.General) (tRet task.General, content []byte, err error) {
	glog.Infof("started to GetYml")

	// #1 get user info from userId
	var user api.User02
	if user, err = w.Uc.GetUser(t.Run.UserId); err != nil {
		glog.Errorf("failed to GetUser: userId(%s) err(%v)", t.Run.UserId, err)
		return
	}
	t.User = &user

	// #2 set token for github-client
	w.Gc.SetToken(t.User.Auth.Token)

	// #3 get project info from (userId,projectId)
	var project api.Project
	if project, err = w.Pc.GetProject(t.Run.UserId, t.Run.ProjectId); err != nil {
		glog.Errorf("failed to GetProject: userId(%s) projectId(%s) err(%v)",
			t.Run.ProjectId, t.Run.UserId, err)
		return
	}
	t.Project = model.NewProjectFrom(&project)

	// return this tRet which is filled with User and Project data
	tRet = t

	// #4 download project yml file's from composed http download url
	content, err = w.Gc.GetYmlContent(t.Project.Meta.FullName, t.Run.Branch)
	glog.Infof("succeeded to GetYmlContent: content(%s)", content)
	return
}

func (w *Parse) ParseYml(content []byte) (t definition.Template, m definition.Manifest, err error) {
	glog.Infof("started to ParseYml")
	return
}

func (w *Parse) MoveOn(ta task.General, te definition.Template, manifest definition.Manifest) (err error) {
	glog.Infof("started to MoveOn")

	// TODO: add other required info, such as manifest and template)
	var nextTasks []interface{}
	nextTasks = append(nextTasks, ta)

	// generate next task for the following worker, TODO: reopen later after scheduler worker is in
	/*
		for _, topic := range w.OutputTopic {
			glog.Infof("try to WriteMsgs to topic(%s)", topic)
			if err = w.Queue.WriteMsgs(topic, nextTasks); err != nil {
				glog.Errorf("call WriteMsgs error: err(%v)", err)
			}
		}
	*/

	// send to finalize worker for debugging, TODO: delete later
	glog.Infof("try to WriteMsgs to topic(%s)", queue.TopicFinalizeGroup1Name)
	if err = w.Queue.WriteMsgs(queue.TopicFinalizeGroup1Name, nextTasks); err != nil {
		glog.Errorf("call WriteMsgs error: err(%v)", err)
	}

	time.Sleep(30 * time.Second) // TODO: debug only

	glog.Infof("succeeded to MoveOn: task(%v)", ta)
	return
}
