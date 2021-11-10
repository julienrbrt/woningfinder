package scheduler

import (
	"fmt"
	"strconv"

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

	// check at 0h, 12h, 18h for every corporation (that are not first come first served)
	// a map permits to de-duplicate schedules
	schedules := map[int]cron.Schedule{
		0:    buildSchedule(parser, 0, 0),
		1230: buildSchedule(parser, 12, 30),
		180:  buildSchedule(parser, 18, 0),
	}

	// add specific selection time
	for _, t := range corporation.SelectionTime {
		result, err := strconv.Atoi(fmt.Sprintf("%d%d", t.Hour(), t.Minute()))
		if err != nil {
			// should never happen
			panic(err)
		}
		schedules[result] = buildSchedule(parser, t.Hour(), t.Minute())
	}

	list := make([]cron.Schedule, 0, len(schedules))
	for _, schedule := range schedules {
		list = append(list, schedule)
	}

	return list
}

func buildSchedule(parser cron.Parser, hour, minute int) cron.Schedule {
	schedule, err := parser.Parse(fmt.Sprintf("%d %d * * *", minute, hour))
	if err != nil {
		// should never happen
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
