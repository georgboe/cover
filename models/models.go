package models

type File struct {
	Filename string
	Basename string
	Testname string
	IsTest   bool
}

type FileAndTest struct {
	Sourcefile *File
	Testfile   *File
}
