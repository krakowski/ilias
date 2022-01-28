package ilias

type AuthService service

type User struct {
	Username	string  `schema:"-"`
	Token		string	`schema:"rtoken"`
}
