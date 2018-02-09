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

func compareRepo(p1, p2 model.Repo) (res bool) {
	if p1.RepoMeta != p2.RepoMeta ||
		p1.Enabled != p2.Enabled ||
		p1.BuildIndex != p2.BuildIndex ||
		compareStringArray(p1.Branches, p2.Branches) ||
		compareStringArray(p1.OwnerIds, p2.OwnerIds) {
		res = false
		return
	}
	return true
}

func getRepoByIndex(idx int) (githubRepo GithubRepo) {
	t := time.Now().UTC()
	idxStr := strconv.Itoa(idx)
	githubRepo = GithubRepo{
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
	return githubRepo
}

func TestGithubScm_GetRepoList(t *testing.T) {

	githubRepos := [3]GithubRepo{}
	for i := 0; i < 3; i++ {
		githubRepos[i] = getRepoByIndex(i)
	}
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pBytes, _ := json.Marshal(githubRepos)
		fmt.Println(string(pBytes))
		w.Write(pBytes)
		w.WriteHeader(200)
		return
	}))

	defer testServer.Close()
	githubScm := NewGithubScm("www.hooks_html.com")
	githubScm.SetApiLink(testServer.URL)
	Repos, err := githubScm.GetRepoList()
	if err != nil {
		t.Fatal("get err:", err, testServer.URL)
	}
	if len(Repos) != 3 {
		t.Fatal("len of Repos not 1")
	}
	expRepos := [3]model.Repo{}
	for i := 0; i < 3; i++ {
		githubRepos[i].CopyTo(&expRepos[i])
		if compareRepo(expRepos[i], Repos[i]) {
			t.Fatal(expRepos[i], Repos[i])
		}
	}

	str, _ := json.Marshal(Repos)
	fmt.Println(string(str))
}
