package lucca

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ViBiOh/httputils/v4/pkg/httpjson"
)

type Principal struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Mail      string `json:"mail"`
	ID        int    `json:"id"`
}

func (a App) Principal(ctx context.Context) (output Principal, err error) {
	var response *http.Response
	response, err = a.req.Path("/identity/api/principal").Send(ctx, nil)
	if err != nil {
		err = fmt.Errorf("get principal: %w", err)

		return
	}

	if err = httpjson.Read(response, &output); err != nil {
		err = fmt.Errorf("read principal: %w", err)

		return
	}

	return
}
