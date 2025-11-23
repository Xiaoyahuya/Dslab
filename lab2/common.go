package lab2

// LamportMsg 在网络中传输的消息包
type LamportMsg struct {
	Timestamp int    // 发送者的逻辑时钟
	SenderID  string // 发送者 ID
	Payload   string // 消息内容 (比如 "Hello")
}

// Reply RPC 的响应 (在这个实验中通常为空，或者是 Ack)
type Reply struct {
	// 实际应用中可能包含接收者的确认信息
}