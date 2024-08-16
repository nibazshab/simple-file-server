package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
)

var (
	port *string
	root *string
)

func main() {
	port = flag.String("port", "8080", "server port")
	root = flag.String("path", "./", "server root path")
	flag.Parse()

	fmt.Println("port=" + *port + "\n" + "path=" + *root)

	http.HandleFunc("/", FileServer)
	http.ListenAndServe(":"+*port, nil)
}

func FileServer(w http.ResponseWriter, req *http.Request) {
	filePath := path.Join(*root, req.URL.Path)

	file, err := os.Open(filePath)
	if err != nil {
		http.NotFound(w, req)
		return
	}
	defer file.Close()

	f, _ := file.Stat()

	if f.IsDir() {
		files, _ := file.Readdir(-1)

		w.Write([]byte("<html><pre><a href=\"../\">../</a>" + "\n"))
		for _, f := range files {
			var list string

			fileTime := f.ModTime().Format("2006-01-02 15:04:05")
			lenName := max(45-len(f.Name()), 4)

			if f.IsDir() {
				list = "<a href=\"" + path.Join(req.URL.Path, f.Name()) + "\">" + f.Name() + "/</a>" + strings.Repeat(" ", lenName-1) + fileTime
			} else {
				fileSize := strconv.FormatInt(f.Size(), 10)
				lenSize := max(15-len(fileSize), 4)
				list = "<a href=\"" + path.Join(req.URL.Path, f.Name()) + "\">" + f.Name() + "</a>" + strings.Repeat(" ", lenName) + fileTime + strings.Repeat(" ", lenSize) + fileSize
			}

			w.Write([]byte(list + "\n"))
		}
		w.Write([]byte("</pre></html>"))
	} else {
		http.ServeContent(w, req, f.Name(), f.ModTime(), file)
	}
}
