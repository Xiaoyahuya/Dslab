package lab2

import (
	"fmt"
	"testing"
	"time"
)

func TestDistributedCausality(t *testing.T) {
	// 1. 创建三个节点
	n1 := NewNode("N1")
	n2 := NewNode("N2")
	n3 := NewNode("N3")

	// 2. 启动 RPC Server
	n1.StartRPCServer()
	addr2 := n2.StartRPCServer()
	addr3 := n3.StartRPCServer()

	// 3. 建立全连接网络 (N1 <-> N2 <-> N3)
	// 简单起见，我们只需要 N1->N2->N3 的单向链路即可测试
	if err := n1.ConnectToPeer("N2", addr2); err != nil { t.Fatal(err) }
	if err := n2.ConnectToPeer("N3", addr3); err != nil { t.Fatal(err) }

	fmt.Println("🚀 网络已建立，开始因果测试...")

	// --- 场景开始 ---

	// Step 1: N1 本地发生事件 A
	// 预期 N1.Clock = 1
	n1.LogLocalEvent("Event A")
	time.Sleep(10 * time.Millisecond) // 等待一下

	// Step 2: N1 发消息给 N2
	// 预期 N1 发送时 Clock = 2
	// N2 收到后 Clock 应该变成 max(0, 2) + 1 = 3
	err := n1.SendMessage("N2", "Hello N2")
	if err != nil {
		t.Fatalf("N1 发送失败: %v", err)
	}
	time.Sleep(50 * time.Millisecond) // 确保 RPC 到达

	// Step 3: N2 收到消息后，发消息给 N3
	// 预期 N2 发送时 Clock = 4
	// N3 收到后 Clock 应该变成 max(0, 4) + 1 = 5
	err = n2.SendMessage("N3", "Forward to N3")
	if err != nil {
		t.Fatalf("N2 发送失败: %v", err)
	}
	time.Sleep(50 * time.Millisecond)

	// --- 验证结果 ---

	t1 := n1.GetTime()
	t2 := n2.GetTime()
	t3 := n3.GetTime()

	fmt.Printf("最终时钟状态: N1=%d, N2=%d, N3=%d\n", t1, t2, t3)

	// 核心断言 1: 只要代码没写，N2 的时间戳肯定是 0 或 1 (仅本地自增)，一定小于 N1 的发送时间
	if t2 <= t1 {
		t.Fatalf("❌ 因果违反! N2 收到 N1 消息后，时间戳(%d) 应该大于 N1(%d)。\n原因：你可能没有实现 max(local, msgTime) + 1", t2, t1)
	}

	// 核心断言 2: 传递性 N1 -> N2 -> N3
	if t3 <= t2 {
		t.Fatalf("❌ 因果违反! N3 收到 N2 消息后，时间戳(%d) 应该大于 N2(%d)", t3, t2)
	}

	fmt.Println("✅ 通过：时钟严格递增 (N1 < N2 < N3)")
}

func TestConcurrentEvents(t *testing.T) {
	// 测试并发情况：两个节点互不通信，时间戳应该可能较小
	n_a := NewNode("A")
	n_b := NewNode("B")
	
	n_a.LogLocalEvent("A event") // Time=1
	n_b.LogLocalEvent("B event") // Time=1

	// 因为没有交互，它们的时间戳应该独立增长
	if n_a.GetTime() != 1 || n_b.GetTime() != 1 {
		t.Fatalf("初始状态错误")
	}
	
	fmt.Println("✅ 通过：并发事件独立性测试")
}