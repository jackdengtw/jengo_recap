package dao

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/qetuantuan/jengo_recap/model"
)

var _ = Describe("Semantic Build Mongo dao", func() {
	Describe("Testing Sbuild Operations", func() {
		var (
			sbdao     *SemanticBuildMongoDao
			expected1 *model.SemanticBuild
			expected2 *model.SemanticBuild
			err       error

			inserted bool
		)
		Context("Preparing one Sbuild", func() {
			BeforeEach(func() {
				if inserted {
					return
				}

				sbdao = &SemanticBuildMongoDao{}
				// init data but won't dial to mongo
				sbdao.Init(&MongoDao{Inited: true})
				sbdao.GSession = session

				expected1 = &model.SemanticBuild{
					RepoId:   "repo1",
					CommitId: "commit1",
					Branch:   "branch1",
				}

				expected2 = &model.SemanticBuild{
					RepoId:   "repo2",
					CommitId: "commit2",
					Branch:   "branch2",
				}

				expected1.Id, err = sbdao.CreateSemanticBuild(*expected1)
				inserted = true
			})
			It("Should return success", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(expected1.Id).NotTo(Equal(""))
			})
			Context("Finding the Sbuild", func() {
				var actual model.SemanticBuild
				BeforeEach(func() {
					actual, err = sbdao.FindSemanticBuildByBranchCommit(
						expected1.RepoId, expected1.CommitId, expected1.Branch)
				})
				It("Should return the same expected obj", func() {
					Expect(err).NotTo(HaveOccurred())
					Expect(actual.Id).To(Equal(expected1.Id))
					Expect(actual.RepoId).To(Equal(expected1.RepoId))
				})
			})
			Context("Finding non-existed Sbuild", func() {
				var actual model.SemanticBuild
				BeforeEach(func() {
					actual, err = sbdao.FindSemanticBuildByBranchCommit(
						expected1.RepoId, expected1.CommitId, "nobranch")
				})
				It("Should return error", func() {
					Expect(err).To(Equal(ErrorBuildNotFound))
				})
			})
			Context("Preparing another Sbuild", func() {
				var inserted2 bool
				BeforeEach(func() {
					if inserted2 {
						return
					}
					expected2.Id, err = sbdao.CreateSemanticBuild(*expected2)
					inserted2 = true
				})
				It("Should return success", func() {
					Expect(err).NotTo(HaveOccurred())
					Expect(expected1.Id).NotTo(Equal(""))
				})
				Context("Get two Sbuilds with ID", func() {
					var actuals []model.SemanticBuild
					BeforeEach(func() {
						actuals, err = sbdao.GetSemanticBuilds(
							[]string{expected1.Id, expected2.Id},
						)
					})
					It("Should return success", func() {
						Expect(err).NotTo(HaveOccurred())
					})
					It("Should have input IDs", func() {
						Expect(len(actuals)).To(Equal(2))
						if actuals[0].Id == expected1.Id {
							Expect(actuals[1].Id).To(Equal(expected2.Id))
						} else {
							Expect(actuals[0].Id).To(Equal(expected2.Id))
							Expect(actuals[1].Id).To(Equal(expected1.Id))
						}
					})
				})
				Context("Get Sbuild with filter", func() {
					var actuals []model.SemanticBuild
					BeforeEach(func() {
						actuals, err = sbdao.GetSemanticBuildsByFilter(
							map[string]interface{}{"branch": "branch2"},
							0, 0,
						)
					})
					It("Should return success", func() {
						Expect(err).NotTo(HaveOccurred())
					})
					It("Should have filtered Sbuild", func() {
						Expect(len(actuals)).To(Equal(1))
						Expect(actuals[0].Id).To(Equal(expected2.Id))
						Expect(actuals[0].Branch).To(Equal(expected2.Branch))
					})
				})
				Context("Get Sbuild with repo Ids", func() {
					var actuals []model.SemanticBuild
					BeforeEach(func() {
						actuals, err = sbdao.GetSemanticBuildsByRepoIds(
							[]string{"repo1"},
						)
					})
					It("Should return success", func() {
						Expect(err).NotTo(HaveOccurred())
					})
					It("Should have filtered Sbuild", func() {
						Expect(len(actuals)).To(Equal(1))
						Expect(actuals[0].RepoId).To(Equal("repo1"))
					})
				})

			})
			Context("Preparing build for Sbuild1", func() {
				var (
					buildInserted  bool
					buildExpected1 *model.Build
				)
				BeforeEach(func() {
					buildExpected1 = &model.Build{
						Id:     "build1",
						Status: "statusUp",
					}

					if buildInserted {
						return
					}
					sbdao.InsertBuild(expected1.Id, *buildExpected1)
					buildInserted = true
				})
				It("Should return success", func() {
					Expect(err).NotTo(HaveOccurred())
				})
				It("Should be in Sbuild1", func() {
					actual, err := sbdao.IsBuildExistInSemanticBuild(
						buildExpected1.Id,
						expected1.Id)
					Expect(err).NotTo(HaveOccurred())
					Expect(actual).To(Equal(true))
				})
				Context("Updating build properties", func() {
					var buildUpdated bool
					BeforeEach(func() {

						if buildUpdated {
							return
						}
						err = sbdao.UpdateBuildProperties(
							expected1.Id,
							buildExpected1.Id,
							map[string]interface{}{"status": "statusDown"})
						buildUpdated = true
					})
					It("Should return success", func() {
						Expect(err).NotTo(HaveOccurred())
					})
					It("Should get updated status", func() {
						var actuals []model.SemanticBuild
						actuals, err = sbdao.GetSemanticBuilds([]string{expected1.Id})
						Expect(err).NotTo(HaveOccurred())
						Expect(len(actuals)).To(Equal(1))
						Expect(len(actuals[0].Builds)).To(Equal(1))
						Expect(actuals[0].Builds[0].Status).To(Equal("statusDown"))
					})
				})
				Context("Updating build log", func() {
					var buildlogUpdated bool
					BeforeEach(func() {
						if buildlogUpdated {
							return
						}
						err = sbdao.UpdateBuildLog(expected1.Id, buildExpected1.Id, "loguri1")
						buildlogUpdated = true
					})
					It("Should return success", func() {
						Expect(err).NotTo(HaveOccurred())
					})
					It("Should get updated log uri", func() {
						var actuals []model.SemanticBuild
						actuals, err = sbdao.GetSemanticBuilds([]string{expected1.Id})
						Expect(err).NotTo(HaveOccurred())
						Expect(len(actuals)).To(Equal(1))
						Expect(len(actuals[0].Builds)).To(Equal(1))
						Expect(*actuals[0].Builds[0].LogUri).To(Equal("loguri1"))
					})
				})
			})
		})
	})
})
