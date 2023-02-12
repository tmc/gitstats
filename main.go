package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	flagRepoPath   = flag.String("repo", "", "path to the git repository")
	flagAppendMode = flag.Bool("append", false, "append to existing file")
	flagOutputFile = flag.String("output", "git_history.csv", "output file")
	flagSince      = flag.String("since", "", "git log `--since` argument")

	flagPerformanceProfile = flag.String("profile", "", "write performance profile to file")
	flagVerbose            = flag.Bool("v", false, "enable verbose logging")
)

func main() {
	flag.Parse()
	// if enabled, generate an fgprof profile
	if *flagPerformanceProfile != "" {
		defer profile(*flagPerformanceProfile)()
	}
	config := RunConfig{
		RepoPath:   *flagRepoPath,
		AppendMode: *flagAppendMode,
		OutputFile: *flagOutputFile,
		Since:      *flagSince,
		Verbose:    *flagVerbose,
	}
	if err := runGitStats(config); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
