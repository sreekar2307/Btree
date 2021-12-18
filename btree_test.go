package btree

import (
	"reflect"
	"testing"
)

type Int int

func (ik Int) Less(than Key) bool {
	return ik < than.(Int)
}

func TestBtree_Search(t *testing.T) {
	t.Run("when key searching for is absent", func(t *testing.T) {
		tree := getADummyBTree(2, 6)
		t.Run("returns false", func(t *testing.T) {
			if tree.Search(Int(8)) {
				t.Errorf("%d is not present should have returned false", 8)
			}
		})
	})

	t.Run("when key searching for is present", func(t *testing.T) {
		tree := getADummyBTree(2, 6)
		t.Run("returns false", func(t *testing.T) {
			if !tree.Search(Int(0)) {
				t.Errorf("%d is present should have returned true", 0)
			}
		})
	})
}

func TestBtree_LeftMajorOrder(t *testing.T) {
	t.Run("returns array of keys in increasing order", func(t *testing.T) {
		tree := getADummyBTree(2, 3)
		lmo := tree.LeftMajorOrder()
		expectedOrder := []Int{-3, -2, -1, 0, 1, 2, 3}
		if reflect.DeepEqual(lmo, expectedOrder) {
			t.Errorf("order should have been %v got %v", expectedOrder, lmo)
		}
	})
}

func getADummyBTree(deg, k int) *Btree {
	tree := &Btree{
		degree: deg,
	}
	for i:=-k;i<=k;i++ {
		tree.Insert(Int(i))
	}
	return tree
}