func fibonacci(w http.ResponseWriter, req *http.Request) {
    nValue, _ := strconv.Atoi(req.URL.Query().Get("n"))

    var n = uint(nValue)

    if n <= 1 {
        fmt.Fprint(w, big.NewInt(int64(n)))
    }

    var n2, n1 = big.NewInt(0), big.NewInt(1)

    for i := uint(1); i < n; i++ {
        n2.Add(n2, n1)
        n1, n2 = n2, n1
    }

    fmt.Fprint(w, n1)
}
