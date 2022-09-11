// Declaring a map
var cities map[string]string
// Initializing
cities = make(map[string]string)
cities // nil
// Insert
cities["NY"] = "EUA"
// Retrieve
newYork = cities["NY"]
newYork // "EUA"
// Delete
delete(cities, "NY")
// Check if a key is setted
value, ok := cities["NY"]
ok // false
value // ""
