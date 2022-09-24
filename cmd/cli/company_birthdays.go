package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var companyBirthdaysCmd = &cobra.Command{
	Use:   "company-birthdays",
	Short: "Company Birthdays of the day",
	RunE: func(cmd *cobra.Command, args []string) error {
		now := time.Now().Truncate(time.Hour * 24)
		start := now
		end := now.AddDate(0, 0, 1)

		if now.Weekday() == time.Monday {
			start = now.AddDate(0, 0, -2)
		}

		if now.Month() == time.February && now.Day() == 28 && now.AddDate(0, 0, 1).Month() == time.March {
			end = now.AddDate(0, 0, 1)
		}

		items, err := luccaClient.GetCompanyBirthdays(context.Background(), start, end)
		if err != nil {
			return fmt.Errorf("get company birthdays: %w", err)
		}

		for _, birthday := range items {
			fmt.Printf("%s %s - %s in the company\n", birthday.FirstName, birthday.LastName, humanDuration(birthday.ContractStart, now))
		}

		return nil
	},
}
