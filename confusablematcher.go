package confusablematcher

// #include "ConfusableMatcher/Export.h"
// #include <stdlib.h>
// #cgo LDFLAGS: -L. -lconfusablematcher
import "C"
import (
	"sync"
	"unsafe"
)

// KeyValue Key and value structure
type KeyValue struct {
	Key   string
	Value string
}

// CMHandle confusable matcher instance
type CMHandle struct {
	matcher    C.CMHandle
	ignoreList C.CMListHandle
	lock       sync.Mutex
}

// MappingResponse describes result returned by `AddMapping`
type MappingResponse int

const (
	// Success success
	Success MappingResponse = 0
	// AlreadyExists Key and value combination already exists
	AlreadyExists MappingResponse = 1
	// EmptyKey key is empty
	EmptyKey MappingResponse = 2
	// EmptyValue value is empty
	EmptyValue MappingResponse = 3
	// InvalidKey key starts with 0x00 or 0x01
	InvalidKey MappingResponse = 4
	// InvalidValue value starts with 0x00 or 0x01
	InvalidValue MappingResponse = 5
)

// InitConfusableMatcher Initializes new confusable matcher. If this instance is not used any more, `FreeConfusableMatcher` function must be called.
//
// Parameters:
//
// - `Map` : InputMap Input key to value mapping
// - `AddDefaultValues` : Whether to add default values or not ([a-z] -> [A-Z], [A-Z] -> [A-Z], [0-9] -> [0-9])
//
// Returns:
//
// - Handle to confusable matcher
func InitConfusableMatcher(InputMap []KeyValue, AddDefaultValues bool) CMHandle {
	var cmMap C.CMMap

	var tmp C.CMKV
	var structSz = int(unsafe.Sizeof(tmp))
	cmMap.Kv = (*C.CMKV)(C.malloc((C.ulonglong)(len(InputMap) * structSz * 5)))
	defer C.free(unsafe.Pointer(cmMap.Kv))

	for x, el := range InputMap {
		var cmKV C.CMKV
		cmKV.Key = C.CString(el.Key)
		defer C.free(unsafe.Pointer(cmKV.Key))
		cmKV.Value = C.CString(el.Value)
		defer C.free(unsafe.Pointer(cmKV.Value))

		var ptr = unsafe.Pointer(uintptr(unsafe.Pointer(cmMap.Kv)) + uintptr(x*structSz))
		*((*C.CMKV)(ptr)) = cmKV
	}

	cmMap.Size = C.uint(len(InputMap))

	var handle CMHandle

	var empty []string
	SetIgnoreList(&handle, empty)
	handle.matcher = C.InitConfusableMatcher(cmMap, (C.bool)(AddDefaultValues))
	return handle
}

// FreeConfusableMatcher Frees confusable matcher. Passed confusable matcher handle cannot be used after this method is called.
//
//
// Parameters:
//
// - `Matcher` : Confusable matcher handle
func FreeConfusableMatcher(Handle CMHandle) {
	C.FreeIgnoreList(Handle.ignoreList)
	C.FreeConfusableMatcher(Handle.matcher)
}

// SetIgnoreList sets an array of strings to ignore when performing an indexOf operation.
// These strings do not consume `contains` part of operation.
//
// Parameters:
//
// - `In` : Input array of strings to set to ignore
func SetIgnoreList(Handle *CMHandle, In []string) {
	var tmp *C.char
	var ptrSz = int(unsafe.Sizeof(&tmp))
	var list = (**C.char)(C.malloc((C.ulonglong)(len(In) * ptrSz)))
	defer C.free(unsafe.Pointer(list))

	for x, el := range In {
		var str = C.CString(el)
		defer C.free(unsafe.Pointer(str))

		var ptr = unsafe.Pointer(uintptr(unsafe.Pointer(list)) + uintptr(x*ptrSz))
		*((**C.char)(ptr)) = str
	}

	(*Handle).lock.Lock()
	{
		if (*Handle).ignoreList != nil {
			C.FreeIgnoreList((*Handle).ignoreList)
		}
		(*Handle).ignoreList = C.ConstructIgnoreList(list, (C.int)(len(In)))
	}
	(*Handle).lock.Unlock()
}

// IndexOf Performs an indexOf operation using specified mapping and ignore list
//
// Parameters:
//
// - `Matcher` : Handle to confusable matcher (returned by InitConfusableMatcher)
// - `In` : Input string
// - `Contains` : What input string should contain, aka the needle
// - `MatchRepeating` : Should it match repeating substrings in the mapping (without consuming the 'contains' portion of operation)
// - `StartIndex` : Starting index
//
// Returns:
//
// - Index and length
func IndexOf(Handle CMHandle, In string, Contains string, MatchRepeating bool, StartIndex int) (int, int) {
	var inPtr = C.CString(In)
	defer C.free(unsafe.Pointer(inPtr))
	var containsPtr = C.CString(Contains)
	defer C.free(unsafe.Pointer(containsPtr))

	var ret uint64
	Handle.lock.Lock()
	{
		ret = uint64(C.StringIndexOf(Handle.matcher, inPtr, containsPtr, (C.bool)(MatchRepeating), (C.int)(StartIndex), Handle.ignoreList))
	}
	Handle.lock.Unlock()

	return int(int32(ret & 0xFFFFFFFF)), int(int32(ret >> 32))
}

// AddMapping Adds a new key to value mapping into existing confusable matcher
//
// Parameters:
//
// - `Matcher` : Handle to confusable matcher (returned by InitConfusableMatcher)
// - `Key` : Input key
// - `Value` : Input value
// - `CheckValueDuplicate` : Check if key and value combination already exists
//
// - Operation result
func AddMapping(Handle CMHandle, Key string, Value string, CheckValueDuplicate bool) MappingResponse {
	var keyPtr = C.CString(Key)
	defer C.free(unsafe.Pointer(keyPtr))
	var valPtr = C.CString(Value)
	defer C.free(unsafe.Pointer(valPtr))

	return MappingResponse(C.AddMapping(Handle.matcher, keyPtr, valPtr, C.bool(CheckValueDuplicate)))
}

// RemoveMapping Removes an existing key to value mapping from confusable matcher
//
// Parameters:
//
// - `Key` : Input key
// - `Value` : Input value
//
// - If operation was successful or not
func RemoveMapping(Handle CMHandle, Key string, Value string) bool {
	var keyPtr = C.CString(Key)
	defer C.free(unsafe.Pointer(keyPtr))
	var valPtr = C.CString(Value)
	defer C.free(unsafe.Pointer(valPtr))

	return bool(C.RemoveMapping(Handle.matcher, keyPtr, valPtr))
}
