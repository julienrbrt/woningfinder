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

	now := time.Now()
	schedules := scheduler.CorporationScheduler(corporation)
	a.Len(schedules, 2)
	a.Equal(schedules[0].Next(now).Hour(), 18)
	a.Equal(schedules[0].Next(now).Minute(), 15)
	a.Equal(schedules[1].Next(now).Hour(), 0)
	a.Equal(schedules[1].Next(now).Minute(), 15)
}

func Test_CorporationScheduler_FirstComeFirstServed(t *testing.T) {
	a := assert.New(t)

	corporation := corporation.Corporation{
		SelectionMethod: []corporation.SelectionMethod{
			corporation.SelectionFirstComeFirstServed,
		},
	}

	now := time.Now()
	schedules := scheduler.CorporationScheduler(corporation)
	a.Len(schedules, 3)
	a.Equal(schedules[0].Next(now).Hour(), 18)
	a.Equal(schedules[0].Next(now).Minute(), 15)
	a.Equal(schedules[1].Next(now).Hour(), 0)
	a.Equal(schedules[1].Next(now).Minute(), 15)
}
