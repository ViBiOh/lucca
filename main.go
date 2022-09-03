package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/ViBiOh/httputils/v4/pkg/httpjson"
	"github.com/ViBiOh/httputils/v4/pkg/request"
)

type ownerReponse struct {
	ID int `json:"id"`
}

type leaveAccount struct {
	LeaveAccountID   int     `json:"leaveAccountId"`
	LeaveAccountName string  `json:"leaveAccountName"`
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
	DaysUnit               bool           `json:"daysUnit"`
	Duration               float64        `json:"duration"`
	OwnerID                int            `json:"ownerId"`
	StartOn                string         `json:"startsOn"`
	EndsOn                 string         `json:"endsOn"`
	StartsAM               bool           `json:"startsAM"`
	EndsAM                 bool           `json:"endsAM"`
	IsHalfDay              bool           `json:"isHalfDay"`
	AutoCreate             bool           `json:"autoCreate"`
	Unit                   int            `json:"unit"`
	AvailableAccounts      []string       `json:"availableAccounts"`
	OtherAvailableAccounts []leaveAccount `json:"otherAvailableAccounts"`
	Users                  users          `json:"users"`
}

const (
	dateISOFormat = "2006-01-02"

	baseURL       = "https://%s.ilucca.net"
	csrfTokenName = "__RequestVerificationToken"

	antiForgeryPrefix   = ".AspNetCore.Antiforgery"
	authTokenCookieName = "authToken"
)

var csrfGetter = regexp.MustCompile(`name="` + csrfTokenName + `" type="hidden" value="(.*?)"`)

func main() {
	fs := flag.NewFlagSet("remote-lucca", flag.ExitOnError)

	subdomain := fs.String("subdomain", "", "Sub domain used")
	username := fs.String("username", "", "Username")
	password := fs.String("password", "", "Password")
	leaveType := fs.String("leaveType", "Télétravail", "Type of leave request")
	days := fs.String("days", "", "Days of week, comma separated")
	start := fs.String("start", "", "Start of repetition, in ISO format")
	end := fs.String("end", "", "End of repetition, in ISO format")
	dryRun := fs.Bool("dry-run", false, "Dry run")

	if err := fs.Parse(os.Args[1:]); err != nil {
		log.Fatal(err)
	}

	startDate, err := time.Parse(dateISOFormat, *start)
	if err != nil {
		log.Fatalf("invalid start date: %s", err)
	}

	endDate, err := time.Parse(dateISOFormat, *end)
	if err != nil {
		log.Fatalf("invalid end date: %s", err)
	}

	recurringDays := parseDaysOfWeek(*days)

	req := request.Get(fmt.Sprintf(baseURL, *subdomain))

	authToken, err := getAuthToken(req, *username, *password)
	if err != nil {
		log.Fatal(err)
	}

	req = req.Header("Cookie", fmt.Sprintf("%s=%s", authTokenCookieName, authToken))

	ownerID, err := getOwnerId(req)
	if err != nil {
		log.Fatalf("owner: %s", err)
	}

	leaveRequestType, err := getLeaveRequestType(req, ownerID, startDate.Format(dateISOFormat), *leaveType)
	if err != nil {
		log.Fatalf("type: %s", err)
	}

	for _, day := range recurringDays {
		startDate := startDate

		for {
			next := NextWeekDay(day, startDate)
			if next.After(endDate) {
				break
			}

			log.Printf("Creating `%s` leave on %s...", *leaveType, next.Format(dateISOFormat))

			if *dryRun {
				log.Printf("Dry run, no action taken.")
			} else {
				if err := createLeaveRequest(req, ownerID, leaveRequestType, next); err != nil {
					log.Fatalf("create: %s", err)
				}
			}

			startDate = next
		}
	}
}

func getOwnerId(req request.Request) (int, error) {
	resp, err := req.Path("/identity/api/principal").Send(context.Background(), nil)
	if err != nil {
		return 0, fmt.Errorf("get: %w", err)
	}

	var content ownerReponse
	if err := httpjson.Read(resp, &content); err != nil {
		return 0, fmt.Errorf("read: %w", err)
	}

	return content.ID, nil
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

func getAuthToken(req request.Request, username, password string) (string, error) {
	req = req.Path("/identity/login")
	resp, err := req.Send(context.Background(), nil)
	if err != nil {
		return "", fmt.Errorf("get login: %w", err)
	}

	payload, err := request.ReadBodyResponse(resp)
	if err != nil {
		return "", fmt.Errorf("parse get: %w", err)
	}

	parts := csrfGetter.FindAllStringSubmatch(string(payload), -1)
	if len(parts) == 0 {
		return "", errors.New("crsft token not found")
	}

	loginForm := url.Values{}
	loginForm.Set("UserName", username)
	loginForm.Set("Password", password)
	loginForm.Set(csrfTokenName, parts[0][1])

	var antiforgeryCookie string
	for _, cookie := range resp.Cookies() {
		if strings.HasPrefix(cookie.Name, antiForgeryPrefix) {
			antiforgeryCookie = cookie.Raw
		}
	}

	resp, err = req.Method(http.MethodPost).Header("Cookie", antiforgeryCookie).Form(context.Background(), loginForm)
	if err != nil {
		return "", fmt.Errorf("login: %w", err)
	}

	for _, cookie := range resp.Cookies() {
		if strings.EqualFold(cookie.Name, authTokenCookieName) {
			return cookie.Value, nil
		}
	}

	return "", errors.New("no auth token found")
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

func parseDaysOfWeek(daysOfWeek string) (output []time.Weekday) {
	for _, day := range strings.Split(daysOfWeek, ",") {
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