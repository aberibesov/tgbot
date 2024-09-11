package product

type Client struct {
	Id         int     `xml:"id"`
	Surname    string  `xml:"surname"`
	Name       string  `xml:"name"`
	Patronymic string  `xml:"patronymic"`
	Address    string  `xml:"address"`
	Telephone  string  `xml:"telephone"`
	Email      string  `xml:"email"`
	Login      string  `xml:"login"`
	Password   string  `xml:"password"`
	DateReg    string  `xml:"datereg"`
	Money      float32 `xml:"money"`
	Tariff     string  `xml:"tarif"`
	OldTariff  string  `xml:"oldtarif"`
	StartDate  string  `xml:"startdate"`
	EndDate    string  `xml:"enddate"`
	Extension  int     `xml:"extension"`
	Discount   int     `xml:"discount"`
}

type Response struct {
	Method string `xml:"method"`
	Code   string `xml:"code"`
	Client Client `xml:"user"`
}

type Request struct {
	Method string `xml:"method"`
	Id     int    `xml:"id"`
}
