# workhourcalc
Calculates the number of work hours between two datetimes, given a set of working days and start and end times of business.

Golang lib to help calculate the number of work hours between two given datetimes.

The primary function is GetWorkingHoursBetween(workHours, workDays, start, end). With the start and end datetimes, as well as objects signifying work days and start and end times for those work days, a float64 will be returned denoting hours between the two times.

Another useful function is IsDuringWorkHours(datetime, workDays, workHours). This function simply returns true if the given time is during work hours on a workday, or false if not. Primarily used to prevent querying for times where the start and/or enddate might fall outside of work hours.

Currently not working or tested for start or end datetimes outside of work hours- coming soon.
