package bpmn_engine_store

import (
	"context"
	"github.com/nitram509/lib-bpmn-engine/pkg/bpmn_engine"
	"github.com/nitram509/lib-bpmn-engine/pkg/spec/BPMN20"
	"time"
)

func (state *BpmnEngineState) handleIntermediateTimerCatchEvent(ctx context.Context, instance IProcessInstanceInfo, ice BPMN20.TIntermediateCatchEvent) bool {
	timer := findExistingTimerNotYetTriggered(ctx, instance, ice.Id)
	if timer == nil {
		newTimer, err := instance.CreateTimer(ctx, ice)
		if err != nil {
			return false
		}
		if err != nil {
			// TODO: proper error handling
			return false
		}
		timer = newTimer
	}
	if time.Now().After(timer.GetDueAt()) {
		if err := timer.SetState(ctx, bpmn_engine.TimerTriggered); err != nil {
			return false
		}
		return true
	}
	return false
}

func findExistingTimerNotYetTriggered(ctx context.Context, instance IProcessInstanceInfo, id string) ITimer {
	var t ITimer
	timers, err := instance.FindTimers(ctx)
	if err != nil {
		return nil
	}
	for _, timer := range timers {
		timerState, err := timer.GetState(ctx)
		if err != nil {
			return nil
		}
		if timer.GetElementId() == id && timerState == bpmn_engine.TimerCreated {
			t = timer
			break
		}
	}
	return t
}

func checkDueTimersAndFindIntermediateCatchEvent(ctx context.Context, instance IProcessInstanceInfo) *BPMN20.BaseElement {
	intermediateCatchEvents := instance.GetProcessInfo().GetDefinitions().Process.IntermediateCatchEvent
	timers, err := instance.FindTimers(ctx)
	if err != nil {
		timers = []ITimer{}
	}
	for _, timer := range timers {
		timerState, err := timer.GetState(ctx)
		if err != nil {
			return nil
		}
		if timer.GetProcessInstanceKey() == instance.GetInstanceKey() && timerState == bpmn_engine.TimerCreated {
			if time.Now().After(timer.GetDueAt()) {
				for _, ice := range intermediateCatchEvents {
					if ice.Id == timer.GetElementId() {
						be := BPMN20.BaseElement(ice)
						return &be
					}
				}
			}
		}
	}
	return nil
}
