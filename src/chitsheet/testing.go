// main.go
func Sum(x, y int) int {
    return x + y
}
// main_test.go
import ( 
    "testing"
    "reflect"
)
func TestSum(t *testing.T) {
    x, y := 2, 4
    expected := 2 + 4
    if !reflect.DeepEqual(Sum(x, y), expected) {
        t.Fatalf("Function Sum not working as expected")
    }
}
