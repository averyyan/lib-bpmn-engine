package store

import (
	"context"
	"github.com/nitram509/lib-bpmn-engine/pkg/bpmn_engine_store"
	"github.com/nitram509/lib-bpmn-engine/pkg/spec/BPMN20/activity"
	"time"
)

type MessageSubscription struct {
	engineState        bpmn_engine_store.IBpmnEngine
	ElementId          string
	ElementInstanceKey bpmn_engine_store.IMessageSubscriptionKey
	ProcessInstanceKey bpmn_engine_store.IProcessInstanceKey
	Name               string
	State              activity.LifecycleState
	CreatedAt          time.Time
}

func (ms *MessageSubscription) GetProcessInstanceKey() bpmn_engine_store.IProcessInstanceKey {
	return ms.ProcessInstanceKey
}

func (ms *MessageSubscription) SetState(ctx context.Context, state activity.LifecycleState) error {
	return ms.engineState.GetStore().SetMessageSubscriptionState(ctx, ms.ProcessInstanceKey, ms.ElementInstanceKey, state)
}

func (ms *MessageSubscription) GetState(ctx context.Context) (activity.LifecycleState, error) {
	return ms.engineState.GetStore().GetMessageSubscriptionState(ctx, ms.ProcessInstanceKey, ms.ElementInstanceKey)
}

func (ms MessageSubscription) GetElementId() string {
	return ms.ElementId
}
