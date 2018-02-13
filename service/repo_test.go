package service

import (
	"github.com/qetuantuan/jengo_recap/util"
	"github.com/qetuantuan/jengo_recap/vo"
	// "sort"
	// "testing"

	"fmt"
	"sort"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/qetuantuan/jengo_recap/model"
	"github.com/qetuantuan/jengo_recap/scm"
)

type mockBuildDao_testRepo struct {
	Sb *model.SemanticBuild // keep one in mem is enough for now
}

func (m *mockBuildDao_testRepo) FindSemanticBuildByBranchCommit(repoId, commitId, branch string) (sbuild model.SemanticBuild, err error) {
	return
}

func (m *mockBuildDao_testRepo) GetSemanticBuildsByRepoIds(repoIds []string) (sbuilds model.SemanticBuilds, err error) {
	sbuilds = model.SemanticBuilds{
		model.SemanticBuild{
			Id:     "buildId1",
			RepoId: "repoId1",
		},
	}
	return
}

func (m *mockBuildDao_testRepo) GetSemanticBuilds(sbuildIds []string) (sbuilds model.SemanticBuilds, err error) {
	return
}

func (m *mockBuildDao_testRepo) GetSemanticBuildsByFilter(filter map[string]interface{}, limitCount, offset int) (sbuilds model.SemanticBuilds, err error) {
	return
}

func (m *mockBuildDao_testRepo) IsBuildExistInSemanticBuild(buildId, sBuildId string) (res bool, err error) {
	return
}

type mockRepoDao_testRepo struct {
	Rs      []model.Repo
	repo1on bool
}

func (m *mockRepoDao_testRepo) GetRepo(id string) (Repo model.Repo, err error) {
	Repo.Id = "repoId1"
	s := "repoName1"
	Repo.Name = &s
	s1 := "hooksUrl1"
	Repo.HooksUrl = &s1
	return
}

func (m *mockRepoDao_testRepo) GetReposByFilter(filter map[string]interface{}, limitCount, offset int) (Repos []model.Repo, err error) {
	Repos = m.Rs
	return
}

func (m *mockRepoDao_testRepo) GetReposByScms(userId string, scms []string) (Repos []model.Repo, err error) {
	Repos = m.Rs
	return
}

func (m *mockRepoDao_testRepo) GetRepos(userId string, limitCount, offset int) (Repos []model.Repo, err error) {
	return
}
func (m *mockRepoDao_testRepo) GetBuildIndex(id string) (idx int, err error) { return }

func (m *mockRepoDao_testRepo) UpsertRepoMeta(Repos []model.Repo, userId string) (err error) {
	m.Rs = Repos
	for i := range m.Rs {
		m.Rs[i].OwnerIds = []string{userId}
	}
	return
}

func (m *mockRepoDao_testRepo) UpdateDynamicRepoInfo(id, branch string) (err error) { return }
func (m *mockRepoDao_testRepo) SwitchRepo(id string, enableStatus bool) (err error) {
	m.repo1on = enableStatus
	return
}

func (m *mockRepoDao_testRepo) UnlinkRepos(Repos []model.Repo, userId string) (err error) { return }

type mockScm_testRepo struct {
}

func (m *mockScm_testRepo) GetUser() (user scm.GithubUser, err error)             { return }
func (m *mockScm_testRepo) SetHook(string) (hook model.GithubHook, err error)     { return }
func (m *mockScm_testRepo) EditHook(string) (hook model.GithubHook, err error)    { return }
func (m *mockScm_testRepo) GetHook(url string) (hook model.GithubHook, err error) { return }
func (m *mockScm_testRepo) DeleteHook(string) (err error)                         { return }
func (m *mockScm_testRepo) GetRepoList() (repos []model.Repo, err error) {
	repos = []model.Repo{
		model.Repo{
			Id: "repoId1",
		},
		model.Repo{
			Id: "repoId2",
		},
	}
	return
}
func (m *mockScm_testRepo) SetToken(string)             { return }
func (m *mockScm_testRepo) GetToken() (token string)    { return }
func (m *mockScm_testRepo) SetUserName(string)          { return }
func (m *mockScm_testRepo) GetUserName() (token string) { return }
func (m *mockScm_testRepo) GetYmlContent(repo string, branch string) (content []byte, err error) {
	return
}

type mockUserService_testRepo struct {
}

