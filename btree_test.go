package btree

import (
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
		expectedOrder := []int{-3, -2, -1, 0, 1, 2, 3}
		if !matchWithIntArray(lmo, expectedOrder) {
			t.Errorf("order should have been %v got %v", expectedOrder, lmo)
		}
	})
}

func TestBtree_GetMax(t *testing.T) {
	t.Run("returns max of the array keys", func(t *testing.T) {
		tree := getADummyBTree(2, 3)
		if max := tree.getMax(tree.root, new(pageEntryStack)); max != Int(3) {
			t.Errorf("max should have been %d got %d", 3, max)
		}
	})
}

func TestBtree_GetMin(t *testing.T) {
	t.Run("returns min of the array keys", func(t *testing.T) {
		tree := getADummyBTree(2, 3)
		if min := tree.getMin(tree.root, new(pageEntryStack)); min != Int(-3) {
			t.Errorf("max should have been %d got %d", 3, min)
		}
	})
}

func TestBtree_Delete(t *testing.T) {
	t.Run("remove a safe key", func(t *testing.T) {
		tree := getADummyBTree(2, 3)
		tree.Delete(Int(1))
		expectedOrder := []int{-3, -2, -1, 0, 2, 3}
		lmo := tree.LeftMajorOrder()
		if !matchWithIntArray(lmo, expectedOrder) {
			t.Errorf("order should have been %v got %v", expectedOrder, lmo)
		}
	})

	t.Run("remove an unsafe key", func(t *testing.T) {
		tree := getADummyBTree(2, 5)
		tree.Delete(Int(1))
		tree.Delete(Int(2))
		tree.Delete(Int(3))
		lmo := tree.LeftMajorOrder()
		expectedOrder := []int{-5, -4, -3, -2, -1, 0, 4, 5}
		if !matchWithIntArray(lmo, expectedOrder) {
			t.Errorf("order should have been %v got %v", expectedOrder, lmo)
		}
	})
}

func getADummyBTree(deg, k int) *Btree {
	tree := &Btree{
		degree: deg,
	}
	for i := -k; i <= k; i++ {
		tree.Insert(Int(i))
	}
	return tree
}

func matchWithIntArray(got []Key, expected []int) bool {
	match := len(got) == len(expected)
	if !match {
		return false
	}
	for i, v := range got {
		if v.(Int) != Int(expected[i]) {
			match = false
			break
		}
	}
	return match
}
