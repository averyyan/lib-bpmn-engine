package store_test

import (
	"github.com/corbym/gocrest/has"
	"github.com/corbym/gocrest/is"
	"github.com/corbym/gocrest/then"
	"github.com/nitram509/lib-bpmn-engine/pkg/spec/BPMN20/activity"
	"github.com/nitram509/lib-bpmn-engine/pkg/spec/BPMN20/process_instance"
	"testing"
)

func Test_creating_a_process_sets_state_to_READY(t *testing.T) {
	// setup

	// given
	process, _ := bpmnEngine.LoadFromFile(ctx, "../../test-cases/message-intermediate-catch-event.bpmn")

	// when
	pi, _ := bpmnEngine.CreateInstance(ctx, process.GetProcessKey(), nil)
	// then
	state, err := pi.GetState(ctx)
	if err != nil {
		return
	}
	then.AssertThat(t, state, is.EqualTo(process_instance.READY))
}

func Test_running_a_process_sets_state_to_ACTIVE(t *testing.T) {
	// setup

	// given
	process, _ := bpmnEngine.LoadFromFile(ctx, "../../test-cases/message-intermediate-catch-event.bpmn")

	// when
	pi, _ := bpmnEngine.CreateInstance(ctx, process.GetProcessKey(), nil)

	procInst, _ := bpmnEngine.RunOrContinueInstance(ctx, pi.GetInstanceKey())

	// then
	state, err := pi.GetState(ctx)
	if err != nil {
		return
	}
	then.AssertThat(t, state, is.EqualTo(process_instance.ACTIVE).
		Reason("Since the BPMN contains an intermediate catch event, the process instance must be active and can't complete."))
	procInstState, err := procInst.GetState(ctx)
	if err != nil {
		return
	}
	then.AssertThat(t, procInstState, is.EqualTo(process_instance.ACTIVE))
}

func Test_IntermediateCatchEvent_received_message_completes_the_instance(t *testing.T) {
	// setup

	// given
	process, _ := bpmnEngine.LoadFromFile(ctx, "../../test-cases/message-intermediate-catch-event.bpmn")
	pi, _ := bpmnEngine.CreateAndRunInstance(ctx, process.GetProcessKey(), nil)

	// when
	bpmnEngine.PublishEventForInstance(ctx, pi.GetInstanceKey(), "globalMsgRef", nil)
	bpmnEngine.RunOrContinueInstance(ctx, pi.GetInstanceKey())

	// then
	state, err := pi.GetState(ctx)
	if err != nil {
		return
	}
	then.AssertThat(t, state, is.EqualTo(process_instance.COMPLETED))
}

func Test_IntermediateCatchEvent_message_can_be_published_before_running_the_instance(t *testing.T) {
	// setup

	// given
	process, _ := bpmnEngine.LoadFromFile(ctx, "../../test-cases/message-intermediate-catch-event.bpmn")
	pi, _ := bpmnEngine.CreateInstance(ctx, process.GetProcessKey(), nil)

	// when
	bpmnEngine.PublishEventForInstance(ctx, pi.GetInstanceKey(), "globalMsgRef", nil)
	bpmnEngine.RunOrContinueInstance(ctx, pi.GetInstanceKey())

	// then
	state, err := pi.GetState(ctx)
	if err != nil {
		return
	}
	then.AssertThat(t, state, is.EqualTo(process_instance.COMPLETED))
}

func Test_IntermediateCatchEvent_a_catch_event_produces_an_active_subscription(t *testing.T) {
	// setup

	// given
	process, _ := bpmnEngine.LoadFromFile(ctx, "../../test-cases/message-intermediate-catch-event.bpmn")
	instance, err := bpmnEngine.CreateAndRunInstance(ctx, process.GetProcessKey(), nil)
	if err != nil {
		return
	}
	subscriptions, err := instance.FindMessageSubscriptions(ctx)
	if err != nil {
		return
	}
	//subscriptions := bpmnEngine.GetMessageSubscriptions()

	then.AssertThat(t, subscriptions, has.Length(1))
	subscription := subscriptions[0]
	then.AssertThat(t, subscription.GetName(), is.EqualTo("event-1"))
	then.AssertThat(t, subscription.GetElementId(), is.EqualTo("id-1"))
	state, _ := subscription.GetState(ctx)
	then.AssertThat(t, state, is.EqualTo(activity.Active))
}

