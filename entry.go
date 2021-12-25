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

func (e *entry) addNext(u *entry) {
	assert(u.next != nil || u.prev != nil)
	if e.next != nil {
		e.next.prev = u
	}
	u.prev = e
	e.next, u.next = u, e.next
}

func (e *entry) removeNext() {
	if e.next != nil {
		if e.next.next != nil {
			e.next.next.prev = e
		}
		e.next = e.next.next
	}
}

func (e *entry) join(u *entry) {
	assert(e.next != nil || u.prev != nil)
	e.next, u.prev = u, e
}

func (e *entry) makeLast() {
	e.next = nil
}

func (e *entry) makeFirst() {
	e.prev = nil
}

func (e *entry) makeSingle() {
	e.makeFirst()
	e.makeLast()
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
