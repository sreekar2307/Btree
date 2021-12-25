package btree

import "testing"

func TestPage_insert(t *testing.T) {
	t.Run("maintains increasing order within page", func(t *testing.T) {
		page := getADummyIntPage([]int{2, 19, 13, 20})
		firstEntry := page.head.next
		secondEntry := firstEntry.next
		thirdEntry := secondEntry.next
		fourthEntry := thirdEntry.next
		if firstEntry.key != Int(2) {
			t.Errorf("the key supposed to be 2 got %d", firstEntry.key)
		}
		if secondEntry.key != Int(13) {
			t.Errorf("the key supposed to be 13 got %d", secondEntry.key)
		}
		if thirdEntry.key != Int(19) {
			t.Errorf("the key supposed to be 19 got %d", thirdEntry.key)
		}
		if fourthEntry.key != Int(20) {
			t.Errorf("the key supposed to be 20 got %d", fourthEntry.key)
		}
	})

	t.Run("increases keys count when new keys are inserted", func(t *testing.T) {
		page := getADummyIntPage([]int{2, 13})
		if page.keys != 2 {
			t.Errorf("expected keys to be %d got %d", 2, page.keys)
		}
	})
}

func TestPage_splitMiddle(t *testing.T) {
	t.Run("decrease the keys count to half", func(t *testing.T) {
		page := getADummyIntPage([]int{2, 19, 13})
		page.splitMiddle()
		if page.keys != 1 {
			t.Errorf("expected keys to be %d got %d", 1, page.keys)
		}
	})

	t.Run("newly created page size should be half of current page", func(t *testing.T) {
		page := getADummyIntPage([]int{2, 19, 13})
		_, rightPage := page.splitMiddle()
		if rightPage.keys != 1 {
			t.Errorf("expected keys to be %d got %d", 1, rightPage.keys)
		}
	})

	t.Run("returns the middle entry", func(t *testing.T) {
		page := getADummyIntPage([]int{2, 19, 13})
		middleEntry, _ := page.splitMiddle()
		if middleEntry.key != Int(13) {
			t.Errorf("expected middle entry to be %d got %d", middleEntry.key, 13)
		}
	})

	t.Run("next of middle but one entry should be nil", func(t *testing.T) {
		page := getADummyIntPage([]int{2, 19, 13})
		page.splitMiddle()
		firstEntry := page.head.next
		if firstEntry.next != nil {
			t.Errorf("expected middle but one entry's next to be nil %p", firstEntry.next)
		}
	})
}

func Test_concatSibling(t *testing.T) {
	/*
	       H      <->     6      <->      10
	       |              |               |
	 H <-> 3 <-> 4  H <-> 8        H <-> 11 <-> 12

	*/
	t.Run("concat sibling left", func(t *testing.T) {
		page := getADummyIntPage([]int{8})
		pageParent := getADummyIntPage([]int{6, 10})
		pageLeftSibling := getADummyIntPage([]int{3, 4})
		leftSiblingTail := pageLeftSibling.tail
		pageRightSibling := getADummyIntPage([]int{11, 12})
		rightSiblingHead := page.head.next
		pageParent.head.pagePtr = pageLeftSibling
		pageParent.head.next.pagePtr = page
		pageParent.head.next.next.pagePtr = pageRightSibling

		concatSiblingLeft(pageParent, pageParent.head.next)

		t.Run("left's sibling's next should be entry", func(t *testing.T) {
			if leftSiblingTail.next.key != Int(6) {
				t.Errorf("expected %d got %d", 6, leftSiblingTail.next.key)
			}
		})

		t.Run("entry's next should be right's sibling's head key", func(t *testing.T) {
			if rightSiblingHead.prev.key != Int(6) {
				t.Errorf("expected %d got %d", 6, rightSiblingHead.prev.key)
			}
		})

		t.Run("size of parent page should reduce by 1", func(t *testing.T) {
			if pageParent.keys != 1 {
				t.Errorf("expected %d got %d", 1, pageParent.keys)
			}
		})

		t.Run("size of left sibling page should inc by 1 + keys in right sibling", func(t *testing.T) {
			if pageLeftSibling.keys != 4 {
				t.Errorf("expected %d got %d", 4, pageLeftSibling.keys)
			}
		})

	})

	t.Run("concat sibling right", func(t *testing.T) {
		page := getADummyIntPage([]int{8})
		leftSiblingTail := page.tail
		pageParent := getADummyIntPage([]int{6, 10})
		pageLeftSibling := getADummyIntPage([]int{3, 4})
		pageRightSibling := getADummyIntPage([]int{11, 12})
		rightSiblingHead := pageRightSibling.head.next
		pageParent.head.pagePtr = pageLeftSibling
		pageParent.head.next.pagePtr = page
		pageParent.head.next.next.pagePtr = pageRightSibling

		concatSiblingRight(pageParent, pageParent.head.next)

		t.Run("left's sibling's next should be entry", func(t *testing.T) {
			if leftSiblingTail.next.key != Int(6) {
				t.Errorf("expected %d got %d", 6, leftSiblingTail.next.key)
			}
		})

		t.Run("entry's next should be right's sibling's head key", func(t *testing.T) {
			if rightSiblingHead.prev.key != Int(6) {
				t.Errorf("expected %d got %d", 6, rightSiblingHead.prev.key)
			}
		})

		t.Run("size of parent page should reduce by 1", func(t *testing.T) {
			if pageParent.keys != 1 {
				t.Errorf("expected %d got %d", 1, pageParent.keys)
			}
		})

		t.Run("size of left sibling page should inc by 1 + keys in right sibling", func(t *testing.T) {
			if page.keys != 4 {
				t.Errorf("expected %d got %d", 4, pageLeftSibling.keys)
			}
		})

	})
}

