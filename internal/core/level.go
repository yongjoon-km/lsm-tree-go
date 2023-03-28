package core

type Level int8

const (
	C1 Level = iota
	C2
	C3
)

func getFilePrefixPerLevel(level Level) string {
	switch level {
	case C1:
		return "C1"
	case C2:
		return "C2"
	case C3:
		return "C3"
	}
	return "C9"
}
