package is

import (
	"errors"
	"fmt"
	"reflect"

	"testing"
)

var numberTypes = []reflect.Type{
	reflect.TypeOf(int(0)),
	reflect.TypeOf(int8(0)),
	reflect.TypeOf(int16(0)),
	reflect.TypeOf(int32(0)),
	reflect.TypeOf(int64(0)),
	reflect.TypeOf(uint(0)),
	reflect.TypeOf(uint8(0)),
	reflect.TypeOf(uint16(0)),
	reflect.TypeOf(uint32(0)),
	reflect.TypeOf(uint64(0)),
	reflect.TypeOf(float32(0)),
	reflect.TypeOf(float64(0)),
}

type testStruct struct {
	v int
}

var tests = []struct {
	a      interface{}
	b      interface{}
	c      interface{}
	d      interface{}
	e      interface{}
	cTypes []reflect.Type
}{
	{
		a:      0,
		b:      0,
		c:      1,
		d:      0,
		e:      1,
		cTypes: numberTypes,
	},
	{
		a: "test",
		b: "test",
		c: "testing",
		d: "",
		e: "testing",
	},
	{
		a: struct{}{},
		b: struct{}{},
		c: struct{ v int }{v: 1},
		d: testStruct{},
		e: testStruct{v: 1},
	},
	{
		a: &struct{}{},
		b: &struct{}{},
		c: &struct{ v int }{v: 1},
		d: &testStruct{},
		e: &testStruct{v: 1},
	},
	{
		a: []int64{0, 1},
		b: []int64{0, 1},
		c: []int64{0, 2},
		d: []int64{},
		e: []int64{0, 2},
	},
	{
		a: map[string]int64{"answer": 42},
		b: map[string]int64{"answer": 42},
		c: map[string]int64{"answer": 43},
		d: map[string]int64{},
		e: map[string]int64{"answer": 42},
	},
	{
		a: true,
		b: true,
		c: false,
		d: false,
		e: true,
	},
}

func TestIs(t *testing.T) {
	is := New(t)

	for i, test := range tests {
		for _, cType := range test.cTypes {
			is.fail = func(format string, args ...interface{}) {
				fmt.Print(decorate(fmt.Sprintf(fmt.Sprintf("(test #%d) - ", i)+format, args...)))
				t.FailNow()
			}
			is.Equal(test.a, reflect.ValueOf(test.b).Convert(cType).Interface())
		}
		is.Equal(test.a, test.b)
	}

	for i, test := range tests {
		for _, cType := range test.cTypes {
			is.fail = func(format string, args ...interface{}) {
				fmt.Print(decorate(fmt.Sprintf(fmt.Sprintf("(test #%d) - ", i)+format, args...)))
				t.FailNow()
			}
			is.NotEqual(test.a, reflect.ValueOf(test.c).Convert(cType).Interface())
		}
		is.NotEqual(test.a, test.c)
	}

	for i, test := range tests {
		for _, cType := range test.cTypes {
			is.fail = func(format string, args ...interface{}) {
				fmt.Print(decorate(fmt.Sprintf(fmt.Sprintf("(test #%d) - ", i)+format, args...)))
				t.FailNow()
			}
			is.Zero(reflect.ValueOf(test.d).Convert(cType).Interface())
		}
		is.Zero(test.d)
	}

	for i, test := range tests {
		for _, cType := range test.cTypes {
			is.fail = func(format string, args ...interface{}) {
				fmt.Print(decorate(fmt.Sprintf(fmt.Sprintf("(test #%d) - ", i)+format, args...)))
				t.FailNow()
			}
			is.NotZero(reflect.ValueOf(test.e).Convert(cType).Interface())
		}
		is.NotZero(test.e)
	}

	is.fail = func(format string, args ...interface{}) {
		fmt.Print(decorate(fmt.Sprintf(format, args...)))
		t.FailNow()
	}
	is.Nil(nil)
	is.NotNil(&testStruct{v: 1})
	is.Err(errors.New("error"))
	is.NotErr(nil)
	is.True(true)
	is.False(false)
}

type testAfter struct {
	success bool
	calls   int
}

// ensure after implements After interface
var _ After = (*testAfter)(nil)

func (a *testAfter) Msg(format string, args ...interface{}) {
	if !a.success {
		a.calls++
	}
}

func TestIsAfter(t *testing.T) {
	is := New(t)
	is.fail = func(format string, args ...interface{}) {
	}
	newAfter = func(success bool) After {
		return &testAfter{success: success}
	}

	a := is.True(false)
	a.Msg("testing after")

	if a.(*testAfter).calls == 0 {
		t.Logf("should have called the Msg method after failure")
	}

	a = is.True(true)
	a.Msg("testing after")
	if a.(*testAfter).calls == 1 {
		t.Logf("should not have called the Msg method after failure")
	}

}
