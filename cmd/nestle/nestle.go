package main

import (
	"fmt"
	"io/ioutil"

	"github.com/tarunKoyalwar/nestle"
)

func main() {
	// prematch := `(query|mutation)\s+[a-zA-Z]+[0-9]*[a-zA-Z]+(\([^(\(|\))]+\))*\s*`

	// postmatch := `(nested)`

	// nested_regex := prematch + `[{:nested:}]` + postmatch

	// nested_regex := `(attributes)\s[{:nested:}]`

	nested_regex := `[{:nested:}]nested`

	n, err := nestle.MustCompile(nested_regex)
	if err != nil {
		panic(err)
	}

	bin, _ := ioutil.ReadFile("../../testcase_graphql.txt")

	data := string(bin)

	// res := n.FindAllStringIndex(data)

	// for _, v := range res {
	// 	fmt.Printf("\n%v\n", data[v[0]:v[1]])
	// }

	res := n.FindAllString(data)

	for _, v := range res {
		fmt.Println(v)
	}

}
