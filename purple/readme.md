用go封装libpurple接口API

或者再抛出其他脚本语言的API，像python, ruby等。

go run examples/bot.go -u yourname@irc.freenode.net

前3个commit，分别实现基本登陆，收发群消息，dbus事件功能。

初期代码比较直接，比较粗糙，可能逻辑上并不是很好。以后再调整代码结构。



