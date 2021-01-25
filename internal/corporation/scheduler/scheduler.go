package scheduler

import (
	"fmt"
	"time"

	"github.com/woningfinder/woningfinder/internal/domain/entity"

	"github.com/robfig/cron/v3"
)

var (
	minutes = []int{0, 1, 2, 3}
	seconds = []int{0, 15, 30, 45}
)

// CorporationScheduler creates schedules (when to fetch their offer) given a selection time for a housing corporation
func CorporationScheduler(corporation entity.Corporation) []cron.Schedule {
	var schedules []cron.Schedule

	parser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	defaultSched := buildSchedule(parser, 17, 0, 0)

	// add schedule at selection time
	if corporation.SelectionTime == (time.Time{}) {
		schedules = append(schedules, defaultSched)
	} else {
		sched := buildSchedule(parser, corporation.SelectionTime.Hour(), corporation.SelectionTime.Minute(), corporation.SelectionTime.Second())
		schedules = append(schedules, sched)
	}

	// checks before the selection time and after the selection time
	// only checks multiple time if the selection time is defined
	if hasFirstComeFirstServed(corporation) && corporation.SelectionTime != (time.Time{}) {
		for _, minute := range minutes {
			for _, second := range seconds {
				// skip the first one as already set above
				if minute == 0 && second == 0 {
					continue
				}

				newTime := corporation.SelectionTime.Add(time.Duration(minute) * time.Minute).Add(time.Duration(second) * time.Second)
				sched := buildSchedule(parser, newTime.Hour(), newTime.Minute(), newTime.Second())
				schedules = append(schedules, sched)
			}
		}
	}

	// always check at midnight
	sched, _ := parser.Parse("@midnight")
	schedules = append(schedules, sched)

	return schedules
}

func buildSchedule(parser cron.Parser, hour, minute, second int) cron.Schedule {
	schedule, err := parser.Parse(fmt.Sprintf("%d %d %d * * *", second, minute, hour))
	if err != nil {
		panic(err)
	}

	return schedule
}

// hasFirstComeFirstServed returns true if a housing corporation select by first come first served
func hasFirstComeFirstServed(corporation entity.Corporation) bool {
	for _, s := range corporation.SelectionMethod {
		if s.Method == entity.SelectionFirstComeFirstServed {
			return true
		}
	}

	return false
}
