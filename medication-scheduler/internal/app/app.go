package app

import (
	"fmt"
	"time"

	"github.com/mrvin/tasks-go/medication-scheduler/internal/storage"
)

const (
	fromTime = "08:00:00"       // Время начало приема лекарств.
	toTime   = "22:00:00"       // Время окончания приема лекарств.
	m        = 15 * time.Minute // Рассчитанное время приема лекарств будет кратно этому значению.
)

// GenerateTimes генерирует слайс времени приема лекарств на день от fromTime до toTime.
func GenerateTimes(numPerDay int16) []storage.TimeOnly {
	result := make([]storage.TimeOnly, numPerDay)
	from, _ := time.Parse(time.TimeOnly, fromTime)
	to, _ := time.Parse(time.TimeOnly, toTime)
	if numPerDay == 1 {
		period := to.Sub(from) / 2 //nolint:mnd
		result[0] = storage.TimeOnly(from.Add(period).Round(m))
		return result
	}

	period := to.Sub(from) / time.Duration(numPerDay-1)

	taking := from
	for i := range result {
		result[i] = storage.TimeOnly(taking.Round(m))
		taking = taking.Add(period)
	}

	return result
}

// SelectNextTakings выбираем из всех приемов лекарств только те, которые необходимо принять в ближайший период.
func SelectNextTakings(allTaking []storage.AllTaking, now time.Time, period time.Duration) []storage.Taking {
	result := make([]storage.Taking, 0)
	begin, _ := time.Parse(time.TimeOnly, fmt.Sprintf("%02d:%02d:%02d", now.Hour(), now.Minute(), now.Second()))
	end := begin.Add(period)
	for _, taking := range allTaking {
		for i := range taking.Times {
			t := time.Time(taking.Times[i])
			if t.After(begin) && t.Before(end) {
				result = append(result, storage.Taking{NameMedicine: taking.NameMedicine, Time: taking.Times[i]})
			}
		}
	}

	return result
}
