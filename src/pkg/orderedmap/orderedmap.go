package orderedmap

import (
	"sync"
)

type keyValue struct {
	key   string
	value string
}

type listNode struct {
	key  string
	next *listNode
	prev *listNode
}

type OrderedMap struct {
	sync.RWMutex
	head   *listNode
	tail   *listNode
	values map[string]*keyValue
	nodes  map[string]*listNode // To keep track of list nodes for each key
}

func NewOrderedMap() *OrderedMap {
	return &OrderedMap{
		values: make(map[string]*keyValue),
		nodes:  make(map[string]*listNode),
	}
}

func (om *OrderedMap) Add(key, value string) {
	om.Lock()
	defer om.Unlock()

	if _, exists := om.values[key]; !exists {
		node := &listNode{key: key}
		if om.head == nil {
			om.head = node
			om.tail = node
		} else {
			om.tail.next = node
			node.prev = om.tail
			om.tail = node
		}
		om.nodes[key] = node
	}
	om.values[key] = &keyValue{key, value}
}

func (om *OrderedMap) Delete(key string) {
	om.Lock()
	defer om.Unlock()

	if node, exists := om.nodes[key]; exists {
		delete(om.values, key)
		delete(om.nodes, key)
		if node.prev != nil {
			node.prev.next = node.next
		} else {
			om.head = node.next
		}
		if node.next != nil {
			node.next.prev = node.prev
		} else {
			om.tail = node.prev
		}
	}
}

func (om *OrderedMap) Get(key string) (string, bool) {
	om.RLock()
	defer om.RUnlock()

	kv, exists := om.values[key]
	if !exists {
		return "", false
	}
	return kv.value, true
}

func (om *OrderedMap) GetAll() []keyValue {
	om.RLock()
	defer om.RUnlock()

	var items []keyValue
	for node := om.head; node != nil; node = node.next {
		if kv, exists := om.values[node.key]; exists {
			items = append(items, *kv)
		}
	}
	return items
}
