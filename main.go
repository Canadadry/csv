package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"sort"
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
	newSeparator := ";"
	start := 0
	end := 0
	size := 0
	lazy := false
	insertPos := 0
	insertName := ""
	insertValue := ""
	swapX := 0
	swapY := 0
	deleteIndex := 0
	mergeX := 0
	mergeY := 0
	mergeSperator := "-"
	sortIndex := 0
	fset := flag.NewFlagSet("csv", flag.ContinueOnError)
	fset.StringVar(&action, "action", action, "what action shoud be performed")
	fset.StringVar(&separator, "separator", separator, "what separator use the file")
	fset.StringVar(&newSeparator, "convert-to-separator", newSeparator, "what separator you want")
	fset.BoolVar(&lazy, "lazy", lazy, "should use lazy quoted")
	fset.IntVar(&start, "slice-start", start, "where slicing should start")
	fset.IntVar(&end, "slice-end", end, "where slicing should end")
	fset.IntVar(&size, "check-size", size, "what width must have the csv")
	fset.IntVar(&insertPos, "insert-pos", insertPos, "where to insert column")
	fset.StringVar(&insertName, "insert-name", insertName, "how to name the inserted column")
	fset.StringVar(&insertValue, "insert-value", insertValue, "default insterted value")
	fset.IntVar(&swapX, "swap-x", swapX, "col 1 to swap")
	fset.IntVar(&swapY, "swap-y", swapY, "col 2 to swap")
	fset.IntVar(&deleteIndex, "delete-index", deleteIndex, "col to delete")
	fset.IntVar(&mergeX, "merge-x", mergeX, "col 1 to merge")
	fset.IntVar(&mergeY, "merge-y", mergeY, "col 2 to merge")
	fset.StringVar(&mergeSperator, "merge-separator", mergeSperator, "join between merged columns")
	fset.IntVar(&sortIndex, "sort-by", sortIndex, "sort by column index")
	err := fset.Parse(args)
	if err != nil {
		return fmt.Errorf("cannot read args %w", err)
	}

	if len(separator) != 1 {
		return fmt.Errorf("separator should only by one char")
	}

	csvReader := csv.NewReader(os.Stdin)
	csvReader.Comma = []rune(separator)[0]
	csvReader.LazyQuotes = lazy
	if size > 0 {
		csvReader.FieldsPerRecord = size
	}

	switch action {
	case "head":
		line, err := csvReader.Read()
		if err != nil {
			return fmt.Errorf("cannot read head %w", err)
		}
		fmt.Printf("%d : %#v\n", len(line), line)
	case "size":
		line, err := csvReader.Read()
		if err != nil {
			return fmt.Errorf("cannot read head %w", err)
		}
		fmt.Println(len(line))
	case "check":
		_, err = csvReader.ReadAll()
		if err != nil {
			return fmt.Errorf("cannot read all %w", err)
		}
	case "convert":
		lines, err := csvReader.ReadAll()
		if err != nil {
			return fmt.Errorf("cannot read all %w", err)
		}
		w := csv.NewWriter(os.Stdout)
		w.Comma = []rune(newSeparator)[0]
		return w.WriteAll(lines)
	case "insert":
		lines, err := csvReader.ReadAll()
		if err != nil {
			return fmt.Errorf("cannot read all %w", err)
		}
		lines[0] = insert(lines[0], insertPos, insertName)
		for i := 1; i < len(lines); i++ {
			lines[i] = insert(lines[i], insertPos, insertValue)
		}
		w := csv.NewWriter(os.Stdout)
		w.Comma = []rune(newSeparator)[0]
		return w.WriteAll(lines)
	case "swap":
		lines, err := csvReader.ReadAll()
		if err != nil {
			return fmt.Errorf("cannot read all %w", err)
		}
		for i := 0; i < len(lines); i++ {
			lines[i][swapX], lines[i][swapY] = lines[i][swapY], lines[i][swapX]
		}
		w := csv.NewWriter(os.Stdout)
		w.Comma = []rune(newSeparator)[0]
		return w.WriteAll(lines)
	case "delete":
		lines, err := csvReader.ReadAll()
		if err != nil {
			return fmt.Errorf("cannot read all %w", err)
		}
		for i := 0; i < len(lines); i++ {
			lines[i] = remove(lines[i], deleteIndex)
		}
		w := csv.NewWriter(os.Stdout)
		w.Comma = []rune(newSeparator)[0]
		return w.WriteAll(lines)
	case "merge":
		lines, err := csvReader.ReadAll()
		if err != nil {
			return fmt.Errorf("cannot read all %w", err)
		}
		for i := 0; i < len(lines); i++ {
			lines[i][mergeX] = lines[i][mergeX] + mergeSperator + lines[i][mergeY]
		}
		w := csv.NewWriter(os.Stdout)
		w.Comma = []rune(newSeparator)[0]
		return w.WriteAll(lines)
	case "sort":
		lines, err := csvReader.ReadAll()
		if err != nil {
			return fmt.Errorf("cannot read all %w", err)
		}
		sort.Slice(lines, func(i, j int) bool { return lines[i][sortIndex] < lines[j][sortIndex] })
		w := csv.NewWriter(os.Stdout)
		w.Comma = []rune(newSeparator)[0]
		return w.WriteAll(lines)
	}
	return nil
}

func insert[T any](a []T, index int, value T) []T {
	if index >= len(a) {
		return append(a, value)
	}
	a = append(a[:index+1], a[index:]...)
	a[index] = value
	return a
}

func remove[T any](slice []T, s int) []T {
	return append(slice[:s], slice[s+1:]...)
}
