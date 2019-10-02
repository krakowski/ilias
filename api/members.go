package api

import (
	"github.com/PuerkitoBio/goquery"
	"strings"
)

const (
	MembersBaseUrl string = "ilias.php?cmdClass=ilcoursemembershipgui&cmdNode=hb:t3:su&baseClass=ilrepositorygui"
	CourseIdParam  string = "ref_id"
)

type MemberInfo struct {
	UserId         string
	Firstname      string
	Lastname	   string
	Role 		   string
}

func (s *MemberInfo) ToRow() []string {
	return []string{s.UserId, s.Firstname, s.Lastname, s.Role}
}

func (c *Client) GetMembers(courseId string, includeEmpty bool) ([]MemberInfo, error) {
	// Create the request URL
	u, err := c.BaseURL.Parse(MembersBaseUrl + "&" + CourseIdParam + "=" + courseId)
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

	table := findMembersTable(doc)
	if table == nil {
		return nil, nil
	}

	return readMembers(table, includeEmpty), nil
}


func findMembersTable(doc *goquery.Document) *goquery.Selection {
	selection := doc.Find("#participants");
	if selection == nil {
		return nil;
	}

	return selection.Find(".table-responsive > table > tbody > tr")
}

func readMembers(selction *goquery.Selection, includeEmpty bool) []MemberInfo {
	var members []MemberInfo
	selction.Each(func(i int, selection *goquery.Selection) {
		nodes := selection.Find("td")
		splitName := strings.Split(nodes.Eq(1).Text(), ",")
		member := MemberInfo{
			Lastname: strings.TrimSpace(splitName[0]),
			Firstname:   strings.TrimSpace(splitName[1]),
			UserId: strings.TrimSpace(nodes.Eq(2).Text()),
			Role: strings.TrimSpace(nodes.Eq(3).Text()),
		}

		members = append(members, member)
	})
	return members
}