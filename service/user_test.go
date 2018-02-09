package service

import (
	"errors"
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/qetuantuan/jengo_recap/vo"
	"github.com/qetuantuan/jengo_recap/model"
	"github.com/qetuantuan/jengo_recap/scm"
	"github.com/qetuantuan/jengo_recap/util"
	"gopkg.in/mgo.v2"
)

type mockUserDao_testUser struct {
	U *model.User
}

func (u *mockUserDao_testUser) GetUser(userId string) (user model.User, err error) {
	// fmt.Println("get user: %v", u.U)
	if u.U != nil {
		user = *u.U
		return
	} else {
		err = mgo.ErrNotFound
		return
	}
}

func (u *mockUserDao_testUser) UpsertUser(user *model.User) (err error) {
	u.U = user
	// fmt.Println("inserted user: %v", u.U)
	return
}

func (u *mockUserDao_testUser) GetUserByLogin(loginName string, auth string) (user model.User, err error) {
	// fmt.Println("get by user: %v", u.U)
	if u.U != nil {
		user = *u.U
		return
	} else {
		err = mgo.ErrNotFound
		return
	}
}

func (u *mockUserDao_testUser) UpdateToken(userId, sourceType, scmId string, token []byte) (err error) {
	// fmt.Println("update token: %v", u.U)
	if u.U != nil {
		u.U.Scms[0].Token = token
		return
	} else {
		return
	}
}

type mockScm_testUser struct {
	U *scm.GithubUser
}

func (s *mockScm_testUser) GetUser(token string) (user scm.GithubUser, err error) {
	if s.U != nil {
		return *s.U, nil
	} else {
		err = errors.New("get scm user error!")
		return
	}
}

func (s *mockScm_testUser) SetHook(string) (hook model.GithubHook, err error)  { return }
func (s *mockScm_testUser) EditHook(string) (hook model.GithubHook, err error) { return }
func (s *mockScm_testUser) DeleteHook(string) error                            { return nil }
func (s *mockScm_testUser) GetRepoList() ([]model.Repo, error)                 { return nil, nil }
func (s *mockScm_testUser) SetToken(string)                                    { return }
func (s *mockScm_testUser) GetYmlContent(repo string, branch string) (content []byte, err error) {
	return nil, nil
}

var _ = Describe("Test User Service", func() {
	Context("User with specific id in Mock Dao", func() {
		var (
			service *LocalUserService
			userId  string = "123"
		)
		BeforeEach(func() {
			service = &LocalUserService{
				Md: &mockUserDao_testUser{
					U: &model.User{
						Id: userId,
						Auths: []model.Auth{
							model.Auth{
								AuthBase: vo.AuthBase{
									LoginName: "existed",
									Primary:   true,
								},
								AuthSourceId: "github.com",
							},
						},
						Scms: []model.Scm{
							model.Scm{
								ScmBase: vo.ScmBase{
									Id: "u_github_123",
								},
							},
						},
					},
				},
			}
		})
		It("Should return an conflict error when creating the same user", func() {
			_, err := service.CreateUser("existed", "github.com", "token1")
			Expect(err).To(Equal(CreateConflictError))
		})
		It("Should be able to update auth token", func() {
			err := service.UpdateScmToken(userId, "u_github_123", "NewToken")
			Expect(err).NotTo(HaveOccurred())
			u, err := service.GetUser(userId)
			Expect(u.Scms).To(HaveLen(1))
			Expect(u.Scms[0].GetDecryptedToken(util.KeyCoder)).To(Equal("NewToken"))
		})
	})

	Context("With a mock github server", func() {
		var (
			userId  string = "456"
			service *LocalUserService
		)
		BeforeEach(func() {
			scmUser := scm.GithubUser{
				Id:    456,
				Login: "myuser",
			}

			service = &LocalUserService{
				Md:        &mockUserDao_testUser{},
				GithubScm: &mockScm_testUser{U: &scmUser},
			}

		})
		It("Should get the user just created", func() {
			var err error
			userId, err = service.CreateUser("myuser", "github", "token1")
			Expect(err).NotTo(HaveOccurred())
			u, err := service.GetUser(userId)
			Expect(err).NotTo(HaveOccurred())
			Expect(u.Id).NotTo(Equal(""))
			Expect(u.PrimaryAuth().OriginId).To(Equal(strconv.Itoa(456)))
			Expect(u.PrimaryAuth().AuthSourceId).To(Equal("github"))
			Expect(u.PrimaryAuth().LoginName).To(Equal("myuser"))

			// Note: usually we should test one method per case.
			//       sometimes we may take chance of create/get
			//       since we are not blind testing.
		})
	})
})
