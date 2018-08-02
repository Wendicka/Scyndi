/*
	Scyndi
	Lua Export
	
	
	
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
			"os"
			"fmt"
			"strings"
			"trickyunits/mkl"
			"trickyunits/qff"
			"trickyunits/qstr"
)

const luabig=`function CRASH(A) 
	error("Scyndi-RuntimeError:\t"..A) 
end
	
	function scynipairs(a)
		if type(a)~="table" then CRASH("Table expected for iteration!") end
		i=0
		t={}
		for j=0,#a do t[#t+1]=a[j] end
		return function()
			i=i+1
			if not t[i] then return nil,nil end
			return i-1,t[i]
		end
	end
	
	function scynmpairs(a)
		if type(a)~="table" then CRASH("Table expected for iteration!") end
		local tk = {}
		local tv = {}
		for k,v in pairs(a) do
			tk[#tk+1]=k
			tv[#tv+1]=v
		end
		i=0
		return function()
			i=i+1
			return tk[i],tv[i]
		end
	end
	
	function i2scyntab(a)
		if type(a)~="table" then CRASH("Table expected for conversion!") end
		ret = {}
		for i=1,#a do ret[i-1]=a[i] end
		return ret
	end
	
	function forloop(start,einde,stappen,atu)
		local tu = atu or "to"
		local i,e,s=start,einde,stappen
		return function()
			if tu:lower()=="until" and i>=e then return nil,nil end
			if i>e then return nil,nil end
			local o=i
			i = i + s
			return o
		end
	end
	
`

func init(){
mkl.Lic    ("Scyndi Programming Language - zztrans-lua.go","GNU General Public License 3")
mkl.Version("Scyndi Programming Language - zztrans-lua.go","18.08.02")

	
	TransMod["Lua"] = &T_TransMod {}
	tmw:=TransMod["Lua"]
	
	tmw.extension="lua"
	tmw.constantsupport = false
	
	tmw.simpleif="if %s then"
	tmw.simpleelif="elseif %s then"
	tmw.simpleelse="else"
	tmw.simpleendif="end"
	
	tmw.simplewhile="while %s do"
	tmw.simpleendwhile="end"
	
	tmw.simpleloop="repeat"
	tmw.simpleinfloop="until false"
	tmw.simpleuntilloop="until %s"
	
	tmw.simpleendfor="end"
	
	tmw.operators = defoperators // We need to have the basic. As this is not a pointer assignment, I can just modify all this :)
	tmw.operators[  "="]="=="
	tmw.operators[ "=="]="=="
	tmw.operators[ "~="]="~="
	tmw.operators[ "!="]="~="
	tmw.operators[ "<>"]="~="
	tmw.operators["AND"]="and"
	tmw.operators[ "OR"]="or"
	tmw.operators["NOT"]="not"
	
	tmw.operators["concat"]=".."


	tmw.procnoneedbracket=false
	
	tmw.TransVars = func(src *tsource) string{

		ret:="\n\n-- Global Variables (and constants)\n"
		for vname,vdata:=range src.identifiers {
			if vdata.idtype=="VAR" {
				if src.orilinerem { ret += "\t\t-- VAR "+vname+"\n"}
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
				if qstr.Prefixed(vdata.dttype,"ARRAY ") || qstr.Prefixed(vdata.dttype,"ARRAY ")  { vdata.defaultvalue="{}" }
				//if vdata.defstring { ret += "\""+vdata.defaultvalue+"\"\n" } else { ret+=vdata.defaultvalue+"\n" }
				if vdata.defstring { ret += tmw.setstring(vdata.defaultvalue)+"\n" } else { ret+=vdata.defaultvalue+"\n" }
			}
		}
		return ret
	}
	
	tmw.plusone  = func(i interface{}) string { 
		if gt(i)=="*scynt.tidentifier" {
			return i.(*tidentifier).translateto+" = "+i.(*tidentifier).translateto+" + 1"
		} else if gt(i)=="string" {
			return i.(string)+" = "+i.(string)+" + 1"
		} else {
			throw("INTERNAL ERROR: What the hell must a ++ request do with type: "+gt(i))
			return "error" // does nothing, but Go requires it!
		}
		
	}
	tmw.minusone = func(i interface{}) string { //return i.translateto+":-1 "}
		if gt(i)=="*scynt.tidentifier" {
			return i.(*tidentifier).translateto+" = "+i.(*tidentifier).translateto+" - 1"
		} else if gt(i)=="string" {
			return i.(string)+" = "+i.(string)+" - 1"
		} else {
			throw("INTERNAL ERROR: What the hell must a -- request do with type: "+gt(i))
			return "error" // does nothing, but Go requires it!
		}
	}

	tmw.createindexvar = func(indexedvariable string,indexedidentifier *tidentifier,sex string) (ivar string,iid *tidentifier){
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


	tmw.FuncHeaderRem = func() string { return "-- Translated functions\n\n" }

	tmw.FuncHeader = func(s *tsource,ch *tchunk) string {
		ret := "function "+ch.translateto+" "
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
		ret +=" )\n"
		return ret
	}
	
	
	tmw.EndFunc = func(s *tsource,ch *tchunk,trueend bool) string {
		ret:=""
		if len(ch.varpars)>0 {
			throw("Variable parameters in function not yet supported for Lua, but I am working on it.\nPlease see issue: https://github.com/Wendicka/Scyndi/issues/43")
		}
		if trueend { ret+="end\n\n" }
		return ret
	}
	
	tmw.Merge = func(b map[string]string) string {
		ret:="-- Code generated by Scyndi\n\n\n"
		//ret:=b["USE"]+"\n\n"
		ret+=b["VAR"]+"\n\n"
		ret+=b["FUN"]+"\n\n"
		return ret
	}

	tmw.FuncEndless = func(s *tsource,ol *tori,c *tchunk,epos *int, a *targ ,retargs []string) (ret []string) {
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
		rets:="{ "
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
			rets += fmt.Sprintf("[%d] = ",cnt)+ee; cnt++
			*epos=ep
			//(*epos)++
			if *epos>=len(ol.sline) { break }
			if ol.sline[*epos].Word==")" { break }
			if ol.sline[*epos].Word!="," { ol.throw("Comma expected between two arguments") }
			(*epos)++
		}
		rets+=" }"
		ret=append(retargs,rets)
		return
	}

	tmw.setstring = func(str string) string {
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
	
	tmw.definevar = func(s *tsource,id *tidentifier,ex string) string{
		return id.translateto+" = "+ex
	}


	tmw.savetrans = func(s *tsource,translation,outp string){
		switch s.srctype{
			case "PROGRAM":
				bigsource:="--[[\n\tLua Source Generated by Scyndi\n\tPlease read the original coder's documentation about the\n\tCopyright notices and under which terms this code may\n\tbe distributed (if it may be distributed at al)\n]]\n\n\n"
				bigsource+=luabig+"\n\n"
				for i,us:=range s.used {
					sumdot(i)
					doingln("Merging: ",us.srcname)
					bigsource +="-- Module: "+us.srcname+"\n"
					file:=us.filename
					src,e:=qff.EGetString(file+".scyndi.translation."+TARGET+"."+tmw.extension)
					if e!=nil { ethrow(e) }
					bigsource+=src+"\n-- End module: "+us.srcname+"\n\n"
				}
				doingln("Merging: ","Main Program")
				bigsource+="-- ************\n"
				bigsource+="-- MAIN PROGRAM\n"
				bigsource+="-- ************\n"
				bigsource+=translation+"\n\n"
				if m,ok:=s.identifiers["MAIN"];ok{
					if m.idtype!="PROCEDURE" && m.idtype!="PROC" && m.idtype!="VOID" { throw("Identifier MAIN had to be a PROCEDURE, not a "+m.idtype) }
					bigsource+="-- Call main\n"+m.translateto+"()\n"
				} else {
					throw("No MAIN procedure in the program found")
				}
				outfile:=qstr.StripExt(s.filename)+".lua"
				doingln("Saving: ",outfile)
				e:=qff.WriteStringToFile(outfile, bigsource)
				if e!=nil { ethrow(e) }
			case "SCRIPT":
				outdir:=qstr.MyTrim(outp)
				if outdir=="" { outdir=qstr.StripExt(s.filename)+".luabuild" }
				doingln("Creating build bundle: ",outdir)
				if qff.IsFile(outdir) { throw("Target bundle exists as a file!") }
				if qff.IsDir(outdir) { 
					warn("Target bundle does exist! Content may get overwritten!") 
				} else {
					e:=os.MkdirAll(outdir,0777)
					if e!=nil { ethrow(e) }
				}
				mods:=[]string{}
				dotty:=0
				for i,us:=range s.used {
					sumdot(i)
					doingln("Bundling: ",us.srcname)
					//bigsource +="-- Module: "+us.srcname+"\n"
					file:=us.filename
					src,e:=qff.EGetString(file+".scyndi.translation."+TARGET+"."+tmw.extension)
					if e!=nil { ethrow(e) }
					//bigsource+=src+"\n-- End module: "+s.srcname+"\n\n"
					outfile := "used_"+us.srcname 
					mods=append(mods,outfile)
					e=qff.WriteStringToFile(outdir+"/"+outfile+".lua",src)
					if e!=nil { ethrow(e) }
					dotty=i+1
				}
				sumdot(dotty)
				doingln("Bundling: ",s.srcname)
				outfile := "main_"+s.srcname+".lua"
				//mods=append(mods,outfile)
				for _,m:=range mods { translation="require \""+m+"\"\n"+translation }
				er:=qff.WriteStringToFile(outdir+"/"+outfile,translation)
				if er!=nil { ethrow(er) }
				
		}
	}
	
	tmw.StartFor = func(fortype string,index *tidentifier,sxu,exu,step string,stepconstant bool) string {
		wto:="To"
		if fortype=="FORU" { wto="Until" }
		ret:= fmt.Sprintf("for %s in forloop(%s,%s,%s,'%s') do",index.translateto,sxu,exu,step,wto)
		return ret
	}
	
	tmw.startforeach = func(eachi, fkey,fvalue *tidentifier,arrayORmap string,self *tsource,chf *tchunk,ol *tori) string{
		f:=""
		switch arrayORmap{
			case "array":	f="scynipairs"
			case "map":		f="scynmpairs"
			default:		ol.throw("Internal error! ForEach misdefinition! ("+arrayORmap+"). Please report!")
		}
		ret:=fmt.Sprintf("for %s,%s in %s(%s) do ",fkey.translateto,fvalue.translateto,f,eachi.translateto)
		return ret
	}


}
