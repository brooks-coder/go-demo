package common

import (
	"container/list"
	"fmt"
	"reflect"
	"testing"
)

func TestLRUCache_Add(t *testing.T) {
	type fields struct {
		max  int
		Call func(key interface{}, value interface{})
	}
	type args struct {
		key   interface{}
		value interface{}
	}
	callFunc := func(key interface{}, value interface{}) { fmt.Println("key:", key, "value", value) }

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{name: "normal", fields: fields{max: 3}, args: args{key: "1", value: 1}, wantErr: false},
		{name: "no-limit", fields: fields{max: 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewLRUCache(tt.fields.max, callFunc)
			if err := c.Add(tt.args.key, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Add() error = %v, wantErr %v", err, tt.wantErr)
			}
			c.Add("2",2)
			c.Add("3",3)
			c.Add("4",4)
		})
	}
}

func TestLRUCache_Get(t *testing.T) {
	type fields struct {
		max   int
		cache map[interface{}]*list.Element
		Call  func(key interface{}, value interface{})
		l     *list.List
	}
	type args struct {
		key interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   interface{}
		want1  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &LRUCache{
				max:   tt.fields.max,
				cache: tt.fields.cache,
				Call:  tt.fields.Call,
				l:     tt.fields.l,
			}
			got, got1 := c.Get(tt.args.key)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Get() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestLRUCache_Len(t *testing.T) {
	type fields struct {
		max   int
		cache map[interface{}]*list.Element
		Call  func(key interface{}, value interface{})
		l     *list.List
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &LRUCache{
				max:   tt.fields.max,
				cache: tt.fields.cache,
				Call:  tt.fields.Call,
				l:     tt.fields.l,
			}
			if got := c.Len(); got != tt.want {
				t.Errorf("Len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLRUCache_Remove(t *testing.T) {
	type fields struct {
		max   int
		cache map[interface{}]*list.Element
		Call  func(key interface{}, value interface{})
		l     *list.List
	}
	type args struct {
		key interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &LRUCache{
				max:   tt.fields.max,
				cache: tt.fields.cache,
				Call:  tt.fields.Call,
				l:     tt.fields.l,
			}
			if err := c.Remove(tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("Remove() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
