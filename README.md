# ino-chat
简单的web聊天

* 使用ws连接初始化并且仅用作接收消息
* 使用http处理房间信息维护和发送消息并推送到nsq队列处理
* 使用redis维护房间信息
* mongodb持久化数据
