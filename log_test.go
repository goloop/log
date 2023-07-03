package log

import (
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/goloop/log/level"
)

// TestNew tests the New function.
func TestNew(t *testing.T) {
	tests := []struct {
		name     string
		prefixes []string
		want     string
	}{
		{
			name:     "No prefix",
			prefixes: []string{},
			want:     "",
		},
		{
			name:     "One prefix",
			prefixes: []string{"test"},
			want:     "test",
		},
		{
			name:     "Multiple prefixes",
			prefixes: []string{"test", "logger"},
			want:     "test-logger",
		},
		{
			name:     "Empty prefixes",
			prefixes: []string{"", " "},
			want:     "",
		},
		{
			name:     "With marker",
			prefixes: []string{"myapp", ":"},
			want:     "myapp:",
		},
		{
			name:     "Two words with marker",
			prefixes: []string{"my", "app", ":"},
			want:     "my-app:",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			logger := New(test.prefixes...)
			if logger.prefix != test.want {
				t.Errorf("%s got = %v; want %v",
					test.name, logger.prefix, test.want)
			}

			if logger.skipStackFrames != skipStackFrames {
				t.Errorf("%s skipStackFrames got = %d; want %d",
					test.name, logger.skipStackFrames, skipStackFrames)
			}

			if logger.fatalStatusCode != fatalStatusCode {
				t.Errorf("%s fatalStatusCode got = %d; want %d",
					test.name, logger.fatalStatusCode, fatalStatusCode)
			}

			if len(logger.outputs) != 2 ||
				logger.outputs["stdout"] == nil ||
				logger.outputs["stderr"] == nil {
				t.Errorf("%s outputs are not properly set", test.name)
			}
		})
	}
}

// TestCopy tests the Copy function.
func TestCopy(t *testing.T) {
	copy := Copy() // copy self logger

	// Perform an in-depth comparison of objects.
	if !reflect.DeepEqual(copy, self) {
		t.Errorf("Copy method doesn't make a complete copy of the object.")
	}
}

// TestSetSkipStackFrames tests the SetSkipStackFrames function.
func TestSetSkipStackFrames(t *testing.T) {
	skip := SkipStackFrames()
	defer SetSkipStackFrames(skip)

	tests := []struct {
		name string
		skip int
		want int
	}{
		{
			name: "Seat skip as -1",
			skip: -1,
			want: SkipStackFrames(), // current value
		},
		{
			name: "Seat skip as 1",
			skip: 1,
			want: 1,
		},
		{
			name: "Seat skip as 2",
			skip: 2,
			want: 2,
		},
		{
			name: "Very high value",
			skip: 32,
			want: -1, // does not have to match
		},
	}

	// Don't use parallel tests here.
	for _, tt := range tests {
		SetSkipStackFrames(tt.skip)
		if tt.want < 0 {
			if s := SkipStackFrames(); s == tt.skip {
				t.Errorf("%s: skip limit protection did not work", tt.name)
			}
		} else {
			if s := SkipStackFrames(); s != tt.want {
				t.Errorf("%s: failed, got %d, want %d",
					tt.name, s, tt.want)
			}
		}
	}
}

// TestSetPrefix tests the SetPrefix function.
func TestSetPrefix(t *testing.T) {
	prefix := Prefix()
	defer SetPrefix(prefix)

	tests := []struct {
		name string
		in   string
	}{
		{
			name: "Set: `hello`",
			in:   "hello",
		},
		{
			name: "Set: `hello-world`",
			in:   "hello-world",
		},
	}

	// Don't use parallel tests here.
	for _, tt := range tests {
		SetPrefix(tt.in)
		if s := Prefix(); s != tt.in {
			t.Errorf("%s: failed, got %s, want %s",
				tt.name, s, tt.in)
		}
	}
}

