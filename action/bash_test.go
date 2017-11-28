package action

import (
	"testing"
)

type fakeContext struct {
}

func (c *fakeContext) Init(initMap map[string]string) {}

func (c *fakeContext) Version() string { return "" }

func (c *fakeContext) Cwd() string { return "" }

func (c *fakeContext) ContextId() string { return "" }

func (c *fakeContext) LogDir() string { return "/tmp" }

func (c *fakeContext) ActionResult(key string, asc bool) (value string) { return "" }

func (c *fakeContext) ActionResultAll(key string) (values []string) { return []string{""} }

func (c *fakeContext) SetCurrentResult(map[string]string) error { return nil }

func TestBashSuccess(t *testing.T) {
	bash := Bash{
		Noop: Noop{
			Id: "testSuccess",
		},
		Timeout:       3600,
		ScriptContent: "ls",
	}
	bash.SetContext(&fakeContext{})
	bash.Do()
	if bash.State() != SUCCESS {
		t.Fatalf("Action failed!")
	}
}

func TestBashFailed(t *testing.T) {
	bash := Bash{
		Noop: Noop{
			Id: "testFailed",
		},
		Timeout:       3600,
		ScriptContent: "cat not_existed",
	}
	bash.SetContext(&fakeContext{})
	bash.Do()
	if bash.State() != FAILED {
		t.Fatalf("Action not failed!")
	}
}

func TestBashTimout(t *testing.T) {
	bash := Bash{
		Noop: Noop{
			Id: "testTimout",
		},
		Timeout:       5,
		ScriptContent: "sleep 10",
	}
	bash.SetContext(&fakeContext{})
	bash.Do()
	if bash.State() != FAILED {
		t.Fatalf("Action not failed!")
	}
}
