# go websocket chatroom

初步完成go websocket json的聊天室，websocket模块使用了github.com/gorilla/websocket，因此需要先安装此模块

```shell
go get github.com/gorilla/websocket
```

在linux下面可以使用如下命令编译出可执行程序chat，然后再运行chat

```shell
# sh build.sh
# ./chat
2018/01/09 17:55:00 websocket-chat.go:126: start server listen in http://localhost:11888!
```

程序中使用了go的chan和goroutine，以达到期望中的高并发性能。

暂时没有使用到可持久化，后续有兴趣再加!