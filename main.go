package main

import (
	"flag"
	"fmt"
	"io/fs"
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

type InfoShow struct {
	Name string
	Time string
	Size string
	Path string
}

type FileInfo struct {
	Path string
	Obj  *os.File
	Info fs.FileInfo
}

func main() {
	port = flag.String("port", "8080", "server port")
	root = flag.String("path", "./", "server root path")
	flag.Parse()

	fmt.Println("start HTTP server @ 0.0.0.0:" + *port + "\n" + "load storage @ " + *root)

	http.HandleFunc("/", FileServer)
	http.ListenAndServe(":"+*port, nil)
}

func FileServer(w http.ResponseWriter, req *http.Request) {
	var file FileInfo
	var info InfoShow

	file.Path = path.Join(*root, req.URL.Path)

	var err error
	file.Obj, err = os.Open(file.Path)
	if err != nil {
		http.NotFound(w, req)
		return
	}
	defer file.Obj.Close()

	file.Info, _ = file.Obj.Stat()

	if file.Info.IsDir() {
		fileLists, _ := file.Obj.Readdir(-1)

		w.Write([]byte("<html><pre><a href=\"../\">../</a>" + "\n"))
		for _, f := range fileLists {
			var list string

			info.Time = f.ModTime().Format("2006-01-02 15:04:05")
			info.Name = f.Name()
			info.Path = path.Join(req.URL.Path, info.Name)

			lenName := max(45-len(info.Name), 4)

			if f.IsDir() {
				list = "<a href=\"" + info.Path + "\">" + info.Name + "/</a>" + strings.Repeat(" ", lenName-1) + info.Time
			} else {
				infoSize := f.Size()
				if infoSize > 10240 {
					infoSize >>= 10
					info.Size = strconv.FormatInt(infoSize, 10) + "kb"
				} else {
					info.Size = strconv.FormatInt(infoSize, 10)
				}

				lenSize := max(15-len(info.Size), 4)

				list = "<a href=\"" + info.Path + "\">" + info.Name + "</a>" + strings.Repeat(" ", lenName) + info.Time + strings.Repeat(" ", lenSize) + info.Size
			}

			w.Write([]byte(list + "\n"))
		}
		w.Write([]byte("</pre></html>"))
	} else {
		http.ServeContent(w, req, file.Info.Name(), file.Info.ModTime(), file.Obj)
	}
}
