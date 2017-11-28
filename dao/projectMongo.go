package dao

import (
	"errors"
	"fmt"

	"github.com/golang/glog"
	"github.com/qetuantuan/jengo_recap/algo"
	"github.com/qetuantuan/jengo_recap/model"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var ErrorBuildNotFind = errors.New("build not found")

var projectDbName string = "projects"
var projectCol string = "project"

//var runCol string = "run"
var buildCol string = "build"
var logCol string = "log"
var hookCol string = "hook"

type ProjectMongoDao struct {
	Url      string
	GSession *mgo.Session
}

func (md *ProjectMongoDao) Init() error {
	var err error
	md.GSession, err = mgo.Dial(md.Url)
	return err
}

func (md *ProjectMongoDao) UpdateProjects(projects []model.Project, userId string) (err error) {
	session := md.GSession.Copy()
	defer session.Close()
	pc := session.DB(projectDbName).C(projectCol)
	for _, project := range projects {
		storeId, tmpErr := algo.To16Bytes(project.Meta.Id)
		if tmpErr != nil {
			err = tmpErr
			glog.Warningf("Update project failed! get hash id failed![%v]", err)
			break
		}
		err = pc.UpdateId(storeId, bson.M{"$set": bson.M{"project.meta": project.Meta}, "$addToSet": bson.M{"project.users": userId}})
		if err != nil {
			glog.Errorf("Update failed for %v %v %v", project.Meta.Id, string(storeId), err)
			break
		}
		glog.Errorf("Update Success for %v %v", project.Meta.Id, string(storeId))
	}
	return
}

// todo: maybe can merge InsertProjects&UpdateProjects to UpsertProjects function
func (md *ProjectMongoDao) InsertProjects(projects []model.Project, userId string) (err error) {
	session := md.GSession.Copy()
	defer session.Close()
	pc := session.DB(projectDbName).C(projectCol)
	for _, project := range projects {
		storeId, tmpErr := algo.To16Bytes(project.Meta.Id)
		if tmpErr != nil {
			err = tmpErr
			glog.Warningf("insert project failed! get hash id failed![%v]", err)
			break
		}
		_, err = pc.UpsertId(storeId, bson.M{"$set": bson.M{"project.meta": project.Meta}, "$addToSet": bson.M{"project.users": userId}})
		if err != nil {
			glog.Warningf("insert project failed for %v %v %v", project.Meta.Id, string(storeId), err)
			break
		}
		glog.Infof("insert project success for %v %v", project.Meta.Id, string(storeId))
	}
	return
}

func (md *ProjectMongoDao) DeleteProjects(projects []model.Project, userId string) (err error) {
	session := md.GSession.Copy()
	defer session.Close()
	pc := session.DB(projectDbName).C(projectCol)
	for _, project := range projects {
		storeId, tmpErr := algo.To16Bytes(project.Meta.Id)
		if tmpErr != nil {
			err = tmpErr
			glog.Warningf("Delete project failed! get hash id failed![%v]", err)
			break
		}
		//todo: use removeall
		//tmpErr := pc.RemoveId(storeId)
		err = pc.UpdateId(storeId, bson.M{"$pull": bson.M{"project.users": userId}})
		if err != nil {
			glog.Warningf("Delete project failed for %v %v %v", project.Meta.Id, string(storeId), err)
		} else {
			glog.Infof("Delete project success for %v %v", project.Meta.Id, string(storeId))
		}
	}
	return
}

func (md *ProjectMongoDao) GetProjectsByFilter(filter map[string]interface{}, limitCount, offset int) (projects []model.Project, err error) {
	session := md.GSession.Copy()
	defer session.Close()
	pc := session.DB(projectDbName).C(projectCol)
	projectFilter := bson.M{}
	for key, value := range filter {
		filterKey := "project." + key
		switch value.(type) {
		case bool:
			if value.(bool) {
				projectFilter[filterKey] = value
			} else {
				projectFilter[filterKey] = bson.M{"$ne": true}
			}
		case []string:
			projectFilter[filterKey] = bson.M{"$in": value}
		default:
			projectFilter[filterKey] = value
		}
	}
	fmt.Println(projectFilter)
	err = pc.Find(&projectFilter).
		Sort("-project.meta.createdat").Skip(offset).Limit(limitCount).All(&projects)
	return
}

func (md *ProjectMongoDao) GetProjectsByScms(userId string, scms []string) (projects []model.Project, err error) {
	session := md.GSession.Copy()
	defer session.Close()
	pc := session.DB(projectDbName).C(projectCol)
	err = pc.Find(&bson.M{"project.users": userId, "project.meta.scm": &bson.M{"$in": scms}}).
		Sort("-project.meta.createdat").All(&projects)
	return
}

func (md *ProjectMongoDao) GetProjects(userId string, limitCount, offset int) (projects []model.Project, err error) {
	session := md.GSession.Copy()
	defer session.Close()
	pc := session.DB(projectDbName).C(projectCol)
	err = pc.Find(bson.M{"project.users": userId}).Sort("-project.meta.createdat").Skip(offset).Limit(limitCount).All(&projects)
	return
}

func (md *ProjectMongoDao) GetBuildIndex(projectId string) (idx int, err error) {
	session := md.GSession.Copy()
	defer session.Close()
	pc := session.DB(projectDbName).C(projectCol)
	var storeId []byte
	if storeId, err = algo.To16Bytes(projectId); err != nil {
		return
	}
	var p model.Project
	change := mgo.Change{
		Update:    bson.M{"$inc": bson.M{"project.runindex": 1}},
		ReturnNew: true,
	}

	_, err = pc.Find(bson.M{"_id": storeId}).Apply(change, &p)
	idx = p.Project.RunIndex
	return
}

func (md *ProjectMongoDao) GetProject(id string) (project model.Project, err error) {
	session := md.GSession.Copy()
	defer session.Close()
	pc := session.DB(projectDbName).C(projectCol)
	var storeId []byte
	if storeId, err = algo.To16Bytes(id); err != nil {
		return
	}
	var p model.Project
	err = pc.FindId(storeId).One(&p)
	return p, err
}

func (md *ProjectMongoDao) SwitchProject(projectId string, enableStatus bool) (err error) {
	session := md.GSession.Copy()
	defer session.Close()
	pc := session.DB(projectDbName).C(projectCol)
	var storeId []byte
	if storeId, err = algo.To16Bytes(projectId); err != nil {
		return
	}
	err = pc.UpdateId(storeId, bson.M{"$set": bson.M{"project.enable": enableStatus}})
	return
}

/*
	Update project dynamic info. including: state latestBuildId, branch
*/
func (md *ProjectMongoDao) UpdateDynamicProjectInfo(projectId, state, latestBuildId, branch string) (err error) {
	session := md.GSession.Copy()
	defer session.Close()
	pc := session.DB(projectDbName).C(projectCol)
	var storeId []byte
	if storeId, err = algo.To16Bytes(projectId); err != nil {
		return
	}
	updateMap := bson.M{}

	setMap := bson.M{}
	if state != "" {
		setMap["project.state"] = state
	}
	if latestBuildId != "" {
		setMap["project.latestbuildid"] = latestBuildId
	}
	if len(setMap) != 0 {
		updateMap["$set"] = setMap
	}

	addToSetMap := bson.M{}
	if branch != "" {
		addToSetMap["project.branches"] = branch
	}

	if len(addToSetMap) != 0 {
		updateMap["$addToSet"] = addToSetMap
	}
	if len(updateMap) == 0 {
		err = errors.New("nothing to update")
		return
	}
	err = pc.UpdateId(storeId, updateMap)
	return
}

func (md *ProjectMongoDao) InsertBuild(build model.Build) (id string, err error) {
	session := md.GSession.Copy()
	defer session.Close()
	bc := session.DB(projectDbName).C(buildCol)
	oid := bson.NewObjectId()
	id = oid.Hex()
	build.Id = id
	err = bc.Insert(build)
	if err != nil {
		return
	}
	return
}

func (md *ProjectMongoDao) FindBuildByBranchCommit(projectId, commitId, branch string) (build model.Build, err error) {
	session := md.GSession.Copy()
	defer session.Close()
	bc := session.DB(projectDbName).C(buildCol)
	builds := model.Builds{}
	err = bc.Find(bson.M{"projectid": projectId, "commitid": commitId, "branch": branch}).All(&builds)
	if err != nil {
		return
	}
	if len(builds) <= 0 {
		err = ErrorBuildNotFind
		return
	}
	build = builds[0]
	return
}

func (md *ProjectMongoDao) RunExistInBuild(buildId, runId string) (res bool, err error) {
	session := md.GSession.Copy()
	defer session.Close()
	bc := session.DB(projectDbName).C(buildCol)
	num, err := bc.Find(bson.M{"_id": buildId, "runs._id": runId}).Count()
	res = num > 0
	return

}

func (md *ProjectMongoDao) UpdateRun(buildId string, run model.Run) (err error) {
	session := md.GSession.Copy()
	defer session.Close()
	bc := session.DB(projectDbName).C(buildCol)
	err = bc.Update(bson.M{"_id": buildId, "runs._id": run.Id}, bson.M{"$set": bson.M{"runs.$": run}})
	return
}

func (md *ProjectMongoDao) UpdateRunProperties(buildId string, runId string, p map[string]interface{}) (err error) {
	// transform p to map[string]interface{}
	// https://docs.mongodb.com/manual/reference/operator/update/set/
	var runInterface = make(map[string]interface{})
	for k, v := range p {
		runInterface["runs.$."+k] = v
	}
	glog.Infof("runInterface: %v", runInterface)

	// db.Update
	session := md.GSession.Copy()
	defer session.Close()
	bc := session.DB(projectDbName).C(buildCol)
	err = bc.Update(bson.M{"_id": buildId, "runs._id": runId}, bson.M{"$set": runInterface})
	if err != nil {
		msg := fmt.Sprintf("partial update p failed! error:%v", err)
		glog.Warning(msg)
		return
	}
	return
}

func (md *ProjectMongoDao) InsertRun(buildId string, run model.Run) (err error) {
	session := md.GSession.Copy()
	defer session.Close()
	bc := session.DB(projectDbName).C(buildCol)
	err = bc.UpdateId(buildId, bson.M{"$push": bson.M{"runs": run}})
	return
}

func (md *ProjectMongoDao) GetLatestBuild(projectIds []string) (latestBuilds model.Builds, err error) {
	session := md.GSession.Copy()
	defer session.Close()
	pc := session.DB(projectDbName).C(projectCol)
	var storeIds [][]byte
	for _, projectId := range projectIds {
		storeId, errt := algo.To16Bytes(projectId)
		if errt != nil {
			err = errt
			return
		}
		storeIds = append(storeIds, storeId)
	}
	pipe := pc.Pipe([]bson.M{{"$match": bson.M{"projectid": bson.M{"$in": storeIds}}},
		{"$lookup": bson.M{"from": "build", "localField": "latestbuildid", "foreignField": "_id", "as": "latestbuild"}},
		{"$out": "latestbuild"}})
	var out []model.Builds
	if err = pipe.All(&out); err != nil {
		return
	}
	for _, builds := range out {
		for _, build := range builds {
			latestBuilds = append(latestBuilds, build)
		}
	}
	return
}

func (md *ProjectMongoDao) GetBuilds(buildIds []string) (builds model.Builds, err error) {
	session := md.GSession.Copy()
	defer session.Close()
	rc := session.DB(projectDbName).C(buildCol)
	err = rc.Find(bson.M{"_id": bson.M{"$in": buildIds}}).All(&builds)
	return
}

func (md *ProjectMongoDao) GetBuildsByFilter(filter map[string]string, limitCount, offset int) (builds model.Builds, err error) {
	session := md.GSession.Copy()
	defer session.Close()
	bc := session.DB(projectDbName).C(buildCol)
	err = bc.Find(filter).Sort("-numero").Skip(offset).Limit(limitCount).All(&builds)
	return
}

func (md *ProjectMongoDao) UpdateRunLog(buildId, runId string, logId string) (err error) {
	session := md.GSession.Copy()
	defer session.Close()
	bc := session.DB(projectDbName).C(buildCol)
	err = bc.Update(bson.M{"_id": buildId, "runs._id": runId}, bson.M{"$set": bson.M{"runs.$.logid": logId}})
	return
}

func (md *ProjectMongoDao) AddLog(logs []byte) (id string, err error) {
	session := md.GSession.Copy()
	defer session.Close()
	lc := session.DB(projectDbName).C(logCol)
	oId := bson.NewObjectId()
	id = oId.Hex()
	err = lc.Insert(&model.RunLog{Id: id, Content: string(logs)})
	return
}

func (md *ProjectMongoDao) GetLog(id string) (log model.RunLog, err error) {
	session := md.GSession.Copy()
	defer session.Close()
	lc := session.DB(projectDbName).C(logCol)
	err = lc.FindId(id).One(&log)
	return
}

func (md *ProjectMongoDao) AddHook(hook model.GithubHook) (err error) {
	session := md.GSession.Copy()
	defer session.Close()
	hc := session.DB(projectDbName).C(hookCol)
	_, err = hc.UpsertId(hook.ProjectId, &hook)
	if err != nil {
		glog.Error("err when Create/add Hook because of :", err)
	}
	return
}

func (md *ProjectMongoDao) GetHook(projectId string) (hook model.GithubHook, err error) {
	session := md.GSession.Copy()
	defer session.Close()
	hc := session.DB(projectDbName).C(hookCol)
	err = hc.FindId(projectId).One(&hook)
	return
}

func (md *ProjectMongoDao) DeleteHook(projectId string) (err error) {
	session := md.GSession.Copy()
	defer session.Close()
	hc := session.DB(projectDbName).C(hookCol)
	err = hc.RemoveId(projectId)
	return
}
