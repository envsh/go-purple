

libpurple's go/golang binding. Now can write pidgin plugin with go/golang.

## Features
* irc 连接
* title bot
* tox 连接
* irc <=> tox 互通
* 更多 IM 互通

## Tech
也许嵌入python脚本来做逻辑
也许使用独立的python脚本做逻辑

## Design

#### 非tox官方群组
从tox方可以加入任意的irc群组，采用invite #channel命令方式。

如果当前有这个tox群组，则邀请命令发送方进入现有的群里。

如果当前还没有这个tox群组，则创建一个新群组，邀请命令发送方进入，并连接到irc 对应的 #channel中。

创建之后，该程序会一直维持该群组，直接该程序掉线退出，同时保持与irc对应#channel的同步关系。

#### tox官方群组

自动进入官方群组，并通过映射打通对应的irc #channel，一直保持到退出该群组，并进程重新连接。

#### 同步用户名问题

本bot中存在两种账号，一种是root账号，另一种非root账号，同时包括bot与irc。

这样就能够在转发的时候，使用对应的名字了，在做转发时显示的名字更友好。

不过由于tox的资源占用内存资源比较多，可能每1G内存只能创建50实例/账号。

所以非root账号必须设计为可以回收重用的。

比如，当一个irc用户离开了群级

