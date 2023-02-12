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

	flagHttp = flag.String("http", "", "start http server on given address")
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

	runMode := "git2csv"
	if *flagHttp != "" {
		runMode = "http"
	}

	switch runMode {
	case "git2csv":
		if err := runGitStats(config); err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
	case "http":
		if err := runHttpServer(*flagHttp, config); err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
		}
	default:
		fmt.Fprintln(os.Stderr, "unknown run mode:", runMode)
	}
}
