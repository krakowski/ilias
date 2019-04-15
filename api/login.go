package api

import (
	"net/url"
)

const (
	FormParamLogin string = "username"
	FormParamPassword string = "password"
	FormParamCommand string = "cmd[doStandardAuthentication]"
	FormParamCommandValue string = "Anmelden"
)

const (
	EmployeLoginURL string = "ilias.php?lang=de&cmd=post&cmdClass=ilstartupgui&cmdNode=oa&baseClass=ilStartUpGUI&rtoken="
)

func (c *Client) Login(username string, password string) error {
	// Parse the request URL
	u, err := c.BaseURL.Parse(EmployeLoginURL)
	if err != nil {
		return err
	}

	// Send a POST request containing the credentials
	_, err = c.httpClient.PostForm(u.String(), url.Values {
		FormParamLogin : {username},
		FormParamPassword: {password},
		FormParamCommand: {FormParamCommandValue},
	})

	if err != nil {
		return err
	}

	return nil
}