func Test_Having_IntermediateCatchEvent_and_ServiceTask_in_parallel_the_process_state_is_maintained(t *testing.T) {
	// setup
	cp := CallPath{}

	// given
	process, _ := bpmnEngine.LoadFromFile(ctx, "../../test-cases/message-intermediate-catch-event-and-parallel-tasks.bpmn")
	instance, _ := bpmnEngine.CreateInstance(ctx, process.GetProcessKey(), nil)
	bpmnEngine.AddTaskHandler("task-1", cp.CallPathHandler)
	bpmnEngine.AddTaskHandler("task-2", cp.CallPathHandler)

	// when
	bpmnEngine.RunOrContinueInstance(ctx, instance.GetInstanceKey())

	// then
	state, _ := instance.GetState(ctx)
	then.AssertThat(t, state, is.EqualTo(process_instance.ACTIVE))

	// when
	bpmnEngine.PublishEventForInstance(ctx, instance.GetInstanceKey(), "event-1", nil)
	bpmnEngine.RunOrContinueInstance(ctx, instance.GetInstanceKey())

	// then
	then.AssertThat(t, cp.CallPath, is.EqualTo("task-2,task-1"))
	getState, _ := instance.GetState(ctx)

	then.AssertThat(t, getState, is.EqualTo(process_instance.COMPLETED))
}

func Test_multiple_intermediate_catch_events_possible(t *testing.T) {
	// setup
	cp := CallPath{}

	// given
	process, _ := bpmnEngine.LoadFromFile(ctx, "../../test-cases/message-multiple-intermediate-catch-events.bpmn")
	bpmnEngine.AddTaskHandler("task1", cp.CallPathHandler)
	bpmnEngine.AddTaskHandler("task2", cp.CallPathHandler)
	bpmnEngine.AddTaskHandler("task3", cp.CallPathHandler)
	instance, err := bpmnEngine.CreateAndRunInstance(ctx, process.GetProcessKey(), nil)
	then.AssertThat(t, err, is.Nil())

	// when
	err = bpmnEngine.PublishEventForInstance(ctx, instance.GetInstanceKey(), "msg-event-2", nil)
	then.AssertThat(t, err, is.Nil())
	bpmnEngine.RunOrContinueInstance(ctx, instance.GetInstanceKey())

	// then
	then.AssertThat(t, cp.CallPath, is.EqualTo("task2"))
	// then still active, since there's an implicit fork
	state, _ := instance.GetState(ctx)
	then.AssertThat(t, state, is.EqualTo(process_instance.ACTIVE))
}

func Test_multiple_intermediate_catch_events_implicit_fork_and_merged_COMPLETED(t *testing.T) {
	// setup

	// given
	process, _ := bpmnEngine.LoadFromFile(ctx, "../../test-cases/message-multiple-intermediate-catch-events-merged.bpmn")
	instance, err := bpmnEngine.CreateAndRunInstance(ctx, process.GetProcessKey(), nil)
	then.AssertThat(t, err, is.Nil())

	// when
	err = bpmnEngine.PublishEventForInstance(ctx, instance.GetInstanceKey(), "msg-event-1", nil)
	err = bpmnEngine.PublishEventForInstance(ctx, instance.GetInstanceKey(), "msg-event-2", nil)
	err = bpmnEngine.PublishEventForInstance(ctx, instance.GetInstanceKey(), "msg-event-3", nil)
	then.AssertThat(t, err, is.Nil())
	bpmnEngine.RunOrContinueInstance(ctx, instance.GetInstanceKey())

	// then
	state, err := instance.GetState(ctx)
	if err != nil {
		return
	}
	then.AssertThat(t, state, is.EqualTo(process_instance.COMPLETED))
}

func Test_multiple_intermediate_catch_events_implicit_fork_and_merged_ACTIVE(t *testing.T) {
	// setup

	// given
	process, _ := bpmnEngine.LoadFromFile(ctx, "../../test-cases/message-multiple-intermediate-catch-events-merged.bpmn")
	instance, err := bpmnEngine.CreateAndRunInstance(ctx, process.GetProcessKey(), nil)
	then.AssertThat(t, err, is.Nil())

	// when
	err = bpmnEngine.PublishEventForInstance(ctx, instance.GetInstanceKey(), "msg-event-2", nil)
	then.AssertThat(t, err, is.Nil())
	bpmnEngine.RunOrContinueInstance(ctx, instance.GetInstanceKey())

	// then
	state, err := instance.GetState(ctx)
	if err != nil {
		return
	}
	then.AssertThat(t, state, is.EqualTo(process_instance.ACTIVE))
}

