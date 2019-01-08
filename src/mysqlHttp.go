package main

// #include <stdio.h>
// #include <stdlib.h>
// #include <string.h>
// #include <mysql.h>
// #cgo CFLAGS: -I/usr/include/mysql
import "C"

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"unsafe"
)

func msg(message *C.char, s string) {
	m := C.CString(s)
	defer C.free(unsafe.Pointer(m))

	C.strcpy(message, m)
}

func argToGostrings(length C.uint, args **C.char) []string {
	argslice := (*[1 << 30]*C.char)(unsafe.Pointer(args))[:length:length]
	gostrings := make([]string, length)

	for i, s := range argslice {
		gostrings[i] = C.GoString(s)
	}

	return gostrings
}

//export jsonObject_init
func jsonObject_init(initid *C.UDF_INIT, args *C.UDF_ARGS, message *C.char) C.my_bool {
	return 0
}

//export jsonObject
func jsonObject(initid *C.UDF_INIT, args *C.UDF_ARGS, result *C.char, length *C.ulong, isNull *C.char, error *C.char) *C.char {
	data := ""

	if int(args.arg_count) % 2 == 0 {
		argsString := argToGostrings(args.arg_count, args.args)
		hashData := map[string]interface{}{}

		for k, v := range argsString {
			if k > 0 && k % 2 != 0 {
				hashData[argsString[k-1]] = v
			}
		}

		jsonData, err := json.Marshal(hashData)
		if err == nil {
			data = string(jsonData)
		}
	}

	result = C.CString(data)
	*length = C.ulong(len(data))

	return result
}

//export jsonObject_deinit
func jsonObject_deinit(initid *C.UDF_INIT) {
}

//export httpPost_init
func httpPost_init(initid *C.UDF_INIT, args *C.UDF_ARGS, message *C.char) C.my_bool {
	if int(args.arg_count) != 2 {
		msg(message, "error params")
		return 1
	}

	return 0
}

//export httpPost
func httpPost(initid *C.UDF_INIT, args *C.UDF_ARGS, result *C.char, length *C.ulong, isNull *C.char, error *C.char) *C.char {
	argsString := argToGostrings(args.arg_count, args.args)

	data := ""

	response, err := http.Post(argsString[0], "application/octet-stream", strings.NewReader(argsString[1]))
	if err != nil {
		data = err.Error()
	} else {
		defer response.Body.Close()
		content, err := ioutil.ReadAll(response.Body)
		if err != nil {
			data = err.Error()
		} else {
			data = string(content)
		}
	}

	result = C.CString(data)
	*length = C.ulong(len(data))

	return result
}

//export httpPost_deinit
func httpPost_deinit(initid *C.UDF_INIT) {
}

func main()  {}
