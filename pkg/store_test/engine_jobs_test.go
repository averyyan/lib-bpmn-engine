package store_test

import (
	"github.com/nitram509/lib-bpmn-engine/pkg/bpmn_engine_store"
	"testing"

	"github.com/corbym/gocrest/has"
	"github.com/corbym/gocrest/is"
	"github.com/corbym/gocrest/then"
	"github.com/nitram509/lib-bpmn-engine/pkg/spec/BPMN20/activity"
	"github.com/nitram509/lib-bpmn-engine/pkg/spec/BPMN20/process_instance"
)

const (
	varCounter                  = "counter"
	varEngineValidationAttempts = "engineValidationAttempts"
	varFoobar                   = "foobar"
)

func Test_a_job_can_fail_and_keeps_the_instance_in_active_state(t *testing.T) {
	// setup
	process, _ := bpmnEngine.LoadFromFile(ctx, "../../test-cases/simple_task.bpmn")
	bpmnEngine.AddTaskHandler("id", jobFailHandler)

	instance, _ := bpmnEngine.CreateAndRunInstance(ctx, process.GetProcessKey(), nil)
	state, err := instance.GetState(ctx)
	if err != nil {
		return
	}
	then.AssertThat(t, state, is.EqualTo(process_instance.ACTIVE))
}

func Test_simple_count_loop(t *testing.T) {
	// setup
	process, _ := bpmnEngine.LoadFromFile(ctx, "../../test-cases/simple-count-loop.bpmn")
	bpmnEngine.AddTaskHandler("id-increaseCounter", increaseCounterHandler)

	vars := map[string]interface{}{}
	vars[varCounter] = 0
	instance, _ := bpmnEngine.CreateAndRunInstance(ctx, process.GetProcessKey(), vars)
	variable, err := instance.GetVariable(ctx, varCounter)
	if err != nil {
		return
	}
	then.AssertThat(t, variable, is.EqualTo(4))
	state, err := instance.GetState(ctx)
	if err != nil {
		return
	}
	then.AssertThat(t, state, is.EqualTo(process_instance.COMPLETED))
}

func Test_simple_count_loop_with_message(t *testing.T) {
	// setup
	process, _ := bpmnEngine.LoadFromFile(ctx, "../../test-cases/simple-count-loop-with-message.bpmn")

	vars := map[string]interface{}{}
	vars[varEngineValidationAttempts] = 0
	bpmnEngine.AddTaskHandler("do-nothing", jobCompleteHandler)
	bpmnEngine.AddTaskHandler("validate", func(job bpmn_engine_store.IActivatedJob) {
		variable, err := job.GetVariable(ctx, varEngineValidationAttempts)
		if err != nil {
			return
		}
		attempts := variable.(int)

		foobar := attempts >= 1
		attempts++
		job.SetVariable(ctx, varEngineValidationAttempts, attempts)
		job.SetVariable(ctx, varFoobar, foobar)
		job.Complete(ctx)
	})

	instance, _ := bpmnEngine.CreateAndRunInstance(ctx, process.GetProcessKey(), vars) // should stop at the intermediate message catch event

	_ = bpmnEngine.PublishEventForInstance(ctx, instance.GetInstanceKey(), "msg", nil)
	_, _ = bpmnEngine.RunOrContinueInstance(ctx, instance.GetInstanceKey()) // again, should stop at the intermediate message catch event
	// validation happened
	_ = bpmnEngine.PublishEventForInstance(ctx, instance.GetInstanceKey(), "msg", nil)
	_, _ = bpmnEngine.RunOrContinueInstance(ctx, instance.GetInstanceKey()) // should finish
	// validation happened
	variable, err := instance.GetVariable(ctx, varFoobar)
	if err != nil {
		return
	}
	then.AssertThat(t, variable, is.EqualTo(true))
	getVariable, err := instance.GetVariable(ctx, varEngineValidationAttempts)
	if err != nil {
		return
	}
	then.AssertThat(t, getVariable, is.EqualTo(2))
	state, err := instance.GetState(ctx)
	if err != nil {
		return
	}
	then.AssertThat(t, state, is.EqualTo(process_instance.COMPLETED))

	// internal state expected
	subscriptions, err := instance.FindMessageSubscriptions(ctx)
	if err != nil {
		return
	}

	then.AssertThat(t, subscriptions, has.Length(2))
	getState, err := subscriptions[0].GetState(ctx)
	if err != nil {
		return
	}
	then.AssertThat(t, getState, is.EqualTo(activity.Completed))
	getState2, err := subscriptions[1].GetState(ctx)
	if err != nil {
		return
	}
	then.AssertThat(t, getState2, is.EqualTo(activity.Completed))
}

