package service

import (
	"fmt"

	"github.com/golang/glog"
	"github.com/qetuantuan/jengo_recap/api"
	"github.com/qetuantuan/jengo_recap/client"
	"github.com/qetuantuan/jengo_recap/dao"
	"github.com/qetuantuan/jengo_recap/model"
	"github.com/qetuantuan/jengo_recap/scm"
	"gopkg.in/mgo.v2"
)

type User api.User02
type ProjectServiceInterface interface {
	UpdateProjects(userId string) (projects []api.Project, err error)
	SwitchProject(userId, projectId string, enableStatus bool) (err error)
	//GetProjects(userId string) (project *dao.Project, err error)
	GetProject(projectId string) (project api.Project, err error)

	GetProjectsByFilter(filter map[string]interface{}, limitCount, offset int) (projects []api.Project, err error)
}

type ProjectService struct {
	Md        *dao.ProjectMongoDao
	GithubScm *scm.GithubScm
	UsClient  client.UserStoreClientInterface
}

/*
func getUserRemotely(userId string) (user User, err error) {
	uri := usHost + "/v0.2/internal_user/" + userId
	resp, err := http.Get(uri)
	if err != nil {
		return
	}
	if resp.StatusCode != http.StatusOK {
		err = errors.New("response is not 200!" + strconv.Itoa(resp.StatusCode))
		return
	}
	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(respBytes, &user)

	return
}
*/

func (p *ProjectService) GetProject(projectId string) (project api.Project, err error) {
	var proj model.Project
	proj, err = p.Md.GetProject(projectId)

	if err != nil {
		if err == mgo.ErrNotFound {
			glog.Warning("project not found: " + projectId)
			err = NotFoundError
			return
		} else {
			glog.Error("get project failed! error: ", err)
			err = MongoError
			return
		}
	}
	return *proj.ToApiObj(), nil
}

// syncProjectSet sync old projects status to new projects, missing ones as returning.
// n: new projects map, whose key is project.Meta.Id.
// o: old projects.
// d: delete projects (missing ones in new projects)
// s: projects need update (in two sets)
// i: insert projects (missing in old projects)
func syncProjectSet(n map[string]*model.Project, o map[string]*model.Project) (d, s, i []model.Project) {
	// assuming new project.Enable = False by default
	for _, oldProject := range o {
		newProject, ok := n[oldProject.Meta.Id]
		if ok {
			s = append(s, *newProject)
		} else {
			d = append(d, *oldProject)
		}
	}
	for _, newProject := range n {
		if _, ok := o[newProject.Meta.Id]; !ok {
			i = append(i, *newProject)
		}
	}
	return
}
func (p *ProjectService) UpdateProjects(userId string) (projects []api.Project, err error) {
	user, err := p.UsClient.GetUser(userId)
	if err != nil {
		glog.Errorf("update projects: get user from us failed! user_id: %v, error: %v", userId, err)
		return
	}
	p.GithubScm.User = user.Auth.LoginName
	p.GithubScm.Token = user.Auth.Token
	// TODO:  if auth.AuthSource == 'github'
	var projs []model.Project
	projs, err = p.GithubScm.GetProjectList()
	if err != nil {
		glog.Error("get projects from github for user failed", p.GithubScm.User)
		err = ScmError
		return
	}
	glog.Infof("get %v projects from github for user %v", len(projs), p.GithubScm.User)

	// get projects in mongodb to update enable field
	oldProjects, tmpErr := p.Md.GetProjects(userId, 0, 0)
	if tmpErr != nil {
		glog.Errorf("Query mongo failed: %v", err)
		err = MongoError
		return
	}

	glog.Infof("get %v projects from mongo for user %v", len(oldProjects), userId)

	newProjectMap := make(map[string]*model.Project)
	oldProjectMap := make(map[string]*model.Project)
	for i := range projs {
		newProjectMap[projs[i].Meta.Id] = &projs[i]
		glog.Infof("one new:%s", projs[i].Meta.Id)
	}
	for i := range oldProjects {
		oldProjectMap[oldProjects[i].Meta.Id] = &oldProjects[i]
		glog.Infof("one old:%s", oldProjects[i].Meta.Id)
	}
	glog.Infof("new:%d, old:%d", len(newProjectMap), len(oldProjectMap))
	deleteProjects, updateProjects, insertProjects := syncProjectSet(newProjectMap, oldProjectMap)
	deleteNum := len(deleteProjects)
	updateNum := len(updateProjects)
	insertNum := len(insertProjects)
	glog.Infof("Will delete %d projects, update %d projects, insert %d projects", deleteNum, updateNum, insertNum)
	if deleteNum != 0 {
		glog.Infof("Deleting %v non-existed projects...", deleteNum)
		err = p.Md.DeleteProjects(deleteProjects, userId)
		if err != nil {
			glog.Errorf(" delete projects from mongo failed: %v", err)
			err = MongoError
			return
		}
	}
	if updateNum != 0 {
		glog.Infof("Updating %v projects...", updateNum)
		err = p.Md.UpdateProjects(updateProjects, userId)
		if err != nil {
			glog.Errorf(" update projects to mongo failed: %v", err)
			err = MongoError
			return
		}
	}
	if insertNum != 0 {
		glog.Infof("Inserting %v projects...", insertNum)
		err = p.Md.InsertProjects(insertProjects, userId)
		if err != nil {
			glog.Errorf(" insert projects to mongo failed: %v", err)
			err = MongoError
			return
		}
	}
	for _, p := range projs {
		projects = append(projects, *p.ToApiObj())
	}
	return
}

