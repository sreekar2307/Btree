package btree

import (
	"fmt"
	"strings"
)

type Key interface {
	Less(than Key) bool
}

type page struct {
	head *Entry
	keys int
}

// returns the first Entry's page pointer which is prev to an entry which is just greater than or equal
// If matched Entry's pagePtr is nil then return current

func (p *page) scan(e *Entry) (foundPtr *page) {
	if p.head != nil {
		itr := p.iterator()
		prevItrEntry, itrEntry := itr.next(), itr.next()
		for itrEntry != nil {
			if e.key.Less(itrEntry.key) {
				foundPtr = prevItrEntry.pagePtr
				break
			}
			prevItrEntry = itrEntry
			itrEntry = itr.next()
		}
		if foundPtr == nil {
			foundPtr = prevItrEntry.pagePtr
		}
	}
	if foundPtr == nil {
		foundPtr = p
	}
	return foundPtr
}

// returns true if it finds the entry which is equal to or else return's false

func (p *page) find(e *Entry) bool {
	itr := p.iterator()
	itr.next()
	itrEntry := itr.next()
	for itrEntry != nil {
		if !e.key.Less(itrEntry.key) && !itrEntry.key.Less(e.key) {
			return true
		}
		itrEntry = itr.next()
	}
	return false
}

// inserts an Entry at pos which is just greater than or equal
// takes linear time to insert an entry

func (p *page) insert(entry *Entry) {
	assert(entry == nil)
	if p.head == nil {
		p.head = &Entry{}
		p.head.addNext(entry)
	} else {
		itr := p.iterator()
		prevItrEntry, itrEntry := itr.next(), itr.next()
		for itrEntry != nil {
			if entry.key.Less(itrEntry.key) {
				break
			}
			prevItrEntry = itrEntry
			itrEntry = itr.next()
		}
		prevItrEntry.addNext(entry)
	}
	p.keys++
}

func (p *page) splitMiddle() (*Entry, *page) {
	itr := p.iterator()
	mid, index := p.keys/2, 0
	prevItrEntry, itrEntry := itr.next(), itr.next()
	for itrEntry != nil && index != mid {
		prevItrEntry = itrEntry
		itrEntry = itr.next()
		index++
	}

	// copy mid+1 to n entries to right Page
	rightPage := &page{}
	rightPage.insert(itr.next())

	prevItrEntry.makeLast()
	itrEntry.makeSingle()

	// set the no of keys is current page and in the new page
	rightPage.keys = p.keys - index - 1
	p.keys = index

	// reattach the middle entry page ptr to pageRight head entry
	rightPage.head.pagePtr = itrEntry.pagePtr
	itrEntry.pagePtr = rightPage
	return itrEntry, rightPage
}

func (p *page) iterator() EntryIterator {
	assert(p.head == nil)
	return p.head.iterator()
}

func (p *page) String() string {
	itr := p.iterator()
	var str strings.Builder
	itr.next()
	itrEntry := itr.next()
	for itrEntry != nil {
		_, _ = fmt.Fprintf(&str, "%v->", itrEntry.key)
		itrEntry = itr.next()
	}
	str.WriteString("NIL")
	return str.String()
}

func assert(cond bool) {
	if cond {
		panic("Assertion failed")
	}
}
