package service

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/qetuantuan/jengo_recap/dao"
	"github.com/qetuantuan/jengo_recap/model"
	"github.com/qetuantuan/jengo_recap/scm"
)

var _ = Describe("User", func() {
	Describe("With a User in DB", func() {
		var (
			service      *UserService
			userInserted bool
			userId       string = "123"
		)
		Context("After User with id 123 inserted into DB", func() {
			BeforeEach(func() {
				if userInserted {
					return
				}
				service = &UserService{
					Md: &dao.MongoDao{
						GSession: session,
					},
				}

				var user *model.User = &model.User{
					UserId: userId,
					Auths: []model.Auth{
						model.Auth{
							LoginName:    "existed",
							AuthSourceId: "github.com",
							Primary:      true,
						},
					},
					Scms: []model.Scm{
						model.Scm{
							Id: "u_github_123",
						},
					},
				}
				uc := session.DB("users").C("user02")
				_, err := uc.UpsertId(userId, &user)
				Expect(err).NotTo(HaveOccurred())
				err = uc.FindId(userId).One(user)
				Expect(err).NotTo(HaveOccurred())
				Expect(user.UserId).To(Equal(userId))

				userInserted = true
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
				Expect(u.Scms[0].Token).To(Equal("NewToken"))
			})
		})
	})

	Describe("Create user successfully", func() {
		Context("With a fake github server", func() {
			var (
				testServer *httptest.Server
				service    *UserService
			)
			BeforeEach(func() {
				user := model.GithubUser{
					Id:    123,
					Login: "myuser",
				}
				testServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					pBytes, _ := json.Marshal(user)
					By("Fake response:")
					By(string(pBytes))
					w.WriteHeader(200)
					w.Write(pBytes)
					return
				}))

				service = &UserService{
					Md: &dao.MongoDao{
						GSession: session,
					},
					GithubScm: scm.NewGithubScm(),
				}
				service.GithubScm.ApiLink = testServer.URL

			})
			Context("Calling CreateUser", func() {
				var userId string
				BeforeEach(func() {
					var err error
					userId, err = service.CreateUser("myuser", "github.com", "token1")
					Expect(err).NotTo(HaveOccurred())
				})
				It("Should get the user just created", func() {
					u, err := service.GetUser(userId)
					Expect(err).NotTo(HaveOccurred())
					Expect(u.UserId).NotTo(Equal(""))
					Expect(u.Auth.LoginName).To(Equal("myuser"))
					Expect(u.Auth.AuthSource.Id).To(Equal("github.com"))
					Expect(u.Auth.OriginId).To(Equal(123))

					// Note: usually we should test one method per case.
					//       sometimes we may take chance of create/get
					//       since we are not blind testing.
				})
			})
		})
	})
})
