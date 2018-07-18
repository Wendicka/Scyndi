package scynt
import(
							"strings"
							"fmt"
)

var fortotal = 0




func (self *tsource) callfunction() {
}

func (self *tsource)  translatefunctions() string{
	trans:=TransMod[TARGET]
	ret:=trans.FuncHeaderRem()
	for _,chf := range self.chunks {
		chf.forid = map[string]*tidentifier{}
		chf.fors  = map[int]bool{}
		ret += trans.FuncHeader(self,chf)
		//doingln("DEBUG: Translating func to: ",chf.translateto) // debug only
		for _,ins:=range chf.instructions {
			ol:=ins.ori
			pt:=ol.sline[0]
			for tab:=0;tab<ins.level;tab++{ 
				if !((pt.Word=="END" || pt.Word=="UNTIL" || pt.Word=="FOREVER" || pt.Word=="LOOP" || pt.Word=="ELSEIF" || pt.Word=="ELIF" || pt.Word=="ELSE") && tab==ins.level-1) {ret+="\t" }
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
							ret+=trans.plusone(idname)+"\n"
						case "--":
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
							ol.throw("Function calls not yet implemented! (coming soon)")
					}
				}
			} else if pt.Word=="IF" || pt.Word=="ELSEIF" || pt.Word=="ELIF" {
				if pt.Word!="IF" && ins.state.openinstruct!="IF" && ins.state.openinstruct!="IF block" { ol.throw(pt.Word+" can only be used within an IF block") }
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
			} else if pt.Word=="WHILE" {
				exp,exu:=self.translateExpressions("boolean", chf, ol,1,0)
				if exp<len(ol.sline) { 
					//echat("\tEnd:",exp,len(ol.sline))
					ol.throw("Separator expected") 
				}
				ret+=fmt.Sprintf(trans.simplewhile,exu)+"\n"
			} else if pt.Word=="FOR" || pt.Word=="FORU" {
				fortotal++
				if len(ol.sline)<2 { ol.throw("Index variable expected") }
				indexw:=ol.sline[1]
				indexn:=indexw.Word
				i:=2
				if indexw.Wtype!="identifier" { ol.throw("Unexpected "+indexw.Wtype+": "+indexw.Word) }
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
				fmt.Println(i,len(ol.sline)) // debug
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
				chf.forid[indexn]=&index
				ins.state.startfor=fortotal
				chf.fors[fortotal]=true
				ret+=trans.StartFor(pt.Word,&index,sxu,exu,step,stepconstant)+"\n"
			} else if pt.Word=="REPEAT" || pt.Word=="DO" {
				ret+=trans.simpleloop+"\n"
				if len(ol.sline)>1 { ol.throw("REPEAT and DO do not accept parameters") }
			} else if pt.Word=="END" || pt.Word=="FOREVER" || pt.Word=="LOOP" || pt.Word=="UNTIL" {
				//doingln("Ending:",ins.state.openinstruct) // debug only
				if (ins.state.openinstruct!="REPEAT" && ins.state.openinstruct!="REPEAT block" && ins.state.openinstruct!="DO" && ins.state.openinstruct!="DO block") && pt.Word!="END"{
					ol.throw("Keyword "+pt.Word+" may ONLY be used to close REPEAT/DO loops")
				}
				switch ins.state.openinstruct {
					case "PROCEDURE","PROC","VOID":
						ret += trans.EndFunc(self,chf,true)
					case "DEF","FUNCTION","FUNC":
						// Make sure there is a return in the end and add one if not.
						ret += trans.EndFunc(self,chf,true)
					case "FOR","FORU","FOREACH","FOR block","FORU block","FOREACH block":
						ret += trans.simpleendfor+"\n"
					case "IF","IF block":
						ret += trans.simpleendif+"\n"
					case "WHILE","WHILE block":
						ret += trans.simpleendwhile+"\n"
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
	}
	return ret // I must have this asa temp measyre or Go won't work (figures).
}
