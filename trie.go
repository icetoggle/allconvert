package main

import (
	"log"
	"strconv"
	"strings"
)

type Trie struct {
	root *TrieNode
	id   int
}

type TrieNode struct {
	key       string
	id        int
	childNode []*TrieNode
}

func (t *TrieNode) searchKey(key string) *TrieNode {
	for _, child := range t.childNode {
		if child.key == key {
			return child
		}
	}
	return nil
}

func newTrie() *Trie {
	root := new(TrieNode)
	root.key = ""
	root.childNode = make([]*TrieNode, 0)

	tree := new(Trie)
	tree.root = root
	tree.id = 0
	return tree
}

func (t *Trie) pushXmlNode(root *TrieNode, obj interface{}) {
	switch obj.(type) {
	case []*KeyValue:
		set := obj.([]*KeyValue)

		for _, child := range set {
			trieNode := root.searchKey(child.key)
			if trieNode == nil {
				trieNode = new(TrieNode)
				trieNode.key = child.key
				trieNode.childNode = make([]*TrieNode, 0)
				root.childNode = append(root.childNode, trieNode)

			}
			t.pushXmlNode(trieNode, child.value)
		}

	case []interface{}:
		list := obj.([]interface{})
		for i, child := range list {
			trieNode := root.searchKey(strconv.Itoa(i))
			if trieNode == nil {
				trieNode = new(TrieNode)
				trieNode.key = strconv.Itoa(i)
				trieNode.childNode = make([]*TrieNode, 0)
				root.childNode = append(root.childNode, trieNode)
			}
			t.pushXmlNode(trieNode, child)
		}

	}
}

func (t *Trie) recordId(root *TrieNode) {
	if len(root.childNode) == 0 {
		root.id = t.id
		t.id = t.id + 1
		return
	}

	for _, child := range root.childNode {
		t.recordId(child)
	}
}

func (t *Trie) findKey(key string) int {
	keys := strings.Split(key, ".")
	return t._findKey(t.root, keys, 0)
}

func (t *Trie) _findKey(root *TrieNode, keys []string, offset int) int {
	if len(root.childNode) == 0 {
		return root.id
	}
	key := root.searchKey(keys[offset])
	if key == nil {
		log.Fatalln("key can't be nil")
	}

	return t._findKey(key, keys, offset+1)

}

func (t *Trie) dfs() {
	dfsTrie(t.root)
}

func dfsTrie(root *TrieNode) {
	log.Println(root.key + "={")
	for _, child := range root.childNode {
		dfsTrie(child)
		log.Println(",")
	}
	log.Println("}")
}
