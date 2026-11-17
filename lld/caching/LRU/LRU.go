package LRU

import "fmt"

type Node struct {
	key   string
	value string

	prev *Node
	next *Node
}

type LRUCache struct {
	capacity int
	cache    map[string]*Node

	head *Node // most recently used
	tail *Node // least recently used
}

// Constructor
func NewLRUCache(capacity int) *LRUCache {
	if capacity < 0 {
		capacity = 0
	}

	return &LRUCache{
		capacity: capacity,
		cache:    make(map[string]*Node),
	}
}

// GET - O(1) average
func (l *LRUCache) Get(key string) (string, bool) {
	node, exists := l.cache[key]
	if !exists {
		return "", false
	}

	// Move accessed node to head
	l.removeNode(node)
	l.addToHead(node)

	return node.value, true
}

// PUT - O(1) average
func (l *LRUCache) Put(key, value string) {

	// Case 1: Key already exists
	if node, exists := l.cache[key]; exists {
		node.value = value

		l.removeNode(node)
		l.addToHead(node)
		return
	}

	// Case 2: New key
	node := &Node{
		key:   key,
		value: value,
	}

	l.addToHead(node)
	l.cache[key] = node

	// Evict least recently used node
	if len(l.cache) > l.capacity {
		lru := l.tail
		l.removeNode(lru)
		delete(l.cache, lru.key)
	}
}

// Remove a node from doubly linked list
func (l *LRUCache) removeNode(node *Node) {

	if node.prev != nil {
		node.prev.next = node.next
	} else {
		l.head = node.next
	}

	if node.next != nil {
		node.next.prev = node.prev
	} else {
		l.tail = node.prev
	}

	node.prev = nil
	node.next = nil
}

// Add node to head (most recently used)
func (l *LRUCache) addToHead(node *Node) {

	node.prev = nil
	node.next = l.head

	if l.head != nil {
		l.head.prev = node
	}

	l.head = node

	// If list was empty
	if l.tail == nil {
		l.tail = node
	}
}

// Print cache from HEAD to TAIL
func (l *LRUCache) Print() {
	current := l.head

	for current != nil {
		fmt.Printf("[%s:%s]", current.key, current.value)

		if current.next != nil {
			fmt.Print(" <-> ")
		}

		current = current.next
	}
	fmt.Println()
}

func main() {
	lru := NewLRUCache(3)

	lru.Put("A", "10")
	lru.Put("B", "20")
	lru.Put("C", "30")
	lru.Print()

	value, found := lru.Get("A")
	fmt.Println("Get A:", value, found)
	lru.Print()

	lru.Put("D", "40")
	lru.Print()

	lru.Put("C", "300")
	lru.Print()

	value, found = lru.Get("B")
	fmt.Println("Get B:", value, found)
}
