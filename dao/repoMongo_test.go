package dao

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/qetuantuan/jengo_recap/model"
)

var _ = Describe("Repo Mongo dao", func() {
	Describe("Prepare multiple Repos", func() {
		var (
			rmdao     *RepoMongoDao
			expected1 *model.Repo
			expected2 *model.Repo
			err       error

			inserted bool
		)
		BeforeEach(func() {
			expected1 = &model.Repo{
				UserIds: []string{"user1"},
			}
			expected1.Id = "id1"
			expected1.Name = new(string)
			*expected1.Name = "name1"

			expected2 = &model.Repo{
				UserIds: []string{"user2"},
			}
			expected2.Id = "id2"
			expected2.ScmName = "github"
			expected2.Name = new(string)
			*expected2.Name = "name2"

			rmdao = &RepoMongoDao{}
			// init data but won't dial to mongo
			rmdao.Init(&MongoDao{Inited: true})
			rmdao.GSession = session

			if inserted {
				return
			}
			err = rmdao.UpsertRepoMeta(
				[]model.Repo{*expected1, *expected2},
				"user1")
			inserted = true
		})
		It("Should return success", func() {
			Expect(err).NotTo(HaveOccurred())
		})
		Context("Getting Repos", func() {
			var actual []model.Repo
			BeforeEach(func() {
				actual, err = rmdao.GetRepos("user1", 20, 0)
			})
			It("Should return the same expected obj", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(len(actual)).To(Equal(2))
				Expect(actual[0].Id).To(Equal("id1"))
				Expect(actual[1].Id).To(Equal("id2"))
				Expect(actual[0].OwnerIds).To(ContainElement("user1"))
				Expect(actual[1].OwnerIds).To(ContainElement("user1"))
			})
		})
		Describe("Testing a multiple repos flow", func() {
			var (
				updated bool
				deleted bool
			)
			Context("Updating Repos", func() {
				BeforeEach(func() {
					updatedRepo := *expected1
					updatedRepo.Enabled = true
					*updatedRepo.Name = "updated1"
					if updated {
						return
					}
					err = rmdao.UpsertRepoMeta([]model.Repo{updatedRepo}, "user3")
					updated = true
				})
				It("Should return success", func() {
					Expect(err).NotTo(HaveOccurred())
				})
				It("Should update repo meta only", func() {
					var actual []model.Repo
					actual, err = rmdao.GetRepos("user3", 20, 0)
					Expect(len(actual)).To(Equal(1))
					Expect(actual[0].Id).To(Equal("id1"))
					Expect(actual[0].Enabled).To(Equal(false))
					Expect(*actual[0].Name).To(Equal("updated1"))
					Expect(actual[0].OwnerIds).To(ContainElement("user1"))
				})
				Context("Unlinking Repos", func() {
					BeforeEach(func() {
						if deleted {
							return
						}
						rmdao.UnlinkRepos([]model.Repo{
							*expected1,
							*expected2,
						}, "user3")
						deleted = true
					})
					It("Should return success", func() {
						Expect(err).NotTo(HaveOccurred())
					})
					It("Should get NO repos", func() {
						actual, err := rmdao.GetRepos("user3", 0, 0)
						Expect(err).NotTo(HaveOccurred())
						Expect(len(actual)).To(Equal(0))
					})
				})
			})
		})
		Describe("Testing various gets", func() {
			It("Should get repos by scm filter", func() {
				actual, err := rmdao.GetReposByScms("user1", []string{"github"})
				Expect(err).NotTo(HaveOccurred())
				Expect(len(actual)).To(Equal(1))
				Expect(actual[0].Id).To(Equal("id2"))
			})
			It("Should not get repos by strange scm filter", func() {
				actual, err := rmdao.GetReposByScms("user1", []string{"", "blahblah"})
				Expect(err).NotTo(HaveOccurred())
				Expect(len(actual)).To(Equal(1))
				Expect(actual[0].Id).To(Equal("id1"))
			})
			It("Should get repos by equal filter", func() {
				actual, err := rmdao.GetReposByFilter(
					map[string]interface{}{
						"_id": "id2",
					}, 20, 0)
				Expect(err).NotTo(HaveOccurred())
				Expect(len(actual)).To(Equal(1))
				Expect(actual[0].Id).To(Equal("id2"))
			})
			It("Should still get repos by limits = 0", func() {
				actual, err := rmdao.GetReposByFilter(
					map[string]interface{}{
						"_id": "id2",
					}, 0, 0)
				Expect(err).NotTo(HaveOccurred())
				Expect(len(actual)).To(Equal(1))
				Expect(actual[0].Id).To(Equal("id2"))
			})
			It("Should get repos by bool filter", func() {
				actual, err := rmdao.GetReposByFilter(
					map[string]interface{}{
						"enabled": false,
					}, 0, 0)
				Expect(err).NotTo(HaveOccurred())
				Expect(len(actual)).To(Equal(2))
			})
			It("Should get repos by in_array filter", func() {
				actual, err := rmdao.GetReposByFilter(
					map[string]interface{}{
						"owner_ids": []string{"user1"},
					}, 0, 0)
				Expect(err).NotTo(HaveOccurred())
				Expect(len(actual)).To(Equal(2))
			})
		})
		Describe("Testing a single repo flow", func() {
			Context("Testing Partial Updating", func() {
				var (
					partialUpdated bool
				)
				BeforeEach(func() {
					if partialUpdated {
						return
					}
					err = rmdao.UpdateDynamicRepoInfo(
						"id1", "branch1")
					partialUpdated = true
				})
				It("Should return success", func() {
					Expect(err).NotTo(HaveOccurred())
				})
				It("Should be able to get updated values", func() {
					actual, err := rmdao.GetRepo("id1")
					Expect(err).NotTo(HaveOccurred())
					Expect(actual.Branches).To(ContainElement("branch1"))
				})
				It("Should be able to update a few of values", func() {
					err = rmdao.UpdateDynamicRepoInfo(
						"id1", "branch2")
					Expect(err).NotTo(HaveOccurred())
					actual, err := rmdao.GetRepo("id1")
					Expect(err).NotTo(HaveOccurred())
					Expect(actual.Branches).To(ContainElement("branch1"))
					Expect(actual.Branches).To(ContainElement("branch2"))
				})
			})
			Context("Testing SwitchRepo", func() {
				var (
					switched bool
				)
				BeforeEach(func() {
					if switched {
						return
					}
					err = rmdao.SwitchRepo("id1", true)
					switched = true
				})
				It("Should return success", func() {
					Expect(err).NotTo(HaveOccurred())
				})
				It("Should be able to get updated values", func() {
					actual, err := rmdao.GetRepo("id1")
					Expect(err).NotTo(HaveOccurred())
					Expect(actual.Enabled).To(Equal(true))
				})
				It("Should be able to switch back", func() {
					err = rmdao.SwitchRepo("id1", false)
					Expect(err).NotTo(HaveOccurred())
					actual, err := rmdao.GetRepo("id1")
					Expect(err).NotTo(HaveOccurred())
					Expect(actual.Enabled).To(Equal(false))
				})
			})
			Context("Testing GetBuildIndex", func() {
				It("Should return sequential number", func() {
					idx, err := rmdao.GetBuildIndex("id1")
					Expect(err).NotTo(HaveOccurred())
					Expect(idx).To(Equal(1))
					idx, err = rmdao.GetBuildIndex("id1")
					Expect(err).NotTo(HaveOccurred())
					Expect(idx).To(Equal(2))
					actual, err := rmdao.GetRepo("id1")
					Expect(err).NotTo(HaveOccurred())
					Expect(actual.BuildIndex).To(Equal(2))
				})
			})
		})
	})
})
