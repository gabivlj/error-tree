package rrors

import (
	"errors"
	"fmt"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

type rror struct {
	wrap error
	tmpl string
}

func (r rror) Error() string {
	v, _ := formatMultipleErr(r)
	return v
}

func (r rror) Unwrap() error {
	return r.wrap
}

type rrors struct {
	msg   string
	wraps []error
}

func (r rrors) Unwrap() []error {
	return r.wraps
}

func (r rrors) Error() string {
	v, _ := formatMultipleErr(r)
	return v
}

func formatMultipleErr(mainErr error) (v string, replace string) {
	if ror, ok := mainErr.(rror); ok {
		v, _ := formatMultipleErr(ror.wrap)
		return ConnectString(ror.tmpl, []string{v}), ""
	}

	rors, ok := mainErr.(rrors)
	if !ok {
		err := mainErr
		if multiple, ok := err.(interface{ Unwrap() []error }); ok {
			v, _ := formatMultipleErr(rrors{msg: ">", wraps: multiple.Unwrap()})
			return v, err.Error()
		}

		unwrap := errors.Unwrap(err)
		if unwrap == nil {
			return err.Error(), err.Error()
		}

		v, replace := formatMultipleErr(unwrap)
		return ConnectString(strings.Replace(err.Error(), replace, "", 1), []string{v}), err.Error()
	}

	strs := []string{}
	for _, err := range rors.wraps {
		v, _ := formatMultipleErr(err)
		strs = append(strs, v)
	}

	return ConnectString(rors.msg, strs), ""
}

func iterateErrsWithArgs(tmpl string, args []any, cb func(arg any)) (string, []any) {
	i := 0
	currArg := 0
	skipIndexes := []int{}
	nonErrArgs := []any{}
	for i < len(tmpl) {
		for i < len(tmpl) && tmpl[i] != '%' {
			i++
		}

		if i >= len(tmpl) {
			continue
		}

		// found %
		j := i + 1
		if j < len(tmpl) && tmpl[j] != '%' && tmpl[j] != ' ' {
			j++
		}

		if j < len(tmpl) && tmpl[j] == '%' {
			for j < len(tmpl) && tmpl[j] != '%' {
				j++
			}
		}

		if tmpl[i:j] == "%w" {
			cb(args[currArg])
			currArg++
			skipIndexes = append(skipIndexes, i)
		} else {
			nonErrArgs = append(nonErrArgs, args[currArg])
			currArg++
		}

		i = j
	}

	builder := strings.Builder{}
	prev := 0
	for _, idx := range skipIndexes {
		builder.WriteString(tmpl[prev:idx])
		prev = idx + 2
	}

	builder.WriteString(tmpl[prev:])
	return builder.String(), nonErrArgs
}

var basePath = ""

func split(p string) []string {
	s := []string{}
	for p != "/" && len(p) != 0 {
		path, rest := filepath.Split(p)
		s = append([]string{rest}, s...)
		p = path[:len(path)-1]
	}

	return s
}

func common(target, path2 string) string {
	pathSplit, pathSplit2 := split(target), split(path2)
	for i := range pathSplit {
		if i != len(pathSplit)-1 && i < len(pathSplit2) && pathSplit[i] == pathSplit2[i] {
			continue
		}

		return path.Join(pathSplit[i:]...)
	}

	return path.Join(pathSplit[len(pathSplit)-1:]...)
}

func Initialise() {
	_, file, _, _ := runtime.Caller(1)
	basePath = file
}

func Errorf(tmpl string, args ...any) error {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "?"
		line = 0
	}

	file = common(file, basePath)
	lineStr := "at " + file + fmt.Sprintf(":%d", line)
	idx := strings.Index(tmpl, "%w")
	if idx < 0 {
		return fmt.Errorf(tmpl+" "+lineStr, args...)
	}

	var errs []error
	newTmpl, argsNoErr := iterateErrsWithArgs(tmpl, args, func(arg any) {
		if err, ok := arg.(error); ok {
			errs = append(errs, err)
		} else {
			panic("so... how do we handle this?")
		}
	})

	if len(errs) == 1 {
		er := errs[0]
		if mul, ok := er.(interface{ Unwrap() []error }); ok {
			return rrors{
				wraps: mul.Unwrap(),
				msg:   newTmpl + lineStr,
			}
		}

		return rror{er, fmt.Sprintf(newTmpl+lineStr, argsNoErr...)}
	}

	return rrors{
		wraps: errs,
		msg:   fmt.Sprintf(newTmpl, argsNoErr...),
	}
}
