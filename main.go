package main

import (
	"net/http"
	"time"
	"errors"
	"strings"
	"fmt"
)

const baseUrl = "https://en.wikipedia.org"
const nWorker = 5

func work(g *Graph, done chan bool) {
	fmt.Println("Worker started")
	for {
		select {
		case <-g.stop:
			fmt.Println("Worker ended")
			done <- true
			return
		default:
			// Fetch the url and parse it
			url, err := g.pop()
			if err != nil {
				time.Sleep(time.Second / 2)
			} else {
				response, err := http.Get(baseUrl + url)
				if err != nil {
					panic(err)
				}
				children := parse(response.Body)
				g.addChildren(url, children)
			}
		}
	}
}

func buildResult(end *Node) []string {

	r := []string{end.value}
	if len(end.parents) == 0 {
		return make([]string, 0)
	}

	node := end.parents[0]
	for {
		r = append([]string{node.value}, r...)
		if len(node.parents) == 0 {
			return r
		} else {
			node = node.parents[0]
		}
	}

	return r
}

func testInput(link string) error {

	if strings.ContainsAny(link, " /:.&?=") {
		return errors.New("Forbidden character found")
	}

	_, err := http.Get(baseUrl + "/wiki/" + link)
	if err != nil {
		return errors.New("Cannot GET the URL")
	}
	
	return nil
}

func run(from, to string) ([]string, error) {

	if err := testInput(from); err != nil {
		return make([]string, 0), err
	}
	if err := testInput(to); err != nil {
		return make([]string, 0), err
	}

	graph := newGraph()
	graph.add("/wiki/" + from, 0)
	graph.endLabel = "/wiki/" + to

	stop := make(chan bool)
	// Run workers
	for i := 0; i < nWorker; i++ {
		go work(graph, stop)
	}

	for i := 0; i < nWorker; i++ {
		<-stop
	}

	return buildResult(graph.endNode), nil
}

func main() {
	res, err := run("Taylor_Swift", "Doc_Gynéco")

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, v := range res {
		fmt.Println(v)
	}
}
