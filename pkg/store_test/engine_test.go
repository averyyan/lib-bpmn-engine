package store_test

import (
	"context"
	"github.com/corbym/gocrest/is"
	"github.com/corbym/gocrest/then"
	"github.com/nitram509/lib-bpmn-engine/pkg/bpmn_engine_store"
	"github.com/nitram509/lib-bpmn-engine/pkg/store"
	"testing"
	"time"
)

var bpmnEngine = bpmn_engine_store.New("name", store.NewMemoryStore())
var ctx = context.Background()

type CallPath struct {
	CallPath string
}

func (callPath *CallPath) CallPathHandler(job bpmn_engine_store.IActivatedJob) {
	if len(callPath.CallPath) > 0 {
		callPath.CallPath += ","
	}
	callPath.CallPath += job.GetElementId()
	job.Complete(context.TODO())
}

func TestAllInterfacesImplemented(t *testing.T) {
	var _ bpmn_engine_store.IBpmnEngine = &bpmn_engine_store.BpmnEngineState{}
}

func TestRegisterHandlerByTaskIdGetsCalled(t *testing.T) {
	// setup
	ctx := context.Background()
	process, _ := bpmnEngine.LoadFromFile(ctx, "../../test-cases/simple_task.bpmn")
	wasCalled := false
	handler := func(job bpmn_engine_store.IActivatedJob) {
		wasCalled = true
		job.Complete(ctx)
	}

	// given
	bpmnEngine.AddTaskHandler("id", handler)

	// when
	bpmnEngine.CreateAndRunInstance(ctx, process.GetProcessKey(), nil)

	then.AssertThat(t, wasCalled, is.True())
}

func TestRegisteredHandlerCanMutateVariableContext(t *testing.T) {
	// setup
	variableName := "variable_name"
	taskId := "id"
	process, _ := bpmnEngine.LoadFromFile(ctx, "../../test-cases/simple_task.bpmn")
	variableContext := make(map[string]interface{}, 1)
	variableContext[variableName] = "oldVal"

	handler := func(job bpmn_engine_store.IActivatedJob) {
		v, _ := job.GetVariable(ctx, variableName)
		then.AssertThat(t, v, is.EqualTo("oldVal").Reason("one should be able to read variables"))
		job.SetVariable(ctx, variableName, "newVal")
		job.Complete(ctx)
	}

	// given
	bpmnEngine.AddTaskHandler(taskId, handler)

	// when
	instance, _ := bpmnEngine.CreateAndRunInstance(ctx, process.GetProcessKey(), variableContext)

	variable, _ := instance.GetVariable(ctx, variableName)

	// then
	then.AssertThat(t, variable, is.EqualTo("newVal"))
}

func TestMetadataIsGivenFromLoadedXmlFile(t *testing.T) {
	// setup
	metadata, _ := bpmnEngine.LoadFromFile(ctx, "../../test-cases/simple_task.bpmn")

	then.AssertThat(t, metadata.GetVersion(), is.EqualTo(int32(1)))
	then.AssertThat(t, metadata.GetProcessKey(), is.GreaterThan(1))
	then.AssertThat(t, metadata.GetBpmnProcessId(), is.EqualTo("Simple_Task_Process"))
}

func TestLoadingTheSameFileWillNotIncreaseTheVersionNorChangeTheProcessKey(t *testing.T) {
	// setup

	metadata, _ := bpmnEngine.LoadFromFile(ctx, "../../test-cases/simple_task.bpmn")
	keyOne := metadata.GetProcessKey()
	then.AssertThat(t, metadata.GetVersion(), is.EqualTo(int32(1)))

	metadata, _ = bpmnEngine.LoadFromFile(ctx, "../../test-cases/simple_task.bpmn")
	keyTwo := metadata.GetProcessKey()
	then.AssertThat(t, metadata.GetVersion(), is.EqualTo(int32(1)))

	then.AssertThat(t, keyOne, is.EqualTo(keyTwo))
}

