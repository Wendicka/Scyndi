package scynt


//import "fmt"
import (
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

func init(){

	
	TransMod["BlitzMax"] = &T_TransMod {}
	tmw:=TransMod["BlitzMax"]
	
	tmw.extension="bmx"
	tmw.constantsupport = true

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
	
	tmw.savetrans = func(s *tsource,translation,outp string){
		bigsource:="rem\n\tBlitzMax Source Generated by Scyndi\n\tPlease read the original coder's documentation about the\n\tCopyright notices and under which terms this code may\n\tbe distributed (if it may be distributed at al)\nend rem\n\n\n"
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
		ret:="' Code generated by Scyndi\n\nStrict\n"
		//ret:=b["USE"]+"\n\n"
		ret+=b["VAR"]+"\n\n"
		ret+=b["FUN"]+"\n\n"
		return ret
	}
}
