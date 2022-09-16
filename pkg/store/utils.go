package store

import (
	"errors"
	"fmt"
	"github.com/nitram509/lib-bpmn-engine/pkg/spec/BPMN20"
	"github.com/senseyeio/duration"
)

func findDurationValue(ice BPMN20.TIntermediateCatchEvent) (duration.Duration, error) {
	durationStr := &ice.TimerEventDefinition.TimeDuration.XMLText
	if durationStr == nil {
		return duration.Duration{}, errors.New(fmt.Sprintf("Can't find 'timeDuration' value for INTERMEDIATE_CATCH_EVENT with id=%s", ice.Id))
	}
	d, err := duration.ParseISO8601(*durationStr)
	return d, err
}
