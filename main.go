package main

import (
	"flag"
	"fmt"
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

func main() {
	Port = flag.String("port", "8080", "server port")
	Root = flag.String("path", "./", "server root path")
	flag.Parse()

	fmt.Printf("start HTTP server @ 0.0.0.0:%s\nload storage @ %s", *Port, *Root)

	http.HandleFunc("/", fileServer)
	err := http.ListenAndServe(":"+*Port, nil)
	if err != nil {
		fmt.Printf("start error: %v", err)
		os.Exit(1)
	}
}

func fileServer(w http.ResponseWriter, req *http.Request) {
	fPath := path.Join(*Root, req.URL.Path)
	fObj, err := os.Open(fPath)
	if err != nil {
		http.NotFound(w, req)
		return
	}
	defer fObj.Close()
	fInfo, _ := fObj.Stat()

	if fInfo.IsDir() {
		fList, _ := fObj.Readdir(-1)
		sort.SliceStable(fList, func(i, j int) bool {
			if fList[i].IsDir() == fList[j].IsDir() {
				return fList[i].Name() < fList[j].Name()
			}
			return fList[i].IsDir() && !fList[j].IsDir()
		})

		var oLastIndex string
		oIndex := req.URL.Path
		if oIndex == "/" {
			oLastIndex = oIndex
		} else {
			oLastIndex = path.Dir(oIndex)
			oIndex += "/"
		}

		head := fmt.Sprintf("<html><h1>Index of %s</h1><hr/><pre><a href=\"%s\">../</a>\n", oIndex, oLastIndex)
		_, _ = w.Write([]byte(head))

		for _, _f := range fList {
			var li string
			mName := _f.Name()
			mUrl := path.Join(req.URL.Path, mName)
			mTime := _f.ModTime().Format("2006-01-02 15:04:05")
			mNameLength := max(50-len(mName), 1)
			const mSizeLength = 19

			if _f.IsDir() {
				sn := strings.Repeat(" ", mNameLength)
				sl := strings.Repeat(" ", mSizeLength)
				li = fmt.Sprintf("<a href=\"%s\">%s/</a>%s%s%s-\n", mUrl, mName, sn, mTime, sl)
			} else {
				var mSize string
				_size := _f.Size()

				if _size > 10240 {
					_size >>= 10
					mSize = strconv.FormatInt(_size, 10) + "kb"
				} else {
					mSize = strconv.FormatInt(_size, 10)
				}

				sn := strings.Repeat(" ", mNameLength+1)
				sl := strings.Repeat(" ", max(mSizeLength-len(mSize), 1)+1)
				li = fmt.Sprintf("<a href=\"%s\">%s</a>%s%s%s%s\n", mUrl, mName, sn, mTime, sl, mSize)
			}

			_, _ = w.Write([]byte(li))
		}
		_, _ = w.Write([]byte("</pre><hr/></html>"))
	} else {
		http.ServeContent(w, req, fInfo.Name(), fInfo.ModTime(), fObj)
	}
}