func Test_activated_job_data(t *testing.T) {
	process, _ := bpmnEngine.LoadFromFile(ctx, "../../test-cases/simple_task.bpmn")
	bpmnEngine.AddTaskHandler("id", func(aj bpmn_engine_store.IActivatedJob) {
		then.AssertThat(t, aj.GetElementId(), is.Not(is.Empty()))
		then.AssertThat(t, aj.GetCreatedAt(), is.Not(is.Nil()))
		//state, _ := aj.GetState(ctx)
		//then.AssertThat(t, state, is.Not(is.EqualTo(activity.Active)))
		then.AssertThat(t, aj.GetProcessInstanceKey(), is.Not(is.EqualTo(int64(0))))
		//then.AssertThat(t, aj.GetKey(), is.Not(is.EqualTo(int64(0))))
		//then.AssertThat(t, aj.GetBpmnProcessId(), is.Not(is.Empty()))
		//then.AssertThat(t, aj.GetProcessDefinitionKey(), is.Not(is.EqualTo(int64(0))))
		//then.AssertThat(t, aj.GetProcessDefinitionVersion(), is.Not(is.EqualTo(int32(0))))
		then.AssertThat(t, aj.GetProcessInstanceKey(), is.Not(is.EqualTo(int64(0))))
	})

	instance, _ := bpmnEngine.CreateAndRunInstance(ctx, process.GetProcessKey(), nil)
	state, err := instance.GetState(ctx)
	if err != nil {
		return
	}
	then.AssertThat(t, state, is.EqualTo(process_instance.ACTIVE))
}

//
func increaseCounterHandler(job bpmn_engine_store.IActivatedJob) {
	variable, err := job.GetVariable(ctx, varCounter)
	if err != nil {
		return
	}
	counter := variable.(int)
	counter++
	job.SetVariable(ctx, varCounter, counter)
	job.Complete(ctx)
}

func jobFailHandler(job bpmn_engine_store.IActivatedJob) {
	job.Fail(ctx, "just because I can")
}

func jobCompleteHandler(job bpmn_engine_store.IActivatedJob) {
	job.Complete(ctx)
}

func Test_task_InputOutput_mapping_happy_path(t *testing.T) {
	// setup
	cp := CallPath{}

	// give
	process, _ := bpmnEngine.LoadFromFile(ctx, "../../test-cases/service-task-input-output.bpmn")
	bpmnEngine.AddTaskHandler("service-task-1", cp.CallPathHandler)
	bpmnEngine.AddTaskHandler("user-task-2", cp.CallPathHandler)

	// when
	pi, err := bpmnEngine.CreateAndRunInstance(ctx, process.GetProcessKey(), nil)
	then.AssertThat(t, err, is.Nil())

	// then
	jobs, _ := pi.FindJobs(ctx)
	for _, job := range jobs {
		state, _ := job.GetState(ctx)
		then.AssertThat(t, state, is.EqualTo(activity.Completed))
	}
	then.AssertThat(t, cp.CallPath, is.EqualTo("service-task-1,user-task-2"))
	variable, _ := pi.GetVariable(ctx, "id")
	then.AssertThat(t, variable, is.EqualTo(1))
	variable2, _ := pi.GetVariable(ctx, "orderId")
	then.AssertThat(t, variable2, is.EqualTo(1234))
	variable3, _ := pi.GetVariable(ctx, "order")
	then.AssertThat(t, variable3, is.EqualTo(map[string]interface{}{
		"name": "order1",
		"id":   "1234",
	}))
	variable4, _ := pi.GetVariable(ctx, "orderName")
	then.AssertThat(t, variable4.(string), is.EqualTo("order1"))
}

func Test_instance_fails_on_Invalid_Input_mapping(t *testing.T) {
	// setup
	cp := CallPath{}

	// give
	process, _ := bpmnEngine.LoadFromFile(ctx, "../../test-cases/service-task-invalid-input.bpmn")
	bpmnEngine.AddTaskHandler("invalid-input", cp.CallPathHandler)

	// when
	pi, err := bpmnEngine.CreateAndRunInstance(ctx, process.GetProcessKey(), nil)
	then.AssertThat(t, err, is.Nil())

	// then
	then.AssertThat(t, cp.CallPath, is.EqualTo(""))
	variable, _ := pi.GetVariable(ctx, "id")

	jobs, _ := pi.FindJobs(ctx)
	then.AssertThat(t, variable, is.EqualTo(nil))
	state, _ := jobs[0].GetState(ctx)
	then.AssertThat(t, state, is.EqualTo(activity.Failed))
	getState, _ := pi.GetState(ctx)
	then.AssertThat(t, getState, is.EqualTo(process_instance.FAILED))
}

func Test_job_fails_on_Invalid_Output_mapping(t *testing.T) {
	// setup
	cp := CallPath{}

	// give
	process, _ := bpmnEngine.LoadFromFile(ctx, "../../test-cases/service-task-invalid-output.bpmn")
	bpmnEngine.AddTaskHandler("invalid-output", cp.CallPathHandler)

	// when
	pi, err := bpmnEngine.CreateAndRunInstance(ctx, process.GetProcessKey(), nil)
	then.AssertThat(t, err, is.Nil())

	// then
	then.AssertThat(t, cp.CallPath, is.EqualTo("invalid-output"))
	variable, err := pi.GetVariable(ctx, "order")
	if err != nil {
		return
	}
	jobs, _ := pi.FindJobs(ctx)

	then.AssertThat(t, variable, is.EqualTo(nil))
	state, _ := jobs[0].GetState(ctx)
	then.AssertThat(t, state, is.EqualTo(activity.Failed))
	getState, _ := pi.GetState(ctx)
	then.AssertThat(t, getState, is.EqualTo(process_instance.FAILED))
}
