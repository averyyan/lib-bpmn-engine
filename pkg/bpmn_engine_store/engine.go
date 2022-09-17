package bpmn_engine_store

import (
	"context"
	"errors"
	"fmt"
	"github.com/nitram509/lib-bpmn-engine/pkg/bpmn_engine/exporter"
	"github.com/nitram509/lib-bpmn-engine/pkg/spec/BPMN20"
	"github.com/nitram509/lib-bpmn-engine/pkg/spec/BPMN20/activity"
	"github.com/nitram509/lib-bpmn-engine/pkg/spec/BPMN20/process_instance"
)

type IProcessInfoKey int64
type IProcessInstanceKey int64
type IJobKey int64
type ICatchEventKey int64
type IMessageSubscriptionKey int64
type ITimerKey int64
type IScheduledFlowKey string

type IBpmnEngine interface {
	GetStore() IBpmnEngineStore
	GenerateKey() int64
}

const continueNextElement = true

// IBpmnEngineStore context.Context is for db ACID

func New(name string, store IBpmnEngineStore) BpmnEngineState {
	snowflakeIdGenerator := initializeSnowflakeIdGenerator()
	return BpmnEngineState{
		name:        name,
		engineStore: store,
		handlers:    map[string]func(job IActivatedJob){},
		snowflake:   snowflakeIdGenerator,
		exporters:   []exporter.EventExporter{},
	}
}

// CreateInstance creates a new instance for a process with given processKey
// will return (nil, nil), when no process with given was found
func (state *BpmnEngineState) CreateInstance(ctx context.Context, processKey IProcessInfoKey, variableContext map[string]interface{}) (IProcessInstanceInfo, error) {
	return state.GetStore().CreateProcessInstance(
		ctx,
		processKey,
		IProcessInstanceKey(state.generateKey()),
		variableContext,
	)
}

func (state *BpmnEngineState) CreateAndRunInstance(ctx context.Context, processKey IProcessInfoKey, variableContext map[string]interface{}) (IProcessInstanceInfo, error) {
	instance, err := state.CreateInstance(ctx, processKey, variableContext)
	if err != nil {
		return nil, err
	}
	if instance == nil {
		return nil, errors.New(fmt.Sprint("can't find process with processKey=", processKey, "."))
	}
	err = state.run(ctx, instance)
	return instance, err
}

// RunOrContinueInstance runs or continues a process instance by a given processInstanceKey.
// returns the process instances, when found
// does nothing, if process is already in ProcessInstanceCompleted State
// returns nil, when no process instance was found
// Additionally, every time this method is called, former completed instances are 'garbage collected'.
func (state *BpmnEngineState) RunOrContinueInstance(ctx context.Context, processInstanceKey IProcessInstanceKey) (IProcessInstanceInfo, error) {
	pi, err := state.engineStore.FindProcessInstanceInfo(ctx, processInstanceKey)
	if err != nil {
		return nil, err
	}
	return pi, state.run(ctx, pi)
}

