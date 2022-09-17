package store

import (
	"context"
	"errors"
	"fmt"
	"github.com/nitram509/lib-bpmn-engine/pkg/bpmn_engine_store"
	"github.com/nitram509/lib-bpmn-engine/pkg/spec/BPMN20/activity"
	"time"
)

func (store *EngineMemoryStore) CreateJob(ctx context.Context, engineStore bpmn_engine_store.IBpmnEngineStore, processInstanceKey bpmn_engine_store.IProcessInstanceKey, elementId string, elementInstanceKey bpmn_engine_store.IJobKey) (bpmn_engine_store.IJob, error) {
	if store.jobs[processInstanceKey] == nil {
		store.jobs[processInstanceKey] = map[bpmn_engine_store.IJobKey]*job{}
	}
	job := &job{
		engineStore:        engineStore,
		ProcessInstanceKey: processInstanceKey,
		ElementId:          elementId,
		ElementInstanceKey: elementInstanceKey,
		JobKey:             int64(elementInstanceKey) + 1,
		State:              activity.Active,
		CreatedAt:          time.Now(),
	}
	store.jobs[processInstanceKey][elementInstanceKey] = job
	return store.jobs[processInstanceKey][elementInstanceKey], nil
}

//func (store *EngineMemoryStore) FindOrCreateJob(ctx context.Context, processInstanceKey bpmn_engine_store.IProcessInstanceKey, elementId bpmn_engine_store.IJobKey, elementInstanceKey int64) (bpmn_engine_store.IJob, error) {
//	job, err := store.findJob(processInstanceKey, elementId)
//	if err != nil {
//		return store.CreateJob(ctx, processInstanceKey, elementId, elementInstanceKey)
//	}
//	return job, nil
//}

func (store *EngineMemoryStore) FindJob(ctx context.Context, processInstanceKey bpmn_engine_store.IProcessInstanceKey, jobKey bpmn_engine_store.IJobKey) (bpmn_engine_store.IJob, error) {
	return store.findJob(processInstanceKey, jobKey)
}

func (store *EngineMemoryStore) FindJobs(ctx context.Context, processInstanceKey bpmn_engine_store.IProcessInstanceKey) ([]bpmn_engine_store.IJob, error) {
	jobs, err := store.findJobs(processInstanceKey)
	if err != nil {
		return nil, err
	}
	var result []bpmn_engine_store.IJob
	for _, job := range jobs {
		result = append(result, job)
	}
	return result, nil
}

func (store *EngineMemoryStore) SetJobState(
	ctx context.Context,
	processInstanceKey bpmn_engine_store.IProcessInstanceKey,
	jobKey bpmn_engine_store.IJobKey,
	state activity.LifecycleState,
	reason string,
) error {
	job, err := store.findJob(processInstanceKey, jobKey)
	if err != nil {
		return err
	}
	job.State = state
	job.Reason = reason
	return nil
}

func (store *EngineMemoryStore) GetJobState(ctx context.Context, processInstanceKey bpmn_engine_store.IProcessInstanceKey, jobKey bpmn_engine_store.IJobKey) (activity.LifecycleState, error) {
	job, err := store.findJob(processInstanceKey, jobKey)
	if err != nil {
		return "", err
	}
	return job.State, nil
}

func (store *EngineMemoryStore) findJobs(processInstanceKey bpmn_engine_store.IProcessInstanceKey) (map[bpmn_engine_store.IJobKey]*job, error) {
	jobs, isExistJobs := store.jobs[processInstanceKey]
	if !isExistJobs {
		return nil, errors.New(fmt.Sprintf("jobs: ProcessInstanceKey=%d  is not exist", processInstanceKey))
	}
	return jobs, nil
}

func (store *EngineMemoryStore) findJob(processInstanceKey bpmn_engine_store.IProcessInstanceKey, jobKey bpmn_engine_store.IJobKey) (*job, error) {
	jobs, err := store.findJobs(processInstanceKey)
	if err != nil {
		return nil, err
	}
	job, isJobExist := jobs[jobKey]
	if !isJobExist {
		return nil, err
	}
	return job, nil
}
