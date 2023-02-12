package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// RunConfig is the configuration for running gitstats
type RunConfig struct {
	RepoPath   string
	OutputFile string
	AppendMode bool
	Since      string
	Verbose    bool
}

// runGitStats runs the gitstats command
func runGitStats(config RunConfig) error {
	gitArgs := []string{"log", "--pretty=format:%H	%ai	%an	%ae	%s", "--numstat"}
	if config.Since != "" {
		gitArgs = append(gitArgs, fmt.Sprintf("--since=%s", config.Since))
	}
	if config.Verbose {
		fmt.Printf("Running:\n\tgit %s\n", strings.Join(gitArgs, " "))
	}
	cmd := exec.Command("git", gitArgs...)
	if config.RepoPath == "" {
		return fmt.Errorf("repo path is required")
	}
	repo := filepath.Base(config.RepoPath)
	cmd.Dir = config.RepoPath
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return err
	}

	var file *os.File
	// Open a file for writing the CSV
	if config.AppendMode {
		file, err = os.OpenFile(config.OutputFile, os.O_APPEND|os.O_WRONLY, 0600)
	} else {
		file, err = os.Create(config.OutputFile)
	}
	if err != nil {
		return err
	}
	defer file.Close()
	cw := csv.NewWriter(file)

	// Write the header to the file (if not in append mode)
	if !config.AppendMode {
		if err := cw.Write([]string{"repo", "sha", "date", "author name", "author email", "subject", "filename", "lines_added", "lines_removed"}); err != nil {
			return err
		}
	}

	scanner := bufio.NewScanner(stdout)
	// Write the commit history to the file
	i := 0
	var hash, date, authorName, authorEmail, subject string
	for scanner.Scan() {
		i++
		line := scanner.Text()
		if line == "" {
			continue
		}
		parts := strings.Split(line, "\t")
		newCommit := len(parts) >= 5
		if newCommit {
			hash, date, authorName, authorEmail, subject = parts[0], parts[1], parts[2], parts[3], parts[4]
		} else {
			if err := cw.Write([]string{repo, hash, date, authorName, authorEmail, subject, parts[2], parts[0], parts[1]}); err != nil {
				return err
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		return err
	}
	cw.Flush()
	return cw.Error()
}

// type Nodes map[string]Node

type Node struct {
	Name     string  `json:"name,omitempty"`
	Value    float64 `json:"value,omitempty"`
	Children []*Node `json:"children,omitempty"`
}

type FileCommitStats struct {
	AllTimeRank   float64 `json:"all_time_rank,omitempty"`
	Filename      string  `json:"filename,omitempty"`
	RecentCommits string  `json:"recent_commits,omitempty"`
	RecentRank    float64 `json:"recent_rank,omitempty"`
	TotalCommits  float64 `json:"total_commits,omitempty"`
}

// func (n Nodes) AddFileStats(stats FileCommitStats) error {
// 	dir := filepath.Dir(stats.Filename)
// 	en, ok := n[dir]
// 	if !ok {
// 		en = Node{
// 			Children: Nodes{},
// 		}
// 		n[dir] = en
// 	}
// 	return en.AddFileStats(stats)
// }

/*
func (n Node) AddFileStats(stats FileCommitStats) error {
	dir := filepath.Dir(stats.Filename)
	if dir == "." || dir == "/" {
		n.Children[stats.Filename] = Node{
			Name:     stats.Filename,
			Value:    stats.TotalCommits,
			Children: Nodes{},
		}
		return nil
	}
	return n.Children.AddFileStats(stats)
}
*/

func (n *Node) AddFileStats(stats FileCommitStats) error {
	pathParts := strings.Split(stats.Filename, "/")
	isLeaf := len(pathParts) == 1

	if isLeaf {
		n.Children = append(n.Children, &Node{
			Name:  stats.Filename,
			Value: stats.TotalCommits,
		})
		return nil
	} else {
		dir := pathParts[0]
		stats.Filename = strings.Join(pathParts[1:], "/")

		var child *Node
		var found bool
		for i, c := range n.Children {
			if c.Name == dir {
				child = n.Children[i]
				found = true
			}
		}
		if !found {
			child = &Node{
				Name:  dir,
				Value: 0,
			}
			n.Children = append(n.Children, child)
		}
		return child.AddFileStats(stats)
	}
}
