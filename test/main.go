package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
)

func main() {

	s := `
a: 
  - 1
  - 2 
c: d`

	node := yaml.Node{}
	err := yaml.Unmarshal([]byte(s), &node)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(node)

}
