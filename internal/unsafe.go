package internal

import "unsafe"

func UnsafeStringToBytes(s string) []byte {
	if s == "" {
		return nil
	}
	var p *byte = unsafe.StringData(s)
	return unsafe.Slice(p, len(s))
}

func UnsafeBytesToString(b []byte) string {
	if len(b) == 0 {
		return ""
	}
	var p *byte = unsafe.SliceData(b)
	return unsafe.String(p, len(b))
}

/* Usage:
 * var stack_mem [100]byte // or stack_mem = []byte{99:0}
 * var mem_ne = UnsafeNonEscape(stack_mem[:])
 * ....
 * doSomething(mem_ne)
 * ....
 * stack_mem[0] = 0 // or _ = stack_mem[0], can be sure that memory is not collected by gc
 */
func UnsafeNonEscape[T any](a []T) []T {
	var sliceHeader struct {
		Data uintptr
		Len  int
		Cap  int
	}

	sliceHeader.Data = uintptr(unsafe.Pointer(&a[0])) ^ 0
	sliceHeader.Len = len(a)
	sliceHeader.Cap = cap(a)

	return *(*[]T)(unsafe.Pointer(&sliceHeader))
}

// Hides a pointer from escape analysis.
// This was copied from the runtime and [strings.Builder].
//
//go:nosplit
//go:nocheckptr
func noescape(p unsafe.Pointer) unsafe.Pointer {
	var x = uintptr(p)
	return unsafe.Pointer(x ^ 0)
}

/* Usage:
 *   var stack_obj = some_struct{...}
 *   var p = UnsafeNonEscapePtr(&stack_obj)
 *   ...
 *   doSomething(p)
 *   ...
 */
func UnsafeNonEscapePtr[T any](a *T) *T {
	var p = noescape(unsafe.Pointer(a))
	return (*T)(p)
}

/* Usage:
 *   var stack_obj = some_struct{...}
 *   var p = UnsafeNonEscapePtr(&stack_obj)
 *   ...
 *   doSomething(p)
 *   ...
 */
func UnsafeNonEscapeConv[TSrc any, TRet any](a *TSrc) *TRet {
	var p = noescape(unsafe.Pointer(a))
	return (*TRet)(p)
}
