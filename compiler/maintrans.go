/*
	Scyndi
	Main
	
	
	
	(c) Jeroen P. Broks, 2018, All rights reserved
	
		This program is free software: you can redistribute it and/or modify
		it under the terms of the GNU General Public License as published by
		the Free Software Foundation, either version 3 of the License, or
		(at your option) any later version.
		
		This program is distributed in the hope that it will be useful,
		but WITHOUT ANY WARRANTY; without even the implied warranty of
		MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
		GNU General Public License for more details.
		You should have received a copy of the GNU General Public License
		along with this program.  If not, see <http://www.gnu.org/licenses/>.
		
	Exceptions to the standard GNU license are available with Jeroen's written permission given prior 
	to the project the exceptions are needed for.
Version: 18.08.11
*/
package scynt

import(
			"sort"
			"fmt"
			"trickyunits/mkl"
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
	TransTypes func(t *tsource) string
	TransTypeDefinition func(t *tsource,dtype *tidentifier,did *tidentifier) string
	TransTypeKill func(t *tsource,dtype *tidentifier,did *tidentifier) string
	Merge func(b map[string]string) string
	FuncHeaderRem func() string
	FuncHeader func(s *tsource,ch *tchunk) string
	EndFunc func(s *tsource,ch *tchunk,trueend bool) string
	StartFor func(fortype string,index *tidentifier,sxu,exu,step string,stepconstant bool) string
	startforeach func(eachi, fkey,fvalue *tidentifier,arrayORmap string,self *tsource,chf *tchunk,ol *tori) string
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
	simpleif string
	simpleelif string
	simpleelse string
	simpleendif string
	simplewhile string
	simpleendwhile string
	simpleloop string
	simpleinfloop string
	simpleuntilloop string
	simpleendfor string
	procnoneedbracket bool
	AltFuncCall func() (string,int) // more stuff to be added later!
	createindexvar func(indexedvariable string,indexedidentifier *tidentifier,sex string) (ivar string,iid *tidentifier)
	FuncEndless func(s *tsource,ol *tori,c *tchunk, epos *int,a *targ ,retargs []string) (r []string)
	FormVoidReturn func(s *tsource, c *tchunk, ol *tori) string
	FormFuncReturn func(s *tsource, c *tchunk, ol *tori, expression string) string
	nocasesupported bool // if set to "true", Scyndi will replace the casing sessions with if statements
	casefallsthrough bool // Scyndi will NOT support this, as most programming language don't, but Scyndi must be able to detect it properly and this this variable
	simplecasestart string
	simplecaseoneitem string
	simplecasemultiitem func(want string,ol *tori) string
	simplecaseend string
	MyNil string
	AltGetTypedIdentifier func() *tidentifier
	SimpleTypeSeparator string
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
mkl.Lic    ("Scyndi Programming Language - maintrans.go","GNU General Public License 3")
mkl.Version("Scyndi Programming Language - maintrans.go","18.08.11")

	dfo:=&defoperators
	for _,k := range ([]string{"==","+","-","/","*","^","!=","<",">","<=",">="}){
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
	
	(*dfo)["true"] ="true"
	(*dfo)["false"]="false"
}
