package datetime


// DayDiff returns the number of days between the two dates, taking into consideration only the day of the month.
// So 1 am on the 2nd and 11 pm on the 1st are a day apart.
// If dt1 is before dt2, the result will be negative
func DayDiff(dt1, dt2 DateTime) int {
	d1 := dt1.YearDay() + NumLeaps(dt1.Year() - 1) + dt1.Year() * 365
	d2 := dt2.YearDay() + NumLeaps(dt2.Year() - 1) + dt2.Year() * 365
	return d1 - d2
}
