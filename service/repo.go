package service

import (
	"fmt"

	"github.com/golang/glog"
	"github.com/qetuantuan/jengo_recap/vo"
	"github.com/qetuantuan/jengo_recap/dao"
	"github.com/qetuantuan/jengo_recap/model"
	"github.com/qetuantuan/jengo_recap/scm"
	"github.com/qetuantuan/jengo_recap/util"
	"gopkg.in/mgo.v2"
)

type RepoReader interface {
	GetRepo(RepoId string) (Repo model.Repo, err error)
	GetReposByFilter(filter map[string]interface{}, limitCount, offset int) (Repos []vo.Repo, err error)
	//GetRepos(userId string) (Repo *dao.Repo, err error)
}

type RepoWriter interface {
	UpdateRepos(userId string) (Repos []model.Repo, err error)
	SwitchRepo(userId, repoId string, enableStatus bool) (err error)
}

type RepoService interface {
	RepoReader
	RepoWriter
}

type LocalRepoService struct {
	Md             dao.RepoDao
	BuildMd        dao.SemanticBuildDao
	GithubScm      *scm.GithubScm
	HttpUserClient UserService
}

var _ RepoService = &LocalRepoService{}

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

func (p *LocalRepoService) GetRepo(RepoId string) (Repo model.Repo, err error) {
	var repo model.Repo
	repo, err = p.Md.GetRepo(RepoId)

	if err != nil {
		if err == mgo.ErrNotFound {
			glog.Warning("Repo not found: " + RepoId)
			err = NotFoundError
			return
		} else {
			glog.Error("get Repo failed! error: ", err)
			err = MongoError
			return
		}
	}
	return repo, nil
}

// syncRepoSet sync old Repos status to new Repos, missing ones as returning.
// n: new Repos map, whose key is Repo.Meta.Id.
// o: old Repos.
// d: delete Repos (missing ones in new Repos)
// s: Repos need update (in two sets)
// i: insert Repos (missing in old Repos)
func syncRepoSet(n map[string]*model.Repo, o map[string]*model.Repo) (d, s, i []model.Repo) {
	// assuming new Repo.Enable = False by default
	for _, oldRepo := range o {
		newRepo, ok := n[oldRepo.Id]
		if ok {
			s = append(s, *newRepo)
		} else {
			d = append(d, *oldRepo)
		}
	}
	for _, newRepo := range n {
		if _, ok := o[newRepo.Id]; !ok {
			i = append(i, *newRepo)
		}
	}
	return
}
func (p *LocalRepoService) UpdateRepos(userId string) (Repos []model.Repo, err error) {
	user, err := p.HttpUserClient.GetUser(userId)
	if err != nil {
		glog.Errorf("update Repos: get user from us failed! user_id: %v, error: %v", userId, err)
		return
	}
	auth := user.PrimaryAuth()
	p.GithubScm.User = auth.LoginName
	p.GithubScm.Token = auth.GetDecryptedToken(util.KeyCoder)
	// TODO:  if auth.AuthSource == 'github'
	var repos []model.Repo
	repos, err = p.GithubScm.GetRepoList()
	if err != nil {
		glog.Error("get Repos from github for user failed", p.GithubScm.User)
		err = ScmError
		return
	}
	glog.Infof("get %v Repos from github for user %v", len(repos), p.GithubScm.User)

	// get Repos in mongodb to update enable field
	oldRepos, tmpErr := p.Md.GetRepos(userId, 0, 0)
	if tmpErr != nil {
		glog.Errorf("Query mongo failed: %v", err)
		err = MongoError
		return
	}

	glog.Infof("get %v Repos from mongo for user %v", len(oldRepos), userId)

	newRepoMap := make(map[string]*model.Repo)
	oldRepoMap := make(map[string]*model.Repo)
	for i := range repos {
		newRepoMap[repos[i].Id] = &repos[i]
		glog.Infof("one new:%s", repos[i].Id)
	}
	for i := range oldRepos {
		oldRepoMap[oldRepos[i].Id] = &oldRepos[i]
		glog.Infof("one old:%s", oldRepos[i].Id)
	}
	glog.Infof("new:%d, old:%d", len(newRepoMap), len(oldRepoMap))
	deleteRepos, updateRepos, insertRepos := syncRepoSet(newRepoMap, oldRepoMap)
	deleteNum := len(deleteRepos)
	updateNum := len(updateRepos)
	insertNum := len(insertRepos)
	glog.Infof("Will delete %d Repos, update %d Repos, insert %d Repos", deleteNum, updateNum, insertNum)
	if deleteNum != 0 {
		glog.Infof("Deleting %v non-existed Repos...", deleteNum)
		err = p.Md.UnlinkRepos(deleteRepos, userId)
		if err != nil {
			glog.Errorf(" delete Repos from mongo failed: %v", err)
			err = MongoError
			return
		}
	}
	if updateNum != 0 {
		glog.Infof("Updating %v Repos...", updateNum)
		err = p.Md.UpsertRepoMeta(updateRepos, userId)
		if err != nil {
			glog.Errorf(" update Repos to mongo failed: %v", err)
			err = MongoError
			return
		}
	}
	if insertNum != 0 {
		glog.Infof("Inserting %v Repos...", insertNum)
		err = p.Md.UpsertRepoMeta(insertRepos, userId)
		if err != nil {
			glog.Errorf(" insert Repos to mongo failed: %v", err)
			err = MongoError
			return
		}
	}
	for _, p := range repos {
		Repos = append(Repos, p)
	}
	return
}

