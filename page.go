package btree

import (
	"fmt"
	"strings"
)

type page struct {
	head *entry
	tail *entry
	keys int
}

// returns the entry's which is prev to the first entry which is just greater than or equal to the input entry

func (p *page) scan(k Key) (entry *entry) {
	if p.head != nil {
		itr := p.iterator()
		prevItrEntry, itrEntry := itr.next(), itr.next()
		for itrEntry != nil {
			if k.Less(itrEntry.key) || equalKeys(k, itrEntry.key) {
				entry = prevItrEntry
				break
			}
			prevItrEntry = itrEntry
			itrEntry = itr.next()
		}
		if entry == nil {
			entry = prevItrEntry
		}
	}
	return entry
}

// returns true if it finds the entry which is equal to or else return's false and the output returned by the scan

func (p *page) find(k Key) (*entry, bool) {
	fitPos := p.scan(k)
	var matchedEntry *entry
	{
		itr := fitPos.iterator()
		itr.next()
		matchedEntry = itr.next()
	}
	if matchedEntry == nil {
		return fitPos, false
	}
	if equalKeys(matchedEntry.key, k) {
		return matchedEntry, true
	}
	return fitPos, false
}

// inserts an entry at pos which is just greater than or equal
// takes linear time to insert an entry

func (p *page) insert(e *entry) {
	assert(e == nil)
	e.chopBoth()
	if p.head == nil {
		p.head = &entry{}
		p.head.friendRight(e)
		p.tail = e
	} else {
		itr := p.iterator()
		prevItrEntry, itrEntry := itr.next(), itr.next()
		for itrEntry != nil {
			if e.key.Less(itrEntry.key) {
				break
			}
			prevItrEntry = itrEntry
			itrEntry = itr.next()
		}
		if itrEntry == nil {
			p.tail = e
		}
		prevItrEntry.friendRight(e)
	}
	p.keys++
}

func (p *page) remove(e *entry) bool {
	assert(e.key == nil)
	itr := p.iterator()
	prevItrEntry, itrEntry := itr.next(), itr.next()
	for itrEntry != nil {
		if itrEntry == e {
			prevItrEntry.unFriendRight()
			if p.tail == itrEntry {
				p.tail = prevItrEntry
			}
			p.keys--
			return true
		}
		prevItrEntry = itrEntry
		itrEntry = itr.next()
	}
	return false
}

func (p *page) splitMiddle() (*entry, *page) {
	itr := p.iterator()
	mid, index := p.keys/2, 0
	prevItrEntry, itrEntry := itr.next(), itr.next()
	for itrEntry != nil && index != mid {
		prevItrEntry = itrEntry
		itrEntry = itr.next()
		index++
	}

	// copy mid+1 to n entries to right Page
	midPlusOneEntry := itr.next()
	midPlusOneEntry.chopLeft()
	rightPageHead := &entry{}
	rightPageHead.holdHands(midPlusOneEntry)
	rightPage := &page{
		head: rightPageHead,
		tail: p.tail,
	}

	prevItrEntry.chopRight()
	itrEntry.chopBoth()

	// patch tails
	p.tail = prevItrEntry

	// set the no of keys is current page and in the new page
	rightPage.keys = p.keys - index - 1
	p.keys = index

	// reattach the middle entry page ptr to pageRight head entry
	rightPage.head.pagePtr = itrEntry.pagePtr
	itrEntry.pagePtr = rightPage
	return itrEntry, rightPage
}

func rightSibling(e *entry) *page {
	entryIterator := e.iterator()
	entryIterator.next()
	nextEntry := entryIterator.next()
	if nextEntry == nil {
		return nil
	}
	return nextEntry.pagePtr
}

func leftSibling(e *entry) *page {
	entryIterator := e.iterator()
	entryIterator.prev()
	prevEntry := entryIterator.prev()
	if prevEntry == nil {
		return nil
	}
	return prevEntry.pagePtr
}

func concatSiblingLeft(p *page, e *entry) bool {
	leftPageSibling := leftSibling(e)
	if leftPageSibling == nil {
		return false
	}
	rightPageSibling := e.pagePtr
	lspt := leftPageSibling.tail
	var rightEntry *entry
	{
		itr := rightPageSibling.iterator()
		itr.next()
		rightEntry = itr.next()
	}
	e.pagePtr = rightPageSibling.head.pagePtr
	assert(!p.remove(e))
	e.chopBoth()
	combine(lspt, e, rightEntry)

	leftPageSibling.keys += 1 + rightPageSibling.keys
	leftPageSibling.tail = rightPageSibling.tail
	return true
}

func concatSiblingRight(p *page, e *entry) bool {
	leftPageSibling := e.pagePtr
	rightPageSibling := rightSibling(e)
	if rightPageSibling == nil {
		return false
	}
	var (
		reAttachablePagePtr *page
		rightEntry          *entry
	)
	{
		itr := rightPageSibling.iterator()
		reAttachablePagePtr = itr.next().pagePtr
		rightEntry = itr.next()
	}
	leftEntry := e.pagePtr.tail

	e.pagePtr = reAttachablePagePtr
	assert(!p.remove(e))
	e.chopBoth()
	combine(leftEntry, e, rightEntry)

	leftPageSibling.keys += 1 + rightPageSibling.keys
	leftPageSibling.tail = rightPageSibling.tail
	return true
}

func combine(left, middle, right *entry) {
	{
		middle.chopBoth()
		left.chopRight()
		right.chopLeft()
	}

	left.holdHands(middle)
	middle.holdHands(right)
}

func transferFromLeftSibling(e *entry, p *page) bool {
	leftPageSibling := leftSibling(e)
	rightPageSibling := e.pagePtr
	if leftPageSibling == nil {
		return false
	}
	// remove e from page
	assert(!p.remove(e))

	// remove lspt from left sibling lspt left sibling page's tail
	var lspt *entry
	{
		lspt = leftPageSibling.tail
		leftPageSibling.remove(lspt)
	}

	// add lspt to p
	lspt.pagePtr = e.pagePtr
	p.insert(lspt)

	// add e to right page sibling
	e.pagePtr = nil
	rightPageSibling.insert(e)

	return true
}

func transferFromRightSibling(e *entry, p *page) bool {
	rightPageSibling := rightSibling(e)
	lefPageSibling := e.pagePtr
	if rightPageSibling == nil {
		return false
	}
	var rsphk, rightToE *entry // rsphk right sibling page's head key
	{
		// get right to e
		itr := e.iterator()
		itr.next()
		rightToE = itr.next()
	}

	// remove rightToE from p
	assert(!p.remove(rightToE))

	{
		// remove rsphk form rightPageSibling
		itr := rightPageSibling.iterator()
		itr.next()
		rsphk = itr.next()
		rightPageSibling.remove(rsphk)
	}

	// add rsphk to p
	rsphk.pagePtr = rightToE.pagePtr
	p.insert(rsphk)

	// add rightToE to lefPageSibling
	rightToE.pagePtr = nil
	lefPageSibling.insert(rightToE)
	return true
}
func (p *page) iterator() entryIterator {
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
