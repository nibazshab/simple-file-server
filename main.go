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
	Port *string
	Root *string
)

var Version string

type File struct {
	Path string
	Obj  *os.File
	Info fs.FileInfo
	List []os.FileInfo
}

type Meta struct {
	Name       string
	Url        string
	Time       string
	Size       string
	NameLength int
	SizeLength int
}

type IndexOf struct {
	Index     string
	LastIndex string
}

func main() {
	Port = flag.String("port", "8080", "server port")
	Root = flag.String("path", "./", "server root path")
	flag.Parse()

	fmt.Println("Version: " + Version)
	fmt.Println("start HTTP server @ 0.0.0.0:" + *Port + "\n" + "load storage @ " + *Root)

	http.HandleFunc("/", fileServer)
	err := http.ListenAndServe(":"+*Port, nil)
	if err != nil {
		fmt.Printf("start error: %v", err)
		os.Exit(1)
	}
}

func fileServer(w http.ResponseWriter, req *http.Request) {
	var f File
	var m Meta
	var o IndexOf

	f.Path = path.Join(*Root, req.URL.Path)

	var err error
	f.Obj, err = os.Open(f.Path)
	if err != nil {
		http.NotFound(w, req)
		return
	}
	defer f.Obj.Close()

	f.Info, _ = f.Obj.Stat()
	if f.Info.IsDir() {
		f.List, _ = f.Obj.Readdir(-1)
		sort.SliceStable(f.List, func(i, j int) bool {
			if f.List[i].IsDir() == f.List[j].IsDir() {
				return f.List[i].Name() < f.List[j].Name()
			}
			return f.List[i].IsDir() && !f.List[j].IsDir()
		})

		o.Index = req.URL.Path
		if o.Index == "/" {
			o.LastIndex = o.Index
		} else {
			o.LastIndex = path.Dir(o.Index)
			o.Index += "/"
		}

		w.Write([]byte("<html><h1>Index of " + o.Index + "</h1><hr/><pre><a href=\"" + o.LastIndex + "\">../</a>" + "\n"))
		for _, _f := range f.List {
			var li string
			var sn string
			var sl string

			m.Name = _f.Name()
			m.Url = path.Join(req.URL.Path, m.Name)
			m.Time = _f.ModTime().Format("2006-01-02 15:04:05")
			m.NameLength = max(50-len(m.Name), 1)
			m.SizeLength = 19

			if _f.IsDir() {
				sn = strings.Repeat(" ", m.NameLength)
				sl = strings.Repeat(" ", m.SizeLength)
				li = "<a href=\"" + m.Url + "\">" + m.Name + "/</a>" + sn + m.Time + sl + "-"
			} else {
				_size := _f.Size()
				if _size > 10240 {
					_size >>= 10
					m.Size = strconv.FormatInt(_size, 10) + "kb"
				} else {
					m.Size = strconv.FormatInt(_size, 10)
				}

				sn = strings.Repeat(" ", m.NameLength+1)
				sl = strings.Repeat(" ", max(m.SizeLength-len(m.Size), 1)+1)
				li = "<a href=\"" + m.Url + "\">" + m.Name + "</a>" + sn + m.Time + sl + m.Size
			}
			w.Write([]byte(li + "\n"))
		}
		w.Write([]byte("</pre><hr/></html>"))
	} else {
		http.ServeContent(w, req, f.Info.Name(), f.Info.ModTime(), f.Obj)
	}
}
