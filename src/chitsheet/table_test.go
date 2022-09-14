func testAppend(t *testing.T) {
	list := NewSinglyLinkedList()

	tableTests := []struct {
		list     SinglyLinkedList
		val      int
		expected []int
	}{
		{list: list, val: 1, expected: []int{1}},    // case: empty list
		{list: list, val: 2, expected: []int{1, 2}}, // case: non-empty list
	}

	for _, tt := range tableTests {
		tt.list.Append(tt.val)
		require.Equal(t, tt.expected, tt.list.Values())
	}
}

func testValues(t *testing.T) {
	tableTests := []struct {
		list     SinglyLinkedList
		expected []int
	}{
		{list: NewSinglyLinkedList(), expected: []int{}},         // case: empty list
		{list: NewSinglyLinkedList(1), expected: []int{1}},       // case: one-item list
		{list: NewSinglyLinkedList(1, 2), expected: []int{1, 2}}, // case: multiple item list
	}

	for _, tt := range tableTests {
		require.Equal(t, tt.expected, tt.list.Values())
	}
}
