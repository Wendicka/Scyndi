package scynt

type T_TransMod struct {
	NameIdentifiers func(t *TPackage)
}

var TransMod = map[string] *T_TransMod{}