func (state *BpmnEngineState) run(ctx context.Context, instance IProcessInstanceInfo) (err error) {
	var instanceGetStateError error
	var instanceState process_instance.State
	type queueElement struct {
		inboundFlowId string
		baseElement   BPMN20.BaseElement
	}
	queue := make([]queueElement, 0)
	process := instance.GetProcessInfo()
	instanceState, instanceGetStateError = instance.GetState(ctx)
	if instanceGetStateError != nil {
		return instanceGetStateError
	}
	switch instanceState {
	case process_instance.READY:
		// use start events to start the instance
		for _, event := range process.GetDefinitions().Process.StartEvents {
			queue = append(queue, queueElement{
				inboundFlowId: "",
				baseElement:   event,
			})
		}
		if err := instance.SetState(ctx, process_instance.ACTIVE); err != nil {
			return err
		}
	case process_instance.ACTIVE:

		userTasks := findActiveUserTasksForContinuation(ctx, instance)
		for _, userTask := range userTasks {
			queue = append(queue, queueElement{
				inboundFlowId: "",
				baseElement:   *userTask,
			})
		}
		intermediateCatchEvents := findIntermediateCatchEventsForContinuation(ctx, instance)
		for _, ice := range intermediateCatchEvents {
			queue = append(queue, queueElement{
				inboundFlowId: "",
				baseElement:   *ice,
			})
		}
	case process_instance.COMPLETED:
		return nil
	case process_instance.FAILED:
		return nil
	default:
		panic(any("Unknown process instance state."))
	}

	for len(queue) > 0 {
		element := queue[0].baseElement
		inboundFlowId := queue[0].inboundFlowId
		queue = queue[1:]

		continueNextElement := state.handleElement(ctx, instance, element)

		if continueNextElement {
			state.exportElementEvent(process, instance, element, exporter.ElementCompleted)

			if inboundFlowId != "" {
				err := instance.RemoveScheduledFlow(ctx, IScheduledFlowKey(inboundFlowId))
				if err != nil {
					return err
				}
			}
			definitions := process.GetDefinitions()
			nextFlows := BPMN20.FindSequenceFlows(&definitions.Process.SequenceFlows, element.GetOutgoingAssociation())
			if element.GetType() == BPMN20.ExclusiveGateway {
				variableContext, getVariableContextErr := instance.GetVariableContext(ctx)
				if getVariableContextErr != nil {
					return getVariableContextErr
				}
				nextFlows, err = exclusivelyFilterByConditionExpression(nextFlows, variableContext)
				if err != nil {
					if err := instance.SetState(ctx, process_instance.FAILED); err != nil {
						return err
					}
					break
				}
			}
			for _, flow := range nextFlows {

				state.exportSequenceFlowEvent(process, instance, flow)

				// TODO: create test for that
				// if len(flows) < 1 {
				//	panic(fmt.Sprintf("Can't find 'sequenceFlow' element with ID=%s. "+
				//		"This is likely because your BPMN is invalid.", flows[0]))
				// }

				if err := instance.CreateScheduledFlow(ctx, IScheduledFlowKey(flow.Id)); err != nil {
					return err
				}
				baseElements := BPMN20.FindBaseElementsById(process.GetDefinitions(), flow.TargetRef)
				// TODO: create test for that
				// if len(baseElements) < 1 {
				//	panic(fmt.Sprintf("Can't find flow element with ID=%s. "+
				//		"This is likely because there are elements in the definition, "+
				//		"which this engine does not support (yet).", flow.Id))
				// }
				targetBaseElement := baseElements[0]
				queue = append(queue, queueElement{
					inboundFlowId: flow.Id,
					baseElement:   targetBaseElement,
				})
			}
		}
	}
	instanceState, instanceGetStateError = instance.GetState(ctx)
	if instanceGetStateError != nil {
		return instanceGetStateError
	}
	if instanceState == process_instance.COMPLETED || instanceState == process_instance.FAILED {
		// TODO need to send failed state
		state.exportEndProcessEvent(process, instance)
	}
	return err
}

func (state *BpmnEngineState) handleElement(ctx context.Context, instance IProcessInstanceInfo, element BPMN20.BaseElement) bool {
	process := instance.GetProcessInfo()

	state.exportElementEvent(process, instance, element, exporter.ElementActivated)

	switch element.GetType() {
	case BPMN20.StartEvent:
		return true
	case BPMN20.ServiceTask:
		taskElement := element.(BPMN20.TaskElement)
		return state.handleServiceTask(ctx, instance, taskElement)
	case BPMN20.UserTask:
		taskElement := element.(BPMN20.TaskElement)
		return state.handleUserTask(ctx, instance, taskElement)
	case BPMN20.ParallelGateway:
		return state.handleParallelGateway(ctx, instance, element)
	case BPMN20.EndEvent:
		state.handleEndEvent(ctx, instance)
		state.exportElementEvent(process, instance, element, exporter.ElementCompleted) // special case here, to end the instance
		return false
	case BPMN20.IntermediateCatchEvent:
		return state.handleIntermediateCatchEvent(ctx, instance, element.(BPMN20.TIntermediateCatchEvent))
	case BPMN20.EventBasedGateway:
		// TODO improve precondition tests
		// simply proceed
		return true
	default:
		// do nothing
		// TODO: should we print a warning?
	}
	return true
}

