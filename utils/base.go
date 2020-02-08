package utils

import (
	"github.com/astaxie/beego/logs"
	"fmt"
	"reflect"
	"unsafe"
)

func UnsafeBytesToString(bytes []byte) string {
	hdr := &reflect.StringHeader{
		Data: uintptr(unsafe.Pointer(&bytes[0])),
		Len:  len(bytes),
	}
	return *(*string)(unsafe.Pointer(hdr))
}

func UnsafeStringToBytes(str string) *[]byte {
	string_header := (*reflect.StringHeader)(unsafe.Pointer(&str))
	bytes_header := &reflect.SliceHeader{
		Data:string_header.Data,
		Len: string_header.Len,
		Cap: string_header.Len,
	}
	return (*[]byte)(unsafe.Pointer(bytes_header))
}

func UnsafeStringsToBytes(res *[][]byte, strs ...string) {
	//res := make([][]byte, len(strs))
	if len(*res) < len(strs) {
		*res = make([][]byte, len(strs))
	}
	for i, str := range strs {
		logs.Debug(str)
		(*res)[i] = *UnsafeStringToBytes(str)
		//res = append(res, UnsafeStringToBytes(str))
	}
	logs.Debug(res)
}

func ConcateString(bytes []byte,strs ...string) (res string) {
	var string_header *reflect.StringHeader
	var bytes_header *reflect.SliceHeader
	//index := 0
	for _,str := range strs {
		//string_header = (*reflect.StringHeader)(unsafe.Pointer(&str))
		//bytes_header = &reflect.SliceHeader{
		//	Data:string_header.Data,
		//	Len:string_header.Len,
		//	Cap:string_header.Len,
		//}
		//bytes = append(bytes, *(*[]byte)(unsafe.Pointer(bytes_header))...)
		bytes = append(bytes, str...)
		//temp := *(*[]byte)(unsafe.Pointer(bytes_header))
		//for _, c := range temp {
		//	bytes[index] = c
		//	index++
		//}
	}
	bytes_header = (*reflect.SliceHeader)(unsafe.Pointer(&bytes))
	string_header = &reflect.StringHeader{
		Data:bytes_header.Data,
		Len:bytes_header.Len,
	}
	res = *(*string)(unsafe.Pointer(string_header))
	return
}

func ConcateString2(strs ...string) (res string,clean_fn func()) {
	length := 0
	for _,str := range strs {
		length += len(str)
	}
	//bytes := make([]byte,length)
	bytes := GetBytes(length)
	clean_fn = func() {
		PutBytes(&bytes)
	}
	var string_header *reflect.StringHeader
	var bytes_header *reflect.SliceHeader
	for _,str := range strs {
		bytes = append(bytes, str...)
	}
	bytes_header = (*reflect.SliceHeader)(unsafe.Pointer(&bytes))
	string_header = &reflect.StringHeader{
		Data:bytes_header.Data,
		Len:bytes_header.Len,
	}
	res = *(*string)(unsafe.Pointer(string_header))
	return
}

func PrintBytesPointer(bytes *[]byte) {
	slice_header := (*reflect.SliceHeader)(unsafe.Pointer(bytes))
	fmt.Printf("address: %v ", slice_header.Data)
	logs.Warn(*bytes)
}

func PrintStringPointer(str *string) {
	string_header := (*reflect.StringHeader)(unsafe.Pointer(str))
	fmt.Printf("address: %v ", string_header.Data)
	logs.Warn(*str)
}