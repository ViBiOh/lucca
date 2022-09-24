package cli

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/ViBiOh/lucca/pkg/lucca"
	"github.com/spf13/cobra"
)

type Birthday struct {
	lucca.User
	BirthdayThisYear time.Time
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
		start := now
		end := now

		if now.Weekday() == time.Monday {
			start = now.AddDate(0, 0, -2)
		}

		if now.Month() == time.February && now.Day() == 28 && now.AddDate(0, 0, 1).Month() == time.March {
			end = now.AddDate(0, 0, 1)
		}

		birthdays, err := luccaClient.GetBirthdays(context.Background(), start, end)
		if err != nil {
			return fmt.Errorf("get birthdays: %w", err)
		}

		items := make([]Birthday, len(birthdays))
		for index, birthday := range birthdays {
			var birthdayThisYear time.Time

			if !birthday.Date.IsZero() {
				birthdayThisYear = time.Date(now.Year(), birthday.Date.Month(), birthday.Date.Day(), 0, 0, 0, 0, time.UTC)
			}

			items[index] = Birthday{
				User:             birthday,
				BirthdayThisYear: birthdayThisYear,
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

			if !birthday.ContractStart.IsZero() {
				fmt.Printf(" - %s in the company", humanDuration(birthday.ContractStart, now))
			}

			fmt.Printf("\n")
		}

		return nil
	},
}

func humanDuration(start, now time.Time) string {
	if start.After(now) {
		return "not yet"
	}

	var output strings.Builder

	years := now.Year() - start.Year()

	nowMonth := now.Month()
	if nowMonth < start.Month() {
		nowMonth += 12
		years -= 1
	}

	if years > 0 {
		output.WriteString(fmt.Sprintf("%d year", years))
		if years > 1 {
			output.WriteString("s")
		}
	}

	if months := nowMonth - start.Month(); months > 0 {
		if output.Len() != 0 {
			output.WriteString(", ")
		}

		output.WriteString(fmt.Sprintf("%d month", months))
		if months > 1 {
			output.WriteString("s")
		}
	}

	if output.Len() == 0 {
		return "just arrived"
	}

	return output.String()
}