func (m *mockUserService_testRepo) GetUser(userId string) (user model.User, err error) {
	auth := model.Auth{
		AuthBase: vo.AuthBase{
			Id:        "authId1",
			Primary:   true,
			LoginName: "user1",
		},
	}
	user.Auths = append(user.Auths, auth)
	user.SetTokenEncrypted("authId1", util.KeyCoder, "token1")
	return
}

func (m *mockUserService_testRepo) GetUserByLogin(loginName string, auth string) (user model.User, err error) {
	return
}

var _ = Describe("Test Repo Service", func() {
	Describe("Test SyncRepoSet", func() {
		tdata := []repoTestData{
			repoTestData{
				N: map[string]*model.Repo{
					"e1": &model.Repo{Id: "e1"},
					"d1": &model.Repo{Id: "d1"},
					"d3": &model.Repo{Id: "d3"},
				},
				O: map[string]*model.Repo{
					"e1": &model.Repo{Id: "e1"},
					"d1": &model.Repo{Id: "d1"},
					"e2": &model.Repo{Id: "e2"},
					"d2": &model.Repo{Id: "d2"},
				},
				ExpectD: []model.Repo{
					model.Repo{Id: "e2"},
					model.Repo{Id: "d2"},
				},
				ExpectU: []model.Repo{
					model.Repo{Id: "e1"},
					model.Repo{Id: "d1"},
				},
				ExpectI: []model.Repo{
					model.Repo{Id: "d3"},
				},
			},
		}

		for _, info := range tdata {
			d, u, i := syncRepoSet(info.N, info.O)
			// Not doing an ID comparison to save a loop search or a map var
			// Using naming convension that "e*" means enabled while "d*" vice verse
			if !CompareRepoSet(d, info.ExpectD) {
				fmt.Errorf("project to be deleted not as expected.\n d: %v\n expect: %v", d, info.ExpectD)
				Expect(true).Should(BeFalse())
			}
			if !CompareRepoSet(u, info.ExpectU) {
				fmt.Errorf("project to be updated not as expected.\n d: %v\n expect: %v", u, info.ExpectU)
				Expect(true).Should(BeFalse())
			}
			if !CompareRepoSet(i, info.ExpectI) {
				fmt.Errorf("project to be insert not as expected.\n d: %v\n expect: %v", i, info.ExpectI)
				Expect(true).Should(BeFalse())
			}
		}

	})
	Describe("Get two Repos from Scm", func() {
		var (
			service   *LocalRepoService
			err       error
			userId1   string = "userId1"
			initRepos []model.Repo
			md        *mockRepoDao_testRepo = &mockRepoDao_testRepo{}
		)
		BeforeEach(func() {
			service = &LocalRepoService{
				Md:      md,
				BuildMd: &mockBuildDao_testRepo{},
				Scm:     &mockScm_testRepo{},
				Us:      &mockUserService_testRepo{},
			}
			initRepos, err = service.UpdateRepos(userId1)
		})
		It("Should return success", func() {
			Expect(err).NotTo(HaveOccurred())
			Expect(len(initRepos)).To(Equal(2))
		})
		It("Should be able to get the repos", func() {
			repos, err := service.GetReposByFilter(map[string]interface{}{
				"owner_ids": []string{userId1},
			}, 0, 0)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(repos)).To(Equal(2))
			Expect(repos[0].OwnerIds).To(ContainElement("userId1"))
		})
		It("Should be able to get one repo", func() {
			repo, err := service.GetRepo("repoId1")
			Expect(err).NotTo(HaveOccurred())
			Expect(repo.Id).To(Equal("repoId1"))
		})
		It("Should be able to enable one repo", func() {
			err = service.SwitchRepo("userId1", "repoId1", true)
			Expect(err).NotTo(HaveOccurred())
			Expect(md.repo1on).Should(BeTrue())
		})
		It("Should be able to disable one repo", func() {
			err = service.SwitchRepo("userId1", "repoId1", false)
			Expect(err).NotTo(HaveOccurred())
			Expect(md.repo1on).Should(BeFalse())
		})
	})
})

type repoTestData struct {
	N       map[string]*model.Repo
	O       map[string]*model.Repo
	ExpectD []model.Repo
	ExpectU []model.Repo
	ExpectI []model.Repo
}

func CompareRepoSet(actual, expected []model.Repo) bool {
	if len(expected) != len(actual) {
		return false
	}
	sort.Sort(model.ById(actual))
	sort.Sort(model.ById(expected))
	for i := range actual {
		if actual[i].Id != expected[i].Id {
			return false
		}
	}
	return true
}
