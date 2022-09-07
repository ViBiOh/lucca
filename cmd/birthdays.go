package cmd

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/ViBiOh/httputils/v4/pkg/httpjson"
	"github.com/spf13/cobra"
)

type Birthday struct {
	BirthDate        string `json:"birthDate"`
	BirthdayThisYear time.Time
	Date             time.Time
	Name             string `json:"name"`
	FirstName        string `json:"firstName"`
	LastName         string `json:"lastName"`
	ID               int    `json:"id"`
}

type Birthdays struct {
	Data struct {
		Items []Birthday `json:"items"`
	} `json:"data"`
}

// BirthdaysByDate sort Birthday by Date
type BirthdaysByDate []Birthday

func (a BirthdaysByDate) Len() int      { return len(a) }
func (a BirthdaysByDate) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a BirthdaysByDate) Less(i, j int) bool {
	return a[i].BirthdayThisYear.Before(a[j].BirthdayThisYear)
}

var birthdaysCmd = &cobra.Command{
	Use:   "birthdays",
	Short: "Birthdays of the day",
	RunE: func(cmd *cobra.Command, args []string) error {
		now := time.Now()
		start := now.Format(dateISOFormat)
		end := now.Format(dateISOFormat)

		if now.Weekday() == time.Monday {
			start = now.AddDate(0, 0, -2).Format(dateISOFormat)
		}

		response, err := req.Path("/api/v3/users/birthday?fields=id,name,firstname,lastname,birthDate&startsOn=%s&endsOn=%s", start, end).Send(context.Background(), nil)
		if err != nil {
			return fmt.Errorf("get birthdays: %w", err)
		}

		var birthdays Birthdays
		if err = httpjson.Read(response, &birthdays); err != nil {
			return fmt.Errorf("read birthdays: %w", err)
		}

		items := birthdays.Data.Items
		for index := range items {
			if date, err := time.Parse(dateTimeFormat, items[index].BirthDate); err == nil {
				items[index].Date = date
				items[index].BirthdayThisYear = time.Date(now.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
			}
		}

		sort.Sort(BirthdaysByDate(items))

		var currentDay string

		for _, birthday := range items {
			if start != end {
				if birthDayOfWeek := birthday.BirthdayThisYear.Format("Monday"); currentDay != birthDayOfWeek {
					if len(currentDay) > 0 {
						fmt.Printf("\n")
					}

					currentDay = birthDayOfWeek
					fmt.Println(currentDay)
				}
			}

			fmt.Printf("%s %s", birthday.FirstName, birthday.LastName)

			if !birthday.Date.IsZero() {
				fmt.Printf(" - %d years", now.Year()-birthday.Date.Year())
			}

			fmt.Printf("\n")
		}

		return nil
	},
}
