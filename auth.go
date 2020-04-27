package ilias

type AuthService service

type User struct {
	Username	string  `schema:"-"`
	Firstname	string	`schema:"-"`
	Lastname	string	`schema:"-"`
	Token		string	`schema:"rtoken"`
}
