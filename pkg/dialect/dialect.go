package dialect

var dialectsMap = map[string]Dialect{}

type Dialect interface {
	CreateOrder(action string, price float64, qty float64) error
	RevokeOrder(OrderIds string) error
	CloseOrder(OrderIds string) error
}

func RegisterDialect(name string, dialect Dialect) {
	dialectsMap[name] = dialect
}

func GetDialect(name string) (dialect Dialect, ok bool) {
	dialect, ok = dialectsMap[name]
	return
}
