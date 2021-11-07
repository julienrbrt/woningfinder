package scheduler

import (
	"fmt"

	"github.com/robfig/cron/v3"
	"github.com/woningfinder/woningfinder/internal/corporation"
)

var parser = cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)

// CorporationScheduler creates schedules (when to fetch the offers) for a housing corporation
func CorporationScheduler(corporation corporation.Corporation) []cron.Schedule {
	// for corporation that has first come first served we only need to check often
	if hasFirstComeFirstServed(corporation) {
		// checks every 30 minutes
		schedule, err := parser.Parse("*/30 * * * *")
		if err != nil {
			// should never happens
			panic(err)
		}

		return []cron.Schedule{schedule}
	}

	schedules := []cron.Schedule{
		buildSchedule(parser, 0, 0),  // check at 00:00 for every corporation (that are not first come first served)
		buildSchedule(parser, 18, 0), // check at 18:00 for every corporation (that are not first come first served)
	}

	// add specific selection time
	for _, t := range corporation.SelectionTime {
		schedules = append(schedules, buildSchedule(parser, t.Hour(), t.Minute()))
	}

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
