package util

import (
	"testing"
)

type idTest struct {
	scm      ScmSymbol
	rawId    interface{}
	expected string
}

func TestGetUserIdOk(t *testing.T) {
	td := []idTest{
		{GithubSymbol, 123, "u_github_123"},
		{GithubSymbol, "abc", "u_github_abc"},
	}
	for _, d := range td {
		if GetUserId(d.scm, d.rawId) != d.expected {
			t.Fatalf("Id from %v %v is not %v", d.scm, d.rawId, d.expected)
		}
	}
}

func TestGetProjectIdOk(t *testing.T) {
	td := []idTest{
		{GithubSymbol, 123, "p_github_123"},
		{GithubSymbol, "abc", "p_github_abc"},
	}
	for _, d := range td {
		if GetProjectId(d.scm, d.rawId) != d.expected {
			t.Fatalf("Id from %v %v is not %v", d.scm, d.rawId, d.expected)
		}
	}
}
