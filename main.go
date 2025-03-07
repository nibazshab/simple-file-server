package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
)

var rootfs *os.Root

type info struct {
	name string
	size string
	time string
	i    int
	j    int
}

func main() {
	port := flag.String("port", "8080", "server port")
	root := flag.String("path", "./", "server root path")
	flag.Parse()

	*root, _ = filepath.Abs(*root)
	rootfs, _ = os.OpenRoot(*root)

	fmt.Printf("@ 0.0.0.0:%s\n@ %s\n", *port, *root)
	http.HandleFunc("/", server)
	http.ListenAndServe(":"+*port, nil)
}

func server(w http.ResponseWriter, req *http.Request) {
	fmt.Printf("%s\n", req.URL.Path)

	p := strings.TrimPrefix(req.URL.Path, "/")
	if p == "" {
		p = "."
	}

	fi, err := rootfs.Stat(p)
	if err != nil {
		http.NotFound(w, req)
		return
	}

	f, _ := rootfs.Open(p)
	defer f.Close()

	if !fi.IsDir() {
		http.ServeContent(w, req, fi.Name(), fi.ModTime(), f)
		return
	}

	fl, _ := f.Readdir(-1)

	slices.SortStableFunc(fl, func(a, b os.FileInfo) int {
		if a.IsDir() != b.IsDir() {
			if a.IsDir() {
				return -1
			}
			return 1
		}
		return 0
	})

	var b strings.Builder

	of := "Index of " + req.URL.Path
	b.WriteString(fmt.Sprintf("<title>%s</title><h1>%s</h1>", of, of))
	b.WriteString("<hr><pre><a href=\"../\">../</a>\n")

	var i info
	for _, l := range fl {
		if l.IsDir() {
			i.name = l.Name() + "/"
			i.size = "-"
		} else {
			i.name = l.Name()
			i.size = strconv.FormatInt(l.Size(), 10)
		}
		i.time = l.ModTime().Format("02-Jan-2006 15:04")
		i.i = max(51-len(i.name), 1)
		i.j = max(20-len(i.size), 1)

		b.WriteString(fmt.Sprintf("<a href=\"%s\">%s</a>", i.name, i.name))
		b.WriteString(strings.Repeat(" ", i.i))
		b.WriteString(i.time)
		b.WriteString(strings.Repeat(" ", i.j))
		b.WriteString(i.size)
		b.WriteString("\n")
	}
	b.WriteString("</pre><hr>")
	w.Write([]byte(b.String()))
}
