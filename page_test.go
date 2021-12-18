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

	t.Run("returns the middle Entry", func(t *testing.T) {
		page := getADummyIntPage([]int{2, 19, 13})
		middleEntry, _ := page.splitMiddle()
		if middleEntry.key != Int(13) {
			t.Errorf("expected middle entry to be %d got %d", middleEntry.key, 13)
		}
	})

	t.Run("next of middle but one Entry should be nil", func(t *testing.T) {
		page := getADummyIntPage([]int{2, 19, 13})
		page.splitMiddle()
		firstEntry := page.head.next
		if firstEntry.next != nil {
			t.Errorf("expected middle but one Entry's next to be nil %p", firstEntry.next)
		}
	})
}

func TestPage_scan(t *testing.T) {
	t.Run("when page is not a leaf", func(t *testing.T) {
		t.Run("returns the first Entry's page pointer which is just >=", func(t *testing.T) {
			page := getADummyIntPage([]int{2, 19, 13, 20})
			firstEntry := page.head.next
			firstEntry.pagePtr = getADummyIntPage([]int{10, 11})
			pagePtr := page.scan(&Entry{
				key: Int(12),
			})
			if pagePtr != firstEntry.pagePtr {
				t.Errorf("expected page ptr to be %p got %p", firstEntry.pagePtr, pagePtr)
			}
		})
	})

	t.Run("when a max key is scanned for", func(t *testing.T) {
		t.Run("when max entry page ptr is not nil", func(t *testing.T) {
			page := getADummyIntPage([]int{19})
			firstEntry := page.head.next
			firstEntry.pagePtr = getADummyIntPage([]int{21, 23})
			pagePtr := page.scan(&Entry{
				key: Int(20),
			})
			if pagePtr != firstEntry.pagePtr {
				t.Errorf("expected page ptr to be %p got %p", firstEntry.pagePtr, pagePtr)
			}
		})
		t.Run("when max entry page ptr is nil", func(t *testing.T) {
			page := getADummyIntPage([]int{19})
			pagePtr := page.scan(&Entry{
				key: Int(20),
			})
			if pagePtr != page {
				t.Errorf("expected page ptr to be %p got %p", page, pagePtr)
			}
		})
	})

	t.Run("when a min key is scanned for", func(t *testing.T) {
		t.Run("when min entry page ptr is not nil", func(t *testing.T) {
			page := getADummyIntPage([]int{19})
			page.head.pagePtr = getADummyIntPage([]int{13, 14})
			pagePtr := page.scan(&Entry{
				key: Int(17),
			})
			if pagePtr != page.head.pagePtr {
				t.Errorf("expected page ptr to be %p got %p", page.head.pagePtr, pagePtr)
			}
		})
		t.Run("when in entry page ptr is nil", func(t *testing.T) {
			page := getADummyIntPage([]int{19})
			pagePtr := page.scan(&Entry{
				key: Int(17),
			})
			if pagePtr != page {
				t.Errorf("expected page ptr to be %p got %p", page, pagePtr)
			}

		})
	})

	t.Run("when page is a leaf", func(t *testing.T) {
		t.Run("should return the current page", func(t *testing.T) {
			page := getADummyIntPage([]int{2, 19, 13, 20})
			pagePtr := page.scan(&Entry{
				key: Int(12),
			})
			if pagePtr != page {
				t.Errorf("expected page ptr to be %p got %p", page, pagePtr)
			}
		})
	})
}

func getADummyIntPage(arr []int) *page {
	page := &page{}
	for _, v := range arr {
		page.insert(&Entry{
			key:     Int(v),
			pagePtr: nil,
		})
	}
	return page
}
