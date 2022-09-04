package cmd

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ViBiOh/httputils/v4/pkg/httpjson"
	"github.com/ViBiOh/httputils/v4/pkg/request"
)

type Principal struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Mail      string `json:"mail"`
}

func getPrincipal(req request.Request) (output Principal, err error) {
	var response *http.Response
	response, err = req.Path("/identity/api/principal").Send(context.Background(), nil)
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
