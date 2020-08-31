package game

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

//     A      B      C      D      E      F      G      H
// 8 (0,0)  (1,0)  (2,0)  (3,0)  (4,0)  (5,0)  (6,0)  (7,0)
// 7 (0,1)  (1,1)  (2,1)  (3,1)  (4,1)  (5,1)  (6,1)  (7,1)
// 6 (0,2)  (1,2)  (2,2)  (3,2)  (4,2)  (5,2)  (6,2)  (7,2)
// 5 (0,3)  (1,3)  (2,3)  (3,3)  (4,3)  (5,3)  (6,3)  (7,3)
// 4 (0,4)  (1,4)  (2,4)  (3,4)  (4,4)  (5,4)  (6,4)  (7,4)
// 3 (0,5)  (1,5)  (2,5)  (3,5)  (4,5)  (5,5)  (6,5)  (7,5)
// 2 (0,6)  (1,6)  (2,6)  (3,6)  (4,6)  (5,6)  (6,6)  (7,6)
// 1 (0,7)  (1,7)  (2,7)  (3,7)  (4,7)  (5,7)  (6,7)  (7,7)
func Test_PositionToModel_Should_Return_Valid_Model_When_Given_A_Position(t *testing.T) {
	tables := []struct {
		position      Position
		expectedModel string
	}{
		{Position{0, 0}, "a8"},
		{Position{0, 1}, "a7"},
		{Position{0, 2}, "a6"},
		{Position{0, 3}, "a5"},
		{Position{0, 4}, "a4"},
		{Position{0, 5}, "a3"},
		{Position{0, 6}, "a2"},
		{Position{0, 7}, "a1"},

		{Position{1, 0}, "b8"},
		{Position{1, 1}, "b7"},
		{Position{1, 2}, "b6"},
		{Position{1, 3}, "b5"},
		{Position{1, 4}, "b4"},
		{Position{1, 5}, "b3"},
		{Position{1, 6}, "b2"},
		{Position{1, 7}, "b1"},

		{Position{2, 2}, "c6"},
		{Position{3, 3}, "d5"},
		{Position{4, 4}, "e4"},
		{Position{5, 5}, "f3"},
		{Position{6, 6}, "g2"},
		{Position{7, 7}, "h1"},
	}

	for _, tt := range tables {
		modelActual := tt.position.positionToModel()
		assert.Equal(t, tt.expectedModel, modelActual, "should be equal")
	}
}

func Test_ModelToPosition_Should_Return_Valid_Position_When_Given_A_Position(t *testing.T) {
	tables := []struct {
		model            string
		expectedPosition Position
	}{
		{"a1", Position{0, 7}},
		{"a2", Position{0, 6}},
		{"a3", Position{0, 5}},
		{"a4", Position{0, 4}},
		{"a5", Position{0, 3}},
		{"a6", Position{0, 2}},
		{"a7", Position{0, 1}},
		{"a8", Position{0, 0}},
		{"b7", Position{1, 1}},
		{"c6", Position{2, 2}},
		{"d5", Position{3, 3}},
		{"e4", Position{4, 4}},
		{"f3", Position{5, 5}},
		{"g2", Position{6, 6}},
		{"h1", Position{7, 7}},
	}

	for _, tt := range tables {
		positionActual := modelToPosition(tt.model)
		assert.True(t, reflect.DeepEqual(tt.expectedPosition, positionActual), "should be equal")
	}
}
