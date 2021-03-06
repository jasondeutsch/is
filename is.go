package is

import (
	"fmt"
	"log"
	"reflect"
	"testing"
	"time"
)

// Is provides methods that leverage the existing testing capabilities found
// in the Go test framework. The methods provided allow for a more natural,
// efficient and expressive approach to writing tests. The goal is to write
// fewer lines of code while improving communication of intent.
type Is struct {
	TB         testing.TB
	strict     bool
	failFormat string
	failArgs   []interface{}
}

// New creates a new instance of the Is object and stores a reference to the
// provided testing object.
func New(tb testing.TB) *Is {
	if tb == nil {
		log.Fatalln("You must provide a testing object.")
	}
	return &Is{TB: tb, strict: true}
}

// New creates a new copy of your Is object and replaces the internal testing
// object with the provided testing object. This is useful for re-initializing
// your `is` instance inside a subtest so that it doesn't panic when using
// Strict mode.
//
// For example, creating your initial instance as such
//  is := is.New(t)
// is the convention, but this obviously shadows the `is` package namespace.
// Inside your subtest, you can do the exact same thing to initialize a locally scoped
// variable that uses the subtest's testing.T object.
func (is *Is) New(tb testing.TB) *Is {
	return &Is{
		TB:         tb,
		strict:     is.strict,
		failFormat: is.failFormat,
		failArgs:   is.failArgs,
	}
}

// Msg defines a message to print in the event of a failure. This allows you
// to print out additional information about a failure if it happens.
func (is *Is) Msg(format string, args ...interface{}) *Is {
	return &Is{
		TB:         is.TB,
		strict:     is.strict,
		failFormat: format,
		failArgs:   args,
	}
}

// AddMsg appends a message to print in the event of a failure. This allows
// you to build a failure message in multiple steps. If no message was
// previously set, simply sets the message.
//
// This method is most useful as a way of setting a default error message,
// then adding additional information to the output for specific assertions.
// For example:
//
// is := is.New(t).Msg("User ID: %d",u.ID)
// /*do things*/
// is.AddMsg("Raw Response: %s",body).Equal(res.StatusCode, http.StatusCreated)
func (is *Is) AddMsg(format string, args ...interface{}) *Is {
	if is.failFormat == "" {
		return is.Msg(format, args...)
	}
	return &Is{
		TB:         is.TB,
		strict:     is.strict,
		failFormat: fmt.Sprintf("%s - %s", is.failFormat, format),
		failArgs:   append(is.failArgs, args...),
	}
}

// Lax returns a copy of this instance of Is which does not abort the test if
// a failure occurs. Use this to run a set of tests and see all the failures
// at once.
func (is *Is) Lax() *Is {
	return &Is{
		TB:         is.TB,
		strict:     false,
		failFormat: is.failFormat,
		failArgs:   is.failArgs,
	}
}

// Strict returns a copy of this instance of Is which aborts the test if a
// failure occurs. This is the default behavior, thus this method has no
// effect unless it is used to reverse a previous call to Lax.
func (is *Is) Strict() *Is {
	return &Is{
		TB:         is.TB,
		strict:     true,
		failFormat: is.failFormat,
		failArgs:   is.failArgs,
	}
}

// Equal performs a deep compare of the provided objects and fails if they are
// not equal.
//
// Equal does not respect type differences. If the types are different and
// comparable (eg int32 and int64), they will be compared as though they are
// the same type.
func (is *Is) Equal(actual interface{}, expected interface{}) {
	is.TB.Helper()
	if !isEqual(actual, expected) {
		fail(is, "got %v (%s). expected %v (%s)",
			actual, objectTypeName(actual),
			expected, objectTypeName(expected))
	}
}

// NotEqual performs a deep compare of the provided objects and fails if they are
// equal.
//
// NotEqual does not respect type differences. If the types are different and
// comparable (eg int32 and int64), they will be compared as though they are
// the same type.
func (is *Is) NotEqual(a interface{}, b interface{}) {
	is.TB.Helper()
	if isEqual(a, b) {
		fail(is, "expected objects '%s' and '%s' not to be equal",
			objectTypeName(a),
			objectTypeName(b))
	}
}

// OneOf performs a deep compare of the provided object and an array of
// comparison objects. It fails if the first object is not equal to one of the
// comparison objects.
//
// OneOf does not respect type differences. If the types are different and
// comparable (eg int32 and int64), they will be compared as though they are
// the same type.
func (is *Is) OneOf(a interface{}, b ...interface{}) {
	is.TB.Helper()
	result := false
	for _, o := range b {
		result = isEqual(a, o)
		if result {
			break
		}
	}
	if !result {
		fail(is, "expected object '%s' to be equal to one of '%s', but got: %v and %v",
			objectTypeName(a),
			objectTypeNames(b), a, b)
	}
}

