package scynt

import (
			"strings"
			"trickyunits/qff"
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
	if file=="" { throw("No way found to use module: "+module) }
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

