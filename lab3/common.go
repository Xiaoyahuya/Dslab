package lab3

// LogEntry 日志条目
type LogEntry struct {
	Index   int    // 日志在切片中的索引 (类似于 ID)
	Command string // 具体命令，如 "SET X=1"
}

// AppendArgs RPC 请求参数
type AppendArgs struct {
	Entry LogEntry
}

// AppendReply RPC 响应参数
type AppendReply struct {
	Success bool
}