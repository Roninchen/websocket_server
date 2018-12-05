package main

import "fmt"


func (r *Person) hello() string{
	return r.name
}

func main(){
	fmt.Println(&Person{"liu",10})
	var r  = &Person{"liu",10}
	fmt.Println(r)
	fmt.Println(*r)
	fmt.Println(&r)
	var r2 *Person =  &Person{"liu2",10}
	fmt.Println(r2)
}