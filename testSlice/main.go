package main

import "log"

func printSlice1(s []string) {
	log.Println(s)
}

func printSlice2(s *[]string) {
	log.Println(s)
}

func main() {
	var names []string
	names = append(names, "Tom")
	names = append(names, "Julia")
	log.Println("using original slice")
	printSlice1(names)
	log.Println("using pointer of slice")
	printSlice2(&names)
}
