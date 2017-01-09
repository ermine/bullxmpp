package main

import (
	"fmt"
	"flag"
	"os"
	"unicode"
	"path/filepath"
		"code.google.com/p/gcfg"
)

func readDir(path string) (files []string, err error) {
	d, err := os.Open(path)
	if err != nil {
		return
	}
	defer d.Close()
	fi, err := d.Readdir(-1)
	if err != nil {
		return
	}
	for _, fi := range fi {
		if fi.Mode().IsRegular() {
			files = append(files, fi.Name())
		}
	}
	return
}

type Config struct {
	Main struct {
		Indir string
	}
	Golang struct {
		Outdir string
	}
	Java struct {
		Outdir string
		Package_prefix string
	}
	Kotlin struct {
		Outdir string
		Package_prefix string
		Package_prefix_data string
	}
	Ocaml struct {
		Outdir string
		Extensionfile string
		Mlpack string
	}
	C struct {
		Outdir string
		Extensionfile string
		ExtensionTypes string
	}	
}

var cfg Config

var schemas []*Schema

func main() {
	cfgfile := flag.String("config", "generator.conf", "path to configuration file")
	var err error
	if _, err = os.Stat(*cfgfile); os.IsNotExist(err) {
		fmt.Println("No config file")
		os.Exit(1)
	}
	file := flag.String("file", "", "parse only this file")
	lang := flag.String("language", "golang", "Code language (golang, java, kotlin, ocaml or c)")
	flag.Parse()
	err = gcfg.ReadFileInto(&cfg, *cfgfile)
	checkError(err)
	if cfg.Main.Indir  == "" {
		fmt.Println("No indir")
		os.Exit(1)
	}

	var files []string
	if file != nil && *file != ""  {
		files = []string{*file}
	}  else {
		fmt.Println("Reading directory ", cfg.Main.Indir)
		files, err = readDir(cfg.Main.Indir)
		checkError(err)
	}
	for _, file := range files {
		schema := &Schema{
			Props: make(map[string]string),
		}
		err = schema.ParseFile(filepath.Join(cfg.Main.Indir, file))
		checkError(err)
		schemas = append(schemas, schema)
	}
	fmt.Println("Generating source code for ", *lang)
	switch *lang {
	case "golang":
		err = GolangGenerate()
	case "java":
		err = JavaGenerate()
	case "kotlin":
		err = KotlinGenerate()
	case "ocaml":
		err = OcamlGenerate()
	case "c":
		err = CGenerate()
	}
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func uppercase(s string) string {
	var result string
	for i, v := range s {
		result = string(unicode.ToUpper(v)) + s[i+1:]
		break
	}
	return result
}

