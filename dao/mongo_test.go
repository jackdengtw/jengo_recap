package dao

import (
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type normalStruct struct {
	s  string `json:"ss,omitempty"`
	i  int
	ps *string
}

var _ = Describe("MongoDao", func() {
	Describe("Testing initField", func() {
		var (
			dao MongoDao
			ns  normalStruct
		)
		It("Should have fields inited", func() {
			dao.initFieldJsonTypes(reflect.TypeOf(ns))
			Expect(len(dao.fieldJsonTypes)).To(Equal(3))
			Expect(dao.fieldJsonTypes["ss"].Kind()).To(Equal(reflect.String))
			Expect(dao.fieldJsonTypes["i"].Kind()).To(Equal(reflect.Int))
			Expect(dao.fieldJsonTypes["ps"].Kind()).To(Equal(reflect.String))
		})
		It("Should pass with correct json types", func() {
			sValue := "i am string"
			dao.initFieldJsonTypes(reflect.TypeOf(ns))
			pass := dao.checkFieldType(map[string]interface{}{
				"ss": sValue,
				"i":  1,
				"ps": &sValue,
			})
			Expect(pass).To(Equal(true))
		})
		It("Should pass with correct json types except pointers", func() {
			sValue := "i am string"
			dao.initFieldJsonTypes(reflect.TypeOf(ns))
			pass := dao.checkFieldType(map[string]interface{}{
				"ss": &sValue,
				"i":  1,
				"ps": sValue,
			})
			Expect(pass).To(Equal(true))
		})
		It("Should NOT pass with wrong json types", func() {
			dao.initFieldJsonTypes(reflect.TypeOf(ns))
			pass := dao.checkFieldType(map[string]interface{}{
				"ss": 0,
				"i":  "string",
				"ps": "string",
			})
			Expect(pass).To(Equal(false))
		})
	})
})
