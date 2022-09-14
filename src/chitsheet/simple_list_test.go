func TestSinglyLinkedList(t *testing.T) {
	for title, fn := range map[string]func(t *testing.T){
		"append item to list":  testAppend, // test: append functionality
		"return items in list": testValues, // test: values functionality
	} {
		t.Run(title, func(t *testing.T) {
			fn(t)
		})
	}
}
