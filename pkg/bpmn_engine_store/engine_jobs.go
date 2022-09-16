package bpmn_engine_store

import (
	"context"
	"fmt"
	"github.com/nitram509/lib-bpmn-engine/pkg/spec/BPMN20"
	"github.com/nitram509/lib-bpmn-engine/pkg/spec/BPMN20/activity"
	"github.com/nitram509/lib-bpmn-engine/pkg/spec/BPMN20/process_instance"
)

func (state *BpmnEngineState) handleServiceTask(ctx context.Context, instance IProcessInstanceInfo, element BPMN20.TaskElement) bool {
	id := element.GetId()
	jobs, err := state.GetStore().FindJobs(ctx, instance.GetInstanceKey())
	if err != nil {
		jobs = []IJob{}
	}
	job := state.findOrCreateJob(ctx, jobs, id, instance)
	fmt.Println("findOrCreateJob activatedJob", nil != state.handlers && nil != state.handlers[id])
	if nil != state.handlers && nil != state.handlers[id] {
		err := state.GetStore().SetJobState(ctx, instance.GetInstanceKey(), job.GetElementInstanceKey(), activity.Active, "")

		if err != nil {
			return false
		}
		activatedJob := state.GetStore().CreateActivatedJob(state, job)

		// TODO here is able to use transaction
		if err := evaluateVariableMapping(ctx, instance, element.GetInputMapping()); err != nil {
			if err := job.SetState(ctx, activity.Failed, ""); err != nil {
				return false
			}
			if err := instance.SetState(ctx, process_instance.FAILED); err != nil {
				return false
			}
			return false
		}
		state.handlers[id](activatedJob)
		if err := evaluateVariableMapping(ctx, instance, element.GetOutputMapping()); err != nil {
			if err := job.SetState(ctx, activity.Failed, ""); err != nil {
				return false
			}
			if err := instance.SetState(ctx, process_instance.FAILED); err != nil {
				return false
			}
			return false
		}
	}

	jobState, err := job.GetState(ctx)
	if err != nil {
		return false
	}
	return jobState == activity.Completed
}

func (state *BpmnEngineState) findOrCreateJob(ctx context.Context, jobs []IJob, id string, instance IProcessInstanceInfo) IJob {
	for _, job := range jobs {
		if job.GetElementId() == id {
			return job
		}
	}
	job, err := instance.CreateJob(ctx, id)
	if err != nil {
		return nil
	}
	return job
}
