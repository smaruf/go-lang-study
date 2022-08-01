func sha256Worker(c chan string, wg *sync.WaitGroup) {
    h := sha256.New()
    h.Write([]byte("nodejs-golang"))
    sha256_hash := hex.EncodeToString(h.Sum(nil))

    c <- sha256_hash

    wg.Done()
}

func sha256Array(w http.ResponseWriter, req *http.Request) {
    n, _ := strconv.Atoi(req.URL.Query().Get("n"))

    c := make(chan string, n)
    var wg sync.WaitGroup

    for i := 0; i < n; i++ {
        wg.Add(1)
        go sha256Worker(c, &wg)
    }

    wg.Wait()

    fmt.Fprint(w, n)
}
