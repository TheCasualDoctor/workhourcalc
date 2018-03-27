package workhourcalc

import (
	"time"
	"errors"
)

type WorkHours struct {
	StartHour int
	StartMinute int
	EndHour int
	EndMinute int
}

type WorkDays []time.Weekday

func IsDuringWorkHours(day time.Time, workDays WorkDays, workHours WorkHours) bool {
	if !isWorkDay(day.Weekday(), workDays) {
		return false
	}

	if !isInsideWorkHours(day, workHours) {
		return false
	}

	return true
}

func GetWorkingHoursBetween(workHours WorkHours, workDays WorkDays, start time.Time, end time.Time) (float64, error) {
	if !start.Before(end) {
		return 0, errors.New("start date must be before end date")
	}

	if areSameDay(start, end) {
		return getHoursBetween(start, end), nil
	} else {
		return getWorkHoursBetween(workHours, workDays, start, end), nil
	}
}

//Private Functions
func isInsideWorkHours(day time.Time, workHours WorkHours) bool {
	dayStart := changeHourAndMinute(day, workHours.StartHour, workHours.StartMinute)

	if day.Before(dayStart) {
		return false
	}

	dayEnd := changeHourAndMinute(day, workHours.EndHour, workHours.EndMinute)

	if day.After(dayEnd) {
		return false
	}

	return true
}

func getWorkHoursOnDaysBetween(workHours WorkHours, workDays WorkDays, start time.Time, end time.Time) float64 {
	workDaysBetween := getWorkDaysBetween(workDays, start, end)
	hoursPerWorkday := getHoursPerWorkday(workHours)

	return float64(workDaysBetween) * hoursPerWorkday
}

func getWorkDaysBetween(workDays WorkDays, start time.Time, end time.Time) int {
	//check consecutive days
	if(areConsecutiveDays(start, end)) || areSameDay(start, end) {
		return 0
	}

	currentDay := start.AddDate(0, 0, 1) //Don't check first day
	workDayCount := 0

	for !areSameDay(currentDay, end) {
		if isWorkDay(currentDay.Weekday(), workDays) {
			workDayCount++
		}

		//Done, add 1 day
		currentDay = currentDay.AddDate(0, 0, 1)
	}

	return workDayCount
}

func isWorkDay(currentDay time.Weekday, workDays WorkDays) bool {
	for _, dayInt := range workDays {
		if currentDay == dayInt {
			return true
		}
	}

	return false
}

func getHoursPerWorkday(workHours WorkHours) float64 {
	start := time.Date(2018, 03, 18, workHours.StartHour, workHours.StartMinute, 0, 0, time.Local)
	end := time.Date(2018, 03, 18, workHours.EndHour, workHours.EndMinute, 0, 0, time.Local)
	return getHoursBetween(start, end)
}


func areSameDay(start time.Time, end time.Time) bool {
	return start.YearDay() == end.YearDay() && start.Year() == end.Year()
}

func areConsecutiveDays(start time.Time, end time.Time) bool {
	incrementedDay := start.AddDate(0, 0, 1)
	return incrementedDay.YearDay() == end.YearDay() && incrementedDay.Year() == end.Year()
}

func getHoursBetween(start time.Time, end time.Time) float64 {
	return end.Sub(start).Hours()
}

func getWorkHoursBetween(workHours WorkHours, workDays WorkDays, start time.Time, end time.Time) float64 {
	//startDay
	hoursToEndOfDay := getHoursUntilEndOfDay(workHours, workDays, start)
	//endDay
	hoursFromBeginningOfDay := getHoursFromStartOfDay(workHours, workDays, end)
	//days Between
	workHoursBetween := getWorkHoursOnDaysBetween(workHours, workDays, start, end)
	return hoursToEndOfDay + hoursFromBeginningOfDay + workHoursBetween
}

func getHoursUntilEndOfDay(workHours WorkHours, workDays WorkDays, day time.Time) float64 {
	end := changeHourAndMinute(day, workHours.EndHour, workHours.EndMinute)

	return getHoursBetween(day, end)
}

func getHoursFromStartOfDay(workHours WorkHours, workDays WorkDays, day time.Time) float64 {
	start := changeHourAndMinute(day, workHours.StartHour, workHours.StartMinute)

	return getHoursBetween(start, day)
}

func changeHourAndMinute(day time.Time, newHour int, newMinute int) time.Time {
	return time.Date(day.Year(), day.Month(), day.Day(), newHour, newMinute, 0, 0, day.Location())
}

func parseTime(dateTimeString string) time.Time {
	layout := "2006-01-02T15:04:05.000Z"
	t, _ := time.Parse(layout, dateTimeString)

	return t
}
