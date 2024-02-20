package query

import (
	"github.com/hashicorp/go-set"
	"strings"
)

type Info struct {
	match   []string
	exclude []string
	one     []string
	hash    string
}

func NewInfo(match *set.Set[string], exclude *set.Set[string], one *set.Set[string]) *Info {
	var hashBuilder strings.Builder
	addWordsToHash(match, &hashBuilder)
	addWordsToHash(exclude, &hashBuilder)
	addWordsToHash(one, &hashBuilder)

	return &Info{
		match:   match.Slice(),
		exclude: exclude.Slice(),
		one:     one.Slice(),
		hash:    hashBuilder.String(),
	}
}

func addWordsToHash(words *set.Set[string], strBuilder *strings.Builder) {
	strBuilder.WriteString("/")
	words.ForEach(func(word string) bool {
		strBuilder.WriteString(word)
		strBuilder.WriteString(".")
		return true
	})
}
