package workhourcalc

import (
	"testing"
	"time"
)

func TestIsInsideWorkHours(t *testing.T) {
	workHours := WorkHours{
		StartHour: 7,
		StartMinute: 45,
		EndHour: 18,
		EndMinute: 25,
	}

	day := parseTime("2017-03-30T09:35:00.000Z")
	actual := isInsideWorkHours(day, workHours)
	if true != actual {
		t.Errorf("Incorrect, wanted: %d, got: %d.", true, actual)
	}

	day = parseTime("2017-03-30T01:35:00.000Z")
	actual = isInsideWorkHours(day, workHours)
	if false != actual {
		t.Errorf("Incorrect, wanted: %d, got: %d.", false, actual)
	}

	day = parseTime("2017-03-30T23:35:00.000Z")
	actual = isInsideWorkHours(day, workHours)
	if false != actual {
		t.Errorf("Incorrect, wanted: %d, got: %d.", false, actual)
	}
}

func TestIsDuringWorkHours (t *testing.T) {
	workHours := WorkHours{
		StartHour: 7,
		StartMinute: 45,
		EndHour: 18,
		EndMinute: 25,
	}

	workDays := []time.Weekday{time.Monday,time.Tuesday,time.Thursday,time.Friday}

	day := parseTime("2018-03-30T09:35:00.000Z")
	actual := IsDuringWorkHours(day, workDays, workHours)
	if true != actual {
		t.Errorf("Incorrect, wanted: %d, got: %d.", true, actual)
	}

	day = parseTime("2018-03-31T09:35:00.000Z")
	actual = IsDuringWorkHours(day, workDays, workHours)
	if false != actual {
		t.Errorf("Incorrect, wanted: %d, got: %d.", false, actual)
	}

	day = parseTime("2018-03-31T09:35:00.000Z")
	actual = IsDuringWorkHours(day, workDays, workHours)
	if false != actual {
		t.Errorf("Incorrect, wanted: %d, got: %d.", false, actual)
	}

	day = parseTime("2018-03-29T21:35:00.000Z")
	actual = IsDuringWorkHours(day, workDays, workHours)
	if false != actual {
		t.Errorf("Incorrect, wanted: %d, got: %d.", false, actual)
	}

	day = parseTime("2018-03-28T13:35:00.000Z")
	actual = IsDuringWorkHours(day, workDays, workHours)
	if false != actual {
		t.Errorf("Incorrect, wanted: %d, got: %d.", false, actual)
	}
}

func TestStartAfterEnd(t *testing.T) {
	start := parseTime("2017-03-30T09:35:00.000Z")
	end := parseTime("2017-03-29T13:05:00.000Z")

	workHours := WorkHours{
		StartHour: 7,
		StartMinute: 45,
		EndHour: 18,
		EndMinute: 25,
	}

	workDays := []time.Weekday{time.Monday,time.Tuesday,time.Wednesday,time.Thursday,time.Friday}

	_, err := GetWorkingHoursBetween(workHours, workDays, start, end)

	if err == nil {
		t.Errorf("Expected error, got none")
	}
}

func TestHoursBetweenSameDay(t *testing.T) {
	start := parseTime("2017-03-29T09:35:00.000Z")
	end := parseTime("2017-03-29T13:05:00.000Z")

	workHours := WorkHours{
		StartHour: 7,
		StartMinute: 45,
		EndHour: 18,
		EndMinute: 25,
	}

	workDays := []time.Weekday{time.Monday,time.Tuesday,time.Wednesday,time.Thursday,time.Friday}

	hours, _ := GetWorkingHoursBetween(workHours, workDays, start, end)

	expectedHours := 3.5

	if hours != expectedHours {
		t.Errorf("Incorrect, wanted: %d h, got: %d h.", expectedHours, hours)
	}
}

func TestHoursConsecutiveWorkdays(t *testing.T) {
	start := parseTime("2017-03-29T09:45:00.000Z") //7.25
	end := parseTime("2017-03-30T13:15:00.000Z") //5.25

	workHours := WorkHours{
		StartHour: 8,
		StartMinute: 00,
		EndHour: 17,
		EndMinute: 00,
	}

	workDays := []time.Weekday{time.Monday,time.Tuesday,time.Wednesday,time.Thursday,time.Friday}

	hours, _ := GetWorkingHoursBetween(workHours, workDays, start, end)

	expectedHours := 12.5

	if hours != expectedHours {
		t.Errorf("Incorrect, wanted: %d, got: %d.", expectedHours, hours)
	}
}

