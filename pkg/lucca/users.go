package lucca

import (
	"context"
	"fmt"
	"time"

	"github.com/ViBiOh/httputils/v4/pkg/httpjson"
)

type User struct {
	Date              time.Time
	ContractStart     time.Time
	BirthDate         string `json:"birthDate"`
	ContractStartDate string `json:"dtContractStart"`
	Name              string `json:"name"`
	FirstName         string `json:"firstName"`
	LastName          string `json:"lastName"`
	ID                int    `json:"id"`
}

type users struct {
	Data struct {
		Items []User `json:"items"`
	} `json:"data"`
}

func (a App) GetCompanyBirthdays(ctx context.Context, start, end time.Time) ([]User, error) {
	appInstanceID, operationID, err := a.GetAppInstance(ctx, "DIRECTORY")
	if err != nil {
		return nil, fmt.Errorf("get app instance: %w", err)
	}

	response, err := a.req.Path("/api/v3/users/scope?appInstanceId=%d&operations=%d&fields=id,name,firstname,lastname,dtContractStart&paging=0,200", appInstanceID, operationID).Send(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}

	var output users
	if err = httpjson.Read(response, &output); err != nil {
		return nil, fmt.Errorf("read: %w", err)
	}

	var companyBirthdays []User

	for _, user := range output.Data.Items {
		if date, err := time.Parse(dateTimeFormat, user.ContractStartDate); err == nil {
			anniversaryThisyear := time.Date(start.Year(), date.Month(), date.Day(), 12, 0, 0, 0, date.Location())

			if !anniversaryThisyear.After(end) && !anniversaryThisyear.Before(start) {
				user.ContractStart = date
				companyBirthdays = append(companyBirthdays, user)
			}
		}
	}

	return companyBirthdays, nil
}
