package dialect

type binance struct{}

var _ Dialect = (*binance)(nil)

func init() {
	RegisterDialect("binance", &binance{})
}

func (s *binance) CreateOrder(action string, price float64, qty float64) error {
	//TODO
	return nil
}

func (s *binance) RevokeOrder(OrderIds string) error {
	//TODO
	return nil
}

func (s *binance) CloseOrder(OrderIds string) error {
	//TODO
	return nil
}
