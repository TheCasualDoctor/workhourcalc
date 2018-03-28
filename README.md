# workhourcalc
Calculates the number of work hours between two datetimes, given a set of working days and start and end times of business.

Golang lib to help calculate the number of work hours between two given datetimes.

Helpful functions:
GetWorkingHoursBetween(workHours, workDays, start, end).
With the start and end datetimes, as well as objects signifying work days and start and end times for those work days, a float64 will be returned denoting hours between the two times.
If given start or end dates outside of work datetimes, start dates will be moved forward to the next valid start datetime, and end dates will be moved back to the last valid datetime.

IsDuringWorkHours(datetime, workDays, workHours).
This function simply returns true if the given time is during work hours on a workday, or false if not. Primarily used to prevent querying for times where the start and/or enddate might fall outside of work hours.

SubtractWorkHours(start, hoursToSubtract, workDays, workHours) & AddWorkHours(start, hoursToSubtract, workDays, workHours)
Similar functions that add or remove a float64 of hours to the given start date, taking into account work hours. ie. Adding to hours to 4pm on a Friday will give you 9am on a Monday (assuming Mon-Fri 8:00 to 17:00 work times).

