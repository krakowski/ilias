package api

import (
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"syscall"
)

const (
	baseUrl        string = "https://ilias.hhu.de/ilias/"
)

type Client struct {
	BaseURL *url.URL
	httpClient *http.Client
}

func NewClient(client *http.Client) (*Client, error) {
	// Create a default client if none is specified
	if client == nil {
		client = http.DefaultClient
	}

	// Attach a cookie jar to the client
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	client.Jar = jar

	// Parse the base url and create the http client
	base, _ := url.Parse(baseUrl)
	ret := &Client{httpClient: client, BaseURL: base}

	// Get the username
	user, present := os.LookupEnv("ILIAS_USER")
	if present == false {
		fmt.Fprint(os.Stderr, "Please specify your ILIAS username: ")
		inputBytes, err := terminal.ReadPassword(syscall.Stdin)
		if err != nil {
			log.Fatal(err)
		}
		user = string(inputBytes)
		fmt.Fprintln(os.Stderr)
	}

	// Get the password
	password, present := os.LookupEnv("ILIAS_PASS")
	if present == false {
		fmt.Fprint(os.Stderr, "Please specify your ILIAS password: ")
		inputBytes, err := terminal.ReadPassword(syscall.Stdin)
		if err != nil {
			log.Fatal(err)
		}
		password = string(inputBytes)
		fmt.Fprintln(os.Stderr)
	}

	// Login using the client
	err = ret.Login(user, password)
	if err != nil {
		log.Fatal(err)
	}

	return ret, nil
}
