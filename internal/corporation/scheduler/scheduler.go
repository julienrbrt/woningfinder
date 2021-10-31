package scheduler

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/woningfinder/woningfinder/internal/corporation"
)

var parser = cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)

// CorporationScheduler creates schedules (when to fetch their offer) given a selection time for a housing corporation
func CorporationScheduler(corporation corporation.Corporation) []cron.Schedule {
	schedules := []cron.Schedule{
		// always check at 00:15 for every corporation
		buildSchedule(parser, 0, 15),
	}

	// for corporation that has first come first served
	if hasFirstComeFirstServed(corporation) {
		// check by default every 30 minutes from 9-21 hours
		crontab := "*/30 9-21 * * *"

		// check every 2 minutes for a 20 minutes range if selection time defined
		if corporation.SelectionTime != (time.Time{}) {
			crontab = fmt.Sprintf("%d-20/2 %d * * *", corporation.SelectionTime.Minute(), corporation.SelectionTime.Hour())
		}

		schedule, err := parser.Parse(crontab)
		if err != nil {
			// should never happens
			panic(err)
		}

		schedules = append(schedules, schedule)
	} else {
		// always check at 18:15 for random
		schedules = append(schedules, buildSchedule(parser, 18, 15))
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
