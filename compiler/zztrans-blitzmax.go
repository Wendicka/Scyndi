package scynt


//import "fmt"
import (
			"fmt"
			"strings"
			"trickyunits/qstr"
			"trickyunits/qff"
)


func bmx_transtype(dttype string,o *tori) string{
				ret:=""
				switch dttype {
					case "STRING":  ret +=":String"
					case "INTEGER": ret +=":Long"
					case "BOOLEAN": ret +=":Byte"
					case "FLOAT":   ret +=":Double"
					case "VARIANT": o.throw("Variant types not allowed in BlitzMax translation!!")
					default:        o.throw("Alternate types ("+dttype+") are not yet supported for BlitzMax translation. Please come back later!")
				}
				return ret
}

const bmxbig=`Function CRASH(A$) 
	Print("Scyndi-RuntimeError:~t"+A) 
End Function
	
Function FIRange:Long[](start:Long,eind:Long,steps:Long=1,tu$="to")
	Local ret:Long[] 
	Local i:Long=start
	Local ok=True
	If steps=0 crash("Step value zero")
	Repeat
		'Print "Range: "+start+tu+eind+"; cyc:"+i
		Select tu.tolower()
			Case "until"
				If steps>0 And i>=eind Return ret
				If steps<0 And i<=eind Return ret
			Case "to"
				If steps>0 And i> eind Return ret
				If steps<0 And i< eind Return ret
		End Select
		ret=ret[..Len(ret)+1]
		ret[Len(ret)-1]=i
		i:+steps
	Forever	
End Function

Function FFRange:Double[](start:Double,eind:Double,steps:Double=1,tu$="to")
	Local ret:Double[] 
	Local i:Double=start
	Local ok=True
	If steps=0 crash("Step value zero")
	Repeat
		'Print "Range: "+start+tu+eind+"; cyc:"+i
		Select tu.tolower()
			Case "until"
				If steps>0 And i>=eind Return ret
				If steps<0 And i<=eind Return ret
			Case "to"
				If steps>0 And i> eind Return ret
				If steps<0 And i< eind Return ret
		End Select
		ret=ret[..Len(ret)+1]
		ret[Len(ret)-1]=i
		i:+steps
	Forever	
End Function
`

