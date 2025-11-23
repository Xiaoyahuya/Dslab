package lab2

import (
	"fmt"
	"net"
	"net/rpc"
	"sync"
)

// LogEntry 本地日志条目
type LogEntry struct {
	Clock int
	Msg   string
	Node  string
}

type Node struct {
	mu    sync.Mutex
	ID    string
	Clock int
	Logs  []LogEntry

	// RPC 相关
	listener net.Listener
	peers    map[string]*rpc.Client // 存储其他节点的连接: NodeID -> Client
}

func NewNode(id string) *Node {
	return &Node{
		ID:    id,
		Clock: 0,
		Logs:  make([]LogEntry, 0),
		peers: make(map[string]*rpc.Client),
	}
}

// StartRPCServer 启动 RPC 监听 (已提供，无需修改)
func (n *Node) StartRPCServer() string {
	rpcServer := rpc.NewServer()
	rpcServer.Register(n)

	l, err := net.Listen("tcp", ":0") // 随机端口
	if err != nil {
		panic(err)
	}
	n.listener = l
	
	go func() {
		for {
			conn, err := n.listener.Accept()
			if err != nil {
				return
			}
			go rpcServer.ServeConn(conn)
		}
	}()
	return l.Addr().String()
}

// ConnectToPeer 连接到另一个节点 (已提供，无需修改)
func (n *Node) ConnectToPeer(peerID, addr string) error {
	client, err := rpc.Dial("tcp", addr)
	if err != nil {
		return err
	}
	n.mu.Lock()
	n.peers[peerID] = client
	n.mu.Unlock()
	return nil
}

// ---------------------------------------------------------
// TODO: 请完成下面三个核心方法
// ---------------------------------------------------------

// LogLocalEvent 本地发生了一个事件 (如打印日志)
func (n *Node) LogLocalEvent(msg string) {
	n.mu.Lock()
	defer n.mu.Unlock()

	// TODO: 1. 逻辑时钟自增
	
	// 记录日志
	n.Logs = append(n.Logs, LogEntry{Clock: n.Clock, Msg: msg, Node: n.ID})
}

// SendMessage 给指定节点发送消息
func (n *Node) SendMessage(targetID string, msgContent string) error {
	n.mu.Lock()
	// 注意：RPC 调用可能耗时，通常不建议在持有锁时进行网络 IO。
	// 但为了简化保护 Clock 的逻辑，这里先持有锁获取 Time，解锁后再发送，
	// 或者全程持有锁（简单但性能低）。
	// 建议：获取当前 Clock 并自增，然后解锁，再去发 RPC。
	
	// TODO: 2. 逻辑时钟自增 (发送事件)
	
	client, ok := n.peers[targetID]
	n.mu.Unlock() // 解锁以进行网络请求

	if !ok {
		return fmt.Errorf("peer %s not found", targetID)
	}

	// TODO: 3. 构造 RPC 参数，args和reply


	// TODO: 4. 调用 RPC (使用 client.Call)
	// 方法名是 "Node.HandleMessage"
	
	return nil // 占位
}

// HandleMessage RPC 处理器：接收消息
// 注意：RPC 方法必须导出 (首字母大写)，且签名必须符合 (args, reply) error
func (n *Node) HandleMessage(args *LamportMsg, reply *Reply) error {
	n.mu.Lock()
	defer n.mu.Unlock()

	// TODO: 5. 实现 Lamport 接收规则

	// 记录一条日志表示收到了消息
	logMsg := fmt.Sprintf("Recv from %s: %s", args.SenderID, args.Payload)
	n.Logs = append(n.Logs, LogEntry{Clock: n.Clock, Msg: logMsg, Node: n.ID})

	return nil
}

// Helper: 获取最新时间 (用于测试)
func (n *Node) GetTime() int {
	n.mu.Lock()
	defer n.mu.Unlock()
	return n.Clock
}