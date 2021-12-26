package btree

type entry struct {
	key     Key
	pagePtr *page
	next    *entry
	prev    *entry
}

type entryIterator struct {
	next    func() *entry
	prev    func() *entry
	current *entry
}

func (e *entry) friendRight(u *entry) {
	assert(u.next != nil || u.prev != nil)
	if e.next != nil {
		e.next.prev = u
	}
	u.prev = e
	e.next, u.next = u, e.next
}

func (e *entry) unFriendRight() {
	if e.next != nil {
		if e.next.next != nil {
			e.next.next.prev = e
		}
		e.next = e.next.next
	}
}

func (e *entry) holdHands(u *entry) {
	assert(e.next != nil || u.prev != nil)
	e.next, u.prev = u, e
}

func (e *entry) chopRight() {
	e.next = nil
}

func (e *entry) chopLeft() {
	e.prev = nil
}

func (e *entry) chopBoth() {
	e.chopLeft()
	e.chopRight()
}

func (e *entry) iterator() entryIterator {
	it := entryIterator{
		current: e,
	}

	it.next = func() (currPos *entry) {
		currPos = it.current
		if it.current != nil {
			it.current = it.current.next
		}
		return currPos
	}

	it.prev = func() (currPos *entry) {
		currPos = it.current
		if it.current != nil {
			it.current = it.current.prev
		}
		return currPos
	}
	return it
}
