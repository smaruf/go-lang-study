var count int

func GetPdfUrl(ctx context.Context) (string, error) {
	count++

	if count <= 3 {
		return "", errors.New("boom")
	} else {
		return "https://linktopdf.com", nil
	}
}

func main() {
	r := Retry(GetPdfUrl, 5, 2*time.Second)

	res, err := r(context.Background())

	fmt.Println(res, err)
}
