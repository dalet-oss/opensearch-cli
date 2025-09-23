package api

import (
	"context"
	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
	"log"
)

// CreateReplication creates the replication task
func (api *OpensearchWrapper) CreateReplication() {
	ctx, cancelFunc := context.WithTimeout(context.TODO(), LightOperationTimeout)
	defer cancelFunc()
	var result interface{}
	if _, err := api.Client.Do(ctx, opensearchapi.InfoReq{}, &result); err != nil {
		log.Fatal(err)
	} else {
		// todo: do something
	}
}

func (api *OpensearchWrapper) PauseReplication()      {}
func (api *OpensearchWrapper) ResumeReplication()     {}
func (api *OpensearchWrapper) StopReplication()       {}
func (api *OpensearchWrapper) StatusReplication()     {}
func (api *OpensearchWrapper) TaskStatusReplication() {}
