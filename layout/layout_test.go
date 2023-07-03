package layout

import (
	"errors"
	"testing"
)

// TestIsValid tests the IsValid method of the Layout type.
func TestIsValid(t *testing.T) {
	tests := []struct {
		name   string
		Layout Layout
		expect bool
	}{
		{
			name:   "FullFilePath Layout",
			Layout: FullFilePath,
			expect: true,
		},
		{
			name:   "ShortFilePath Layout",
			Layout: ShortFilePath,
			expect: true,
		},
		{
			name:   "FuncName Layout",
			Layout: FuncName,
			expect: true,
		},
		{
			name:   "Overflow Layout",
			Layout: overflowLayoutValue + 1,
			expect: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.Layout.IsValid(); got != test.expect {
				t.Errorf("IsValid = %v; want %v", got, test.expect)
			}
		})
	}
}

// TestIsSingle tests the IsSingle method of the Layout type.
func TestIsSingle(t *testing.T) {
	tests := []struct {
		name   string
		Layout Layout
		expect bool
	}{
		{
			name:   "FullFilePath Layout",
			Layout: FullFilePath,
			expect: true,
		},
		{
			name:   "Combined Layout",
			Layout: FullFilePath | ShortFilePath,
			expect: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.Layout.IsSingle(); got != test.expect {
				t.Errorf("IsSingle = %v; want %v", got, test.expect)
			}
		})
	}
}

// TestContains tests the Contains method of the Layout type.
func TestContains(t *testing.T) {
	tests := []struct {
		name   string
		Layout Layout
		flag   Layout
		expect bool
		err    error
	}{
		{
			name:   "FullFilePath Layout contains FullFilePath",
			Layout: FullFilePath,
			flag:   FullFilePath,
			expect: true,
			err:    nil,
		},
		{
			name:   "FullFilePath Layout contains FuncName",
			Layout: FullFilePath,
			flag:   FuncName,
			expect: false,
			err:    nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.Layout.Contains(test.flag)
			if err != test.err || got != test.expect {
				t.Errorf("Contains = %v, FuncName = %v; "+
					"want %v, FuncName = %v", got, err, test.expect, test.err)
			}
		})
	}
}

// TestLayoutFlags tests the methods that check for flags in the Layout type.
func TestLayoutFlags(t *testing.T) {
	tests := []struct {
		name    string
		Layout  Layout
		method  func(*Layout) bool
		contain bool
	}{
		{
			name:    "FullFilePath Flag in FullFilePath Layout",
			Layout:  FullFilePath,
			method:  (*Layout).FullFilePath,
			contain: true,
		},
		{
			name:    "FullFilePath Flag in FuncName Layout",
			Layout:  FuncName,
			method:  (*Layout).FullFilePath,
			contain: false,
		},
		{
			name:    "ShortFilePath Flag in ShortFilePath Layout",
			Layout:  ShortFilePath,
			method:  (*Layout).ShortFilePath,
			contain: true,
		},
		{
			name:    "ShortFilePath Flag in FuncName Layout",
			Layout:  FuncName,
			method:  (*Layout).ShortFilePath,
			contain: false,
		},
		{
			name:    "FuncName Flag in FuncName Layout",
			Layout:  FuncName,
			method:  (*Layout).FuncName,
			contain: true,
		},
		{
			name:    "FuncName Flag in FuncAddress Layout",
			Layout:  FuncAddress,
			method:  (*Layout).FuncName,
			contain: false,
		},
		{
			name:    "FuncAddress Flag in FuncAddress Layout",
			Layout:  FuncAddress,
			method:  (*Layout).FuncAddress,
			contain: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.method(&test.Layout); got != test.contain {
				t.Errorf("%s got = %v; want %v", test.name, got, test.contain)
			}
		})
	}
}

// TestSet tests the Set method of the Layout type.
func TestSet(t *testing.T) {
	tests := []struct {
		name   string
		Layout Layout
		flags  []Layout
		expect Layout
		err    error
	}{
		{
			name:   "Set FullFilePath and FuncName flags",
			Layout: Default,
			flags:  []Layout{FullFilePath, FuncName},
			expect: FullFilePath | FuncName,
			err:    nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.Layout.Set(test.flags...)
			if err != test.err || got != test.expect {
				t.Errorf("Set = %v, FuncName = %v; want %v, FuncName = %v",
					got, err, test.expect, test.err)
			}
		})
	}
}

// TestLayoutAdd tests the Add method of the Layout type.
func TestLayoutAdd(t *testing.T) {
	tests := []struct {
		name  string
		start Layout
		add   []Layout
		want  Layout
		err   error
	}{
		{
			name:  "Add FullFilePath to Empty Layout",
			start: 0,
			add:   []Layout{FullFilePath},
			want:  FullFilePath,
			err:   nil,
		},
		{
			name:  "Add FullFilePath to FullFilePath Layout",
			start: FullFilePath,
			add:   []Layout{FullFilePath},
			want:  FullFilePath,
			err:   nil,
		},
		{
			name:  "Add FuncName to FullFilePath Layout",
			start: FullFilePath,
			add:   []Layout{FuncName},
			want:  FullFilePath | FuncName,
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

// TestLayoutDelete tests the Delete method of the Layout type.
func TestLayoutDelete(t *testing.T) {
	tests := []struct {
		name   string
		start  Layout
		delete []Layout
		want   Layout
		err    error
	}{
		{
			name:   "Delete FullFilePath from FullFilePath Layout",
			start:  FullFilePath,
			delete: []Layout{FullFilePath},
			want:   0,
			err:    nil,
		},
		{
			name:   "Delete FullFilePath from FuncName Layout",
			start:  FuncName,
			delete: []Layout{FullFilePath},
			want:   FuncName,
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

// TestLayoutAll tests the All method of the Layout type.
func TestLayoutAll(t *testing.T) {
	tests := []struct {
		name  string
		start Layout
		all   []Layout
		want  bool
	}{
		{
			name:  "All FullFilePath and FuncName in FullFilePath Layout",
			start: FullFilePath,
			all:   []Layout{FullFilePath, FuncName},
			want:  false,
		},
		{
			name: "All FullFilePath and FuncName in " +
				"FullFilePath and FuncName Layout",
			start: FullFilePath | FuncName,
			all:   []Layout{FullFilePath, FuncName},
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

// TestLayoutAny tests the Any method of the Layout type.
func TestLayoutAny(t *testing.T) {
	tests := []struct {
		name  string
		start Layout
		any   []Layout
		want  bool
	}{
		{
			name:  "Any FullFilePath and FuncName in FullFilePath Layout",
			start: FullFilePath,
			any:   []Layout{FullFilePath, FuncName},
			want:  true,
		},
		{
			name:  "Any FullFilePath and FuncName in FuncAddress Layout",
			start: FuncAddress,
			any:   []Layout{FullFilePath, FuncName},
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
