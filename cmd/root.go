package cmd

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/ViBiOh/httputils/v4/pkg/request"
	"github.com/spf13/cobra"
)

const (
	baseURL = "https://%s.ilucca.net"

	csrfTokenName       = "__RequestVerificationToken"
	antiForgeryPrefix   = ".AspNetCore.Antiforgery"
	authTokenCookieName = "authToken"
)

var (
	req request.Request

	subdomain string
	username  string
	password  string
	dryRun    bool

	csrfGetter = regexp.MustCompile(`name="` + csrfTokenName + `" type="hidden" value="(.*?)"`)
)

var rootCmd = &cobra.Command{
	Use:   "lucca",
	Short: "Run Lucca action fro the CLI",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		if parent := cmd.Parent(); parent != nil && parent.Name() == "completion" {
			return
		}

		req = request.Get(fmt.Sprintf(baseURL, subdomain))

		authToken, err := getAuthToken(req, username, password)
		if err != nil {
			return fmt.Errorf("get token: %w", err)
		}

		req = req.Header("Cookie", fmt.Sprintf("%s=%s", authTokenCookieName, authToken))

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		principal, err := getPrincipal(req)
		if err != nil {
			return fmt.Errorf("principal: %w", err)
		}

		fmt.Printf("Hello %s\n", principal.FirstName)

		return nil
	},
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

func init() {
	flags := rootCmd.PersistentFlags()

	flags.StringVarP(&subdomain, "subdomain", "", "", "Subdomain")
	flags.StringVarP(&username, "username", "", "", "Username")
	flags.StringVarP(&password, "password", "", "", "Password")
	flags.BoolVarP(&dryRun, "dry-run", "", false, "Dry run")

	rootCmd.AddCommand(birthdaysCmd)

	rootCmd.AddCommand(leaveCmd)
	initLeave()
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
	}
}
