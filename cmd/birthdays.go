package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/ViBiOh/httputils/v4/pkg/httpjson"
	"github.com/spf13/cobra"
)

type Birthdays struct {
	Data struct {
		Items []struct {
			BirthDate string `json:"birthDate"`
			ID        int    `json:"id"`
			Name      string `json:"name"`
			FirstName string `json:"firstName"`
			LastName  string `json:"lastName"`
		} `json:"items"`
	} `json:"data"`
}

var birthdaysCmd = &cobra.Command{
	Use:   "birthdays",
	Short: "Birthdays of the day",
	RunE: func(cmd *cobra.Command, args []string) error {
		now := time.Now()
		today := now.Format(dateISOFormat)

		response, err := req.Path("/api/v3/users/birthday?fields=id,name,firstname,lastname,birthDate&startsOn=%s&endsOn=%s", today, today).Send(context.Background(), nil)
		if err != nil {
			return fmt.Errorf("get birthdays: %w", err)
		}

		var birthdays Birthdays
		if err = httpjson.Read(response, &birthdays); err != nil {
			return fmt.Errorf("read birthdays: %w", err)
		}

		for _, birthday := range birthdays.Data.Items {
			fmt.Printf("%s %s", birthday.FirstName, birthday.LastName)

			if len(birthday.BirthDate) > 0 {
				if dateOfBirth, err := time.Parse(dateISOFormat+"T00:00:00", birthday.BirthDate); err == nil {
					fmt.Printf("- %d years", now.Year()-dateOfBirth.Year())
				}
			}

			fmt.Printf("\n")
		}

		return nil
	},
}
