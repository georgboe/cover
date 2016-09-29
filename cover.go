package main

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strings"

	"github.com/olekukonko/tablewriter"
)

type file struct {
	filename string
	sortname string
	isTest   bool
	testFile *file
	keep     bool
}

var (
	buffer = make([]byte, 255, 255)
)

func getFiles(reader io.Reader) []file {
	scanner := bufio.NewScanner(reader)
	scanner.Buffer(buffer, 255)
	var files []file
	var path string
	for scanner.Scan() {
		path = scanner.Text()
		if len(path) == 0 {
			continue
		}
		f := file{
			filename: filepath.Base(path),
		}
		basename := strings.ToLower(strings.TrimSuffix(f.filename, filepath.Ext(f.filename)))
		f.isTest = strings.HasSuffix(basename, "test")
		f.sortname = strings.ToLower(strings.Replace(f.filename, "Test.", ".", 1))
		files = append(files, f)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	return files
}

func createPairs(files []file) {
	numFiles := len(files)
	var skip bool
	var f *file
	for i := 0; i < numFiles; i++ {
		f = &files[i]

		if skip {
			skip = false
			continue
		}

		var otherFile *file
		if i+1 < numFiles {
			otherFile = &files[i+1]
		}

		if !f.isTest && otherFile != nil && f.sortname == otherFile.sortname {
			f.testFile = otherFile
			f.keep = true
			skip = true
		} else if f.isTest && otherFile != nil && !otherFile.isTest && f.sortname == otherFile.sortname {
			otherFile.testFile = f
			otherFile.keep = true
		} else {
			f.keep = true
		}
	}
}

func getData(reader io.Reader) []file {
	files := getFiles(reader)
	qsort(files)
	createPairs(files)
	return files
}

func qsort(a []file) {
	if len(a) < 2 {
		return
	}

	left, right := 0, len(a)-1

	// Pick a pivot
	pivotIndex := rand.Int() % len(a)

	// Move the pivot to the right
	a[pivotIndex], a[right] = a[right], a[pivotIndex]

	// Pile elements smaller than the pivot on the left
	for i := range a {
		if strings.Compare(a[i].sortname, a[right].sortname) == -1 {
			a[i], a[left] = a[left], a[i]
			left++
		}
	}

	// Place the pivot after the last smaller element
	a[left], a[right] = a[right], a[left]

	// Go down the rabbit hole
	qsort(a[:left])
	qsort(a[left+1:])
}

func main() {
	pairs := getData(os.Stdin)

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Source", "Test"})

	for _, v := range pairs {

		if !v.keep {
			continue
		}

		var sourceFile string
		var testFile string

		if !v.isTest {
			sourceFile = v.filename
			if v.testFile != nil {
				testFile = v.testFile.filename
			}
		} else {
			testFile = v.filename
		}

		table.Append([]string{sourceFile, testFile})
	}

	table.Render()
}
