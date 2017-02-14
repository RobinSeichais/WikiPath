package main

import (
	"net/http"
	"time"
	"fmt"
)


const baseUrl = "https://en.wikipedia.org"

func main() {
	// run("Taylor_Swift", "Doc_Gynéco")
	run("Binary_heap", "Function_(mathematics)")
}

func monitor(queue *PriorityQueue, graph *Graph) {
	tick := time.Tick(10 * time.Second)
	for _ = range tick {
		fmt.Printf("[%s] %d queued items, %d nodes registered\n", time.Now().Format(time.UnixDate), queue.Len(), graph.Len())
	}
}

func run(start, target string) {

	start = "/wiki/" + start
	target = "/wiki/" + target

	var response *http.Response
	var err error

	stopWords := []string{"a", "about", "above", "after", "again", "against", "all", "am", "an", "and", "any", "are", "aren't", "as", "at", "be", "because", "been", "before", "being", "below", "between", "both", "but", "by", "can't", "cannot", "could", "couldn't", "did", "didn't", "do", "does", "doesn't", "doing", "don't", "down", "during", "each", "few", "for", "from", "further", "had", "hadn't", "has", "hasn't", "have", "haven't", "having", "he", "he'd", "he'll", "he's", "her", "here", "here's", "hers", "herself", "him", "himself", "his", "how", "how's", "i", "i'd", "i'll", "i'm", "i've", "if", "in", "into", "is", "isn't", "it", "it's", "its", "itself", "let's", "me", "more", "most", "mustn't", "my", "myself", "no", "nor", "not", "of", "off", "on", "once", "only", "or", "other", "ought", "our", "ours", "ourselves", "out", "over", "own", "same", "shan't", "she", "she'd", "she'll", "she's", "should", "shouldn't", "so", "some", "such", "than", "that", "that's", "the", "their", "theirs", "them", "themselves", "then", "there", "there's", "these", "they", "they'd", "they'll", "they're", "they've", "this", "those", "through", "to", "too", "under", "until", "up", "very", "was", "wasn't", "we", "we'd", "we'll", "we're", "we've", "were", "weren't", "what", "what's", "when", "when's", "where", "where's", "which", "while", "who", "who's", "whom", "why", "why's", "with", "won't", "would", "wouldn't", "you", "you'd", "you'll", "you're", "you've", "your", "yours", "yourself", "yourselves",}
	scorer := NewScorer()
	
	if response, err = http.Get(baseUrl + target); err != nil {
		panic(err)
	}
	scorer.Initialize(response.Body, stopWords)
	
	if response, err = http.Get(baseUrl + start); err != nil {
		panic(err)
	}
	scorer.Exclude(response.Body)

	// queueMut = new(sync.Mutex)
	queue := NewPriorityQueue()

	// graphMut = new(sync.Mutex)
	graph := NewGraph()
	maxDepth := -1

	queue.Push(NewItem(start, 0, 0))

	go monitor(queue, graph)

	for queue.Len() > 0 {

		page, depth := queue.Pop()
		
		if maxDepth >= 0 && depth > maxDepth {
			continue // Ignore nodes too deep if target was already reached
		}

		if response, err = http.Get(baseUrl + page); err != nil {
			queue.Push(NewItem(page, depth, 0))
			fmt.Println(err)
			continue
		}
		
		score, children := Parse(response.Body, scorer)
		// fmt.Println(children)
		children = graph.Set(page, depth, children) // Only keep children not parsed yet	
		for _, child := range children {
			// fmt.Println(child)
			if child == target {
				fmt.Println("Found !")
				maxDepth = depth+1
				fmt.Println(graph.GetPath(target))
				return
			} else {
				queue.Push(NewItem(child, depth+1, score))
			}
		}
	}

	fmt.Println(graph.GetPath(target))
}
