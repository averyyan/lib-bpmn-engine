package store

import (
	"context"
	"errors"
	"fmt"
	"github.com/nitram509/lib-bpmn-engine/pkg/bpmn_engine"
	"github.com/nitram509/lib-bpmn-engine/pkg/bpmn_engine_store"
	"github.com/nitram509/lib-bpmn-engine/pkg/spec/BPMN20"
	"time"
)

func (store *EngineMemoryStore) CreateTimer(
	ctx context.Context,
	engineStore bpmn_engine_store.IBpmnEngineStore,
	processInstanceKey bpmn_engine_store.IProcessInstanceKey,
	timerKey bpmn_engine_store.ITimerKey,
	ice BPMN20.TIntermediateCatchEvent,
) (bpmn_engine_store.ITimer, error) {
	if store.timers[processInstanceKey] == nil {
		store.timers[processInstanceKey] = map[bpmn_engine_store.ITimerKey]*Timer{}
	}
	now := time.Now().Local()
	durationVal, err := findDurationValue(ice)
	if err != nil {
		return nil, err
	}
	timer := &Timer{
		engineStore:        engineStore,
		ElementId:          ice.Id,
		ElementInstanceKey: timerKey,
		ProcessInstanceKey: processInstanceKey,
		State:              bpmn_engine.TimerCreated,
		CreatedAt:          now,
		DueAt:              durationVal.Shift(now),
		Duration:           time.Duration(durationVal.TS) * time.Second,
	}
	store.timers[processInstanceKey][timerKey] = timer
	return store.timers[processInstanceKey][timerKey], nil
}

func (store *EngineMemoryStore) FindTimers(ctx context.Context, processInstanceKey bpmn_engine_store.IProcessInstanceKey) ([]bpmn_engine_store.ITimer, error) {
	timers, err := store.findTimers(processInstanceKey)
	if err != nil {
		return nil, err
	}
	var result []bpmn_engine_store.ITimer
	for _, timer := range timers {
		result = append(result, timer)
	}
	return result, nil
}
func (store *EngineMemoryStore) GetTimerState(ctx context.Context, processInstanceKey bpmn_engine_store.IProcessInstanceKey, timerKey bpmn_engine_store.ITimerKey) (bpmn_engine.TimerState, error) {
	timer, err := store.findTimer(processInstanceKey, timerKey)
	if err != nil {
		return "", err
	}
	return timer.State, nil
}

func (store *EngineMemoryStore) SetTimerState(ctx context.Context, processInstanceKey bpmn_engine_store.IProcessInstanceKey, timerKey bpmn_engine_store.ITimerKey, state bpmn_engine.TimerState) error {
	timer, err := store.findTimer(processInstanceKey, timerKey)
	if err != nil {
		return err
	}
	timer.State = state
	return nil
}

func (store *EngineMemoryStore) findTimers(processInstanceKey bpmn_engine_store.IProcessInstanceKey) (map[bpmn_engine_store.ITimerKey]*Timer, error) {
	timers, isExist := store.timers[processInstanceKey]
	if !isExist {
		return nil, errors.New(fmt.Sprintf("timers: ProcessInstanceKey=%d  is not exist", processInstanceKey))
	}
	return timers, nil
}

func (store *EngineMemoryStore) findTimer(processInstanceKey bpmn_engine_store.IProcessInstanceKey, timerKey bpmn_engine_store.ITimerKey) (*Timer, error) {
	timers, err := store.findTimers(processInstanceKey)
	if err != nil {
		return nil, err
	}
	timer, isExist := timers[timerKey]
	if !isExist {
		return nil, errors.New(fmt.Sprintf("timer: ProcessInstanceKey=%d && TimerKey=%d  is not exist", processInstanceKey, timerKey))
	}
	return timer, nil
}
