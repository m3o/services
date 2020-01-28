package main

import(
	"fmt"
	"github.com/kytra-app/helpers/textenhancer"
)

func main() {
	str := "Hey @Ben and @Sam, how are you?"
	srv := textenhancer.Service{}
	fmt.Println(srv.ListTaggedUsers(str))
}