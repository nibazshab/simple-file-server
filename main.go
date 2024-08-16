package main

import (
	"flag"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
)

var (
	port *string
	root *string
)

var Version string

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

	fmt.Println("Version: " + Version)
	fmt.Println("start HTTP server @ 0.0.0.0:" + *port + "\n" + "load storage @ " + *root)

	http.HandleFunc("/", fileServer)
	http.ListenAndServe(":"+*port, nil)
}

func fileServer(w http.ResponseWriter, req *http.Request) {
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

		sort.SliceStable(fileLists, func(i, j int) bool {
			if fileLists[i].IsDir() == fileLists[j].IsDir() {
				return fileLists[i].Name() < fileLists[j].Name()
			}
			return fileLists[i].IsDir() && !fileLists[j].IsDir()
		})

		indexof := req.URL.Path
		var lastdir string

		if indexof == "/" {
			lastdir = indexof
		} else {
			lastdir = path.Dir(indexof)
			indexof += "/"
		}

		w.Write([]byte("<html><h1>Index of " + indexof + "</h1><hr/><pre><a href=\"" + lastdir + "\">../</a>" + "\n"))
		for _, f := range fileLists {
			var list string

			info.Time = f.ModTime().Format("2006-01-02 15:04:05")
			info.Name = f.Name()
			info.Path = path.Join(req.URL.Path, info.Name)

			lenName := max(51-len(info.Name), 1)
			lenSize := 20

			if f.IsDir() {
				list = "<a href=\"" + info.Path + "\">" + info.Name + "/</a>" + strings.Repeat(" ", lenName-1) + info.Time + strings.Repeat(" ", lenSize-1) + "-"
			} else {
				infoSize := f.Size()
				if infoSize > 10240 {
					infoSize >>= 10
					info.Size = strconv.FormatInt(infoSize, 10) + "kb"
				} else {
					info.Size = strconv.FormatInt(infoSize, 10)
				}

				lenSize = max(lenSize-len(info.Size), 1)
				list = "<a href=\"" + info.Path + "\">" + info.Name + "</a>" + strings.Repeat(" ", lenName) + info.Time + strings.Repeat(" ", lenSize) + info.Size
			}

			w.Write([]byte(list + "\n"))
		}
		w.Write([]byte("</pre><hr/></html>"))
	} else {
		http.ServeContent(w, req, file.Info.Name(), file.Info.ModTime(), file.Obj)
	}
}
