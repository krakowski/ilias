package ilias

import (
	"bytes"
	"errors"
	"github.com/gorilla/schema"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/textproto"
	"net/url"
	"strings"
)

var (
	decoder = schema.NewDecoder()
	encoder = schema.NewEncoder()

	ErrCredentials = errors.New("wrong username or password")
	ErrToken = errors.New("token could not be found")
	ErrFullName = errors.New("full name could not be found")
	ErrUpdate = errors.New("update failed")
	ErrFileHash = errors.New("file hash could not be found")
)

const (
	baseUrl        string = "https://ilias.hhu.de/"
	defaultHost	   string = "ilias.hhu.de"
)

type Credentials struct {
	Username string
	Password string
}

type Client struct {

	// The current user's login name.
	User		*User

	// Base URL for requests. Should end with a dash.
	BaseURL 	*url.URL

	// Host field set within request headers
	Host		string

	// HTTP Client used for making requests against the ILIAS platform.
	client 		*http.Client

	common 		service

	Auth 		*AuthService
	Exercise 	*ExerciseService
	Members		*MemberService
	Tables		*TableService
}

type service struct {
	client *Client
}

func NewClient(client *http.Client, credentials *Credentials) (*Client, error) {
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
	ret := &Client{ BaseURL: base, Host: defaultHost, client: client }
	ret.common.client = ret
	ret.Auth = (*AuthService)(&ret.common)
	ret.Exercise = (*ExerciseService)(&ret.common)
	ret.Members = (*MemberService)(&ret.common)
	ret.Tables = (*TableService)(&ret.common)

	// Login using the client

	user, err := ret.Auth.Login(credentials.Username, credentials.Password)
	if err != nil {
		return nil, err
	}

	ret.User = user
	return ret, nil
}

func (c *Client) NewRequest(method string, path string, body url.Values) (*http.Request, error) {
	target, err := c.BaseURL.Parse(path)
	if err != nil {
		return nil, err
	}

	// https://github.com/golang/go/issues/32897
	var request *http.Request
	if body != nil {
		request, err = http.NewRequest(method, target.String(), strings.NewReader(body.Encode()))
	} else {
		request, err = http.NewRequest(method, target.String(), nil)
	}

	if err != nil {
		return nil, err
	}

	if body != nil {
		request.Header.Set("content-type", "application/x-www-form-urlencoded")
	}

	request.Host = c.Host
	return request, nil
}

type UploadFile struct {
	Header	    textproto.MIMEHeader
	Content		*bytes.Buffer
}

func (c *Client) NewMultipartRequest(method string, path string, body url.Values, upload *UploadFile) (*http.Request, error) {

	target, err := c.BaseURL.Parse(path)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Create form field for the file
	fieldWriter, err := writer.CreatePart(upload.Header)
	if err != nil {
		return nil, err
	}

	// Copy file contents into request
	if _, err = io.Copy(fieldWriter, upload.Content); err != nil {
		return nil, err
	}

	// Copy additional parameters
	for key, array := range body {
		for _, value := range array {
			_ = writer.WriteField(key, value)
		}
	}

	writer.Close()

	request, err := http.NewRequest(method, target.String(), &buf)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", writer.FormDataContentType())
	request.Host = c.Host
	return request, nil
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	return c.client.Do(req)
}
