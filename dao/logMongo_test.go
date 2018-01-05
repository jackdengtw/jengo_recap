package dao

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/qetuantuan/jengo_recap/model"
)

var _ = Describe("Log Mongo dao", func() {
	Describe("Testing a basic flow", func() {
		var (
			lmdao      *LogMongoDao
			expected   *model.BuildLog
			contentStr string = "i am content"
			err        error

			inserted bool
			got      bool
		)
		Context("Inserting Log", func() {
			BeforeEach(func() {
				if inserted {
					return
				}
				lmdao = &LogMongoDao{}
				// init data but won't dial to mongo
				lmdao.Init(&MongoDao{Inited: true})
				lmdao.GSession = session

				expected = &model.BuildLog{
					Content: contentStr,
				}

				fmt.Println("Inserting")
				expected.Id, err = lmdao.AddLog([]byte(contentStr))
				inserted = true
			})
			It("Should return success", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(expected.Id).NotTo(Equal(""))
			})
			Context("Getting Log", func() {
				actual := &model.BuildLog{}
				BeforeEach(func() {
					if got {
						return
					}

					fmt.Println("Getting")
					*actual, err = lmdao.GetLog(expected.Id)
					got = true
				})
				It("Should return the same expected obj", func() {
					Expect(err).NotTo(HaveOccurred())
					Expect(actual.Id).To(Equal(expected.Id))
					Expect(actual.Content).To(
						Equal(expected.Content))
				})
			})
		})
	})
})
