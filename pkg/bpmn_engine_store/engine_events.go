package bpmn_engine_store

import (
	"context"
	"github.com/nitram509/lib-bpmn-engine/pkg/spec/BPMN20"
	"github.com/nitram509/lib-bpmn-engine/pkg/spec/BPMN20/activity"
	"github.com/nitram509/lib-bpmn-engine/pkg/spec/BPMN20/process_instance"
)

func (state *BpmnEngineState) handleIntermediateMessageCatchEvent(ctx context.Context, instance IProcessInstanceInfo, ice BPMN20.TIntermediateCatchEvent) bool {

	process := instance.GetProcessInfo()
	subscriptions, err := instance.FindMessageSubscriptions(ctx)
	if err != nil {
		return false
	}

	messageSubscription := findMatchingActiveSubscriptions(ctx, subscriptions, ice.Id)

	if messageSubscription == nil {
		messageSubscription, err = state.GetStore().CreateMessageSubscription(
			ctx,
			state,
			instance.GetInstanceKey(),
			IMessageSubscriptionKey(state.generateKey()),
			ice,
		)
		if err != nil {
			return false
		}
	}

	messages := state.findMessagesByProcessKey(ctx, process.GetProcessKey())

	caughtEvent := findMatchingCaughtEvent(ctx, messages, instance, ice)

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

func (state *BpmnEngineState) findMessagesByProcessKey(ctx context.Context, processKey IProcessInfoKey) *[]BPMN20.TMessage {
	process, err := state.GetStore().FindProcessInfo(ctx, processKey)
	if err != nil {
		return nil
	}
	definitions := process.GetDefinitions()
	return &definitions.Messages
}

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
