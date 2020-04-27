package ilias

type ExerciseService service

type Correction struct {
	Student		string	`yaml:"student"`
	Points		int		`yaml:"points"`
	Corrected	bool	`yaml:"corrected"`
	Correction	string	`yaml:"correction,flow"`
}

