package scheduler_test

import (
	"testing"
	"time"

	"github.com/woningfinder/woningfinder/internal/corporation/scheduler"

	"github.com/stretchr/testify/assert"
	"github.com/woningfinder/woningfinder/internal/domain/entity"
)

func Test_CorporationScheduler_NoSelectionTime(t *testing.T) {
	a := assert.New(t)
	corporation := entity.Corporation{
		SelectionMethod: []entity.SelectionMethod{
			{
				Method: entity.SelectionFirstComeFirstServed,
			},
		},
	}

	now := time.Now()
	schedules := scheduler.CorporationScheduler(corporation)
	a.Len(schedules, 2)
	a.Equal(schedules[0].Next(now).Hour(), 17)
	a.Equal(schedules[0].Next(now).Minute(), 0)
	a.Equal(schedules[1].Next(now).Hour(), 0)
	a.Equal(schedules[1].Next(now).Minute(), 0)
}

func Test_CorporationScheduler_WithSelectionTime(t *testing.T) {
	a := assert.New(t)
	corporation := entity.Corporation{
		SelectionTime: scheduler.CreateSelectionTime(12, 55, 0),
		SelectionMethod: []entity.SelectionMethod{
			{
				Method: entity.SelectionRandom,
			},
		},
	}

	now := time.Now()
	schedules := scheduler.CorporationScheduler(corporation)
	a.Len(schedules, 2)
	a.Equal(schedules[0].Next(now).Hour(), 12)
	a.Equal(schedules[0].Next(now).Minute(), 55)
	a.Equal(schedules[1].Next(now).Hour(), 0)
	a.Equal(schedules[1].Next(now).Minute(), 0)
}

func Test_CorporationScheduler_FirstComeFirstServed(t *testing.T) {
	a := assert.New(t)

	corporation := entity.Corporation{
		SelectionTime: scheduler.CreateSelectionTime(17, 59, 15),
		SelectionMethod: []entity.SelectionMethod{
			{
				Method: entity.SelectionFirstComeFirstServed,
			},
		},
	}

	now := time.Now()
	schedules := scheduler.CorporationScheduler(corporation)
	a.Len(schedules, 11)
	a.Equal(schedules[0].Next(now).Hour(), 17)
	a.Equal(schedules[0].Next(now).Minute(), 59)
	a.Equal(schedules[0].Next(now).Second(), 15)
	a.Equal(schedules[1].Next(now).Hour(), 17)
	a.Equal(schedules[1].Next(now).Minute(), 59)
	a.Equal(schedules[1].Next(now).Second(), 30)
	a.Equal(schedules[3].Next(now).Hour(), 18)
	a.Equal(schedules[3].Next(now).Minute(), 0)
	a.Equal(schedules[3].Next(now).Second(), 0)
	a.Equal(schedules[len(schedules)-1].Next(now).Hour(), 0)
	a.Equal(schedules[len(schedules)-1].Next(now).Minute(), 0)
}
