package storage

import (
	"bytes"
	"encoding/csv"
	"os"
)

type Storage struct {
	n *node
}

func (s *Storage) Find(phone string) (code string) {
	return s.n.find(phone)
}

func (s *Storage) Update() {
	// при необходимости обновлять список кодов,
	// можно создавать новое дерево и заменять на него текущее при помощи atomic
}

type node struct {
	code     string
	children [10]*node
}

func (n *node) find(s string) string {
	if s == "" {
		return n.code
	}
	next := n.children[s[0]-'0']
	if next == nil {
		return n.code
	}
	return next.find(s[1:])
}

func (n *node) add(prefix string, code string) {
	if prefix == "" {
		if n.code != "" {
			return
		}
		n.code = code
		return
	}
	i := prefix[0] - '0'
	if n.children[i] == nil {
		n.children[i] = &node{}
	}
	n.children[i].add(prefix[1:], code)
}

func New(file string) *Storage {
	data, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}
	records, err := csv.NewReader(bytes.NewReader(data)).ReadAll()
	if err != nil {
		panic(err)
	}
	records = records[1:]
	return &Storage{n: constructTree(records)}
}

func constructTree(rows [][]string) *node {
	root := &node{
		code:     "",
		children: [10]*node{},
	}
	for _, r := range rows {
		root.add(r[0], r[1])
	}
	return root
}
