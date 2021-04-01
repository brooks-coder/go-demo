package common

import (
	"container/list"
	"errors"
)

type LRUCache struct {
	max   int
	cache map[interface{}]*list.Element
	Call  func(key interface{}, value interface{})
	l     *list.List
}
type node struct {
	key   interface{}
	value interface{}
}

func NewLRUCache(max int, call func(key interface{}, value interface{})) *LRUCache {
	return &LRUCache{max: max, cache: make(map[interface{}]*list.Element), Call: call, l: list.New()}
}

func (c *LRUCache) Get(key interface{}) (interface{}, bool) {
	if c.cache == nil || c.l == nil {
		return nil, false
	}
	if ele, ok := c.cache[key]; ok {
		c.l.MoveToFront(ele)
		return ele.Value.(*node).value, true
	}
	return nil, false
}

func (c *LRUCache) Add(key interface{}, value interface{}) error {
	if c.cache == nil || c.l == nil {
		return errors.New("not init")
	}
	val := &node{
		key:   key,
		value: value,
	}
	if ele, ok := c.cache[key]; ok {
		ele.Value = val
		c.l.MoveToFront(ele)
		return nil
	}
	ele := c.l.PushFront(val)
	c.cache[key] = ele
	if c.max > 0 && c.l.Len() > c.max {
		c.removeOldest()
	}
	return nil
}

func (c *LRUCache) Remove(key interface{}) error {
	if c.cache == nil || c.l == nil {
		return errors.New("not init")
	}
	if ele, ok := c.cache[key]; ok {
		c.remove(ele)
		return nil
	}
	return nil
}

func (c *LRUCache) Len() int {
	if c.l == nil {
		return 0
	}
	return c.l.Len()
}

func (c *LRUCache) remove(ele *list.Element) {
	n := ele.Value.(*node)
	c.l.Remove(ele)
	delete(c.cache, n.key)
	if c.Call != nil {
		c.Call(n.key, n.value)
	}
}

func (c *LRUCache) removeOldest() {
	ele := c.l.Back()
	c.remove(ele)
}
