package dao

import (
	"fmt"

	"github.com/golang/glog"
	"github.com/qetuantuan/jengo_recap/model"
	"gopkg.in/mgo.v2/bson"
)

type SemanticBuildReader interface {
	FindSemanticBuildByBranchCommit(repoId, commitId, branch string) (sbuild model.SemanticBuild, err error)
	GetSemanticBuilds(sbuildIds []string) (sbuilds model.SemanticBuilds, err error)
	GetSemanticBuildsByFilter(filter map[string]interface{}, limitCount, offset int) (sbuilds model.SemanticBuilds, err error)
	IsBuildExistInSemanticBuild(buildId, sBuildId string) (res bool, err error)
}

type SemanticBuildWriter interface {
	CreateSemanticBuild(b model.SemanticBuild) (id string, err error)
}

type BuildWriter interface {
	InsertBuild(sbuildId string, build model.Build) (err error)
	UpdateBuildProperties(sBuildId string, buildId string, p map[string]interface{}) (err error)
	UpdateBuildLog(sbuildId, buildId string, logId string) (err error)
}

type SemanticBuildDao interface {
	SemanticBuildReader
	SemanticBuildWriter
	BuildWriter
}

type SemanticBuildMongoDao struct {
	MongoDao
}

var _ SemanticBuildDao = &SemanticBuildMongoDao{}

func (md *SemanticBuildMongoDao) Init(d *MongoDao) (err error) {
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

func (md *SemanticBuildMongoDao) CreateSemanticBuild(b model.SemanticBuild) (id string, err error) {
	session := md.GSession.Copy()
	defer session.Close()
	bc := session.DB(repoDbName).C(buildCol)
	// TODO: Id is hash value from repoId, branch and commitId
	// to save a round trip to Mongo
	var one model.SemanticBuild
	if one, err = md.FindSemanticBuildByBranchCommit(b.RepoId, b.CommitId, b.Branch); err != nil {
		if err == ErrorBuildNotFound {
			oid := bson.NewObjectId()
			id = oid.Hex()
			b.Id = id
			err = bc.Insert(b)
		}
	} else {
		id = one.Id
		err = ErrorAlreadyExisted
		return
	}
	return
}

func (md *SemanticBuildMongoDao) FindSemanticBuildByBranchCommit(
	repoId, commitId, branch string) (sbuild model.SemanticBuild, err error) {

	session := md.GSession.Copy()
	defer session.Close()
	bc := session.DB(repoDbName).C(buildCol)
	sbuilds := model.SemanticBuilds{}
	err = bc.Find(bson.M{"repoid": repoId, "commitid": commitId, "branch": branch}).All(&sbuilds)
	if err != nil {
		return
	}
	if len(sbuilds) <= 0 {
		err = ErrorBuildNotFound
		return
	} else if len(sbuilds) > 1 {
		glog.Errorf("More than one build found for one id: %v", sbuilds)
		err = ErrorMoreThanOneBuildExisted
		return
	}
	sbuild = sbuilds[0]
	return
}

func (md *SemanticBuildMongoDao) GetSemanticBuilds(sbuildIds []string) (sbuilds model.SemanticBuilds, err error) {
	session := md.GSession.Copy()
	defer session.Close()
	rc := session.DB(repoDbName).C(buildCol)
	err = rc.Find(bson.M{"_id": bson.M{"$in": sbuildIds}}).All(&sbuilds)
	return
}

func (md *SemanticBuildMongoDao) GetSemanticBuildsByFilter(
	filter map[string]interface{}, limitCount, offset int) (sbuilds model.SemanticBuilds, err error) {

	session := md.GSession.Copy()
	defer session.Close()
	bc := session.DB(repoDbName).C(buildCol)

	// TODO: valid filters with Sbuild definition
	err = bc.Find(filter).Sort("-numero").Skip(offset).Limit(limitCount).All(&sbuilds)
	return
}

// latest build presented by last build of Sbuild.Builds
/*
func (md *SemanticBuildMongoDao) GetLatestBuild(sBuildIds []string) (latestBuilds model.SemanticBuilds, err error) {
	session := md.GSession.Copy()
	defer session.Close()
	pc := session.DB(repoDbName).C(semanticBuildCol)
	pipe := pc.Pipe([]bson.M{{"$match": bson.M{"SemanticBuildid": bson.M{"$in": sBuildIds}}},
		{"$lookup": bson.M{"from": "build", "localField": "latestbuildid", "foreignField": "_id", "as": "latestbuild"}},
		{"$out": "latestbuild"}})
	var out []model.SemanticBuilds
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
*/

func (md *SemanticBuildMongoDao) IsBuildExistInSemanticBuild(buildId, sBuildId string) (res bool, err error) {
	session := md.GSession.Copy()
	defer session.Close()
	bc := session.DB(repoDbName).C(buildCol)
	var num int
	num, err = bc.Find(bson.M{"_id": sBuildId, "builds._id": buildId}).Count()
	res = num > 0
	return
}

//
// Build
//

func (md *SemanticBuildMongoDao) InsertBuild(sbuildId string, build model.Build) (err error) {
	session := md.GSession.Copy()
	defer session.Close()
	bc := session.DB(repoDbName).C(buildCol)
	err = bc.UpdateId(sbuildId, bson.M{"$push": bson.M{"builds": build}})
	return
}

/*
func (md *SemanticBuildMongoDao) UpdateBuild(sbuildId string, Build model.Build) (err error) {
	session := md.GSession.Copy()
	defer session.Close()
	bc := session.DB(repoDbName).C(buildCol)
	err = bc.Update(bson.M{"_id": sbuildId, "Builds._id": Build.Id}, bson.M{"$set": bson.M{"Builds.$": Build}})
	return
}
*/

func (md *SemanticBuildMongoDao) UpdateBuildProperties(
	sBuildId string, buildId string, p map[string]interface{}) (err error) {

	// TODO: valid key existence/value type with Sbuild definition

	// https://docs.mongodb.com/manual/reference/operator/update/set/
	var BuildInterface = make(map[string]interface{})
	for k, v := range p {
		BuildInterface["builds.$."+k] = v
	}
	glog.Infof("BuildInterface: %v", BuildInterface)

	// db.Update
	session := md.GSession.Copy()
	defer session.Close()
	bc := session.DB(repoDbName).C(buildCol)
	err = bc.Update(bson.M{"_id": sBuildId, "builds._id": buildId}, bson.M{"$set": BuildInterface})
	if err != nil {
		msg := fmt.Sprintf("partial update p failed! error:%v", err)
		glog.Warning(msg)
		return
	}
	return
}

func (md *SemanticBuildMongoDao) UpdateBuildLog(sbuildId, buildId string, logId string) (err error) {
	session := md.GSession.Copy()
	defer session.Close()
	bc := session.DB(repoDbName).C(buildCol)

	// TODO: use UpdateBuildProperties
	err = bc.Update(bson.M{"_id": sbuildId, "builds._id": buildId}, bson.M{"$set": bson.M{"builds.$.loguri": logId}})
	return
}
