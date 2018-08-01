/*
	Scyndi
	Use Parsing and Processing
	
	
	
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
Version: 18.08.01
*/
package scynt

import (
			"strings"
			"trickyunits/qff"
			"trickyunits/qstr"
)

func use_from_interface(trans map[string] *T_TransMod,s *tsource, blocks *map[string]string,module string){
}

func use(transm map[string] *T_TransMod,s *tsource, blocks *map[string]string,module string){
	trans:=transm[TARGET]
	file:=""
	if strings.ToLower(module)=="system" {
		SYSTEMDIR = strings.Replace(SYSTEMDIR,"\\","/",-1)
		file=SYSTEMDIR
		if !strings.HasSuffix(file,"/") { file += "/" }
		file += TARGET+".scf"
	} else {
		if strings.ToLower(module)!=module { warn( module+" appears to contain some uppercase characters. In order to assure full compatibility with case sensitive filesystems (such as Linux and other Unix based systems) it is strongly recommended to make module names in lower case only") }
		for _,epath:=range USEPATH {
			path:=strings.Replace(epath,"\\","/",-1)
			if qstr.Right(path,1)!="/" { path+="/" }
			if qff.IsFile(path+module+".ssf") {
				file=path+module+".ssf"
			} else if qff.IsFile(path+module+"/"+TARGET+".ssf") {
				file=path+module+"/"+TARGET+".ssf"
			} else if qff.IsFile(path+module+"/"+TARGET+".scf") {
				file=path+module+"/"+TARGET+".scf"
			} else if qff.IsFile(path+module+"/_other.scf") {
				file=path+module+"/_other.scf"
			} else if qff.IsFile(path+module+"/_other.ssf") {
				file=path+module+"/_other.ssf"
			} else if qff.IsFile(path+module+"/"+module+".ssf") {
				file=path+module+"/"+module+".ssf"
			} else if qff.IsFile(path+module+"/"+module+".scf") {
				file=path+module+"/"+module+".scf"
			}
		}
	}
	// Already used or not
	cmodule:=strings.ToUpper(module)
	if s.usedmap==nil { s.usedmap=&smap{}; s.usedmap.m=map[string]*tsource{} }
	m:=s.usedmap.m
	if us,ok:=m[cmodule];ok {
		for id,i:=range us.identifiers {
			if !i.private { s.allid[id]=i } //else { doing("IGNORED PRIVATE: ",id) }
		}
		for id,i:=range us.allid {
			if !i.private { s.allid[id]=i }
		}
		return
	} 
	// What, we don't have the file?
	if file=="" {
		doingln("The next directories are available for finding modules:","")
		for i,mp:=range USEPATH{
			sumdot(i+1); doingln("",mp)
		} 
		throw("No way found to use module: "+module) 
	}
	// Let's do it
	usetranslation,usesource:=CompileFile(file,"MODULE")
	s.usedmap.m[cmodule]=usesource
	s.used=append(s.used,usesource)
	e:=qff.WriteStringToFile(file+".scyndi.translation."+TARGET+"."+trans.extension, usetranslation)
	if e!=nil { ethrow(e) }
	usesource.usedmap=s.usedmap
	us:=usesource
	for id,i:=range us.identifiers {
		if !i.private { s.allid[id]=i } //else { doing("IGNORED PRIVATE: ",id) }
	}
	for id,i:=range us.allid {
		if !i.private { s.allid[id]=i }
	}

}

func useblock(trans map[string] *T_TransMod,s *tsource, blocks *map[string]string){
	//if s.usedmap==nil { s.usedmap=&smap{} }
	res:=[]string{}
	if s.srcname!="SYSTEM" { res=append(res,"SYSTEM") }
	for _,u:=range s.userequested { res=append(res,u) }
	for _,u:=range res { use(trans,s,blocks,u) }
}

