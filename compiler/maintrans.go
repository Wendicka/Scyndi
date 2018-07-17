package scynt

import(
			"sort"
			"fmt"
)

var defoperators = map[string] string {}


func gt(i interface{}) string {
	return fmt.Sprintf("%T",i)
}


type T_TransMod struct {
	constantsupport bool
	extension string
	//NameIdentifiers func(t *tsource)
	UseBlocks func(s *tsource, b *map[string]string)
	SealBlocks func(b *map[string]string)
	TransVars func(t *tsource) string
	Merge func(b map[string]string) string
	FuncHeaderRem func() string
	FuncHeader func(s *tsource,ch *tchunk) string
	EndFunc func(s *tsource,ch *tchunk,trueend bool) string
	savetrans func(s *tsource,trans, outp string)
	plusone func(i interface{}) string
	minusone func(i interface{}) string
	setstring func(str string) string
	setint func(str string) string
	transexp func (expect string,source *tsource, c *tchunk, ol *tori,start,level int) (endpos int,ex string)
	//expressiontrans func(ex *tex) string
	//definevar func(s *tsource,id *tidentifier,ex *tex) string
	definevar func(s *tsource,id *tidentifier,ex string) string
	UsePureCode byte // 0 = Purecode PRIOR to translated code; 1 = PurseCode AFTER translated code; 2 = Let the translation module handle it by itself
	DontInterface bool // If set no interface files will be written, meaning the entire module will always be compiled whole
	ProcessUsed func(s *tsource, b *map[string]string,translation string) // If no function set, the imported code will just be added at the top of the translated file
	operators map[string] string
	endlessargs bool
	int2float string
	float2int string
	iint2string string
	iflt2string string
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



func init(){
	dfo:=&defoperators
	for _,k := range ([]string{"==","+","-","/","*","^","!="}){
		(*dfo)[k]=k
	}
	(*dfo)["NOT"]="!"
	(*dfo)["AND"]="&&"
	(*dfo)["OR" ]="||"
	(*dfo)["MOD"]="%"
	
	(*dfo)["="] ="=="
	(*dfo)["<>"]="!="
	(*dfo)["~="]="!="
	
	(*dfo)["concat"]="+"
}
