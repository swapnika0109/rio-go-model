package model

type EmailRequest struct{
    Name string `json:"name"`
    Email string `json:"email"`
    Message string `json:"message"`
}

func (e *EmailRequest) ToMap(m map[string]interface{}) map[string]interface{} {
    if m == nil { m = make(map[string]interface{}) }
    m["name"] = e.Name
    m["email"] = e.Email
    m["message"] = e.Message
    return m
}
