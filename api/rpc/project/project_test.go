package project

import (
	"context"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/appcelerator/amp/data/storage/etcd"
	"github.com/gogo/protobuf/proto"
)

const (
	etcdDefaultEndpoint = "http://localhost:2379"
)

var (
	etcdEndpoints = []string{etcdDefaultEndpoint}
	proj          *Proj
	sampleProject = &Project{RepoId: 12345, OwnerName: "amp", RepoName: "amp-repo", Token: "FakeToken"}
)

func TestMain(m *testing.M) {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Lshortfile)
	log.SetPrefix("test: ")
	proj = createProjectServer()
	os.Exit(m.Run())
}

func create() (*CreateReply, error) {
	req := &CreateRequest{Project: sampleProject}
	return proj.Create(context.Background(), req)
}

func delete() (*DeleteReply, error) {
	return proj.Delete(context.Background(), &DeleteRequest{RepoId: sampleProject.RepoId})
}

func TestCreate(t *testing.T) {
	delete()
	resp, err := create()
	if err != nil {
		t.Error(err)
	}
	if proto.Equal(resp.Project, sampleProject) {
		t.Errorf("expected %v, got %v", sampleProject, resp.Project)
	}
}

func TestCreateAlreadyExists(t *testing.T) {
	delete()
	create()
	resp, err := create()
	// Should result in a duplicate entry
	if err == nil {
		t.Errorf("Expected Duplicate Entry, got %v \n", resp)
	}
}

func TestDelete(t *testing.T) {
	resp, err := delete()
	if err != nil {
		t.Error(err)
	}
	if !proto.Equal(resp.Project, sampleProject) {
		t.Errorf("expected %v, got %v", sampleProject, resp.Project)
	}
}

// Create the ProjectServer as a local struct that can be excercised directly over the call stack
func createProjectServer() *Proj {
	//Create the config
	var proj = &Proj{}

	if endpoints := os.Getenv("endpoints"); endpoints != "" {
		etcdEndpoints = strings.Split(endpoints, ",")
	}
	proj.Store = etcd.New(etcdEndpoints, "amp")
	if err := proj.Store.Connect(5 * time.Second); err != nil {
		panic(err)
	}
	return proj
}
