package scheduler

import (
	"fmt"

	"github.com/julienrbrt/woningfinder/internal/corporation"
	"github.com/robfig/cron/v3"
)

var parser = cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)

// CorporationScheduler creates schedules (when to fetch the offers) for a housing corporation
func CorporationScheduler(corporation corporation.Corporation) []cron.Schedule {
	var schedules []cron.Schedule

	// checks every 30 minutes
	schedule, err := parser.Parse("*/30 * * * *")
	if err != nil {
		// should never happens
		panic(err)
	}

	schedules = append(schedules, schedule)

	// as we check every 30 minutes, every hours, we only need to add a selection if the minutes aren't 0 or 30
	if corporation.SelectionTime.Minute() != 0 && corporation.SelectionTime.Minute() != 30 {
		schedules = append(schedules, buildSchedule(parser, corporation.SelectionTime.Hour(), corporation.SelectionTime.Minute()))
	}

	return schedules
}

func buildSchedule(parser cron.Parser, hour, minute int) cron.Schedule {
	schedule, err := parser.Parse(fmt.Sprintf("%d %d * * *", minute, hour))
	if err != nil {
		// should never happen
		panic(err)
	}

	return schedule
}
