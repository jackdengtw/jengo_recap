package service

import (
	"sort"
	"testing"

	"github.com/qetuantuan/jengo_recap/api"
	"github.com/qetuantuan/jengo_recap/model"
)

type projectsTestData struct {
	N       map[string]*model.Project
	O       map[string]*model.Project
	ExpectD []model.Project
	ExpectU []model.Project
	ExpectI []model.Project
}

func CompareProjectSet(actual, expected []model.Project) bool {
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
func TestSyncProjectSet(t *testing.T) {
	tdata := []projectsTestData{
		projectsTestData{
			N: map[string]*model.Project{
				"e1": &model.Project{Project: api.Project{Meta: api.ProjectMeta{Id: "e1"}}},
				"d1": &model.Project{Project: api.Project{Meta: api.ProjectMeta{Id: "d1"}}},
				"d3": &model.Project{Project: api.Project{Meta: api.ProjectMeta{Id: "d3"}}},
			},
			O: map[string]*model.Project{
				"e1": &model.Project{Project: api.Project{Meta: api.ProjectMeta{Id: "e1"}}},
				"d1": &model.Project{Project: api.Project{Meta: api.ProjectMeta{Id: "d1"}}},
				"e2": &model.Project{Project: api.Project{Meta: api.ProjectMeta{Id: "e2"}}},
				"d2": &model.Project{Project: api.Project{Meta: api.ProjectMeta{Id: "d2"}}},
			},
			ExpectD: []model.Project{
				model.Project{Project: api.Project{Meta: api.ProjectMeta{Id: "e2"}}},
				model.Project{Project: api.Project{Meta: api.ProjectMeta{Id: "d2"}}},
			},
			ExpectU: []model.Project{
				model.Project{Project: api.Project{Meta: api.ProjectMeta{Id: "e1"}}},
				model.Project{Project: api.Project{Meta: api.ProjectMeta{Id: "d1"}}},
			},
			ExpectI: []model.Project{
				model.Project{Project: api.Project{Meta: api.ProjectMeta{Id: "d3"}}},
			},
		},
	}

	for _, info := range tdata {
		d, u, i := syncProjectSet(info.N, info.O)
		// Not doing an ID comparison to save a loop search or a map var
		// Using naming convension that "e*" means enabled while "d*" vice verse
		if !CompareProjectSet(d, info.ExpectD) {
			t.Errorf("project to be deleted not as expected.\n d: %v\n expect: %v", d, info.ExpectD)
		}
		if !CompareProjectSet(u, info.ExpectU) {
			t.Errorf("project to be updated not as expected.\n d: %v\n expect: %v", u, info.ExpectU)
		}
		if !CompareProjectSet(i, info.ExpectI) {
			t.Errorf("project to be insert not as expected.\n d: %v\n expect: %v", i, info.ExpectI)
		}
	}
}
