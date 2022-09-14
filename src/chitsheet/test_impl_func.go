// Node represents a node in singly linked list.
type Node struct {
	Val  int
	Next *Node
}

// List implements a singly linked list.
type List struct {
	first *Node
}

func newList(vals ...int) *List {
	if len(vals) < 1 {
		return &List{}
	}

	pseudo := &Node{}
	current := pseudo
	for _, v := range vals {
		current.Next = &Node{Val: v}
		current = current.Next
	}

	return &List{first: pseudo.Next}
}

// Append adds an item to the list.
func (l *List) Append(val int) {
	newNode := &Node{Val: val}

	if l.first == nil {
		l.first = newNode
		return
	}

	current := l.first
	for current.Next != nil {
		current = current.Next
	}

	current.Next = newNode
}

// Values returns the sequence of items in the list.
func (l *List) Values() []int {
	values := []int{}

	current := l.first
	for current != nil {
		values = append(values, current.Val)
		current = current.Next
	}

	return values
}
