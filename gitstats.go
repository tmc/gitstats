package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	flagRepoPath   = flag.String("repo", "", "path to the git repository")
	flagAppendMode = flag.Bool("append", false, "append to existing file")
	flagOutputFile = flag.String("output", "git_history.csv", "output file")
)

func main() {
	flag.Parse()
	if err := run(*flagRepoPath, *flagAppendMode, *flagOutputFile); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

// run runs the gitstats command
func run(repoPath string, appendMode bool, outputFile string) error {
	cmd := exec.Command("git", "log", "--pretty=format:%H	%ad	%an	%ae	%s", "--numstat")
	if repoPath == "" {
		return fmt.Errorf("repo path is required")
	}
	repo := filepath.Base(repoPath)
	cmd.Dir = repoPath
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	var file *os.File
	// Open a file for writing the CSV
	if appendMode {
		file, err = os.OpenFile(outputFile, os.O_APPEND|os.O_WRONLY, 0600)
	} else {
		file, err = os.Create(outputFile)
	}
	if err != nil {
		return err
	}
	defer file.Close()
	cw := csv.NewWriter(file)

	// Write the header to the file (if not in append mode)
	if !appendMode {
		if err := cw.Write([]string{"Repo", "Hash", "Date", "Author Name", "Author Email", "Subject", "Filename", "Lines Added", "Lines Removed"}); err != nil {
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