// TestSetOutputs tests the SetOutputs function.
func TestSetOutputs(t *testing.T) {
	outputs := Outputs()
	defer SetOutputs(outputs...)

	tests := []struct {
		name   string
		in     []Output
		names  []string
		hasErr bool
	}{
		{
			name:   "Empty outputs",
			in:     []Output{},
			names:  []string{},
			hasErr: true,
		},
		{
			name: "One short output",
			in: []Output{
				{
					Name:   "test",
					Writer: os.Stdout,
				},
			},
			names:  []string{"test"},
			hasErr: false,
		},
		{
			name: "One short output with defaults",
			in: []Output{
				{
					Name:   "test",
					Writer: os.Stdout,
				},
				Stdout,
				Stderr,
			},
			names:  []string{"test", "stdout", "stderr"},
			hasErr: false,
		},
		{
			name: "Noname output",
			in: []Output{
				{
					Writer: os.Stdout,
				},
			},
			hasErr: true,
		},
		{
			name: "Nowriter output",
			in: []Output{
				{
					Name: "test",
				},
			},
			hasErr: true,
		},
		{
			name: "Duplicates outputs",
			in: []Output{
				{
					Name:   "test",
					Writer: os.Stdout,
				},
				{
					Name:   "test",
					Writer: os.Stdout,
				},
			},
			hasErr: true,
		},
	}

	// Don't use parallel tests here.
	for _, tt := range tests {
		// Error check.
		err := SetOutputs(tt.in...)
		if tt.hasErr {
			if err == nil {
				t.Errorf("%s: an error was expected", tt.name)
			}

			continue
		}

		if err != nil {
			t.Errorf("%s: an error occurred: %s", tt.name, err.Error())
		}

		// Len check.
		out := Outputs()
		if len(out) != len(tt.in) {
			t.Errorf("%s: %d items passed but %d items sets",
				tt.name, len(out), len(tt.in))
		}

		// Check names.
		for _, n := range tt.names {
			if _, ok := self.outputs[n]; !ok {
				t.Errorf("%s: %s output not found", tt.name, n)
			}
		}
	}
}

// TestEditOutputs tests the EditOutputs function.
func TestEditOutputs(t *testing.T) {
	outputs := Outputs()
	defer SetOutputs(outputs...)

	tests := []struct {
		name   string
		in     []Output
		hasErr bool
	}{
		{
			name:   "Empty outputs",
			in:     []Output{},
			hasErr: true,
		},
		{
			name: "One short output",
			in: []Output{
				{
					Name:   "some-not-real-name",
					Writer: os.Stdout,
				},
			},
			hasErr: true,
		},
		{
			name:   "Update stdout and stderr",
			in:     []Output{Stdout, Stderr},
			hasErr: false,
		},
		{
			name: "Update stdout levels",
			in: []Output{
				{
					Name:   Stdout.Name,
					Levels: Stdout.Levels | level.Error,
				},
			},
			hasErr: false,
		},
	}

	// Don't use parallel tests here.
	for _, tt := range tests {
		// Error check.
		SetOutputs(Stdout, Stderr)
		err := EditOutputs(tt.in...)
		if tt.hasErr {
			if err == nil {
				t.Errorf("%s: an error was expected", tt.name)
			}

			continue
		}

		if err != nil {
			t.Errorf("%s: an error occurred: %s", tt.name, err.Error())
		}

		// Check names.
		for _, o := range tt.in {
			oo, ok := self.outputs[o.Name]
			if !ok {
				t.Errorf("%s: %s output not found", tt.name, o.Name)
			}

			if o.Levels != oo.Levels {
				t.Errorf("%s: %s output not updated correctly",
					tt.name, o.Name)
			}
		}
	}
}

// TestDeleteOutputs tests the DeleteOutputs function.
func TestDeleteOutputs(t *testing.T) {
	outputs := Outputs()
	defer SetOutputs(outputs...)

	DeleteOutputs(Stdout.Name)
	if _, ok := self.outputs[Stdout.Name]; ok {
		t.Errorf("DeleteOutputs did not delete the output")
	}
}

// TestOutputs tests the Outputs function.
func TestOutputs(t *testing.T) {
	outputs := Outputs()
	defer SetOutputs(outputs...)

	output := Output{
		Name:   "test",
		Writer: os.Stdout,
	}

	err := SetOutputs(output)
	if err != nil {
		t.Fatal(err)
	}

	out := Outputs()
	if len(out) != 1 {
		t.Errorf("Outputs did not return the correct number of outputs")
	}

	if out[0].Name != output.Name {
		t.Errorf("Outputs did not return the correct output")
	}
}

