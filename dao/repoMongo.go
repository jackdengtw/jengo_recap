package dao

import (
	"errors"
	"fmt"

	"github.com/golang/glog"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/qetuantuan/jengo_recap/algo"
	"github.com/qetuantuan/jengo_recap/model"
)

type RepoMongoDao struct {
	MongoDao
}

func (md *RepoMongoDao) Init(d *MongoDao) (err error) {
	if d == nil {
		err = md.MongoDao.Init()
	} else {
		md.MongoDao = *d
		if !d.Inited {
			err = d.Init()
		}
	}
	return err
}

func (md *RepoMongoDao) UpdateRepos(repos []model.Repo, userId string) (err error) {
	session := md.GSession.Copy()
	defer session.Close()
	pc := session.DB(repoDbName).C(repoCol)
	for _, repo := range repos {
		storeId, tmpErr := algo.To16Bytes(repo.RepoMeta.Id)
		if tmpErr != nil {
			err = tmpErr
			glog.Warningf("Update Repo failed! get hash id failed![%v]", err)
			break
		}
		err = pc.UpdateId(storeId, bson.M{"$set": bson.M{"repo.RepoMeta": repo.RepoMeta}, "$addToSet": bson.M{"Repo.users": userId}})
		if err != nil {
			glog.Errorf("Update failed for %v %v %v", repo.RepoMeta.Id, string(storeId), err)
			break
		}
		glog.Errorf("Update Success for %v %v", repo.RepoMeta.Id, string(storeId))
	}
	return
}

// todo: maybe can merge InsertRepos&UpdateRepos to UpsertRepos function
func (md *RepoMongoDao) InsertRepos(Repos []model.Repo, userId string) (err error) {
	session := md.GSession.Copy()
	defer session.Close()
	pc := session.DB(repoDbName).C(repoCol)
	for _, repo := range Repos {
		storeId, tmpErr := algo.To16Bytes(repo.RepoMeta.Id)
		if tmpErr != nil {
			err = tmpErr
			glog.Warningf("insert Repo failed! get hash id failed![%v]", err)
			break
		}
		_, err = pc.UpsertId(storeId, bson.M{"$set": bson.M{"repo.RepoMeta": repo.RepoMeta}, "$addToSet": bson.M{"Repo.users": userId}})
		if err != nil {
			glog.Warningf("insert Repo failed for %v %v %v", repo.RepoMeta.Id, string(storeId), err)
			break
		}
		glog.Infof("insert Repo success for %v %v", repo.RepoMeta.Id, string(storeId))
	}
	return
}

func (md *RepoMongoDao) DeleteRepos(Repos []model.Repo, userId string) (err error) {
	session := md.GSession.Copy()
	defer session.Close()
	pc := session.DB(repoDbName).C(repoCol)
	for _, repo := range Repos {
		storeId, tmpErr := algo.To16Bytes(repo.RepoMeta.Id)
		if tmpErr != nil {
			err = tmpErr
			glog.Warningf("Delete Repo failed! get hash id failed![%v]", err)
			break
		}
		//todo: use removeall
		//tmpErr := pc.RemoveId(storeId)
		err = pc.UpdateId(storeId, bson.M{"$pull": bson.M{"Repo.users": userId}})
		if err != nil {
			glog.Warningf("Delete Repo failed for %v %v %v", repo.RepoMeta.Id, string(storeId), err)
		} else {
			glog.Infof("Delete Repo success for %v %v", repo.RepoMeta.Id, string(storeId))
		}
	}
	return
}

func (md *RepoMongoDao) GetReposByFilter(filter map[string]interface{}, limitCount, offset int) (Repos []model.Repo, err error) {
	session := md.GSession.Copy()
	defer session.Close()
	pc := session.DB(repoDbName).C(repoCol)
	RepoFilter := bson.M{}
	for key, value := range filter {
		filterKey := "Repo." + key
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
	fmt.Println(RepoFilter)
	err = pc.Find(&RepoFilter).
		Sort("-repo.RepoMeta.createdat").Skip(offset).Limit(limitCount).All(&Repos)
	return
}

func (md *RepoMongoDao) GetReposByScms(userId string, scms []string) (Repos []model.Repo, err error) {
	session := md.GSession.Copy()
	defer session.Close()
	pc := session.DB(repoDbName).C(repoCol)
	err = pc.Find(&bson.M{"Repo.users": userId, "repo.RepoMeta.scm": &bson.M{"$in": scms}}).
		Sort("-repo.RepoMeta.createdat").All(&Repos)
	return
}

func (md *RepoMongoDao) GetRepos(userId string, limitCount, offset int) (Repos []model.Repo, err error) {
	session := md.GSession.Copy()
	defer session.Close()
	pc := session.DB(repoDbName).C(repoCol)
	err = pc.Find(bson.M{"Repo.users": userId}).Sort("-repo.RepoMeta.createdat").Skip(offset).Limit(limitCount).All(&Repos)
	return
}

func (md *RepoMongoDao) GetBuildIndex(RepoId string) (idx int, err error) {
	session := md.GSession.Copy()
	defer session.Close()
	pc := session.DB(repoDbName).C(repoCol)
	var storeId []byte
	if storeId, err = algo.To16Bytes(RepoId); err != nil {
		return
	}
	var p model.Repo
	change := mgo.Change{
		Update:    bson.M{"$inc": bson.M{"Repo.runindex": 1}},
		ReturnNew: true,
	}

	_, err = pc.Find(bson.M{"_id": storeId}).Apply(change, &p)
	idx = p.BuildIndex
	return
}

func (md *RepoMongoDao) GetRepo(id string) (Repo model.Repo, err error) {
	session := md.GSession.Copy()
	defer session.Close()
	pc := session.DB(repoDbName).C(repoCol)
	var storeId []byte
	if storeId, err = algo.To16Bytes(id); err != nil {
		return
	}
	var p model.Repo
	err = pc.FindId(storeId).One(&p)
	return p, err
}

func (md *RepoMongoDao) SwitchRepo(RepoId string, enableStatus bool) (err error) {
	session := md.GSession.Copy()
	defer session.Close()
	pc := session.DB(repoDbName).C(repoCol)
	var storeId []byte
	if storeId, err = algo.To16Bytes(RepoId); err != nil {
		return
	}
	err = pc.UpdateId(storeId, bson.M{"$set": bson.M{"Repo.enable": enableStatus}})
	return
}

/*
	Update Repo dynamic info. including: state latestBuildId, branch
*/
func (md *RepoMongoDao) UpdateDynamicRepoInfo(RepoId, state, latestBuildId, branch string) (err error) {
	session := md.GSession.Copy()
	defer session.Close()
	pc := session.DB(repoDbName).C(repoCol)
	var storeId []byte
	if storeId, err = algo.To16Bytes(RepoId); err != nil {
		return
	}
	updateMap := bson.M{}

	setMap := bson.M{}
	if state != "" {
		setMap["Repo.state"] = state
	}
	if latestBuildId != "" {
		setMap["Repo.latestbuildid"] = latestBuildId
	}
	if len(setMap) != 0 {
		updateMap["$set"] = setMap
	}

	addToSetMap := bson.M{}
	if branch != "" {
		addToSetMap["Repo.branches"] = branch
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
