package scynt

type T_TransMod struct {
	//NameIdentifiers func(t *tsource)
	TransVars func(t *tsource) string
	Merge func(b map[string]string) string
}

var TransMod = map[string] *T_TransMod{}