func Test_multiple_intermediate_catch_events_implicit_fork_and_parallel_gateway_COMPLETED(t *testing.T) {
	// setup

	// given
	process, _ := bpmnEngine.LoadFromFile(ctx, "../../test-cases/message-multiple-intermediate-catch-events-parallel.bpmn")
	instance, err := bpmnEngine.CreateAndRunInstance(ctx, process.GetProcessKey(), nil)
	then.AssertThat(t, err, is.Nil())

	// when
	err = bpmnEngine.PublishEventForInstance(ctx, instance.GetInstanceKey(), "msg-event-1", nil)
	err = bpmnEngine.PublishEventForInstance(ctx, instance.GetInstanceKey(), "msg-event-2", nil)
	err = bpmnEngine.PublishEventForInstance(ctx, instance.GetInstanceKey(), "msg-event-3", nil)
	then.AssertThat(t, err, is.Nil())
	bpmnEngine.RunOrContinueInstance(ctx, instance.GetInstanceKey())

	// then
	state, err := instance.GetState(ctx)
	if err != nil {
		return
	}
	then.AssertThat(t, state, is.EqualTo(process_instance.COMPLETED))
}

func Test_multiple_intermediate_catch_events_implicit_fork_and_parallel_gateway_ACTIVE(t *testing.T) {
	// setup

	// given
	process, _ := bpmnEngine.LoadFromFile(ctx, "../../test-cases/message-multiple-intermediate-catch-events-parallel.bpmn")
	instance, err := bpmnEngine.CreateAndRunInstance(ctx, process.GetProcessKey(), nil)
	then.AssertThat(t, err, is.Nil())

	// when
	err = bpmnEngine.PublishEventForInstance(ctx, instance.GetInstanceKey(), "msg-event-2", nil)
	then.AssertThat(t, err, is.Nil())
	bpmnEngine.RunOrContinueInstance(ctx, instance.GetInstanceKey())

	// then
	state, err := instance.GetState(ctx)
	if err != nil {
		return
	}
	then.AssertThat(t, state, is.EqualTo(process_instance.ACTIVE))
}

func Test_multiple_intermediate_catch_events_implicit_fork_and_exclusive_gateway_COMPLETED(t *testing.T) {
	// setup

	// given
	process, _ := bpmnEngine.LoadFromFile(ctx, "../../test-cases/message-multiple-intermediate-catch-events-exclusive.bpmn")
	instance, err := bpmnEngine.CreateAndRunInstance(ctx, process.GetProcessKey(), nil)
	then.AssertThat(t, err, is.Nil())

	// when
	err = bpmnEngine.PublishEventForInstance(ctx, instance.GetInstanceKey(), "msg-event-1", nil)
	err = bpmnEngine.PublishEventForInstance(ctx, instance.GetInstanceKey(), "msg-event-2", nil)
	err = bpmnEngine.PublishEventForInstance(ctx, instance.GetInstanceKey(), "msg-event-3", nil)
	then.AssertThat(t, err, is.Nil())
	bpmnEngine.RunOrContinueInstance(ctx, instance.GetInstanceKey())

	// then
	state, err := instance.GetState(ctx)
	if err != nil {
		return
	}
	then.AssertThat(t, state, is.EqualTo(process_instance.COMPLETED))
}

func Test_multiple_intermediate_catch_events_implicit_fork_and_exclusive_gateway_ACTIVE(t *testing.T) {
	// setup

	// given
	process, _ := bpmnEngine.LoadFromFile(ctx, "../../test-cases/message-multiple-intermediate-catch-events-exclusive.bpmn")
	instance, err := bpmnEngine.CreateAndRunInstance(ctx, process.GetProcessKey(), nil)
	then.AssertThat(t, err, is.Nil())

	// when
	err = bpmnEngine.PublishEventForInstance(ctx, instance.GetInstanceKey(), "msg-event-2", nil)
	then.AssertThat(t, err, is.Nil())
	bpmnEngine.RunOrContinueInstance(ctx, instance.GetInstanceKey())

	// then
	state, err := instance.GetState(ctx)
	if err != nil {
		return
	}
	then.AssertThat(t, state, is.EqualTo(process_instance.ACTIVE))
}

