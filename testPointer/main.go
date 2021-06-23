package main

import "fmt"

type Angle struct {
	Name string
	Age  int
}

func main() {
	m := make(map[string]*Angle)
	angs := []Angle{
		{Name: "ka", Age: 18},
		{Name: "kan", Age: 18},
		{Name: "kang", Age: 18},
	}
	for i, ang := range angs {
		fmt.Printf("%d %p\n", i, &ang)
		m[ang.Name] = &ang
	}
	for k, v := range m {
		fmt.Println("key=>"+k, "name=>"+v.Name)
	}
}

// https://jingwei.link/2018/09/06/golang-pointer-duplication.html
