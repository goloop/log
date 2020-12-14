package log

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

// TestEcho tests echo method.
func TestEcho(t *testing.T) {
	type test struct {
		level Level
		data  []interface{}

		showFilePath bool
		showFuncName bool
		showFileLine bool
	}

	var tests = []test{
		{DEBUG, []interface{}{"1 ", "2"}, true, true, true},
		{ERROR, []interface{}{"1 ", "2"}, true, false, false},
		{FATAL, []interface{}{"1 ", "2"}, false, false, true},
		{INFO, []interface{}{"1 ", "2"}, false, false, false},
	}

	trace := getTrace(3)
	l, _ := New(DEBUG, ERROR, FATAL)
	for i, s := range tests {
		var buf = new(bytes.Buffer)
		l.Format(s.showFilePath, s.showFuncName, s.showFileLine)
		l.echo(4, buf, s.level, s.data...)

		res := buf.String()
		exp := getPrefix(trace, "", l.Timestamp, s.level, s.showFilePath,
			s.showFuncName, s.showFileLine) + fmt.Sprint(s.data...)

		if l.Levels.All(s.level) {
			res = res[20:]
			exp = exp[20:]
		} else {
			exp = "" // if level not supportde
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}

// TestEchof tests echof method.
func TestEchof(t *testing.T) {
	type test struct {
		level  Level
		format string
		data   []interface{}

		showFilePath bool
		showFuncName bool
		showFileLine bool
	}

	var tests = []test{
		{DEBUG, "%s-%s", []interface{}{"1", "2"}, true, true, true},
		{ERROR, "%s:%s", []interface{}{"1", "2"}, true, false, false},
		{FATAL, "%s=%s", []interface{}{"1", "2"}, false, false, true},
		{INFO, "%s=%s", []interface{}{"1", "2"}, false, false, true},
	}

	trace := getTrace(3)
	l, _ := New(DEBUG, ERROR, FATAL)
	for i, s := range tests {
		var buf = new(bytes.Buffer)
		l.Format(s.showFilePath, s.showFuncName, s.showFileLine)
		l.echof(4, buf, s.level, s.format, s.data...)

		res := buf.String()
		prefix := getPrefix(trace, s.format, l.Timestamp, s.level,
			s.showFilePath, s.showFuncName, s.showFileLine)
		exp := fmt.Sprintf(prefix, s.data...)

		if l.Levels.All(s.level) {
			res = res[20:]
			exp = exp[20:]
		} else {
			exp = "" // if level not supportde
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}

// TestEcholn tests echoln method.
func TestEcholn(t *testing.T) {
	type test struct {
		level Level
		data  []interface{}

		showFilePath bool
		showFuncName bool
		showFileLine bool
	}

	var tests = []test{
		{DEBUG, []interface{}{"1", "2"}, true, true, true},
		{ERROR, []interface{}{"1", "2"}, true, false, false},
		{FATAL, []interface{}{"1", "2"}, false, false, true},
	}

	trace := getTrace(3)
	l, _ := New()
	for i, s := range tests {
		var buf = new(bytes.Buffer)
		l.Format(s.showFilePath, s.showFuncName, s.showFileLine)
		l.echoln(4, buf, s.level, s.data...)

		res := buf.String()
		exp := getPrefix(trace, "", l.Timestamp, s.level, s.showFilePath,
			s.showFuncName, s.showFileLine) + fmt.Sprintln(s.data...)

		if l.Levels.All(s.level) {
			res = res[20:]
			exp = exp[20:]
		} else {
			exp = "" // if level not supportde
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}

// TestFfatal tests Ffatal method.
func TestFfatal(t *testing.T) {
	type test struct {
		data []interface{}

		showFilePath bool
		showFuncName bool
		showFileLine bool
	}

	var tests = []test{
		{[]interface{}{"1 ", "2"}, true, true, true},
		{[]interface{}{"1 ", "2"}, true, false, false},
		{[]interface{}{"1 ", "2"}, false, false, true},
	}

	trace := getTrace(3)
	l, _ := New()
	for i, s := range tests {
		var buf = new(bytes.Buffer)
		l.Format(s.showFilePath, s.showFuncName, s.showFileLine)
		l.FatalStatusCode = 0
		l.Ffatal(buf, s.data...)

		res := buf.String()
		exp := getPrefix(trace, "", l.Timestamp, FATAL, s.showFilePath,
			s.showFuncName, s.showFileLine) + fmt.Sprint(s.data...)

		if l.Levels.All(FATAL) {
			res = res[20:]
			exp = exp[20:]
		} else {
			exp = "" // if level not supportde
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}

// TestFfatalf tests Ffatalf method.
func TestFfatalf(t *testing.T) {
	type test struct {
		format string
		data   []interface{}

		showFilePath bool
		showFuncName bool
		showFileLine bool
	}

	var tests = []test{
		{"%s-%s", []interface{}{"1", "2"}, true, true, true},
		{"%s %s", []interface{}{"1", "2"}, true, false, false},
		{"%s+%s", []interface{}{"1", "2"}, false, false, true},
	}

	trace := getTrace(3)
	l, _ := New()
	for i, s := range tests {
		var buf = new(bytes.Buffer)
		l.FatalStatusCode = 0
		l.Format(s.showFilePath, s.showFuncName, s.showFileLine)
		l.Ffatalf(buf, s.format, s.data...)

		res := buf.String()
		prefix := getPrefix(trace, s.format, l.Timestamp, FATAL,
			s.showFilePath, s.showFuncName, s.showFileLine)
		exp := fmt.Sprintf(prefix, s.data...)

		if l.Levels.All(FATAL) {
			res = res[20:]
			exp = exp[20:]
		} else {
			exp = "" // if level not supportde
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}

// TestFfatalln tests Ffatalln method.
func TestFfatalln(t *testing.T) {
	type test struct {
		data []interface{}

		showFilePath bool
		showFuncName bool
		showFileLine bool
	}

	var tests = []test{
		{[]interface{}{"1", "2"}, true, true, true},
		{[]interface{}{"1", "2"}, true, false, false},
		{[]interface{}{"1", "2"}, false, false, true},
	}

	trace := getTrace(3)
	l, _ := New()
	for i, s := range tests {
		var buf = new(bytes.Buffer)
		l.FatalStatusCode = 0
		l.Format(s.showFilePath, s.showFuncName, s.showFileLine)
		l.Ffatalln(buf, s.data...)

		res := buf.String()
		exp := getPrefix(trace, "", l.Timestamp, FATAL, s.showFilePath,
			s.showFuncName, s.showFileLine) + fmt.Sprintln(s.data...)

		if l.Levels.All(FATAL) {
			res = res[20:]
			exp = exp[20:]
		} else {
			exp = "" // if level not supportde
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}

// TestFatal tests Fatal method.
func TestFatal(t *testing.T) {
	type test struct {
		data []interface{}

		showFilePath bool
		showFuncName bool
		showFileLine bool
	}

	var tests = []test{
		{[]interface{}{"1 ", "2"}, true, true, true},
		{[]interface{}{"1 ", "2"}, true, false, false},
		{[]interface{}{"1 ", "2"}, false, false, true},
	}

	trace := getTrace(3)
	l, _ := New()
	for i, s := range tests {
		var buf = new(bytes.Buffer)
		l.Writer = buf
		l.FatalStatusCode = 0
		l.Format(s.showFilePath, s.showFuncName, s.showFileLine)
		l.Fatal(s.data...)

		res := buf.String()
		exp := getPrefix(trace, "", l.Timestamp, FATAL, s.showFilePath,
			s.showFuncName, s.showFileLine) + fmt.Sprint(s.data...)

		if l.Levels.All(FATAL) {
			res = res[20:]
			exp = exp[20:]
		} else {
			exp = "" // if level not supportde
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}

// TestFatalf tests Fatalf method.
func TestFatalf(t *testing.T) {
	type test struct {
		format string
		data   []interface{}

		showFilePath bool
		showFuncName bool
		showFileLine bool
	}

	var tests = []test{
		{"%s-%s", []interface{}{"1", "2"}, true, true, true},
		{"%s %s", []interface{}{"1", "2"}, true, false, false},
		{"%s+%s", []interface{}{"1", "2"}, false, false, true},
	}

	trace := getTrace(3)
	l, _ := New()
	for i, s := range tests {
		var buf = new(bytes.Buffer)
		l.Writer = buf
		l.FatalStatusCode = 0
		l.Format(s.showFilePath, s.showFuncName, s.showFileLine)
		l.Fatalf(s.format, s.data...)

		res := buf.String()
		prefix := getPrefix(trace, s.format, l.Timestamp, FATAL,
			s.showFilePath, s.showFuncName, s.showFileLine)
		exp := fmt.Sprintf(prefix, s.data...)

		if l.Levels.All(FATAL) {
			res = res[20:]
			exp = exp[20:]
		} else {
			exp = "" // if level not supportde
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}

// TestFatalln tests Fatalln method.
func TestFatalln(t *testing.T) {
	type test struct {
		data []interface{}

		showFilePath bool
		showFuncName bool
		showFileLine bool
	}

	var tests = []test{
		{[]interface{}{"1", "2"}, true, true, true},
		{[]interface{}{"1", "2"}, true, false, false},
		{[]interface{}{"1", "2"}, false, false, true},
	}

	trace := getTrace(3)
	l, _ := New()
	for i, s := range tests {
		var buf = new(bytes.Buffer)
		l.Writer = buf
		l.FatalStatusCode = 0
		l.Format(s.showFilePath, s.showFuncName, s.showFileLine)
		l.Fatalln(s.data...)

		res := buf.String()
		exp := getPrefix(trace, "", l.Timestamp, FATAL, s.showFilePath,
			s.showFuncName, s.showFileLine) + fmt.Sprintln(s.data...)

		if l.Levels.All(FATAL) {
			res = res[20:]
			exp = exp[20:]
		} else {
			exp = "" // if level not supportde
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}

// TestFerror tests Ferror method.
func TestFerror(t *testing.T) {
	type test struct {
		data []interface{}

		showFilePath bool
		showFuncName bool
		showFileLine bool
	}

	var tests = []test{
		{[]interface{}{"1 ", "2"}, true, true, true},
		{[]interface{}{"1 ", "2"}, true, false, false},
		{[]interface{}{"1 ", "2"}, false, false, true},
	}

	trace := getTrace(3)
	l, _ := New()
	for i, s := range tests {
		var buf = new(bytes.Buffer)
		l.Format(s.showFilePath, s.showFuncName, s.showFileLine)
		l.Ferror(buf, s.data...)

		res := buf.String()
		exp := getPrefix(trace, "", l.Timestamp, ERROR, s.showFilePath,
			s.showFuncName, s.showFileLine) + fmt.Sprint(s.data...)

		if l.Levels.All(ERROR) {
			res = res[20:]
			exp = exp[20:]
		} else {
			exp = "" // if level not supportde
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}

// TestFerrorf tests Ferrorf method.
func TestFerrorf(t *testing.T) {
	type test struct {
		format string
		data   []interface{}

		showFilePath bool
		showFuncName bool
		showFileLine bool
	}

	var tests = []test{
		{"%s-%s", []interface{}{"1", "2"}, true, true, true},
		{"%s %s", []interface{}{"1", "2"}, true, false, false},
		{"%s+%s", []interface{}{"1", "2"}, false, false, true},
	}

	trace := getTrace(3)
	l, _ := New()
	for i, s := range tests {
		var buf = new(bytes.Buffer)
		l.Format(s.showFilePath, s.showFuncName, s.showFileLine)
		l.Ferrorf(buf, s.format, s.data...)

		res := buf.String()
		prefix := getPrefix(trace, s.format, l.Timestamp, ERROR,
			s.showFilePath, s.showFuncName, s.showFileLine)
		exp := fmt.Sprintf(prefix, s.data...)

		if l.Levels.All(ERROR) {
			res = res[20:]
			exp = exp[20:]
		} else {
			exp = "" // if level not supportde
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}

// TestFerrorln tests Ferrorln method.
func TestFerrorln(t *testing.T) {
	type test struct {
		data []interface{}

		showFilePath bool
		showFuncName bool
		showFileLine bool
	}

	var tests = []test{
		{[]interface{}{"1", "2"}, true, true, true},
		{[]interface{}{"1", "2"}, true, false, false},
		{[]interface{}{"1", "2"}, false, false, true},
	}

	trace := getTrace(3)
	l, _ := New()
	for i, s := range tests {
		var buf = new(bytes.Buffer)
		l.Format(s.showFilePath, s.showFuncName, s.showFileLine)
		l.Ferrorln(buf, s.data...)

		res := buf.String()
		exp := getPrefix(trace, "", l.Timestamp, ERROR, s.showFilePath,
			s.showFuncName, s.showFileLine) + fmt.Sprintln(s.data...)

		if l.Levels.All(ERROR) {
			res = res[20:]
			exp = exp[20:]
		} else {
			exp = "" // if level not supportde
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}

// TestError tests Error method.
func TestError(t *testing.T) {
	type test struct {
		data []interface{}

		showFilePath bool
		showFuncName bool
		showFileLine bool
	}

	var tests = []test{
		{[]interface{}{"1 ", "2"}, true, true, true},
		{[]interface{}{"1 ", "2"}, true, false, false},
		{[]interface{}{"1 ", "2"}, false, false, true},
	}

	trace := getTrace(3)
	l, _ := New()
	for i, s := range tests {
		var buf = new(bytes.Buffer)
		l.Writer = buf
		l.Format(s.showFilePath, s.showFuncName, s.showFileLine)
		l.Error(s.data...)

		res := buf.String()
		exp := getPrefix(trace, "", l.Timestamp, ERROR, s.showFilePath,
			s.showFuncName, s.showFileLine) + fmt.Sprint(s.data...)

		if l.Levels.All(ERROR) {
			res = res[20:]
			exp = exp[20:]
		} else {
			exp = "" // if level not supportde
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}

// TestErrorf tests Errorf method.
func TestErrorf(t *testing.T) {
	type test struct {
		format string
		data   []interface{}

		showFilePath bool
		showFuncName bool
		showFileLine bool
	}

	var tests = []test{
		{"%s-%s", []interface{}{"1", "2"}, true, true, true},
		{"%s %s", []interface{}{"1", "2"}, true, false, false},
		{"%s+%s", []interface{}{"1", "2"}, false, false, true},
	}

	trace := getTrace(3)
	l, _ := New()
	for i, s := range tests {
		var buf = new(bytes.Buffer)
		l.Writer = buf
		l.Format(s.showFilePath, s.showFuncName, s.showFileLine)
		l.Errorf(s.format, s.data...)

		res := buf.String()
		prefix := getPrefix(trace, s.format, l.Timestamp, ERROR,
			s.showFilePath, s.showFuncName, s.showFileLine)
		exp := fmt.Sprintf(prefix, s.data...)

		if l.Levels.All(ERROR) {
			res = res[20:]
			exp = exp[20:]
		} else {
			exp = "" // if level not supportde
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}

// TestErrorln tests Errorln method.
func TestErrorln(t *testing.T) {
	type test struct {
		data []interface{}

		showFilePath bool
		showFuncName bool
		showFileLine bool
	}

	var tests = []test{
		{[]interface{}{"1", "2"}, true, true, true},
		{[]interface{}{"1", "2"}, true, false, false},
		{[]interface{}{"1", "2"}, false, false, true},
	}

	trace := getTrace(3)
	l, _ := New()
	for i, s := range tests {
		var buf = new(bytes.Buffer)
		l.Writer = buf
		l.Format(s.showFilePath, s.showFuncName, s.showFileLine)
		l.Errorln(s.data...)

		res := buf.String()
		exp := getPrefix(trace, "", l.Timestamp, ERROR, s.showFilePath,
			s.showFuncName, s.showFileLine) + fmt.Sprintln(s.data...)

		if l.Levels.All(ERROR) {
			res = res[20:]
			exp = exp[20:]
		} else {
			exp = "" // if level not supportde
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}

// TestFwarn tests Fwarn method.
func TestFwarn(t *testing.T) {
	type test struct {
		data []interface{}

		showFilePath bool
		showFuncName bool
		showFileLine bool
	}

	var tests = []test{
		{[]interface{}{"1 ", "2"}, true, true, true},
		{[]interface{}{"1 ", "2"}, true, false, false},
		{[]interface{}{"1 ", "2"}, false, false, true},
	}

	trace := getTrace(3)
	l, _ := New()
	for i, s := range tests {
		var buf = new(bytes.Buffer)
		l.Format(s.showFilePath, s.showFuncName, s.showFileLine)
		l.Fwarn(buf, s.data...)

		res := buf.String()
		exp := getPrefix(trace, "", l.Timestamp, WARN, s.showFilePath,
			s.showFuncName, s.showFileLine) + fmt.Sprint(s.data...)

		if l.Levels.All(WARN) {
			res = res[20:]
			exp = exp[20:]
		} else {
			exp = "" // if level not supportde
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}

// TestFwarnf tests Fwarnf method.
func TestFwarnf(t *testing.T) {
	type test struct {
		format string
		data   []interface{}

		showFilePath bool
		showFuncName bool
		showFileLine bool
	}

	var tests = []test{
		{"%s-%s", []interface{}{"1", "2"}, true, true, true},
		{"%s %s", []interface{}{"1", "2"}, true, false, false},
		{"%s+%s", []interface{}{"1", "2"}, false, false, true},
	}

	trace := getTrace(3)
	l, _ := New()
	for i, s := range tests {
		var buf = new(bytes.Buffer)
		l.Format(s.showFilePath, s.showFuncName, s.showFileLine)
		l.Fwarnf(buf, s.format, s.data...)

		res := buf.String()
		prefix := getPrefix(trace, s.format, l.Timestamp, WARN,
			s.showFilePath, s.showFuncName, s.showFileLine)
		exp := fmt.Sprintf(prefix, s.data...)

		if l.Levels.All(WARN) {
			res = res[20:]
			exp = exp[20:]
		} else {
			exp = "" // if level not supportde
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}

// TestFwarnln tests Fwarnln method.
func TestFwarnln(t *testing.T) {
	type test struct {
		data []interface{}

		showFilePath bool
		showFuncName bool
		showFileLine bool
	}

	var tests = []test{
		{[]interface{}{"1", "2"}, true, true, true},
		{[]interface{}{"1", "2"}, true, false, false},
		{[]interface{}{"1", "2"}, false, false, true},
	}

	trace := getTrace(3)
	l, _ := New()
	for i, s := range tests {
		var buf = new(bytes.Buffer)
		l.Format(s.showFilePath, s.showFuncName, s.showFileLine)
		l.Fwarnln(buf, s.data...)

		res := buf.String()
		exp := getPrefix(trace, "", l.Timestamp, WARN, s.showFilePath,
			s.showFuncName, s.showFileLine) + fmt.Sprintln(s.data...)

		if l.Levels.All(WARN) {
			res = res[20:]
			exp = exp[20:]
		} else {
			exp = "" // if level not supportde
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}

// TestWarn tests Warn method.
func TestWarn(t *testing.T) {
	type test struct {
		data []interface{}

		showFilePath bool
		showFuncName bool
		showFileLine bool
	}

	var tests = []test{
		{[]interface{}{"1 ", "2"}, true, true, true},
		{[]interface{}{"1 ", "2"}, true, false, false},
		{[]interface{}{"1 ", "2"}, false, false, true},
	}

	trace := getTrace(3)
	l, _ := New()
	for i, s := range tests {
		var buf = new(bytes.Buffer)
		l.Writer = buf
		l.Format(s.showFilePath, s.showFuncName, s.showFileLine)
		l.Warn(s.data...)

		res := buf.String()
		exp := getPrefix(trace, "", l.Timestamp, WARN, s.showFilePath,
			s.showFuncName, s.showFileLine) + fmt.Sprint(s.data...)

		if l.Levels.All(WARN) {
			res = res[20:]
			exp = exp[20:]
		} else {
			exp = "" // if level not supportde
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}

// TestWarnf tests Warnf method.
func TestWarnf(t *testing.T) {
	type test struct {
		format string
		data   []interface{}

		showFilePath bool
		showFuncName bool
		showFileLine bool
	}

	var tests = []test{
		{"%s-%s", []interface{}{"1", "2"}, true, true, true},
		{"%s %s", []interface{}{"1", "2"}, true, false, false},
		{"%s+%s", []interface{}{"1", "2"}, false, false, true},
	}

	trace := getTrace(3)
	l, _ := New()
	for i, s := range tests {
		var buf = new(bytes.Buffer)
		l.Writer = buf
		l.Format(s.showFilePath, s.showFuncName, s.showFileLine)
		l.Warnf(s.format, s.data...)

		res := buf.String()
		prefix := getPrefix(trace, s.format, l.Timestamp, WARN,
			s.showFilePath, s.showFuncName, s.showFileLine)
		exp := fmt.Sprintf(prefix, s.data...)

		if l.Levels.All(WARN) {
			res = res[20:]
			exp = exp[20:]
		} else {
			exp = "" // if level not supportde
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}

// TestWarnln tests Warnln method.
func TestWarnln(t *testing.T) {
	type test struct {
		data []interface{}

		showFilePath bool
		showFuncName bool
		showFileLine bool
	}

	var tests = []test{
		{[]interface{}{"1", "2"}, true, true, true},
		{[]interface{}{"1", "2"}, true, false, false},
		{[]interface{}{"1", "2"}, false, false, true},
	}

	trace := getTrace(3)
	l, _ := New()
	for i, s := range tests {
		var buf = new(bytes.Buffer)
		l.Writer = buf
		l.Format(s.showFilePath, s.showFuncName, s.showFileLine)
		l.Warnln(s.data...)

		res := buf.String()
		exp := getPrefix(trace, "", l.Timestamp, WARN, s.showFilePath,
			s.showFuncName, s.showFileLine) + fmt.Sprintln(s.data...)

		if l.Levels.All(WARN) {
			res = res[20:]
			exp = exp[20:]
		} else {
			exp = "" // if level not supportde
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}

// TestFinfo tests Finfo method.
func TestFinfo(t *testing.T) {
	type test struct {
		data []interface{}

		showFilePath bool
		showFuncName bool
		showFileLine bool
	}

	var tests = []test{
		{[]interface{}{"1 ", "2"}, true, true, true},
		{[]interface{}{"1 ", "2"}, true, false, false},
		{[]interface{}{"1 ", "2"}, false, false, true},
	}

	trace := getTrace(3)
	l, _ := New()
	for i, s := range tests {
		var buf = new(bytes.Buffer)
		l.Format(s.showFilePath, s.showFuncName, s.showFileLine)
		l.Finfo(buf, s.data...)

		res := buf.String()
		exp := getPrefix(trace, "", l.Timestamp, INFO, s.showFilePath,
			s.showFuncName, s.showFileLine) + fmt.Sprint(s.data...)

		if l.Levels.All(INFO) {
			res = res[20:]
			exp = exp[20:]
		} else {
			exp = "" // if level not supportde
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}

// TestFinfof tests Finfof method.
func TestFinfof(t *testing.T) {
	type test struct {
		format string
		data   []interface{}

		showFilePath bool
		showFuncName bool
		showFileLine bool
	}

	var tests = []test{
		{"%s-%s", []interface{}{"1", "2"}, true, true, true},
		{"%s %s", []interface{}{"1", "2"}, true, false, false},
		{"%s+%s", []interface{}{"1", "2"}, false, false, true},
	}

	trace := getTrace(3)
	l, _ := New()
	for i, s := range tests {
		var buf = new(bytes.Buffer)
		l.Format(s.showFilePath, s.showFuncName, s.showFileLine)
		l.Finfof(buf, s.format, s.data...)

		res := buf.String()
		prefix := getPrefix(trace, s.format, l.Timestamp, INFO,
			s.showFilePath, s.showFuncName, s.showFileLine)
		exp := fmt.Sprintf(prefix, s.data...)

		if l.Levels.All(INFO) {
			res = res[20:]
			exp = exp[20:]
		} else {
			exp = "" // if level not supportde
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}

// TestFinfoln tests Finfoln method.
func TestFinfoln(t *testing.T) {
	type test struct {
		data []interface{}

		showFilePath bool
		showFuncName bool
		showFileLine bool
	}

	var tests = []test{
		{[]interface{}{"1", "2"}, true, true, true},
		{[]interface{}{"1", "2"}, true, false, false},
		{[]interface{}{"1", "2"}, false, false, true},
	}

	trace := getTrace(3)
	l, _ := New()
	for i, s := range tests {
		var buf = new(bytes.Buffer)
		l.Format(s.showFilePath, s.showFuncName, s.showFileLine)
		l.Finfoln(buf, s.data...)

		res := buf.String()
		exp := getPrefix(trace, "", l.Timestamp, INFO, s.showFilePath,
			s.showFuncName, s.showFileLine) + fmt.Sprintln(s.data...)

		if l.Levels.All(INFO) {
			res = res[20:]
			exp = exp[20:]
		} else {
			exp = "" // if level not supportde
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}

// TestInfo tests Info method.
func TestInfo(t *testing.T) {
	type test struct {
		data []interface{}

		showFilePath bool
		showFuncName bool
		showFileLine bool
	}

	var tests = []test{
		{[]interface{}{"1 ", "2"}, true, true, true},
		{[]interface{}{"1 ", "2"}, true, false, false},
		{[]interface{}{"1 ", "2"}, false, false, true},
	}

	trace := getTrace(3)
	l, _ := New()
	for i, s := range tests {
		var buf = new(bytes.Buffer)
		l.Writer = buf
		l.Format(s.showFilePath, s.showFuncName, s.showFileLine)
		l.Info(s.data...)

		res := buf.String()
		exp := getPrefix(trace, "", l.Timestamp, INFO, s.showFilePath,
			s.showFuncName, s.showFileLine) + fmt.Sprint(s.data...)

		if l.Levels.All(INFO) {
			res = res[20:]
			exp = exp[20:]
		} else {
			exp = "" // if level not supportde
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}

// TestInfof tests Infof method.
func TestInfof(t *testing.T) {
	type test struct {
		format string
		data   []interface{}

		showFilePath bool
		showFuncName bool
		showFileLine bool
	}

	var tests = []test{
		{"%s-%s", []interface{}{"1", "2"}, true, true, true},
		{"%s %s", []interface{}{"1", "2"}, true, false, false},
		{"%s+%s", []interface{}{"1", "2"}, false, false, true},
	}

	trace := getTrace(3)
	l, _ := New()
	for i, s := range tests {
		var buf = new(bytes.Buffer)
		l.Writer = buf
		l.Format(s.showFilePath, s.showFuncName, s.showFileLine)
		l.Infof(s.format, s.data...)

		res := buf.String()
		prefix := getPrefix(trace, s.format, l.Timestamp, INFO,
			s.showFilePath, s.showFuncName, s.showFileLine)
		exp := fmt.Sprintf(prefix, s.data...)

		if l.Levels.All(INFO) {
			res = res[20:]
			exp = exp[20:]
		} else {
			exp = "" // if level not supportde
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}

// TestInfoln tests Infoln method.
func TestInfoln(t *testing.T) {
	type test struct {
		data []interface{}

		showFilePath bool
		showFuncName bool
		showFileLine bool
	}

	var tests = []test{
		{[]interface{}{"1", "2"}, true, true, true},
		{[]interface{}{"1", "2"}, true, false, false},
		{[]interface{}{"1", "2"}, false, false, true},
	}

	trace := getTrace(3)
	l, _ := New()
	for i, s := range tests {
		var buf = new(bytes.Buffer)
		l.Writer = buf
		l.Format(s.showFilePath, s.showFuncName, s.showFileLine)
		l.Infoln(s.data...)

		res := buf.String()
		exp := getPrefix(trace, "", l.Timestamp, INFO, s.showFilePath,
			s.showFuncName, s.showFileLine) + fmt.Sprintln(s.data...)

		if l.Levels.All(INFO) {
			res = res[20:]
			exp = exp[20:]
		} else {
			exp = "" // if level not supportde
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}

// TestFdebug tests Fdebug method.
func TestFdebug(t *testing.T) {
	type test struct {
		data []interface{}

		showFilePath bool
		showFuncName bool
		showFileLine bool
	}

	var tests = []test{
		{[]interface{}{"1 ", "2"}, true, true, true},
		{[]interface{}{"1 ", "2"}, true, false, false},
		{[]interface{}{"1 ", "2"}, false, false, true},
	}

	trace := getTrace(3)
	l, _ := New()
	for i, s := range tests {
		var buf = new(bytes.Buffer)
		l.Format(s.showFilePath, s.showFuncName, s.showFileLine)
		l.Fdebug(buf, s.data...)

		res := buf.String()
		exp := getPrefix(trace, "", l.Timestamp, DEBUG, s.showFilePath,
			s.showFuncName, s.showFileLine) + fmt.Sprint(s.data...)

		if l.Levels.All(DEBUG) {
			res = res[20:]
			exp = exp[20:]
		} else {
			exp = "" // if level not supportde
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}

// TestFdebugf tests Fdebugf method.
func TestFdebugf(t *testing.T) {
	type test struct {
		format string
		data   []interface{}

		showFilePath bool
		showFuncName bool
		showFileLine bool
	}

	var tests = []test{
		{"%s-%s", []interface{}{"1", "2"}, true, true, true},
		{"%s %s", []interface{}{"1", "2"}, true, false, false},
		{"%s+%s", []interface{}{"1", "2"}, false, false, true},
	}

	trace := getTrace(3)
	l, _ := New()
	for i, s := range tests {
		var buf = new(bytes.Buffer)
		l.Format(s.showFilePath, s.showFuncName, s.showFileLine)
		l.Fdebugf(buf, s.format, s.data...)

		res := buf.String()
		prefix := getPrefix(trace, s.format, l.Timestamp, DEBUG,
			s.showFilePath, s.showFuncName, s.showFileLine)
		exp := fmt.Sprintf(prefix, s.data...)

		if l.Levels.All(DEBUG) {
			res = res[20:]
			exp = exp[20:]
		} else {
			exp = "" // if level not supportde
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}

// TestFdebugln tests Fdebugln method.
func TestFdebugln(t *testing.T) {
	type test struct {
		data []interface{}

		showFilePath bool
		showFuncName bool
		showFileLine bool
	}

	var tests = []test{
		{[]interface{}{"1", "2"}, true, true, true},
		{[]interface{}{"1", "2"}, true, false, false},
		{[]interface{}{"1", "2"}, false, false, true},
	}

	trace := getTrace(3)
	l, _ := New()
	for i, s := range tests {
		var buf = new(bytes.Buffer)
		l.Format(s.showFilePath, s.showFuncName, s.showFileLine)
		l.Fdebugln(buf, s.data...)

		res := buf.String()
		exp := getPrefix(trace, "", l.Timestamp, DEBUG, s.showFilePath,
			s.showFuncName, s.showFileLine) + fmt.Sprintln(s.data...)

		if l.Levels.All(DEBUG) {
			res = res[20:]
			exp = exp[20:]
		} else {
			exp = "" // if level not supportde
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}

// TestDebug tests Debug method.
func TestDebug(t *testing.T) {
	type test struct {
		data []interface{}

		showFilePath bool
		showFuncName bool
		showFileLine bool
	}

	var tests = []test{
		{[]interface{}{"1 ", "2"}, true, true, true},
		{[]interface{}{"1 ", "2"}, true, false, false},
		{[]interface{}{"1 ", "2"}, false, false, true},
	}

	trace := getTrace(3)
	l, _ := New()
	for i, s := range tests {
		var buf = new(bytes.Buffer)
		l.Writer = buf
		l.Format(s.showFilePath, s.showFuncName, s.showFileLine)
		l.Debug(s.data...)

		res := buf.String()
		exp := getPrefix(trace, "", l.Timestamp, DEBUG, s.showFilePath,
			s.showFuncName, s.showFileLine) + fmt.Sprint(s.data...)

		if l.Levels.All(DEBUG) {
			res = res[20:]
			exp = exp[20:]
		} else {
			exp = "" // if level not supportde
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}

// TestDebugf tests Debugf method.
func TestDebugf(t *testing.T) {
	type test struct {
		format string
		data   []interface{}

		showFilePath bool
		showFuncName bool
		showFileLine bool
	}

	var tests = []test{
		{"%s-%s", []interface{}{"1", "2"}, true, true, true},
		{"%s %s", []interface{}{"1", "2"}, true, false, false},
		{"%s+%s", []interface{}{"1", "2"}, false, false, true},
	}

	trace := getTrace(3)
	l, _ := New()
	for i, s := range tests {
		var buf = new(bytes.Buffer)
		l.Writer = buf
		l.Format(s.showFilePath, s.showFuncName, s.showFileLine)
		l.Debugf(s.format, s.data...)

		res := buf.String()
		prefix := getPrefix(trace, s.format, l.Timestamp, DEBUG,
			s.showFilePath, s.showFuncName, s.showFileLine)
		exp := fmt.Sprintf(prefix, s.data...)

		if l.Levels.All(DEBUG) {
			res = res[20:]
			exp = exp[20:]
		} else {
			exp = "" // if level not supportde
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}

// TestDebugln tests Debugln method.
func TestDebugln(t *testing.T) {
	type test struct {
		data []interface{}

		showFilePath bool
		showFuncName bool
		showFileLine bool
	}

	var tests = []test{
		{[]interface{}{"1", "2"}, true, true, true},
		{[]interface{}{"1", "2"}, true, false, false},
		{[]interface{}{"1", "2"}, false, false, true},
	}

	trace := getTrace(3)
	l, _ := New()
	for i, s := range tests {
		var buf = new(bytes.Buffer)
		l.Writer = buf
		l.Format(s.showFilePath, s.showFuncName, s.showFileLine)
		l.Debugln(s.data...)

		res := buf.String()
		exp := getPrefix(trace, "", l.Timestamp, DEBUG, s.showFilePath,
			s.showFuncName, s.showFileLine) + fmt.Sprintln(s.data...)

		if l.Levels.All(DEBUG) {
			res = res[20:]
			exp = exp[20:]
		} else {
			exp = "" // if level not supportde
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}

// TestFtrace tests Ftrace method.
func TestFtrace(t *testing.T) {
	type test struct {
		data []interface{}

		showFilePath bool
		showFuncName bool
		showFileLine bool
	}

	var tests = []test{
		{[]interface{}{"1 ", "2"}, true, true, true},
		{[]interface{}{"1 ", "2"}, true, false, false},
		{[]interface{}{"1 ", "2"}, false, false, true},
	}

	trace := getTrace(3)
	l, _ := New()
	for i, s := range tests {
		var buf = new(bytes.Buffer)
		l.Format(s.showFilePath, s.showFuncName, s.showFileLine)
		l.Ftrace(buf, s.data...)

		res := buf.String()
		exp := getPrefix(trace, "", l.Timestamp, TRACE, s.showFilePath,
			s.showFuncName, s.showFileLine) + fmt.Sprint(s.data...)

		if l.Levels.All(TRACE) {
			res = res[20:]
			exp = exp[20:]
		} else {
			exp = "" // if level not supportde
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}

// TestFtracef tests Ftracef method.
func TestFtracef(t *testing.T) {
	type test struct {
		format string
		data   []interface{}

		showFilePath bool
		showFuncName bool
		showFileLine bool
	}

	var tests = []test{
		{"%s-%s", []interface{}{"1", "2"}, true, true, true},
		{"%s %s", []interface{}{"1", "2"}, true, false, false},
		{"%s+%s", []interface{}{"1", "2"}, false, false, true},
	}

	trace := getTrace(3)
	l, _ := New()
	for i, s := range tests {
		var buf = new(bytes.Buffer)
		l.Format(s.showFilePath, s.showFuncName, s.showFileLine)
		l.Ftracef(buf, s.format, s.data...)

		res := buf.String()
		prefix := getPrefix(trace, s.format, l.Timestamp, TRACE,
			s.showFilePath, s.showFuncName, s.showFileLine)
		exp := fmt.Sprintf(prefix, s.data...)

		if l.Levels.All(TRACE) {
			res = res[20:]
			exp = exp[20:]
		} else {
			exp = "" // if level not supportde
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}

// TestFtraceln tests Ftraceln method.
func TestFtraceln(t *testing.T) {
	type test struct {
		data []interface{}

		showFilePath bool
		showFuncName bool
		showFileLine bool
	}

	var tests = []test{
		{[]interface{}{"1", "2"}, true, true, true},
		{[]interface{}{"1", "2"}, true, false, false},
		{[]interface{}{"1", "2"}, false, false, true},
	}

	trace := getTrace(3)
	l, _ := New()
	for i, s := range tests {
		var buf = new(bytes.Buffer)
		l.Format(s.showFilePath, s.showFuncName, s.showFileLine)
		l.Ftraceln(buf, s.data...)

		res := buf.String()
		exp := getPrefix(trace, "", l.Timestamp, TRACE, s.showFilePath,
			s.showFuncName, s.showFileLine) + fmt.Sprintln(s.data...)

		if l.Levels.All(TRACE) {
			res = res[20:]
			exp = exp[20:]
		} else {
			exp = "" // if level not supportde
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}

// TestTrace tests Trace method.
func TestTrace(t *testing.T) {
	type test struct {
		data []interface{}

		showFilePath bool
		showFuncName bool
		showFileLine bool
	}

	var tests = []test{
		{[]interface{}{"1 ", "2"}, true, true, true},
		{[]interface{}{"1 ", "2"}, true, false, false},
		{[]interface{}{"1 ", "2"}, false, false, true},
	}

	trace := getTrace(3)
	l, _ := New()
	for i, s := range tests {
		var buf = new(bytes.Buffer)
		l.Writer = buf
		l.Format(s.showFilePath, s.showFuncName, s.showFileLine)
		l.Trace(s.data...)

		res := buf.String()
		exp := getPrefix(trace, "", l.Timestamp, TRACE, s.showFilePath,
			s.showFuncName, s.showFileLine) + fmt.Sprint(s.data...)

		if l.Levels.All(TRACE) {
			res = res[20:]
			exp = exp[20:]
		} else {
			exp = "" // if level not supportde
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}

// TestTracef tests Tracef method.
func TestTracef(t *testing.T) {
	type test struct {
		format string
		data   []interface{}

		showFilePath bool
		showFuncName bool
		showFileLine bool
	}

	var tests = []test{
		{"%s-%s", []interface{}{"1", "2"}, true, true, true},
		{"%s %s", []interface{}{"1", "2"}, true, false, false},
		{"%s+%s", []interface{}{"1", "2"}, false, false, true},
	}

	trace := getTrace(3)
	l, _ := New()
	for i, s := range tests {
		var buf = new(bytes.Buffer)
		l.Writer = buf
		l.Format(s.showFilePath, s.showFuncName, s.showFileLine)
		l.Tracef(s.format, s.data...)

		res := buf.String()
		prefix := getPrefix(trace, s.format, l.Timestamp, TRACE,
			s.showFilePath, s.showFuncName, s.showFileLine)
		exp := fmt.Sprintf(prefix, s.data...)

		if l.Levels.All(TRACE) {
			res = res[20:]
			exp = exp[20:]
		} else {
			exp = "" // if level not supportde
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}

// TestTraceln tests Traceln method.
func TestTraceln(t *testing.T) {
	type test struct {
		data []interface{}

		showFilePath bool
		showFuncName bool
		showFileLine bool
	}

	var tests = []test{
		{[]interface{}{"1", "2"}, true, true, true},
		{[]interface{}{"1", "2"}, true, false, false},
		{[]interface{}{"1", "2"}, false, false, true},
	}

	trace := getTrace(3)
	l, _ := New()
	for i, s := range tests {
		var buf = new(bytes.Buffer)
		l.Writer = buf
		l.Format(s.showFilePath, s.showFuncName, s.showFileLine)
		l.Traceln(s.data...)

		res := buf.String()
		exp := getPrefix(trace, "", l.Timestamp, TRACE, s.showFilePath,
			s.showFuncName, s.showFileLine) + fmt.Sprintln(s.data...)

		if l.Levels.All(TRACE) {
			res = res[20:]
			exp = exp[20:]
		} else {
			exp = "" // if level not supportde
		}

		if !strings.HasSuffix(res, exp) {
			t.Errorf("test %d is failed, expected `%s` but `%s`", i, exp, res)
		}
	}
}
