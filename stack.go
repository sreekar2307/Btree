package btree

type Stack []*page

func (s *Stack) push(p *page) {
	*s = append(*s, p)
}
func (s *Stack) Len() int {
	return len(*s)
}

func (s *Stack) pop() *page {
	if *s == nil || len(*s) == 0 {
		return nil
	}
	curr := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return curr
}
