package service

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/mgo.v2/bson"

	"github.com/qetuantuan/jengo_recap/dao"
	"github.com/qetuantuan/jengo_recap/model"
)

type mockBuildDao_testBuild struct {
	Sb *model.SemanticBuild // keep one in mem is enough for now
}

func (m *mockBuildDao_testBuild) FindSemanticBuildByBranchCommit(repoId, commitId, branch string) (sbuild model.SemanticBuild, err error) {
	if m.Sb != nil {
		sbuild = *m.Sb
	} else {
		err = dao.ErrorBuildNotFound
	}
	return
}

func (m *mockBuildDao_testBuild) GetSemanticBuilds(sbuildIds []string) (sbuilds model.SemanticBuilds, err error) {
	if m.Sb != nil {
		sbuilds = append(sbuilds, *m.Sb)
	}
	return
}

func (m *mockBuildDao_testBuild) GetSemanticBuildsByRepoIds(repoIds []string) (sbuilds model.SemanticBuilds, err error) {
	return
}

func (m *mockBuildDao_testBuild) GetSemanticBuildsByFilter(filter map[string]interface{}, limitCount, offset int) (sbuilds model.SemanticBuilds, err error) {
	return
}

func (m *mockBuildDao_testBuild) IsBuildExistInSemanticBuild(buildId, sBuildId string) (res bool, err error) {
	return
}

func (m *mockBuildDao_testBuild) CreateSemanticBuild(b model.SemanticBuild) (id string, err error) {
	m.Sb = &model.SemanticBuild{
		Id: bson.NewObjectId().Hex(),
	}
	id = m.Sb.Id
	return
}

func (m *mockBuildDao_testBuild) InsertBuild(sbuildId string, build model.Build) (err error) {
	if m.Sb != nil {
		m.Sb.Builds = append(m.Sb.Builds, build)
	} else {
		err = dao.ErrorBuildNotFound
	}
	return
}

func (m *mockBuildDao_testBuild) UpdateBuildProperties(sBuildId string, buildId string, p map[string]interface{}) (err error) {
	return
}

func (m *mockBuildDao_testBuild) UpdateBuildLog(sbuildId, buildId string, logId string) (err error) {
	return
}

var _ = Describe("Test Build Service", func() {
	Describe("Create a build", func() {
		var (
			service   *LocalBuildService
			sbuild1   model.SemanticBuild
			commitID1 string = "commitID1"
			err       error
		)
		BeforeEach(func() {
			service = &LocalBuildService{
				Md: &mockBuildDao_testBuild{},
			}
			b := model.Build{
				HeadCommit: &model.PushEventCommit{},
			}
			b.HeadCommit.ID = &commitID1

			sbuild1, err = service.InsertBuild(b)
		})
		It("Should return success", func() {
			Expect(err).NotTo(HaveOccurred())
		})
		It("Should be able to get the build", func() {
			builds, err := service.GetSemanticBuildsByIds([]string{sbuild1.Id})
			Expect(err).NotTo(HaveOccurred())
			Expect(len(builds)).To(Equal(1))
		})
	})
})
