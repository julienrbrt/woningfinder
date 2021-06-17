package scheduler

import (
	"fmt"

	"github.com/robfig/cron/v3"
	"github.com/woningfinder/woningfinder/internal/corporation"
)

var (
	parser = cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
)

// CorporationScheduler creates schedules (when to fetch their offer) given a selection time for a housing corporation
func CorporationScheduler(corporation corporation.Corporation) []cron.Schedule {
	var schedules []cron.Schedule

	// add schedule at selection time (always check at 18:15 and 00:15)
	schedules = append(schedules, buildSchedule(parser, 18, 15))
	schedules = append(schedules, buildSchedule(parser, 0, 15))

	// for corporation that has first come first served, check every 30 minutes
	if hasFirstComeFirstServed(corporation) {
		schedule, err := parser.Parse("*/30 9-21 * * *")
		if err != nil {
			// should never happens
			panic(err)
		}

		schedules = append(schedules, schedule)
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
