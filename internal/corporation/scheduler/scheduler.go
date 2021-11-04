package scheduler

import (
	"fmt"

	"github.com/robfig/cron/v3"
	"github.com/woningfinder/woningfinder/internal/corporation"
)

var parser = cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)

// CorporationScheduler creates schedules (when to fetch the offers) for a housing corporation
func CorporationScheduler(corporation corporation.Corporation) []cron.Schedule {
	schedules := []cron.Schedule{
		buildSchedule(parser, 0, 00), // check at 00:00 for every corporation
	}

	// for corporation that has first come first served we only need to check often
	if hasFirstComeFirstServed(corporation) {
		// check by default every 30 minutes from 9-21 hours
		schedule, err := parser.Parse("*/30 9-21 * * *")
		if err != nil {
			// should never happens
			panic(err)
		}

		return append(schedules, schedule)
	}

	// add specific selection time
	for _, t := range corporation.SelectionTime {
		schedules = append(schedules, buildSchedule(parser, t.Hour(), t.Minute()))
	}
	schedules = append(schedules, buildSchedule(parser, 18, 00)) // check at 18:00 for every corporation (that arent's first come first served)

	return schedules
}

func buildSchedule(parser cron.Parser, hour, minute int) cron.Schedule {
	schedule, err := parser.Parse(fmt.Sprintf("%d %d * * *", minute, hour))
	if err != nil {
		// should never happens
		panic(err)
	}

	return schedule
}

// hasFirstComeFirstServed returns true if a housing corporation select by first come first served
func hasFirstComeFirstServed(corp corporation.Corporation) bool {
	for _, s := range corp.SelectionMethod {
		if s == corporation.SelectionFirstComeFirstServed {
			return true
		}
	}

	return false
}
