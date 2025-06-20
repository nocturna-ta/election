package model

type EducationHistory struct {
	InstituteName string `db:"institute_name"`
	Year          string `db:"year"`
}

type WorkHistory struct {
	InstituteName string `db:"institute_name"`
	Position      string `db:"position"`
	Year          string `db:"year"`
}

type WorkProgram struct {
	ProgramName  string   `json:"program_name"`
	ProgramPhoto string   `json:"program_photo"`
	ProgramDesc  []string `json:"program_desc"`
}
