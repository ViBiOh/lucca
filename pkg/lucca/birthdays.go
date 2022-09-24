package lucca

import (
	"context"
	"fmt"
	"time"

	"github.com/ViBiOh/httputils/v4/pkg/httpjson"
)

func (a App) GetBirthdays(ctx context.Context, start, end time.Time) ([]User, error) {
	response, err := a.req.Path("/api/v3/users/birthday?fields=id,name,firstname,lastname,birthDate,dtContractStart&startsOn=%s&endsOn=%s", start.Format(dateISOFormat), end.Format(dateISOFormat)).Send(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}

	var output users
	if err = httpjson.Read(response, &output); err != nil {
		return nil, fmt.Errorf("read: %w", err)
	}

	birthdays := output.Data.Items
	for index := range birthdays {
		if date, err := time.Parse(dateTimeFormat, birthdays[index].BirthDate); err == nil {
			birthdays[index].Date = date
		}

		if date, err := time.Parse(dateTimeFormat, birthdays[index].ContractStartDate); err == nil {
			birthdays[index].ContractStart = date
		}
	}

	return birthdays, nil
}
