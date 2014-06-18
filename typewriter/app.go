package typewriter

import (
	"bytes"
	"fmt"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// +test foo:"Bar" baz:"qux,thing"
type app struct {
	// All typewriter.Type found in the current directory.
	Types []Type
	// All typewriter.TypeWriters registered on init.
	TypeWriters []TypeWriter
	Directive   string
}

// NewApp parses the current directory, collecting Types and their related information.
func NewApp(directive string) (*app, error) {
	return NewAppFiltered(directive, nil)
}

// NewAppNewAppFiltered parses the current directory, collecting Types and their related information. Pass a filter to limit which files are operated on.
func NewAppFiltered(directive string, filter func(os.FileInfo) bool) (*app, error) {
	a := &app{
		Directive: directive,
	}

	typs, err := getTypes(directive, filter)
	if err != nil {
		return a, err
	}

	a.Types = typs
	a.TypeWriters = typeWriters
	return a, nil
}

// Individual TypeWriters register on init, keyed by name
var typeWriters []TypeWriter

// Register allows template packages to make themselves known to a 'parent' package, usually in the init() func.
// Comparable to the approach taken by builtin image package for registration of image types (eg image/png).
// Your program will do something like:
//	import (
//		"github.com/clipperhouse/gen/typewriter"
//		_ "github.com/clipperhouse/gen/typewriters/container"
//	)
func Register(tw TypeWriter) error {
	for _, v := range typeWriters {
		if v.Name() == tw.Name() {
			return fmt.Errorf("A TypeWriter by the name %s has already been registered", tw.Name())
		}
	}
	typeWriters = append(typeWriters, tw)
	return nil
}

// WriteAll writes the generated code for all Types and TypeWriters in the App to respective files.
func (a *app) WriteAll() error {
	// TypeWriters which will write for each type; use string as key because Type is not comparable
	var writersByType = make(map[string][]TypeWriter)

	// validate them all (don't fail halfway)
	for _, t := range a.Types {
		for _, tw := range a.TypeWriters {
			write, err := tw.Validate(t)
			if err != nil {
				return err
			}
			if write {
				writersByType[t.String()] = append(writersByType[t.String()], tw)
			}
		}
	}

	// one buffer for each file, keyed by file name
	buffers := make(map[string]*bytes.Buffer)

	// write the generated code for each Type & TypeWriter into memory
	for _, t := range a.Types {
		for _, tw := range writersByType[t.String()] {
			var b bytes.Buffer
			write(&b, a, t, tw)
			f := strings.ToLower(fmt.Sprintf("%s_%s.go", t.LocalName(), tw.Name()))
			buffers[f] = &b
		}
	}

	// validate generated ast's before committing to files
	for f, b := range buffers {
		if _, err := parser.ParseFile(token.NewFileSet(), f, b.String(), 0); err != nil {
			// TODO: prompt to write (ignored) _file on error? parsing errors are meaningless without.
			return err
		}
	}

	// format and commit to files
	for f, b := range buffers {
		src, err := format.Source(b.Bytes())

		// shouldn't be an error if the ast parsing above succeeded
		if err != nil {
			return err
		}

		if err := writeFile(f, src); err != nil {
			return err
		}
	}

	return nil
}

func write(w io.Writer, a *app, t Type, tw TypeWriter) {
	// start with byline at top, give future readers some background
	// on where the file came from
	bylineFmt := `// Generated by: %s
// TypeWriter: %s
// Directive: %s on %s

`
	caller := filepath.Base(os.Args[0])
	byline := fmt.Sprintf(bylineFmt, caller, tw.Name(), a.Directive, t.String())
	w.Write([]byte(byline))

	// TypeWriter's headers, usually licenses & credits
	tw.WriteHeader(w, t)

	// err on the side of extra line breaks; format will tighten it up later
	w.Write([]byte(strings.Repeat("\n", 2)))

	pkgFmt := `package %s

`
	pkgDecl := fmt.Sprintf(pkgFmt, t.Package.Name())
	w.Write([]byte(pkgDecl))

	importsTmpl.Execute(w, tw.Imports(t))

	tw.WriteBody(w, t)
}

func writeFile(filename string, byts []byte) error {
	w, err := os.Create(filename)

	if err != nil {
		return err
	}

	defer w.Close()

	w.Write(byts)

	// TODO: make this optional or do a proper logging/verbosity thing
	fmt.Printf("  writing %s\n", filename)

	return nil
}

var importsTmpl = template.Must(template.New("imports").Parse(`{{if gt (len .) 0}}
import ({{range .}}
	"{{.}}"{{end}}
)
{{end}}
`))
