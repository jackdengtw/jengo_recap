package scm

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/qetuantuan/jengo_recap/model"
)

func compareStringArray(a1, a2 []string) bool {
	if len(a1) != len(a2) {
		return false
	}
	for i := range a1 {
		if a1[i] != a2[i] {
			return false
		}
	}
	return true
}

func compareProject(pp1, pp2 model.Project) (res bool) {
	p1 := pp1.Project
	p2 := pp2.Project
	if p1.Meta != p2.Meta ||
		p1.Enable != p2.Enable ||
		p1.LatestBuildId != p2.LatestBuildId ||
		p1.RunIndex != p2.RunIndex ||
		p1.State != p2.State ||
		compareStringArray(p1.Branches, p2.Branches) ||
		compareStringArray(p1.Users, p2.Users) {
		res = false
		return
	}
	return true
}

func getProjectByIndex(idx int) (githubProject GithubProject) {
	t := time.Now().UTC()
	idxStr := strconv.Itoa(idx)
	githubProject = GithubProject{
		1000 + idx,
		"testName_" + idxStr,
		"testUser/testName_" + idxStr,
		User{"testUser_" + idxStr, 100 + idx},
		"www.testhtml_" + idxStr + ".com",
		t.Add(time.Duration(idx)),
		t.Add(time.Duration(idx)),
		t.Add(time.Duration(idx)),
		"www.hooks_html.com",
		"www.gittest.com" + idxStr,
		"sdasda" + idxStr,
		"sdasdasdas" + idxStr,
		idx%2 == 1,
		"golang" + idxStr,
	}
	return githubProject
}

func TestGithubScm_GetProjectList(t *testing.T) {

	githubProjects := [3]GithubProject{}
	for i := 0; i < 3; i++ {
		githubProjects[i] = getProjectByIndex(i)
	}
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pBytes, _ := json.Marshal(githubProjects)
		fmt.Println(string(pBytes))
		w.Write(pBytes)
		w.WriteHeader(200)
		return
	}))

	defer testServer.Close()
	githubScm := NewGithubScm("www.hooks_html.com")
	githubScm.SetApiLink(testServer.URL)
	projects, err := githubScm.GetProjectList()
	if err != nil {
		t.Fatal("get err:", err, testServer.URL)
	}
	if len(projects) != 3 {
		t.Fatal("len of projects not 1")
	}
	expProjects := [3]model.Project{}
	for i := 0; i < 3; i++ {
		githubProjects[i].CopyTo(&expProjects[i])
		if compareProject(expProjects[i], projects[i]) {
			t.Fatal(expProjects[i], projects[i])
		}
	}

	str, _ := json.Marshal(projects)
	fmt.Println(string(str))
}
