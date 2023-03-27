package core

type Database interface {
	Insert(int, string)
	Find(int) (string, bool)
	Delete(int)
	PrintBuffer()
}
