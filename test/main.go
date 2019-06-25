package main

import "fmt"
type good struct {
	Name string
}

func main(){
	s := []*good(nil)
	fmt.Printf("%T len = %d, cap = %d\n",s, len(s),cap(s))



}
