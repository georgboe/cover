package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/olekukonko/tablewriter"
)

type file struct {
	path     string
	filename string
	basename string
}

type fileAndTest struct {
	sourcefile file
	testfile   file
}

type byBasename []file

func (s byBasename) Len() int      { return len(s) }
func (s byBasename) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s byBasename) Less(i, j int) bool {
	return strings.Compare(s[i].basename, s[j].basename) == -1
}

func getFiles() []file {
	scanner := bufio.NewScanner(os.Stdin)
	var files []file
	for scanner.Scan() {
		path := scanner.Text()
		if len(path) == 0 {
			continue
		}
		f := file{
			path:     path,
			filename: filepath.Base(path),
		}
		f.basename = strings.ToLower(
			strings.TrimSuffix(f.filename, filepath.Ext(f.filename)))
		files = append(files, f)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	return files
}

func getPairs(files []file) []fileAndTest {
	var pairs []fileAndTest
	fileList := make(map[string]file)
	for _, f := range files {
		fileList[f.basename] = f
	}

	skippedFiles := make(map[string]interface{})
	var ok bool
	for _, f := range files {

		_, ok = skippedFiles[f.path]
		if ok {
			continue
		}

		testFile, ok := fileList[f.basename+"test"]

		if ok {
			pairs = append(pairs, fileAndTest{sourcefile: f, testfile: testFile})
			skippedFiles[testFile.path] = nil
		} else if strings.HasSuffix(f.basename, "test") {
			pairs = append(pairs, fileAndTest{testfile: f})
		} else {
			pairs = append(pairs, fileAndTest{sourcefile: f})
		}
	}
	return pairs
}

func main() {
	files := getFiles()
	sort.Sort(byBasename(files))
	pairs := getPairs(files)

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Source", "Test"})

	for _, v := range pairs {
		table.Append([]string{v.sourcefile.filename, v.testfile.filename})
	}
	table.Render()
}