func (state *BpmnEngineState) handleIntermediateCatchEvent(ctx context.Context, instance IProcessInstanceInfo, ice BPMN20.TIntermediateCatchEvent) bool {
	if ice.MessageEventDefinition.Id != "" {
		return state.handleIntermediateMessageCatchEvent(ctx, instance, ice)
	}
	if ice.TimerEventDefinition.Id != "" {
		return state.handleIntermediateTimerCatchEvent(ctx, instance, ice)
	}
	return false
}

func (state *BpmnEngineState) handleEndEvent(ctx context.Context, instance IProcessInstanceInfo) {
	completedJobs := true
	jobs, err := instance.FindJobs(ctx)
	// TODO if this is db err must stop the process
	if err != nil {
		jobs = []IJob{}
	}
	for _, job := range jobs {
		jobState, err := job.GetState(ctx)
		if err != nil {
			return
		}
		if jobState == activity.Ready || jobState == activity.Active {
			completedJobs = false
			break
		}
	}
	if completedJobs && !hasActiveSubscriptions(ctx, instance) {
		if err := instance.SetState(ctx, process_instance.COMPLETED); err != nil {
			return
		}
	}
}

func hasActiveSubscriptions(ctx context.Context, instance IProcessInstanceInfo) bool {
	process := instance.GetProcessInfo()
	activeSubscriptions := map[string]bool{}
	subscriptions, err := instance.FindMessageSubscriptions(ctx)
	if err != nil {
		return false
	}
	for _, ms := range subscriptions {
		msState, err := ms.GetState(ctx)
		if err != nil {
			return false
		}
		activeSubscriptions[ms.GetElementId()] = msState == activity.Ready || msState == activity.Active
	}
	// eliminate the active subscriptions, which are from one 'parent' EventBasedGateway
	for _, gateway := range process.GetDefinitions().Process.EventBasedGateway {
		definitions := process.GetDefinitions()
		flows := BPMN20.FindSequenceFlows(&definitions.Process.SequenceFlows, gateway.OutgoingAssociation)
		isOneEventCompleted := true
		for _, flow := range flows {
			isOneEventCompleted = isOneEventCompleted && !activeSubscriptions[flow.TargetRef]
		}
		for _, flow := range flows {
			activeSubscriptions[flow.TargetRef] = isOneEventCompleted
		}
	}
	for _, v := range activeSubscriptions {
		if v {
			return true
		}
	}
	return false
}

func (state *BpmnEngineState) handleParallelGateway(ctx context.Context, instance IProcessInstanceInfo, element BPMN20.BaseElement) bool {
	// check incoming flows, if ready, then continue
	allInboundsAreScheduled := true
	for _, inFlowId := range element.GetIncomingAssociation() {
		isExistFlow, err := instance.IsExistScheduledFlow(ctx, IScheduledFlowKey(inFlowId))
		if err != nil {
			return false
		}
		allInboundsAreScheduled = isExistFlow && allInboundsAreScheduled
	}
	return allInboundsAreScheduled
}

func findActiveUserTasksForContinuation(ctx context.Context, instance IProcessInstanceInfo) (ret []*BPMN20.BaseElement) {
	process := instance.GetProcessInfo()
	for _, userTask := range process.GetDefinitions().Process.UserTasks {
		jobs, err := instance.FindJobs(ctx)
		if err != nil {
			return nil
		}
		for _, job := range jobs {
			jobState, err := job.GetState(ctx)
			if err != nil {
				return nil
			}
			if jobState == activity.Active && job.GetElementId() == userTask.GetId() {
				_userTask := BPMN20.BaseElement(userTask)
				ret = append(ret, &_userTask)
			}
		}

	}
	return ret
}

