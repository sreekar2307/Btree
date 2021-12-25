package btree

type Key interface {
	Less(than Key) bool
}

func equalKeys(a, b Key) bool {
	return !a.Less(b) && !b.Less(a)
}
