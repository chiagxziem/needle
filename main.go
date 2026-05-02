package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/chiagxziem/needle/internal/search"
)

func main() {
	// define flags
	ignoreCase := flag.Bool("i", false, "ignore case distinctions in patterns")
	showLineNumbers := flag.Bool("n", false, "print line number with output lines")
	printCountPerFIle := flag.Bool("c", false, "print only a count of matching lines per file")
	printFilesWithMatches := flag.Bool("l", false, "print only filenames with matches")
	useFixedStrings := flag.Bool("F", false, "use patterns as strings instead of regular expressions")
	recursiveSearch := flag.Bool("r", false, "search files & directories recursively")

	// parse the command line into the defined flags
	flag.Parse()

	// show usage & help message if no pattern is passed
	if len(flag.Args()) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: needle [OPTION]... PATTERNS [FILE]...")
		fmt.Fprintln(os.Stderr, "Try 'needle --help' for more information.")
		os.Exit(1)
	}

	// get pattern and paths, if given
	pattern, paths := flag.Arg(0), flag.Args()[1:]

	// define opts from flags
	opts := search.Options{
		IgnoreCase:            *ignoreCase,
		ShowLineNumbers:       *showLineNumbers,
		PrintCountPerFile:     *printCountPerFIle,
		PrintFilesWithMatches: *printFilesWithMatches,
		UseFixedStrings:       *useFixedStrings,
		RecursiveSearch:       *recursiveSearch,
	}

	hasAnyMatch := false

	// recursive mode
	if opts.RecursiveSearch {
		var roots []string
		if len(paths) == 0 {
			roots = append(roots, ".")
		} else {
			roots = paths
		}

		for _, root := range roots {
			results, err := search.SearchDir(root, pattern, opts)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			for _, result := range results {
				if result.HasMatch {
					hasAnyMatch = true
				}

				getOutput(result, opts, true)
			}
		}

		// if no file matches the pattern, exit the program with code 1
		if !hasAnyMatch {
			os.Exit(1)
		}
		return
	}

	// Stdin mode
	if len(paths) == 0 {
		hasMatch, err := search.SearchStdin(pattern, opts)

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if !hasMatch {
			os.Exit(1)
		}

		return
	}

	// file mode
	multipleFiles := len(paths) > 1

	for _, p := range paths {
		result, err := search.SearchFile(p, pattern, opts)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if result.HasMatch {
			hasAnyMatch = true
		}

		getOutput(result, opts, multipleFiles)
	}

	// if no file matches the pattern, exit the program with code 1
	if !hasAnyMatch {
		os.Exit(1)
	}

}

func getOutput(r search.Result, opts search.Options, multipleFiles bool) {
	if opts.PrintFilesWithMatches {
		if r.HasMatch {
			fmt.Println(r.Path)
		}
	} else if opts.PrintCountPerFile {
		if multipleFiles {
			fmt.Printf("%s:%d\n", r.Path, r.Count)
		} else {
			fmt.Println(r.Count)
		}
	} else {
		for _, m := range r.Matches {
			if multipleFiles {
				fmt.Printf("%s:%s\n", r.Path, m.Format(opts.ShowLineNumbers))
			} else {
				fmt.Println(m.Format(opts.ShowLineNumbers))
			}
		}
	}
}
