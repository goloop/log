package log

import (
	"os"
	"strings"
	"testing"

	"github.com/goloop/log/level"
	"github.com/goloop/trit"
)

// TestEcho tests the echo method of the Logger.
func TestEcho(t *testing.T) {
	logger := New()

	r, w, _ := os.Pipe()
	err := logger.SetOutputs(Output{
		Name:   "test",
		Writer: w,
		Levels: level.Default,
	})
	if err != nil {
		t.Fatal(err)
	}

	logger.echo(nil, level.Debug, "test %s", "message")
	outC := make(chan string)
	go ioCopy(r, outC)
	w.Close()
	out := <-outC

	if !strings.Contains(out, "test message") {
		t.Errorf("echo did not write the correct TEXT message: %s", out)
	}

	// As JSON.
	r, w, _ = os.Pipe()
	logger.SetOutputs(Output{
		Name:   "test",
		Writer: w,
		Levels: level.Default,
		Text:   trit.False,
	})

	logger.echo(nil, level.Debug, "test %s", "message")
	outC = make(chan string)
	go ioCopy(r, outC)
	w.Close()
	out = <-outC

	if !strings.Contains(out, "\"level\":\"DEBUG\"") {
		t.Errorf("echo did not write the correct JSON message: %s", out)
	}

	// Disabled.
	r, w, _ = os.Pipe()
	logger.SetOutputs(Output{
		Name:    "test",
		Writer:  w,
		Enabled: trit.False,
	})

	logger.echo(nil, level.Debug, "test %s", "message")
	outC = make(chan string)
	go ioCopy(r, outC)
	w.Close()
	out = <-outC

	if len(out) != 0 {
		t.Errorf("should not write anything: %s", out)
	}
}

