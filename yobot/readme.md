
a tox to irc channel bridge.

## Features
* irc 连接
* title bot
* tox 连接
* irc <=> tox 互通
* 更多 IM 互通

## Tech
也许嵌入python脚本来做逻辑
也许使用独立的python脚本做逻辑

## TODO
- [x] rename problem
- [ ] preconnection irc
- [ ] multiple tox instance
- [ ] 逻辑与连接分开，可以重加载逻辑模块而不丢失连接。
- [ ] 消息颜色格式转换
- [x] action同步处理
- [x] 消息颜色格式过滤
- [ ] features collection for channels
- [ ] 发送失败的消息重发
- [x] 快速掉线检测与重连(目前有时要15分钟才能检测到，虽然irc实现自带了ping功能？)
- [x] 进出群事件通知
- [ ] 运行时配置修改。PART某群。
- [ ] 错过的消息的处理。短暂离开tox群，在超时前再次进入时的情况。
- [ ] 记录被invite时的data信息
- [x] 消息的变换。替换掉转发bot信息。
- [ ] 主bot账号的私信消息处理。
- [ ] info 中显示群人数
- [x] welcome
- [ ] 两个bot进同一群的冲突，循环转发同一条消息。
- [ ] tox 端进群账号立即分配一个irc连接。
- [ ] hlp响应消息到不同协议对应的真实用户名。


### 运行

### 自动邀请进群功能
前提是tox机器人的好友。
在好友上线事件，但没有在指定的群组时，主动inivite进该群。
这个群应包含接通的irc群和tox官方bot的群。
这个功能，主要还是需要存储状态数据，存储好友指定自动进入哪些群。

### welcome功能
用户分类，