// TestFpanic tests the Fpanic function.
func TestFpanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic as expected")
		}
	}()

	r, w, _ := os.Pipe()
	SetOutputs(Output{
		Name:   "test",
		Writer: w,
		Levels: level.Panic,
	})

	Fpanic(w, "Test fatal")
	outC := make(chan string)
	go ioCopy(r, outC)
	w.Close()
	out := <-outC

	expected := "Test fatal"
	n := level.Labels[level.Panic]
	if !strings.Contains(out, expected) || !strings.Contains(out, n) {
		t.Errorf("Result `%s` doesn't contains `%s` and `%s`",
			out, expected, n)
	}
}

// TestFpanicf tests the Fpanicf function.
func TestFpanicf(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic as expected")
		}
	}()

	r, w, _ := os.Pipe()
	SetOutputs(Output{
		Name:   "test",
		Writer: w,
		Levels: level.Panic,
	})

	Fpanicf(w, "Test fatal %s", "formatted")
	outC := make(chan string)
	go ioCopy(r, outC)
	w.Close()
	out := <-outC

	expected := "Test fatal formatted"
	n := level.Labels[level.Panic]
	if !strings.Contains(out, expected) || !strings.Contains(out, n) {
		t.Errorf("Result `%s` doesn't contains `%s` and `%s`",
			out, expected, n)
	}
}

// TestFpanicln tests the Fpanicln function.
func TestFpanicln(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic as expected")
		}
	}()

	r, w, _ := os.Pipe()
	SetOutputs(Output{
		Name:   "test",
		Writer: w,
		Levels: level.Panic,
	})

	Fpanicln(w, "Test fatalln")
	outC := make(chan string)
	go ioCopy(r, outC)
	w.Close()
	out := <-outC

	expected := "Test fatalln"
	n := level.Labels[level.Panic]
	if !strings.Contains(out, expected) || !strings.Contains(out, n) {
		t.Errorf("Result `%s` doesn't contains `%s` and `%s`",
			out, expected, n)
	}
}

// TestPanic tests the Panic function.
func TestPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic as expected")
		}
	}()

	r, w, _ := os.Pipe()
	SetOutputs(Output{
		Name:   "test",
		Writer: w,
		Levels: level.Default,
	})

	Panic("Test fatal")
	outC := make(chan string)
	go ioCopy(r, outC)
	w.Close()
	out := <-outC

	expected := "Test fatal"
	n := level.Labels[level.Panic]
	if !strings.Contains(out, expected) || !strings.Contains(out, n) {
		t.Errorf("Result `%s` doesn't contains `%s` and `%s`",
			out, expected, n)
	}
}

// TestPanicf tests the Panicf function.
func TestPanicf(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic as expected")
		}
	}()

	r, w, _ := os.Pipe()
	SetOutputs(Output{
		Name:   "test",
		Writer: w,
		Levels: level.Default,
	})

	Panicf("Test fatal %s", "formatted")
	outC := make(chan string)
	go ioCopy(r, outC)
	w.Close()
	out := <-outC

	expected := "Test fatal formatted"
	n := level.Labels[level.Panic]
	if !strings.Contains(out, expected) || !strings.Contains(out, n) {
		t.Errorf("Result `%s` doesn't contains `%s` and `%s`",
			out, expected, n)
	}
}

// TestPanicln tests the Panicln function.
func TestPanicln(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic as expected")
		}
	}()

	r, w, _ := os.Pipe()
	SetOutputs(Output{
		Name:   "test",
		Writer: w,
		Levels: level.Default,
	})

	Panicln("Test fatalln")
	outC := make(chan string)
	go ioCopy(r, outC)
	w.Close()
	out := <-outC

	expected := "Test fatalln"
	n := level.Labels[level.Panic]
	if !strings.Contains(out, expected) || !strings.Contains(out, n) {
		t.Errorf("Result `%s` doesn't contains `%s` and `%s`",
			out, expected, n)
	}
}

// TestFfatal tests the Ffatal function.
func TestFfatal(t *testing.T) {
	exit = func(i int) {}
	defer func() {
		exit = os.Exit
	}()

	r, w, _ := os.Pipe()
	Ffatal(w, "Test Ffatal")
	outC := make(chan string)
	go ioCopy(r, outC)
	w.Close()
	out := <-outC

	expected := "Test Ffatal"
	n := level.Labels[level.Fatal]
	if !strings.Contains(out, expected) || !strings.Contains(out, n) {
		t.Errorf("Result `%s` doesn't contains `%s` and `%s`",
			out, expected, n)
	}
}

