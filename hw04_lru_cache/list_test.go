package hw04lrucache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())

		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, 2, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)
	})

	t.Run("moveToFront", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		l.PushBack(40)  // [10, 20, 30, 40]
		l.PushBack(50)  // [10, 20, 30, 40, 50]
		require.Equal(t, 5, l.Len())

		// Перемещаем элемент из середины списка
		l.MoveToFront(l.Back().Prev.Prev) // [30, 10, 20, 40, 50]
		require.Equal(t, 30, l.Front().Value)
		require.Equal(t, 40, l.Back().Prev.Value)
		require.Equal(t, 20, l.Back().Prev.Prev.Value)

		// Перемещаем элемент из начала списка
		l.MoveToFront(l.Front()) // [30, 10, 20, 40, 50]
		require.Equal(t, 30, l.Front().Value)

		// Перемещаем элемент из конца списка
		l.MoveToFront(l.Back()) // [50, 30, 10, 20, 40]
		require.Equal(t, 50, l.Front().Value)
		require.Equal(t, 40, l.Back().Value)
	})

	t.Run("removeElements", func(t *testing.T) {
		l := NewList()
		item := l.PushFront(10) // [10]
		require.Equal(t, 1, l.Len())

		l.Remove(item)
		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())

		item = l.PushBack(20) // [20]
		l.PushBack(30)        // [20, 30]
		require.Equal(t, 2, l.Len())

		l.Remove(l.Back()) // [20]
		require.Equal(t, 1, l.Len())
		require.Equal(t, item, l.Front())
		require.Equal(t, item, l.Back())

		item = l.PushBack(40) // [20, 40]
		l.PushBack(50)        // [20, 40, 50]
		require.Equal(t, 3, l.Len())

		l.Remove(l.Front())
		require.Equal(t, 2, l.Len())
		require.Equal(t, item, l.Front())
	})
}
