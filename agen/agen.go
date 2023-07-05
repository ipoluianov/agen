package agen

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"

	"github.com/yuin/goldmark"
)

func ProcessDirectory(path string) (err error) {
	srcDirectory := path + "/src"
	outDirectory := path + "/out"
	tmp, err := os.ReadFile(path + "/tmp.html")
	if err != nil {
		return
	}
	return processDirectory(outDirectory, srcDirectory, string(tmp), 0)
}

func processDirectory(outDirectory string, path string, tmp string, depth int) (err error) {
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
			processDirectory(outDirectory+"/"+item.Name(), path+"/"+item.Name(), tmp, depth+1)
		} else {
			if ext != ".md" {
				continue
			}
			shortName := strings.TrimSuffix(item.Name(), filepath.Ext(item.Name()))
			outFile := outDirectory + "/" + shortName + ".html"
			err = ProcessFile(tmp, outFile, path+"/"+shortName+".md")
		}
	}

	return
}

func ProcessFile(tmp string, outFile string, path string) error {
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
	os.WriteFile(outFile, []byte(htmlFile), 0666)
	return nil
}