// TestFfatalf tests the Ffatalf function.
func TestFfatalf(t *testing.T) {
	exit = func(i int) {}
	defer func() {
		exit = os.Exit
	}()

	r, w, _ := os.Pipe()
	Ffatalf(w, "Test fatal %s", "formatted")
	outC := make(chan string)
	go ioCopy(r, outC)
	w.Close()
	out := <-outC

	expected := "Test fatal formatted"
	n := level.Labels[level.Fatal]
	if !strings.Contains(out, expected) || !strings.Contains(out, n) {
		t.Errorf("Result `%s` doesn't contains `%s` and `%s`",
			out, expected, n)
	}
}

// TestFfatalln tests the Ffatalln function.
func TestFfatalln(t *testing.T) {
	exit = func(i int) {}
	defer func() {
		exit = os.Exit
	}()

	r, w, _ := os.Pipe()
	Ffatalln(w, "Test fatalln")
	outC := make(chan string)
	go ioCopy(r, outC)
	w.Close()
	out := <-outC

	expected := "Test fatalln"
	n := level.Labels[level.Fatal]
	if !strings.Contains(out, expected) || !strings.Contains(out, n) {
		t.Errorf("Result `%s` doesn't contains `%s` and `%s`",
			out, expected, n)
	}
}

// TestFatal tests the Fatal function.
func TestFatal(t *testing.T) {
	exit = func(i int) {}
	defer func() {
		exit = os.Exit
	}()

	r, w, _ := os.Pipe()
	SetOutputs(Output{
		Name:   "test",
		Writer: w,
		Levels: level.Default,
	})

	Fatal("Test fatal")
	outC := make(chan string)
	go ioCopy(r, outC)
	w.Close()
	out := <-outC

	expected := "Test fatal"
	n := level.Labels[level.Fatal]
	if !strings.Contains(out, expected) || !strings.Contains(out, n) {
		t.Errorf("Result `%s` doesn't contains `%s` and `%s`",
			out, expected, n)
	}
}

// TestFatalf tests the Fatalf function.
func TestFatalf(t *testing.T) {
	exit = func(i int) {}
	defer func() {
		exit = os.Exit
	}()

	r, w, _ := os.Pipe()
	SetOutputs(Output{
		Name:   "test",
		Writer: w,
		Levels: level.Default,
	})

	Fatalf("Test fatal %s", "formatted")
	outC := make(chan string)
	go ioCopy(r, outC)
	w.Close()
	out := <-outC

	expected := "Test fatal formatted"
	n := level.Labels[level.Fatal]
	if !strings.Contains(out, expected) || !strings.Contains(out, n) {
		t.Errorf("Result `%s` doesn't contains `%s` and `%s`",
			out, expected, n)
	}
}

// TestFatalln tests the Fatalln function.
func TestFatalln(t *testing.T) {
	exit = func(i int) {}
	defer func() {
		exit = os.Exit
	}()

	r, w, _ := os.Pipe()
	SetOutputs(Output{
		Name:   "test",
		Writer: w,
		Levels: level.Default,
	})

	Fatalln("Test fatalln")
	outC := make(chan string)
	go ioCopy(r, outC)
	w.Close()
	out := <-outC

	expected := "Test fatalln"
	n := level.Labels[level.Fatal]
	if !strings.Contains(out, expected) || !strings.Contains(out, n) {
		t.Errorf("Result `%s` doesn't contains `%s` and `%s`",
			out, expected, n)
	}
}

// TestFerror tests the Ferror function.
func TestFerror(t *testing.T) {
	r, w, _ := os.Pipe()
	Ferror(w, "Test Ferror")
	outC := make(chan string)
	go ioCopy(r, outC)
	w.Close()
	out := <-outC

	expected := "Test Ferror"
	n := level.Labels[level.Error]
	if !strings.Contains(out, expected) || !strings.Contains(out, n) {
		t.Errorf("Result `%s` doesn't contains `%s` and `%s`",
			out, expected, n)
	}
}