func TestLoadingTheSameProcessWithModificationWillCreateNewVersion(t *testing.T) {
	// setup

	process1, _ := bpmnEngine.LoadFromFile(ctx, "../../test-cases/simple_task.bpmn")
	process2, _ := bpmnEngine.LoadFromFile(ctx, "../../test-cases/simple_task_modified_taskId.bpmn")
	process3, _ := bpmnEngine.LoadFromFile(ctx, "../../test-cases/simple_task.bpmn")

	then.AssertThat(t, process1.GetBpmnProcessId(), is.EqualTo(process2.GetBpmnProcessId()).Reason("both prepared files should have equal IDs"))
	then.AssertThat(t, int64(process2.GetProcessKey()), is.GreaterThan(int64(process1.GetProcessKey())).Reason("Because later created"))
	then.AssertThat(t, process3.GetProcessKey(), is.EqualTo(process1.GetProcessKey()).Reason("Same processKey return for same input file, means already registered"))

	then.AssertThat(t, process1.GetVersion(), is.EqualTo(int32(1)))
	then.AssertThat(t, process2.GetVersion(), is.EqualTo(int32(2)))
	then.AssertThat(t, process3.GetVersion(), is.EqualTo(int32(1)))

	then.AssertThat(t, process1.GetProcessKey(), is.Not(is.EqualTo(process2.GetProcessKey())))
}

func TestMultipleInstancesCanBeCreated(t *testing.T) {
	// setup
	beforeCreation := time.Now()

	// given
	process, _ := bpmnEngine.LoadFromFile(ctx, "../../test-cases/simple_task.bpmn")

	// when
	instance1, _ := bpmnEngine.CreateInstance(ctx, process.GetProcessKey(), nil)
	instance2, _ := bpmnEngine.CreateInstance(ctx, process.GetProcessKey(), nil)

	// then
	then.AssertThat(t, instance1.GetCreatedAt().UnixNano(), is.GreaterThanOrEqualTo(beforeCreation.UnixNano()).Reason("make sure we have creation time set"))
	then.AssertThat(t, instance1.GetProcessInfo().GetProcessKey(), is.EqualTo(instance2.GetProcessInfo().GetProcessKey()))
	then.AssertThat(t, int64(instance2.GetInstanceKey()), is.GreaterThan(int64(instance1.GetInstanceKey())).Reason("Because later created"))
}

func TestSimpleAndUncontrolledForkingTwoTasks(t *testing.T) {
	// setup
	cp := CallPath{}

	// given
	process, _ := bpmnEngine.LoadFromFile(ctx, "../../test-cases/forked-flow.bpmn")
	bpmnEngine.AddTaskHandler("id-a-1", cp.CallPathHandler)
	bpmnEngine.AddTaskHandler("id-b-1", cp.CallPathHandler)
	bpmnEngine.AddTaskHandler("id-b-2", cp.CallPathHandler)

	// when
	bpmnEngine.CreateAndRunInstance(ctx, process.GetProcessKey(), nil)

	// then
	then.AssertThat(t, cp.CallPath, is.EqualTo("id-a-1,id-b-1,id-b-2"))
}

func TestParallelGateWayTwoTasks(t *testing.T) {
	// setup
	cp := CallPath{}

	// given
	process, _ := bpmnEngine.LoadFromFile(ctx, "../../test-cases/parallel-gateway-flow.bpmn")
	bpmnEngine.AddTaskHandler("id-a-1", cp.CallPathHandler)
	bpmnEngine.AddTaskHandler("id-b-1", cp.CallPathHandler)
	bpmnEngine.AddTaskHandler("id-b-2", cp.CallPathHandler)

	// when
	bpmnEngine.CreateAndRunInstance(ctx, process.GetProcessKey(), nil)

	// then
	then.AssertThat(t, cp.CallPath, is.EqualTo("id-a-1,id-b-1,id-b-2"))
}
