package lib

type Msg struct {
	Errmsg  string       `json:"errmsg"`
	Errcode COMMONSTATUS `json:"errcode"`
}

type COMMONSTATUS int64

const (
	STATUS_OK          COMMONSTATUS = 0
	STATUS_BAD_REQUEST COMMONSTATUS = 1
)

func WXError(errmsg string, errcode COMMONSTATUS) *Msg {
	var msg *Msg
	msg.Errmsg = errmsg
	msg.Errcode = errcode
	return msg
}
