package dao

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/qetuantuan/jengo_recap/model"
)

var _ = Describe("EngineBuildDao", func() {
	Describe("Testing a basic flow", func() {
		var (
			strangeStatus string = "strange"
			ebdao         *EngineBuildMongoDao
			expected      *model.Build
			err           error

			inserted   bool
			got        bool
			updated    bool
			gotUpdated bool
		)
		Context("Inserting Build", func() {
			BeforeEach(func() {
				if inserted {
					return
				}
				ebdao = &EngineBuildMongoDao{}
				// init data but won't dial to mongo
				ebdao.Init(&MongoDao{Inited: true})
				ebdao.GSession = session

				expected = &model.Build{
					Id: "123",
				}

				fmt.Println("Inserting")
				_, err = ebdao.InsertBuild(expected)
				inserted = true
			})
			It("Should return success", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			Context("Getting Build", func() {
				var actual *model.Build
				BeforeEach(func() {
					if got {
						return
					}

					fmt.Println("Getting")
					actual, err = ebdao.GetBuild(expected.Id)
					got = true
				})
				It("Should return the same expected obj", func() {
					Expect(err).NotTo(HaveOccurred())
					Expect(actual.Id).To(Equal(expected.Id))
				})
				Context("Updating Properties", func() {
					BeforeEach(func() {
						if updated {
							return
						}

						fmt.Println("Updating")
						err = ebdao.UpdateBuildProperties(
							expected.Id,
							map[string]interface{}{"status": strangeStatus})
						updated = true
					})
					It("Should return updated success", func() {
						Expect(err).NotTo(HaveOccurred())
					})
					It("Should Not update with wrong type", func() {
						err = ebdao.UpdateBuildProperties(
							expected.Id,
							map[string]interface{}{"result": 1})
						Expect(err).To(Equal(ErrorTypeNotMatch))
					})
					Context("Getting updated Build", func() {
						var actual *model.Build
						BeforeEach(func() {
							if gotUpdated {
								return
							}

							fmt.Println("Getting updated")
							actual, err = ebdao.GetBuild(expected.Id)
							gotUpdated = true
						})
						It("Should return the same expected obj", func() {
							Expect(err).NotTo(HaveOccurred())
							Expect(actual.Status).To(Equal(strangeStatus))
						})
					})
					Context("Listing Build", func() {
						var actual model.Builds
						BeforeEach(func() {
							fmt.Println("Listing build")
							actual, err = ebdao.ListBuilds(
								map[string]interface{}{"status": strangeStatus},
								10,
								0,
							)
						})
						It("Should return a list of strange builds", func() {
							Expect(err).NotTo(HaveOccurred())
							Expect(len(actual)).To(Equal(1))
							Expect(actual[0].Status).To(Equal("strange"))
						})
					})
				})
			})
		})
	})
})