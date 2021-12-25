package btree

type pageStack []*page

func (s *pageStack) push(p *page) {
	*s = append(*s, p)
}
func (s *pageStack) Len() int {
	return len(*s)
}

func (s *pageStack) pop() *page {
	if *s == nil || len(*s) == 0 {
		return nil
	}
	curr := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return curr
}

func (s *pageStack) peek() *page {
	if *s == nil || len(*s) == 0 {
		return nil
	}
	return (*s)[len(*s)-1]
}

type pageEntryTuple struct {
	p *page
	e *entry
}

type pageEntryStack []*pageEntryTuple

func (s *pageEntryStack) push(p *pageEntryTuple) {
	*s = append(*s, p)
}
func (s *pageEntryStack) Len() int {
	return len(*s)
}

func (s *pageEntryStack) pop() *pageEntryTuple {
	if *s == nil || len(*s) == 0 {
		return nil
	}
	curr := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return curr
}

func (s *pageEntryStack) peek() *pageEntryTuple {
	if *s == nil || len(*s) == 0 {
		return nil
	}
	return (*s)[len(*s)-1]
}