func Test_publishing_a_random_message_does_no_harm(t *testing.T) {
	// setup

	// given
	process, _ := bpmnEngine.LoadFromFile(ctx, "../../test-cases/message-intermediate-catch-event.bpmn")
	instance, err := bpmnEngine.CreateInstance(ctx, process.GetProcessKey(), nil)
	then.AssertThat(t, err, is.Nil())

	// when
	err = bpmnEngine.PublishEventForInstance(ctx, instance.GetInstanceKey(), "random-message", nil)
	then.AssertThat(t, err, is.Nil())
	_, err = bpmnEngine.RunOrContinueInstance(ctx, instance.GetInstanceKey())

	// then
	then.AssertThat(t, err, is.Nil())
	state, err := instance.GetState(ctx)
	if err != nil {
		return
	}
	then.AssertThat(t, state, is.EqualTo(process_instance.ACTIVE))
}

func Test_eventBasedGateway_just_fires_one_event_and_instance_COMPLETED(t *testing.T) {
	// setup
	cp := CallPath{}

	// given
	process, _ := bpmnEngine.LoadFromFile(ctx, "../../test-cases/message-EventBasedGateway.bpmn")
	instance, _ := bpmnEngine.CreateInstance(ctx, process.GetProcessKey(), nil)
	bpmnEngine.AddTaskHandler("task-a", cp.CallPathHandler)
	bpmnEngine.AddTaskHandler("task-b", cp.CallPathHandler)

	// when
	_ = bpmnEngine.PublishEventForInstance(ctx, instance.GetInstanceKey(), "msg-b", nil)
	_, _ = bpmnEngine.RunOrContinueInstance(ctx, instance.GetInstanceKey())

	// then
	then.AssertThat(t, cp.CallPath, is.EqualTo("task-b"))
	state, err := instance.GetState(ctx)
	if err != nil {
		return
	}
	then.AssertThat(t, state, is.EqualTo(process_instance.COMPLETED))
}

func Test_intermediate_message_catch_event_publishes_variables_into_instance(t *testing.T) {
	// setup

	// given
	process, _ := bpmnEngine.LoadFromFile(ctx, "../../test-cases/simple-intermediate-message-catch-event.bpmn")
	instance, _ := bpmnEngine.CreateInstance(ctx, process.GetProcessKey(), nil)

	// when
	vars := map[string]interface{}{"foo": "bar"}
	_ = bpmnEngine.PublishEventForInstance(ctx, instance.GetInstanceKey(), "msg", vars)
	_, _ = bpmnEngine.RunOrContinueInstance(ctx, instance.GetInstanceKey())

	// then
	state, err := instance.GetState(ctx)
	if err != nil {
		return
	}
	then.AssertThat(t, state, is.EqualTo(process_instance.COMPLETED))
	variable, err := instance.GetVariable(ctx, "foo")
	if err != nil {
		return
	}
	then.AssertThat(t, variable, is.EqualTo("bar"))
	getVariable, err := instance.GetVariable(ctx, "mappedFoo")
	if err != nil {
		return
	}
	then.AssertThat(t, getVariable, is.EqualTo("bar"))
}

func Test_intermediate_message_catch_event_output_mapping_failed(t *testing.T) {
	// setup

	// given
	process, _ := bpmnEngine.LoadFromFile(ctx, "../../test-cases/simple-intermediate-message-catch-event-broken.bpmn")
	instance, _ := bpmnEngine.CreateInstance(ctx, process.GetProcessKey(), nil)

	// when
	_ = bpmnEngine.PublishEventForInstance(ctx, instance.GetInstanceKey(), "msg", nil)
	_, _ = bpmnEngine.RunOrContinueInstance(ctx, instance.GetInstanceKey())

	// then
	state, err := instance.GetState(ctx)
	if err != nil {
		return
	}
	then.AssertThat(t, state, is.EqualTo(process_instance.FAILED))
	variable, err := instance.GetVariable(ctx, "mappedFoo")
	if err != nil {
		return
	}
	then.AssertThat(t, variable, is.Nil())
	messageSubscriptions, _ := instance.FindMessageSubscriptions(ctx)
	getState, err := messageSubscriptions[0].GetState(ctx)
	if err != nil {
		return
	}
	then.AssertThat(t, getState, is.EqualTo(activity.Failed))
}
