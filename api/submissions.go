package api

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

const (
	ExerciseBaseUrl string = "ilias.php?cmd=members&cmdClass=ilexercisemanagementgui&cmdNode=ul:us:um&baseClass=ilExerciseHandlerGUI"
	ExerciseParam string = "ref_id"
	AssignmentParam string = "ass_id"
)

type SubmissionInfo struct {
	Firstname      string
	Lastname	   string
	UserId         string
	Date		   string
}

func (s *SubmissionInfo) ToRow() []string {
	return []string{s.UserId, s.Lastname, s.Firstname, s.Date}
}

func (c *Client) GetSubmissions(exerciseId string, assignmentId string, includeEmpty bool) ([]SubmissionInfo, error) {
	// Create the request URL
	u, err := c.BaseURL.Parse(ExerciseBaseUrl + "&" + selectAssignment(exerciseId, assignmentId))
	if err != nil {
		return nil, err
	}

	// Retrieve the HTML source
	resp, err := c.httpClient.Get(u.String())
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)

	table := findSubmissionTable(doc)
	if table == nil {
		return nil, nil
	}

	return readSubmissions(table, includeEmpty), nil
}


func findSubmissionTable(doc *goquery.Document) *goquery.Selection {
	return doc.Find("#exc_mem tbody tr")
}

func readSubmissions(selction *goquery.Selection, includeEmpty bool) []SubmissionInfo {
	var submissions []SubmissionInfo
	selction.Each(func(i int, selection *goquery.Selection) {
		nodes := selection.Find("td")
		splitName := strings.Split(nodes.Eq(1).Text(), ",")
		submission := SubmissionInfo{
			Lastname: strings.TrimSpace(splitName[0]),
			Firstname:   strings.TrimSpace(splitName[1]),
			UserId: strings.TrimSpace(nodes.Eq(2).Text()),
			Date: strings.TrimSpace(nodes.Eq(3).Text()),
		}

		if len(submission.Date) == 0 && !includeEmpty {
			return
		}

		submissions = append(submissions, submission)
	})
	return submissions
}

func selectAssignment(exerciseId string, assignmentId string) string {
	return fmt.Sprintf("%s=%s&%s=%s", ExerciseParam, exerciseId, AssignmentParam, assignmentId)
}