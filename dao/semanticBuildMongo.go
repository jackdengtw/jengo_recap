package dao

import (
	"fmt"

	"github.com/golang/glog"
	"github.com/qetuantuan/jengo_recap/algo"
	"github.com/qetuantuan/jengo_recap/model"
	"gopkg.in/mgo.v2/bson"
)

type SemanticBuildMongoDao struct {
	MongoDao
}

func (md *SemanticBuildMongoDao) Init(d *MongoDao) (err error) {
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

func (md *SemanticBuildMongoDao) CreateBuild(build model.Build) (id string, err error) {
	session := md.GSession.Copy()
	defer session.Close()
	bc := session.DB(repoDbName).C(buildCol)
	oid := bson.NewObjectId()
	id = oid.Hex()
	build.Id = id
	err = bc.Insert(build)
	if err != nil {
		return
	}
	return
}

func (md *SemanticBuildMongoDao) FindBuildByBranchCommit(SemanticBuildId, commitId, branch string) (build model.Build, err error) {
	session := md.GSession.Copy()
	defer session.Close()
	bc := session.DB(repoDbName).C(buildCol)
	builds := model.Builds{}
	err = bc.Find(bson.M{"SemanticBuildid": SemanticBuildId, "commitid": commitId, "branch": branch}).All(&builds)
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

func (md *SemanticBuildMongoDao) BuildExistInBuild(buildId, BuildId string) (res bool, err error) {
	session := md.GSession.Copy()
	defer session.Close()
	bc := session.DB(repoDbName).C(buildCol)
	num, err := bc.Find(bson.M{"_id": buildId, "Builds._id": BuildId}).Count()
	res = num > 0
	return

}

func (md *SemanticBuildMongoDao) UpdateBuild(buildId string, Build model.Build) (err error) {
	session := md.GSession.Copy()
	defer session.Close()
	bc := session.DB(repoDbName).C(buildCol)
	err = bc.Update(bson.M{"_id": buildId, "Builds._id": Build.Id}, bson.M{"$set": bson.M{"Builds.$": Build}})
	return
}

func (md *SemanticBuildMongoDao) UpdateBuildProperties(buildId string, BuildId string, p map[string]interface{}) (err error) {
	// transform p to map[string]interface{}
	// https://docs.mongodb.com/manual/reference/operator/update/set/
	var BuildInterface = make(map[string]interface{})
	for k, v := range p {
		BuildInterface["Builds.$."+k] = v
	}
	glog.Infof("BuildInterface: %v", BuildInterface)

	// db.Update
	session := md.GSession.Copy()
	defer session.Close()
	bc := session.DB(repoDbName).C(buildCol)
	err = bc.Update(bson.M{"_id": buildId, "Builds._id": BuildId}, bson.M{"$set": BuildInterface})
	if err != nil {
		msg := fmt.Sprintf("partial update p failed! error:%v", err)
		glog.Warning(msg)
		return
	}
	return
}

func (md *SemanticBuildMongoDao) InsertBuild(buildId string, Build model.Build) (err error) {
	session := md.GSession.Copy()
	defer session.Close()
	bc := session.DB(repoDbName).C(buildCol)
	err = bc.UpdateId(buildId, bson.M{"$push": bson.M{"Builds": Build}})
	return
}

func (md *SemanticBuildMongoDao) GetLatestBuild(SemanticBuildIds []string) (latestBuilds model.Builds, err error) {
	session := md.GSession.Copy()
	defer session.Close()
	pc := session.DB(repoDbName).C(semanticBuildCol)
	var storeIds [][]byte
	for _, SemanticBuildId := range SemanticBuildIds {
		storeId, errt := algo.To16Bytes(SemanticBuildId)
		if errt != nil {
			err = errt
			return
		}
		storeIds = append(storeIds, storeId)
	}
	pipe := pc.Pipe([]bson.M{{"$match": bson.M{"SemanticBuildid": bson.M{"$in": storeIds}}},
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

func (md *SemanticBuildMongoDao) GetBuilds(buildIds []string) (builds model.Builds, err error) {
	session := md.GSession.Copy()
	defer session.Close()
	rc := session.DB(repoDbName).C(buildCol)
	err = rc.Find(bson.M{"_id": bson.M{"$in": buildIds}}).All(&builds)
	return
}

func (md *SemanticBuildMongoDao) GetBuildsByFilter(filter map[string]string, limitCount, offset int) (builds model.Builds, err error) {
	session := md.GSession.Copy()
	defer session.Close()
	bc := session.DB(repoDbName).C(buildCol)
	err = bc.Find(filter).Sort("-numero").Skip(offset).Limit(limitCount).All(&builds)
	return
}

func (md *SemanticBuildMongoDao) UpdateBuildLog(buildId, BuildId string, logId string) (err error) {
	session := md.GSession.Copy()
	defer session.Close()
	bc := session.DB(repoDbName).C(buildCol)
	err = bc.Update(bson.M{"_id": buildId, "Builds._id": BuildId}, bson.M{"$set": bson.M{"Builds.$.logid": logId}})
	return
}
