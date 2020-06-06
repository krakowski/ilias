package ilias

type ExerciseService service

type Correction struct {
	Student		string	`yaml:"student"`
	Points		float64	`yaml:"points"`
	Corrected	bool	`yaml:"corrected"`
	Correction	string	`yaml:"correction,flow"`
}