/*

// TestSetSkipStackFramesMethod tests the SetSkipStackFrames
// method of the Logger.
func TestSetSkipStackFramesMethod(t *testing.T) {
	logger := New()
	skip := 5

	logger.SetSkipStackFrames(skip)
	if logger.skipStackFrames != skip {
		t.Errorf("SetSkipStackFrames failed, got %d, want %d",
			logger.skipStackFrames, skip)
	}
}

// TestSetPrefixMethod tests the SetPrefix method of the Logger.
func TestSetPrefixMethod(t *testing.T) {
	logger := New()
	prefix := "test"

	logger.SetPrefix(prefix)
	if logger.prefix != prefix {
		t.Errorf("SetPrefix failed, got %s, want %s", logger.prefix, prefix)
	}
}

// TestSetOutputsMethod tests the SetOutputs method of the Logger.
func TestSetOutputsMethod(t *testing.T) {
	logger := New()
	output := Output{
		Name:   "test",
		Writer: os.Stdout,
	}

	err := logger.SetOutputs(output)
	if err != nil {
		t.Errorf("SetOutputs failed with error: %v", err)
	}

	if len(logger.outputs) != 1 {
		t.Errorf("SetOutputs failed, got %d outputs, want 1",
			len(logger.outputs))
	}

	if _, ok := logger.outputs["test"]; !ok {
		t.Errorf("SetOutputs failed, did not find output with name 'test'")
	}
}

// TestSetOutputsMethodNoName tests the SetOutputs method with
// an output with no name.
func TestSetOutputsMethodNoName(t *testing.T) {
	logger := New()
	output := Output{
		Writer: os.Stdout,
	}

	err := logger.SetOutputs(output)
	if err == nil {
		t.Errorf("SetOutputs did not fail with error for output with no name")
	}
}

// TestSetOutputsMethodDuplicateName tests the SetOutputs method
// with duplicate output names.
func TestSetOutputsMethodDuplicateName(t *testing.T) {
	logger := New()
	output1 := Output{
		Name:   "test",
		Writer: os.Stdout,
	}
	output2 := Output{
		Name:   "test",
		Writer: os.Stdout,
	}

	err := logger.SetOutputs(output1, output2)
	if err == nil {
		t.Errorf("SetOutputs didn't fail with error for dupl. output names")
	}
}

// TestSetOutputsMethodNilWriter tests the SetOutputs method with
// an output with a nil writer and empty list.
func TestSetOutputsMethodNil(t *testing.T) {
	logger := New()
	output := Output{
		Name: "test",
	}

	err := logger.SetOutputs(output)
	if err == nil {
		t.Errorf("SetOutputs didn't fail with error " +
			"for output with nil writer")
	}

	err = logger.SetOutputs()
	if err == nil {
		t.Errorf("SetOutputs didn't fail with error " +
			"for empty list of outputs")
	}
}

// TestEditOutputsMethod tests the EditOutputs method of the Logger.
func TestEditOutputsMethod(t *testing.T) {
	logger := New()
	output := Output{
		Name:   "test",
		Writer: os.Stdout,
		Levels: level.Error,
	}

	err := logger.SetOutputs(output)
	if err != nil {
		t.Fatal(err)
	}

	output.Levels = level.Debug
	err = logger.EditOutputs(output)
	if err != nil {
		t.Errorf("EditOutputs failed with error: %v", err)
	}

	if logger.outputs["test"].Levels != level.Debug {
		t.Errorf("EditOutputs did not update the output levels correctly")
	}
}

// TestDeleteOutputsMethod tests the DeleteOutputs method of the Logger.
func TestDeleteOutputsMethod(t *testing.T) {
	logger := New()
	output := Output{
		Name:   "test",
		Writer: os.Stdout,
	}

	err := logger.SetOutputs(output)
	if err != nil {
		t.Fatal(err)
	}

	logger.DeleteOutputs(output.Name)

	if _, ok := logger.outputs["test"]; ok {
		t.Errorf("DeleteOutputs did not delete the output")
	}
}

// TestOutputsMethod tests the Outputs method of the Logger.
func TestOutputsMethod(t *testing.T) {
	logger := New()
	output := Output{
		Name:   "test",
		Writer: os.Stdout,
	}

	err := logger.SetOutputs(output)
	if err != nil {
		t.Fatal(err)
	}

	outputs := logger.Outputs()

	if len(outputs) != 1 {
		t.Errorf("Outputs did not return the correct number of outputs")
	}

	if outputs[0].Name != output.Name {
		t.Errorf("Outputs did not return the correct output")
	}
}

// TestCopyMethod tests the Copy method of the Logger.
func TestCopyMethod(t *testing.T) {
	logger := New()
	logger.SetPrefix("test")
	copy := logger.Copy()
	if copy.prefix != "test" {
		t.Errorf("Copy did not correctly copy the prefix")
	}
}



// TestPanicMethods tests the Panic, Panicf, and Panicln methods of the Logger.
func TestPanicMethods(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic as expected")
		}
	}()

	logger := New()
	_, w, _ := os.Pipe()
	logger.SetOutputs(Output{
		Name:   "test",
		Writer: w,
		Levels: level.Default,
	})

	logger.Panic("Test panic")

	// Panicf
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic as expected")
		}
	}()
	logger.Panicf("Test panicf %s", "formatted")

	// Panicln
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic as expected")
		}
	}()
	logger.Panicln("Test panicln")
}
*/
// //////////////////////////////////////////////////////////////////////////////////////////////////////////////

