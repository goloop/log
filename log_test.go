package log

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

const (
	testSkip     = 3
	testSkipSeek = 1
)

// TestNew tests Log.New method.
func TestNew(t *testing.T) {
	var (
		log *Log
		err error
	)

	// Create log with custom leveles.
	log, err = New(Info, Debug)
	if err != nil {
		t.Error(err)
	}

	if ok, _ := log.Config.Levels.All(Info, Debug); !ok {
		t.Error("the Info and Debug levels must be active")
	}

	if ok, _ := log.Config.Levels.All(Info, Error); ok {
		t.Error("the Error level shouldn't be active")
	}

	// Create log with defaults leveles.
	log, err = New()
	if err != nil {
		t.Error(err)
	}

	if ok, _ := log.Config.Levels.All(Info, Debug); !ok {
		t.Error("the Info and Debug levels must be active")
	}

	if ok, _ := log.Config.Levels.All(Panic, Fatal, Error, Warn); !ok {
		t.Error("all log levels must be set")
	}
}

// TestCopy tests Log.Copy method.
func TestCopy(t *testing.T) {
	var log, err = New(Info, Debug)
	if err != nil {
		t.Error(err)
	}

	log.Config.Levels.Add(Error)
	log.Config.Formats.Set(FuncName, LineNumber)
	log.Config.FatalStatusCode = 7

	clone := log.Copy()
	clone.Config.FatalStatusCode = 3

	if log.Config.Levels != clone.Config.Levels {
		t.Error("log levels don't match")
	}

	if log.Config.Formats != clone.Config.Formats {
		t.Error("format styles don't match")
	}

	if log.Config.FatalStatusCode == clone.Config.FatalStatusCode {
		t.Error("the FatalStatusCode must be different")
	}
}

