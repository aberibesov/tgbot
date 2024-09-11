package product

var allProducts = []Product{
	{Title: "one"},
	{Title: "two"},
	{Title: "tree"},
	{Title: "four"},
	{Title: "five"},
}

type Product struct {
	Title string `json:"title"`
}
