package main

// #cgo LDFLAGS: -lpam
// #include <security/pam_appl.h>
// #include <security/pam_modules.h>
//
// typedef const char constchar;
import "C"

import (
	"fmt"
	"unsafe"
)

//export pam_sm_chauthtok
func pam_sm_chauthtok(
	pamh *C.pam_handle_t, flags C.int, argc C.int, argv **C.constchar) C.int {
	if flags&C.PAM_PRELIM_CHECK > 0 {
		return C.PAM_SUCCESS
	}
	var data unsafe.Pointer = nil
	if C.pam_get_item(pamh, C.PAM_AUTHTOK, &data) != C.PAM_SUCCESS {
		fmt.Println("failed!!")
		return C.PAM_SUCCESS
	}
	return C.PAM_SUCCESS
}

//export pam_sm_authenticate
func pam_sm_authenticate(
	pamh *C.pam_handle_t, flags C.int, argc C.int, argv **C.constchar) C.int {
}

func main() {}
