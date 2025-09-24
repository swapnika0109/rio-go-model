package model

import "rio-go-model/internal/util"

type UserProfile struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	Country   string `json:"country"`
	City      string `json:"city"`
	Preferences []string `json:"preferences"`
	Religions []string `json:"religions"`
}



func (u *UserProfile) FromMap(m map[string]interface{}) *UserProfile {
	return &UserProfile{
		Username: m["username"].(string),
		Email: m["email"].(string),
		Country: m["country"].(string),
		City: m["city"].(string),
        Preferences: util.SafeStringSlice(m["preferences"]),
        Religions: util.SafeStringSlice(m["religions"]),
	}
}