package ilias

import (
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"net/url"
	"strings"
)

const (
	loginPath string = "ilias.php?lang=de&client_id=UniRZ&cmd=post&cmdClass=ilstartupgui&cmdNode=bx&baseClass=ilStartUpGUI&rtoken="
)

func (auth *AuthService) Login(username string, password string) (*User, error) {
	request, err := auth.client.NewRequest(http.MethodPost, loginPath, url.Values{
		"username": { username },
		"password": { password },
		"cmd[doStandardAuthentication]": { "Anmelden" },
	})

	if err != nil {
		return nil, err
	}

	resp, err := auth.client.client.Do(request)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	if containsDangerAlert(doc) {
		return nil, ErrCredentials
	}

	token, err := findToken(doc)
	if err != nil {
		return nil, err
	}

	fullName, err := findFullName(doc)
	if err != nil {
		return nil, err
	}

	split := strings.Split(fullName, " ")
	return &User{
		Username: username,
		Firstname: strings.Join(split[:len(split) - 1], " "),
		Lastname: split[len(split) - 1],
		Token:    token,
	}, nil
}

func findToken(doc *goquery.Document) (string, error) {
	action, exists := doc.Find("#mm_search_form").Eq(0).Attr("action")
	if !exists {
		return "", ErrToken
	}

	values, err := url.ParseQuery(action)
	if err != nil {
		return "", nil
	}

	return values.Get("rtoken"), nil
}

func findFullName(doc *goquery.Document) (string, error) {
	title, exists := doc.Find("#userlog img").Eq(0).Attr("title")
	if !exists {
		return "", ErrFullName
	}

	return title, nil
}
