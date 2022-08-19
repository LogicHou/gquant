package mystrategy

type Strategy struct {
	Name      string
	Dialect   string
	AccessKey string
	SecretKey string
	Symbol    string
}

func init() {
	// cc := &Strategy{
	// 	Name:    "demo",
	// 	Dialect: cst.BINANCE,
	// 	// AccessKey: config.,
	// 	// SecretKey:,
	// 	// Symbol:,
	// }
	// fmt.Println(cc)
}

func OnKlineUpdate() {

}

func OnTickerUpdate() {

}