// TestFerrorf tests the Ferrorf function.
func TestFerrorf(t *testing.T) {
	r, w, _ := os.Pipe()
	Ferrorf(w, "Test %s", "Ferrorf")
	outC := make(chan string)
	go ioCopy(r, outC)
	w.Close()
	out := <-outC

	expected := "Test Ferrorf"
	n := level.Labels[level.Error]
	if !strings.Contains(out, expected) || !strings.Contains(out, n) {
		t.Errorf("Result `%s` doesn't contains `%s` and `%s`",
			out, expected, n)
	}
}

// TestFerrorln tests the Ferrorln function.
func TestFerrorln(t *testing.T) {
	r, w, _ := os.Pipe()
	Ferrorln(w, "Test Ferrorln")
	outC := make(chan string)
	go ioCopy(r, outC)
	w.Close()
	out := <-outC

	expected := "Test Ferrorln"
	n := level.Labels[level.Error]
	if !strings.Contains(out, expected) || !strings.Contains(out, n) {
		t.Errorf("Result `%s` doesn't contains `%s` and `%s`",
			out, expected, n)
	}
}

// TestError tests the Error function.
func TestError(t *testing.T) {
	r, w, _ := os.Pipe()
	SetOutputs(Output{
		Name:   "test",
		Writer: w,
		Levels: level.Default,
	})

	Error("Test Error")
	outC := make(chan string)
	go ioCopy(r, outC)
	w.Close()
	out := <-outC

	expected := "Test Error"
	n := level.Labels[level.Error]
	if !strings.Contains(out, expected) || !strings.Contains(out, n) {
		t.Errorf("Result `%s` doesn't contains `%s` and `%s`",
			out, expected, n)
	}
}

// TestErrorf tests the Errorf function.
func TestErrorf(t *testing.T) {
	r, w, _ := os.Pipe()
	SetOutputs(Output{
		Name:   "test",
		Writer: w,
		Levels: level.Default,
	})

	Errorf("Test %s", "Errorf")
	outC := make(chan string)
	go ioCopy(r, outC)
	w.Close()
	out := <-outC

	expected := "Test Errorf"
	n := level.Labels[level.Error]
	if !strings.Contains(out, expected) || !strings.Contains(out, n) {
		t.Errorf("Result `%s` doesn't contains `%s` and `%s`",
			out, expected, n)
	}
}

// TestErrorln tests the Errorln function.
func TestErrorln(t *testing.T) {
	r, w, _ := os.Pipe()
	SetOutputs(Output{
		Name:   "test",
		Writer: w,
		Levels: level.Default,
	})

	Errorln("Test Errorln")
	outC := make(chan string)
	go ioCopy(r, outC)
	w.Close()
	out := <-outC

	expected := "Test Errorln"
	n := level.Labels[level.Error]
	if !strings.Contains(out, expected) || !strings.Contains(out, n) {
		t.Errorf("Result `%s` doesn't contains `%s` and `%s`",
			out, expected, n)
	}
}

// TestFwarn tests the Fwarn function.
func TestFwarn(t *testing.T) {
	r, w, _ := os.Pipe()
	Fwarn(w, "Test Fwarn")
	outC := make(chan string)
	go ioCopy(r, outC)
	w.Close()
	out := <-outC

	expected := "Test Fwarn"
	n := level.Labels[level.Warn]
	if !strings.Contains(out, expected) || !strings.Contains(out, n) {
		t.Errorf("Result `%s` doesn't contains `%s` and `%s`",
			out, expected, n)
	}
}

// TestFwarnf tests the Fwarnf function.
func TestFwarnf(t *testing.T) {
	r, w, _ := os.Pipe()
	Fwarnf(w, "Test %s", "Fwarnf")
	outC := make(chan string)
	go ioCopy(r, outC)
	w.Close()
	out := <-outC

	expected := "Test Fwarnf"
	n := level.Labels[level.Warn]
	if !strings.Contains(out, expected) || !strings.Contains(out, n) {
		t.Errorf("Result `%s` doesn't contains `%s` and `%s`",
			out, expected, n)
	}
}

