package lucca

import (
	"context"
	"fmt"
	"time"

	"github.com/ViBiOh/httputils/v4/pkg/httpjson"
)

type Birthday struct {
	Date              time.Time
	ContractStart     time.Time
	BirthDate         string `json:"birthDate"`
	ContractStartDate string `json:"dtContractStart"`
	Name              string `json:"name"`
	FirstName         string `json:"firstName"`
	LastName          string `json:"lastName"`
	ID                int    `json:"id"`
}

type birthdays struct {
	Data struct {
		Items []Birthday `json:"items"`
	} `json:"data"`
}

func (a App) GetBirthdays(ctx context.Context, start, end time.Time) ([]Birthday, error) {
	response, err := a.req.Path("/api/v3/users/birthday?fields=id,name,firstname,lastname,birthDate,dtContractStart&startsOn=%s&endsOn=%s", start.Format(dateISOFormat), end.Format(dateISOFormat)).Send(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}

	var birthdays birthdays
	if err = httpjson.Read(response, &birthdays); err != nil {
		return nil, fmt.Errorf("read: %w", err)
	}

	output := birthdays.Data.Items
	for index := range output {
		if date, err := time.Parse(dateTimeFormat, output[index].BirthDate); err == nil {
			output[index].Date = date
		}

		if date, err := time.Parse(dateTimeFormat, output[index].ContractStartDate); err == nil {
			output[index].ContractStart = date
		}
	}

	return output, nil
}