func Test_transferSibling(t *testing.T) {
	/*
	       H       <->     6      <->      10
	       |               |               |
	 H <-> 3 <-> 4 <-> 5   H <-> 8       H <-> 11 <-> 12

	*/
	t.Run("transfer from left sibling", func(t *testing.T) {
		page := getADummyIntPage([]int{8})
		pageParent := getADummyIntPage([]int{6, 10})
		pageLeftSibling := getADummyIntPage([]int{3, 4, 5})
		pageRightSibling := getADummyIntPage([]int{11, 12})
		pageParent.head.pagePtr = pageLeftSibling
		pageParent.head.next.pagePtr = page
		pageParent.head.next.next.pagePtr = pageRightSibling

		transferFromLeftSibling(pageParent.head.next, pageParent)

		t.Run("pageLeftSibling should have one less key", func(t *testing.T) {
			if pageLeftSibling.keys != 2 {
				t.Errorf("expected left page sibling to have %d keys got %d", 2, pageLeftSibling.keys)
			}
		})

		t.Run("pageParent should same keys", func(t *testing.T) {
			if pageParent.keys != 2 {
				t.Errorf("expected pageParent to have %d keys got %d", 2, pageParent.keys)
			}
		})

		t.Run("page should have one extra keys", func(t *testing.T) {
			if page.keys != 2 {
				t.Errorf("expected page to have %d keys got %d", 2, page.keys)
			}
		})

		t.Run("pageLeftSibling tail should be 4", func(t *testing.T) {
			if pageLeftSibling.tail.key != Int(4) {
				t.Errorf("expected pageLeftSibling tail to have %d keys got %d", 4,
					pageLeftSibling.tail.key)
			}
		})

		t.Run("pageParent head next should be 5", func(t *testing.T) {
			if pageParent.head.next.key != Int(5) {
				t.Errorf("expected pageParent.head.next to be %d keys got %d", 5,
					pageParent.head.next.key)
			}
		})

		t.Run("pageParent head next's page ptr should be pageRightSibling", func(t *testing.T) {
			if pageParent.head.next.pagePtr != page {
				t.Errorf("expected pageParent.head.next.pagePtr to be %p got %p", page,
					pageParent.head.next.pagePtr)
			}
		})

		t.Run("pageParent head next's next should be 10", func(t *testing.T) {
			if pageParent.head.next.next.key != Int(10) {
				t.Errorf("expected pageParent.head.next.next.key to be %d got %d", 10,
					pageParent.head.next.next.key)
			}
		})

		t.Run("page head next key should be 6", func(t *testing.T) {
			if page.head.next.key != Int(6) {
				t.Errorf("expected page.head.next.key to be %d got %d", 6, page.head.next.key)
			}
		})

		t.Run("page head next pagePtr to be nil", func(t *testing.T) {
			if page.head.next.pagePtr != nil {
				t.Errorf("expected page.head.next.pagePtr to be nil got %p", page.head.next.pagePtr)
			}
		})

		t.Run("page head next next key should be 8", func(t *testing.T) {
			if page.head.next.next.key != Int(8) {
				t.Errorf("expected page.head.next.next.key to be %d got %d", 8, page.head.next.next.key)
			}
		})
	})

	/*
	       H       <->     6      <->      10
	       |               |               |
	 H <-> 3 <->4      H <-> 8       H <-> 11 <-> 12 <-> 13

	*/

	t.Run("transfer from right sibling", func(t *testing.T) {
		page := getADummyIntPage([]int{8})
		pageParent := getADummyIntPage([]int{6, 10})
		pageLeftSibling := getADummyIntPage([]int{3, 4})
		pageRightSibling := getADummyIntPage([]int{11, 12, 13})
		pageParent.head.pagePtr = pageLeftSibling
		pageParent.head.next.pagePtr = page
		pageParent.head.next.next.pagePtr = pageRightSibling

		transferFromRightSibling(pageParent.head.next, pageParent)

		t.Run("pageRightSibling should have one less key", func(t *testing.T) {
			if pageRightSibling.keys != 2 {
				t.Errorf("expected right page sibling to have %d keys got %d", 2, pageRightSibling.keys)
			}
		})

		t.Run("pageParent should same keys", func(t *testing.T) {
			if pageParent.keys != 2 {
				t.Errorf("expected pageParent to have %d keys got %d", 2, pageParent.keys)
			}
		})

		t.Run("page should have one extra keys", func(t *testing.T) {
			if page.keys != 2 {
				t.Errorf("expected page to have %d keys got %d", 2, page.keys)
			}
		})

		t.Run("page's tail should be 10", func(t *testing.T) {
			if page.tail.key != Int(10) {
				t.Errorf("expected page's tail to be %d keys got %d", 10, page.tail.key)
			}
		})

		t.Run("page's tail's next should be nil", func(t *testing.T) {
			if page.tail.next != nil {
				t.Errorf("expected page's tail's next to be nil got %p", page.tail.next)
			}
		})

		t.Run("page's tail's pagePtr should be nil", func(t *testing.T) {
			if page.tail.pagePtr != nil {
				t.Errorf("expected page's pagePtr to be nil got %p", page.tail.pagePtr)
			}
		})

		t.Run("pageParent.head.next.next's key should be 11", func(t *testing.T) {
			if pageParent.head.next.next.key != Int(11) {
				t.Errorf("expected pageParent.head.next.next.key to be %d got %d", 11,
					pageParent.head.next.next.key)
			}
		})

		t.Run("pageParent.head.next.next's PagePtr should be rightSibling", func(t *testing.T) {
			if pageParent.head.next.next.pagePtr != pageRightSibling {
				t.Errorf("expected pageParent.head.next.next.pagePtr to be %p got %p", pageRightSibling,
					pageParent.head.next.next.pagePtr)
			}
		})

		t.Run("pageParent.head.next.next's next should be nil", func(t *testing.T) {
			if pageParent.head.next.next.next != nil {
				t.Errorf("expected pageParent.head.next.next.next to be nil got %p",
					pageParent.head.next.next.next)
			}
		})

		t.Run("pageRightSibling's head key should be 12", func(t *testing.T) {
			if pageRightSibling.head.next.key != Int(12) {
				t.Errorf("expected pageRightSibling.head.next.key to be %d got %d", 12,
					pageRightSibling.head.next.key)
			}
		})
	})

}

