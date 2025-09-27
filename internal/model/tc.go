package model

type Tc struct {
	Accepted bool   `json:"accepted"`
	Email     string `json:"email"`
}

func NewTc(accepted bool, email string) *Tc {
	return &Tc{
		Accepted: accepted,
		Email: email,
	}
}

