package ilias

import (
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"net/url"
	"strconv"
)


const (
	gradesPath = "https://ilias.hhu.de/ilias/ilias.php?cmd=post&cmdClass=ilexercisemanagementgui&cmdNode=at:ld:au&baseClass=ilExerciseHandlerGUI"
)

type GradesQuery struct {
	Reference	string  `schema:"ref_id"`
	Assignment	string	`schema:"ass_id"`
	Token	    string	`schema:"rtoken"`
}

func (exercise *ExerciseService) UpdateGrades(params *GradesQuery, corrections []Correction) error {

	// Prepare request url
	path, err := addQueryParams(gradesPath, params)
	if err != nil {
		return err
	}


	// Magic voodo parameters
	values := url.Values{
		"selected_cmd": {"saveStatusSelected"},
		"selected_cmd2": {"saveStatusSelected"},
		"select_cmd2": {"Ausführen"},
	}

	for _, correction := range corrections {
		values.Add("member[" + correction.Student + "]", "1")
		values.Add("id[" + correction.Student + "]", "1")
		values.Add("idlid[" + correction.Student + "]", "")
		values.Add("status[" + correction.Student + "]", "passed")
		values.Add("mark[" + correction.Student + "]", strconv.FormatFloat(correction.Points, 'f', -1, 64) + " Punkte")
		values.Add("notice[" + correction.Student + "]", "")
	}

	req, err := exercise.client.NewRequest(http.MethodPost, path, values)
	if err != nil {
		return err
	}

	resp, err := exercise.client.Do(req)
	if err != nil {
		return err
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}

	if !containsSuccessAlert(doc) {
		return ErrUpdate
	}

	return nil
}