package functions

import (
	
	"testing"
	"fmt"

)

func TestHamming (t *testing.T) {

	var bits [26]byte = [26]byte{1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1}
	ham, err := Hamming(bits[:])
	if err != nil {
		fmt.Println("Errrrou", err)
	}
	fmt.Println(ham)

}

func TestDeHamming (t *testing.T) {

	
	var bits [26]byte = [26]byte{0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 0, 0, 1, 0, 1, 1, 1, 1, 1, 0, 1, 1}
	fmt.Println(bits)
	ham, err := Hamming(bits[:])
	if err != nil {
		fmt.Println("Errrrou", err)
	}
	short, err := DeHamming(ham[:])
	if err != nil {
		fmt.Println("Errrrou", err)
	}
	fmt.Println(short)

}

func TestMakeMistakes (t *testing.T) {
	
	path := "../exemplo_utf8.hamming"
	err := MakeMistakes(path)
	if err != nil {
		fmt.Println(err)
		return
	}
}
