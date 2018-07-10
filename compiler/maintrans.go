package scynt

import(
			"sort"
)

type T_TransMod struct {
	constantsupport bool
	extension string
	//NameIdentifiers func(t *tsource)
	UseBlocks func(s *tsource, b *map[string]string)
	SealBlocks func(b *map[string]string)
	TransVars func(t *tsource) string
	Merge func(b map[string]string) string
	UsePureCode byte // 0 = Purecode PRIOR to translated code; 1 = PurseCode AFTER translated code; 2 = Let the translation module handle it by itself
	DontInterface bool // If set no interface files will be written, meaning the entire module will always be compiled whole
	ProcessUsed func(s *tsource, b *map[string]string,translation string) // If no function set, the imported code will just be added at the top of the translated file
	endlessargs bool
}

var TransMod = map[string] *T_TransMod{}


func CodeProcessUsed (s *tsource, b *map[string]string,translation string) {
	(*b)["USE"] += translation
}

func TargetsSupported() string {
	ts:=[]string{}
	for k,_ := range TransMod {
		ts=append(ts,k)
	}
	sort.Strings(ts)
	ret:=""
	for i:=0;i<len(ts);i++{
		if i<len(ts)-1 && i!=0 { ret +=", " } else if i>=len(ts)-1 { ret += " and " }
		ret += ts[i]
	}
	return ret
}
