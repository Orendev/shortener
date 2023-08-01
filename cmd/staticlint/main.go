// staticlint is a tool for static analysis of Go programs.
// Assemble the staticline go build command.
//
// Build and run go build -o staticlint cmd/staticlint/main.go
// To check all packages beneath the current directory: ./staticlint -test=false ./...
//
// Package errcheck.Analyzer is a program for checking for unchecked errors in Go code.
// Package osexitcheck.Analyzer check whether to use a direct call to os.Exit.
// Package shadow.Analyzer check for possible unintended shadowing of variables.
// Package asmdecl.Analyzer defines an Analyzer that reports mismatches between assembly files and Go declarations.
// Package structtag.Analyzer check that struct field tags conform to reflect.Struct Tag.Get Also report certain struct tags (json, xml) used with unexported fields.
// Package printf.Analyzer check consistency of Printf format strings and arguments.
// Package assign.Analyzer defines an Analyzer that detects useless assignments.
// Package atomicalign defines an Analyzer that checks for non-64-bit-aligned arguments to sync/atomic functions.
// Package atomic defines an Analyzer that checks for common mistakes using the sync/atomic package.
// Package bools defines an Analyzer that detects common mistakes involving boolean operators.
// Package buildssa defines an Analyzer that constructs the SSA representation of an error-free package and returns the set of all functions within it.
// Package buildtag defines an Analyzer that checks build tags.
// Package cgocall defines an Analyzer that detects some violations of the cgo pointer passing rules.
// Package composite defines an Analyzer that checks for unkeyed composite literals.
// Package copylock defines an Analyzer that checks for locks erroneously passed by value.
// Package ctrlflow is an analysis that provides a syntactic control-flow graph (CFG) for the body of a function.
// Package deepequalerrors defines an Analyzer that checks for the use of reflect.DeepEqual with error values.
// Package defers defines an Analyzer that checks for common mistakes in defer statements.
// Package directive defines an Analyzer that checks known Go toolchain directives.
// The errorsas package defines an Analyzer that checks that the second argument to errors.As is a pointer to a type implementing error.
// Package framepointer defines an Analyzer that reports assembly code that clobbers the frame pointer before saving it.
// Package httpresponse defines an Analyzer that checks for mistakes using HTTP responses.
// Package inspect defines an Analyzer that provides an AST inspector (golang.org/x/tools/go/ast/inspector.Inspector) for the syntax trees of a package.
// Package loopclosure defines an Analyzer that checks for references to enclosing loop variables from within nested functions.
// Package nilfunc defines an Analyzer that checks for useless comparisons against nil.
// The pkgfact package is a demonstration and test of the package fact mechanism.
// Package reflectvaluecompare defines an Analyzer that checks for accidentally using == or reflect.DeepEqual to compare reflect.Value values.
// Package shift defines an Analyzer that checks for shifts that exceed the width of an integer.
// Package sigchanyzer defines an Analyzer that detects misuse of unbuffered signal as argument to signal.Notify.
// Package sortslice defines an Analyzer that checks for calls to sort.Slice that do not use a slice type as first argument.
// Package stdmethods defines an Analyzer that checks for misspellings in the signatures of methods similar to well-known interfaces.
package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/Orendev/shortener/pkg/osexitcheck"
	"github.com/kisielk/errcheck/errcheck"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/asmdecl"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/atomicalign"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/analysis/passes/buildtag"
	"golang.org/x/tools/go/analysis/passes/cgocall"
	"golang.org/x/tools/go/analysis/passes/composite"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/ctrlflow"
	"golang.org/x/tools/go/analysis/passes/deepequalerrors"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/framepointer"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/pkgfact"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/reflectvaluecompare"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/sigchanyzer"
	"golang.org/x/tools/go/analysis/passes/sortslice"
	"golang.org/x/tools/go/analysis/passes/stdmethods"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/analysis/facts/directives"
	"honnef.co/go/tools/staticcheck"
)

// Config is the name of the configuration file for staticcheck.
const Config = `config.json`

// ConfigData описывает структуру файла конфигурации Staticcheck.
type ConfigData struct {
	Staticcheck []string
}

func main() {

	appfile, err := os.Executable()
	if err != nil {
		panic(err)
	}
	data, err := os.ReadFile(filepath.Join(filepath.Dir(appfile), Config))
	if err != nil {
		panic(err)
	}
	var cfg ConfigData
	if err = json.Unmarshal(data, &cfg); err != nil {
		panic(err)
	}

	mychecks := []*analysis.Analyzer{
		osexitcheck.Analyzer,
		asmdecl.Analyzer,
		printf.Analyzer,
		shadow.Analyzer,
		structtag.Analyzer,
		errcheck.Analyzer,
		assign.Analyzer,
		atomic.Analyzer,
		atomicalign.Analyzer,
		bools.Analyzer,
		buildssa.Analyzer,
		buildtag.Analyzer,
		cgocall.Analyzer,
		composite.Analyzer,
		copylock.Analyzer,
		ctrlflow.Analyzer,
		deepequalerrors.Analyzer,
		directives.Analyzer,
		errorsas.Analyzer,
		framepointer.Analyzer,
		httpresponse.Analyzer,
		inspect.Analyzer,
		loopclosure.Analyzer,
		nilfunc.Analyzer,
		pkgfact.Analyzer,
		reflectvaluecompare.Analyzer,
		shift.Analyzer,
		sigchanyzer.Analyzer,
		sortslice.Analyzer,
		stdmethods.Analyzer,
	}

	checks := make(map[string]bool)
	for _, v := range cfg.Staticcheck {
		checks[v] = true
	}

	// добавляем анализаторы из staticcheck, которые указаны в файле конфигурации
	for _, v := range staticcheck.Analyzers {
		if strings.HasPrefix(v.Analyzer.Name, "SA") {
			mychecks = append(mychecks, v.Analyzer)
		}
		if checks[v.Analyzer.Name] {
			mychecks = append(mychecks, v.Analyzer)
		}
	}
	multichecker.Main(
		mychecks...,
	)
}
