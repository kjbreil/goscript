package helpers

import (
	"fmt"
	"github.com/adhocore/gronx"
	"sort"
	"strconv"
	"time"
)

func TimeToCron(t time.Time) string {
	minute := t.Minute()
	hour := t.Hour()
	cron := fmt.Sprintf("%d %d * * *", minute, hour)

	return cron
}

func CronToTime(cron string, t time.Time) (time.Time, error) {
	tickTime, err := gronx.NextTickAfter(cron, t, true)

	return tickTime, err
}
func CronToSegments(cron string) ([6]int, error) {
	segments, err := gronx.Segments(cron)
	var segs [6]int
	if err != nil {
		return segs, err
	}
	if len(segments) != 6 {
		return segs, fmt.Errorf("cron segments must be 6")
	}

	for i := range segs {
		if segments[i] == "*" {
			segs[i] = -1
			continue
		}
		var n int
		n, err = strconv.Atoi(segments[i])
		if err != nil {
			return segs, err
		}
		segs[i] = n
	}

	return segs, nil
}

func LastValidCron(crons []string, t time.Time) (string, error) {
	crons = SortCronJobs(crons)

	lastExp := ""
	lastHour := -1
	lastMin := -1
	for _, exp := range crons {
		cronTime, err := CronToTime(exp, t)
		if err != nil {
			return "", err
		}

		if cronTime.Hour() <= t.Hour() && cronTime.Hour() >= lastHour {
			if lastHour != cronTime.Hour() {
				lastExp = exp
				lastHour = cronTime.Hour()
				lastMin = cronTime.Minute()
			} else if cronTime.Minute() <= t.Minute() && t.Minute() >= lastMin {
				if t.Hour() != lastHour {
					lastMin = -1
				}
				lastExp = exp
				lastHour = cronTime.Hour()
				lastMin = cronTime.Minute()
			}
		}
	}

	if lastHour == -1 {
		lastExp = crons[len(crons)-1]
	}

	return lastExp, nil
}
func NextValidCron(crons []string, t time.Time) (string, error) {
	crons = SortCronJobs(crons)

	lastExp := ""
	lastHour := 25
	lastMin := 61
	for _, exp := range crons {
		cronTime, err := CronToTime(exp, t)
		if err != nil {
			return "", err
		}

		if cronTime.Hour() >= t.Hour() && cronTime.Hour() <= lastHour {
			if lastHour != cronTime.Hour() {
				lastExp = exp
				lastHour = cronTime.Hour()
				lastMin = cronTime.Minute()
			} else if cronTime.Minute() >= t.Minute() && t.Minute() <= lastMin {
				if t.Hour() != lastHour {
					lastMin = 61
				}
				lastExp = exp
				lastHour = cronTime.Hour()
				lastMin = cronTime.Minute()
			}
		}
	}

	if lastHour == 25 {
		lastExp = crons[0]
	}

	return lastExp, nil
}

type cronNextTick struct {
	expr string
	t    time.Time
}

func SortCronJobs(expressions []string) []string {
	zeroDay := time.Date(1, 1, 1, 0, 0, 0, 0, time.Local)
	nextTicks := make([]cronNextTick, len(expressions))

	for i := range expressions {
		tickTime, err := gronx.NextTickAfter(expressions[i], zeroDay, true)
		if err != nil {
			panic(err)
		}
		nextTicks[i] = cronNextTick{
			expr: expressions[i],
			t:    tickTime,
		}
	}

	sort.Slice(nextTicks, func(i, j int) bool {
		return nextTicks[i].t.Before(nextTicks[j].t)
	})
	for i := range nextTicks {
		expressions[i] = nextTicks[i].expr
	}

	return expressions
}

func NextTime(expressions []string, t time.Time) (time.Time, error) {
	cron, err := NextValidCron(expressions, t)
	if err != nil {
		return time.Time{}, err
	}

	nextTick, err := gronx.NextTickAfter(cron, t, false)
	if err != nil {
		return time.Time{}, err
	}
	// correct for error in gronx
	if t.Hour() != nextTick.Hour() {
		nextTick = nextTick.Truncate(60 * time.Minute)
	}

	return nextTick, nil
}
