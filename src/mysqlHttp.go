package main

// #include <stdio.h>
// #include <stdlib.h>
// #include <string.h>
// #include <mysql.h>
// #cgo CFLAGS: -I/usr/include/mysql
import "C"

import (
	"io/ioutil"
	"net/http"
	"strings"
	"unicode/utf8"
	"unsafe"
)

func msg(message *C.char, s string) {
	m := C.CString(s)
	defer C.free(unsafe.Pointer(m))

	C.strcpy(message, m)
}

func argToGostrings(count C.uint, args **C.char) []string {
	length := count
	argslice := (*[1 << 30]*C.char)(unsafe.Pointer(args))[:length:length]

	gostrings := make([]string, count)

	for i, s := range argslice {
		gostrings[i] = C.GoString(s)
	}
	return gostrings
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
	*length = C.ulong(utf8.RuneCountInString(C.GoString(result)))

	return result
}

//export httpPost_deinit
func httpPost_deinit(initid *C.UDF_INIT) {
}

func main()  {}
