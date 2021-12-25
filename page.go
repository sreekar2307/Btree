package btree

import (
	"fmt"
	"strings"
)

type Key interface {
	Less(than Key) bool
}

type page struct {
	head *entry
	tail *entry
	keys int
}

// returns the first entry's page pointer which is prev to an entry which is just greater than or equal
// If matched entry's pagePtr is nil then return current

func (p *page) scan(e *entry) (entry *entry) {
	if p.head != nil {
		itr := p.iterator()
		prevItrEntry, itrEntry := itr.next(), itr.next()
		for itrEntry != nil {
			if e.key.Less(itrEntry.key) {
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

// returns true if it finds the entry which is equal to or else return's false

func (p *page) find(e *entry) (*entry, bool) {
	itr := p.iterator()
	itr.next()
	itrEntry := itr.next()
	for itrEntry != nil {
		if !e.key.Less(itrEntry.key) && !itrEntry.key.Less(e.key) {
			return itrEntry, true
		}
		itrEntry = itr.next()
	}
	return nil, false
}

// inserts an entry at pos which is just greater than or equal
// takes linear time to insert an entry

func (p *page) insert(e *entry) {
	assert(e == nil)
	e.makeSingle()
	if p.head == nil {
		p.head = &entry{}
		p.head.addNext(e)
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
		prevItrEntry.addNext(e)
	}
	p.keys++
}

func (p *page) remove(e *entry) {
	assert(e.key == nil)
	_, ok := p.find(e)
	assert(!ok)
	var prevEntry *entry
	{
		itr := e.iterator()
		itr.prev()
		prevEntry = itr.prev()
	}
	prevEntry.removeNext()
	if p.tail == e {
		p.tail = prevEntry
	}
	p.keys--
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
	midPlusOneEntry.makeFirst()
	rightPageHead := &entry{}
	rightPageHead.join(midPlusOneEntry)
	rightPage := &page{
		head: rightPageHead,
		tail: p.tail,
	}

	prevItrEntry.makeLast()
	itrEntry.makeSingle()

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

func concatSiblingLeft(page *page, e *entry) bool {
	leftPageSibling := leftSibling(e)
	if leftPageSibling == nil {
		return false
	}
	rightPageSibling := e.pagePtr
	leftEntry := leftPageSibling.tail
	var rightEntry *entry
	{
		itr := rightPageSibling.iterator()
		itr.next()
		rightEntry = itr.next()
	}
	reAttachablePagePtr := rightPageSibling.head.pagePtr

	combine(leftEntry, e, rightEntry)

	leftPageSibling.keys += 1 + rightPageSibling.keys
	e.pagePtr = reAttachablePagePtr
	page.keys--
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

	combine(leftEntry, e, rightEntry)

	leftPageSibling.keys += 1 + rightPageSibling.keys
	e.pagePtr = reAttachablePagePtr
	p.keys--
	return true
}

func combine(left, middle, right *entry) {
	{
		middle.makeSingle()
		left.makeLast()
		right.makeFirst()
	}

	left.join(middle)
	middle.join(right)
}

func transferFromLeftSibling(e *entry, p *page) bool {
	leftPageSibling := leftSibling(e)
	rightPageSibling := e.pagePtr
	if leftPageSibling == nil {
		return false
	}
	// remove lspt from left sibling lspt left sibling page's tail
	var lspt *entry
	{
		lspt = leftPageSibling.tail
		leftPageSibling.remove(lspt)
	}

	// remove e from page and add lspt to p
	lspt.pagePtr = e.pagePtr
	p.remove(e)
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
		// remove rsphk form rightPageSibling
		itr := rightPageSibling.iterator()
		itr.next()
		rsphk = itr.next()
		rightPageSibling.remove(rsphk)
	}

	{
		// get right to e
		itr := e.iterator()
		itr.next()
		rightToE = itr.next()
	}

	// remove rightToE from p and add rsphk to p
	rsphk.pagePtr = rightToE.pagePtr
	p.remove(rightToE)
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
