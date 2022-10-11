package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type wordCounter struct {
	Word  string
	Count int
}

func Top10(text string) []string {
	words := strings.Fields(text)
	wordCounterMap := make(map[string]int)
	wordsLimit := 10

	for _, word := range words {
		if _, ok := wordCounterMap[word]; !ok {
			wordCounterMap[word] = 0
		}

		wordCounterMap[word]++
	}

	counterSlice := make([]wordCounter, 0, len(wordCounterMap))

	for word, count := range wordCounterMap {
		counterSlice = append(counterSlice, wordCounter{word, count})
	}

	sort.SliceStable(counterSlice, func(i, j int) bool {
		return counterSlice[i].Count > counterSlice[j].Count ||
			(counterSlice[i].Count == counterSlice[j].Count && counterSlice[i].Word < counterSlice[j].Word)
	})

	if len(counterSlice) < 10 {
		wordsLimit = len(counterSlice)
	}

	result := make([]string, 0, wordsLimit)

	for _, wordCounter := range counterSlice[:wordsLimit] {
		result = append(result, wordCounter.Word)
	}

	return result
}
