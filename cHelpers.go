package oci8

// #include "oci8.go.h"
import "C"

import (
	"unsafe"
)

// getInt64 gets int64 from pointer
func getInt64(p unsafe.Pointer) int64 {
	return int64(*(*C.sb8)(p))
}

// getUint64 gets uint64 from pointer
func getUint64(p unsafe.Pointer) uint64 {
	return uint64(*(*C.sb8)(p))
}

// CByte comverts byte slice to oratext.
// must be freed
func CByte(b []byte) *C.oratext {
	p := C.malloc(C.size_t(len(b)))
	pp := (*[1 << 30]byte)(p)
	copy(pp[:], b)
	return (*C.oratext)(p)
}

// CByteN comverts byte slice to C oratext with size.
// must be freed
func CByteN(b []byte, size int) *C.oratext {
	p := C.malloc(C.size_t(size))
	pp := (*[1 << 30]byte)(p)
	copy(pp[:], b)
	return (*C.oratext)(p)
}

// CString coverts string to C oratext.
// must be freed
func CString(s string) *C.oratext {
	p := C.malloc(C.size_t(len(s) + 1))
	pp := (*[1 << 30]byte)(p)
	copy(pp[:], s)
	pp[len(s)] = 0
	return (*C.oratext)(p)
}

// CStringN coverts string to C oratext with size.
// must be freed
func CStringN(s string, size int) *C.oratext {
	p := C.malloc(C.size_t(size))
	pp := (*[1 << 30]byte)(p)
	copy(pp[:], s)
	if len(s) < size {
		pp[len(s)] = 0
	} else {
		pp[size-1] = 0
	}
	return (*C.oratext)(p)
}

// CGoStringN coverts C oratext to Go string
func CGoStringN(s *C.oratext, size int) string {
	if size == 0 {
		return ""
	}
	p := (*[1 << 30]byte)(unsafe.Pointer(s))
	buf := make([]byte, size)
	copy(buf, p[:])
	return *(*string)(unsafe.Pointer(&buf))
}

// freeDefines frees defines
func freeDefines(defines []oci8Define) {
	for _, define := range defines {
		if define.pbuf != nil {
			freeBuffer(define.pbuf, define.dataType)
			define.pbuf = nil
		}
		if define.length != nil {
			C.free(unsafe.Pointer(define.length))
			define.length = nil
		}
		if define.indicator != nil {
			C.free(unsafe.Pointer(define.indicator))
			define.indicator = nil
		}
		define.defineHandle = nil // should be freed by oci statment close
	}
}

// freeBinds frees binds
func freeBinds(binds []oci8Bind) {
	for _, bind := range binds {
		if bind.pbuf != nil {
			freeBuffer(bind.pbuf, bind.dataType)
			bind.pbuf = nil
		}
		if bind.length != nil {
			C.free(unsafe.Pointer(bind.length))
			bind.length = nil
		}
		if bind.indicator != nil {
			C.free(unsafe.Pointer(bind.indicator))
			bind.indicator = nil
		}
		bind.bindHandle = nil // freed by oci statment close
	}
}

// freeBuffer calles OCIDescriptorFree to free double pointer to buffer
// or calles C free to free pointer to buffer
func freeBuffer(buffer unsafe.Pointer, dataType C.ub2) {
	switch dataType {
	case C.SQLT_CLOB, C.SQLT_BLOB:
		C.OCIDescriptorFree(*(*unsafe.Pointer)(buffer), C.OCI_DTYPE_LOB)
	case C.SQLT_TIMESTAMP:
		C.OCIDescriptorFree(*(*unsafe.Pointer)(buffer), C.OCI_DTYPE_TIMESTAMP)
	case C.SQLT_TIMESTAMP_TZ:
		C.OCIDescriptorFree(*(*unsafe.Pointer)(buffer), C.OCI_DTYPE_TIMESTAMP_TZ)
	case C.SQLT_TIMESTAMP_LTZ:
		C.OCIDescriptorFree(*(*unsafe.Pointer)(buffer), C.OCI_DTYPE_TIMESTAMP_LTZ)
	case C.SQLT_INTERVAL_DS:
		C.OCIDescriptorFree(*(*unsafe.Pointer)(buffer), C.OCI_DTYPE_INTERVAL_DS)
	case C.SQLT_INTERVAL_YM:
		C.OCIDescriptorFree(*(*unsafe.Pointer)(buffer), C.OCI_DTYPE_INTERVAL_YM)
	default:
		C.free(buffer)
	}
}
