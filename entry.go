package btree

type Entry struct {
	key     Key
	pagePtr *page
	next    *Entry
	prev    *Entry
}

type EntryIterator struct {
	next    func() *Entry
	current *Entry
}

func (e *Entry) addNext(u *Entry) {
	if e.next != nil {
		e.next.prev = u
	}
	e.next, u.next = u, e.next
}

func (e *Entry) makeLast() {
	e.next = nil
}

func (e *Entry) makeFirst() {
	e.prev = nil
}

func (e *Entry) makeSingle() {
	e.makeFirst()
	e.makeLast()
}

func (e *Entry) iterator() EntryIterator {
	it := EntryIterator{
		current: e,
	}

	it.next = func() (currPos *Entry) {
		currPos = it.current
		if it.current != nil {
			it.current = it.current.next
		}
		return currPos
	}
	return it
}
