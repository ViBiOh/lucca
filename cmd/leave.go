package cmd

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ViBiOh/httputils/v4/pkg/httpjson"
	"github.com/ViBiOh/httputils/v4/pkg/request"
	"github.com/spf13/cobra"
)

var (
	dateISOFormat = "2006-01-02"

	leaveType string
	days      []string
	start     string
	end       string
)

type leaveAccount struct {
	LeaveAccountName string  `json:"leaveAccountName"`
	LeaveAccountID   int     `json:"leaveAccountId"`
	Duration         float64 `json:"duration"`
	Unit             int     `json:"unit"`
}

type leaveRequestTypeResponse struct {
	OtherAvailableAccounts []leaveAccount `json:"otherAvailableAccounts"`
}

type users struct {
	IDs []int `json:"userIds"`
}

type leaveRequestRequest struct {
	StartOn                string         `json:"startsOn"`
	EndsOn                 string         `json:"endsOn"`
	AvailableAccounts      []string       `json:"availableAccounts"`
	Users                  users          `json:"users"`
	OtherAvailableAccounts []leaveAccount `json:"otherAvailableAccounts"`
	Unit                   int            `json:"unit"`
	OwnerID                int            `json:"ownerId"`
	Duration               float64        `json:"duration"`
	EndsAM                 bool           `json:"endsAM"`
	IsHalfDay              bool           `json:"isHalfDay"`
	AutoCreate             bool           `json:"autoCreate"`
	DaysUnit               bool           `json:"daysUnit"`
	StartsAM               bool           `json:"startsAM"`
}

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

		principal, err := getPrincipal(req)
		if err != nil {
			return fmt.Errorf("principal: %w", err)
		}

		leaveRequestType, err := getLeaveRequestType(req, principal.ID, startDate.Format(dateISOFormat), leaveType)
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

					if err := createLeaveRequest(req, principal.ID, leaveRequestType, next); err != nil {
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

func createLeaveRequest(req request.Request, ownerID, leaveRequestTpe int, date time.Time) error {
	payload := leaveRequestRequest{
		DaysUnit:          true,
		Duration:          1,
		OwnerID:           ownerID,
		StartOn:           date.Format(dateISOFormat) + "T00:00:00",
		EndsOn:            date.Format(dateISOFormat) + "T00:00:00",
		StartsAM:          true,
		EndsAM:            false,
		IsHalfDay:         false,
		AutoCreate:        true,
		Unit:              0,
		AvailableAccounts: []string{},
		OtherAvailableAccounts: []leaveAccount{
			{
				LeaveAccountID: leaveRequestTpe,
				Duration:       1,
				Unit:           0,
			},
		},
		Users: users{
			IDs: []int{ownerID},
		},
	}

	resp, err := req.Path("/api/v3/leaveRequestFactory").Method(http.MethodPost).JSON(context.Background(), payload)
	if err != nil {
		return fmt.Errorf("create: %w", err)
	}

	if discardErr := request.DiscardBody(resp.Body); discardErr != nil {
		return fmt.Errorf("discard: %w", err)
	}

	return nil
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

func getLeaveRequestType(req request.Request, ownerID int, date, name string) (int, error) {
	resp, err := req.Path(fmt.Sprintf("/api/v3/services/leaveRequestFactory?ownerId=%d&startsOn=%s&startsAM=true&endsOn=%s&endsAM=false", ownerID, date, date)).Send(context.Background(), nil)
	if err != nil {
		return 0, fmt.Errorf("get: %w", err)
	}

	var leaveTypes leaveRequestTypeResponse
	if err := httpjson.Read(resp, &leaveTypes); err != nil {
		return 0, fmt.Errorf("parse: %w", err)
	}

	for _, leaveType := range leaveTypes.OtherAvailableAccounts {
		if strings.EqualFold(leaveType.LeaveAccountName, name) {
			return leaveType.LeaveAccountID, nil
		}
	}

	return 0, fmt.Errorf("leave type `%s` not found", name)
}
