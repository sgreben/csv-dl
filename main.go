package main

import (
	"bytes"
	"encoding/csv"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"text/template"
)

const appName = "csv-dl"
const (
	columnFunctionName = "column"
	fieldFunctionName  = "field"
)

var nonzeroExit = false
var newline = []byte{'\n'}
var httpClient = http.DefaultClient

var config struct {
	Quiet       bool
	UseHeader   bool
	SkipHeader  bool
	Schema      string
	Links       stringsVar
	Force       bool
	Parallelism int
}

func init() {
	log.SetOutput(os.Stderr)
	config.Parallelism = runtime.NumCPU()
	flag.BoolVar(&config.Quiet, "q", config.Quiet, "(alias for -quiet)")
	flag.BoolVar(&config.Quiet, "quiet", config.Quiet, "suppress all logging")
	flag.BoolVar(&config.Force, "f", config.Force, "(alias for -force-overwrite)")
	flag.BoolVar(&config.Force, "force-overwrite", config.Force, "overwrite existing files")
	flag.BoolVar(&config.UseHeader, "u", config.UseHeader, "(alias for -use-csv-header)")
	flag.BoolVar(&config.UseHeader, "use-csv-header", config.UseHeader, "assume the first row is the CSV header, use it as a schema")
	flag.BoolVar(&config.SkipHeader, "skip-csv-eader", config.SkipHeader, "assume the first row is the CSV header, skip it")
	flag.StringVar(&config.Schema, "s", config.Schema, "(alias for -schema)")
	flag.StringVar(&config.Schema, "schema", config.Schema, "use the given CSV expression as the table schema")
	flag.IntVar(&config.Parallelism, "p", config.Parallelism, "(alias for -parallel)")
	flag.IntVar(&config.Parallelism, "parallel", config.Parallelism, "number of parallel connections")
	flag.Var(&config.Links, "l", "(alias for -link)")
	flag.Var(&config.Links, "link", `a link to download, may use go {{template}} syntax and refer to data columns by index (column i) or name (field "f")`)
	flag.Parse()
	if config.Quiet {
		log.SetOutput(ioutil.Discard)
	}
}

func parseLinkTemplates(root *template.Template) (out []*template.Template) {
	for _, linkTemplateText := range config.Links.Values {
		linkTemplate, err := root.New("link").Parse(linkTemplateText)
		if err != nil {
			log.Fatalf("parse link template %q: %v", linkTemplateText, err)
		}
		out = append(out, linkTemplate)
	}
	return
}

func buildColumnNames(r *csv.Reader) (out []string) {
	switch {
	case config.Schema != "":
		out = strings.Split(config.Schema, ",")
	case config.UseHeader:
		header, err := r.Read()
		if err != nil {
			log.Fatalf("read header: %v", err)
		}
		out = header
	}
	for i, name := range out {
		out[i] = strings.TrimSpace(name)
	}
	return
}

func download(urls chan string) {
	for link := range urls {
		name := filepath.Base(link)
		if !config.Force {
			if _, err := os.Stat(name); err == nil {
				log.Printf("file %q already exists, skipping %q", name, link)
				continue
			}
		}
		req, err := http.NewRequest(http.MethodGet, link, nil)
		if err != nil {
			nonzeroExit = true
			log.Printf("error building request for %q: %v", link, err)
			continue
		}
		resp, err := httpClient.Do(req)
		if err != nil {
			nonzeroExit = true
			log.Printf("error executing request for %q: %v", link, err)
			continue
		}
		file, err := os.OpenFile(name, os.O_CREATE|os.O_RDWR, 0600)
		if err != nil {
			resp.Body.Close()
			nonzeroExit = true
			log.Printf("error opening file %q for %q: %v", name, link, err)
			continue
		}
		if _, err := io.Copy(file, resp.Body); err != nil {
			resp.Body.Close()
			file.Close()
			nonzeroExit = true
			log.Printf("error writing file %q for %q: %v", name, link, err)
			continue
		}
		resp.Body.Close()
		file.Close()
	}
}

func main() {
	links := make(chan string)
	var wg sync.WaitGroup
	for i := 0; i == 0 || i < config.Parallelism; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			download(links)
		}()
	}

	r := csv.NewReader(os.Stdin)
	r.FieldsPerRecord = -1
	columnNames := buildColumnNames(r)
	templateRoot := template.New(appName)
	var columnForName func(string) string
	var columnForIndex func(int) string
	templateFuncs := map[string]interface{}{
		columnFunctionName: func(i int) string { return columnForIndex(i) },
		fieldFunctionName:  func(name string) string { return columnForName(name) },
	}
	columnIndexForName := map[string]int{}
	for i, name := range columnNames {
		if name == "" {
			continue
		}
		columnIndexForName[name] = i
	}
	templateRoot.Funcs(templateFuncs)
	linkTemplates := parseLinkTemplates(templateRoot)

	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("error reading row: %v", err)
			nonzeroExit = true
			continue
		}
		columnForIndex = func(i int) string {
			if i < 0 || i > len(row) {
				return ""
			}
			return row[i]
		}
		columnForName = func(name string) string {
			if i, ok := columnIndexForName[name]; ok {
				return columnForIndex(i)
			}
			return ""
		}
		var buf bytes.Buffer
		for i, t := range linkTemplates {
			buf.Reset()
			err := t.Execute(&buf, row)
			if err != nil {
				log.Printf("error executing template %q: %v", config.Links.Values[i], err)
				nonzeroExit = true
				continue
			}
			link := buf.String()
			log.Println(link)
			links <- link
		}
	}
	close(links)
	wg.Wait()
	if nonzeroExit {
		os.Exit(1)
	}
}
