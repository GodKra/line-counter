# line-counter
A simple program to count lines of a file/project. It only counts the files which are supported (Shown below). Directories starting with "." (like ".vscode", ".git") will not be opened (read). If you want to print the files which are being read, you need to run the program as ```line-counter --v <Path>```. '--v' should be before Path because Go's flag library doesn't allow me to do it any other way (as far as I know.).

## Currently Supported Types
* Golang (.go)
* Rust (.rs)
* Kotlin (.kt)
* Java (.java)
* Markdown (.md)
* C (.c)

More types will be added later.
