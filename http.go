package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

func runHttpServer(addr string, config RunConfig) error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		f, err := os.Open("index.html")
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		defer f.Close()
		io.Copy(w, f)
	})
	http.HandleFunc("/data.json", func(w http.ResponseWriter, r *http.Request) {
		if err := dataToJson(w); err != nil {
			fmt.Println("Error:", err)
			return
		}
	})
	return http.ListenAndServe(addr, nil)
}

func dataToJson(w io.Writer) error {
	stats, err := loadStats("stats.json")
	if err != nil {
		return err
	}
	data, err := statsToNodes(stats)
	if err != nil {
		return err
	}
	if err := json.NewEncoder(w).Encode(data); err != nil {
		return fmt.Errorf("encoding: %w", err)
	}
	return nil
}

func loadStats(file string) ([]FileCommitStats, error) {
	var stats []FileCommitStats
	in, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return stats, json.Unmarshal(in, &stats)
}

func statsToNodes(stats []FileCommitStats) (*Node, error) {
	root := &Node{
		Name: "root",
	}
	for _, stat := range stats {
		if err := root.AddFileStats(stat); err != nil {
			return nil, err
		}
	}
	return root, nil
}