func findIntermediateCatchEventsForContinuation(ctx context.Context, instance IProcessInstanceInfo) (ret []*BPMN20.BaseElement) {
	process := instance.GetProcessInfo()
	messageRef2IntermediateCatchEventMapping := map[string]BPMN20.BaseElement{}
	for _, ice := range process.GetDefinitions().Process.IntermediateCatchEvent {
		messageRef2IntermediateCatchEventMapping[ice.MessageEventDefinition.MessageRef] = ice
	}
	caughtEvents, err := instance.FindCatchEvents(ctx)
	if err != nil {
		return nil
	}
	for _, caughtEvent := range caughtEvents {
		consumed, err := caughtEvent.GetConsumed(ctx)
		if err != nil {
			return nil
		}
		if consumed == true {
			// skip consumed ones
			continue
		}
		for _, msg := range process.GetDefinitions().Messages {
			// find the matching message definition
			if msg.Name == caughtEvent.GetName() {
				// find potential event definitions
				event := messageRef2IntermediateCatchEventMapping[msg.Id]
				if hasActiveMessageSubscriptionForId(ctx, instance, event.GetId()) {
					ret = append(ret, &event)
				}
			}
		}
	}
	ice := checkDueTimersAndFindIntermediateCatchEvent(ctx, instance)
	if ice != nil {
		ret = append(ret, ice)
	}
	return eliminateEventsWhichComeFromTheSameGateway(process.GetDefinitions(), ret)
}

func hasActiveMessageSubscriptionForId(ctx context.Context, instance IProcessInstanceInfo, id string) bool {
	messageSubscriptions, err := instance.FindMessageSubscriptions(ctx)

	if err != nil {
		return false
	}
	for _, subscription := range messageSubscriptions {
		state, err := subscription.GetState(ctx)
		if err != nil {
			return false
		}
		if id == subscription.GetElementId() && (state == activity.Ready || state == activity.Active) {
			return true
		}
	}
	return false
}

func eliminateEventsWhichComeFromTheSameGateway(definitions BPMN20.TDefinitions, events []*BPMN20.BaseElement) (ret []*BPMN20.BaseElement) {
	// a bubble-sort-like approach to find elements, which have the same incoming association
	for len(events) > 0 {
		event := events[0]
		events = events[1:]
		if event == nil {
			continue
		}
		ret = append(ret, event)
		for i := 0; i < len(events); i++ {
			if haveEqualInboundBaseElement(definitions, event, events[i]) && inboundIsEventBasedGateway(definitions, event) {
				events[i] = nil
			}
		}
	}

	return ret
}

func inboundIsEventBasedGateway(definitions BPMN20.TDefinitions, event *BPMN20.BaseElement) bool {
	ref := BPMN20.FindSourceRefs(definitions.Process.SequenceFlows, (*event).GetIncomingAssociation()[0])[0]
	baseElement := BPMN20.FindBaseElementsById(definitions, ref)[0]
	return baseElement.GetType() == BPMN20.EventBasedGateway
}

func haveEqualInboundBaseElement(definitions BPMN20.TDefinitions, event1 *BPMN20.BaseElement, event2 *BPMN20.BaseElement) bool {
	if event1 == nil || event2 == nil {
		return false
	}
	checkOnlyOneAssociationOrPanic(event1)
	checkOnlyOneAssociationOrPanic(event2)
	ref1 := BPMN20.FindSourceRefs(definitions.Process.SequenceFlows, (*event1).GetIncomingAssociation()[0])[0]
	ref2 := BPMN20.FindSourceRefs(definitions.Process.SequenceFlows, (*event2).GetIncomingAssociation()[0])[0]
	baseElement1 := BPMN20.FindBaseElementsById(definitions, ref1)[0]
	baseElement2 := BPMN20.FindBaseElementsById(definitions, ref2)[0]
	return baseElement1.GetId() == baseElement2.GetId()
}

func checkOnlyOneAssociationOrPanic(event *BPMN20.BaseElement) {
	if len((*event).GetIncomingAssociation()) != 1 {
		panic(any(fmt.Sprintf("Element with id=%s has %d incoming associations, but only 1 is supported by this engine.",
			(*event).GetId(), len((*event).GetIncomingAssociation()))))
	}
}
