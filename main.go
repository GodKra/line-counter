package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var verbose = flag.Bool("v", false, "Prints the files scanned. Usage: -v <true/false>")

// All types that are currently supported by the program
type supportedType string

const (
	goFiles  supportedType = "Go"
	rust     supportedType = "Rust"
	comment  supportedType = "Comments"
	markdown supportedType = "Markdown"
	kotlin   supportedType = "Kotlin"
	java     supportedType = "Java"
	c        supportedType = "C"
)

// counter contains information of the line count and number of files counted
type counter struct {
	lineCounts map[supportedType]int
	fileCount  int
}

// fileAdd Adds the number of lines and comments in a file to the counter. If isLang
// (Indicates whether it is a language or not) is false, number of comments will be added
// to the line count It also prints the file using printFiles()
func (c *counter) fileAdd(isLang bool, fileType supportedType, file *os.File) error {
	c.fileCount++
	lines, comments, e := lineCount(file)
	if e != nil {
		return e
	}
	printFiles(fileType, file)
	if isLang {
		c.lineCounts[fileType] = lines
		c.lineCounts[comment] = comments
	} else {
		c.lineCounts[fileType] = lines + comments
	}
	return nil
}

// add Adds the contents of the counter with another counter
func (c *counter) add(countr counter) {
	for k := range countr.lineCounts {
		c.lineCounts[k] += countr.lineCounts[k]
	}
	c.fileCount += countr.fileCount
}

// total Adds up all the lines counted and returns it.
func (c counter) total() int {
	var total int
	for _, f := range c.lineCounts {
		total += f
	}
	return total
}

// getElems Returns a []string which contains everything in the map as "Key: Value", where key is the file type
// and value is line count of that file.
func (c counter) getElems() []string {
	slice := make([]string, len(c.lineCounts))
	i := 0
	for k, v := range c.lineCounts {
		slice[i] = fmt.Sprintf("%v: %v", k, v)
		i++
	}
	return slice
}

func main() {
	flag.Parse()
	var path string
	if len(flag.Args()) > 0 {
		path = flag.Arg(0)
	} else {
		log.Fatal("Please specify a path")
	}
	file, e := os.Open(path)
	checkError(e)
	countr, e := count(file)
	checkError(e)
	fmt.Printf("Files Counted: %v\n", countr.fileCount)
	for _, s := range countr.getElems() {
		fmt.Println(s)
	}
	fmt.Printf("Total: %v\n", countr.total())
}

// Reads the contents of a single file and returns line count (without empty new lines and comments) and
// number of comments.
func lineCount(file *os.File) (lines int, comment int, e error) {
	s := bufio.NewReader(file)
	insideBlockComment := false
	var err error
	for {
		if err == io.EOF {
			return lines, comment, nil
		}
		b, e := s.ReadBytes('\n')
		err = e
		if e != nil && e != io.EOF {
			return 0, 0, e
		}

		str := bytes.TrimSpace(b)

		if bytes.HasPrefix(str, []byte{'/', '*'}) {
			insideBlockComment = true
			continue
		}
		if insideBlockComment {
			if bytes.HasPrefix(str, []byte{'*', '/'}) {
				insideBlockComment = false
				comment++
				continue
			}
			continue
		}

		if len(str) == 0 {
			continue
		}
		if str[0] == '/' {
			comment++
		} else {
			lines++
		}
	}
}

// count recursively reads a directory (excluding directories starting with '.') and whenever it
// encounters a file (with extension and if it is supported), it counts the lines of that file
// using lineCount() and returns a counter.
func count(file *os.File) (counter, error) {
	var countr counter
	countr.lineCounts = make(map[supportedType]int)
	stat, e := file.Stat()
	if e != nil {
		return counter{}, e
	}
	if strings.HasPrefix(stat.Name(), ".") {
		return counter{}, nil
	}

	if stat.IsDir() {
		files, e := file.Readdirnames(0)
		if e != nil {
			return counter{}, e
		}
		for _, filename := range files {
			f, e := os.Open(filepath.Join(file.Name(), filename))
			if e != nil {
				return counter{}, e
			}
			recursiveCountr, e := count(f)
			if e != nil {
				return counter{}, e
			}
			countr.add(recursiveCountr)
		}
		return countr, nil
	}
	switch filepath.Ext(file.Name()) {
	case ".go":
		e = countr.fileAdd(true, goFiles, file)
	case ".md":
		e = countr.fileAdd(false, markdown, file)
	case ".rs":
		e = countr.fileAdd(true, rust, file)
	case ".kt":
		e = countr.fileAdd(true, kotlin, file)
	case ".java":
		e = countr.fileAdd(true, java, file)
	case ".c":
		e = countr.fileAdd(true, c, file)
	}

	if e != nil {
		return counter{}, e
	}
	return countr, nil
}

// Prints the current file being read as "<FileType>: <FilePath>" if verbose is true.
func printFiles(fileType supportedType, file *os.File) {
	if *verbose {
		fmt.Printf("%v: %v\n", fileType, file.Name())
	}
}

func checkError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
