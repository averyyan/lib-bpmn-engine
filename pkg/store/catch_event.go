package store

import (
	"context"
	"github.com/nitram509/lib-bpmn-engine/pkg/bpmn_engine_store"
	"time"
)

type catchEvent struct {
	engineState        bpmn_engine_store.IBpmnEngine
	ProcessInstanceKey bpmn_engine_store.IProcessInstanceKey
	key                bpmn_engine_store.ICatchEventKey
	name               string
	caughtAt           time.Time
	isConsumed         bool
	variables          map[string]interface{}
}

func (ce *catchEvent) GetVariables(ctx context.Context) (map[string]interface{}, error) {
	return ce.engineState.GetStore().GetCatchEventVariables(ctx, ce.ProcessInstanceKey, ce.key)
}

func (ce *catchEvent) SetConsumed(ctx context.Context, consumed bool) error {
	return ce.engineState.GetStore().SetCatchEventConsumed(ctx, ce.ProcessInstanceKey, ce.key, consumed)
}

func (ce *catchEvent) GetConsumed(ctx context.Context) (bool, error) {
	return ce.engineState.GetStore().GetCatchEventConsumed(ctx, ce.ProcessInstanceKey, ce.key)
}

func (ce *catchEvent) GetName() string {
	return ce.name
}
