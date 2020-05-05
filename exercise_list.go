package ilias

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"strings"
)

const (
	listPath string = "ilias.php?exc_mem_trows=800&cmd=members&cmdClass=ilexercisemanagementgui&cmdNode=xn:xk:125&baseClass=ilExerciseHandlerGUI"
)

var (
	memberReplacer = strings.NewReplacer(
		"member[",            "",
		"]",	              "",
	)
)

type SubmissionMeta struct {
	Identifier	string
	Firstname  	string
	Lastname   	string
	UserId     	string
	Date       	string
}

func (s *SubmissionMeta) ToRow() []string {
	return []string{s.Identifier, s.UserId, s.Lastname, s.Firstname, s.Date}
}

type ListParams struct {
	Reference		string	`schema:"ref_id"`
	Assignment		string  `schema:"ass_id"`
	IncludeEmpty 	bool	`schema:"-"`
}

func (exercise *ExerciseService) List(params *ListParams) ([]SubmissionMeta, error) {

	// Prepare request url
	path, err := addQueryParams(listPath, params)
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

	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	table := doc.Find("#exc_mem tbody tr")
	if table == nil {
		return nil, nil
	}

	return readSubmissions(table, params.IncludeEmpty), nil
}

func readSubmissions(selction *goquery.Selection, includeEmpty bool) []SubmissionMeta {
	var submissions []SubmissionMeta
	selction.Each(func(i int, selection *goquery.Selection) {
		nodes := selection.Find("td")

		// Extract the member id
		memberId, exists := nodes.Eq(0).Find("input").Eq(0).Attr("name")
		if !exists {
			log.Fatal("extracting member id failed")
		}

		splitName := strings.Split(nodes.Eq(1).Nodes[0].FirstChild.Data, ",")
		submission := SubmissionMeta{
			Identifier: memberReplacer.Replace(memberId),
			Lastname:   strings.TrimSpace(splitName[0]),
			Firstname:  strings.TrimSpace(splitName[1]),
			UserId:     strings.TrimSpace(nodes.Eq(2).Text()),
			Date:       strings.TrimSpace(nodes.Eq(3).Text()),
		}

		if len(submission.Date) == 0 && !includeEmpty {
			return
		}

		submissions = append(submissions, submission)
	})

	return submissions
}