func TestHoursWithWorkdaysBetween(t *testing.T) {
	start := parseTime("2017-03-28T09:45:00.000Z") //7.25
	end := parseTime("2017-03-30T13:15:00.000Z") //5.25

	workHours := WorkHours{
		StartHour: 8,
		StartMinute: 00,
		EndHour: 17,
		EndMinute: 00,
	}

	workDays := []time.Weekday{time.Monday,time.Tuesday,time.Wednesday,time.Thursday,time.Friday}

	hours, _ := GetWorkingHoursBetween(workHours, workDays, start, end)

	expectedHours := 21.5

	if hours != expectedHours {
		t.Errorf("Incorrect, wanted: %d, got: %d.", expectedHours, hours)
	}
}

func TestHoursBetweenOverWeekend(t *testing.T) {
	start := parseTime("2018-03-29T09:40:00.000Z")
	end := parseTime("2018-04-04T13:10:00.000Z")

	workHours := WorkHours{
		StartHour: 8,
		StartMinute: 00,
		EndHour: 17,
		EndMinute: 00,
	}

	workDays := []time.Weekday{time.Monday,time.Tuesday,time.Wednesday,time.Thursday,time.Friday}

	hours, err := GetWorkingHoursBetween(workHours, workDays, start, end)

	expectedHours := 39.5

	if hours != expectedHours {
		t.Errorf("Incorrect, wanted: %d, got: %d.", expectedHours, hours)
	}

	if err != nil {
		t.Errorf("Was not expecting error, but got one.")
	}
}

func TestGetStartDayHours(t *testing.T) {
	day := parseTime("2017-03-29T09:45:00.000Z")

	workHours := WorkHours{
		StartHour: 7,
		StartMinute: 45,
		EndHour: 18,
		EndMinute: 15,
	}

	hours := getHoursUntilEndOfDay(workHours, day)

	expectedHours := 8.5

	if hours != expectedHours {
		t.Errorf("Incorrect, wanted: %d h, got: %d h.", expectedHours, hours)
	}
}

func TestGetEndDayHours(t *testing.T) {
	day := parseTime("2017-03-29T13:30:00.000Z")

	workHours := WorkHours{
		StartHour: 7,
		StartMinute: 45,
		EndHour: 18,
		EndMinute: 15,
	}

	hours := getHoursFromStartOfDay(workHours, day)

	expectedHours := 5.75

	if hours != expectedHours {
		t.Errorf("Incorrect, wanted: %d h, got: %d h.", expectedHours, hours)
	}
}

func TestChangeHourAndMinute(t *testing.T) {
	day := parseTime("2017-03-29T09:45:00.000Z")
	newHour := 15
	newMinute := 30

	expected := parseTime("2017-03-29T15:30:00.000Z")

	actual := changeHourAndMinute(day, newHour, newMinute)

	if expected != actual {
		t.Errorf("Incorrect, wanted: %d, got: %d.", expected, actual)
	}
}

func TestAreSameDay(t *testing.T) {
	day1 := parseTime("2017-03-29T09:45:00.000Z")
	day2 := parseTime("2017-03-29T10:13:00.000Z")

	actual := areSameDay(day1, day2)

	if true != actual {
		t.Errorf("Incorrect, wanted: %d, got: %d.", true, actual)
	}

	day1 = parseTime("2017-03-29T09:45:00.000Z")
	day2 = parseTime("2017-03-22T10:13:00.000Z")

	actual = areSameDay(day1, day2)

	if false != actual {
		t.Errorf("Incorrect, wanted: %d, got: %d.", false, actual)
	}
}

func TestAreConsecutiveDays(t *testing.T) {
	day1 := parseTime("2017-03-29T09:45:00.000Z")
	day2 := parseTime("2017-03-30T10:13:00.000Z")

	actual := areConsecutiveDays(day1, day2)

	if true != actual {
		t.Errorf("Incorrect, wanted: %d, got: %d.", true, actual)
	}

	day1 = parseTime("2017-03-29T09:45:00.000Z")
	day2 = parseTime("2016-03-30T10:13:00.000Z")

	actual = areConsecutiveDays(day1, day2)

	if false != actual {
		t.Errorf("Incorrect, wanted: %d, got: %d.", false, actual)
	}
}

