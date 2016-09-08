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
	port          string
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

func TestCreate(t *testing.T) {
	// Cleanup previous tests
	//	proj.Delete(context.Background(), &DeleteRequest{RepoId: sampleProject.RepoId})
	req := &CreateRequest{Project: sampleProject}
	resp, err := proj.Create(context.Background(), req)

	if err != nil {
		t.Error(err)
	}
	if proto.Equal(resp.Project, sampleProject) {
		t.Errorf("expected %v, got %v", sampleProject, resp.Project)
	}
}

func TestDelete(t *testing.T) {
	resp, err := proj.Delete(context.Background(), &DeleteRequest{RepoId: sampleProject.RepoId})
	if err != nil {
		t.Error(err)
	}
	if !proto.Equal(resp.Project, sampleProject) {
		t.Errorf("expected %v, got %v", sampleProject, resp.Project)
	}
}

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
