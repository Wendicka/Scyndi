// License Information:
// 	Scyndi
// 	Function processing
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
// Version: 19.02.13
// End License Information
/*
	Scyndi
	Function processing
	
	
	
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
							"trickyunits/qstr"
							"trickyunits/mkl"
							"strings"
							"fmt"
)

var fortotal = 0




func (self *tsource) callfunction(c *tchunk, ol *tori, mustreturn bool, funpos int) (cf string, epos int){
	cf=""
	epos=funpos
	trans:=TransMod[TARGET]
	if trans.AltFuncCall!=nil { return trans.AltFuncCall() }
	id:=rti
	//doing("Call to ",id.translateto)
	if id.idtype!="PROCEDURE" && id.idtype!="FUNCTION" { ol.throw(rtt+" cannot be called as a function") }
	epos++
	//if mustreturn { epos++ } // ignore bracket we don't need here.
	tvargs:=[]string{}
	ending:=len(id.args.a)+funpos+((len(id.args.a)*2)-1)
	nomore:=false
	// regular paramters
	for ac,aa:=range id.args.a {
		//fmt.Println(fmt.Sprintf("%5d: Processing arg %s #%d  (nomore: %t)",ol.ln,id.translateto,ac,nomore))
		want:=""
		switch aa.arg.dttype {
			case "VARIANT":
				want="anything"
			case "INTEGER","FLOAT","BOOLEAN","STRING":
				want=strings.ToLower(aa.arg.dttype)
			default:
				want=aa.arg.dttype
		}
		if nomore {
			if !aa.optional { ol.throw(fmt.Sprintf("Missing parameter! %d expected but I only have %d!",len(id.args.a),ac)) }
			tvargs=append(tvargs,aa.arg.defaultvalue)
		} else {
			ep,eu:=self.translateExpressions(want, c, ol,epos,0)
			//fmt.Println(eu)
			tvargs=append(tvargs,eu)
			epos=ep
			////fmt.Println(epos,ending)
			if epos+1>=len(ol.sline) {
				//fmt.Println("epos+1 tb",epos,len(ol.sline))
				nomore=true
			} else if epos<=ending { 
				c:=ol.sline[epos]
				if c.Word==")" {
					nomore=true
				} else if c.Word!="," { 
					// Define optional paramters if present
					// else
					ol.throw("Comma expected") 
				}
				epos++
			} else { 
				//fmt.Println("???",ending)
				nomore = true 
			}
		}
	}
	// Infinite parameters
	if id.args.endless!=nil{
		tvargs=trans.FuncEndless(self,ol,c,&epos,id.args.endless,tvargs)
	}
	if !mustreturn {cf=id.translateto} // (In "must return" cases the function name is already in the expression!)
	if mustreturn || (!trans.procnoneedbracket) { cf +="(" } else {cf +=" "}
	for ai,at:=range tvargs{
		////fmt.Println(ai,at)
		if ai!=0 { cf += ", " }
		cf += at
	}
	if mustreturn || (!trans.procnoneedbracket) { cf +=")" } 
	
	return 
}

func (self *tsource)  translatefunctions() string{
	trans:=TransMod[TARGET]
	ret:=trans.FuncHeaderRem()
	returned:=false
	casevar:=map[int] *tidentifier { }
	semi:="" // This can be defined by any language requiring ; after every instruction.
	for _,chf := range self.chunks {
		ended:=false
		chf.forid      = map[int] map[string] *tidentifier {}
		chf.fors       = map[int] bool{}
		chf.scopeid    = map[int] map[string] *tidentifier {}
		chf.scopeindex = map[int] uint32 {}
		chf.forline2ins = map[int]*tinstruction{}
		chf.newscope(0);
		ret += trans.FuncHeader(self,chf)
		//doingln("DEBUG: Translating func to: ",chf.translateto) // debug only
		for _,ins:=range chf.instructions {
			ol:=ins.ori
			pt:=ol.sline[0]
			for tab:=0;tab<ins.level;tab++{ 
				if !((pt.Word=="END" || pt.Word=="UNTIL" || pt.Word=="FOREVER" || pt.Word=="LOOP" || pt.Word=="ELSEIF" || pt.Word=="ELIF" || pt.Word=="ELSE" || pt.Word=="DEFAULT" || pt.Word=="CASE") &&  tab==ins.level-1) {ret+="\t" } 
				if (!(pt.Word=="END" || pt.Word=="UNTIL" || pt.Word=="FOREVER" || pt.Word=="LOOP" || pt.Word=="ELSEIF" || pt.Word=="ELIF" || pt.Word=="ELSE" || pt.Word=="CASE" || pt.Word=="DEFAULT")) && returned {
					ol.warn("This line comes after a return command while no proper ending took place.\n\t\tMost programming languages will simply ignore this line, some may even throw an error!")
				} 	
			}
			if pt.Wtype=="identifier" {
				/* old code
				//id,idfound:=self.identifiers[pt.Word]
				//if !idfound { ol.throw("Primary identifier unknown: "+pt.Word) }
				id:=self.GetIdentifier(pt.Word,chf,ol)
				if len(ol.sline)>1 {
					op:=ol.sline[1]
					if op.Wtype=="operator" {
						switch op.Word {
							case "++":
								if id.constant { ol.throw("Constants cannot be redefined") }
								if len(ol.sline)>2 { ol.throw("Invalid increment request") }
								ret += trans.plusone(id)+"\n"
							case "--":
							if id.constant { ol.throw("Constants cannot be redefined") }
								if len(ol.sline)>2 { ol.throw("Invalid decrement request") }
								ret += trans.minusone(id)+"\n"
							case "=",":=":
								if id.constant { ol.throw("Constants cannot be redefined") }
								//einde,ex:=
								self.translateExpressions(strings.ToLower(id.dttype),chf,ol,2,0)
								//if einde<len(ol.sline) { ol.throw("unexpected stuff after definition") }
								//ret+=trans.definevar(self,id,ex)+"\n"
							default: ol.throw("Operator not expected in this particular situation: "+op.Word)
						}
					}
				}
				*/
				pos,idname:=self.translateExpressions("identifier", chf, ol,0,0)
				if pos<len(ol.sline) {
					nxt:=ol.sline[pos]
					switch nxt.Word{
						case "++":
							if rti.dttype!="INTEGER" && rti.dttype!="FLOAT" { ol.throw("Incorrect increment type") }
							if rti.idtype!="VAR" || rti.constant { ol.throw("I can only use ++ on variables") }
							ret+=trans.plusone(idname)+"\n"
						case "--":
							if rti.idtype!="VAR" || rti.constant { ol.throw("I can only use -- on variables") }
							if rti.dttype!="INTEGER" && rti.dttype!="FLOAT" { ol.throw("Incorrect deccrement type") }
							ret+=trans.minusone(idname)+"\n"
						case ":+","+=":
							ol.throw(":+/+= altering not yet supported! (coming soon)")
						case ":-","-=":
							ol.throw(":-/-= altering not yet supported! (coming soon)")
						case "=",":=":
							id:=rti
							if id.constant { ol.throw("Constants cannot be redefined") }
							exp,exu:=self.translateExpressions(strings.ToLower(id.dttype), chf, ol,pos+1,0)
							if exp<len(ol.sline) { 
								echat("\tEnd:",exp,len(ol.sline))
								ol.throw("Separator expected") 
							}
							ret+=trans.definevar(self,id,exu)+"\n"
						default:
							//ol.throw("Function calls not yet implemented! (coming soon)")
							scall,spos:=self.callfunction(chf,ol,false,0)
							ret+=scall+"\n"
							pchat(fmt.Sprintf("%d",spos)) // just compiler distraction... for now
					}
				} else {
					scall,spos:=self.callfunction(chf,ol,false,0)
					ret+=scall+"\n"
					pchat(fmt.Sprintf("%d",spos)) // just compiler distraction... for now
				}
			} else if pt.Word=="NEW" {
				exp,_:=self.translateExpressions("identifier", chf, ol,1,0)
				//fmt.Println(exp,len(ol.sline))
				if exp<len(ol.sline) { ol.throw("Separator expected") }
				id:=rti
				if id.idtype=="INTEGER" || id.idtype=="FLOAT" || id.idtype=="STRING" || id.idtype=="BOOL" || qstr.Prefixed(id.idtype,"ARRAY ") || qstr.Prefixed(id.idtype,"MAP ") { ol.throw("Invalid type for NEW") }
				idtype:=self.GetIdentifier(id.dttype,chf,ol)
				mynew:=trans.TransTypeDefinition(self,idtype,id)
				ret += mynew+semi+"\n"
			} else if pt.Word=="VAR" {
				if len(ol.sline)<2 { ol.throw("Local variable declarations cannot (yet) come in var blocks") }
				n:=ol.sline[1];
				if n.Wtype!="identifier" { ol.throw("Identifier expected for local variable declaration. Not a "+n.Wtype) }
				_,v := self.declarevar(ol.sline[1:])
				v.translateto = fmt.Sprintf("SCYNDI_LOCAL_SCOPE%d_%s",chf.scopeindex[ins.level],n.Word)
				//fmt.Printf("- %s translated to %s should have value \"%s\"\n",n.Word,v.translateto,v.defaultvalue) // debug line
				ret += trans.TransLocal(self,&v)+"\n"; 
			} else if pt.Word=="KILL" {
				exp,_:=self.translateExpressions("identifier", chf, ol,1,0)
				if exp<len(ol.sline) { ol.throw("Separator expected") }
				id:=rti
				if id.idtype=="INTEGER" || id.idtype=="FLOAT" || id.idtype=="STRING" || id.idtype=="BOOL" || qstr.Prefixed(id.idtype,"ARRAY ") || qstr.Prefixed(id.idtype,"MAP ") { ol.throw("Invalid type for KILL") }
				idtype:=self.GetIdentifier(id.dttype,chf,ol)
				mynew:=trans.TransTypeKill(self,idtype,id)
				ret += mynew+semi+"\n"
			} else if pt.Word=="RETURN" {
				returned=true
				if chf.pof==0 {
					if len(ol.sline)>1 { ol.throw("Procedures don't take parameters for return") }
					if trans.FormVoidReturn==nil {
						ret+="return"+semi+"\n"
					} else {
						ret+=trans.FormVoidReturn(self,chf,ol)
					}
				} else {
					if len(ol.sline)==1 { ol.throw("Function returns require values or expressions to be returned") }
					exp,exu:=self.translateExpressions(strings.ToLower(chf.returntype), chf, ol,1,0)
					if exp<len(ol.sline) { 
						ol.throw("Separator expected") 
					}
					if trans.FormVoidReturn==nil {
						ret+="return "+exu+semi+"\n"
					} else {
						ret+=trans.FormFuncReturn(self,chf,ol,exu)+"\n"
					}
				}
			} else if pt.Word=="IF" || pt.Word=="ELSEIF" || pt.Word=="ELIF" {
				if pt.Word!="IF" && ins.state.openinstruct!="IF" && ins.state.openinstruct!="IF block" { ol.throw(pt.Word+" can only be used within an IF block") }
				chf.newscope(ins.level+1)
				exp,exu:=self.translateExpressions("boolean", chf, ol,1,0)
				if exp<len(ol.sline) { 
					//echat("\tEnd:",exp,len(ol.sline))
					ol.throw("Separator expected") 
				}
				if pt.Word=="IF" {
					ret+=fmt.Sprintf(trans.simpleif,exu)+"\n"
				} else {
					ret+=fmt.Sprintf(trans.simpleelif,exu)+"\n"
				}
			} else if pt.Word=="ELSE" {
				chf.newscope(ins.level+1)
				if ins.state.openinstruct!="IF" && ins.state.openinstruct!="IF block" { ol.throw(pt.Word+" can only be used within an IF block") }
				ret += trans.simpleelse+"\n"
			} else if pt.Word=="SELECT" || pt.Word=="SWITCH" {
				chf.newscope(ins.level+1)
				if len(ol.sline)!=2 { ol.throw("Invalid SWITCH/SELECT!") }
				ag:=ol.sline[1]
				if ag.Wtype!="identifier" { ol.throw("Identifier expected for SWTICH/SELECT but I got "+ag.Wtype+": "+ag.Word) }
				id:=self.GetIdentifier(ag.Word,chf,ol)
				if id.idtype!="VAR" { ol.throw(ag.Word+" is an illegal identifier for casing") }
				if id.dttype!="STRING" && id.dttype!="INTEGER" && id.dttype!="FLOAT" && id.dttype!="BOOLEAN" { ol.throw(ag.Word+" is an uncasable variable") }
				if id.constant { ol.warn("I can case "+ag.Word+", but please keep in mind that it's a constant and not a variable!") }
				casevar[ins.level+1]=id //.translateto
				if !trans.nocasesupported { 
					ol.throw("Call it irony, but languages that support casing systems themselves are not yet supported. Come back later") 
				} else {
					ret += fmt.Sprintf(trans.simpleif,trans.operators["false"])+"\n"
				}
			} else if pt.Word=="CASE" {
				returned=false
				if ins.state.openinstruct!="SELECT" && ins.state.openinstruct!="SELECT block" && ins.state.openinstruct!="SWITCH" && ins.state.openinstruct!="SWITCH block" { ol.throw(pt.Word+" can only be used within an SWITCH block") }
				if len(ol.sline)!=2 { ol.throw("Invalid CASE!") }
				id,ok:=casevar[ins.level]
				wantword:=ol.sline[1]
				want:=""
				if !ok { ol.throw("Internal error! Something went wrong with the CASE-variable preserveation! Please report!") }
				switch id.dttype {
					case "STRING": 
						if wantword.Wtype!="string" { ol.throw("I want a string for this CASE!") }
						want=trans.setstring(wantword.Word)
					case "INTEGER":
						if wantword.Wtype!="integer" { ol.throw("I want an integer for this CASE!") }
						want=wantword.Word
					case "FLOAT":
						if wantword.Wtype=="integer" {
							want=wantword.Word+".0"
						} else if wantword.Wtype=="float" {
							want=wantword.Word
						} else { ol.throw("I want a number for this CASE!") }
					case "BOOLEAN":
						switch wantword.Word{
							case "TRUE":	want=trans.operators["true"]
							case "FASLE":	want=trans.operators["false"]
							default:		ol.throw("TRUE of FALSE exptected for this CASE!") 
						}
					default:
						ol.throw("CASE DEFINTION ERROR")
				}
				if !trans.nocasesupported { 
					ol.throw("Call it irony, but languages that support casing systems themselves are not yet supported. Come back later") 
				} else {
					tex := id.translateto + trans.operators["="] + want
					ret += fmt.Sprintf(trans.simpleelif,tex)+"\n"
				}
			} else if pt.Word=="DEFAULT" {
				returned=false
				if ins.state.openinstruct!="SELECT" && ins.state.openinstruct!="SELECT block" && ins.state.openinstruct!="SWITCH" && ins.state.openinstruct!="SWITCH block" { ol.throw(pt.Word+" can only be used within an SWITCH block") }
				if !trans.nocasesupported { 
					ol.throw("Call it irony, but languages that support casing systems themselves are not yet supported. Come back later") 
				} else {
					ret += trans.simpleelse+"\n"
				}				
			} else if pt.Word=="WHILE" {
				chf.newscope(ins.level+1)
				exp,exu:=self.translateExpressions("boolean", chf, ol,1,0)
				if exp<len(ol.sline) { 
					//echat("\tEnd:",exp,len(ol.sline))
					ol.throw("Separator expected") 
				}
				ret+=fmt.Sprintf(trans.simplewhile,exu)+"\n"
			} else if pt.Word=="FOREACH" {
				chf.newscope(ins.level+1)
				fortotal++
				chf.forline2ins[ol.ln]=ins
				//exp,eache:=self.translateExpressions("chain", chf, ol,1,0)
				exp:=1
				eache:=ol.getword(exp)
				if eache.Wtype!="identifier" { ol.throw("identifier expected as chainvar") }
				eachi:=self.GetIdentifier(eache.Word,chf,ol)
				exp++
				if exp>=len(ol.sline) { ol.throw("FOREACH without iteration variables") }
				if ol.sline[exp].Word!="," { ol.throw("Comma expected") }
				exp++
				a1:=ol.getword(exp)
				if a1.Wtype!="identifier" { ol.throw("identifier expected as key var") }
				var a2 *tword
				if exp<len(ol.sline) { 
					exp++
					if ol.sline[exp].Word!="," { ol.throw("Comma or separator expected") }
					exp++
					a2=ol.getword(exp)
					if a1.Wtype!="identifier" { ol.throw("identifier expected as value var") }
				}
				eacht:=eachi.dttype
				tsplit:=strings.Split(eacht," ")
				fkey:=&tidentifier {} //self.declarevar(div)
				fval:=&tidentifier {} //self.declarevar(div)
				fkey.idtype="VAR"
				fkey.translateto=fmt.Sprintf("SCYNDI_FOR%X_KEY",fortotal)
				fval.translateto=fmt.Sprintf("SCYNDI_FOR%X_VAL",fortotal)
				if _,ok:=chf.forid[fortotal];!ok { chf.forid[fortotal] = map[string]*tidentifier {} }
				chf.fors[fortotal]=true
				switch tsplit[0]{
					case "STRING":
						ol.throw("Stringsplace FOREACH not yet supported, but it IS planned")
					case "ARRAY":
						fkey.dttype="INTEGER"
						if a2==nil { 
							chf.forid[fortotal][a1.Word]=fval
						} else {
							chf.forid[fortotal][a1.Word]=fkey
							chf.forid[fortotal][a2.Word]=fval
							//fmt.Println(a1.Word," becomes ",fkey.translateto," in ",fortotal)
						}
						ret += trans.startforeach(eachi,fkey,fval,"array",self,chf,ol)+"\n"
					case "MAP":
						fkey.dttype="STRING"
						if a2==nil { 
							chf.forid[fortotal][a1.Word]=fkey
						} else {
							chf.forid[fortotal][a1.Word]=fkey
							chf.forid[fortotal][a2.Word]=fval
						}
						ret += trans.startforeach(eachi,fkey,fval,"map",self,chf,ol)+"\n"
					default: ol.throw(eacht+" is not a valid type to use for FOREACH")
				}
				
			} else if pt.Word=="FOR" || pt.Word=="FORU" {
				chf.newscope(ins.level+1)
				fortotal++
				chf.forline2ins[ol.ln]=ins
				if len(ol.sline)<2 { ol.throw("Index variable expected") }
				indexw:=ol.sline[1]
				indexn:=indexw.Word
				i:=2
				if indexw.Wtype!="identifier" { ol.throw("Unexpected "+indexw.Wtype+": "+indexw.Word+" in FOR") }
				if len(ol.sline)<3 { ol.throw("Unexpected end of FOR-loop-definition") }
				indextype:="INTEGER"
				dp:=ol.sline[i]
				if dp.Word==":"{ // Alternate type
					i++
					if len(ol.sline)<i+1 { ol.throw("Type required") }
					t:=ol.sline[i]
					indextype=t.Word
				}
				if indextype!="INTEGER" && indextype!="FLOAT" { ol.throw("FOR only accepts INTEGERS or FLOATS as index") }
				// Start value
				i++
				dp=ol.getword(i)
				if dp.Word=="=" || dp.Word==":=" { i++ ; dp=ol.getword(i)}
				sxp,sxu:=self.translateExpressions(strings.ToLower(indextype), chf, ol,i,0)
				i=sxp
				// end value
				dp=ol.getword(i)
				//fmt.Println(dp.Word) // debug
				if dp.Word!="," { ol.throw(", expected") }
				i++
				exp,exu:=self.translateExpressions(strings.ToLower(indextype), chf, ol,i,0)
				i=exp
				// step value if present, if not use 1
				step:="1"
				stepconstant:=true
				//fmt.Println(i,len(ol.sline)) // debug
				if i<len(ol.sline) {
				//} else {
					dp=ol.getword(i)
					if dp.Word!="," { ol.throw(", expected or nothing at all") }
					i++
					if i>=len(ol.sline) { ol.throw("Unexpected end of line. Step value expected") } 
					if i+1==len(ol.sline) {
						stepi:=ol.getword(i)
						if stepi.Wtype==strings.ToLower(indextype) { 
							step=stepi.Word 
						} else if stepi.Wtype=="identifier" {
							v:=self.GetIdentifier(stepi.Word,chf,ol)
							step=v.translateto
							stepconstant=v.constant
						} else { 
							ol.throw("Unexpected "+stepi.Wtype+": "+stepi.Word+"; Step-value was expected")
						}
					} else {
						stxp,stxu := self.translateExpressions(strings.ToLower(indextype), chf, ol,i,0)
						step=stxu
						stepconstant=false
						if stxp>len( ol.sline) { ol.throw("Separator expected") }
					}
				}
				dname:=tword{}; dname.Word=indexn; dname.Wtype="identifier"
				ddpnt:=tword{}; ddpnt.Word=":"; ddpnt.Wtype="operator"
				dtype:=tword{}; dtype.Word=indextype; dtype.Wtype="keyword"				
				div:=[]*tword{&dname,&ddpnt,&dtype}
				_,index:=self.declarevar(div)
				index.translateto=fmt.Sprintf("SCYNDI_FOR%X_INDEX",fortotal)
				if _,ok:=chf.forid[fortotal];!ok { chf.forid[fortotal] = map[string]*tidentifier {} }
				chf.forid[fortotal][indexn]=&index
				ins.state.startfor=fortotal
				chf.fors[fortotal]=true
				ret+=trans.StartFor(pt.Word,&index,sxu,exu,step,stepconstant)+"\n"
			} else if pt.Word=="REPEAT" || pt.Word=="DO" {
				chf.newscope(ins.level+1)
				ret+=trans.simpleloop+"\n"
				if len(ol.sline)>1 { ol.throw("REPEAT and DO do not accept parameters") }
			} else if pt.Word=="END" || pt.Word=="FOREVER" || pt.Word=="LOOP" || pt.Word=="UNTIL" {
				//doingln("Ending:",ins.state.openinstruct) // debug only
				returned=false
				if (ins.state.openinstruct!="REPEAT" && ins.state.openinstruct!="REPEAT block" && ins.state.openinstruct!="DO" && ins.state.openinstruct!="DO block") && pt.Word!="END"{
					ol.throw("Keyword "+pt.Word+" may ONLY be used to close REPEAT/DO loops")
				}
				chf.endscope(ins.level);
				switch ins.state.openinstruct {
					case "PROCEDURE","PROC","VOID":
						ret += trans.EndFunc(self,chf,true)
						ended=true
					case "DEF","FUNCTION","FUNC":
						// Make sure there is a return in the end and add one if not.
						ret += trans.EndFunc(self,chf,true)
						ended=true
					case "FOR","FORU","FOREACH","FOR block","FORU block","FOREACH block":
						ret += trans.simpleendfor+"\n"
						//fmt.Println(ins.state.openline,chf.forline2ins)
						//for i,_:=range chf.forline2ins { fmt.Println("I have:",i) }
						fins:=chf.forline2ins[ins.state.openline]
						chf.fors[fins.state.startfor]=false
						//fmt.Println("Close FOR #",fins.state.startfor)
					case "IF","IF block":
						ret += trans.simpleendif+"\n"
					case "WHILE","WHILE block":
						ret += trans.simpleendwhile+"\n"
					case "SWITCH","SWITCH block","SELECT","SELECT block":
						if !trans.nocasesupported {
							ret += trans.simplecaseend+"\n"
						} else {
							ret += trans.simpleendif+"\n"
						}
					case "DO","REPEAT","DO block","REPEAT block":
						switch pt.Word{
							case "END","LOOP","FOREVER":
								ret += trans.simpleinfloop+"\n"
							case "UNTIL":
								exp,exu:=self.translateExpressions("boolean", chf, ol,1,0)
								if exp<len(ol.sline) { 
									//echat("\tEnd:",exp,len(ol.sline))
									ol.throw("Separator expected") 
								}
								ret += fmt.Sprintf(trans.simpleuntilloop,exu)+"\n"
							default:
								ol.throw("huh?") // This should NEVER be possible to happen but, hey, who knows :P
						}
					default:
						ol.throw("I do not yet know how to end the "+ins.state.openinstruct+"!\nEither a bug or you are still working with an experimental version?")
				}
			} else if pt.Word=="PURECODE" {
				ret += purecode(self,chf,ol)+"\n"
			} else {
				ol.throw("Unexpected "+pt.Wtype+" ("+pt.Word+")")
			}
		}
		if !ended {
			ei:=chf.instructions[len(chf.instructions)-1]
			throw(fmt.Sprintf("%s in line %d not properly closed!",ei.state.openinstruct,ei.state.openline))
		}
		chf.endscope(0)
	}	
	return ret // I must have this asa temp measyre or Go won't work (figures).
}

func init(){
mkl.Version("Scyndi Programming Language - functions.go","19.02.13")
mkl.Lic    ("Scyndi Programming Language - functions.go","GNU General Public License 3")
}
