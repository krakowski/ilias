package ilias

import (
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

const (
	loginPath string = "ilias.php?lang=de&client_id=UniRZ&cmd=post&cmdClass=ilstartupgui&cmdNode=11f&baseClass=ilStartUpGUI&rtoken="
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

	return &User{
		Username: username,
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
	title, exists := doc.Find("#userlog img").Eq(0).Attr("alt")
	if !exists {
		return "", ErrFullName
	}

	return title, nil
}
