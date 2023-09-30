package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Println("failed", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	action := "head"
	separator := ","
	start := 0
	end := 0
	fset := flag.NewFlagSet("csv", flag.ContinueOnError)
	fset.StringVar(&action, "action", action, "what action shoud be performed")
	fset.StringVar(&separator, "separator", separator, "what separator use the file")
	fset.IntVar(&start, "slice-start", start, "where slicing should start")
	fset.IntVar(&end, "slice-end", end, "where slicing should end")
	err := fset.Parse(args)
	if err != nil {
		return fmt.Errorf("cannot read args %w", err)
	}

	if len(separator) != 1 {
		return fmt.Errorf("separator should only by one char")
	}

	csvReader := csv.NewReader(os.Stdin)
	csvReader.Comma = []rune(separator)[0]

	switch action {
	case "head":
		line, err := csvReader.Read()
		if err != nil {
			return fmt.Errorf("cannot read head %w", err)
		}
		fmt.Println(line)
	case "check":
		_, err := csvReader.ReadAll()
		if err != nil {
			return fmt.Errorf("cannot read all %w", err)
		}
	case "slice":
		lines, err := csvReader.ReadAll()
		if err != nil {
			return fmt.Errorf("cannot read all %w", err)
		}
		for end <= 0 {
			end += len(lines)
		}

		if start >= len(lines) || end > len(lines) {
			return fmt.Errorf("slice start and end must be lower than file size %d", len(lines))
		}

		if start >= end {
			return fmt.Errorf("slice start must be before slice end")
		}

		w := csv.NewWriter(os.Stdout)
		return w.WriteAll(lines[start:end])
	}
	return nil
}
