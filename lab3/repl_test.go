package lab3

import (
	"fmt"
	"testing"
	"time"
)

func TestReplication(t *testing.T) {
	// 1. 启动两个节点
	primary := NewNode("Primary", true)
	backup := NewNode("Backup", false)
	defer primary.Close()
	defer backup.Close()

	primary.StartRPCServer()
	bAddr := backup.StartRPCServer()

	// 2. 建立连接 (Primary -> Backup)
	if err := primary.ConnectToPeer(bAddr); err != nil {
		t.Fatal(err)
	}

	fmt.Println("🚀 集群启动，开始测试...")

	// --- 测试 1: 单条写入 ---
	if err := primary.AppendClient("CMD1"); err != nil {
		t.Fatalf("第一次写入失败: %v", err)
	}

	// 验证一致性
	pLog := primary.GetLog()
	bLog := backup.GetLog()

	if len(pLog) != 1 || len(bLog) != 1 {
		t.Fatalf("日志长度不一致: P=%d, B=%d", len(pLog), len(bLog))
	}
	if pLog[0].Command != "CMD1" || bLog[0].Command != "CMD1" {
		t.Fatal("日志内容错误")
	}
	fmt.Println("✅ 单条写入通过")

	// --- 测试 2: 连续写入 ---
	for i := 0; i < 5; i++ {
		cmd := fmt.Sprintf("CMD_BATCH_%d", i)
		if err := primary.AppendClient(cmd); err != nil {
			t.Fatalf("批量写入失败: %v", err)
		}
	}

	if len(backup.GetLog()) != 6 {
		t.Fatal("批量写入后长度错误")
	}
	fmt.Println("✅ 连续写入通过")
}

func TestReplicationFailure(t *testing.T) {
	// 测试原子性：如果 Backup 挂了，Primary 应该回滚，不能自己偷偷写入
	primary := NewNode("Primary", true)
	
	// 注意：我们故意不启动 Backup，或者连一个不存在的地址
	// 这样 Primary 的 Connect 会成功（因为 Dial 只是建立对象），但 Call 会失败
	// 或者我们建立连接后把 Backup 关掉
	
	backup := NewNode("Backup", false)
	bAddr := backup.StartRPCServer()
	primary.ConnectToPeer(bAddr)
	
	// 写入一条成功的数据
	primary.AppendClient("SafeCmd")

	// 💀 模拟 Backup 宕机
	backup.Close()
	time.Sleep(100 * time.Millisecond) // 等 TCP 断开

	fmt.Println("💀 模拟 Backup 宕机，尝试写入...")

	// 尝试写入新数据
	err := primary.AppendClient("UnsafeCmd")

	// 期望：写入失败
	if err == nil {
		t.Fatal("错误：Backup 宕机了，Primary 依然返回成功，违反了一致性！")
	}

	// 核心验证：Primary 的日志里不应该有 "UnsafeCmd"
	// 它应该在检测到 RPC 失败后，把本地已经 append 进去的那条删掉 (Rollback)
	logs := primary.GetLog()
	lastLog := logs[len(logs)-1]

	if lastLog.Command == "UnsafeCmd" {
		t.Fatalf("严重错误：Primary 保存了脏数据！没有回滚！")
	}
	if lastLog.Command != "SafeCmd" {
		t.Fatalf("数据被破坏")
	}

	fmt.Println("✅ 故障回滚测试通过 (Strong Consistency)")
}