func (p *LocalRepoService) GetReposByFilter(filter map[string]interface{}, limitCount, offset int) (Repos []vo.Repo, err error) {
	var repos []model.Repo
	repos, err = p.Md.GetReposByFilter(filter, limitCount, offset)
	if err != nil {
		glog.Error("get Repo failed! error: ", err)
		err = MongoError
		return
	}
	for _, p := range repos {
		Repos = append(Repos, *p.ToApiObj())
	}
	// TODO: query build dao for latest build id and latest state
	//       add global cache layer to optimization later on
	return
}
func (p *LocalRepoService) SwitchRepo(userId, RepoId string, enableStatus bool) (err error) {
	user, err := p.HttpUserClient.GetUser(userId)
	if err != nil {
		glog.Warningf("Get user from us failed! user_id:%s, errpr:%v", userId, err)
		return
	}
	auth := user.PrimaryAuth()
	p.GithubScm.User = auth.LoginName
	p.GithubScm.Token = auth.GetDecryptedToken(util.KeyCoder)
	// TODO:  if auth.AuthSource == 'github'
	Repo, err := p.Md.GetRepo(RepoId)
	if err != nil {
		if err == mgo.ErrNotFound {
			glog.Warning("Repo not found: " + RepoId)
			err = NotFoundError
			return
		} else {
			glog.Error("get Repo failed! error: ", err)
			err = MongoError
			return
		}
	}
	if enableStatus {
		_, err = p.GithubScm.SetHook(*Repo.Name)
		if err != nil && err != scm.HookExistError {
			glog.Warningf("enable hook failed! error:%v", err)
			err = ScmError
			return
		}
	} else { //do disable the hook
		hook, errTmp := p.GithubScm.GetHook(*Repo.HooksUrl) //github exist this hook or not
		if errTmp != nil {
			err = errTmp
			if errTmp != scm.HookNonExistError { // err returned but not hooknonexist
				glog.Warningf("get hook failed from github! error:%v", err)
				err = ScmError
				return
			} else { //err is HookNonExistError, github no such hook
				glog.Warning("found no such hook from github when delete hook")
				// not return since we need to update Repo status
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

	err1 := p.Md.SwitchRepo(RepoId, enableStatus)
	if err1 != nil {
		glog.Warningf("switch Repo failed! error:%v", err)
		err = MongoError
		return
	}
	err = nil //in case of dao.HookNonExistError or dao.HookExistError
	glog.Info(fmt.Sprintf("swith Repo %s to %s success", RepoId, enableStatus))
	return
}
