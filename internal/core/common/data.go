package common

type Data struct {
	Key   string
	Value string
}

func CreateData(key string, value string) Data {
	return Data{Key: key, Value: value}
}
