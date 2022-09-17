package store_test

import (
	"github.com/corbym/gocrest/has"
	"github.com/corbym/gocrest/is"
	"github.com/corbym/gocrest/then"
	"testing"
	"time"
)

func TestEventBasedGatewaySelectsPathWhereTimerOccurs(t *testing.T) {
	// setup
	cp := CallPath{}

	// given
	process, _ := bpmnEngine.LoadFromFile(ctx, "../../test-cases/message-intermediate-timer-event.bpmn")
	bpmnEngine.AddTaskHandler("task-for-message", cp.CallPathHandler)
	bpmnEngine.AddTaskHandler("task-for-timer", cp.CallPathHandler)
	instance, _ := bpmnEngine.CreateAndRunInstance(ctx, process.GetProcessKey(), nil)

	// when
	bpmnEngine.PublishEventForInstance(ctx, instance.GetInstanceKey(), "message", nil)
	bpmnEngine.RunOrContinueInstance(ctx, instance.GetInstanceKey())

	// then
	then.AssertThat(t, cp.CallPath, is.EqualTo("task-for-message"))
}

func TestInvalidTimer_will_stop_continue_execution(t *testing.T) {
	// setup
	cp := CallPath{}

	// given
	process, _ := bpmnEngine.LoadFromFile(ctx, "../../test-cases/message-intermediate-invalid-timer-event.bpmn")
	bpmnEngine.AddTaskHandler("task-for-timer", cp.CallPathHandler)
	instance, _ := bpmnEngine.CreateAndRunInstance(ctx, process.GetProcessKey(), nil)

	// when
	bpmnEngine.PublishEventForInstance(ctx, instance.GetInstanceKey(), "message", nil)
	bpmnEngine.RunOrContinueInstance(ctx, instance.GetInstanceKey())

	// then
	then.AssertThat(t, cp.CallPath, is.EqualTo(""))
}

func TestEventBasedGatewaySelectsPathWhereMessageReceived(t *testing.T) {
	// setup
	cp := CallPath{}

	// given
	process, _ := bpmnEngine.LoadFromFile(ctx, "../../test-cases/message-intermediate-timer-event.bpmn")
	bpmnEngine.AddTaskHandler("task-for-message", cp.CallPathHandler)
	bpmnEngine.AddTaskHandler("task-for-timer", cp.CallPathHandler)
	instance, _ := bpmnEngine.CreateAndRunInstance(ctx, process.GetProcessKey(), nil)

	// when
	time.Sleep((1 * time.Second) + (1 * time.Millisecond))
	bpmnEngine.RunOrContinueInstance(ctx, instance.GetInstanceKey())

	// then
	then.AssertThat(t, cp.CallPath, is.EqualTo("task-for-timer"))
}

func TestEventBasedGatewaySelectsJustOnePath(t *testing.T) {
	// setup
	cp := CallPath{}

	// given
	process, _ := bpmnEngine.LoadFromFile(ctx, "../../test-cases/message-intermediate-timer-event.bpmn")
	bpmnEngine.AddTaskHandler("task-for-message", cp.CallPathHandler)
	bpmnEngine.AddTaskHandler("task-for-timer", cp.CallPathHandler)
	instance, _ := bpmnEngine.CreateAndRunInstance(ctx, process.GetProcessKey(), nil)

	// when
	time.Sleep((1 * time.Second) + (1 * time.Millisecond))
	bpmnEngine.PublishEventForInstance(ctx, instance.GetInstanceKey(), "message", nil)
	bpmnEngine.RunOrContinueInstance(ctx, instance.GetInstanceKey())

	// then
	then.AssertThat(t, cp.CallPath, is.AllOf(
		has.Prefix("task-for"),
		is.Not(is.ValueContaining(","))),
	)
}
