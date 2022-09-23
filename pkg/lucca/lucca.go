package lucca

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/ViBiOh/flags"
	"github.com/ViBiOh/httputils/v4/pkg/request"
)

const (
	baseURL = "https://%s.ilucca.net"

	csrfTokenName       = "__RequestVerificationToken"
	antiForgeryPrefix   = ".AspNetCore.Antiforgery"
	authTokenCookieName = "authToken"
)

var csrfGetter = regexp.MustCompile(`name="` + csrfTokenName + `" type="hidden" value="(.*?)"`)

type App struct {
	req request.Request
}

type Config struct {
	subdomain *string
	username  *string
	password  *string
}

func Flags(fs *flag.FlagSet, prefix string, overrides ...flags.Override) Config {
	return Config{
		subdomain: flags.String(fs, prefix, "lucca", "Subdomain", "Subdomain of company", "", overrides),
		username:  flags.String(fs, prefix, "lucca", "Username", "Username", "", overrides),
		password:  flags.String(fs, prefix, "lucca", "Subdomain", "Password", "", overrides),
	}
}

func New(config Config) (App, error) {
	return NewFromValues(*config.subdomain, *config.username, *config.password)
}

func NewFromValues(subdomain, username, password string) (App, error) {
	req := request.Get(fmt.Sprintf(baseURL, subdomain))

	authToken, err := getAuthToken(req, strings.TrimSpace(username), password)
	if err != nil {
		return App{}, fmt.Errorf("get token: %w", err)
	}

	return App{
		req: req.Header("Cookie", fmt.Sprintf("%s=%s", authTokenCookieName, authToken)),
	}, nil
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
