// License Information:
// 	Scyndi
// 	JavaScript export for Node, NW.JS and WebScripts
// 	
// 	
// 	
// 	(c) Jeroen P. Broks, 2018, 2019, All rights reserved
// 	
// 		This program is free software: you can redistribute it and/or modify
// 		it under the terms of the GNU General Public License as published by
// 		the Free Software Foundation, either version 3 of the License, or
// 		(at your option) any later version.
// 		
// 		This program is distributed in the hope that it will be useful,
// 		but WITHOUT ANY WARRANTY; without even the implied warranty of
// 		MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// 		GNU General Public License for more details.
// 		You should have received a copy of the GNU General Public License
// 		along with this program.  If not, see <http://www.gnu.org/licenses/>.
// 		
// 	Exceptions to the standard GNU license are available with Jeroen's written permission given prior 
// 	to the project the exceptions are needed for.
// Version: 19.02.08
// End License Information
package scynt

import(
	"os"
	"fmt"
	"strings"
	"trickyunits/mkl"
	"trickyunits/qff"
	"trickyunits/qstr"
)


func init(){
	mkl.Lic    ("Scyndi Programming Language - zztrans-javascript.go","GNU General Public License 3")
	mkl.Version("Scyndi Programming Language - zztrans-javascript.go","19.02.08")
	var tmw *T_TransMod;
	forNode:=&T_TransMod {}
	forNWJS:=&T_TransMod {}
	forSite:=&T_TransMod {}
	TransMod["NodeJS"]=forNode;
	TransMod["NW.JS"]=forNWJS;
	TransMod["WebJS"]=forSite;
	allJS:=[]*T_TransMod{forNode,forNWJS,forSite}
	
	// functions
	TransVars := func(src *tsource) string{
		ret:="\n\n// Global Variables (and constants)\n"
		for vname,vdata:=range src.identifiers {
			if vdata.idtype=="VAR" && (!vdata.imported) {
				if src.orilinerem { ret += "\t\t// VAR "+vname+"\n"}
				if vdata.constant { ret +="const "; } else { ret+="let "; }
				ret += vdata.translateto
				/*
				if vdata.constant{
					ret += "Const "+vdata.translateto
				} else {
					ret += "Global "+vdata.translateto
				}
				*/ 
				/*
				switch vdata.dttype {
					case "STRING":  ret +=":String"
					case "INTEGER": ret +=":Long"
					case "BOOLEAN": ret +=":Byte"
					case "FLOAT":   ret +=":Double"
					case "VARIANT": throw("Variant types not allowed in BlitzMax translation")
					default:        
									serial:=vdata.dttype
									sersp:=strings.Split(serial," ")
									if sersp[0]=="ARRAY"{
										if len(sersp)<=1 { throw("Array of what?") }
										switch sersp[1]{
											case "STRING":	ret += "ArrayOfString";	vdata.defaultvalue="New ArrayOfString"
											case "INTEGER":	ret += "ArrayOfLong";	vdata.defaultvalue="New ArrayOfLong"
											case "BOOLEAN":	ret += "ArrayOfByte";	vdata.defaultvalue="New ArrayOfByte"
											case "FLOAT":	ret += "ArrayOfDouble";	vdata.defaultvalue="New ArrayOfDouble"
											default:		ret += "ArrayOfObject";	vdata.defaultvalue="New ArrayOfObject"
										}
									} else {
										throw("Alternate types ("+vdata.dttype+") are not yet supported for BlitzMax translation. Please come back later!!")
									}
				}
				*/
				ret += " = "
				if qstr.Prefixed(vdata.dttype,"ARRAY ") || qstr.Prefixed(vdata.dttype,"ARRAY ")  { vdata.defaultvalue="[]" }
				if qstr.Prefixed(vdata.dttype,"MAP ") || qstr.Prefixed(vdata.dttype,"MAP ")  { vdata.defaultvalue="{}" }
				//if vdata.defstring { ret += "\""+vdata.defaultvalue+"\"\n" } else { ret+=vdata.defaultvalue+"\n" }
				if vdata.defstring { ret += tmw.setstring(vdata.defaultvalue)+"\n" } else { ret+=vdata.defaultvalue+"\n" }
			}
		}
		return ret
	}
	
	TransTypes := func(src *tsource) string {
		ret:=""
		for vname,vdata:=range src.identifiers {
			//fmt.Println("Type:",vname,vdata.idtype,vdata.idtype=="TYPE" && (!vdata.imported))
			if vdata.idtype=="TYPE" && (!vdata.imported) {
				ret += "function SCYNDI_NEW_"+vdata.translateto+"() { // type: "+vname+"\n"
				ret += "\tlet ret={};\n"
				for fldname,flddata := range vdata.typeidentifiers {
					switch flddata.idtype {
						case "VAR":
							dta := flddata.defaultvalue
							if flddata.dttype=="STRING" { dta=tmw.setstring(flddata.defaultvalue) }
							ret += fmt.Sprintf("\tret.%s = %s  // field %s \n",flddata.translateto,dta,fldname)
						case "PROCEDURE","FUNCTION":
							ret += fmt.Sprintf("\tret.%s = %s  // method %s\n",flddata.translateto,flddata.defaultvalue,fldname)
						default:
							throw("I cannot use a "+flddata.idtype+" as type field")

					}
				}
				ret += "\treturn ret\n}//end\n\n"
			}
		}
		return ret
	}

	TransTypeDefinition := func(t *tsource,dtype *tidentifier,did *tidentifier) string{
		if dtype.imported { throw("NEW not possible for imported types") }
		if did==nil { return "SCYNDI_NEW_"+dtype.translateto+"()" }
		return did.translateto+" = SCYNDI_NEW_"+dtype.translateto+"()"
	}

	TransTypeKill := func(t *tsource,dtype *tidentifier,did *tidentifier) string{
		if dtype.imported { throw("KILL not possible for imported types") }
		return did.translateto+" = undefined;"
	}
	
	plusone  := func(i interface{}) string { 
		if gt(i)=="*scynt.tidentifier" {
			return i.(*tidentifier).translateto+"++;"// = "+i.(*tidentifier).translateto+" + 1"
		} else if gt(i)=="string" {
			return i.(string)+"++;"// = "+i.(string)+" + 1"
		} else {
			throw("INTERNAL ERROR: What the hell must a ++ request do with type: "+gt(i))
			return "error" // does nothing, but Go requires it!
		}
		
	}
	minusone := func(i interface{}) string { //return i.translateto+":-1 "}
		if gt(i)=="*scynt.tidentifier" {
			return i.(*tidentifier).translateto+"--;"// = "+i.(*tidentifier).translateto+" - 1"
		} else if gt(i)=="string" {
			return i.(string)+"--;"// = "+i.(string)+" - 1"
		} else {
			throw("INTERNAL ERROR: What the hell must a -- request do with type: "+gt(i))
			return "error" // does nothing, but Go requires it!
		}
	}

	createindexvar := func(indexedvariable string,indexedidentifier *tidentifier,sex string) (ivar string,iid *tidentifier){
		//doingln("Indexing: ",indexedvariable) // debug
		iid = &tidentifier{}
		if indexedidentifier.dttype=="STRING" {
			ivar=indexedvariable+"["+sex+"]"
			iid.translateto=ivar
			iid.indexed=true
			iid.dttype="STRING"
			iid.idtype="VAR"
			return
		}
		st:=strings.Split(indexedidentifier.dttype," ")
		at:=""
		for _,t:=range st[1:]{
			if at!="" { at+=" " }
			at += t
		}
		/*
		switch st[0] {
			case "MAP":
				throw("Map indexing not yet supported in Lua")
			case "ARRAY":*/
				ivar=indexedvariable+"["+sex+"]"
				//doingln("Created: ",ivar)
				// Some more stuff will be needed when more complex stuff comes into play
				
				// Make fake identifier
				iid.translateto=ivar
				iid.dttype=at
				iid.indexed=true
				iid.idtype="VAR"
		//}
		return
	}

	FuncHeaderRem := func() string { return "// Translated functions\n\n" }
	
	FuncHeader := func(s *tsource,ch *tchunk) string {
		ret := "const "+ch.translateto+" = "
		//o   := ch.from
		ch.varpars = []int{}
		if ch.pof==1 {
			//doingln("pof1","check") // debug only
			//ret+=bmx_transtype(ch.returntype,o)+" "
		}
		ret += "( "
		for i,a:=range ch.args.a {
			if i>0 { ret+=" , " }
			v:=a.arg
			//doing("arg","check") // debug only
			ret += v.translateto //+bmx_transtype(v.dttype,o)
			if !v.constant { 
				ch.xreturn += v.translateto+","; ch.varpars=append(ch.varpars,i) 
				if ch.pof==1 { throw("Variable arguments not supported in functions when translating to Lua (you can use the when using procedures)") }
			}
		}
		if ch.args.endless!=nil {
			if len(ch.args.a)>0 { ret+=", " }
			ret+=ch.args.endless.arg.translateto+" "
		}
		ret +=" ) => {\n"
		return ret
	}
	
	EndFunc := func(s *tsource,ch *tchunk,trueend bool) string {
		ret:=""
		if len(ch.varpars)>0 {
			throw("Variable parameters in function not yet supported for JavaScript, but I am working on it.\nPlease see issue: https://github.com/Wendicka/Scyndi/issues/43")
		}
		if trueend { ret+="}\n\n" }
		return ret
	}

	Merge := func(b map[string]string) string {
		ret:="// Code generated by Scyndi\n\n\n"
		//ret:=b["USE"]+"\n\n"
		ret+=b["HEADPURECODE"]+"\n\n"
		ret+=b["TYPES"]+"\n\n"
		ret+=b["VAR"]+"\n\n"
		ret+=b["FUN"]+"\n\n"
		return ret
	}

	// Now JavaScript supports this by itself, but I am not full sure how to suppor that in combination with other supported languages
	FuncEndless := func(s *tsource,ol *tori,c *tchunk,epos *int, a *targ ,retargs []string) (ret []string) {
		/*
		ret=retargs
		add:="ArrayOf"
		t:=""
		switch a.a.dtype{
			case "INTEGER":	t = "Long"
			case "FLOAT":	t = "Double"
			case "BOOLEAN":	t = "Byte"
			case "STRING":	t = "String"
			default:
				ol.throw("No support for this kind of endless argument for the BlitzMax target (yet)!")
		}
		add += t + ".CreateArray("+t+"["
		for{
			if epos>=len(ol.sline) { break }
			
		}
		*/
		rets:="[ "
		want:=""
		switch a.arg.dttype {
			case "VARIANT":
				want="anything"
			case "INTEGER","FLOAT","BOOLEAN","STRING":
				want=strings.ToLower(a.arg.dttype)
			default:
				want=a.arg.dttype
		}
		first:=true
		cnt:=0
		for {
			if *epos>=len(ol.sline) { break }
			if !first {
				rets += ","
			} else { 
				first=false
			}
			ep,ee:=s.translateExpressions(want, c, ol,*epos,0)
			rets += fmt.Sprintf("/*[%d] =*/ ",cnt)+ee; cnt++
			*epos=ep
			//(*epos)++
			if *epos>=len(ol.sline) { break }
			if ol.sline[*epos].Word==")" { break }
			if ol.sline[*epos].Word!="," { ol.throw("Comma expected between two arguments") }
			(*epos)++
		}
		rets+=" ]"
		ret=append(retargs,rets)
		return
	}

	setstring := func(str string) string {
		ret:=[]byte{'"'}
		bb:=[]byte(str)
		for _,b:=range bb {
			switch b{
				case  13: ret=append(ret,'\\'); ret=append(ret,'r')
				case  10: ret=append(ret,'\\'); ret=append(ret,'n')
				case   8: ret=append(ret,'\\'); ret=append(ret,'b')
				case   9: ret=append(ret,'\\'); ret=append(ret,'t')
				case '"': ret=append(ret,'\\'); ret=append(ret,'"')
				case '\\': ret=append(ret,'\\'); ret=append(ret,'\\')
				default:
					if b>31 && b<127 {
						ret = append(ret,b)
					} else {
						q:=fmt.Sprintf("%x",b)
						if len(q)<2 { q="0"+q }
						tb:=[]byte(fmt.Sprintf("\\x%s",q))
						for _,b2:=range(tb) { ret=append(ret,b2) }
					}
			}
		}
		ret = append(ret,'"')
		return string(ret)
	}
	
	definevar := func(s *tsource,id *tidentifier,ex string) string{
		return id.translateto+" = "+ex
	}
	
	StartFor := func(fortype string,index *tidentifier,sxu,exu,step string,stepconstant bool) string {
		wto:="<="
		if fortype=="FORU" { wto="<" }
		//ret:= fmt.Sprintf("for %s in forloop(%s,%s,%s,'%s') do",index.translateto,sxu,exu,step,wto)
		ret:=fmt.Sprintf("for(let %s=%s;%s%s%s; %s+=%s){",index.translateto,sxu,  index.translateto,wto,exu,   index.translateto,step)
		return ret
	}
	
	startforeach := func(eachi, fkey,fvalue *tidentifier,arrayORmap string,self *tsource,chf *tchunk,ol *tori) string{
		/*
		f:=""
		switch arrayORmap{
			case "array":	f="scynipairs"
			case "map":		f="scynmpairs"
			default:		ol.throw("Internal error! ForEach misdefinition! ("+arrayORmap+"). Please report!")
		}
		ret:=fmt.Sprintf("for %s,%s in %s(%s) do ",fkey.translateto,fvalue.translateto,f,eachi.translateto)
		return ret
		*/
		
		return fmt.Sprintf("for (let %s in %s) { let %s=%s[%s]; ",fkey.translateto,eachi.translateto,fvalue.translateto,eachi.translateto,fkey.translateto)
		
	}


	// The stuff that is the same for all JavaScript output kinds.
	for _,js:=range allJS{
		tmw=js;
		// basis
		js.extension="js"
		js.constantsupport = true
		js.nocasesupported=true
		// if
		js.simpleif="if ( %s ) {"
		js.simpleelif="} else if ( %s ) {"
		js.simpleelse="} else {"
		js.simpleendif="}"
		// while
		js.simplewhile="while ( %s ) {"
		js.simpleendwhile="}"
		// repeat until
		js.simpleloop="do {"
		js.simpleinfloop="} while true"
		js.simpleuntilloop="} while(!( %s ))" // Please note that JavaScript loops as long as the condition is true, and Scyndi is set up to stop when the codition is true.
		// end of a for loop
		js.simpleendfor="}"
		// defoperators
		js.operators = NewDefOperators(); //defoperators
		// functions
		js.TransVars=TransVars;
		js.TransTypes=TransTypes;
		js.TransTypeDefinition=TransTypeDefinition;
		js.TransTypeKill=TransTypeKill;
		js.plusone=plusone
		js.minusone=minusone
		js.createindexvar=createindexvar;
		js.FuncHeaderRem=FuncHeaderRem;
		js.FuncHeader=FuncHeader;
		js.EndFunc=EndFunc;
		js.Merge=Merge;
		js.FuncEndless=FuncEndless;
		js.setstring=setstring;
		js.definevar=definevar;
		js.StartFor=StartFor;
		js.startforeach=startforeach
		

	}
	
	// saving and bundling
	saveprogram:=func(s *tsource,translation,altout string){
				bigsource:="/*\n\tJaveScript Source Generated by Scyndi\n\tPlease read the original coder's documentation about the\n\tCopyright notices and under which terms this code may\n\tbe distributed (if it may be distributed at al)\n*/\n\n\n"
				bigsource+="const TRUE=true; const FALSE=false;\n\n";
				//bigsource+=luabig+"\n\n"
				m:=0
				for i,us:=range s.used {
					m=i+2
					sumdot(i+1)
					doingln("Merging: ",us.srcname)
					bigsource +="// Module: "+us.srcname+"\n"
					file:=us.filename
					src,e:=qff.EGetString(file+".scyndi.translation."+TARGET+"."+tmw.extension)
					if e!=nil { ethrow(e) }
					bigsource+=src+"\n"
					if iinit,iok:=us.identifiers["INIT"];iok{
						if iinit.idtype!="PROCEDURE" && iinit.idtype!="PROC" && iinit.idtype!="VOID" {
							throw("Invalid INIT identifier")
						}
						bigsource+="\n// INIT\n"+iinit.translateto+"()\n\n"
					}
					bigsource+="\n// End module: "+us.srcname+"\n\n"
				}
				sumdot(m)
				doingln("Merging: ","Main Program")
				bigsource+="// ************\n"
				bigsource+="// MAIN PROGRAM\n"
				bigsource+="// ************\n"
				bigsource+=translation+"\n\n"
				if iinit,iok:=s.identifiers["INIT"];iok{
					if iinit.idtype!="PROCEDURE" && iinit.idtype!="PROC" && iinit.idtype!="VOID" {
						throw("Invalid INIT identifier")
					}
					bigsource+="\n// INIT\n"+iinit.translateto+"()\n\n"
				}
				if m,ok:=s.identifiers["MAIN"];ok{
					if m.idtype!="PROCEDURE" && m.idtype!="PROC" && m.idtype!="VOID" { throw("Identifier MAIN had to be a PROCEDURE, not a "+m.idtype) }
					bigsource+="// Call main\n"+m.translateto+"()\n"
				} else {
					throw("No MAIN procedure in the program found")
				}
				outfile:=qstr.StripExt(s.filename)+".js"
				if (altout!="") {outfile=altout}
				doingln("Saving: ",outfile)
				er:=qff.WriteStringToFile(outfile, bigsource)
				if er!=nil { ethrow(er) }
			
	}
	forNode.savetrans=func(s *tsource,translation,outp string){
		if s.srctype!="PROGRAM" { throw("NodeJS requires PROGRAM"); }
		saveprogram(s,translation,"");
	}
	
	forSite.savetrans=func(s *tsource,translation,outp string){
		outdir:=qstr.StripExt(s.filename)+".sout"
		outfile:=outdir+"/"+qstr.StripAll(s.filename)+".js"
		os.MkdirAll(outdir,0777)
		saveprogram(s,translation,outfile)
		doingln("Saving: ",outdir+"/index.html")
		qff.WriteStringToFile(outdir+"/index.html","<html> Please await contain later </html>")
	}
}
