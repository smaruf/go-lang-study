package main
// go get -u github.com/go-chi/chi/v5
// go get -u gorm.io/gorm
// go get -u gorm.io/driver/postgres
import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "github.com/go-chi/chi/v5"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)
func main() {
    port := "8082"
    log.Println("Running on http://localhost:" + port)
    //INITIALIZE DATABASE
    db := InitDatabase()
    // USE DEPENDENCY INJECTION TO AVOID CREATING MULTIPLE DB INSTANCES
    h := New(db)
r := chi.NewRouter()
    r.Get("/", hello)
    r.Get("/products", h.getAllProducts)
    r.Post("/product", h.addProduct)
log.Fatal(http.ListenAndServe(":"+port, r))
}
type Product struct {
    ID          int    `json:"id" gorm:"primaryKey"`
    Name        string `json:"name"`
    Description string `json:"description"`
    IsOnSale    bool   `json:"isOnSale"`
}
func hello(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Hello World!"))
}
func (h httpHandler) getAllProducts(w http.ResponseWriter, r *http.Request) {
    var products []Product
    if result := h.DB.Find(&products); result.Error != nil {
        fmt.Println(result.Error)
    }
    w.Header().Add("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(products)
}
func (h httpHandler) addProduct(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        log.Fatalln(err)
    }
    var product Product
    json.Unmarshal(body, &product)
    // ADD PRODUCT TO DATABASE
    if result := h.DB.Create(&product); result.Error != nil {
        fmt.Println(result.Error)
    }
// Send a 201 created response
    w.Header().Add("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode("Product Added Successfully")
}
func InitDatabase() *gorm.DB {
    dbURL := "postgres://root@localhost:26257/testdb?sslmode=disable"
    db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
    if err != nil {
        log.Fatal(err)
    }
    //AUTOMATICALLY CREATE PRODUCT TABLE IF IT DOESN'T EXIST
    db.AutoMigrate(&Product{})
    // ADD SOME SAMPLE DATA
    db.Create(&Product{ID: 1, Name: "Sample Product", Description: "This is a sample product", IsOnSale: false})
    return db
}
type httpHandler struct {
    DB *gorm.DB
}
func New(db *gorm.DB) httpHandler {
    return httpHandler{db}
}
