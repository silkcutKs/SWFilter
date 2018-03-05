package filter

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	alphaWordNum int32 = 26
)

var TreeRoot struct {
	Node [alphaWordNum]*TrieTree
}

type TrieTree struct {
	char   int
	is_end bool
	Node   []*TrieTree
}

// just match [a-z], for tried tree is easy
type TrieFilter struct {
	words           []string
	replaceWord     string
	config          FilterConfig
	conditionFilter *regexp.Regexp
	conditionMatch  *regexp.Regexp
}

func NewTrieFilter() Filter {
	return &TrieFilter{}
}

func (tf *TrieFilter) Conf(config FilterConfig) {
	tf.config = config
}

// build trie tree, now is alpha word
func (tf *TrieFilter) Build(words []string) error {
	tf.conditionFilter, _ = regexp.Compile("[^a-z]+")
	tf.conditionMatch, _ = regexp.Compile("[a-z]+")
	for _, word := range words {
		filterWord := tf.conditionFilter.ReplaceAllString(strings.ToLower(word), "")
		chars := []rune(filterWord)
		charsLen := len(chars)
		idx := 0
		tmpNode := &TreeRoot.Node[int(chars[idx])-int('a')]
		idx += 1
		preNode := tmpNode
		fmt.Println(word, tmpNode)
		for ; idx < charsLen; idx++ {
			indent := int(chars[idx]) - int('a')
			if *tmpNode == nil {
				*tmpNode = &TrieTree{
					char:   indent,
					is_end: false,
					Node:   make([]*TrieTree, alphaWordNum),
				}
			}

			preNode = tmpNode
			tmpNode = &(*tmpNode).Node[indent]
		}

		(*preNode).is_end = true
	}
	return nil
}

func (tf *TrieFilter) Search(text string) ([]FilterSet, error) {
	return nil, nil
}

func (tf *TrieFilter) Replace() error {
	return nil
}

func (tf *TrieFilter) Ban(text string) bool {
	words := tf.conditionMatch.FindAllString(text, -1)
	wordLen := len(words)
	match := make(chan bool, wordLen)
	matchRes := false
	for _, word := range words {
		go func(word string) {
			match <- tf.searchOne(word)
		}(word)
	}

	for i := 0; i < wordLen; i++ {
		res := <-match
		if res {
			matchRes = true
		}
	}
	return matchRes
}

// no need to lock
func (tf *TrieFilter) searchOne(word string) bool {
	chars := []rune(word)
	charsLen := len(chars)
	i := 0
	tmpNode := TreeRoot.Node[int(chars[i])-int('a')]
	for i = 1; i < charsLen; i++ {
		if tmpNode == nil {
			return false
		}
		tmpNode = tmpNode.Node[int(chars[i])-int('a')]
	}

	if i == charsLen && tmpNode == nil {
		return true
	}

	if tmpNode.is_end == true {
		return true
	} else {
		return false
	}
}
