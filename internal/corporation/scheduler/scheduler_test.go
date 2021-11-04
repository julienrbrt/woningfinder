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
	a.Equal(schedules[0].Next(now).Hour(), 0)
	a.Equal(schedules[0].Next(now).Minute(), 15)
	a.Equal(schedules[1].Next(now).Hour(), 18)
	a.Equal(schedules[1].Next(now).Minute(), 15)
}

func Test_CorporationScheduler_FirstComeFirstServed(t *testing.T) {
	a := assert.New(t)

	corporation := corporation.Corporation{
		SelectionMethod: []corporation.SelectionMethod{
			corporation.SelectionFirstComeFirstServed,
		},
	}

	now := time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC)
	schedules := scheduler.CorporationScheduler(corporation)
	a.Len(schedules, 2)
	a.Equal(schedules[0].Next(now).Hour(), 0)
	a.Equal(schedules[0].Next(now).Minute(), 15)
	a.Equal(schedules[1].Next(now).Hour(), 9)
	a.Equal(schedules[1].Next(now).Minute(), 0)
}

func Test_CorporationScheduler_FirstComeFirstServed_SelectionTime(t *testing.T) {
	a := assert.New(t)

	corporation := corporation.Corporation{
		SelectionMethod: []corporation.SelectionMethod{
			corporation.SelectionFirstComeFirstServed,
		},
		SelectionTime: []time.Time{
			scheduler.CreateSelectionTime(18, 00),
		},
	}

	now := time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC)
	schedules := scheduler.CorporationScheduler(corporation)
	a.Len(schedules, 2)
	a.Equal(schedules[0].Next(now).Hour(), 0)
	a.Equal(schedules[0].Next(now).Minute(), 15)
	a.Equal(schedules[1].Next(now).Hour(), 18)
	a.Equal(schedules[1].Next(now).Minute(), 00)
}
