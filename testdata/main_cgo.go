//go:build ignore

package main

// #include <stdio.h>
// void hello() {
//   printf("hello from C\n");
// }
import "C"

func main() {
	C.hello()
}
