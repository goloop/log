package level

import (
	"errors"
	"testing"
)

// TestIsValid tests the IsValid method of the Level type.
func TestIsValid(t *testing.T) {
	tests := []struct {
		name   string
		level  Level
		expect bool
	}{
		{name: "Panic Level", level: Panic, expect: true},
		{name: "Fatal Level", level: Fatal, expect: true},
		{name: "Error Level", level: Error, expect: true},
		{name: "Overflow Level", level: overflowLevelValue + 1, expect: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.level.IsValid(); got != test.expect {
				t.Errorf("IsValid = %v; want %v", got, test.expect)
			}
		})
	}
}

// TestIsSingle tests the IsSingle method of the Level type.
func TestIsSingle(t *testing.T) {
	tests := []struct {
		name   string
		level  Level
		expect bool
	}{
		{name: "Panic Level", level: Panic, expect: true},
		{name: "Combined Level", level: Panic | Fatal, expect: false},
		// add more tests as needed
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.level.IsSingle(); got != test.expect {
				t.Errorf("IsSingle = %v; want %v", got, test.expect)
			}
		})
	}
}

// TestContains tests the Contains method of the Level type.
func TestContains(t *testing.T) {
	tests := []struct {
		name   string
		level  Level
		flag   Level
		expect bool
		err    error
	}{
		{
			name:   "Panic Level contains Panic",
			level:  Panic,
			flag:   Panic,
			expect: true,
			err:    nil,
		},
		{
			name:   "Panic Level contains Error",
			level:  Panic,
			flag:   Error,
			expect: false,
			err:    nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.level.Contains(test.flag)
			if err != test.err || got != test.expect {
				t.Errorf("Contains = %v, error = %v; want %v, error = %v",
					got, err, test.expect, test.err)
			}
		})
	}
}

// TestLevelFlags tests the level flags.
func TestLevelFlags(t *testing.T) {
	tests := []struct {
		name    string
		level   Level
		method  func(*Level) bool
		contain bool
	}{
		{
			name:    "Panic Flag in Panic Level",
			level:   Panic,
			method:  (*Level).Panic,
			contain: true,
		},
		{
			name:    "Panic Flag in Error Level",
			level:   Error,
			method:  (*Level).Panic,
			contain: false,
		},
		{
			name:    "Fatal Flag in Fatal Level",
			level:   Fatal,
			method:  (*Level).Fatal,
			contain: true,
		},
		{
			name:    "Fatal Flag in Error Level",
			level:   Error,
			method:  (*Level).Fatal,
			contain: false,
		},
		{
			name:    "Error Flag in Error Level",
			level:   Error,
			method:  (*Level).Error,
			contain: true,
		},
		{
			name:    "Error Flag in Info Level",
			level:   Info,
			method:  (*Level).Error,
			contain: false,
		},
		{
			name:    "Info Flag in Info Level",
			level:   Info,
			method:  (*Level).Info,
			contain: true,
		},
		{
			name:    "Info Flag in Debug Level",
			level:   Debug,
			method:  (*Level).Info,
			contain: false,
		},
		{
			name:    "Debug Flag in Debug Level",
			level:   Debug,
			method:  (*Level).Debug,
			contain: true,
		},
		{
			name:    "Debug Flag in Trace Level",
			level:   Trace,
			method:  (*Level).Debug,
			contain: false,
		},
		{
			name:    "Trace Flag in Trace Level",
			level:   Trace,
			method:  (*Level).Trace,
			contain: true,
		},
		{
			name:    "Trace Flag in Panic Level",
			level:   Panic,
			method:  (*Level).Trace,
			contain: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.method(&test.level); got != test.contain {
				t.Errorf("%s got = %v; want %v", test.name, got, test.contain)
			}
		})
	}
}

// TestSet tests the Set method of Level.
func TestSet(t *testing.T) {
	tests := []struct {
		name   string
		level  Level
		flags  []Level
		expect Level
		err    error
	}{
		{
			name:   "Set Panic and Error flags",
			level:  Default,
			flags:  []Level{Panic, Error},
			expect: Panic | Error,
			err:    nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.level.Set(test.flags...)
			if err != test.err || got != test.expect {
				t.Errorf("Set = %v, error = %v; want %v, error = %v",
					got, err, test.expect, test.err)
			}
		})
	}
}

func TestLevelAdd(t *testing.T) {
	tests := []struct {
		name  string
		start Level
		add   []Level
		want  Level
		err   error
	}{
		{
			name:  "Add Panic to Empty Level",
			start: 0,
			add:   []Level{Panic},
			want:  Panic,
			err:   nil,
		},
		{
			name:  "Add Panic to Panic Level",
			start: Panic,
			add:   []Level{Panic},
			want:  Panic,
			err:   nil,
		},
		{
			name:  "Add Error to Panic Level",
			start: Panic,
			add:   []Level{Error},
			want:  Panic | Error,
			err:   nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.start.Add(test.add...)
			if got != test.want || !errors.Is(err, test.err) {
				t.Errorf("%s got = %v, err = %v; want %v, %v",
					test.name, got, err, test.want, test.err)
			}
		})
	}
}

// TestLevelDelete tests the Delete method of Level.
func TestLevelDelete(t *testing.T) {
	tests := []struct {
		name   string
		start  Level
		delete []Level
		want   Level
		err    error
	}{
		{
			name:   "Delete Panic from Panic Level",
			start:  Panic,
			delete: []Level{Panic},
			want:   0,
			err:    nil,
		},
		{
			name:   "Delete Panic from Error Level",
			start:  Error,
			delete: []Level{Panic},
			want:   Error,
			err:    nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.start.Delete(test.delete...)
			if got != test.want || !errors.Is(err, test.err) {
				t.Errorf("%s got = %v, err = %v; want %v, %v",
					test.name, got, err, test.want, test.err)
			}
		})
	}
}

// TestLevelAll tests the All method of Level.
func TestLevelAll(t *testing.T) {
	tests := []struct {
		name  string
		start Level
		all   []Level
		want  bool
	}{
		{
			name:  "All Panic and Error in Panic Level",
			start: Panic,
			all:   []Level{Panic, Error},
			want:  false,
		},
		{
			name:  "All Panic and Error in Panic and Error Level",
			start: Panic | Error,
			all:   []Level{Panic, Error},
			want:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.start.All(test.all...); got != test.want {
				t.Errorf("%s got = %v; want %v", test.name, got, test.want)
			}
		})
	}
}

// TestLevelAny tests the Any method of Level.
func TestLevelAny(t *testing.T) {
	tests := []struct {
		name  string
		start Level
		any   []Level
		want  bool
	}{
		{
			name:  "Any Panic and Error in Panic Level",
			start: Panic,
			any:   []Level{Panic, Error},
			want:  true,
		},
		{
			name:  "Any Panic and Error in Info Level",
			start: Info,
			any:   []Level{Panic, Error},
			want:  false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.start.Any(test.any...); got != test.want {
				t.Errorf("%s got = %v; want %v", test.name, got, test.want)
			}
		})
	}
}
