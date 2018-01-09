package dao

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/qetuantuan/jengo_recap/api"
	"github.com/qetuantuan/jengo_recap/model"
)

var _ = Describe("User Mongo dao", func() {
	Describe("Testing a basic flow", func() {
		var (
			usdao      *UserMongoDao
			expected   *model.User
			expectedId string
			err        error

			inserted bool
		)
		Context("Inserting User", func() {
			BeforeEach(func() {
				expectedId = "user1"
				expected = &model.User{
					Id: expectedId,
					Auths: []model.Auth{
						model.Auth{
							AuthBase: api.AuthBase{
								Id:        "auth1",
								LoginName: "login1",
								Primary:   true,
							},
							AuthSourceId: "sourceId1",
						},
					},
				}

				if inserted {
					return
				}
				usdao = &UserMongoDao{}
				// init data but won't dial to mongo
				usdao.Init(&MongoDao{Inited: true})
				usdao.GSession = session

				err = usdao.UpsertUser(expected)
				inserted = true
			})
			It("Should return success", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			Context("Getting User", func() {
				actual := &model.User{}
				BeforeEach(func() {
					*actual, err = usdao.GetUser(expectedId)
				})
				It("Should return the same expected obj", func() {
					Expect(err).NotTo(HaveOccurred())
					Expect(actual.Id).To(Equal(expectedId))
				})
			})
			Context("Getting User by Login", func() {
				actual := &model.User{}
				BeforeEach(func() {
					*actual, err = usdao.GetUserByLogin(
						"login1",
						"sourceId1",
					)
				})
				It("Should return the same expected obj", func() {
					Expect(err).NotTo(HaveOccurred())
					Expect(actual.Id).To(Equal(expectedId))
				})
			})
			Context("Updating Token", func() {
				var updated bool
				BeforeEach(func() {
					if updated {
						return
					}
					err = usdao.UpdateToken(
						expectedId,
						"auths",
						"auth1",
						[]byte("i am a token"),
					)
					updated = true
				})
				It("Should get the updated token", func() {
					Expect(err).NotTo(HaveOccurred())
					actual, err := usdao.GetUser(expectedId)
					Expect(err).NotTo(HaveOccurred())
					Expect(len(actual.Auths)).To(Equal(1))
					Expect(string(actual.Auths[0].Token)).To(Equal("i am a token"))
				})
			})
		})
	})
})
