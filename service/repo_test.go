package service

import (
	// "sort"
	// "testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/mgo.v2/bson"

	"github.com/qetuantuan/jengo_recap/dao"
	"github.com/qetuantuan/jengo_recap/model"
)

type mockBuildDao_testRepo struct {
	Sb *model.SemanticBuild // keep one in mem is enough for now
}

func (m *mockBuildDao_testRepo) FindSemanticBuildByBranchCommit(repoId, commitId, branch string) (sbuild model.SemanticBuild, err error) {
	if m.Sb != nil {
		sbuild = *m.Sb
	} else {
		err = dao.ErrorBuildNotFound
	}
	return
}

func (m *mockBuildDao_testRepo) GetSemanticBuilds(sbuildIds []string) (sbuilds model.SemanticBuilds, err error) {
	if m.Sb != nil {
		sbuilds = append(sbuilds, *m.Sb)
	}
	return
}

func (m *mockBuildDao_testRepo) GetSemanticBuildsByFilter(filter map[string]interface{}, limitCount, offset int) (sbuilds model.SemanticBuilds, err error) {
	return
}

func (m *mockBuildDao_testRepo) IsBuildExistInSemanticBuild(buildId, sBuildId string) (res bool, err error) {
	return
}

func (m *mockBuildDao_testRepo) CreateSemanticBuild(b model.SemanticBuild) (id string, err error) {
	m.Sb = &model.SemanticBuild{
		Id: bson.NewObjectId().Hex(),
	}
	id = m.Sb.Id
	return
}

func (m *mockBuildDao_testRepo) InsertBuild(sbuildId string, build model.Build) (err error) {
	if m.Sb != nil {
		m.Sb.Builds = append(m.Sb.Builds, build)
	} else {
		err = dao.ErrorBuildNotFound
	}
	return
}

func (m *mockBuildDao_testRepo) UpdateBuildProperties(sBuildId string, buildId string, p map[string]interface{}) (err error) {
	return
}

func (m *mockBuildDao_testRepo) UpdateBuildLog(sbuildId, buildId string, logId string) (err error) {
	return
}

type mockRepoDao_testRepo struct {
	R *model.Repo
}

func (m *mockRepoDao_testRepo) GetRepo(id string) (Repo model.Repo, err error) { return }
func (m *mockRepoDao_testRepo) GetReposByFilter(filter map[string]interface{}, limitCount, offset int) (Repos []model.Repo, err error) {
	return
}

func (m *mockRepoDao_testRepo) GetReposByScms(userId string, scms []string) (Repos []model.Repo, err error) {
	return
}

func (m *mockRepoDao_testRepo) GetRepos(userId string, limitCount, offset int) (Repos []model.Repo, err error) {
	return
}
func (m *mockRepoDao_testRepo) GetBuildIndex(id string) (idx int, err error) { return }

func (m *mockRepoDao_testRepo) UpsertRepoMeta(Repos []model.Repo, userId string) (err error) { return }
func (m *mockRepoDao_testRepo) UpdateDynamicRepoInfo(id, branch string) (err error)          { return }
func (m *mockRepoDao_testRepo) SwitchRepo(id string, enableStatus bool) (err error)          { return }
func (m *mockRepoDao_testRepo) UnlinkRepos(Repos []model.Repo, userId string) (err error)    { return }

type projectsTestData struct {
	N       map[string]*model.Repo
	O       map[string]*model.Repo
	ExpectD []model.Repo
	ExpectU []model.Repo
	ExpectI []model.Repo
}

var _ = Describe("Test Repo Service", func() {
	Describe("Create a Repo", func() {
		var (
			service *LocalRepoService
			err     error
		)
		BeforeEach(func() {
			service = &LocalRepoService{
				Md:      &mockRepoDao_testRepo{},
				BuildMd: &mockBuildDao_testRepo{},
			}
			// sbuild1, err = service.InsertBuild(b)
		})
		It("Should return success", func() {
			Expect(err).NotTo(HaveOccurred())
		})
		It("Should be able to get the build", func() {
			// builds, err := service.GetSemanticBuildsByIds([]string{sbuild1.Id})
			Expect(err).NotTo(HaveOccurred())
		})
	})
})

/*
func CompareRepoSet(actual, expected []model.Repo) bool {
	if len(expected) != len(actual) {
		return false
	}
	sort.Sort(model.ById(actual))
	sort.Sort(model.ById(expected))
	for i, _ := range actual {
		if actual[i].Meta.Id != expected[i].Meta.Id {
			return false
		}
	}
	return true
}
func TestSyncRepoSet(t *testing.T) {
	tdata := []projectsTestData{
		projectsTestData{
			N: map[string]*model.Repo{
				"e1": &model.Repo{Repo: vo.Repo{Meta: vo.RepoMeta{Id: "e1"}}},
				"d1": &model.Repo{Repo: vo.Repo{Meta: vo.RepoMeta{Id: "d1"}}},
				"d3": &model.Repo{Repo: vo.Repo{Meta: vo.RepoMeta{Id: "d3"}}},
			},
			O: map[string]*model.Repo{
				"e1": &model.Repo{Repo: vo.Repo{Meta: vo.RepoMeta{Id: "e1"}}},
				"d1": &model.Repo{Repo: vo.Repo{Meta: vo.RepoMeta{Id: "d1"}}},
				"e2": &model.Repo{Repo: vo.Repo{Meta: vo.RepoMeta{Id: "e2"}}},
				"d2": &model.Repo{Repo: vo.Repo{Meta: vo.RepoMeta{Id: "d2"}}},
			},
			ExpectD: []model.Repo{
				model.Repo{Repo: vo.Repo{Meta: vo.RepoMeta{Id: "e2"}}},
				model.Repo{Repo: vo.Repo{Meta: vo.RepoMeta{Id: "d2"}}},
			},
			ExpectU: []model.Repo{
				model.Repo{Repo: vo.Repo{Meta: vo.RepoMeta{Id: "e1"}}},
				model.Repo{Repo: vo.Repo{Meta: vo.RepoMeta{Id: "d1"}}},
			},
			ExpectI: []model.Repo{
				model.Repo{Repo: vo.Repo{Meta: vo.RepoMeta{Id: "d3"}}},
			},
		},
	}

	for _, info := range tdata {
		d, u, i := syncRepoSet(info.N, info.O)
		// Not doing an ID comparison to save a loop search or a map var
		// Using naming convension that "e*" means enabled while "d*" vice verse
		if !CompareRepoSet(d, info.ExpectD) {
			t.Errorf("project to be deleted not as expected.\n d: %v\n expect: %v", d, info.ExpectD)
		}
		if !CompareRepoSet(u, info.ExpectU) {
			t.Errorf("project to be updated not as expected.\n d: %v\n expect: %v", u, info.ExpectU)
		}
		if !CompareRepoSet(i, info.ExpectI) {
			t.Errorf("project to be insert not as expected.\n d: %v\n expect: %v", i, info.ExpectI)
		}
	}
}
*/
