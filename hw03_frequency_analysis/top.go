package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

type topWord struct {
	word string
	qty  int
}

var re = regexp.MustCompile(`\s+`)

func Top10(text string) []string {
	words := re.Split(text, -1)
	uniqWordsMap := make(map[string]*topWord)
	result := make([]string, 0, 10)
	maxItem := 10

	if strings.TrimSpace(text) == "" {
		return result
	}

	for _, word := range words {
		word = strings.TrimSpace(word)

		if _, ok := uniqWordsMap[word]; !ok {
			uniqWordsMap[word] = &topWord{word, 1}
			continue
		}

		uniqWordsMap[word].qty++
	}

	uniqWordsSlice := make([]*topWord, 0, len(uniqWordsMap))

	for _, value := range uniqWordsMap {
		uniqWordsSlice = append(uniqWordsSlice, value)
	}

	sort.Slice(uniqWordsSlice, func(i, j int) bool {
		if uniqWordsSlice[i].qty != uniqWordsSlice[j].qty {
			return uniqWordsSlice[i].qty > uniqWordsSlice[j].qty
		}

		return uniqWordsSlice[i].word < uniqWordsSlice[j].word
	})

	if len(uniqWordsSlice) < 10 {
		maxItem = len(uniqWordsSlice)
	}

	for i := 0; i < maxItem; i++ {
		result = append(result, uniqWordsSlice[i].word)
	}

	return result
}
