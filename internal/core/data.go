package core

type Data struct {
	key   string
	value string
}

func CreateData(key string, value string) Data {
	return Data{key: key, value: value}
}
