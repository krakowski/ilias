package ilias

import (
	"io/ioutil"
	"net/http"
)

const (
	downloadPath string = "ilias.php?vw=1&cmd=downloadReturned&cmdClass=ilexsubmissionfilegui&cmdNode=b8:lt:b9:bn&baseClass=ilExerciseHandlerGUI"
)

type Submission struct {
	ContentType	   string
	Content 	   []byte
}

type DownloadParams struct {
	Reference	string	`schema:"ref_id"`
	Assignment	string  `schema:"ass_id"`
	Member		string	`schema:"member_id"`
}

func (exercise *ExerciseService) Download(params *DownloadParams) (*Submission, error) {

	// Prepare request url
	path, err := addQueryParams(downloadPath, params)
	if err != nil {
		return nil, err
	}

	// Create request
	req, err := exercise.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	// Retrieve the HTML source
	resp, err := exercise.client.Do(req)
	if err != nil {
		return nil, err
	}

	// Read the responses content
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &Submission{
		ContentType: resp.Header.Get("Content-Type"),
		Content:     body,
	}, nil
}
