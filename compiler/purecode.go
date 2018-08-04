/*
	Scyndi
	PureCode Inline Parsing and Processing
	
	
	
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
Version: 18.08.04
*/
package scynt

import(
//		"fmt"
		"strings"
		"trickyunits/mkl"
		"trickyunits/qff"
)

func purecode(s *tsource,c *tchunk,ol *tori) string {
	 /*
	fmt.Println(len(ol.sline) )
	for i,w:=range(ol.sline){
		fmt.Println(i,"\t",w.Wtype,"\t",w.Word)
	}
	// */
	if len(ol.sline)!=4 { ol.throw("Invalid PURECODE instruction") }
	tar:=ol.sline[1]
	comma:=ol.sline[2]
	code:=ol.sline[3]
	if comma.Word!="," { ol.throw("Comma expected") }
	if tar.Word!=TARGET { return "" } // Only do this if the target is correct
	ret:=code.Word
	if strings.HasPrefix(strings.ToUpper(ret),"IMPORT:"){
		fname:=ret[7:]
		doingln("Importing pure "+TARGET+" code: ",fname)
		r,e:=qff.EGetString(fname)
		if e!=nil { ol.ethrow(e) }
		ret=r
	}
	// replace locals *if* within a function
	if c!=nil {
		for k,i:=range c.locals {
			ret = strings.Replace(ret,"$C{"+k+"}C$",i.translateto,-1)
		}
	}
	// identifier replacement will have to take place here
	for k,i:=range s.identifiers {
		ret = strings.Replace(ret,"$C{"+k+"}C$",i.translateto,-1)
	}
	// Warn if some stuff has been found
	if strings.Index(ret,"$C{")>=0 || strings.Index(ret,"}$C")>=0 {
		ol.warn("Some $C{}$C may not have been properly substituted. Unknown identifiers?")
	}
	
	return ret
}

func init(){
mkl.Lic    ("Scyndi Programming Language - purecode.go","GNU General Public License 3")
mkl.Version("Scyndi Programming Language - purecode.go","18.08.04")
}
