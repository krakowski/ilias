package ilias

type CourseService service

type HeaderConfig struct {
	ShowTitleAndIcon *bool `yaml:"titleAndIcon"`
	ShowActions      *bool `yaml:"actions"`
}

type RegistrationConfig struct {
	Mode       *string `yaml:"mode"`
	Password   *string `yaml:"password"`
	EnableLink *bool   `yaml:"enableLink"`
}

type SortingConfig struct {
	Mode      string `yaml:"mode"`
	Direction string `yaml:"direction"`
}

type PresentationConfig struct {
	Mode    *string         `yaml:"mode"`
	View    *string         `yaml:"view"`
	Sorting SortingConfig  `yaml:"sorting"`
}

type FunctionConfig struct {
	Calendar      *bool `yaml:"calendar"`
	News          *bool `yaml:"news"`
	Metadata      *bool `yaml:"metadata"`
	Ratings       *bool `yaml:"ratings"`
	Badges        *bool `yaml:"badges"`
	Competences   *bool `yaml:"competences"`
	MemberGallery *bool `yaml:"memberGallery"`
	PutOnDesk     *bool `yaml:"putOnDesk"`
	WelcomeMail   *bool `yaml:"welcomeMail"`
}

type CourseSettings struct {
	Title		 string 			`yaml:"title"`
	Description	 *string 			`yaml:"description"`
	Header       HeaderConfig       `yaml:"header"`
	Presentation PresentationConfig `yaml:"presentation"`
	PassCriteria *string             `yaml:"passCriteria"`
	Functions    FunctionConfig     `yaml:"functions"`
	EcsExport	 bool				`yaml:"ecs"`
}
