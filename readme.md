
libpurple's go/golang binding. Now can write pidgin plugin with go/golang.

The go-purple wrapper libpurple API in go/golang. And supplies an OO style API. 

Now implement most API function for write new IM protocol and simple client.

### Install

    go get -u github.com/kitech/go-purple/purple
    
### examples

    import "github.com/kitech/go-purple/purple"

### sub projects

* [purple](purple/) the main binding project
* [purpled](purpled/) might be a purple daemon with a web ui interface
* [tox-prpl](tox-prpl/) tox protocol implement with go-purple
* [wechat-prpl](wechat-prpl/) wechat prpl plugin
* [wechat-prpl/wechat](wechat-prpl/wechat/) wechat web protocol implement
* [yobot](yobot/) an IM robot that relay messages between different protocols, and build on go-purple.

### Notes

* using dbus path different as default pidgin
* using profile directory different as default profile

### TODOs

[ ] full libpurple support
[ ] improve memory management, avoid memory leaks.
[ ] add documents.
