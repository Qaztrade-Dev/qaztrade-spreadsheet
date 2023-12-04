package domain

import (
	"time"
)

var (
	nonWorkingsDays = []string{
		// Декабрь 2023
		"02-12-2023", "03-12-2023",
		"09-12-2023", "10-12-2023",
		"16-12-2023", "17-12-2023",
		"18-12-2023", "23-12-2023", "24-12-2023",
		"30-12-2023", "31-12-2023",

		// Январь 2024
		"01-01-2024", "02-01-2024", "06-01-2024", "07-01-2024",
		"13-01-2024", "14-01-2024",
		"20-01-2024", "21-01-2024",
		"27-01-2024", "28-01-2024",

		// Февраль 2024
		"03-02-2024", "04-02-2024",
		"10-02-2024", "11-02-2024",
		"17-02-2024", "18-02-2024",
		"24-02-2024", "25-02-2024",

		// Март 2024
		"02-03-2024", "03-03-2024",
		"08-03-2024", "09-03-2024", "10-03-2024",
		"16-03-2024", "17-03-2024",
		"21-03-2024", "22-03-2024", "23-03-2024", "24-03-2024",
		"25-03-2024", "30-03-2024", "31-03-2024",

		// Апрель 2024
		"06-04-2024", "07-04-2024",
		"13-04-2024", "14-04-2024",
		"20-04-2024", "21-04-2024",
		"27-04-2024", "28-04-2024",

		// Май 2024
		"01-05-2024", "04-05-2024", "05-05-2024",
		"07-05-2024", "09-05-2024", "11-05-2024", "12-05-2024",
		"18-05-2024", "19-05-2024",
		"25-05-2024", "26-05-2024",

		// Июнь 2024
		"01-06-2024", "02-06-2024",
		"08-06-2024", "09-06-2024",
		"15-06-2024", "16-06-2024",
		"22-06-2024", "23-06-2024",
		"29-06-2024", "30-06-2024",

		// Июль 2024
		"06-07-2024", "07-07-2024",
		"08-07-2024", "13-07-2024", "14-07-2024",
		"20-07-2024", "21-07-2024",
		"27-07-2024", "28-07-2024",

		// Август 2024
		"03-08-2024", "04-08-2024",
		"10-08-2024", "11-08-2024",
		"17-08-2024", "18-08-2024",
		"24-08-2024", "25-08-2024",
		"30-08-2024", "31-08-2024",

		// Сентябрь 2024
		"01-09-2024",
		"07-09-2024", "08-09-2024",
		"14-09-2024", "15-09-2024",
		"21-09-2024", "22-09-2024",
		"28-09-2024", "29-09-2024",

		// Октябрь 2024
		"05-10-2024", "06-10-2024",
		"12-10-2024", "13-10-2024",
		"19-10-2024", "20-10-2024",
		"25-10-2024", "26-10-2024", "27-10-2024",

		// Ноябрь 2024
		"02-11-2024", "03-11-2024",
		"09-11-2024", "10-11-2024",
		"16-11-2024", "17-11-2024",
		"23-11-2024", "24-11-2024",
		"30-11-2024",

		// Декабрь 2024
		"01-12-2024",
		"07-12-2024", "08-12-2024",
		"14-12-2024", "15-12-2024",
		"16-12-2024", "21-12-2024", "22-12-2024",
		"28-12-2024", "29-12-2024",
	}

	nonWorkingsDaysMap map[string]bool
)

func init() {
	nonWorkingsDaysMap = make(map[string]bool)

	for _, day := range nonWorkingsDays {
		nonWorkingsDaysMap[day] = true
	}
}

func isNonWorkingDay(date string) bool {
	return nonWorkingsDaysMap[date]
}

func isNonWorkingDayTime(input time.Time) bool {
	return isNonWorkingDay(input.Format("02-01-2006"))
}

func getNextDayDate(input time.Time) time.Time {
	return input.AddDate(0, 0, 1)
}

func toDate(input time.Time) time.Time {
	return time.Date(input.Year(), input.Month(), input.Day(), 0, 0, 0, 0, input.Location())
}

func CalculateCountdown(input time.Time) *time.Duration {
	var (
		now         = input.UTC().Add(6 * time.Hour) // to be in Astana time
		currentDate = toDate(now)
	)

	var (
		allocatedDaysCount = 7
		nextDay            = getNextDayDate(currentDate)
	)

	for allocatedDaysCount > 0 {
		if !isNonWorkingDayTime(nextDay) {
			allocatedDaysCount--
		}
		nextDay = getNextDayDate(nextDay)
	}

	result := nextDay.Sub(input)
	return &result
}