package models

type MenuEntry struct {
	Name string
	Func func()
}

type Menu struct {
	Entries []MenuEntry
}