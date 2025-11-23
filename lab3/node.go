package lab3

import (
	"errors"
	"fmt"
	"net"
	"net/rpc"
	"sync"
)

type Node struct {
	mu        sync.Mutex
	ID        string
	Log       []LogEntry
	IsPrimary bool

	// 网络相关
	listener  net.Listener
	peerClient *rpc.Client // 假设只有一个 Peer (Backup)
}

func NewNode(id string, isPrimary bool) *Node {
	return &Node{
		ID:        id,
		Log:       make([]LogEntry, 0),
		IsPrimary: isPrimary,
	}
}

// --- 网络基础设施 (已提供，无需修改) ---

func (n *Node) StartRPCServer() string {
	server := rpc.NewServer()
	server.Register(n)
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}
	n.listener = l
	go server.Accept(l)
	return l.Addr().String()
}

func (n *Node) ConnectToPeer(addr string) error {
	client, err := rpc.Dial("tcp", addr)
	if err != nil {
		return err
	}
	n.mu.Lock()
	n.peerClient = client
	n.mu.Unlock()
	return nil
}

// Close 关闭连接 (用于模拟宕机)
func (n *Node) Close() {
	if n.listener != nil {
		n.listener.Close()
	}
	if n.peerClient != nil {
		n.peerClient.Close()
	}
}

// --- 你的任务：核心复制逻辑 ---

// AppendClient 处理客户端请求 (Entry Point)
func (n *Node) AppendClient(cmd string) error {
	n.mu.Lock()
	defer n.mu.Unlock()

	if !n.IsPrimary {
		return errors.New("operation on non-primary node")
	}

	// 1. 构造新日志
	index := len(n.Log)
	entry := LogEntry{Index: index, Command: cmd}

	// TODO: 1. 先乐观地写入本地 Log
	// n.Log = append(n.Log, entry)

	// 如果没有 Peer (单机模式)，直接返回成功
	if n.peerClient == nil {
		return nil
	}

	// 2. 同步调用 Backup
	// 注意：这里我们持有锁进行 RPC 调用。
	// 在生产级系统(Raft)中，为了性能通常会释放锁。
	// 但为了保证简单的强一致性(Sequence)，持有锁是最安全的：
	// 确保上一条没同步完之前，下一条进不来。

	// TODO: 2. 构造RPCargs和reply

	// TODO: 3. 调用 RPC "Node.HandleAppendEntries"

	// TODO: 4. 处理失败情况
	// return fmt.Errorf("replication failed")
	
	return nil
}

// HandleAppendEntries RPC Handler (Backup 接收端)
func (n *Node) HandleAppendEntries(args *AppendArgs, reply *AppendReply) error {
	n.mu.Lock()
	defer n.mu.Unlock()

	if n.IsPrimary {
		return errors.New("primary should not receive append requests")
	}

	// 简单的完整性检查：确保收到的日志索引是连续的
	// 如果 Primary 发来 Index=5，但我本地只有 Index=2，说明中间缺了数据
	if args.Entry.Index != len(n.Log) {
		reply.Success = false
		return fmt.Errorf("log mismatch: expected index %d, got %d", len(n.Log), args.Entry.Index)
	}

	// TODO: 5. 写入本地 Log
	// n.Log = append(n.Log, args.Entry)
	
	reply.Success = true
	return nil
}

// GetLog 辅助函数：获取当前日志副本
func (n *Node) GetLog() []LogEntry {
	n.mu.Lock()
	defer n.mu.Unlock()
	// 返回副本以防止并发读写
	duplicate := make([]LogEntry, len(n.Log))
	copy(duplicate, n.Log)
	return duplicate
}