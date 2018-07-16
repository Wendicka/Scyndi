package scynt


func (self *tsource) callfunction() {
}

func (self *tsource)  translatefunctions() string{
	trans:=TransMod[TARGET]
	ret:=trans.FuncHeaderRem()
	for _,chf := range self.chunks {
		ret += trans.FuncHeader(self,chf)
		//doingln("DEBUG: Translating func to: ",chf.translateto) // debug only
		for _,ins:=range chf.instructions {
			ol:=ins.ori
			pt:=ol.sline[0]
			if pt.Wtype=="identifier" {
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
								einde,ex:=self.translateExpressions(id.dttype,ol,2,0)
								if einde<len(ol.sline) { ol.throw("unexpected stuff after definition") }
								ret+=trans.definevar(self,id,ex)+"\n"
							default: ol.throw("Operator not expected in this particular situation: "+op.Word)
						}
					}
				}
			} else if pt.Word=="END" {
				//doingln("Ending:",ins.state.openinstruct) // debug only
				switch ins.state.openinstruct {
					case "PROCEDURE","PROC","VOID":
						ret += trans.EndFunc(self,chf,true)
					case "DEF","FUNCTION","FUNC":
						// Make sure there is a return in the end and add one if not.
						ret += trans.EndFunc(self,chf,true)
					default:
						ol.throw("I do not yet know how to end the "+ins.state.openinstruct+"; Either a bug or you are still working with an experimental version?")
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