// TestFwarnln tests the Fwarnln function.
func TestFwarnln(t *testing.T) {
	r, w, _ := os.Pipe()
	Fwarnln(w, "Test Fwarnln")
	outC := make(chan string)
	go ioCopy(r, outC)
	w.Close()
	out := <-outC

	expected := "Test Fwarnln"
	n := level.Labels[level.Warn]
	if !strings.Contains(out, expected) || !strings.Contains(out, n) {
		t.Errorf("Result `%s` doesn't contains `%s` and `%s`",
			out, expected, n)
	}
}

// TestWarn tests the Warn function.
func TestWarn(t *testing.T) {
	r, w, _ := os.Pipe()
	SetOutputs(Output{
		Name:   "test",
		Writer: w,
		Levels: level.Default,
	})

	Warn("Test Warn")
	outC := make(chan string)
	go ioCopy(r, outC)
	w.Close()
	out := <-outC

	expected := "Test Warn"
	n := level.Labels[level.Warn]
	if !strings.Contains(out, expected) || !strings.Contains(out, n) {
		t.Errorf("Result `%s` doesn't contains `%s` and `%s`",
			out, expected, n)
	}
}

// TestWarnf tests the Warnf function.
func TestWarnf(t *testing.T) {
	r, w, _ := os.Pipe()
	SetOutputs(Output{
		Name:   "test",
		Writer: w,
		Levels: level.Default,
	})

	Warnf("Test warning %s", "formatted")
	outC := make(chan string)
	go ioCopy(r, outC)
	w.Close()
	out := <-outC

	expected := "Test warning formatted"
	n := level.Labels[level.Warn]
	if !strings.Contains(out, expected) || !strings.Contains(out, n) {
		t.Errorf("Result `%s` doesn't contains `%s` and `%s`",
			out, expected, n)
	}
}

// TestWarnln tests the Warnln function.
func TestWarnln(t *testing.T) {
	r, w, _ := os.Pipe()
	SetOutputs(Output{
		Name:   "test",
		Writer: w,
		Levels: level.Default,
	})

	Warnln("Test warnln")
	outC := make(chan string)
	go ioCopy(r, outC)
	w.Close()
	out := <-outC

	expected := "Test warnln"
	n := level.Labels[level.Warn]
	if !strings.Contains(out, expected) || !strings.Contains(out, n) {
		t.Errorf("Result `%s` doesn't contains `%s` and `%s`",
			out, expected, n)
	}
}

// TestFinfo tests the Finfo function.
func TestFinfo(t *testing.T) {
	r, w, _ := os.Pipe()
	Finfo(w, "Test finfo")
	outC := make(chan string)
	go ioCopy(r, outC)
	w.Close()
	out := <-outC

	expected := "Test finfo"
	n := level.Labels[level.Info]
	if !strings.Contains(out, expected) || !strings.Contains(out, n) {
		t.Errorf("Result `%s` doesn't contains `%s` and `%s`",
			out, expected, n)
	}
}

// TestFinfof tests the Finfof function.
func TestFinfof(t *testing.T) {
	r, w, _ := os.Pipe()
	Finfof(w, "Test %s", "finfof")
	outC := make(chan string)
	go ioCopy(r, outC)
	w.Close()
	out := <-outC

	expected := "Test finfof"
	n := level.Labels[level.Info]
	if !strings.Contains(out, expected) || !strings.Contains(out, n) {
		t.Errorf("Result `%s` doesn't contains `%s` and `%s`",
			out, expected, n)
	}
}

// TestFinfoln tests the Finfoln function.
func TestFinfoln(t *testing.T) {
	r, w, _ := os.Pipe()
	Finfoln(w, "Test finfoln")
	outC := make(chan string)
	go ioCopy(r, outC)
	w.Close()
	out := <-outC

	expected := "Test finfoln"
	n := level.Labels[level.Info]
	if !strings.Contains(out, expected) || !strings.Contains(out, n) {
		t.Errorf("Result `%s` doesn't contains `%s` and `%s`",
			out, expected, n)
	}
}

// TestInfo tests the Info function.
func TestInfo(t *testing.T) {
	r, w, _ := os.Pipe()
	SetOutputs(Output{
		Name:   "test",
		Writer: w,
		Levels: level.Default,
	})

	Info("Test info")
	outC := make(chan string)
	go ioCopy(r, outC)
	w.Close()
	out := <-outC

	expected := "Test info"
	n := level.Labels[level.Info]
	if !strings.Contains(out, expected) || !strings.Contains(out, n) {
		t.Errorf("Result `%s` doesn't contains `%s` and `%s`",
			out, expected, n)
	}
}

