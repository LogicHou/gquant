package dialect

import "fmt"

var dialectsMap = map[string]Dialect{}

type Dialect interface {
	CreateOrder(action string, price float64, qty float64) error
	CloseOrder(OrderIds string) error
}

func Register(name string, dialect Dialect) {
	dialectsMap[name] = dialect
}

func Get(name string) (dialect Dialect, err error) {
	err = nil
	dialect, ok := dialectsMap[name]
	if !ok {
		err = fmt.Errorf("Cannot get dialect: %v", err)
	}
	return
}