func init(){

	
	TransMod["BlitzMax"] = &T_TransMod {}
	tmw:=TransMod["BlitzMax"]
	
	tmw.extension="bmx"
	tmw.constantsupport = true
	
	tmw.int2float="Double(%s)"
	tmw.float2int="Long(%s)"
	tmw.iint2string="psf_int2str(%s)"
	tmw.iflt2string="psf_flt2str(%s)"
	
	tmw.simpleif="If %s"
	tmw.simpleelif="ElseIf %s"
	tmw.simpleendif="End If"
	
	tmw.simplewhile="While %s"
	tmw.simpleendwhile="Wend"
	
	tmw.simpleloop="Repeat"
	tmw.simpleinfloop="Forever"
	tmw.simpleuntilloop="Until %s"
	
	tmw.simpleendfor="Next"

	
	tmw.operators = defoperators // We need to have the basic. As this is not a pointer assignment, I can just modify all this :)
	tmw.operators[  "="]="="
	tmw.operators[ "=="]="="
	tmw.operators[ "!="]="<>"
	tmw.operators[ "<>"]="<>"
	tmw.operators["AND"]="And"
	tmw.operators[ "OR"]="Or"
	tmw.operators["NOT"]="Not"
	tmw.operators["MOD"]="Mod"
	
	tmw.procnoneedbracket=true

	// Trabskate global variables to BlitzMax
	tmw.TransVars = func(src *tsource) string{

		ret:="\n\n' Global Variables\n"
		for vname,vdata:=range src.identifiers {
			if vdata.idtype=="VAR" {
				if src.orilinerem { ret += "\t\t' VAR "+vname+"\n"}
				ret += "Global "+vdata.translateto
				switch vdata.dttype {
					case "STRING":  ret +=":String"
					case "INTEGER": ret +=":Long"
					case "BOOLEAN": ret +=":Byte"
					case "FLOAT":   ret +=":Double"
					case "VARIANT": throw("Variant types not allowed in BlitzMax translation")
					default:        throw("Alternate types ("+vdata.dttype+") are not yet supported for BlitzMax translation. Please come back later!!")
				}
				ret += " = "
				if vdata.defstring { ret += "\""+vdata.defaultvalue+"\"\n" } else { ret+=vdata.defaultvalue+"\n" }
			}
		}
		
		return ret
	}
	
	tmw.FuncHeaderRem = func() string { return "' Translated functions\n\n" }
	tmw.FuncHeader = func(s *tsource,ch *tchunk) string {
		ret := "Function "+ch.translateto+" "
		o   := ch.from
		if ch.pof==1 {
			//doingln("pof1","check") // debug only
			ret+=bmx_transtype(ch.returntype,o)+" "
		}
		ret += "( "
		for i,a:=range ch.args.a {
			if i>0 { ret+=" , " }
			v:=a.arg
			//doing("arg","check") // debug only
			ret += v.translateto+bmx_transtype(v.dttype,o)
			if !v.constant { ret += " var " }
		}
		ret +=" )\n"
		return ret
	}
	
	tmw.EndFunc = func(s *tsource,ch *tchunk,trueend bool) string {
		ret:=""
		if trueend { ret+="End Function\n\n" }
		return ret
	}
	
	tmw.plusone  = func(i interface{}) string { 
		if gt(i)=="*scynt.tidentifier" {
			return i.(*tidentifier).translateto+":+1 "
		} else if gt(i)=="string" {
			return i.(string)+":+1\t"
		} else {
			throw("INTERNAL ERROR: What the hell must a ++ request do with type: "+gt(i))
			return "error" // does nothing, but Go requires it!
		}
		
	}
	tmw.minusone = func(i interface{}) string { //return i.translateto+":-1 "}
		if gt(i)=="*scynt.tidentifier" {
			return i.(*tidentifier).translateto+":-1 "
		} else if gt(i)=="string" {
			return i.(string)+":-1\t"
		} else {
			throw("INTERNAL ERROR: What the hell must a -- request do with type: "+gt(i))
			return "error" // does nothing, but Go requires it!
		}
	}
	
	tmw.setstring = func(str string) string {
		ret:=[]byte{'"'}
		bb:=[]byte(str)
		for _,b:=range bb {
			switch b{
				case  13: ret=append(ret,'~'); ret=append(ret,'r')
				case  10: ret=append(ret,'~'); ret=append(ret,'n')
				case   8: ret=append(ret,'~'); ret=append(ret,'b')
				case   9: ret=append(ret,'~'); ret=append(ret,'t')
				case '"': ret=append(ret,'~'); ret=append(ret,'q')
				case '~': ret=append(ret,'~'); ret=append(ret,'~')
				default:
					if b>31 && b<127 {
						ret = append(ret,b)
					} else {
						tb:=[]byte(fmt.Sprintf("\"+Chr(%3d)+\"",b))
						for _,b2:=range(tb) { ret=append(ret,b2) }
					}
			}
		}
		ret = append(ret,'"')
		return string(ret)
	}
	
	tmw.setint = func(num string) string { return num }
	
	/*
	tmw.expressiontrans = func(ex *tex) string{
		ret:=""
		for _,ei:=range ex.exi {
			switch ei.wtype{
				case "identifier":	ret+=ei.id.translateto
				case "string":  	ret+=tmw.setstring(ei.word)
				case "integer": 	ret+=tmw.setint(ei.word)
				default: 			warn("What is a "+ei.wtype)
			}
		}
		return ret
	}
	
	tmw.definevar = func(s *tsource,id *tidentifier,ex *tex) string{
		return id.translateto+" = "+tmw.expressiontrans(ex)
	}
	*/

	tmw.definevar = func(s *tsource,id *tidentifier,ex string) string{
		return id.translateto+" = "+ex
	}

	
	tmw.savetrans = func(s *tsource,translation,outp string){
		bigsource:="rem\n\tBlitzMax Source Generated by Scyndi\n\tPlease read the original coder's documentation about the\n\tCopyright notices and under which terms this code may\n\tbe distributed (if it may be distributed at al)\nend rem\n\n\nstrict\n\n"
		bigsource+="function psf_int2str$(a:long  )\n\treturn  \"\"+a\nend function\n"
		bigsource+="function psf_flt2str$(a:double)\n\treturn  \"\"+a\nend function\n\n"
		bigsource+=bmxbig+"\n\n"
		for i,us:=range s.used {
			sumdot(i)
			doingln("Merging: ",us.srcname)
			bigsource +="' Module: "+us.srcname+"\n"
			file:=us.filename
			src,e:=qff.EGetString(file+".scyndi.translation."+TARGET+"."+tmw.extension)
			if e!=nil { ethrow(e) }
			bigsource+=src+"\n' End module: "+s.srcname+"\n\n"
		}
		doingln("Merging: ","Main Program")
		bigsource+="' ************\n"
		bigsource+="' MAIN PROGRAM\n"
		bigsource+="' ************\n"
		bigsource+=translation+"\n\n"
		if m,ok:=s.identifiers["MAIN"];ok{
			if m.idtype!="PROCEDURE" && m.idtype!="PROC" && m.idtype!="VOID" { throw("Identifier MAIN had to be a PROCEDURE, not a "+m.idtype) }
			bigsource+="' Call main\n"+m.translateto+"\n"
		} else {
			throw("No MAIN procedure in the program found")
		}
		outfile:=qstr.StripExt(s.filename)+".bmx"
		doingln("Saving: ",outfile)
		e:=qff.WriteStringToFile(outfile, bigsource)
		if e!=nil { ethrow(e) }
	}
	
	// Merge all code together
	// As BlitzMax is set to compile to pure machine language, I don't see much need to put neither USE nor XUSE in separate files... Let's dump it all together.
	tmw.Merge = func(b map[string]string) string {
		ret:="' Code generated by Scyndi\n\n\n"
		//ret:=b["USE"]+"\n\n"
		ret+=b["VAR"]+"\n\n"
		ret+=b["FUN"]+"\n\n"
		return ret
	}
	
	tmw.StartFor = func(fortype string,index *tidentifier,sxu,exu,step string,stepconstant bool) string {
		wto:="To"
		if fortype=="FORU" { wto="Until" }
		ret:="For Local "+index.translateto+bmx_transtype(index.dttype,nil)  +" = "
		// BlitzMax requires a constant for its step value, so this crazy way was needed :-/
		if stepconstant {
			ret += fmt.Sprintf("%s %s %s Step %s",sxu,wto,exu,step) 
		} else {
			rt:="FIRange"
			if index.dttype=="FLOAT" { rt="FFRange" }
			ret += fmt.Sprintf("Eachin %s( %s , %s , %s, \"%s\")",rt,sxu,exu,step,strings.ToLower(wto))
		}
		return ret
	}
}