// TestInfof tests the Infof function.
func TestInfof(t *testing.T) {
	r, w, _ := os.Pipe()
	SetOutputs(Output{
		Name:   "test",
		Writer: w,
		Levels: level.Default,
	})

	Infof("Test %s", "infof")
	outC := make(chan string)
	go ioCopy(r, outC)
	w.Close()
	out := <-outC

	expected := "Test infof"
	n := level.Labels[level.Info]
	if !strings.Contains(out, expected) || !strings.Contains(out, n) {
		t.Errorf("Result `%s` doesn't contains `%s` and `%s`",
			out, expected, n)
	}
}

// TestInfoln tests the Infoln function.
func TestInfoln(t *testing.T) {
	r, w, _ := os.Pipe()
	SetOutputs(Output{
		Name:   "test",
		Writer: w,
		Levels: level.Default,
	})

	Infoln("Test infoln")
	outC := make(chan string)
	go ioCopy(r, outC)
	w.Close()
	out := <-outC

	expected := "Test infoln"
	n := level.Labels[level.Info]
	if !strings.Contains(out, expected) || !strings.Contains(out, n) {
		t.Errorf("Result `%s` doesn't contains `%s` and `%s`",
			out, expected, n)
	}
}

// TestFdebug tests the Fdebug function.
func TestFdebug(t *testing.T) {
	r, w, _ := os.Pipe()
	Fdebug(w, "Test Fdebug")
	outC := make(chan string)
	go ioCopy(r, outC)
	w.Close()
	out := <-outC

	expected := "Test Fdebug"
	n := level.Labels[level.Debug]
	if !strings.Contains(out, expected) || !strings.Contains(out, n) {
		t.Errorf("Result `%s` doesn't contains `%s` and `%s`",
			out, expected, n)
	}
}

// TestFdebugf tests the Fdebugf function.
func TestFdebugf(t *testing.T) {
	r, w, _ := os.Pipe()
	Fdebugf(w, "Test %s", "Fdebugf")
	outC := make(chan string)
	go ioCopy(r, outC)
	w.Close()
	out := <-outC

	expected := "Test Fdebugf"
	n := level.Labels[level.Debug]
	if !strings.Contains(out, expected) || !strings.Contains(out, n) {
		t.Errorf("Result `%s` doesn't contains `%s` and `%s`",
			out, expected, n)
	}
}

// TestFdebugln tests the Fdebugln function.
func TestFdebugln(t *testing.T) {
	r, w, _ := os.Pipe()
	Fdebugln(w, "Test Fdebugln")
	outC := make(chan string)
	go ioCopy(r, outC)
	w.Close()
	out := <-outC

	expected := "Test Fdebugln"
	n := level.Labels[level.Debug]
	if !strings.Contains(out, expected) || !strings.Contains(out, n) {
		t.Errorf("Result `%s` doesn't contains `%s` and `%s`",
			out, expected, n)
	}
}

// TestDebug tests the Debug function.
func TestDebug(t *testing.T) {
	r, w, _ := os.Pipe()
	SetOutputs(Output{
		Name:   "test",
		Writer: w,
		Levels: level.Default,
	})

	Debug("Test Debug")
	outC := make(chan string)
	go ioCopy(r, outC)
	w.Close()
	out := <-outC

	expected := "Test Debug"
	n := level.Labels[level.Debug]
	if !strings.Contains(out, expected) || !strings.Contains(out, n) {
		t.Errorf("Result `%s` doesn't contains `%s` and `%s`",
			out, expected, n)
	}
}

// TestDebugf tests the Debugf function.
func TestDebugf(t *testing.T) {
	r, w, _ := os.Pipe()
	SetOutputs(Output{
		Name:   "test",
		Writer: w,
		Levels: level.Default,
	})

	Debugf("Test %s", "Debugf")
	outC := make(chan string)
	go ioCopy(r, outC)
	w.Close()
	out := <-outC

	expected := "Test Debugf"
	n := level.Labels[level.Debug]
	if !strings.Contains(out, expected) || !strings.Contains(out, n) {
		t.Errorf("Result `%s` doesn't contains `%s` and `%s`",
			out, expected, n)
	}
}

