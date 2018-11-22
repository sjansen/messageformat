package main

import (
	"bytes"
	"fmt"
	"go/format"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	for _, src := range os.Args[1:] {
		if src != filepath.Base(src) || !strings.HasSuffix(src, ".go.tmpl") {
			fmt.Fprintf(os.Stderr, "invalid template filename: %q\n", src)
			os.Exit(1)
		}

		fi, err := os.Stat(src)
		if os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "file not found: %q\n", src)
			os.Exit(1)
		} else if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		dst := strings.TrimSuffix(src, ".tmpl")
		fmt.Println("generating:", dst)

		tmpl := template.New(src).Funcs(template.FuncMap{
			"split": func(s string) []string { return strings.Split(s, ",") },
		})
		if _, err = tmpl.ParseFiles(src); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		var buf bytes.Buffer
		buf.Write([]byte(
			"// Code generated by scripts/tmpl-to-go. DO NOT EDIT.\n",
		))

		if err = tmpl.Execute(&buf, nil); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		code := buf.Bytes()
		if formatted, err := format.Source(code); err == nil {
			code = formatted
		} else {
			_ = ioutil.WriteFile(dst, code, fi.Mode())
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if err = ioutil.WriteFile(dst, code, fi.Mode()); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}