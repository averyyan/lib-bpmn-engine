package bpmn_engine_store

import (
	"context"
	"github.com/nitram509/lib-bpmn-engine/pkg/spec/BPMN20"
	"github.com/nitram509/lib-bpmn-engine/pkg/spec/BPMN20/activity"
	"github.com/nitram509/lib-bpmn-engine/pkg/spec/BPMN20/process_instance"
)

// TODO here need throw a err to use transaction
func (state *BpmnEngineState) handleServiceTask(ctx context.Context, instance IProcessInstanceInfo, element BPMN20.TaskElement) bool {
	id := element.GetId()
	job := findOrCreateJob(ctx, instance, id)
	if nil != state.handlers && nil != state.handlers[id] {
		err := instance.SetJobState(ctx, job.GetElementInstanceKey(), activity.Active, "")

		if err != nil {
			return false
		}
		activatedJob := instance.CreateActivatedJob(job)

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
		// TODO try to use job name to set fun
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

func findOrCreateJob(ctx context.Context, instance IProcessInstanceInfo, id string) IJob {
	jobs, err := instance.FindJobs(ctx)
	if err != nil {
		jobs = []IJob{}
	}
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
