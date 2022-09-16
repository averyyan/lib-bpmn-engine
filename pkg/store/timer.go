package store

import (
	"context"
	"github.com/nitram509/lib-bpmn-engine/pkg/bpmn_engine"
	"github.com/nitram509/lib-bpmn-engine/pkg/bpmn_engine_store"
	"time"
)

type Timer struct {
	engineState        bpmn_engine_store.IBpmnEngine
	ProcessInstanceKey bpmn_engine_store.IProcessInstanceKey
	ElementInstanceKey bpmn_engine_store.ITimerKey
	ProcessKey         bpmn_engine_store.IProcessInfoKey
	ElementId          string
	State              bpmn_engine.TimerState
	CreatedAt          time.Time
	DueAt              time.Time
	Duration           time.Duration
}

func (timer *Timer) SetState(ctx context.Context, state bpmn_engine.TimerState) error {
	return timer.engineState.GetStore().SetTimerState(ctx, timer.ProcessInstanceKey, timer.ElementInstanceKey, state)
}

func (timer *Timer) GetDueAt() time.Time {
	return timer.DueAt
}

func (timer *Timer) GetState(ctx context.Context) (bpmn_engine.TimerState, error) {
	return timer.engineState.GetStore().GetTimerState(ctx, timer.ProcessInstanceKey, timer.ElementInstanceKey)
}

func (timer *Timer) GetProcessInstanceKey() bpmn_engine_store.IProcessInstanceKey {
	return timer.ProcessInstanceKey
}

func (timer *Timer) GetElementId() string {
	return timer.ElementId
}
