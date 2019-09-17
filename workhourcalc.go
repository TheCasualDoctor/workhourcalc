package workhourcalc

import (
	"errors"
	"time"
)

type WorkHours struct {
	StartHour int
	StartMinute int
	EndHour int
	EndMinute int
}

type WorkDays []time.Weekday

func SubtractWorkHours(day time.Time, hoursToSubtract float64, workDays WorkDays, workHours WorkHours) time.Time {
	//Just subtract the time and see if it's valid.
	start := day.Add(-time.Hour * time.Duration(hoursToSubtract))

	timeBetween, _ := GetWorkingHoursBetween(workHours, workDays, start, day)

	if IsDuringWorkHours(start, workDays, workHours) && timeBetween == hoursToSubtract {
		return start
	} else {
		remainingHours := hoursToSubtract - timeBetween
		start = forceMoveToLastValidWorkTime(start, workDays, workHours).Add(-time.Hour * time.Duration(remainingHours))
		return start
	}
}

func AddWorkHours(day time.Time, hoursToAdd float64, workDays WorkDays, workHours WorkHours) time.Time {
	//Just add the time and see if it's valid.
	end := day.Add(time.Hour * time.Duration(hoursToAdd))

	timeBetween, _ := GetWorkingHoursBetween(workHours, workDays, day, end)

	if IsDuringWorkHours(end, workDays, workHours) && timeBetween == hoursToAdd {
		return end
	} else {
		remainingHours := hoursToAdd - timeBetween
		end = forceMoveToNextValidWorkTime(end, workDays, workHours).Add(time.Hour * time.Duration(remainingHours))
		return end
	}
}

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

	start = moveToNextValidWorkTime(start, workDays, workHours)
	end = moveToLastValidWorkTime(end, workDays, workHours)
	if start.After(end) {
		start, end = end, start
	}

	if areSameDay(start, end) {
		return getHoursBetween(start, end), nil
	} else {
		return getWorkHoursBetween(workHours, workDays, start, end), nil
	}
}

func GetNextValidWorkTime(dateTime time.Time, workDays WorkDays, workHours WorkHours) time.Time {
	return moveToNextValidWorkTime(dateTime, workDays, workHours)
}

//Private Functions
func moveToNextValidWorkTime(dateTime time.Time, workDays WorkDays, workHours WorkHours) time.Time {
	if IsDuringWorkHours(dateTime, workDays, workHours) {
		return dateTime
	}

	for !isWorkDay(dateTime.Weekday(), workDays) {
		dateTime = dateTime.AddDate(0, 0, 1)
	}

	dateTime = changeHourAndMinute(dateTime, workHours.StartHour, workHours.StartMinute)

	return dateTime
}

func moveToLastValidWorkTime(dateTime time.Time, workDays WorkDays, workHours WorkHours) time.Time {
	if IsDuringWorkHours(dateTime, workDays, workHours) {
		return dateTime
	}

	for !isWorkDay(dateTime.Weekday(), workDays) {
		dateTime = dateTime.AddDate(0, 0, -1)
	}

	dateTime = changeHourAndMinute(dateTime, workHours.EndHour, workHours.EndMinute)

	return dateTime
}

func forceMoveToLastValidWorkTime(dateTime time.Time, workDays WorkDays, workHours WorkHours) time.Time {
	dateTime = dateTime.AddDate(0, 0, -1)

	for !isWorkDay(dateTime.Weekday(), workDays) {
		dateTime = dateTime.AddDate(0, 0, -1)
	}

	dateTime = changeHourAndMinute(dateTime, workHours.EndHour, workHours.EndMinute)

	return dateTime
}

func forceMoveToNextValidWorkTime(dateTime time.Time, workDays WorkDays, workHours WorkHours) time.Time {
	dateTime = dateTime.AddDate(0, 0, 1)

	for !isWorkDay(dateTime.Weekday(), workDays) {
		dateTime = dateTime.AddDate(0, 0, 1)
	}

	dateTime = changeHourAndMinute(dateTime, workHours.StartHour, workHours.StartMinute)

	return dateTime
}

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
	hoursToEndOfDay := getHoursUntilEndOfDay(workHours, start)
	//endDay
	hoursFromBeginningOfDay := getHoursFromStartOfDay(workHours, end)
	//days Between
	workHoursBetween := getWorkHoursOnDaysBetween(workHours, workDays, start, end)
	return hoursToEndOfDay + hoursFromBeginningOfDay + workHoursBetween
}

func getHoursUntilEndOfDay(workHours WorkHours, day time.Time) float64 {
	end := changeHourAndMinute(day, workHours.EndHour, workHours.EndMinute)

	return getHoursBetween(day, end)
}

func getHoursFromStartOfDay(workHours WorkHours, day time.Time) float64 {
	start := changeHourAndMinute(day, workHours.StartHour, workHours.StartMinute)

	return getHoursBetween(start, day)
}

func changeHourAndMinute(day time.Time, newHour int, newMinute int) time.Time {
	return time.Date(day.Year(), day.Month(), day.Day(), newHour, newMinute, 0, 0, time.Local)
}

func parseTime(dateTimeString string) time.Time {
	layout := "2006-01-02T15:04:05.000Z"
	t, _ := time.ParseInLocation(layout, dateTimeString, time.Local)

	return t
}
