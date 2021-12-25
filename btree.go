package btree

type Btree struct {
	root   *page
	degree int
}

func (bt *Btree) Insert(k Key) {
	if bt.root == nil {
		bt.root = &page{}
	}
	// loop till leaf is reached
	current := bt.root
	var stack pageStack
	u := &entry{key: k}
	for current != nil {
		stack.push(current)
		pos := current.scan(u)
		if pos == nil {
			current = nil
		} else {
			current = pos.pagePtr
		}
	}
	child := stack.pop()

	for {
		child.insert(u)        // insert the key
		if !bt.isSafe(child) { // split and iterate if this page crossed threshold
			middleEntry, _ := child.splitMiddle()
			u = middleEntry
			if child == bt.root {
				page := &page{}
				page.insert(u)
				bt.root, page.head.pagePtr = page, bt.root
			} else {
				child = stack.pop()
				continue
			}
		}
		break
	}
}

func (bt *Btree) Search(k Key) bool {
	// loop till leaf is reached
	var stack pageEntryStack
	ok := bt.search(k, &stack)
	return ok
}

func (bt *Btree) LeftMajorOrder() []Key {
	return bt.leftMajorOrder(bt.root)
}

func (bt *Btree) Delete(k Key) {
	var stack pageEntryStack
	bt.search(k, &stack)
	pet := stack.pop()
	if pet.e.pagePtr != nil {
		bt.getMax(pet.e.pagePtr, &stack)
		max := stack.pop()
		pet.e, max.e = max.e, pet.e
	}
	pet.p.remove(pet.e)
	leafPagePtr := pet.p
label:
	if !bt.isSafe(leafPagePtr) {
		pet = stack.pop()
		if bt.concatSiblingAcross(pet.p, pet.e) {
			leafPagePtr = pet.p
			goto label
		} else {
			bt.transferFromSibling(pet.e, pet.p)
		}
	}

}

func (bt *Btree) leftMajorOrder(curr *page) []Key {
	var keys []Key
	if curr != nil {
		left := bt.leftMajorOrder(curr.head.pagePtr)
		keys = left
		ptr := curr.head.next
		for ptr != nil {
			keys = append(keys, ptr.key)
			right := bt.leftMajorOrder(ptr.pagePtr)
			keys = append(keys, right...)
			ptr = ptr.next
		}
	}
	return keys
}

func (bt *Btree) isSafe(page *page) bool {
	if bt.root == page {
		return page.keys <= 2*bt.degree && page.keys >= 1
	}
	return page.keys <= 2*bt.degree && page.keys >= bt.degree
}

func (bt *Btree) getMax(p *page, stack *pageEntryStack) Key {
	for p != nil {
		pet := &pageEntryTuple{p, p.tail}
		stack.push(pet)
		p = p.tail.pagePtr
	}
	return stack.peek().e.key
}

func (bt *Btree) search(k Key, stack *pageEntryStack) bool {
	u := &entry{key: k}
	curr := bt.root
	var (
		e  *entry
		ok bool
	)
	for !ok && (e == nil || e.pagePtr != nil) {
		e, ok = curr.find(u)
		if !ok {
			e = curr.scan(u)
		}
		pet := &pageEntryTuple{curr, e}
		if !ok {
			curr = e.pagePtr
		}
		stack.push(pet)
	}
	return ok
}

func (bt *Btree) concatSiblingAcross(p *page, e *entry) bool {
	ls := leftSibling(e)
	rs := rightSibling(e)
	if ls != nil && ls.keys+e.pagePtr.keys < 2*bt.degree {
		concatSiblingLeft(p, e)
		return true
	}
	if rs != nil && rs.keys+e.pagePtr.keys < 2*bt.degree {
		concatSiblingRight(p, e)
		return true
	}
	return false
}

func (bt *Btree) transferFromSibling(e *entry, p *page) {
	ls := leftSibling(e)
	rs := rightSibling(e)
	if ls != nil && ls.keys > ls.keys+1 {
		transferFromLeftSibling(e, p)
	}
	if rs != nil && rs.keys > rs.keys+1 {
		transferFromRightSibling(e, p)
	}
}
