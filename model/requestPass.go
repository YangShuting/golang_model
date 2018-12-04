package model

type LoginParams struct {
	Telephone string `form:"tel"`
}

type LoginReturn struct {
	User    *KuaiMaoUser `json:"user"`
	Session string       `json:"session"`
}
