package convertor

type Convertible interface {
	Convert(obj interface{}) ([]byte, error)
}
