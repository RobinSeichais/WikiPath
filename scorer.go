package main

import (
	"io"
	"bufio"
	"strings"
	"github.com/kennygrant/sanitize"
	"bytes"
)

func clean(s string) string {
	str := sanitize.HTML(s)
	str = strings.Replace(str, ".", " ", -1)
	str = strings.Replace(str, ":", " ", -1)
	str = strings.Replace(str, ";", " ", -1)
	str = strings.Replace(str, ",", " ", -1)
	str = strings.Replace(str, "/", " ", -1)
	str = strings.Replace(str, "\\", " ", -1)
	str = strings.Replace(str, "?", " ", -1)
	str = strings.Replace(str, "!", " ", -1)
	str = strings.Replace(str, "'", " ", -1)
	str = strings.Replace(str, "\"", " ", -1)
	str = strings.Replace(str, "-", " ", -1)
	str = strings.Replace(str, "_", " ", -1)
	str = strings.Replace(str, "[", " ", -1)
	str = strings.Replace(str, "]", " ", -1)
	str = strings.Replace(str, "{", " ", -1)
	str = strings.Replace(str, "}", " ", -1)
	str = strings.Replace(str, "(", " ", -1)
	str = strings.Replace(str, ")", " ", -1)
	return str
}

type Scorer map[string]bool

func (s Scorer) Initialize(src io.Reader, stopWords []string) {

	stopWordsIndex := make(map[string]bool)
	for _, word := range stopWords {
		stopWordsIndex[word] = true
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(src)
	str := clean(buf.String())

	scanner := bufio.NewScanner(strings.NewReader(str))
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		word := strings.ToLower(scanner.Text())
		_, isStopWord := stopWordsIndex[word]
		if !isStopWord {
			s[word] = true
		}
	}
}

func (s Scorer) Exclude(src io.Reader) {

	buf := new(bytes.Buffer)
	buf.ReadFrom(src)
	str := clean(buf.String())

	scanner := bufio.NewScanner(strings.NewReader(str))
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		word := strings.ToLower(scanner.Text())
		delete(s, word)
	}
}

func (s Scorer) GetScore(src io.Reader) float64 {

	var score, n_word float64 = 0, 0

	buf := new(bytes.Buffer)
	buf.ReadFrom(src)
	str := clean(buf.String())
	
	scanner := bufio.NewScanner(strings.NewReader(str))
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		word := scanner.Text()
		n_word++
		_, keep := s[word]
		if keep {
			score++
		}
	}

	return score/n_word
}

func (s Scorer) String(n int) string {
	str := make([]string, 0)
	i := 0
	for key := range s {
		str = append(str, key)
		if i++; i == n {
			break
		}
	}
	return strings.Join(str, ", ")
}

func NewScorer() Scorer {
	return make(Scorer, 0)
}