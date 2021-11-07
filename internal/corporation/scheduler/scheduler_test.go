package scheduler_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/corporation/scheduler"
)

func Test_CorporationScheduler_Random(t *testing.T) {
	a := assert.New(t)
	corporation := corporation.Corporation{
		SelectionMethod: []corporation.SelectionMethod{
			corporation.SelectionRandom,
		},
	}

	now := time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC)
	schedules := scheduler.CorporationScheduler(corporation)
	a.Len(schedules, 2)
	a.Equal(0, schedules[0].Next(now).Hour())
	a.Equal(0, schedules[0].Next(now).Minute())
	a.Equal(18, schedules[1].Next(now).Hour())
	a.Equal(0, schedules[1].Next(now).Minute())
}

func Test_CorporationScheduler_FirstComeFirstServed(t *testing.T) {
	a := assert.New(t)

	corporation := corporation.Corporation{
		SelectionMethod: []corporation.SelectionMethod{
			corporation.SelectionFirstComeFirstServed,
			corporation.SelectionRegistrationDate,
		},
	}

	now := time.Date(0, 0, 0, 9, 45, 0, 0, time.UTC)
	schedules := scheduler.CorporationScheduler(corporation)
	a.Len(schedules, 1)
	a.Equal(10, schedules[0].Next(now).Hour())
	a.Equal(0, schedules[0].Next(now).Minute())
}

func Test_CorporationScheduler_SelectionTime(t *testing.T) {
	a := assert.New(t)

	corporation := corporation.Corporation{
		SelectionMethod: []corporation.SelectionMethod{
			corporation.SelectionRandom,
			corporation.SelectionRegistrationDate,
		},
		SelectionTime: []time.Time{
			scheduler.CreateSelectionTime(16, 00),
			scheduler.CreateSelectionTime(21, 00),
		},
	}

	now := time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC)
	schedules := scheduler.CorporationScheduler(corporation)
	a.Len(schedules, 4)
	a.Equal(0, schedules[0].Next(now).Hour())
	a.Equal(0, schedules[0].Next(now).Minute())
	a.Equal(18, schedules[1].Next(now).Hour())
	a.Equal(0, schedules[1].Next(now).Minute())
	a.Equal(16, schedules[2].Next(now).Hour())
	a.Equal(0, schedules[2].Next(now).Minute())
	a.Equal(21, schedules[3].Next(now).Hour())
	a.Equal(0, schedules[3].Next(now).Minute())
}
