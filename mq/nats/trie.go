package nats

import (
	"errors"
	"strings"
	"sync"

	pubsub "github.com/zhangce1999/pubsub/interface"
)

// DefaultSeparator -
var (
	DefaultSeparator = '/'
)

// Trie -
type Trie interface {
	Empty() bool
	Size() int
	Clear()

	// Trie interface, middlewares will be executed in order
	Put(route string, handlers ...pubsub.HandlerFunc) error
	Get(route string) (pubsub.HandlersChain, error)
	Remove(route string)
}

// NewTrie -
func NewTrie(sep rune) Trie {
	if sep == 0 {
		sep = DefaultSeparator
	}

	return &trie{
		node: &node{
			word: "/",
		},
		size: 0,
		sep:  sep,
		mu:   new(sync.Mutex),
	}
}

type trie struct {
	node *node
	size int
	sep  rune // sep represents the separator
	mu   *sync.Mutex
}

type node struct {
	word     string
	parent   *node
	children map[string]*node
	handlers pubsub.HandlersChain
}

func (t *trie) Empty() bool {
	return len(t.node.children) == 0
}

func (t *trie) Size() int {
	return t.size
}

func (t *trie) Clear() {
	t.node = &node{}
	t.size = 0
	t.sep = 0
}

func (t *trie) Put(route string, handlers ...pubsub.HandlerFunc) error {
	query, err := splitWithSeparator(route, t.sep)
	if err != nil {
		return err
	}

	t.mu.Lock()
	curr := t.node
	for _, word := range query {
		child, ok := curr.children[word]
		if !ok {
			child = &node{
				word:     word,
				parent:   curr,
				children: make(map[string]*node),
				handlers: nil,
			}
			curr.children[word] = child
		}
		curr = child
	}

	// Add Connection
	curr.handlers = append(curr.handlers, handlers...)
	t.size++
	t.mu.Unlock()
	return nil
}

func (t *trie) Get(route string) (pubsub.HandlersChain, error) {
	query, err := splitWithSeparator(route, t.sep)
	if err != nil {
		return nil, err
	}

	var res pubsub.HandlersChain

	t.mu.Lock()
	t.get(query, &res, t.node)
	t.mu.Unlock()

	return res, nil
}

func (t *trie) get(query []string, res *pubsub.HandlersChain, node *node) {
	// If we're not yet done, continue
	if len(query) > 0 {
		// Go through the exact match node
		if n, ok := node.children[query[0]]; ok {
			t.get(query[1:], res, n)
		}
	}
}

func (t *trie) Remove(route string) {
	query, err := splitWithSeparator(route, t.sep)
	if err != nil {
		return
	}

	t.mu.Lock()
	curr := t.node
	for _, word := range query {
		child, ok := curr.children[word]
		if !ok {
			t.mu.Unlock()
			return
		}
		curr = child
	}

	// Remove orphans
	if len(curr.handlers) == 0 && len(curr.children) == 0 {
		curr.orphan()
	}

	t.mu.Unlock()
	return
}

func splitWithSeparator(route string, sep rune) ([]string, error) {
	if route == "" {
		return nil, errors.New("[error]: invalid route")
	}

	return strings.FieldsFunc(route, func(r rune) bool {
		return r == sep
	}), nil
}

func (n *node) orphan() {
	if n.parent == nil {
		return
	}

	delete(n.parent.children, n.word)

	if len(n.parent.handlers) == 0 && len(n.parent.children) == 0 {
		n.parent.orphan()
	}
}
