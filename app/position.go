package main

type Position struct {
	x int
	y int
}

// maps (0, 0) -> (A, 8) and (7, 7) -> (H, 1)
//     A      B      C      D      E      F      G      H
// 8 (0,0)  (1,0)  (2,0)  (3,0)  (4,0)  (5,0)  (6,0)  (7,0)
// 7 (0,1)  (1,1)  (2,1)  (3,1)  (4,1)  (5,1)  (6,1)  (7,1)
// 6 (0,2)  (1,2)  (2,2)  (3,2)  (4,2)  (5,2)  (6,2)  (7,2)
// 5 (0,3)  (1,3)  (2,3)  (3,3)  (4,3)  (5,3)  (6,3)  (7,3)
// 4 (0,4)  (1,4)  (2,4)  (3,4)  (4,4)  (5,4)  (6,4)  (7,4)
// 3 (0,5)  (1,5)  (2,5)  (3,5)  (4,5)  (5,5)  (6,5)  (7,5)
// 2 (0,6)  (1,6)  (2,6)  (3,6)  (4,6)  (5,6)  (6,6)  (7,6)
// 1 (0,7)  (1,7)  (2,7)  (3,7)  (4,7)  (5,7)  (6,7)  (7,7)
func (p *Position) positionToModel() string {
	var firstChar, secondChar string

	switch p.x {
	case 0:
		firstChar = "a"
	case 1:
		firstChar = "b"
	case 2:
		firstChar = "c"
	case 3:
		firstChar = "d"
	case 4:
		firstChar = "e"
	case 5:
		firstChar = "f"
	case 6:
		firstChar = "g"
	case 7:
		firstChar = "h"
	}

	switch p.y {
	case 0:
		secondChar = "8"
	case 1:
		secondChar = "7"
	case 2:
		secondChar = "6"
	case 3:
		secondChar = "5"
	case 4:
		secondChar = "4"
	case 5:
		secondChar = "3"
	case 6:
		secondChar = "2"
	case 7:
		secondChar = "1"
	}
	return firstChar + secondChar
}

func modelToPosition(model string) Position {
	var posX, posY int

	// for the first char
	switch model[0] {
	case 'a':
		posX = 0
	case 'b':
		posX = 1
	case 'c':
		posX = 2
	case 'd':
		posX = 3
	case 'e':
		posX = 4
	case 'f':
		posX = 5
	case 'g':
		posX = 6
	case 'h':
		posX = 7
	}

	// for the second char
	switch model[1] {
	case '8':
		posY = 0
	case '7':
		posY = 1
	case '6':
		posY = 2
	case '5':
		posY = 3
	case '4':
		posY = 4
	case '3':
		posY = 5
	case '2':
		posY = 6
	case '1':
		posY = 7
	}
	return Position{posX, posY}
}
