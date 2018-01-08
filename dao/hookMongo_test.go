package dao

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/mgo.v2"

	"github.com/qetuantuan/jengo_recap/model"
)

var _ = Describe("Hook Mongo dao", func() {
	Describe("Testing a basic flow", func() {
		var (
			hmdao    *HookMongoDao
			expected *model.GithubHook
			err      error

			inserted bool
			deleted  bool
		)
		Context("Inserting Hook", func() {
			BeforeEach(func() {
				hmdao = &HookMongoDao{}
				// init data but won't dial to mongo
				hmdao.Init(&MongoDao{Inited: true})
				hmdao.GSession = session

				expected = &model.GithubHook{
					Id: "123",
					Config: model.GithubHookConf{
						ContentType: "json",
					},
				}

				if inserted {
					return
				}
				_, err = hmdao.UpsertHook(*expected)
				inserted = true
			})
			It("Should return success", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			Context("Getting Hook", func() {
				var actual *model.GithubHook = &model.GithubHook{}
				BeforeEach(func() {
					*actual, err = hmdao.GetHook(expected.Id)
				})
				It("Should return the same expected obj", func() {
					Expect(err).NotTo(HaveOccurred())
					Expect(actual.Id).To(Equal(expected.Id))
					Expect(actual.Config.ContentType).To(
						Equal(expected.Config.ContentType))
				})
				Context("Deleting Build", func() {
					BeforeEach(func() {
						if deleted {
							return
						}
						err = hmdao.DeleteHook(expected.Id)
					})
					It("Should return success", func() {
						Expect(err).NotTo(HaveOccurred())
					})
					Context("Getting Hook again", func() {
						It("Should not found", func() {
							_, err := hmdao.GetHook(expected.Id)
							Expect(err).To(Equal(mgo.ErrNotFound))
						})
					})
				})
			})
		})
	})
})
