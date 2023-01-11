package main

import (
   "fmt"
)

func minRemoveToMakeValid(s string) string {

   stack := make([]int, 0)
   toRemove := make([]int, 0)

   for i := 0; i < len(s); i++ {

      if string(s[i]) == "(" {
         stack = append(stack, i)
      }else if string(s[i]) == ")" && len(stack) == 0 {
         toRemove = append(toRemove, i)
      } else if string(s[i]) == ")" && len(stack) > 0 {
         stack = stack[:len(stack) - 1]
      }
   }

   newString := ""
   // --- Initialize the string builder
   var b strings.Builder
   b.Grow(len(s) - len(toRemove))
  // --- Write bytes in the bytes buffer
  for i := 0; i < len(s); i++ {
     if !removeMap[i] {
        b.Write([]byte{s[i]})
     }
  }

   return newString
}

func main() {
   s:= minRemoveToMakeValid("lee(t(c)o)de)")
   fmt.Println(s)
}
