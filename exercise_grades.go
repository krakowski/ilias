package ilias

import (
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)


const (
	gradesPath = "ilias.php?vw=1&cmd=post&cmdClass=ilexercisemanagementgui&cmdNode=c4:ns:c5&baseClass=ilExerciseHandlerGUI&fallbackCmd=saveStatusAll"
	gradesOverviewPath = "ilias.php?cmd=showGradesOverview&cmdClass=ilexercisemanagementgui&cmdNode=c4:ns:c5&baseClass=ilExerciseHandlerGUI"
)

var (

	idRegex = regexp.MustCompile(`(?P<forename>.+),\s(?P<surname>.+)\s\[(?P<id>\w+)\]`)
)

type GradesUpdateQuery struct {
	Reference	string  `schema:"ref_id"`
	Assignment	string	`schema:"ass_id"`
	Token	    string	`schema:"rtoken"`
}

type GradesExportQuery struct {
	Reference	string  `schema:"ref_id"`
}



type Grading struct {
	Id		string
	Forename	string
	Surname		string
	Grades		[]string
}

func (grading *Grading) ToRow() []string {
	row := []string{grading.Id, grading.Forename, grading.Surname}
	for _, value := range grading.Grades {
		row = append(row, value)
	}

	return row
}

func (grading *Grading) ToHeader() []string {
	row := []string{"Kennung", "Vorname", "Nachname"}
	for index, _ := range grading.Grades {
		row = append(row, strconv.Itoa(index))
	}

	return row
}

func (exercise *ExerciseService) UpdateGrades(params *GradesUpdateQuery, corrections []Correction) error {

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

func (exercise *ExerciseService) Export(query *GradesExportQuery) ([]Grading, error) {

	// Prepare request url
	path, err := addQueryParams(gradesOverviewPath, query)
	if err != nil {
		return nil, err
	}

	req, err := exercise.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := exercise.client.Do(req)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	return extractGradings(doc)
}

func extractGradings(doc *goquery.Document) ([]Grading, error) {

	var gradings []Grading

	selection := doc.Find("table[id^='exc_grades_']")
	entries := selection.Find("tbody tr")
	entries.Each(func(i int, selection *goquery.Selection) {

		// Select all columns within the current row
		nodes := selection.Find("td")
		columns := nodes.Length() - 3

		match := idRegex.FindStringSubmatch(strings.TrimSpace(nodes.Eq(0).Text()))
		result := make(map[string]string)
		for i, name := range idRegex.SubexpNames() {
			if i != 0 && name != "" {
				result[name] = match[i]
			}
		}

		grading := Grading{
			Id:       strings.TrimSpace(result["id"]),
			Forename: strings.TrimSpace(result["forename"]),
			Surname:  strings.TrimSpace(result["surname"]),
			Grades: []string{},
		}

		for i := 1; i < columns + 1; i++ {
			points := strings.TrimSpace(nodes.Eq(i).Text())
			if len(points) == 0 {
				points = "0 Punkte"
			}

			grading.Grades = append(grading.Grades, points)
		}

		gradings = append(gradings, grading)
	})

	return gradings, nil
}