/*
// Fatal methods (like Fatal, Fatalf, Fatalln) exit the program,
// which makes them difficult to test directly. They should be covered
// in an integration test, where you can check that a program with a fatal
// log statement exits as expected.
//
// For unit testing, we can at least test that it logs as expected
// before calling os.Exit.
//
// TestFatalMethods tests the Fatal, Fatalf, and Fatalln methods of the Logger.
func TestFatalMethods(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic as expected")
		}
	}()

	logger := New()

	r, w, _ := os.Pipe()
	logger.SetOutputs(Output{
		Name:   "test",
		Writer: w,
		Levels: level.Default,
	})

	logger.Fatal("Test fatal")
	outC := make(chan string)
	go ioCopy(r, outC)
	w.Close()
	out := <-outC

	expected := "Test fatal"
	if !strings.Contains(out, expected) || !strings.Contains(out, "FATAL") {
		t.Errorf("Result `%s` doesn't contains `%s` and `%s`",
            out, expected, n)
	}
}

// TestFatalfMethod tests the Fatalf method of the Logger.
func TestFatalfMethod(t *testing.T) {
	logger := New()

	r, w, _ := os.Pipe()
	logger.SetOutputs(Output{
		Name:   "test",
		Writer: w,
		Levels: level.Default,
	})

	logger.Fatalf("Test fatal %s", "formatted")
	outC := make(chan string)
	go ioCopy(r, outC)
	w.Close()
	out := <-outC

	expected := "Test fatal formatted"
	if !strings.Contains(out, expected) || !strings.Contains(out, "FATAL") {
		t.Errorf("Result `%s` doesn't contains `%s` and `%s`",
            out, expected, n)
	}
}

// TestFatallnMethod tests the Fatalln method of the Logger.
func TestFatallnMethod(t *testing.T) {
	logger := New()

	r, w, _ := os.Pipe()
	logger.SetOutputs(Output{
		Name:   "test",
		Writer: w,
		Levels: level.Default,
	})

	logger.Fatalln("Test fatalln")
	outC := make(chan string)
	go ioCopy(r, outC)
	w.Close()
	out := <-outC

	expected := "Test fatalln"
	if !strings.Contains(out, expected) || !strings.Contains(out, "FATAL") {
		t.Errorf("Result `%s` doesn't contains `%s` and `%s`",
            out, expected, n)
	}
}
*/

// //////////////////////////////////////////////////////////////////////////////////////////////////////////////

