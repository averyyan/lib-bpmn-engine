package bpmn_engine_store

import (
	"context"
	"fmt"
	"github.com/nitram509/lib-bpmn-engine/pkg/spec/BPMN20"
	"github.com/nitram509/lib-bpmn-engine/pkg/spec/BPMN20/activity"
	"github.com/nitram509/lib-bpmn-engine/pkg/spec/BPMN20/process_instance"
)

// PublishEventForInstance publishes a message with a given name and also adds variables to the process instance, which fetches this event
func (state *BpmnEngineState) PublishEventForInstance(ctx context.Context, processInstanceKey IProcessInstanceKey, messageName string, variables map[string]interface{}) error {
	processInstance, err := state.engineStore.FindProcessInstanceInfo(ctx, processInstanceKey)
	if err != nil {
		return err
	}
	if processInstance != nil {
		if _, err := processInstance.CreateCatchEvent(ctx, messageName, variables); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("no process instance with key=%d found", processInstanceKey)
	}
	return nil
}

func (state *BpmnEngineState) handleIntermediateMessageCatchEvent(ctx context.Context, instance IProcessInstanceInfo, ice BPMN20.TIntermediateCatchEvent) bool {

	subscriptions, err := instance.FindMessageSubscriptions(ctx)
	if err != nil {
		subscriptions = []IMessageSubscription{}
	}
	messageSubscription := findMatchingActiveSubscriptions(ctx, subscriptions, ice.Id)

	if messageSubscription == nil {
		messageSubscription, err = instance.CreateMessageSubscription(
			ctx,
			ice,
		)
		if err != nil {
			return false
		}
	}

	messages := instance.GetProcessInfo().GetDefinitions().Messages

	caughtEvent := findMatchingCaughtEvent(ctx, &messages, instance, ice)

	if caughtEvent != nil {
		if err := messageSubscription.SetState(ctx, activity.Completed); err != nil {
			return false
		}
		if err := caughtEvent.SetConsumed(ctx, true); err != nil {
			return false
		}
		variables, err := caughtEvent.GetVariables(ctx)
		if err != nil {
			return false
		}
		for k, v := range variables {
			if err := instance.SetVariable(ctx, k, v); err != nil {
				return false
			}
		}
		if err := evaluateVariableMapping(ctx, instance, ice.Output); err != nil {
			if err := messageSubscription.SetState(ctx, activity.Failed); err != nil {
				return false
			}
			if err := instance.SetState(ctx, process_instance.FAILED); err != nil {
				return false
			}
			return false
		}
		return continueNextElement
	}
	return !continueNextElement
}

//func (state *BpmnEngineState) findMessagesByProcessKey(ctx context.Context, instance IProcessInstanceInfo) *[]BPMN20.TMessage {
//	process, err := state.GetStore().FindProcessInfo(ctx, processKey)
//	if err != nil {
//		return nil
//	}
//	definitions := process.GetDefinitions()
//	return &definitions.Messages
//}

// find first matching catchEvent
func findMatchingCaughtEvent(ctx context.Context, messages *[]BPMN20.TMessage, instance IProcessInstanceInfo, ice BPMN20.TIntermediateCatchEvent) ICatchEvent {
	msgName := findMessageNameById(messages, ice.MessageEventDefinition.MessageRef)
	caughtEvents, err := instance.FindCatchEvents(ctx)
	if err != nil {
		return nil
	}
	for i := 0; i < len(caughtEvents); i++ {
		var caughtEvent = caughtEvents[i]
		consumed, err := caughtEvent.GetConsumed(ctx)
		if err != nil {
			return nil
		}
		if !consumed && msgName == caughtEvent.GetName() {
			return caughtEvent
		}
	}
	return nil
}

func findMessageNameById(messages *[]BPMN20.TMessage, msgId string) string {
	for _, message := range *messages {
		if message.Id == msgId {
			return message.Name
		}
	}
	return ""
}

func findMatchingActiveSubscriptions(ctx context.Context, messageSubscriptions []IMessageSubscription, id string) IMessageSubscription {
	var existingSubscription IMessageSubscription
	for _, ms := range messageSubscriptions {
		state, err := ms.GetState(ctx)
		if err != nil {
			return nil
		}
		if state == activity.Active && ms.GetElementId() == id {
			existingSubscription = ms
			return existingSubscription
		}
	}
	return nil
}