// TestEcho tests Log.echo method.
func TestEcho(t *testing.T) {
	type test struct {
		level   LevelFlag
		levels  []LevelFlag
		formats []FormatFlag
		data    []interface{}
	}

	var tests = []test{
		{
			Error,
			[]LevelFlag{Info, Debug, Trace},
			[]FormatFlag{FilePath},
			[]interface{}{"1", "2", "3"},
		},
		{
			Info,
			[]LevelFlag{Info, Debug, Trace},
			[]FormatFlag{FilePath, FuncName, LineNumber},
			[]interface{}{"1", "2", "3"},
		},
		{
			Debug,
			[]LevelFlag{Info, Debug, Trace},
			[]FormatFlag{FuncName},
			[]interface{}{"1", "2", "3"},
		},
	}

	for i, s := range tests {
		var buf = new(bytes.Buffer)
		l, _ := New(s.levels...)
		l.skip += testSkipSeek
		l.Config.Formats.Set(s.formats...)
		l.echo(skip, s.level, buf, s.data...)

		exp := ""
		res := buf.String()
		ss := getStackSlice(testSkip)

		if ok, _ := l.Config.Levels.All(s.level); ok {
			res = res[19:]
			exp = getPrefix(s.level, l.Config, ss) +
				fmt.Sprint(s.data...)
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}

// TestEchof tests Log.echof method.
func TestEchof(t *testing.T) {
	type test struct {
		level   LevelFlag
		levels  []LevelFlag
		formats []FormatFlag
		data    []interface{}
		layout  string
	}

	var tests = []test{
		{
			Error,
			[]LevelFlag{Info, Debug, Trace},
			[]FormatFlag{FilePath},
			[]interface{}{"1", "2", "3"},
			"%s + %s - %s",
		},
		{
			Info,
			[]LevelFlag{Info, Debug, Trace},
			[]FormatFlag{FilePath, FuncName, LineNumber},
			[]interface{}{"1", "2", "3"},
			"%s / %s + %s",
		},
		{
			Debug,
			[]LevelFlag{Info, Debug, Trace},
			[]FormatFlag{FuncName},
			[]interface{}{"1", "2", "3"},
			"%s * %s + %s",
		},
	}

	for i, s := range tests {
		var buf = new(bytes.Buffer)
		l, _ := New(s.levels...)
		l.skip += testSkipSeek
		l.Config.Formats.Set(s.formats...)
		l.echof(skip, s.level, buf, s.layout, s.data...)

		exp := ""
		res := buf.String()
		ss := getStackSlice(testSkip)

		if ok, _ := l.Config.Levels.All(s.level); ok {
			res = res[19:]
			exp = getPrefix(s.level, l.Config, ss) +
				fmt.Sprintf(s.layout, s.data...)
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}

// TestEcholn tests Log.echoln method.
func TestEcholn(t *testing.T) {
	type test struct {
		level   LevelFlag
		levels  []LevelFlag
		formats []FormatFlag
		data    []interface{}
	}

	var tests = []test{
		{
			Error,
			[]LevelFlag{Info, Debug, Trace},
			[]FormatFlag{FilePath},
			[]interface{}{"1", "2", "3"},
		},
		{
			Info,
			[]LevelFlag{Info, Debug, Trace},
			[]FormatFlag{FilePath, FuncName, LineNumber},
			[]interface{}{"1", "2", "3"},
		},
		{
			Debug,
			[]LevelFlag{Info, Debug, Trace},
			[]FormatFlag{FuncName},
			[]interface{}{"1", "2", "3"},
		},
	}

	for i, s := range tests {
		var buf = new(bytes.Buffer)
		l, _ := New(s.levels...)
		l.skip += testSkipSeek
		l.Config.Formats.Set(s.formats...)
		l.echoln(skip, s.level, buf, s.data...)

		exp := ""
		res := buf.String()
		ss := getStackSlice(testSkip)

		if ok, _ := l.Config.Levels.All(s.level); ok {
			res = res[19:]
			exp = getPrefix(s.level, l.Config, ss) +
				fmt.Sprintln(s.data...)
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}

// TestFpanic tests Fpanic method.
func TestFpanic(t *testing.T) {
	type test struct {
		data    []interface{}
		levels  []LevelFlag
		formats []FormatFlag
	}

	var tests = []test{
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Panic, Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Error, Panic, Debug},
			[]FormatFlag{FilePath, LineNumber},
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
		},
	}

	// Ignore panic for tests.
	defer func() { recover() }()

	for i, s := range tests {
		var buf = new(bytes.Buffer)

		l, _ := New(s.levels...)
		l.Config.Formats.Set(s.formats...)
		l.skip = 5

		l.Fpanic(buf, s.data...)

		exp := ""
		res := buf.String()
		ss := getStackSlice(testSkip)

		if ok, _ := l.Config.Levels.Has(Panic); ok {
			res = res[19:]
			exp = getPrefix(Panic, l.Config, ss) +
				fmt.Sprint(s.data...)
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}

// TestFpanicf tests Fpanicf method.
func TestFpanicf(t *testing.T) {
	type test struct {
		data    []interface{}
		levels  []LevelFlag
		formats []FormatFlag
		format  string
	}

	var tests = []test{
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Panic, Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
			"%s %s",
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Error, Panic, Debug},
			[]FormatFlag{FilePath, LineNumber},
			"%s %s",
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
			"%s %s",
		},
	}

	// Ignore panic for tests.
	defer func() { recover() }()

	for i, s := range tests {
		var buf = new(bytes.Buffer)

		l, _ := New(s.levels...)
		l.Config.Formats.Set(s.formats...)
		l.skip = 5

		l.Fpanicf(buf, s.format, s.data...)

		exp := ""
		res := buf.String()
		ss := getStackSlice(testSkip)

		if ok, _ := l.Config.Levels.Has(Panic); ok {
			res = res[19:]
			exp = getPrefix(Panic, l.Config, ss) +
				fmt.Sprintf(s.format, s.data...)
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}

// TestFpanicln tests Fpanicln method.
func TestFpanicln(t *testing.T) {
	type test struct {
		data    []interface{}
		levels  []LevelFlag
		formats []FormatFlag
	}

	var tests = []test{
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Panic, Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Error, Panic, Debug},
			[]FormatFlag{FilePath, LineNumber},
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
		},
	}

	// Ignore panic for tests.
	defer func() { recover() }()

	for i, s := range tests {
		var buf = new(bytes.Buffer)

		l, _ := New(s.levels...)
		l.Config.Formats.Set(s.formats...)
		l.skip = 5

		l.Fpanicln(buf, s.data...)

		exp := ""
		res := buf.String()
		ss := getStackSlice(testSkip)

		if ok, _ := l.Config.Levels.Has(Panic); ok {
			res = res[19:]
			exp = getPrefix(Panic, l.Config, ss) +
				fmt.Sprintln(s.data...)
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}

// TestPanic tests Panic method.
func TestPanic(t *testing.T) {
	type test struct {
		data    []interface{}
		levels  []LevelFlag
		formats []FormatFlag
	}

	var tests = []test{
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Panic, Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Error, Panic, Debug},
			[]FormatFlag{FilePath, LineNumber},
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
		},
	}

	// Ignore panic for tests.
	defer func() { recover() }()

	for i, s := range tests {
		var buf = new(bytes.Buffer)

		l, _ := New(s.levels...)
		l.Writer = buf
		l.Config.Formats.Set(s.formats...)
		l.skip = 5

		l.Panic(s.data...)

		exp := ""
		res := buf.String()
		ss := getStackSlice(testSkip)

		if ok, _ := l.Config.Levels.Has(Panic); ok {
			res = res[19:]
			exp = getPrefix(Panic, l.Config, ss) +
				fmt.Sprint(s.data...)
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}

// TestPanicf tests Panicf method.
func TestPanicf(t *testing.T) {
	type test struct {
		data    []interface{}
		levels  []LevelFlag
		formats []FormatFlag
		format  string
	}

	var tests = []test{
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Panic, Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
			"%s %s",
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Error, Panic, Debug},
			[]FormatFlag{FilePath, LineNumber},
			"%s %s",
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
			"%s %s",
		},
	}

	// Ignore panic for tests.
	defer func() { recover() }()

	for i, s := range tests {
		var buf = new(bytes.Buffer)

		l, _ := New(s.levels...)
		l.Writer = buf
		l.Config.Formats.Set(s.formats...)
		l.skip = 5

		l.Panicf(s.format, s.data...)

		exp := ""
		res := buf.String()
		ss := getStackSlice(testSkip)

		if ok, _ := l.Config.Levels.Has(Panic); ok {
			res = res[19:]
			exp = getPrefix(Panic, l.Config, ss) +
				fmt.Sprintf(s.format, s.data...)
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}

// TestPanicln tests Panicln method.
func TestPanicln(t *testing.T) {
	type test struct {
		data    []interface{}
		levels  []LevelFlag
		formats []FormatFlag
	}

	var tests = []test{
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Panic, Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Error, Panic, Debug},
			[]FormatFlag{FilePath, LineNumber},
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
		},
	}

	// Ignore panic for tests.
	defer func() { recover() }()

	for i, s := range tests {
		var buf = new(bytes.Buffer)

		l, _ := New(s.levels...)
		l.Writer = buf
		l.Config.Formats.Set(s.formats...)
		l.skip = 5

		l.Panicln(s.data...)

		exp := ""
		res := buf.String()
		ss := getStackSlice(testSkip)

		if ok, _ := l.Config.Levels.Has(Panic); ok {
			res = res[19:]
			exp = getPrefix(Panic, l.Config, ss) +
				fmt.Sprintln(s.data...)
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}

// TestFfatal tests Ffatal method.
func TestFfatal(t *testing.T) {
	type test struct {
		data    []interface{}
		levels  []LevelFlag
		formats []FormatFlag
	}

	var tests = []test{
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Fatal, Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Error, Fatal, Debug},
			[]FormatFlag{FilePath, LineNumber},
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
		},
	}

	for i, s := range tests {
		var buf = new(bytes.Buffer)

		l, _ := New(s.levels...)
		l.Config.Formats.Set(s.formats...)
		l.Config.FatalStatusCode = 0 // ignore force exit for tests
		l.skip = 5

		l.Ffatal(buf, s.data...)

		exp := ""
		res := buf.String()
		ss := getStackSlice(testSkip)

		if ok, _ := l.Config.Levels.Has(Fatal); ok {
			res = res[19:]
			exp = getPrefix(Fatal, l.Config, ss) +
				fmt.Sprint(s.data...)
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}

// TestFfatalf tests Ffatalf method.
func TestFfatalf(t *testing.T) {
	type test struct {
		data    []interface{}
		levels  []LevelFlag
		formats []FormatFlag
		format  string
	}

	var tests = []test{
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Fatal, Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
			"%s %s",
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Error, Fatal, Debug},
			[]FormatFlag{FilePath, LineNumber},
			"%s %s",
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
			"%s %s",
		},
	}

	for i, s := range tests {
		var buf = new(bytes.Buffer)

		l, _ := New(s.levels...)
		l.Config.Formats.Set(s.formats...)
		l.Config.FatalStatusCode = 0 // ignore force exit for tests
		l.skip = 5

		l.Ffatalf(buf, s.format, s.data...)

		exp := ""
		res := buf.String()
		ss := getStackSlice(testSkip)

		if ok, _ := l.Config.Levels.Has(Fatal); ok {
			res = res[19:]
			exp = getPrefix(Fatal, l.Config, ss) +
				fmt.Sprintf(s.format, s.data...)
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}

// TestFfatalln tests Ffatalln method.
func TestFfatalln(t *testing.T) {
	type test struct {
		data    []interface{}
		levels  []LevelFlag
		formats []FormatFlag
	}

	var tests = []test{
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Fatal, Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Error, Fatal, Debug},
			[]FormatFlag{FilePath, LineNumber},
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
		},
	}

	for i, s := range tests {
		var buf = new(bytes.Buffer)

		l, _ := New(s.levels...)
		l.Config.Formats.Set(s.formats...)
		l.Config.FatalStatusCode = 0 // ignore force exit for tests
		l.skip = 5

		l.Ffatalln(buf, s.data...)

		exp := ""
		res := buf.String()
		ss := getStackSlice(testSkip)

		if ok, _ := l.Config.Levels.Has(Fatal); ok {
			res = res[19:]
			exp = getPrefix(Fatal, l.Config, ss) +
				fmt.Sprintln(s.data...)
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}

// TestFatal tests Fatal method.
func TestFatal(t *testing.T) {
	type test struct {
		data    []interface{}
		levels  []LevelFlag
		formats []FormatFlag
	}

	var tests = []test{
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Fatal, Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Error, Fatal, Debug},
			[]FormatFlag{FilePath, LineNumber},
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
		},
	}

	for i, s := range tests {
		var buf = new(bytes.Buffer)

		l, _ := New(s.levels...)
		l.Writer = buf
		l.Config.Formats.Set(s.formats...)
		l.Config.FatalStatusCode = 0 // ignore force exit for tests
		l.skip = 5

		l.Fatal(s.data...)

		exp := ""
		res := buf.String()
		ss := getStackSlice(testSkip)

		if ok, _ := l.Config.Levels.Has(Fatal); ok {
			res = res[19:]
			exp = getPrefix(Fatal, l.Config, ss) +
				fmt.Sprint(s.data...)
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}

// TestFatalf tests Fatalf method.
func TestFatalf(t *testing.T) {
	type test struct {
		data    []interface{}
		levels  []LevelFlag
		formats []FormatFlag
		format  string
	}

	var tests = []test{
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Fatal, Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
			"%s %s",
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Error, Fatal, Debug},
			[]FormatFlag{FilePath, LineNumber},
			"%s %s",
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
			"%s %s",
		},
	}

	for i, s := range tests {
		var buf = new(bytes.Buffer)

		l, _ := New(s.levels...)
		l.Writer = buf
		l.Config.Formats.Set(s.formats...)
		l.Config.FatalStatusCode = 0 // ignore force exit for tests
		l.skip = 5

		l.Fatalf(s.format, s.data...)

		exp := ""
		res := buf.String()
		ss := getStackSlice(testSkip)

		if ok, _ := l.Config.Levels.Has(Fatal); ok {
			res = res[19:]
			exp = getPrefix(Fatal, l.Config, ss) +
				fmt.Sprintf(s.format, s.data...)
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}

// TestFatalln tests Fatalln method.
func TestFatalln(t *testing.T) {
	type test struct {
		data    []interface{}
		levels  []LevelFlag
		formats []FormatFlag
	}

	var tests = []test{
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Fatal, Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Error, Fatal, Debug},
			[]FormatFlag{FilePath, LineNumber},
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
		},
	}

	for i, s := range tests {
		var buf = new(bytes.Buffer)

		l, _ := New(s.levels...)
		l.Writer = buf
		l.Config.Formats.Set(s.formats...)
		l.Config.FatalStatusCode = 0 // ignore force exit for tests
		l.skip = 5

		l.Fatalln(s.data...)

		exp := ""
		res := buf.String()
		ss := getStackSlice(testSkip)

		if ok, _ := l.Config.Levels.Has(Fatal); ok {
			res = res[19:]
			exp = getPrefix(Fatal, l.Config, ss) +
				fmt.Sprintln(s.data...)
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}

// TestFerror tests Ferror method.
func TestFerror(t *testing.T) {
	type test struct {
		data    []interface{}
		levels  []LevelFlag
		formats []FormatFlag
	}

	var tests = []test{
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Error, Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Error, Error, Debug},
			[]FormatFlag{FilePath, LineNumber},
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
		},
	}

	for i, s := range tests {
		var buf = new(bytes.Buffer)

		l, _ := New(s.levels...)
		l.Config.Formats.Set(s.formats...)
		l.skip = 5

		l.Ferror(buf, s.data...)

		exp := ""
		res := buf.String()
		ss := getStackSlice(testSkip)

		if ok, _ := l.Config.Levels.Has(Error); ok {
			res = res[19:]
			exp = getPrefix(Error, l.Config, ss) +
				fmt.Sprint(s.data...)
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}

// TestFerrorf tests Ferrorf method.
func TestFerrorf(t *testing.T) {
	type test struct {
		data    []interface{}
		levels  []LevelFlag
		formats []FormatFlag
		format  string
	}

	var tests = []test{
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Error, Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
			"%s %s",
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Error, Error, Debug},
			[]FormatFlag{FilePath, LineNumber},
			"%s %s",
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
			"%s %s",
		},
	}

	for i, s := range tests {
		var buf = new(bytes.Buffer)

		l, _ := New(s.levels...)
		l.Config.Formats.Set(s.formats...)
		l.skip = 5

		l.Ferrorf(buf, s.format, s.data...)

		exp := ""
		res := buf.String()
		ss := getStackSlice(testSkip)

		if ok, _ := l.Config.Levels.Has(Error); ok {
			res = res[19:]
			exp = getPrefix(Error, l.Config, ss) +
				fmt.Sprintf(s.format, s.data...)
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}

// TestFerrorln tests Ferrorln method.
func TestFerrorln(t *testing.T) {
	type test struct {
		data    []interface{}
		levels  []LevelFlag
		formats []FormatFlag
	}

	var tests = []test{
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Error, Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Error, Error, Debug},
			[]FormatFlag{FilePath, LineNumber},
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
		},
	}

	for i, s := range tests {
		var buf = new(bytes.Buffer)

		l, _ := New(s.levels...)
		l.Config.Formats.Set(s.formats...)
		l.skip = 5

		l.Ferrorln(buf, s.data...)

		exp := ""
		res := buf.String()
		ss := getStackSlice(testSkip)

		if ok, _ := l.Config.Levels.Has(Error); ok {
			res = res[19:]
			exp = getPrefix(Error, l.Config, ss) +
				fmt.Sprintln(s.data...)
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}

// TestError tests Error method.
func TestError(t *testing.T) {
	type test struct {
		data    []interface{}
		levels  []LevelFlag
		formats []FormatFlag
	}

	var tests = []test{
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Error, Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Error, Error, Debug},
			[]FormatFlag{FilePath, LineNumber},
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
		},
	}

	for i, s := range tests {
		var buf = new(bytes.Buffer)

		l, _ := New(s.levels...)
		l.Writer = buf
		l.Config.Formats.Set(s.formats...)
		l.skip = 5

		l.Error(s.data...)

		exp := ""
		res := buf.String()
		ss := getStackSlice(testSkip)

		if ok, _ := l.Config.Levels.Has(Error); ok {
			res = res[19:]
			exp = getPrefix(Error, l.Config, ss) +
				fmt.Sprint(s.data...)
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}

// TestErrorf tests Errorf method.
func TestErrorf(t *testing.T) {
	type test struct {
		data    []interface{}
		levels  []LevelFlag
		formats []FormatFlag
		format  string
	}

	var tests = []test{
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Error, Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
			"%s %s",
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Error, Error, Debug},
			[]FormatFlag{FilePath, LineNumber},
			"%s %s",
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
			"%s %s",
		},
	}

	for i, s := range tests {
		var buf = new(bytes.Buffer)

		l, _ := New(s.levels...)
		l.Writer = buf
		l.Config.Formats.Set(s.formats...)
		l.skip = 5

		l.Errorf(s.format, s.data...)

		exp := ""
		res := buf.String()
		ss := getStackSlice(testSkip)

		if ok, _ := l.Config.Levels.Has(Error); ok {
			res = res[19:]
			exp = getPrefix(Error, l.Config, ss) +
				fmt.Sprintf(s.format, s.data...)
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}

// TestErrorln tests Errorln method.
func TestErrorln(t *testing.T) {
	type test struct {
		data    []interface{}
		levels  []LevelFlag
		formats []FormatFlag
	}

	var tests = []test{
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Error, Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Error, Error, Debug},
			[]FormatFlag{FilePath, LineNumber},
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
		},
	}

	for i, s := range tests {
		var buf = new(bytes.Buffer)

		l, _ := New(s.levels...)
		l.Writer = buf
		l.Config.Formats.Set(s.formats...)
		l.skip = 5

		l.Errorln(s.data...)

		exp := ""
		res := buf.String()
		ss := getStackSlice(testSkip)

		if ok, _ := l.Config.Levels.Has(Error); ok {
			res = res[19:]
			exp = getPrefix(Error, l.Config, ss) +
				fmt.Sprintln(s.data...)
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}

// TestFwarn tests Fwarn method.
func TestFwarn(t *testing.T) {
	type test struct {
		data    []interface{}
		levels  []LevelFlag
		formats []FormatFlag
	}

	var tests = []test{
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Warn, Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Error, Warn, Debug},
			[]FormatFlag{FilePath, LineNumber},
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
		},
	}

	for i, s := range tests {
		var buf = new(bytes.Buffer)

		l, _ := New(s.levels...)
		l.Config.Formats.Set(s.formats...)
		l.skip = 5

		l.Fwarn(buf, s.data...)

		exp := ""
		res := buf.String()
		ss := getStackSlice(testSkip)

		if ok, _ := l.Config.Levels.Has(Warn); ok {
			res = res[19:]
			exp = getPrefix(Warn, l.Config, ss) +
				fmt.Sprint(s.data...)
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}

// TestFwarnf tests Fwarnf method.
func TestFwarnf(t *testing.T) {
	type test struct {
		data    []interface{}
		levels  []LevelFlag
		formats []FormatFlag
		format  string
	}

	var tests = []test{
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Warn, Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
			"%s %s",
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Error, Warn, Debug},
			[]FormatFlag{FilePath, LineNumber},
			"%s %s",
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
			"%s %s",
		},
	}

	for i, s := range tests {
		var buf = new(bytes.Buffer)

		l, _ := New(s.levels...)
		l.Config.Formats.Set(s.formats...)
		l.skip = 5

		l.Fwarnf(buf, s.format, s.data...)

		exp := ""
		res := buf.String()
		ss := getStackSlice(testSkip)

		if ok, _ := l.Config.Levels.Has(Warn); ok {
			res = res[19:]
			exp = getPrefix(Warn, l.Config, ss) +
				fmt.Sprintf(s.format, s.data...)
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}

// TestFwarnln tests Fwarnln method.
func TestFwarnln(t *testing.T) {
	type test struct {
		data    []interface{}
		levels  []LevelFlag
		formats []FormatFlag
	}

	var tests = []test{
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Warn, Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Error, Warn, Debug},
			[]FormatFlag{FilePath, LineNumber},
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
		},
	}

	for i, s := range tests {
		var buf = new(bytes.Buffer)

		l, _ := New(s.levels...)
		l.Config.Formats.Set(s.formats...)
		l.skip = 5

		l.Fwarnln(buf, s.data...)

		exp := ""
		res := buf.String()
		ss := getStackSlice(testSkip)

		if ok, _ := l.Config.Levels.Has(Warn); ok {
			res = res[19:]
			exp = getPrefix(Warn, l.Config, ss) +
				fmt.Sprintln(s.data...)
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}

// TestWarn tests Warn method.
func TestWarn(t *testing.T) {
	type test struct {
		data    []interface{}
		levels  []LevelFlag
		formats []FormatFlag
	}

	var tests = []test{
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Warn, Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Error, Warn, Debug},
			[]FormatFlag{FilePath, LineNumber},
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
		},
	}

	for i, s := range tests {
		var buf = new(bytes.Buffer)

		l, _ := New(s.levels...)
		l.Writer = buf
		l.Config.Formats.Set(s.formats...)
		l.skip = 5

		l.Warn(s.data...)

		exp := ""
		res := buf.String()
		ss := getStackSlice(testSkip)

		if ok, _ := l.Config.Levels.Has(Warn); ok {
			res = res[19:]
			exp = getPrefix(Warn, l.Config, ss) +
				fmt.Sprint(s.data...)
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}

// TestWarnf tests Warnf method.
func TestWarnf(t *testing.T) {
	type test struct {
		data    []interface{}
		levels  []LevelFlag
		formats []FormatFlag
		format  string
	}

	var tests = []test{
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Warn, Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
			"%s %s",
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Error, Warn, Debug},
			[]FormatFlag{FilePath, LineNumber},
			"%s %s",
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
			"%s %s",
		},
	}

	for i, s := range tests {
		var buf = new(bytes.Buffer)

		l, _ := New(s.levels...)
		l.Writer = buf
		l.Config.Formats.Set(s.formats...)
		l.skip = 5

		l.Warnf(s.format, s.data...)

		exp := ""
		res := buf.String()
		ss := getStackSlice(testSkip)

		if ok, _ := l.Config.Levels.Has(Warn); ok {
			res = res[19:]
			exp = getPrefix(Warn, l.Config, ss) +
				fmt.Sprintf(s.format, s.data...)
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}

// TestWarnln tests Warnln method.
func TestWarnln(t *testing.T) {
	type test struct {
		data    []interface{}
		levels  []LevelFlag
		formats []FormatFlag
	}

	var tests = []test{
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Warn, Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Error, Warn, Debug},
			[]FormatFlag{FilePath, LineNumber},
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
		},
	}

	for i, s := range tests {
		var buf = new(bytes.Buffer)

		l, _ := New(s.levels...)
		l.Writer = buf
		l.Config.Formats.Set(s.formats...)
		l.skip = 5

		l.Warnln(s.data...)

		exp := ""
		res := buf.String()
		ss := getStackSlice(testSkip)

		if ok, _ := l.Config.Levels.Has(Warn); ok {
			res = res[19:]
			exp = getPrefix(Warn, l.Config, ss) +
				fmt.Sprintln(s.data...)
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}

// TestFinfo tests Finfo method.
func TestFinfo(t *testing.T) {
	type test struct {
		data    []interface{}
		levels  []LevelFlag
		formats []FormatFlag
	}

	var tests = []test{
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Info, Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Error, Info, Debug},
			[]FormatFlag{FilePath, LineNumber},
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
		},
	}

	for i, s := range tests {
		var buf = new(bytes.Buffer)

		l, _ := New(s.levels...)
		l.Config.Formats.Set(s.formats...)
		l.skip = 5

		l.Finfo(buf, s.data...)

		exp := ""
		res := buf.String()
		ss := getStackSlice(testSkip)

		if ok, _ := l.Config.Levels.Has(Info); ok {
			res = res[19:]
			exp = getPrefix(Info, l.Config, ss) +
				fmt.Sprint(s.data...)
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}

// TestFinfof tests Finfof method.
func TestFinfof(t *testing.T) {
	type test struct {
		data    []interface{}
		levels  []LevelFlag
		formats []FormatFlag
		format  string
	}

	var tests = []test{
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Info, Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
			"%s %s",
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Error, Info, Debug},
			[]FormatFlag{FilePath, LineNumber},
			"%s %s",
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
			"%s %s",
		},
	}

	for i, s := range tests {
		var buf = new(bytes.Buffer)

		l, _ := New(s.levels...)
		l.Config.Formats.Set(s.formats...)
		l.skip = 5

		l.Finfof(buf, s.format, s.data...)

		exp := ""
		res := buf.String()
		ss := getStackSlice(testSkip)

		if ok, _ := l.Config.Levels.Has(Info); ok {
			res = res[19:]
			exp = getPrefix(Info, l.Config, ss) +
				fmt.Sprintf(s.format, s.data...)
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}

// TestFinfoln tests Finfoln method.
func TestFinfoln(t *testing.T) {
	type test struct {
		data    []interface{}
		levels  []LevelFlag
		formats []FormatFlag
	}

	var tests = []test{
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Info, Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Error, Info, Debug},
			[]FormatFlag{FilePath, LineNumber},
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
		},
	}

	for i, s := range tests {
		var buf = new(bytes.Buffer)

		l, _ := New(s.levels...)
		l.Config.Formats.Set(s.formats...)
		l.skip = 5

		l.Finfoln(buf, s.data...)

		exp := ""
		res := buf.String()
		ss := getStackSlice(testSkip)

		if ok, _ := l.Config.Levels.Has(Info); ok {
			res = res[19:]
			exp = getPrefix(Info, l.Config, ss) +
				fmt.Sprintln(s.data...)
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}

// TestInfo tests Info method.
func TestInfo(t *testing.T) {
	type test struct {
		data    []interface{}
		levels  []LevelFlag
		formats []FormatFlag
	}

	var tests = []test{
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Info, Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Error, Info, Debug},
			[]FormatFlag{FilePath, LineNumber},
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
		},
	}

	for i, s := range tests {
		var buf = new(bytes.Buffer)

		l, _ := New(s.levels...)
		l.Writer = buf
		l.Config.Formats.Set(s.formats...)
		l.skip = 5

		l.Info(s.data...)

		exp := ""
		res := buf.String()
		ss := getStackSlice(testSkip)

		if ok, _ := l.Config.Levels.Has(Info); ok {
			res = res[19:]
			exp = getPrefix(Info, l.Config, ss) +
				fmt.Sprint(s.data...)
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}

// TestInfof tests Infof method.
func TestInfof(t *testing.T) {
	type test struct {
		data    []interface{}
		levels  []LevelFlag
		formats []FormatFlag
		format  string
	}

	var tests = []test{
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Info, Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
			"%s %s",
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Error, Info, Debug},
			[]FormatFlag{FilePath, LineNumber},
			"%s %s",
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
			"%s %s",
		},
	}

	for i, s := range tests {
		var buf = new(bytes.Buffer)

		l, _ := New(s.levels...)
		l.Writer = buf
		l.Config.Formats.Set(s.formats...)
		l.skip = 5

		l.Infof(s.format, s.data...)

		exp := ""
		res := buf.String()
		ss := getStackSlice(testSkip)

		if ok, _ := l.Config.Levels.Has(Info); ok {
			res = res[19:]
			exp = getPrefix(Info, l.Config, ss) +
				fmt.Sprintf(s.format, s.data...)
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}

// TestInfoln tests Infoln method.
func TestInfoln(t *testing.T) {
	type test struct {
		data    []interface{}
		levels  []LevelFlag
		formats []FormatFlag
	}

	var tests = []test{
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Info, Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Error, Info, Debug},
			[]FormatFlag{FilePath, LineNumber},
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
		},
	}

	for i, s := range tests {
		var buf = new(bytes.Buffer)

		l, _ := New(s.levels...)
		l.Writer = buf
		l.Config.Formats.Set(s.formats...)
		l.skip = 5

		l.Infoln(s.data...)

		exp := ""
		res := buf.String()
		ss := getStackSlice(testSkip)

		if ok, _ := l.Config.Levels.Has(Info); ok {
			res = res[19:]
			exp = getPrefix(Info, l.Config, ss) +
				fmt.Sprintln(s.data...)
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}

// TestFdebug tests Fdebug method.
func TestFdebug(t *testing.T) {
	type test struct {
		data    []interface{}
		levels  []LevelFlag
		formats []FormatFlag
	}

	var tests = []test{
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Debug, Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Error, Debug, Debug},
			[]FormatFlag{FilePath, LineNumber},
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
		},
	}

	for i, s := range tests {
		var buf = new(bytes.Buffer)

		l, _ := New(s.levels...)
		l.Config.Formats.Set(s.formats...)
		l.skip = 5

		l.Fdebug(buf, s.data...)

		exp := ""
		res := buf.String()
		ss := getStackSlice(testSkip)

		if ok, _ := l.Config.Levels.Has(Debug); ok {
			res = res[19:]
			exp = getPrefix(Debug, l.Config, ss) +
				fmt.Sprint(s.data...)
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}

// TestFdebugf tests Fdebugf method.
func TestFdebugf(t *testing.T) {
	type test struct {
		data    []interface{}
		levels  []LevelFlag
		formats []FormatFlag
		format  string
	}

	var tests = []test{
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Debug, Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
			"%s %s",
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Error, Debug, Debug},
			[]FormatFlag{FilePath, LineNumber},
			"%s %s",
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
			"%s %s",
		},
	}

	for i, s := range tests {
		var buf = new(bytes.Buffer)

		l, _ := New(s.levels...)
		l.Config.Formats.Set(s.formats...)
		l.skip = 5

		l.Fdebugf(buf, s.format, s.data...)

		exp := ""
		res := buf.String()
		ss := getStackSlice(testSkip)

		if ok, _ := l.Config.Levels.Has(Debug); ok {
			res = res[19:]
			exp = getPrefix(Debug, l.Config, ss) +
				fmt.Sprintf(s.format, s.data...)
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}

// TestFdebugln tests Fdebugln method.
func TestFdebugln(t *testing.T) {
	type test struct {
		data    []interface{}
		levels  []LevelFlag
		formats []FormatFlag
	}

	var tests = []test{
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Debug, Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Error, Debug, Debug},
			[]FormatFlag{FilePath, LineNumber},
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
		},
	}

	for i, s := range tests {
		var buf = new(bytes.Buffer)

		l, _ := New(s.levels...)
		l.Config.Formats.Set(s.formats...)
		l.skip = 5

		l.Fdebugln(buf, s.data...)

		exp := ""
		res := buf.String()
		ss := getStackSlice(testSkip)

		if ok, _ := l.Config.Levels.Has(Debug); ok {
			res = res[19:]
			exp = getPrefix(Debug, l.Config, ss) +
				fmt.Sprintln(s.data...)
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}

// TestDebug tests Debug method.
func TestDebug(t *testing.T) {
	type test struct {
		data    []interface{}
		levels  []LevelFlag
		formats []FormatFlag
	}

	var tests = []test{
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Debug, Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Error, Debug, Debug},
			[]FormatFlag{FilePath, LineNumber},
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
		},
	}

	for i, s := range tests {
		var buf = new(bytes.Buffer)

		l, _ := New(s.levels...)
		l.Writer = buf
		l.Config.Formats.Set(s.formats...)
		l.skip = 5

		l.Debug(s.data...)

		exp := ""
		res := buf.String()
		ss := getStackSlice(testSkip)

		if ok, _ := l.Config.Levels.Has(Debug); ok {
			res = res[19:]
			exp = getPrefix(Debug, l.Config, ss) +
				fmt.Sprint(s.data...)
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}

// TestDebugf tests Debugf method.
func TestDebugf(t *testing.T) {
	type test struct {
		data    []interface{}
		levels  []LevelFlag
		formats []FormatFlag
		format  string
	}

	var tests = []test{
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Debug, Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
			"%s %s",
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Error, Debug, Debug},
			[]FormatFlag{FilePath, LineNumber},
			"%s %s",
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
			"%s %s",
		},
	}

	for i, s := range tests {
		var buf = new(bytes.Buffer)

		l, _ := New(s.levels...)
		l.Writer = buf
		l.Config.Formats.Set(s.formats...)
		l.skip = 5

		l.Debugf(s.format, s.data...)

		exp := ""
		res := buf.String()
		ss := getStackSlice(testSkip)

		if ok, _ := l.Config.Levels.Has(Debug); ok {
			res = res[19:]
			exp = getPrefix(Debug, l.Config, ss) +
				fmt.Sprintf(s.format, s.data...)
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}

// TestDebugln tests Debugln method.
func TestDebugln(t *testing.T) {
	type test struct {
		data    []interface{}
		levels  []LevelFlag
		formats []FormatFlag
	}

	var tests = []test{
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Debug, Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Error, Debug, Debug},
			[]FormatFlag{FilePath, LineNumber},
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
		},
	}

	for i, s := range tests {
		var buf = new(bytes.Buffer)

		l, _ := New(s.levels...)
		l.Writer = buf
		l.Config.Formats.Set(s.formats...)
		l.skip = 5

		l.Debugln(s.data...)

		exp := ""
		res := buf.String()
		ss := getStackSlice(testSkip)

		if ok, _ := l.Config.Levels.Has(Debug); ok {
			res = res[19:]
			exp = getPrefix(Debug, l.Config, ss) +
				fmt.Sprintln(s.data...)
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}

// TestFtrace tests Ftrace method.
func TestFtrace(t *testing.T) {
	type test struct {
		data    []interface{}
		levels  []LevelFlag
		formats []FormatFlag
	}

	var tests = []test{
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Trace, Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Error, Trace, Debug},
			[]FormatFlag{FilePath, LineNumber},
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
		},
	}

	for i, s := range tests {
		var buf = new(bytes.Buffer)

		l, _ := New(s.levels...)
		l.Config.Formats.Set(s.formats...)
		l.skip = 5

		l.Ftrace(buf, s.data...)

		exp := ""
		res := buf.String()
		ss := getStackSlice(testSkip)

		if ok, _ := l.Config.Levels.Has(Trace); ok {
			res = res[19:]
			exp = getPrefix(Trace, l.Config, ss) +
				fmt.Sprint(s.data...)
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}

// TestFtracef tests Ftracef method.
func TestFtracef(t *testing.T) {
	type test struct {
		data    []interface{}
		levels  []LevelFlag
		formats []FormatFlag
		format  string
	}

	var tests = []test{
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Trace, Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
			"%s %s",
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Error, Trace, Debug},
			[]FormatFlag{FilePath, LineNumber},
			"%s %s",
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
			"%s %s",
		},
	}

	for i, s := range tests {
		var buf = new(bytes.Buffer)

		l, _ := New(s.levels...)
		l.Config.Formats.Set(s.formats...)
		l.skip = 5

		l.Ftracef(buf, s.format, s.data...)

		exp := ""
		res := buf.String()
		ss := getStackSlice(testSkip)

		if ok, _ := l.Config.Levels.Has(Trace); ok {
			res = res[19:]
			exp = getPrefix(Trace, l.Config, ss) +
				fmt.Sprintf(s.format, s.data...)
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}

// TestFtraceln tests Ftraceln method.
func TestFtraceln(t *testing.T) {
	type test struct {
		data    []interface{}
		levels  []LevelFlag
		formats []FormatFlag
	}

	var tests = []test{
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Trace, Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Error, Trace, Debug},
			[]FormatFlag{FilePath, LineNumber},
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
		},
	}

	for i, s := range tests {
		var buf = new(bytes.Buffer)

		l, _ := New(s.levels...)
		l.Config.Formats.Set(s.formats...)
		l.skip = 5

		l.Ftraceln(buf, s.data...)

		exp := ""
		res := buf.String()
		ss := getStackSlice(testSkip)

		if ok, _ := l.Config.Levels.Has(Trace); ok {
			res = res[19:]
			exp = getPrefix(Trace, l.Config, ss) +
				fmt.Sprintln(s.data...)
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}

// TestTrace tests Trace method.
func TestTrace(t *testing.T) {
	type test struct {
		data    []interface{}
		levels  []LevelFlag
		formats []FormatFlag
	}

	var tests = []test{
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Trace, Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Error, Trace, Debug},
			[]FormatFlag{FilePath, LineNumber},
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
		},
	}

	for i, s := range tests {
		var buf = new(bytes.Buffer)

		l, _ := New(s.levels...)
		l.Writer = buf
		l.Config.Formats.Set(s.formats...)
		l.skip = 5

		l.Trace(s.data...)

		exp := ""
		res := buf.String()
		ss := getStackSlice(testSkip)

		if ok, _ := l.Config.Levels.Has(Trace); ok {
			res = res[19:]
			exp = getPrefix(Trace, l.Config, ss) +
				fmt.Sprint(s.data...)
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}

// TestTracef tests Tracef method.
func TestTracef(t *testing.T) {
	type test struct {
		data    []interface{}
		levels  []LevelFlag
		formats []FormatFlag
		format  string
	}

	var tests = []test{
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Trace, Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
			"%s %s",
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Error, Trace, Debug},
			[]FormatFlag{FilePath, LineNumber},
			"%s %s",
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
			"%s %s",
		},
	}

	for i, s := range tests {
		var buf = new(bytes.Buffer)

		l, _ := New(s.levels...)
		l.Writer = buf
		l.Config.Formats.Set(s.formats...)
		l.skip = 5

		l.Tracef(s.format, s.data...)

		exp := ""
		res := buf.String()
		ss := getStackSlice(testSkip)

		if ok, _ := l.Config.Levels.Has(Trace); ok {
			res = res[19:]
			exp = getPrefix(Trace, l.Config, ss) +
				fmt.Sprintf(s.format, s.data...)
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}

// TestTraceln tests Traceln method.
func TestTraceln(t *testing.T) {
	type test struct {
		data    []interface{}
		levels  []LevelFlag
		formats []FormatFlag
	}

	var tests = []test{
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Trace, Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Error, Trace, Debug},
			[]FormatFlag{FilePath, LineNumber},
		},
		{
			[]interface{}{"1", "2"},
			[]LevelFlag{Debug},
			[]FormatFlag{FilePath, FuncName, LineNumber},
		},
	}

	for i, s := range tests {
		var buf = new(bytes.Buffer)

		l, _ := New(s.levels...)
		l.Writer = buf
		l.Config.Formats.Set(s.formats...)
		l.skip = 5

		l.Traceln(s.data...)

		exp := ""
		res := buf.String()
		ss := getStackSlice(testSkip)

		if ok, _ := l.Config.Levels.Has(Trace); ok {
			res = res[19:]
			exp = getPrefix(Trace, l.Config, ss) +
				fmt.Sprintln(s.data...)
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}
