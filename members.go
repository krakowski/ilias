package ilias

type MemberService service

type CourseMember struct {
	Identifier string
	Username   string
	Firstname  string
	Lastname   string
	Role       string
}

func (s *CourseMember) ToRow() []string {
	return []string{s.Identifier, s.Username, s.Firstname, s.Lastname, s.Role}
}

