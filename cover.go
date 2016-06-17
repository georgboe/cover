package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime/debug"
	"sort"
	"strings"

	"github.com/olekukonko/tablewriter"
)

type file struct {
	filename string
	basename string
	testname string
	isTest   bool
}

type fileAndTest struct {
	sourcefile *file
	testfile   *file
}

type byBasename []file

func (s byBasename) Len() int      { return len(s) }
func (s byBasename) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s byBasename) Less(i, j int) bool {
	return strings.Compare(s[i].basename, s[j].basename) == -1
}

func getFiles(reader io.Reader) []file {
	scanner := bufio.NewScanner(reader)
	files := make([]file, 0, 10)
	for scanner.Scan() {
		path := scanner.Text()
		if len(path) == 0 {
			continue
		}
		f := file{
			filename: filepath.Base(path),
		}
		f.basename = strings.ToLower(
			strings.TrimSuffix(f.filename, filepath.Ext(f.filename)))
		f.testname = f.basename + "test"
		f.isTest = strings.HasSuffix(f.basename, "test")
		files = append(files, f)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	return files
}

func getPairs(files []file) []fileAndTest {
	fileCount := len(files)
	pairs := make([]fileAndTest, 0, fileCount)
	fileList := make(map[string]*file, fileCount)
	for i := range files {
		f := &files[i]
		fileList[f.basename] = f
	}

	skippedFiles := make(map[*file]bool)
	for i := range files {
		f := &files[i]

		_, ok := skippedFiles[f]
		if ok {
			continue
		}

		testFile, ok := fileList[f.testname]

		if ok {
			pairs = append(pairs, fileAndTest{sourcefile: f, testfile: testFile})
			skippedFiles[testFile] = false
		} else if f.isTest {
			pairs = append(pairs, fileAndTest{testfile: f})
		} else {
			pairs = append(pairs, fileAndTest{sourcefile: f})
		}
	}
	return pairs
}

func getData(reader io.Reader) []fileAndTest {
	files := getFiles(reader)
	sort.Sort(byBasename(files))
	return getPairs(files)
}

func main() {
	debug.SetGCPercent(-1)
	pairs := getData(os.Stdin)

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Source", "Test"})

	for _, v := range pairs {

		var sourceFile string
		var testFile string

		if v.sourcefile != nil {
			sourceFile = v.sourcefile.filename
		}

		if v.testfile != nil {
			testFile = v.testfile.filename
		}

		table.Append([]string{sourceFile, testFile})
	}
	table.Render()
}
