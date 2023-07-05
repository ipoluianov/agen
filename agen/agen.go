package agen

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	_ "embed"

	"github.com/yuin/goldmark"
)

//go:embed "tmp.html"
var tmpDefault []byte

func ProcessDirectory(path string) (err error) {
	srcDirectory := path + "/src"
	outDirectory := path + "/out"

	err = os.MkdirAll(srcDirectory, 0777)
	if err != nil {
		return
	}

	tmp, err := os.ReadFile(path + "/tmp.html")
	if err != nil {
		fmt.Println("warning: template file is not found (tmp.html)")
		tmp = tmpDefault
		fmt.Println("warning: using default template file (tmp.html)")
		err = os.WriteFile(path+"/tmp.html", tmpDefault, 0666)
		if err != nil {
			fmt.Println("warning: cannot write default template file (tmp.html):", err)
		} else {
			fmt.Println("warning: default template file (tmp.html) has been written")
		}
	}
	return processDirectory(outDirectory, srcDirectory, string(tmp), 0)
}

func processDirectory(outDirectory string, path string, tmp string, depth int) (err error) {
	fmt.Println("process directory:", path)

	if depth > 10 {
		return nil
	}
	err = os.MkdirAll(outDirectory, 0777)
	if err != nil {
		return
	}

	items, err := os.ReadDir(path)
	if err != nil {
		return
	}

	for _, item := range items {
		ext := filepath.Ext(item.Name())
		if item.IsDir() {
			err = processDirectory(outDirectory+"/"+item.Name(), path+"/"+item.Name(), tmp, depth+1)
			if err != nil {
				return
			}
		} else {
			if ext != ".md" {
				continue
			}
			shortName := strings.TrimSuffix(item.Name(), filepath.Ext(item.Name()))
			outFile := outDirectory + "/" + shortName + ".html"
			err = ProcessFile(tmp, outFile, path+"/"+shortName+".md")
			if err != nil {
				return
			}
		}
	}

	return
}

func ProcessFile(tmp string, outFile string, path string) error {
	fmt.Println("process file:", path)
	// Read the file
	bs, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	// Convert MD to HTML
	var buf bytes.Buffer
	err = goldmark.Convert(bs, &buf)
	if err != nil {
		return err
	}

	// Apply Template
	htmlFile := strings.ReplaceAll(tmp, "%CONTENT%", buf.String())
	err = os.WriteFile(outFile, []byte(htmlFile), 0666)
	return err
}
