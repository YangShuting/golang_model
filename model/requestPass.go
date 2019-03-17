package model

type LoginParams struct {
	Code      string `form:"jscode"`
	Telephone string `form:tel`
}

type LoginReturn struct {
	User    *KuaiMaoUser `json:"user"`
	Session string       `json:"session"`
}
