package msglist

import (
	"sync"
)

type node struct {
	msg  interface{}
	next *node
}
type MsgList struct {
	head *node
	tail *node

	count int
	// 条件变量
	cond *sync.Cond
	// 互斥锁
	l sync.Mutex
}

func NewMsgList() *MsgList {
	msgList := &MsgList{}
	msgList.cond = sync.NewCond(&msgList.l)
	return msgList
}

func (m *MsgList) Push(msg interface{}) {
	m.l.Lock()
	defer m.l.Unlock()

	n := &node{msg: msg}
	if m.head == nil {
		m.head = n
		m.tail = n
	} else {
		m.tail.next = n
		m.tail = n
	}
	m.count++
	m.cond.Signal()
}

func (m *MsgList) Pop() []interface{} {
	m.l.Lock()
	defer m.l.Unlock()

	for m.count == 0 {
		m.cond.Wait()
	}

	var msgs []interface{}
	for m.count > 0 {
		msgs = append(msgs, m.head.msg)
		m.head = m.head.next
		m.count--
	}
	return msgs
}

// Clear 清空消息列表
func (m *MsgList) Clear() {
	m.l.Lock()
	defer m.l.Unlock()

	m.head = nil
	m.tail = nil
	m.count = 0
}