/*
// TestErrorMethod tests the Error method of the Logger.
func TestErrorMethod(t *testing.T) {
	logger := New()

	r, w, _ := os.Pipe()
	logger.SetOutputs(Output{
		Name:   "test",
		Writer: w,
		Levels: level.Default,
	})

	logger.Error("Test Error")
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

// TestErrorfMethod tests the Errorf method of the Logger.
func TestErrorfMethod(t *testing.T) {
	logger := New()

	r, w, _ := os.Pipe()
	logger.SetOutputs(Output{
		Name:   "test",
		Writer: w,
		Levels: level.Default,
	})

	logger.Errorf("Test %s", "Errorf")
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

// TestErrorlnMethod tests the Errorln method of the Logger.
func TestErrorlnMethod(t *testing.T) {
	logger := New()

	r, w, _ := os.Pipe()
	logger.SetOutputs(Output{
		Name:   "test",
		Writer: w,
		Levels: level.Default,
	})

	logger.Errorln("Test Errorln")
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

// TestFerrorMethod tests the Ferror method of the Logger.
func TestFerrorMethod(t *testing.T) {
	logger := New()

	r, w, _ := os.Pipe()
	logger.SetOutputs(Output{
		Name:   "test",
		Writer: w,
		Levels: level.Default,
	})

	logger.Ferror(w, "Test Ferror")
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

// TestFerrorfMethod tests the Ferrorf method of the Logger.
func TestFerrorfMethod(t *testing.T) {
	logger := New()

	r, w, _ := os.Pipe()
	logger.SetOutputs(Output{
		Name:   "test",
		Writer: w,
		Levels: level.Default,
	})

	logger.Ferrorf(w, "Test %s", "Ferrorf")
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

// TestFerrorlnMethod tests the Ferrorln method of the Logger.
func TestFerrorlnMethod(t *testing.T) {
	logger := New()

	r, w, _ := os.Pipe()
	logger.SetOutputs(Output{
		Name:   "test",
		Writer: w,
		Levels: level.Default,
	})

	logger.Ferrorln(w, "Test Ferrorln")
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

// TestFwarnMethod tests the Fwarn method of the Logger.
func TestFwarnMethod(t *testing.T) {
	logger := New()

	r, w, _ := os.Pipe()
	logger.SetOutputs(Output{
		Name:   "test",
		Writer: w,
		Levels: level.Default,
	})

	logger.Fwarn(w, "Test Fwarn")
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

// TestFwarnfMethod tests the Fwarnf method of the Logger.
func TestFwarnfMethod(t *testing.T) {
	logger := New()

	r, w, _ := os.Pipe()
	logger.SetOutputs(Output{
		Name:   "test",
		Writer: w,
		Levels: level.Default,
	})

	logger.Fwarnf(w, "Test %s", "Fwarnf")
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

// TestFwarnlnMethod tests the Fwarnln method of the Logger.
func TestFwarnlnMethod(t *testing.T) {
	logger := New()

	r, w, _ := os.Pipe()
	logger.SetOutputs(Output{
		Name:   "test",
		Writer: w,
		Levels: level.Default,
	})

	logger.Fwarnln(w, "Test Fwarnln")
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

// TestWarnfMethod tests the Warnf method of the Logger.
func TestWarnfMethod(t *testing.T) {
	logger := New()

	r, w, _ := os.Pipe()
	logger.SetOutputs(Output{
		Name:   "test",
		Writer: w,
		Levels: level.Default,
	})

	logger.Warnf("Test warning %s", "formatted")
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

// TestWarnlnMethod tests the Warnln method of the Logger.
func TestWarnlnMethod(t *testing.T) {
	logger := New()

	r, w, _ := os.Pipe()
	logger.SetOutputs(Output{
		Name:   "test",
		Writer: w,
		Levels: level.Default,
	})

	logger.Warnln("Test warnln")
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

// TestFinfoMethod tests the Finfo method of the Logger.
func TestFinfoMethod(t *testing.T) {
	logger := New()

	r, w, _ := os.Pipe()
	logger.SetOutputs(Output{
		Name:   "test",
		Writer: w,
		Levels: level.Default,
	})

	logger.Finfo(w, "Test finfo")
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

// TestFinfofMethod tests the Finfof method of the Logger.
func TestFinfofMethod(t *testing.T) {
	logger := New()

	r, w, _ := os.Pipe()
	logger.SetOutputs(Output{
		Name:   "test",
		Writer: w,
		Levels: level.Default,
	})

	logger.Finfof(w, "Test %s", "finfof")
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

// TestFinfolnMethod tests the Finfoln method of the Logger.
func TestFinfolnMethod(t *testing.T) {
	logger := New()

	r, w, _ := os.Pipe()
	logger.SetOutputs(Output{
		Name:   "test",
		Writer: w,
		Levels: level.Default,
	})

	logger.Finfoln(w, "Test finfoln")
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

// TestInfoMethod tests the Info method of the Logger.
func TestInfoMethod(t *testing.T) {
	logger := New()

	r, w, _ := os.Pipe()
	logger.SetOutputs(Output{
		Name:   "test",
		Writer: w,
		Levels: level.Default,
	})

	logger.Info("Test info")
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

// TestInfofMethod tests the Infof method of the Logger.
func TestInfofMethod(t *testing.T) {
	logger := New()

	r, w, _ := os.Pipe()
	logger.SetOutputs(Output{
		Name:   "test",
		Writer: w,
		Levels: level.Default,
	})

	logger.Infof("Test %s", "infof")
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

// TestInfolnMethod tests the Infoln method of the Logger.
func TestInfolnMethod(t *testing.T) {
	logger := New()

	r, w, _ := os.Pipe()
	logger.SetOutputs(Output{
		Name:   "test",
		Writer: w,
		Levels: level.Default,
	})

	logger.Infoln("Test infoln")
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

// TestFdebugMethod tests the Fdebug method of the Logger.
func TestFdebugMethod(t *testing.T) {
	logger := New()

	r, w, _ := os.Pipe()
	logger.SetOutputs(Output{
		Name:   "test",
		Writer: w,
		Levels: level.Default,
	})

	logger.Fdebug(w, "Test Fdebug")
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

// TestFdebugfMethod tests the Fdebugf method of the Logger.
func TestFdebugfMethod(t *testing.T) {
	logger := New()

	r, w, _ := os.Pipe()
	logger.SetOutputs(Output{
		Name:   "test",
		Writer: w,
		Levels: level.Default,
	})

	logger.Fdebugf(w, "Test %s", "Fdebugf")
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

// TestFdebuglnMethod tests the Fdebugln method of the Logger.
func TestFdebuglnMethod(t *testing.T) {
	logger := New()

	r, w, _ := os.Pipe()
	logger.SetOutputs(Output{
		Name:   "test",
		Writer: w,
		Levels: level.Default,
	})

	logger.Fdebugln(w, "Test Fdebugln")
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

// TestDebugMethod tests the Debug method of the Logger.
func TestDebugMethod(t *testing.T) {
	logger := New()

	r, w, _ := os.Pipe()
	logger.SetOutputs(Output{
		Name:   "test",
		Writer: w,
		Levels: level.Default,
	})

	logger.Debug("Test Debug")
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

// TestDebugfMethod tests the Debugf method of the Logger.
func TestDebugfMethod(t *testing.T) {
	logger := New()

	r, w, _ := os.Pipe()
	logger.SetOutputs(Output{
		Name:   "test",
		Writer: w,
		Levels: level.Default,
	})

	logger.Debugf("Test %s", "Debugf")
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

// TestDebuglnMethod tests the Debugln method of the Logger.
func TestDebuglnMethod(t *testing.T) {
	logger := New()

	r, w, _ := os.Pipe()
	logger.SetOutputs(Output{
		Name:   "test",
		Writer: w,
		Levels: level.Default,
	})

	logger.Debugln("Test Debugln")
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

// TestFtraceMethod tests the Ftrace method of the Logger.
func TestFtraceMethod(t *testing.T) {
	logger := New()

	r, w, _ := os.Pipe()
	logger.SetOutputs(Output{
		Name:   "test",
		Writer: w,
		Levels: level.Default,
	})

	logger.Ftrace(w, "Test Ftrace")
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

// TestFtracefMethod tests the Ftracef method of the Logger.
func TestFtracefMethod(t *testing.T) {
	logger := New()

	r, w, _ := os.Pipe()
	logger.SetOutputs(Output{
		Name:   "test",
		Writer: w,
		Levels: level.Default,
	})

	logger.Ftracef(w, "Test %s", "Ftracef")
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

// TestFtracelnMethod tests the Ftraceln method of the Logger.
func TestFtracelnMethod(t *testing.T) {
	logger := New()

	r, w, _ := os.Pipe()
	logger.SetOutputs(Output{
		Name:   "test",
		Writer: w,
		Levels: level.Default,
	})

	logger.Ftraceln(w, "Test Ftraceln")
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

// TestTraceMethod tests the Trace method of the Logger.
func TestTraceMethod(t *testing.T) {
	logger := New()

	r, w, _ := os.Pipe()
	logger.SetOutputs(Output{
		Name:   "test",
		Writer: w,
		Levels: level.Default,
	})

	logger.Trace("Test Trace")
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

// TestTracefMethod tests the Tracef method of the Logger.
func TestTracefMethod(t *testing.T) {
	logger := New()

	r, w, _ := os.Pipe()
	logger.SetOutputs(Output{
		Name:   "test",
		Writer: w,
		Levels: level.Default,
	})

	logger.Tracef("Test %s", "Tracef")
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

// TestTracelnMethod tests the Traceln method of the Logger.
func TestTracelnMethod(t *testing.T) {
	logger := New()

	r, w, _ := os.Pipe()
	logger.SetOutputs(Output{
		Name:   "test",
		Writer: w,
		Levels: level.Default,
	})

	logger.Traceln("Test Traceln")
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
*/
