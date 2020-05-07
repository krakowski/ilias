package ilias

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	commentPath = "ilias.php?vw=1&cmd=saveCommentForLearners&cmdClass=ilexercisemanagementgui&cmdNode=at:ld:au&baseClass=ilExerciseHandlerGUI&cmdMode=asynch"
)

type CommentParams struct {
	Reference	string  `schema:"ref_id"`
	Assignment	string	`schema:"ass_id"`
}

func (exercise *ExerciseService) UpdateComment(params *CommentParams, correction Correction) error {

	// Prepare request url
	path, err := addQueryParams(commentPath, params)
	if err != nil {
		return err
	}

	values := url.Values{
		"ass_id": {params.Assignment},
		"mem_id": {correction.Student},
		"comm": {correction.Correction},
	}

	req, err := exercise.client.NewRequest(http.MethodPost, path, values)
	if err != nil {
		return err
	}

	// Retrieve the HTML source
	resp, err := exercise.client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if !strings.HasPrefix(string(body), "{\"result\":true") {
		return ErrUpdate
	}

	return nil
}