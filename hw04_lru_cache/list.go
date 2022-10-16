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
	Length    int
	FrontItem *ListItem
	BackItem  *ListItem
}

func (list *list) Len() int {
	return list.Length
}

func (list *list) Front() *ListItem {
	return list.FrontItem
}

func (list *list) Back() *ListItem {
	return list.BackItem
}

func (list *list) PushFront(v interface{}) *ListItem {
	item := &ListItem{
		Value: v,
		Next:  nil,
		Prev:  nil,
	}

	if list.FrontItem != nil {
		list.FrontItem.Prev = item
		item.Next = list.FrontItem
	}

	list.FrontItem = item

	if list.BackItem == nil {
		list.BackItem = item
	}

	list.Length++

	return list.FrontItem
}

func (list *list) PushBack(v interface{}) *ListItem {
	item := &ListItem{
		Value: v,
		Next:  nil,
		Prev:  nil,
	}

	if list.BackItem != nil {
		list.BackItem.Next = item
		item.Prev = list.BackItem
	}

	list.BackItem = item

	if list.FrontItem == nil {
		list.FrontItem = item
	}

	list.Length++

	return list.BackItem
}

func (list *list) Remove(i *ListItem) {
	if i == nil {
		return
	}

	if list.FrontItem == i {
		list.FrontItem = i.Next
	}

	if list.BackItem == i {
		list.BackItem = i.Prev
	}

	if i.Next != nil {
		i.Next.Prev = i.Prev
	}

	if i.Prev != nil {
		i.Prev.Next = i.Next
	}

	list.Length--
}

func (list *list) MoveToFront(i *ListItem) {
	if i == nil {
		return
	}

	list.Remove(i)

	if list.FrontItem != nil {
		list.FrontItem.Prev = i
		i.Next = list.FrontItem
	}

	list.FrontItem = i

	if list.BackItem == nil {
		list.BackItem = i
	}

	list.Length++
}

func NewList() List {
	return new(list)
}
