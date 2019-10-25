package trie

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
	Add(route string, middlewares ...pubsub.HandlerFunc) error
	Search(route string) (pubsub.HandlersChain, error)
	Delete(route string)
}

// NewTrie -
func NewTrie(separator rune) Trie {
	if separator == 0 {
		separator = DefaultSeparator
	}

	return &trie{
		node: &node{
			word: "/",
		},
		size:      0,
		separator: separator,
		mu:        new(sync.Mutex),
	}
}

type trie struct {
	node      *node
	size      int
	separator rune
	mu        *sync.Mutex
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
	t.separator = 0
}

func (t *trie) Add(route string, handlers ...pubsub.HandlerFunc) error {
	query, err := splitWithSeparator(route, t.separator)
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
			}
			curr.children[word] = child
		}
		curr = child
	}

	// Add the handlers
	curr.handlers = handlers
	t.mu.Unlock()
	return nil
}

func (t *trie) Search(route string) (pubsub.HandlersChain, error) {
	query, err := splitWithSeparator(route, t.separator)
	if err != nil {
		return nil, err
	}

	var res pubsub.HandlersChain

	t.mu.Lock()
	t.search(query, &res, t.node)
	t.mu.Unlock()

	return res, nil
}

func (t *trie) search(query []string, res *pubsub.HandlersChain, node *node) {
	// Add handlers from the current node
	if len(*res) > 0 {
		*res = append(*res, node.handlers...)
	}

	// If we're not yet done, continue
	if len(query) > 0 {
		// Go through the exact match node
		if n, ok := node.children[query[0]]; ok {
			t.search(query[1:], res, n)
		}
	}
}

func (t *trie) Delete(route string) {
	query, err := splitWithSeparator(route, t.separator)
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

func splitWithSeparator(route string, separator rune) ([]string, error) {
	if route == "" {
		return nil, errors.New("[error]: invalid route")
	}

	return strings.FieldsFunc(route, func(r rune) bool {
		return r == separator
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
