/*
	Scyndi
	Constants translator
	
	
	
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
Version: 18.08.02
*/
package scynt

import(
			"fmt"
			"trickyunits/mkl"
astr		"trickyunits/ansistring"
sc			"strconv"
)

var enumv int64 // value
var enums int64 // step
var enumb bool  // bitwise?
var enumo bool  // overflow

func (s *tsource) defconst(line []*tword){
	if len(line)!=3 { 
		fmt.Println(astr.SCol("HUH?",astr.A_Red,0))
		for i:=0;i<len(line);i++ {
			c:=astr.A_Magenta
			switch i{
				case 0:
					c=astr.A_Yellow
				case 1:
					c=astr.A_Red
				case 2:
					c=astr.A_Cyan
			}
			fmt.Print(astr.SCol(line[i].Word,c,0)+" ")
		}
		fmt.Println()
		throw("Invalid constant definition!") 
	}
	if line[1].Word!="=" && line[1].Word!=":=" { 
		
		throw("= or := expected in the middle of a constant definition") 
	}
	idname:=line[0]
	value:=line[2]
	if idname.Wtype!="identifier" { throw("I cannot translate "+idname.Word+" as a constant because it's a "+idname.Wtype) }
	name:="SCYNDI_S"+s.srcname+"_CONSTANT_"+idname.Word
	id:=&tidentifier{}
	id.translateto=name
	id.constant=true
	id.defaultvalue=value.Word
	id.private=s.private
	id.idtype="VAR" // Looks odd, but this is the easiest way for the rest of the parsing routines to handle this... The constant field already makes sure redefinition isn't possible!
	switch value.Wtype{
		case "string":
			id.dttype="STRING"
			id.defstring=true
		case "float":
			id.dttype="FLOAT"
		case "integer":
			id.dttype="INTEGER"
		case "keyword":
			if value.Word=="TRUE" || value.Word=="FALSE" {
				id.dttype="BOOLEAN"
				break
			}
			fallthrough
		default:
			throw("I cannot accept a "+value.Wtype+" as a value for constant "+idname.Word)
	}
	if _,ok:=s.identifiers[idname.Word];ok { throw("I cannot create constant "+idname.Word+"; Duplicate identifier") }
	s.identifiers[idname.Word]=id
}

func (s *tsource) enumstart(ol *tori){
	enumo=false
	enums=1
	switch len(ol.sline){
		case 1:
			enumv=1
			enums=1
			enumb=false
			break
		case 3:
			if ol.sline[2].Wtype!="integer" { ol.throw("Unexpected "+ol.sline[2].Wtype+" ("+ol.sline[2].Word+"). Integer expected for enum step value") }
			enums,_=sc.ParseInt(ol.sline[2].Word,10,64)
			fallthrough
		case 2:
			if ol.sline[1].Word=="BITWISE" {
					enumv=1
					enumb=true
					if enums!=1 { ol.throw("Steps cannot be set in bitwise enum") }
					break
			}
			if ol.sline[1].Wtype!="integer" { ol.throw("Unexpected "+ol.sline[1].Wtype+" ("+ol.sline[1].Word+"). Integer expected for enum start value") }
			enumv,_=sc.ParseInt(ol.sline[1].Word,10,64)
		default:
			ol.throw("I don't understand this enum request")
	}
}

func (s *tsource) enumid(ol *tori){
	w1:=ol.sline[0]
	if w1.Word=="ENUM" { s.enumstart(ol); return }
	if w1.Word=="END" { return }
	for _,h:=range ol.sline {
		if enumo { ol.throw("Enum Overflow") }
		if h.Wtype!="identifier" { ol.throw("Unexpected "+h.Wtype+" ("+h.Word+") in enumbering sequence") }
		idname:=h
		name:="SCYNDI_S"+s.srcname+"_ENUMCONSTANT_"+idname.Word
		id:=&tidentifier{}
		id.dttype="INTEGER"
		id.translateto=name
		id.constant=true
		id.defaultvalue=fmt.Sprintf("%9d",enumv)
		id.private=s.private
		id.idtype="VAR" // Looks odd, but this is the easiest way for the rest of the parsing routines to handle this... The constant field already makes sure redefinition isn't possible!
		if _,ok:=s.identifiers[idname.Word];ok { ol.throw("I cannot create constant "+idname.Word+" for enummeration; Duplicate identifier") }
		s.identifiers[idname.Word]=id
		if enumb {
			if enumv>=2147483648 { ol.warn("Bitwise enumberation cannot go past 2147483648 and this value has now been assigned to: "+idname.Word); enumo=true }
			enumv += enumv
		} else {
			if 2147483648-enumv < enums { ol.warn("Enumberation overflow will be there after constant: "+idname.Word); enumo=true }
			enumv += enums
		}
	}
}

func init(){
mkl.Lic    ("Scyndi Programming Language - constants.go","GNU General Public License 3")
mkl.Version("Scyndi Programming Language - constants.go","18.08.02")
}
