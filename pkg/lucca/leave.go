package lucca

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ViBiOh/httputils/v4/pkg/httpjson"
	"github.com/ViBiOh/httputils/v4/pkg/request"
)

var (
	dateISOFormat  = "2006-01-02"
	dateTimeFormat = "2006-01-02T00:00:00"
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

func (a App) GetLeaveRequestType(ctx context.Context, ownerID int, date time.Time, name string) (int, error) {
	resp, err := a.req.Path("/api/v3/services/leaveRequestFactory?ownerId=%d&startsOn=%s&startsAM=true&endsOn=%s&endsAM=false", ownerID, date.Format(dateISOFormat), date.Format(dateISOFormat)).Send(ctx, nil)
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

func (a App) CreateLeaveRequest(ctx context.Context, ownerID, leaveRequestTpe int, date time.Time) error {
	payload := leaveRequestRequest{
		DaysUnit:          true,
		Duration:          1,
		OwnerID:           ownerID,
		StartOn:           date.Format(dateTimeFormat),
		EndsOn:            date.Format(dateTimeFormat),
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

	resp, err := a.req.Path("/api/v3/leaveRequestFactory").Method(http.MethodPost).JSON(ctx, payload)
	if err != nil {
		return fmt.Errorf("create: %w", err)
	}

	if discardErr := request.DiscardBody(resp.Body); discardErr != nil {
		return fmt.Errorf("discard: %w", err)
	}

	return nil
}