// NotOneOf performs a deep compare of the provided object and an array of
// comparison objects. It fails if the first object is equal to one of the
// comparison objects.
//
// NotOneOf does not respect type differences. If the types are different and
// comparable (eg int32 and int64), they will be compared as though they are
// the same type.
func (is *Is) NotOneOf(a interface{}, b ...interface{}) {
	is.TB.Helper()
	result := false
	for _, o := range b {
		result = isEqual(a, o)
		if result {
			break
		}
	}
	if result {
		fail(is, "expected object '%s' not to be equal to one of '%s', but got: %v and %v",
			objectTypeName(a),
			objectTypeNames(b), a, b)
	}
}

// Err checks the provided error object to determine if an error is present.
func (is *Is) Err(e error) {
	is.TB.Helper()
	if isNil(e) {
		fail(is, "expected error")
	}
}

// NotErr checks the provided error object to determine if an error is not
// present.
func (is *Is) NotErr(e error) {
	is.TB.Helper()
	if !isNil(e) {
		fail(is, "expected no error, but got: %v", e)
	}
}

// Nil checks the provided object to determine if it is nil.
func (is *Is) Nil(o interface{}) {
	is.TB.Helper()
	if !isNil(o) {
		fail(is, "expected object '%s' to be nil, but got: %v", objectTypeName(o), o)
	}
}

// NotNil checks the provided object to determine if it is not nil.
func (is *Is) NotNil(o interface{}) {
	is.TB.Helper()
	if isNil(o) {
		fail(is, "expected object '%s' not to be nil", objectTypeName(o))
	}
}

// True checks the provided boolean to determine if it is true.
func (is *Is) True(b bool) {
	is.TB.Helper()
	if !b {
		fail(is, "expected boolean to be true")
	}
}

// False checks the provided boolean to determine if is false.
func (is *Is) False(b bool) {
	is.TB.Helper()
	if b {
		fail(is, "expected boolean to be false")
	}
}

// Zero checks the provided object to determine if it is the zero value
// for the type of that object. The zero value is the same as what the object
// would contain when initialized but not assigned.
//
// This method, for example, would be used to determine if a string is empty,
// an array is empty or a map is empty. It could also be used to determine if
// a number is 0.
//
// In cases such as slice, map, array and chan, a nil value is treated the
// same as an object with len == 0
func (is *Is) Zero(o interface{}) {
	is.TB.Helper()
	if !isZero(o) {
		fail(is, "expected object '%s' to be zero value, but it was: %v", objectTypeName(o), o)
	}
}

// NotZero checks the provided object to determine if it is not the zero
// value for the type of that object. The zero value is the same as what the
// object would contain when initialized but not assigned.
//
// This method, for example, would be used to determine if a string is not
// empty, an array is not empty or a map is not empty. It could also be used
// to determine if a number is not 0.
//
// In cases such as slice, map, array and chan, a nil value is treated the
// same as an object with len == 0
func (is *Is) NotZero(o interface{}) {
	is.TB.Helper()
	if isZero(o) {
		fail(is, "expected object '%s' not to be zero value", objectTypeName(o))
	}
}

// Len checks the provided object to determine if it is the same length as the
// provided length argument.
//
// If the object is not one of type array, slice or map, it will fail.
func (is *Is) Len(o interface{}, l int) {
	is.TB.Helper()
	t := reflect.TypeOf(o)
	if o == nil ||
		(t.Kind() != reflect.Array &&
			t.Kind() != reflect.Slice &&
			t.Kind() != reflect.Map) {
		fail(is, "expected object '%s' to be of length '%d', but the object is not one of array, slice or map", objectTypeName(o), l)
		return
	}

	rLen := reflect.ValueOf(o).Len()
	if rLen != l {
		fail(is, "expected object '%s' to be of length '%d' but it was: %d", objectTypeName(o), l, rLen)
	}
}

// ShouldPanic expects the provided function to panic. If the function does
// not panic, this assertion fails.
func (is *Is) ShouldPanic(f func()) {
	is.TB.Helper()
	defer func() {
		r := recover()
		if r == nil {
			fail(is, "expected function to panic")
		}
	}()
	f()
}

// EqualType checks the type of the two provided objects and
// fails if they are not the same.
func (is *Is) EqualType(expected, actual interface{}) {
	is.TB.Helper()
	if reflect.TypeOf(expected) != reflect.TypeOf(actual) {
		fail(is, "expected objects '%s' to be of the same type as object '%s'", objectTypeName(expected), objectTypeName(actual))
	}
}

// WaitForTrue waits until the provided func returns true. If the timeout is
// reached before the function returns true, the test will fail.
func (is *Is) WaitForTrue(timeout time.Duration, f func() bool) {
	is.TB.Helper()
	after := time.After(timeout)
	for {
		select {
		case <-after:
			fail(is, "function did not return true within the timeout of %v", timeout)
			return
		default:
			if f() {
				return
			}
			time.Sleep(100 * time.Millisecond)
		}
	}
}
