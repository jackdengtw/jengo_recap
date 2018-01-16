package service

import (
	"github.com/golang/glog"
	"github.com/qetuantuan/jengo_recap/api"
	"github.com/qetuantuan/jengo_recap/dao"
	"github.com/qetuantuan/jengo_recap/model"
	"github.com/qetuantuan/jengo_recap/scm"
	"github.com/qetuantuan/jengo_recap/util"
	"gopkg.in/mgo.v2"
)

type UserReader interface {
	GetUser(userId string) (model.User, error)
	GetUserByLogin(loginName string, auth string) (model.User, error)
}

type UserWriter interface {
	CreateUser(loginName string, auth string, token string) (userId string, err error)
	UpdateScmToken(userId, scmId, tokenStr string) error
}

type UserService interface {
	UserReader
	UserWriter
}

type LocalUserService struct {
	Md        dao.UserDao
	GithubScm scm.Scm
}

var _ UserService = &LocalUserService{}

/*
 * Service layer encapsulates biz logic, leaving only HTTP protocol input/output to handler layer
 *
 * It makes possible
 * 1. functional test for service method in UT form as regression test.
 * 2. RESTful api protocol test in handler level with mocked services.
 *
 * It requires
 * 1. An error translation
 * 2. service logging what happened internals
 */

func (u *LocalUserService) CreateUser(loginName string, auth string, token string) (userId string, err error) {
	_, err = u.Md.GetUserByLogin(loginName, auth)
	if err == nil {
		glog.Warningf("user EXIST ALREADY: %v@%v", loginName, auth)
		err = CreateConflictError
		return
	} else if err != mgo.ErrNotFound {
		// TODO: may have selective retry here
		glog.Errorf("Query mongo failed: %v", err)
		err = MongoError
		return
	} else {
		if auth == api.AUTH_SOURCE_GITHUB {
			userScm, err1 := u.GithubScm.GetUser(token)
			if err1 != nil {
				glog.Errorf("failed, err when get user from scm", loginName, err1)
				err = ScmError
				return
			}
			glog.Info("Got user from scm is ", userScm.Id)
			user := userScm.ToMongoUser()
			err = user.SetTokenEncrypted(api.AUTH_SOURCE_GITHUB, util.KeyCoder, token)
			if err != nil {
				glog.Errorf("Encrypt token error: %v", err)
				err = EncryptError
				return
			}

			err = u.Md.UpsertUser(user)
			if err != nil {
				glog.Errorf("failed, err when insert to Mongo", loginName, err)
				err = MongoError
				return
			}
			userId = user.Id
			glog.Infof("insert user %s to db success", userId)
		} else {
			glog.Warningf("Not supported auth: %v", auth)
			err = NotSupportedAuthError
			return
		}
	}
	return
}

func (u *LocalUserService) GetUser(userId string) (user model.User, err error) {
	user, err = u.Md.GetUser(userId)
	if err != nil {
		if err == mgo.ErrNotFound {
			glog.Warning("User not found: " + userId)
			err = UserNotFoundError
			return
		} else {
			glog.Error("get user failed! error: ", err)
			err = MongoError
			return
		}
	}
	glog.Info("user found")
	return
}

func (u *LocalUserService) GetUserByLogin(loginName string, auth string) (user model.User, err error) {
	user, err = u.Md.GetUserByLogin(loginName, auth)
	if err != nil {
		if err == mgo.ErrNotFound {
			glog.Warningf("User %v@%v not found.", loginName, auth)
			err = UserNotFoundError
			return
		} else {
			glog.Errorf("failed to get user: %", err)
			err = MongoError
			return
		}
	}

	glog.Info("get user by login success!")
	return
}

func (u *LocalUserService) UpdateScmToken(userId, scmId, tokenStr string) (err error) {
	var user model.User
	user, err = u.Md.GetUser(userId)
	if err == mgo.ErrNotFound {
		glog.Errorf("userId %v not found!", userId)
		err = UserNotFoundError
		return
	} else if err != nil {
		// TODO: may have selective retry here
		glog.Errorf("Query mongo failed: %v", err)
		err = MongoError
		return
	} else {
		var scm *model.Scm
		for _, s := range user.Scms {
			if s.Id == scmId {
				scm = &s
				break
			}
		}
		if scm == nil {
			glog.Errorf("scm Id %v not found!", scmId)
			err = ScmNotFoundError
			return
		}
		var token []byte
		if token, err = util.AESEncode([]byte(util.KeyCoder), []byte(tokenStr)); err != nil {
			glog.Errorf("encypt token failed: %v", err)
			err = EncryptError
			return
		} else if err = u.Md.UpdateToken(userId, dao.SCM_COLUMN, scmId, token); err != nil {
			glog.Errorf("update token in DB error: %v", err)
			err = MongoError
			return
		} else {
			glog.Info("update token success!")
			return
		}
	}
}
