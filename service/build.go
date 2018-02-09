package service

import (
	"github.com/golang/glog"

	"github.com/qetuantuan/jengo_recap/dao"
	"github.com/qetuantuan/jengo_recap/model"
)

type SemanticBuildReader interface {
	GetSemanticBuildsByFilter(filter map[string]interface{}, maxCount, offset int) (sBuilds model.SemanticBuilds, err error)
	GetSemanticBuildsByIds(buildIds []string) (builds model.SemanticBuilds, err error)
}

type SemanticBuildWriter interface {
	PartialUpdateBuild(sbuildId, buildId, repoId, state string, updateKv map[string]interface{}) (err error)
}

type BuildInserter interface {
	InsertBuild(Build model.Build) (build model.SemanticBuild, err error)
}

type BuildService interface {
	SemanticBuildReader
	SemanticBuildWriter
	BuildInserter
}

type LocalBuildService struct {
	Md dao.SemanticBuildDao
	// RepoMd *dao.RepoMongoDao
}

var _ BuildService = &LocalBuildService{}

// var _ BuildService = &client.HttpBuildService{}

func (r *LocalBuildService) GetSemanticBuildsByFilter(filter map[string]interface{}, maxCount, offset int) (builds model.SemanticBuilds, err error) {
	builds, err = r.Md.GetSemanticBuildsByFilter(filter, maxCount, offset)
	if err != nil {
		glog.Error("get build failed! error: ", err)
		err = MongoError
		return
	}
	return
}

func (r *LocalBuildService) GetSemanticBuildsByIds(buildIds []string) (builds model.SemanticBuilds, err error) {
	builds, err = r.Md.GetSemanticBuilds(buildIds)
	if err != nil {
		glog.Error("get build failed! error: ", err)
		err = MongoError
		return
	}
	return
}

func (r *LocalBuildService) PartialUpdateBuild(
	sbuildId, buildId, repoId, state string, updateKv map[string]interface{}) (err error) {
	err = r.Md.UpdateBuildProperties(sbuildId, buildId, updateKv)
	if err != nil {
		glog.Warningf("partial update Build failed! error:%v", err)
		err = MongoError
		return
	}

	// TODO: update a global cache here
	// 		cache is initialzied from db query when key missed
	/*
		if updateKv["status"] == constant.BUILD_FINISHED {
			err = r.RepoMd.UpdateDynamicRepoInfo(repoId, state, "", "")
			if err != nil {
				glog.Warningf("update Repo:%s state to %s failed! error:%v", repoId, state, err)
				err = MongoError
				return
			}
			glog.Infof("update Repo:%s state to %s success!", repoId, state)
		}
	*/
	return
}

/*
func (r *LocalBuildService) UpdateBuild(buildId string, Build model.Build) (err error) {
	mr := model.NewBuildFrom(&Build)
	err = r.Md.UpdateBuild(buildId, *mr)
	if err != nil {
		glog.Warningf("update Build failed! error:%v", err)
		err = MongoError
		return
	}

	if Build.Status == constant.Build_FINISHED {
		err = r.Md.UpdateDynamicRepoInfo(Build.RepoId, Build.State, "", "")
		if err != nil {
			glog.Warningf("update Repo:%s state to %s failed! error:%v", Build.RepoId, Build.State, err)
			err = MongoError
			return
		}
		glog.Infof("update Repo:%s state to %s success!", Build.RepoId, Build.State)
	}
	return
}
*/

func (r *LocalBuildService) InsertBuild(build model.Build) (sbuild model.SemanticBuild, err error) {
	commitId := *build.HeadCommit.ID
	branch := build.Branch
	var b model.SemanticBuild
	b, err = r.Md.FindSemanticBuildByBranchCommit(build.RepoId, commitId, branch)
	if err == dao.ErrorBuildNotFound {
		var buildNo = build.Index
		// remove this logic
		//       1. TODO: remove the concept of SemanticBuild and its build index
		// 		    use builds group by commit and branch instead.
		// 		 3. build group id can be the first build id of this group
		// 		    there is chance group id is not the first build id in compete condition. but should be okay
		/*
			buildNo, err = r.RepoMd.GetBuildIndex(build.RepoId)
			if err != nil {
				glog.Warningf("get build index from Repo when insert build failed!err:%v", err)
				err = MongoError
				return
			}
		*/
		b = model.SemanticBuild{
			RepoId:   build.RepoId,
			UserId:   build.UserId,
			Numero:   &buildNo,
			Builds:   []model.Build{model.Build(build)},
			CommitId: commitId,
			Branch:   branch,
		}

		var buildId = ""
		buildId, err = r.Md.CreateSemanticBuild(b)
		if err != nil {
			glog.Warningf("insert new build to mongodb failed! error:%v", err)
			err = MongoError
			return
		}
		b.Id = buildId

		// TODO: update a global cache here
		/*
			state := ""
			if build.Status == constant.BUILD_FINISHED {
				state = build.Result
			}
				err = r.RepoMd.UpdateDynamicRepoInfo(b.RepoId, state, buildId, branch)
				if err != nil {
					glog.Warningf("update Repo:%s latest build id :%s failed! error:%v", b.RepoId, buildId, err)
					err = MongoError
					return
				}
		*/
		return
	}

	if err != nil {
		glog.Warningf("find build by commit from mongodb failed! error:%v", err)
		err = MongoError
		return
	}

	err = r.Md.InsertBuild(b.Id, build)
	if err != nil {
		glog.Warningf("insert Build to mongodb failed! error:%v", err)
		err = MongoError
		return
	}

	// TODO: update a global cache here
	/*
		if build.Status == constant.BUILD_FINISHED {
			err = r.RepoMd.UpdateDynamicRepoInfo(build.RepoId, build.Result, "", build.Branch)
			if err != nil {
				glog.Warningf("insert Repo:%s set state to %s failed! error:%v", build.RepoId, build.Result, err)
				err = MongoError
				return
			}
			glog.Infof("update Repo:%s state to %s success!", build.RepoId, build.Result)
		}
	*/
	return
}
