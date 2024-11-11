package main

import (
    "gorm.io/gorm"
    "gorm.io/driver/sqlite"
    "fmt"
)

// Define a Product model
type Product struct {
    gorm.Model
    Code  string
    Price uint
}

func main() {
    // Connect to the database
    db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
    if err != nil {
        panic("failed to connect database")
    }

    // Migrate the schema
    db.AutoMigrate(&Product{})

    // Create a product
    db.Create(&Product{Code: "D42", Price: 100})

    // Read some data
    var product Product
    db.First(&product, 1)  // find product with integer primary key
    fmt.Println(product.Code, product.Price)

    // Update - update product's price to 200
    db.Model(&product).Update("Price", 200)

    // Delete - delete product
    db.Delete(&product)
}
