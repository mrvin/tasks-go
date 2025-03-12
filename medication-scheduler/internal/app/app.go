package app

import (
	"fmt"
	"time"

	"github.com/mrvin/tasks-go/medication-scheduler/internal/storage"
)

const (
	minutes = 60               // Количество минут в одном часе.
	m       = 15 * time.Minute // Рассчитанное время приема лекарств будет кратно этому значению.
)

func GenerateTimeTaking(from, to time.Duration, numPerDay int16) []int64 {
	result := make([]int64, numPerDay)
	if numPerDay == 1 {
		result[0] = rounUp(to-from/2, m, to)
		return result
	}

	period := (to - from) / time.Duration(numPerDay-1)

	taking := from
	for i := range result {
		result[i] = rounUp(taking, m, to)
		taking += period
	}

	return result
}

// rounUp округление до большего кратного m, но не больше to.
func rounUp(d, m, to time.Duration) int64 {
	r := d % m
	if r > time.Minute {
		d += m - r
	}
	if d > to {
		return int64(to)
	}

	return int64(d)
}

func ConvertTimeToStr(times []int64) []string {
	result := make([]string, len(times))

	for i := range times {
		result[i] = sprintfTime(time.Duration(times[i]))
	}

	return result
}

func SelectNextTakings(allTaking []storage.AllTaking, now time.Time, period time.Duration) []storage.Taking {
	result := make([]storage.Taking, 0)
	y, m, d := now.Date()
	begin := now.Sub(time.Date(y, m, d, 0, 0, 0, 0, now.Location()))
	end := begin + period
	for _, taking := range allTaking {
		for i := range taking.Times {
			t := time.Duration(taking.Times[i])
			if begin <= t && end >= t {
				timeStr := sprintfTime(t)
				result = append(result, storage.Taking{NameMedicine: taking.NameMedicine, Time: timeStr})
			}
		}
	}

	return result
}

func sprintfTime(t time.Duration) string {
	hour := int(t.Hours())
	minute := int(t.Minutes()) % minutes

	return fmt.Sprintf("%02d:%02d", hour, minute)
}
