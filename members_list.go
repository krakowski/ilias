package ilias

import (
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	memberListPath string = "ilias.php?cmdClass=ilcoursemembershipgui&cmdNode=vc:ks:8k&baseClass=ilrepositorygui"
)

type MemberParams struct {
	Reference		string	`schema:"ref_id"`
}

func (members *MemberService) List(params *MemberParams) ([]CourseMember, error) {

	// Prepare request url
	path, err := addQueryParams(memberListPath, params)
	if err != nil {
		return nil, err
	}

	// Create request
	req, err := members.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	// Retrieve the HTML source
	resp, err := members.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	return readMembers(doc), nil
}

func readMembers(doc *goquery.Document) []CourseMember {

	// Find participant table
	selection := doc.Find("#participants .table-responsive table tbody tr")

	// Iterate over all table rows
	var members []CourseMember
	selection.Each(func(i int, selection *goquery.Selection) {

		// Select all columns within the current row
		nodes := selection.Find("td")

		identifier, exists := nodes.Eq(0).Find("input[type=checkbox]").Eq(0).Attr("value")
		if !exists {
			return
		}

		// Split the name and extract information
		splitName := strings.Split(nodes.Eq(1).Text(), ",")
		member := CourseMember{
			Identifier: identifier,
			Lastname:  strings.TrimSpace(splitName[0]),
			Firstname: strings.TrimSpace(splitName[1]),
			Username:  strings.TrimSpace(nodes.Eq(2).Text()),
			Role:      strings.TrimSpace(nodes.Eq(3).Text()),
		}

		members = append(members, member)
	})

	return members
}
