package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	firstItem *ListItem
	lastItem  *ListItem
	length    int
}

func (l *list) Len() int {
	return l.length
}

func (l *list) Front() *ListItem {
	return l.firstItem
}

func (l *list) Back() *ListItem {
	return l.lastItem
}

func (l *list) PushFront(v interface{}) *ListItem {
	defer func() { l.length++ }()

	newItem := &ListItem{
		Value: v,
		Next:  nil,
		Prev:  nil,
	}

	if l.length == 0 {
		l.firstItem = newItem
		l.lastItem = newItem
		return newItem
	}

	newItem.Next = l.firstItem
	l.firstItem.Prev = newItem
	l.firstItem = newItem

	return newItem
}

func (l *list) PushBack(v interface{}) *ListItem {
	defer func() { l.length++ }()
	newItem := &ListItem{
		Value: v,
		Next:  nil,
		Prev:  nil,
	}

	if l.length == 0 {
		l.firstItem = newItem
		l.lastItem = newItem
		return newItem
	}

	newItem.Prev = l.lastItem
	l.lastItem.Next = newItem
	l.lastItem = newItem

	return newItem
}

func (l *list) Remove(i *ListItem) {
	defer func() { l.length-- }()

	if l.length == 1 {
		l.firstItem = nil
		l.lastItem = nil
		return
	}

	if i.Prev == nil {
		i.Next.Prev = nil
		l.firstItem = i.Next
		return
	}

	if i.Next == nil {
		i.Prev.Next = nil
		l.lastItem = i.Prev
		return
	}

	i.Next.Prev = i.Prev
	i.Prev.Next = i.Next
}

func (l *list) MoveToFront(i *ListItem) {
	switch i {
	case l.firstItem:
		return
	case l.lastItem:
		i.Prev.Next = nil
		l.lastItem = i.Prev
	default:
		i.Prev.Next = i.Next
		i.Next.Prev = i.Prev
	}

	i.Next = l.firstItem
	l.firstItem.Prev = i
	l.firstItem = i
}

func NewList() List {
	return new(list)
}
