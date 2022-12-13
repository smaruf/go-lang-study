# the service
package main

import (
    "database/sql"
    "fmt"
    "github.com/gin-gonic/gin"
    _ "github.com/go-sql-driver/mysql"
    "net/http"
)

type Fruit struct {
    Id  int `json:"id"`
    Name string `json:"name"`
}

var con *sql.DB

func init(){
  //opening a mysql connection pool with another container 
   db, err := sql.Open("mysql", "root:password@tcp(host.docker.internal:3306)/payments")
   if err != nil {
       panic("failed to open a mysql connection")
   }
   con = db
}

func main() {
    r := gin.Default()
    r.GET("/fruits", fruits)
    r.Run() //server up on 8080
}

// THE REQUEST HANDLER
func fruits(c *gin.Context) {
    fruits := getFruits()
    c.JSON(http.StatusOK, fruits)
}

func getFruits() []Fruit {
    rows, _ := con.Query("SELECT * FROM fruits")
    fruits := []Fruit{}
    for rows.Next() {
        var r Fruit
        rows.Scan(&r.Id, &r.Name)
        fruits = append(fruits, r)
    }
    return fruits
}
