package service

import (
	// "sort"
	// "testing"

	// "github.com/qetuantuan/jengo_recap/api"
	"github.com/qetuantuan/jengo_recap/model"
)

type projectsTestData struct {
	N       map[string]*model.Repo
	O       map[string]*model.Repo
	ExpectD []model.Repo
	ExpectU []model.Repo
	ExpectI []model.Repo
}

/*
func CompareRepoSet(actual, expected []model.Repo) bool {
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
func TestSyncRepoSet(t *testing.T) {
	tdata := []projectsTestData{
		projectsTestData{
			N: map[string]*model.Repo{
				"e1": &model.Repo{Repo: api.Repo{Meta: api.RepoMeta{Id: "e1"}}},
				"d1": &model.Repo{Repo: api.Repo{Meta: api.RepoMeta{Id: "d1"}}},
				"d3": &model.Repo{Repo: api.Repo{Meta: api.RepoMeta{Id: "d3"}}},
			},
			O: map[string]*model.Repo{
				"e1": &model.Repo{Repo: api.Repo{Meta: api.RepoMeta{Id: "e1"}}},
				"d1": &model.Repo{Repo: api.Repo{Meta: api.RepoMeta{Id: "d1"}}},
				"e2": &model.Repo{Repo: api.Repo{Meta: api.RepoMeta{Id: "e2"}}},
				"d2": &model.Repo{Repo: api.Repo{Meta: api.RepoMeta{Id: "d2"}}},
			},
			ExpectD: []model.Repo{
				model.Repo{Repo: api.Repo{Meta: api.RepoMeta{Id: "e2"}}},
				model.Repo{Repo: api.Repo{Meta: api.RepoMeta{Id: "d2"}}},
			},
			ExpectU: []model.Repo{
				model.Repo{Repo: api.Repo{Meta: api.RepoMeta{Id: "e1"}}},
				model.Repo{Repo: api.Repo{Meta: api.RepoMeta{Id: "d1"}}},
			},
			ExpectI: []model.Repo{
				model.Repo{Repo: api.Repo{Meta: api.RepoMeta{Id: "d3"}}},
			},
		},
	}

	for _, info := range tdata {
		d, u, i := syncRepoSet(info.N, info.O)
		// Not doing an ID comparison to save a loop search or a map var
		// Using naming convension that "e*" means enabled while "d*" vice verse
		if !CompareRepoSet(d, info.ExpectD) {
			t.Errorf("project to be deleted not as expected.\n d: %v\n expect: %v", d, info.ExpectD)
		}
		if !CompareRepoSet(u, info.ExpectU) {
			t.Errorf("project to be updated not as expected.\n d: %v\n expect: %v", u, info.ExpectU)
		}
		if !CompareRepoSet(i, info.ExpectI) {
			t.Errorf("project to be insert not as expected.\n d: %v\n expect: %v", i, info.ExpectI)
		}
	}
}
*/
