package store

import (
	"context"
	"github.com/nitram509/lib-bpmn-engine/pkg/bpmn_engine_store"
	"time"
)

type catchEvent struct {
	engineStore        bpmn_engine_store.IBpmnEngineStore
	ProcessInstanceKey bpmn_engine_store.IProcessInstanceKey
	ElementInstanceKey bpmn_engine_store.ICatchEventKey
	name               string
	caughtAt           time.Time
	isConsumed         bool
	variables          map[string]interface{}
}

func (ce *catchEvent) GetVariables(ctx context.Context) (map[string]interface{}, error) {
	return ce.engineStore.GetCatchEventVariables(ctx, ce.ProcessInstanceKey, ce.ElementInstanceKey)
}

func (ce *catchEvent) SetConsumed(ctx context.Context, consumed bool) error {
	return ce.engineStore.SetCatchEventConsumed(ctx, ce.ProcessInstanceKey, ce.ElementInstanceKey, consumed)
}

func (ce *catchEvent) GetConsumed(ctx context.Context) (bool, error) {
	return ce.engineStore.GetCatchEventConsumed(ctx, ce.ProcessInstanceKey, ce.ElementInstanceKey)
}

func (ce *catchEvent) GetName() string {
	return ce.name
}
