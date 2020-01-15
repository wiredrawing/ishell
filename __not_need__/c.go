package main

/*
#include <stdlib.h>
extern void func_to_export(const char*);
*/
import "C"
import "unsafe"

func main() {
  p := C.CString("lorem ipsum")
  defer C.free(unsafe.Pointer(p))
  C.func_to_export(p)
}