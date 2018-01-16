package service

/*
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"bytes"
	"encoding/json"

	"github.com/qetuantuan/jengo_recap/api"
	"github.com/qetuantuan/jengo_recap/dao"
	"github.com/qetuantuan/jengo_recap/queue"
*/

type mockProjectStoreClient struct {
}

/*
func (m *mockProjectStoreClient) InsertBuild(build *api.Build) (*http.Response, []byte, error) {
	body, _ := json.Marshal(map[string]string{"build_id": "1234"})
	resp := &http.Response{
		Status:        "200 OK",
		StatusCode:    200,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Body:          ioutil.NopCloser(bytes.NewBufferString(string(body))),
		ContentLength: int64(len(body)),
		Header:        make(http.Header, 0),
	}

	return resp, body, nil
}

func (m *mockProjectStoreClient) PutLog(userId string, projectId string, buildId string, buildId string, buildLog string) (*http.Response, []byte, error) {
	return &http.Response{}, nil, nil
}

func (m *mockProjectStoreClient) UpdateBuild(build *api.Build, buildId string) (*http.Response, []byte, error) {
	return &http.Response{}, nil, nil
}

func (p *mockProjectStoreClient) UpdatePatchBuild(build *api.PatchBuild, buildId string) (*http.Response, []byte, error) {
	return &http.Response{}, nil, nil
}

func (p *mockProjectStoreClient) GetProject(userId string, projectId string) (project api.Project, err error) {
	return
}

func TestEngineBuildService_CreateBuild(t *testing.T) {
	d, server, session := SetupMongo()
	defer TearDown(d, server, session)

	s := &EngineBuildService{&mockProjectStoreClient{}, &dao.BuildDao{GSession: session}, queue.NewNativeTaskQueue()}

	// test CreateBuild
	commit := api.PushEventCommit{}
	repoId := 123
	params := &api.EngineCreateBuildParams{
		UserId:    "user_id",
		ProjectId: "project_id",
		Repo:      &api.PushEventRepository{ID: &repoId},
		Commits:   []api.PushEventCommit{commit},
	}
	build, err := s.CreateBuild(params)
	if err != nil {
		t.Error(fmt.Sprintf("[TestEngineBuildService_CreateBuild][CreateBuild] error in CreateBuild"))
	}
	if !checkBuildIsSame(build, params) {
		t.Error(fmt.Sprintf("[TestEngineBuildService_CreateBuild][CreateBuild] the build is not the same"))
	}

	dParams := &api.EngineDescribeBuildsParams{
		UserId:    build.UserId,
		ProjectId: build.EventId,
	}

	// test DescribeBuilds
	buildsAct, err := s.DescribeBuilds(dParams)
	if err != nil {
		t.Error(fmt.Sprintf("[TestEngineBuildService_CreateBuild][DescribeBuilds] error in DescribeBuilds"))
	}
	if len(buildsAct) != 1 {
		t.Error(fmt.Printf("[TestEngineBuildService_CreateBuild][DescribeBuilds] the length of builds is not 1: %v", buildsAct))
	}
	if !checkBuildIsSame(&(buildsAct[0]), params) {
		t.Error(fmt.Printf("[TestEngineBuildService_CreateBuild][DescribeBuilds] the user_id is %s, the project_id is %s", buildsAct[0].UserId, buildsAct[0].ProjectId))
	}

	// test DescribeBuild
	buildAct, err := s.DescribeBuild(build.Id)
	if err != nil {
		t.Error(fmt.Sprintf("[TestEngineBuildService_CreateBuild][DescribeBuild] error in DescribeBuild"))
	}
	if !checkBuildIsSame(buildAct, params) {
		t.Error(fmt.Printf("[TestEngineBuildService_CreateBuild][DescribeBuild] the user_id is %s, the project_id is %s", buildAct.UserId, buildAct.ProjectId))
	}
}

func checkBuildIsSame(build *api.Build, params *api.EngineCreateBuildParams) bool {
	if build.ProjectId != params.ProjectId {
		return false
	} else if build.UserId != params.UserId {
		return false
	} else {
		return true
	}
}

func TestEngineBuildService_DescribeBuilds(t *testing.T) {
	d, server, session := SetupMongo()
	defer TearDown(d, server, session)

	s := &EngineBuildService{&mockProjectStoreClient{}, &dao.BuildDao{GSession: session}, queue.NewNativeTaskQueue()}
	p := &api.EngineDescribeBuildsParams{
		ProjectId: "p_123456789_123456789",
	}
	builds, err := s.DescribeBuilds(p)
	if err != nil {
		t.Error(fmt.Sprintf("[TestEngineBuildService_DescribeBuilds] error in DescribeBuilds"))
	}
	if len(builds) != 0 {
		t.Error(fmt.Sprintf("[TestEngineBuildService_DescribeBuilds] the length of builds is not 0"))
	}
}

func TestEngineBuildService_DescribeBuild(t *testing.T) {
	d, server, session := SetupMongo()
	defer TearDown(d, server, session)

	s := &EngineBuildService{&mockProjectStoreClient{}, &dao.BuildDao{GSession: session}, queue.NewNativeTaskQueue()}
	build, err := s.DescribeBuild("build_123456789_123456789")
	if err == nil {
		t.Error(fmt.Sprintf("[TestEngineBuildService_DescribeBuild] No error in DescribeBuild"))
	}
	if build != nil {
		t.Error(fmt.Sprintf("[TestEngineBuildService_DescribeBuild] the build is not nil"))
	}
}
*/
