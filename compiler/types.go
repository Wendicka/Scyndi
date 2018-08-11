/*
	Scyndi Programming Language
	Type Organizer
	
	
	
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

import (
			"trickyunits/mkl"
)

var translatingtype *tidentifier


func init(){
mkl.Lic    ("Scyndi Programming Language - types.go","GNU General Public License 3")
mkl.Version("Scyndi Programming Language - types.go","18.08.11")
}


func (s *tsource) starttype(ol *tori) bool {
	translatingtype = &tidentifier{}
	t:=translatingtype // I am lazy, I know
	if len(ol.sline)<2 { ol.throw("Type without identifier!") }
	iname:=ol.sline[1]
	if iname.Wtype!="identifier" { ol.throw("Indifier expected to name the type. Not a "+iname.Wtype) }
	if _,ok:=s.identifiers[iname.Word];ok { ol.throw("I can't create type "+iname.Word+"; Duplicate identifier") }
	if len(ol.sline)>2 { ol.throw("more advance types are not yet supported! Please come back later for those") }
	s.identifiers[iname.Word]=t
	t.idtype="TYPE"
	t.typeidentifiers=map[string]*tidentifier{}
	t.translateto="SCYNDITYPE"+s.srcname+"_"+iname.Word
	return true
}

func (s *tsource) totype(ol *tori) {
	t:=translatingtype // I am lazy, I know
	fname:=""
	//ftype:="" // temp not needed and Go is a whining language!
	if len(ol.sline)<3 { ol.throw("Illegal field declaration") }
	// Parse checkup
	for i,w:=range ol.sline{
		switch i{
			case 0:	fname=w.Word; if w.Wtype!="identifier" { ol.throw("Illegal field name") }
			case 1: if w.Word!=":" { ol.throw("':' expected") }
			default:
			/* The var declaration routine can deal with this itself
				if i<len(ol.sline)-1{
					if w.Word!="ARRAY" && w.Word!="MAP" { ol.throw("Invalid field type") }
					ftype += w.Word+" "
				} else {
					if w.Word=="ARRAY" && w.Word=="MAP" { ol.throw("Invalid field array/map type") }
					ftype += w.Word
				}
				*/
		}
	}
	if _,ok:=t.typeidentifiers[fname]; ok { ol.throw("Duplicate field identifier: "+fname) }
	_,v:=s.declarevar(ol.sline)
	v.translateto=t.translateto+"__"+fname
	t.typeidentifiers[fname]=&v
	//doingln("field: ",fname)
}

