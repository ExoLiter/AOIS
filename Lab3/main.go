package main

import (
	"fmt"
	"lab3/circuits"
)

func printEquations(title string, eqs []circuits.Equation, showSDNF bool) {
	fmt.Println(title)
	for _, eq := range eqs {
		if showSDNF {
			fmt.Printf("%s:\nSDNF: %s\nMinimized: %s\n\n", eq.Name, eq.SDNF, eq.Minimized)
			continue
		}
		fmt.Printf("%s = %s\n", eq.Name, eq.Minimized)
	}
	fmt.Println()
}

func main() {
	printEquations("Part 1: 1-bit subtractor (ODV-3)", circuits.GetSubtractorEquations(), true)
	printEquations("Part 2.1: Decoder 5421 -> Binary", circuits.GetDecoder5421Equations(), false)
	printEquations("Part 2.2: Encoder Binary -> 5421", circuits.GetEncoder5421Equations(), false)
	printEquations("Part 3: 16-state down counter (T flip-flop)", circuits.GetCounterEquations(), false)
}
