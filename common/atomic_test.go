package common

import (
	"fmt"
	"testing"
)

func TestBool_CompareAndSwap(t *testing.T) {
	type args struct {
		old bool
		new bool
	}
	tests := []struct {
		name        string
		b           Bool
		args        args
		wantSwapped bool
	}{
		// 比较old值与b是否相等，若相等则用new替换b
		{name: "origin false compare swap false false", b: 0, args: args{old: false, new: false}, wantSwapped: true},
		{name: "origin false compare swap false true", b: 0, args: args{old: false, new: true}, wantSwapped: true},
		{name: "origin false compare swap true false", b: 0, args: args{old: true, new: false}, wantSwapped: false},
		{name: "origin false compare swap true true", b: 0, args: args{old: true, new: true}, wantSwapped: false},
		{name: "origin true compare swap false false", b: 1, args: args{old: false, new: false}, wantSwapped: false},
		{name: "origin true compare swap false true", b: 1, args: args{old: false, new: true}, wantSwapped: false},
		{name: "origin true compare swap true false", b: 1, args: args{old: true, new: false}, wantSwapped: true},
		{name: "origin true compare swap true true", b: 1, args: args{old: true, new: true}, wantSwapped: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotSwapped := tt.b.CompareAndSwap(tt.args.old, tt.args.new); gotSwapped != tt.wantSwapped {
				t.Errorf("CompareAndSwap() = %v, want %v", gotSwapped, tt.wantSwapped)
			}
		})
	}
}

func TestBool_Load(t *testing.T) {
	tests := []struct {
		name    string
		b       Bool
		wantVal bool
	}{
		{name: "load false", b: 0, wantVal: false},
		{name: "load true", b: 1, wantVal: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotVal := tt.b.Load(); gotVal != tt.wantVal {
				t.Errorf("Load() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestBool_Store(t *testing.T) {
	type args struct {
		val bool
	}
	tests := []struct {
		name string
		b    Bool
		args args
	}{
		{name: "false store false", b: 0, args: args{val: false}},
		{name: "false store true", b: 0, args: args{val: true}},
		{name: "true store false", b: 1, args: args{val: false}},
		{name: "true store true", b: 1, args: args{val: true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.Store(tt.args.val)
			if tt.b.Load() != tt.args.val {
				t.Errorf("Store() failed want %v , but got %v", tt.args.val, tt.b.Load())
			}
		})
	}
}

func TestBool_Swap(t *testing.T) {
	type args struct {
		new bool
	}
	tests := []struct {
		name    string
		b       Bool
		args    args
		wantOld bool
	}{
		{name: "false swap by true", b: 0, args: args{new: true}, wantOld: false},
		{name: "false swap by false", b: 0, args: args{new: false}, wantOld: false},
		{name: "true swap by true", b: 1, args: args{new: true}, wantOld: true},
		{name: "true swap by false", b: 1, args: args{new: false}, wantOld: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotOld := tt.b.Swap(tt.args.new); gotOld != tt.wantOld {
				t.Errorf("Swap() = %v, want %v", gotOld, tt.wantOld)
			}
		})
	}
}

func ExampleBool_Load() {
	var b Bool
	fmt.Println(b.Load())
	b.Store(true)
	fmt.Println(b.Load())

	// Output:
	// false
	// true
}

func ExampleBool_CompareAndSwap() {
	var b Bool
	fmt.Println(b.CompareAndSwap(true, true))
	fmt.Println(b.CompareAndSwap(false, true))

	// Output:
	// false
	// true
}

func ExampleBool_Store() {
	var b Bool
	b.Store(true)

	// Output:
}

func ExampleBool_Swap() {
	var b Bool
	fmt.Println(b.Swap(true))
	fmt.Println(b.Swap(true))

	// Output:
	// false
	// true
}
