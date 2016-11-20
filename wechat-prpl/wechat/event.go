package wechat

type Event struct {
	Type  int
	SType string
	Args  []string
}

// 事后事件
const (
	EVT_NONE int = iota
	EVT_GOT_UUID
	EVT_GOT_QRCODE
	EVT_WAIT_SCAN
	EVT_SCAN_STATUS
	EVT_REDIR_URL
	EVT_LOGIN_STATUS
	EVT_GOT_BASEINFO
	EVT_GOT_CONTACT
	EVT_GOT_MESSAGE
	EVT_LOGOUT
)

func newEvent(evt int, args []string) *Event {
	this := &Event{}
	this.Type = evt
	this.Args = args
	return this
}

func newEvent2(stype string) *Event {
	this := &Event{}
	this.SType = stype
	return this
}
