package lucca

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/ViBiOh/httputils/v4/pkg/httpjson"
)

type Operation struct {
	ID uint `json:"id"`
}

type AppInstance struct {
	ApplicationID string      `json:"applicationID"`
	Permissions   []Operation `json:"permissions"`
	ID            uint        `json:"id"`
}

type appinstances struct {
	Data struct {
		Items []AppInstance `json:"items"`
	} `json:"data"`
}

func (a App) GetAppInstance(ctx context.Context, appID string) (uint, uint, error) {
	response, err := a.req.Path("/api/v3/appinstances?fields=id,name,applicationID,permissions[operationID]").Send(ctx, nil)
	if err != nil {
		return 0, 0, fmt.Errorf("get: %w", err)
	}

	var output appinstances
	if err = httpjson.Read(response, &output); err != nil {
		return 0, 0, fmt.Errorf("read: %w", err)
	}

	for _, appinstance := range output.Data.Items {
		if strings.EqualFold(appinstance.ApplicationID, appID) {
			if len(appinstance.Permissions) > 0 {
				return appinstance.ID, appinstance.Permissions[0].ID, nil
			}
			return 0, 0, errors.New("no permission")
		}
	}

	return 0, 0, errors.New("not found")
}