func TestHoursPerWorkDay(t *testing.T) {
	workHours := WorkHours{
		StartHour: 7,
		StartMinute: 45,
		EndHour: 18,
		EndMinute: 15,
	}

	expected := 10.5

	actual := getHoursPerWorkday(workHours)

	if expected != actual {
		t.Errorf("Incorrect, wanted: %d, got: %d.", expected, actual)
	}
}

func TestGetWorkDaysBetween(t *testing.T) {
	workDays := []time.Weekday{time.Monday,time.Tuesday,time.Thursday,time.Friday}

	day1 := parseTime("2017-03-26T09:45:00.000Z")
	day2 := parseTime("2017-03-30T10:13:00.000Z")
	expected := 2
	actual := getWorkDaysBetween(workDays, day1, day2)
	if expected != actual {
		t.Errorf("Incorrect, wanted: %d, got: %d.", expected, actual)
	}

	day1 = parseTime("2017-03-26T09:45:00.000Z")
	day2 = parseTime("2017-03-27T10:13:00.000Z")
	expected = 0
	actual = getWorkDaysBetween(workDays, day1, day2)
	if expected != actual {
		t.Errorf("Incorrect, wanted: %d, got: %d.", expected, actual)
	}

	day1 = parseTime("2017-03-26T09:45:00.000Z")
	day2 = parseTime("2017-03-26T10:13:00.000Z")
	expected = 0
	actual = getWorkDaysBetween(workDays, day1, day2)
	if expected != actual {
		t.Errorf("Incorrect, wanted: %d, got: %d.", expected, actual)
	}
}

func TestIsWorkDay(t *testing.T) {
	workDays := []time.Weekday{time.Monday,time.Tuesday,time.Thursday,time.Friday}

	actual := isWorkDay(1, workDays)
	if true != actual {
		t.Errorf("Incorrect, wanted: %d, got: %d.", true, actual)
	}

	actual = isWorkDay(4, workDays)
	if true != actual {
		t.Errorf("Incorrect, wanted: %d, got: %d.", true, actual)
	}

	actual = isWorkDay(6, workDays)
	if false != actual {
		t.Errorf("Incorrect, wanted: %d, got: %d.", false, actual)
	}

	actual = isWorkDay(3, workDays)
	if false != actual {
		t.Errorf("Incorrect, wanted: %d, got: %d.", false, actual)
	}

	actual = isWorkDay(7, workDays)
	if false != actual {
		t.Errorf("Incorrect, wanted: %d, got: %d.", false, actual)
	}
}

func TestMoveToNextValidWorkTime(t *testing.T) {
	workHours := WorkHours{
		StartHour: 7,
		StartMinute: 45,
		EndHour: 18,
		EndMinute: 15,
	}
	workDays := []time.Weekday{time.Monday,time.Tuesday,time.Wednesday,time.Thursday,time.Friday}

	day := parseTime("2017-03-31T14:30:00.000Z")
	actual := moveToNextValidWorkTime(day, workDays, workHours)
	expected := day
	if expected != actual {
		t.Errorf("Incorrect, wanted: %d, got: %d.", expected, actual)
	}

	day = parseTime("2018-03-31T23:30:00.000Z")
	actual = moveToNextValidWorkTime(day, workDays, workHours)
	expected = parseTime("2018-04-02T07:45:00.000Z")
	if expected != actual {
		t.Errorf("Incorrect, wanted: %d, got: %d.", expected, actual)
	}
}


func TestMoveToLastValidWorkTime(t *testing.T) {
	workHours := WorkHours{
		StartHour: 7,
		StartMinute: 45,
		EndHour: 18,
		EndMinute: 15,
	}
	workDays := []time.Weekday{time.Monday,time.Tuesday,time.Wednesday,time.Thursday,time.Friday}

	day := parseTime("2017-03-31T14:30:00.000Z")
	actual := moveToLastValidWorkTime(day, workDays, workHours)
	expected := day
	if expected != actual {
		t.Errorf("Incorrect, wanted: %d, got: %d.", expected, actual)
	}

	day = parseTime("2018-04-01T23:30:00.000Z")
	actual = moveToLastValidWorkTime(day, workDays, workHours)
	expected = parseTime("2018-03-30T18:15:00.000Z")
	if expected != actual {
		t.Errorf("Incorrect, wanted: %d, got: %d.", expected, actual)
	}
}