// TestDebugln tests the Debugln function.
func TestDebugln(t *testing.T) {
	r, w, _ := os.Pipe()
	SetOutputs(Output{
		Name:   "test",
		Writer: w,
		Levels: level.Default,
	})

	Debugln("Test Debugln")
	outC := make(chan string)
	go ioCopy(r, outC)
	w.Close()
	out := <-outC

	expected := "Test Debugln"
	n := level.Labels[level.Debug]
	if !strings.Contains(out, expected) || !strings.Contains(out, n) {
		t.Errorf("Result `%s` doesn't contains `%s` and `%s`",
			out, expected, n)
	}
}

// TestFtrace tests the Ftrace function.
func TestFtrace(t *testing.T) {
	r, w, _ := os.Pipe()
	Ftrace(w, "Test Ftrace")
	outC := make(chan string)
	go ioCopy(r, outC)
	w.Close()
	out := <-outC

	expected := "Test Ftrace"
	n := level.Labels[level.Trace]
	if !strings.Contains(out, expected) || !strings.Contains(out, n) {
		t.Errorf("Result `%s` doesn't contains `%s` and `%s`",
			out, expected, n)
	}
}

// TestFtracef tests the Ftracef function.
func TestFtracef(t *testing.T) {
	r, w, _ := os.Pipe()
	Ftracef(w, "Test %s", "Ftracef")
	outC := make(chan string)
	go ioCopy(r, outC)
	w.Close()
	out := <-outC

	expected := "Test Ftracef"
	n := level.Labels[level.Trace]
	if !strings.Contains(out, expected) || !strings.Contains(out, n) {
		t.Errorf("Result `%s` doesn't contains `%s` and `%s`",
			out, expected, n)
	}
}

// TestFtraceln tests the Ftraceln function.
func TestFtraceln(t *testing.T) {
	r, w, _ := os.Pipe()
	Ftraceln(w, "Test Ftraceln")
	outC := make(chan string)
	go ioCopy(r, outC)
	w.Close()
	out := <-outC

	expected := "Test Ftraceln"
	n := level.Labels[level.Trace]
	if !strings.Contains(out, expected) || !strings.Contains(out, n) {
		t.Errorf("Result `%s` doesn't contains `%s` and `%s`",
			out, expected, n)
	}
}

// TestTrace tests the Trace function.
func TestTrace(t *testing.T) {
	r, w, _ := os.Pipe()
	SetOutputs(Output{
		Name:   "test",
		Writer: w,
		Levels: level.Default,
	})

	Trace("Test Trace")
	outC := make(chan string)
	go ioCopy(r, outC)
	w.Close()
	out := <-outC

	expected := "Test Trace"
	n := level.Labels[level.Trace]
	if !strings.Contains(out, expected) || !strings.Contains(out, n) {
		t.Errorf("Result `%s` doesn't contains `%s` and `%s`",
			out, expected, n)
	}
}

// TestTracef tests the Tracef function.
func TestTracef(t *testing.T) {
	r, w, _ := os.Pipe()
	SetOutputs(Output{
		Name:   "test",
		Writer: w,
		Levels: level.Default,
	})

	Tracef("Test %s", "Tracef")
	outC := make(chan string)
	go ioCopy(r, outC)
	w.Close()
	out := <-outC

	expected := "Test Tracef"
	n := level.Labels[level.Trace]
	if !strings.Contains(out, expected) || !strings.Contains(out, n) {
		t.Errorf("Result `%s` doesn't contains `%s` and `%s`",
			out, expected, n)
	}
}

// TestTraceln tests the Traceln function.
func TestTraceln(t *testing.T) {
	r, w, _ := os.Pipe()
	SetOutputs(Output{
		Name:   "test",
		Writer: w,
		Levels: level.Default,
	})

	Traceln("Test Traceln")
	outC := make(chan string)
	go ioCopy(r, outC)
	w.Close()
	out := <-outC

	expected := "Test Traceln"
	n := level.Labels[level.Trace]
	if !strings.Contains(out, expected) || !strings.Contains(out, n) {
		t.Errorf("Result `%s` doesn't contains `%s` and `%s`",
			out, expected, n)
	}
}
