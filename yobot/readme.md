
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
- [ ] preconnection
- [ ] multiple tox instance
- [ ] 逻辑与连接分开，可以重加载逻辑模块而不丢失连接。
- [ ] 消息颜色格式转换
- [x] action同步处理
- [x] 消息颜色格式过滤
- [ ] features collection for channels
- [ ] 发送错误的消息重发
- [ ] 快速掉线检测与重连(目前有时要15分钟才能检测到，虽然irc实现自带了ping功能？)
- [ ] 进出群事件通知
- [ ] 运行时配置修改
- [ ] 错过的消息的处理。短暂离开tox群，在超时前再次进入时的情况。

### 运行


