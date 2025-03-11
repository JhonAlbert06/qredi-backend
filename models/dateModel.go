package models

import "time"

type Date struct {
	Day      int    `json:"day" form:"day"`
	Month    int    `json:"month" form:"month"`
	Year     int    `json:"year" form:"year"`
	Hour     int    `json:"hour" form:"hour"`
	Minute   int    `json:"minute" form:"minute"`
	Second   int    `json:"second" form:"second"`
	Timezone string `json:"timezone" form:"timezone"`
}

func ToDate(date time.Time) Date {
	return Date{
		Day:      date.Day(),
		Month:    int(date.Month()),
		Year:     date.Year(),
		Hour:     date.Hour(),
		Minute:   date.Minute(),
		Second:   date.Second(),
		Timezone: date.Location().String(),
	}
}