func (p *ProjectService) GetProjectsByFilter(filter map[string]interface{}, limitCount, offset int) (projects []api.Project, err error) {
	var projs []model.Project
	projs, err = p.Md.GetProjectsByFilter(filter, limitCount, offset)
	if err != nil {
		glog.Error("get project failed! error: ", err)
		err = MongoError
		return
	}
	for _, p := range projs {
		projects = append(projects, *p.ToApiObj())
	}
	return
}
func (p *ProjectService) SwitchProject(userId, projectId string, enableStatus bool) (err error) {
	user, err := p.UsClient.GetUser(userId)
	if err != nil {
		glog.Warningf("Get user from us failed! user_id:%s, errpr:%v", userId, err)
		return
	}
	p.GithubScm.User = user.Auth.LoginName
	p.GithubScm.Token = user.Auth.Token
	// TODO:  if auth.AuthSource == 'github'
	project, err := p.Md.GetProject(projectId)
	if err != nil {
		if err == mgo.ErrNotFound {
			glog.Warning("project not found: " + projectId)
			err = NotFoundError
			return
		} else {
			glog.Error("get project failed! error: ", err)
			err = MongoError
			return
		}
	}
	if enableStatus {
		_, err = p.GithubScm.SetHook(project.Meta.Name)
		if err != nil && err != scm.HookExistError {
			glog.Warningf("enable hook failed! error:%v", err)
			err = ScmError
			return
		}
	} else { //do disable the hook
		hook, errTmp := p.GithubScm.GetHook(project.Meta.HooksUrl) //github exist this hook or not
		if errTmp != nil {
			err = errTmp
			if errTmp != scm.HookNonExistError { // err returned but not hooknonexist
				glog.Warningf("get hook failed from github! error:%v", err)
				err = ScmError
				return
			} else { //err is HookNonExistError, github no such hook
				glog.Warning("found no such hook from github when delete hook")
				// not return since we need to update project status
				//i'd like to not return response with body
			}
		} else { // found hook from github
			glog.Info("found hook from github ,now begin to delete hook")

			url := hook.Url

			if err = p.GithubScm.DeleteHook(url); err != nil {
				glog.Warningf("Delete hook failed! error:%v", err)
				err = ScmError
				return
			}
		}

	}

	err1 := p.Md.SwitchProject(projectId, enableStatus)
	if err1 != nil {
		glog.Warningf("switch project failed! error:%v", err)
		err = MongoError
		return
	}
	err = nil //in case of dao.HookNonExistError or dao.HookExistError
	glog.Info(fmt.Sprintf("swith project %s to %s success", projectId, enableStatus))
	return
}