func TestPage_scan(t *testing.T) {
	t.Run("when page is not a leaf", func(t *testing.T) {
		t.Run("returns the first entry's page pointer which is just >=", func(t *testing.T) {
			page := getADummyIntPage([]int{2, 19, 13, 20})
			firstEntry := page.head.next
			firstEntry.pagePtr = getADummyIntPage([]int{10, 11})
			e := page.scan(Int(12))
			if e.pagePtr != firstEntry.pagePtr {
				t.Errorf("expected page ptr to be %p got %p", firstEntry.pagePtr, e.pagePtr)
			}
		})
	})

	t.Run("when a max key is scanned for", func(t *testing.T) {
		t.Run("when max entry page ptr is not nil", func(t *testing.T) {
			page := getADummyIntPage([]int{19})
			firstEntry := page.head.next
			firstEntry.pagePtr = getADummyIntPage([]int{21, 23})
			e := page.scan(Int(20))
			if e.pagePtr != firstEntry.pagePtr {
				t.Errorf("expected page ptr to be %p got %p", firstEntry.pagePtr, e.pagePtr)
			}
		})
	})

	t.Run("when a min key is scanned for", func(t *testing.T) {
		t.Run("when min entry page ptr is not nil", func(t *testing.T) {
			page := getADummyIntPage([]int{19})
			page.head.pagePtr = getADummyIntPage([]int{13, 14})
			e := page.scan(Int(17))
			if e.pagePtr != page.head.pagePtr {
				t.Errorf("expected page ptr to be %p got %p", page.head.pagePtr, e.pagePtr)
			}
		})
	})
}

func getADummyIntPage(arr []int) *page {
	page := &page{}
	for _, v := range arr {
		page.insert(&entry{
			key:     Int(v),
			pagePtr: nil,
		})
	}
	return page
}
