package tui

// Survey A simple questionnaire of data being asked for
type Survey struct {
	Name       string
	DOB        string
	Profession string
}

func NewSurvey() Survey {
	return Survey{
		Name:       "unknown",
		DOB:        "01/01/1970",
		Profession: "Professional Bum",
	}
}
