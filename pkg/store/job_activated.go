package store

import (
	"context"
	"github.com/nitram509/lib-bpmn-engine/pkg/bpmn_engine_store"
	"github.com/nitram509/lib-bpmn-engine/pkg/spec/BPMN20/activity"
)

// ActivatedJob did not save to db
type ActivatedJob struct {
	job bpmn_engine_store.IJob
}

func (aj *ActivatedJob) SetVariable(ctx context.Context, key string, value interface{}) error {
	return aj.job.SetVariable(ctx, key, value)
}

func (aj *ActivatedJob) GetVariable(ctx context.Context, key string) (interface{}, error) {
	return aj.job.GetVariable(ctx, key)
}

func (aj *ActivatedJob) GetElementId() string {
	return aj.job.GetElementId()
}

func (aj *ActivatedJob) Fail(ctx context.Context, reason string) error {
	return aj.job.SetState(ctx, activity.Failed, reason)
}

func (aj *ActivatedJob) Complete(ctx context.Context) error {
	return aj.job.SetState(ctx, activity.Completed, "")
}
