package main

type ToxRoom struct {
	ToxName    string
	IrcName    string
	MatrixName string
	XmppName   string
	// Disabled   bool
}

var FixedRooms = []*ToxRoom{
	&ToxRoom{"#tox-cn-test中文", "#tox-cn-test", "", ""},
	&ToxRoom{"#rtab123中文", "#rtab123", "", ""},
	&ToxRoom{"#testk中文", "#testk", "", ""},
	&ToxRoom{"onlytox中文", "", "", ""},
	&ToxRoom{"onlytox2中文", "", "", ""},
	&ToxRoom{"onlytox3中文", "", "", ""},
	// &ToxRoom{},
}
