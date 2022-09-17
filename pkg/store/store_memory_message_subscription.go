package store

import (
	"context"
	"errors"
	"fmt"
	"github.com/nitram509/lib-bpmn-engine/pkg/bpmn_engine_store"
	"github.com/nitram509/lib-bpmn-engine/pkg/spec/BPMN20"
	"github.com/nitram509/lib-bpmn-engine/pkg/spec/BPMN20/activity"
	"time"
)

func (store *EngineMemoryStore) findMessageSubscriptions(processInstanceKey bpmn_engine_store.IProcessInstanceKey) (map[bpmn_engine_store.IMessageSubscriptionKey]*MessageSubscription, error) {
	messageSubscriptions, isExist := store.messageSubscriptions[processInstanceKey]
	if !isExist {
		return nil, errors.New(fmt.Sprintf("messageSubscriptions: ProcessInstanceKey=%d  is not exist", processInstanceKey))
	}
	return messageSubscriptions, nil
}

func (store *EngineMemoryStore) findMessageSubscription(processInstanceKey bpmn_engine_store.IProcessInstanceKey, messageSubscriptionKey bpmn_engine_store.IMessageSubscriptionKey) (*MessageSubscription, error) {
	subscriptions, err := store.findMessageSubscriptions(processInstanceKey)
	if err != nil {
		return nil, err
	}
	messageSubscription, isExist := subscriptions[messageSubscriptionKey]
	if !isExist {
		return nil, errors.New(fmt.Sprintf("messageSubscription: ProcessInstanceKey=%d && messageSubscriptionKey=%d  is not exist", processInstanceKey, messageSubscriptionKey))
	}
	return messageSubscription, nil
}
func (store *EngineMemoryStore) GetMessageSubscriptionState(ctx context.Context, processInstanceKey bpmn_engine_store.IProcessInstanceKey, messageSubscriptionKey bpmn_engine_store.IMessageSubscriptionKey) (activity.LifecycleState, error) {
	subscription, err := store.findMessageSubscription(processInstanceKey, messageSubscriptionKey)
	if err != nil {
		return "", err
	}
	return subscription.State, nil
}
func (store *EngineMemoryStore) SetMessageSubscriptionState(ctx context.Context, processInstanceKey bpmn_engine_store.IProcessInstanceKey, messageSubscriptionKey bpmn_engine_store.IMessageSubscriptionKey, state activity.LifecycleState) error {
	subscription, err := store.findMessageSubscription(processInstanceKey, messageSubscriptionKey)
	if err != nil {
		return err
	}
	subscription.State = state
	return nil
}

func (store *EngineMemoryStore) CreateMessageSubscription(
	ctx context.Context,
	engineStore bpmn_engine_store.IBpmnEngineStore,
	processInstanceKey bpmn_engine_store.IProcessInstanceKey,
	elementInstanceKey bpmn_engine_store.IMessageSubscriptionKey,
	ice BPMN20.TIntermediateCatchEvent,
) (bpmn_engine_store.IMessageSubscription, error) {
	if store.messageSubscriptions[processInstanceKey] == nil {
		store.messageSubscriptions[processInstanceKey] = map[bpmn_engine_store.IMessageSubscriptionKey]*MessageSubscription{}
	}
	messageSubscription := &MessageSubscription{
		engineStore:        engineStore,
		ElementId:          ice.Id,
		ElementInstanceKey: elementInstanceKey,
		ProcessInstanceKey: processInstanceKey,
		Name:               ice.Name,
		CreatedAt:          time.Now(),
		State:              activity.Active,
	}
	store.messageSubscriptions[processInstanceKey][elementInstanceKey] = messageSubscription
	return store.messageSubscriptions[processInstanceKey][elementInstanceKey], nil
}
