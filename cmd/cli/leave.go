package cli

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var (
	dateISOFormat = "2006-01-02"

	leaveType string
	days      []string
	start     string
	end       string
)

var leaveCmd = &cobra.Command{
	Use:   "leave",
	Short: "Create a leave request",
	RunE: func(cmd *cobra.Command, args []string) error {
		startDate, err := time.Parse(dateISOFormat, start)
		if err != nil {
			return fmt.Errorf("invalid start date: %w", err)
		}

		endDate, err := time.Parse(dateISOFormat, end)
		if err != nil {
			return fmt.Errorf("invalid end date: %w", err)
		}

		ctx := context.Background()

		principal, err := luccaClient.Principal(ctx)
		if err != nil {
			return fmt.Errorf("principal: %w", err)
		}

		leaveRequestType, err := luccaClient.GetLeaveRequestType(ctx, principal.ID, startDate, leaveType)
		if err != nil {
			return fmt.Errorf("get leave type: %w", err)
		}

		recurringDays := parseDaysOfWeek(days)

		for _, day := range recurringDays {
			startDate := startDate

			for {
				next := NextWeekDay(day, startDate)
				if next.After(endDate) {
					break
				}

				fmt.Printf("Creating `%s` leave on %s...", leaveType, next.Format(dateISOFormat))

				if dryRun {
					fmt.Printf(" Dry run, no action taken.\n")
				} else {
					fmt.Printf("\n")

					if err := luccaClient.CreateLeaveRequest(ctx, principal.ID, leaveRequestType, next); err != nil {
						return fmt.Errorf("create: %w", err)
					}
				}

				startDate = next
			}
		}

		return nil
	},
}

func initLeave() {
	flags := leaveCmd.PersistentFlags()

	flags.StringVarP(&leaveType, "leaveType", "", "Télétravail", "Type of leave")
	flags.StringSliceVarP(&days, "days", "", nil, "Days of week, for repetition")
	flags.StringVarP(&start, "start", "", "", "Start of leave, in ISO format")
	flags.StringVarP(&end, "end", "", "", "End of leave, in ISO format")
}

func NextWeekDay(day time.Weekday, t time.Time) time.Time {
	diff := int(day) - int(t.Weekday())
	if diff <= 0 {
		diff += 7
	}

	return t.AddDate(0, 0, diff)
}

func parseDaysOfWeek(days []string) (output []time.Weekday) {
	for _, day := range days {
		switch strings.ToLower(day) {
		case "monday":
			output = append(output, time.Monday)
		case "tuesday":
			output = append(output, time.Tuesday)
		case "wednesday":
			output = append(output, time.Wednesday)
		case "thursday":
			output = append(output, time.Thursday)
		case "friday":
			output = append(output, time.Friday)
		}
	}

	return
}
