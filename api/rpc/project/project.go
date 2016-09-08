package project

import (
	"context"
	"fmt"
	"time"

	"github.com/appcelerator/amp/data/storage"
)

const (
	defTimeout = 5 * time.Second
	prefix     = "project"
)

// Proj structure to implement StatsServer interface
type Proj struct {
	Store storage.Interface
}

// Create adds a new entry to the k,v data store
func (p *Proj) Create(ctx context.Context, req *CreateRequest) (*CreateReply, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defTimeout)
	key := fmt.Sprintf("%s/%v", prefix, req.Project.RepoId)
	ttl := int64(0)
	reply := &CreateReply{Project: &Project{}}
	err := p.Store.Create(ctx, key, req.Project, reply.Project, ttl)
	// cancel timeout (release resources) if operation completes before timeout
	defer cancel()
	return reply, err
}

//Delete removes the entry for the specified Key
func (p *Proj) Delete(ctx context.Context, req *DeleteRequest) (*DeleteReply, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defTimeout)
	key := fmt.Sprintf("%s/%v", prefix, req.RepoId)
	reply := &DeleteReply{Project: &Project{}}
	err := p.Store.Delete(ctx, key, reply.Project)
	// cancel timeout (release resources) if operation completes before timeout
	defer cancel()
	return reply, err
}
