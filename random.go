package main

import (
	"math/rand"
	"time"
)

const (
	size int = 50
)

type node struct {
	name string
	min  int
	p    int
	next *node
}
type member struct {
	name  string
	p     int // 本次中奖概率
	email string
}

// randAsListNode 使用链表，实现抽奖功能，时间复杂度O(n^2)
func randAsListNode(members []*member) map[string]int {
	rand.Seed(time.Now().Unix())
	list := make(map[string]int, size)
	head := &node{name: ""}
	cur := head
	length, index := 0, 0
	for _, m := range members {
		newNode := &node{
			name: m.name,
			min:  index,
			p:    m.p,
		}
		cur.next = newNode
		cur = newNode
		length += newNode.p
		index = newNode.min + newNode.p
	}
	i := 0
	for i < size {
		r := rand.Intn(length)
		flag := false // 本次循环是否命中
		// 执行次数未为len(members)
		for cur, pre := head.next, head; cur != nil; {

			if !flag { // 没有命中记录，需要继续判断
				// 判断是否在区间内
				if cur.min <= r && cur.min+cur.p > r {
					//命中，移除本节点
					length -= cur.p
					list[cur.name] = 1
					pre.next, cur = cur.next, cur.next
					flag = true

				} else {
					pre, cur = cur, cur.next
				}

			} else { //有命中记录，修改后续节点的min值
				//修改本节点的min值
				cur.min = pre.min + pre.p
				pre, cur = cur, cur.next
			}
		}
		i++
	}
	return list
}

//randAsArray 反面例子，太慢了
// TODO 二分，优化
func randAsArray(members []*member) map[string]int {
	rand.Seed(time.Now().Unix())
	pool := make([]string, 0, len(members)*100)
	for _, m := range members {
		for i := 0; i < m.p; i++ {
			pool = append(pool, m.name)
		}
	}
	i := 0
	length := len(pool)
	list := make(map[string]int, size)
	for i < size {

		r := rand.Intn(length)
		s := pool[r]
		if _, ok := list[s]; !ok {
			list[s] = 1
			i++
		}
	}
	return list
}
