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
	var stack Stack
	u := &Entry{key: k}
	stack.push(current)
	child := current.scan(u)
	isLeaf := child == current
	for !isLeaf {
		current = child
		child = current.scan(u)
		stack.push(current)
		isLeaf = child == current
	}
	child = stack.pop()

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
	u := &Entry{key: k}
	var prev *page
	current := bt.root
	for current != prev {
		if current.find(u) {
			return true
		}
		prev = current
		current = current.scan(u)
	}
	return false
}

func (bt *Btree) LeftMajorOrder() []Key {
	return bt.leftMajorOrder(bt.root)
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
