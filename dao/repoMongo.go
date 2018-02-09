package dao

import (
	"errors"

	"github.com/golang/glog"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/qetuantuan/jengo_recap/model"
)

type RepoReader interface {
	GetRepo(id string) (Repo model.Repo, err error)
	GetReposByFilter(filter map[string]interface{}, limitCount, offset int) (Repos []model.Repo, err error)
	GetReposByScms(userId string, scms []string) (Repos []model.Repo, err error)
	GetRepos(userId string, limitCount, offset int) (Repos []model.Repo, err error)
	GetBuildIndex(id string) (idx int, err error)
}

type RepoWriter interface {
	UpsertRepoMeta(Repos []model.Repo, userId string) (err error)
	UpdateDynamicRepoInfo(id, branch string) (err error)
	SwitchRepo(id string, enableStatus bool) (err error)
	UnlinkRepos(Repos []model.Repo, userId string) (err error)
}

type RepoDao interface {
	RepoReader
	RepoWriter
}

type RepoMongoDao struct {
	MongoDao
}

var _ RepoDao = &RepoMongoDao{}

func (md *RepoMongoDao) Init(d *MongoDao) (err error) {
	if d == nil {
		err = md.MongoDao.Init()
	} else {
		md.MongoDao = *d
		if !d.Inited {
			err = md.MongoDao.Init()
		}
	}
	return err
}

func (md *RepoMongoDao) UpsertRepoMeta(Repos []model.Repo, userId string) (err error) {
	session := md.GSession.Copy()
	defer session.Close()
	pc := session.DB(repoDbName).C(repoCol)
	for i, repo := range Repos {
		_, err = pc.UpsertId(repo.Id, bson.M{"$set": bson.M{"repometa": repo.RepoMeta},
			"$addToSet": bson.M{"user_ids": userId}})
		if err != nil {
			err = BatchError{FailedIdx: i, RealErr: err}
			glog.Warningf("insert Repo failed for %v %v", repo.Id, err)
			break
		}
		glog.Infof("insert Repo success for %v", repo.Id)
	}
	return
}

func (md *RepoMongoDao) GetRepos(userId string, limitCount, offset int) (Repos []model.Repo, err error) {
	Repos, err = md.GetReposByFilter(
		map[string]interface{}{
			"user_ids": userId,
		}, limitCount, offset,
	)
	return
}

func (md *RepoMongoDao) UnlinkRepos(Repos []model.Repo, userId string) (err error) {
	session := md.GSession.Copy()
	defer session.Close()
	pc := session.DB(repoDbName).C(repoCol)
	for i, repo := range Repos {
		err = pc.UpdateId(repo.Id, bson.M{"$pull": bson.M{"user_ids": userId}})
		if err != nil {
			err = BatchError{FailedIdx: i, RealErr: err}
			glog.Warningf("Unlink Repo failed for %v %v", repo.Id, err)
		}
		glog.Infof("Unlink Repo success for %v", repo.Id)
	}
	return
}

func (md *RepoMongoDao) GetReposByFilter(filter map[string]interface{}, limitCount, offset int) (Repos []model.Repo, err error) {
	session := md.GSession.Copy()
	defer session.Close()
	pc := session.DB(repoDbName).C(repoCol)
	RepoFilter := bson.M{}
	for key, value := range filter {
		filterKey := key
		switch value.(type) {
		case bool:
			if value.(bool) {
				RepoFilter[filterKey] = value
			} else {
				RepoFilter[filterKey] = bson.M{"$ne": true}
			}
		case []string:
			RepoFilter[filterKey] = bson.M{"$in": value}
		default:
			RepoFilter[filterKey] = value
		}
	}
	glog.Infof("Filter of get repos is: %v", RepoFilter)
	err = pc.Find(&RepoFilter).
		Sort("-repometa.createdat").Skip(offset).Limit(limitCount).All(&Repos)
	return
}

func (md *RepoMongoDao) GetReposByScms(userId string, scms []string) (Repos []model.Repo, err error) {
	Repos, err = md.GetReposByFilter(
		map[string]interface{}{
			"user_ids":         userId,
			"repometa.scmname": scms,
		}, 0, 0,
	)
	return
}

func (md *RepoMongoDao) GetBuildIndex(id string) (idx int, err error) {
	session := md.GSession.Copy()
	defer session.Close()
	pc := session.DB(repoDbName).C(repoCol)
	var p model.Repo
	change := mgo.Change{
		Update:    bson.M{"$inc": bson.M{"buildindex": 1}},
		ReturnNew: true,
	}

	_, err = pc.Find(bson.M{"_id": id}).Apply(change, &p)
	idx = p.BuildIndex
	return
}

func (md *RepoMongoDao) GetRepo(id string) (Repo model.Repo, err error) {
	session := md.GSession.Copy()
	defer session.Close()
	pc := session.DB(repoDbName).C(repoCol)
	var p model.Repo
	err = pc.FindId(id).One(&p)
	return p, err
}

func (md *RepoMongoDao) SwitchRepo(id string, enableStatus bool) (err error) {
	session := md.GSession.Copy()
	defer session.Close()
	pc := session.DB(repoDbName).C(repoCol)
	err = pc.UpdateId(id, bson.M{"$set": bson.M{"enabled": enableStatus}})
	return
}

/*
	Update Repo dynamic info. including: state latestBuildId, branch
*/
func (md *RepoMongoDao) UpdateDynamicRepoInfo(id, branch string) (err error) {
	session := md.GSession.Copy()
	defer session.Close()
	pc := session.DB(repoDbName).C(repoCol)
	updateMap := bson.M{}

	/*
		setMap := bson.M{}
		if state != "" {
			setMap["state"] = state
		}
		if latestBuildId != "" {
			setMap["latestbuildid"] = latestBuildId
		}
		if len(setMap) != 0 {
			updateMap["$set"] = setMap
		}

	*/
	addToSetMap := bson.M{}
	if branch != "" {
		addToSetMap["branches"] = branch
	}

	if len(addToSetMap) != 0 {
		updateMap["$addToSet"] = addToSetMap
	}
	if len(updateMap) == 0 {
		err = errors.New("nothing to update")
		return
	}
	err = pc.UpdateId(id, updateMap)
	return
}
