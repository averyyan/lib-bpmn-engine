package store_test

import (
	"github.com/corbym/gocrest/is"
	"github.com/corbym/gocrest/then"
	"github.com/nitram509/lib-bpmn-engine/pkg/bpmn_engine_store"
	"github.com/nitram509/lib-bpmn-engine/pkg/spec/BPMN20/process_instance"
	"testing"
)

func Test_user_tasks_can_be_handled(t *testing.T) {
	// setup
	process, err := bpmnEngine.LoadFromFile(ctx, "../../test-cases/simple-user-task.bpmn")
	then.AssertThat(t, err, is.Nil())
	cp := CallPath{}
	bpmnEngine.AddTaskHandler("user-task", cp.CallPathHandler)

	instance, _ := bpmnEngine.CreateAndRunInstance(ctx, process.GetProcessKey(), nil)

	state, err := instance.GetState(ctx)
	if err != nil {
		return
	}
	then.AssertThat(t, state, is.EqualTo(process_instance.COMPLETED))
	then.AssertThat(t, cp.CallPath, is.EqualTo("user-task"))
}

func Test_user_tasks_can_be_continue(t *testing.T) {
	// setup
	process, err := bpmnEngine.LoadFromFile(ctx, "../../test-cases/simple-user-task.bpmn")
	then.AssertThat(t, err, is.Nil())
	cp := CallPath{}

	instance, _ := bpmnEngine.CreateInstance(ctx, process.GetProcessKey(), nil)

	userConfirm := false
	bpmnEngine.AddTaskHandler("user-task", func(job bpmn_engine_store.IActivatedJob) {
		if userConfirm {
			cp.CallPathHandler(job)
		}
	})
	bpmnEngine.RunOrContinueInstance(ctx, instance.GetInstanceKey())

	userConfirm = true

	bpmnEngine.RunOrContinueInstance(ctx, instance.GetInstanceKey())

	state, err := instance.GetState(ctx)
	if err != nil {
		return
	}
	then.AssertThat(t, state, is.EqualTo(process_instance.COMPLETED))
	then.AssertThat(t, cp.CallPath, is.EqualTo("user-task"))
}
