package scynt

import(
			"strings"
)

type texi struct{
	word string
	wtype string
	convstring bool
	id *tidentifier
	subex *tex
}

type tex struct{
	oriline *tori
	orisrc *tsource
	exi []*texi
	cptype string
	tottype string 
}

func (s *tsource) translateExpressions(expect string,ol *tori,start,level int) (endpos int,ex *tex){
	endpos=start
	ex=&tex{}
	ex.exi=[]*texi{}
	ex.cptype=expect
	ex.tottype=expect
	cov:=false
	ucov:=false
	endex:=false
	if expect=="BOOLEAN" { ex.cptype="" }
	doingln("Expression! ",expect)
	for {
		//cnt:=endpos-start
		if endpos>=len(ol.sline) { ol.throw("Expression expected, but received nothing at all") }
		p:=ol.sline[endpos]
		ucov=p.Wtype=="identifier" || p.Wtype=="string" || p.Wtype=="float" || p.Wtype=="integer"
		if cov && ucov { ol.throw("Operator expected") }
		if p.Word=="(" {
		} else {
			if ucov {
				exi:=texi{}
				exi.word=p.Word
				exi.wtype=p.Wtype
				if p.Wtype=="identifier" {
					exi.id=s.GetIdentifier(p.Word,nil,ol)
					if ex.cptype=="" { ex.cptype=exi.id.dttype }
					switch  ex.cptype {
						case "STRING":
							exi.convstring=true
						case "INTEGER":
							if exi.id.dttype!="INTEGER" { ol.throw("Integer expected. Variable is "+exi.id.dttype) }
					}
				} else {
					if ex.cptype=="" { ex.cptype=strings.ToUpper(p.Wtype) }
					switch  ex.cptype {
						case "STRING":
							exi.wtype="string"
					}
				}
				ex.exi=append(ex.exi,&exi)
				doingln("exa: ",p.Word)		// debug
			}
		}
		// This always be last
		cov=ucov
		endpos++
		if endpos>=len(ol.sline) { 
			endex=true
		} else {
			commacheck:=ol.sline[endpos]
			endex=commacheck.Word==","
		}
		if endex {
			if cov { return }
			ol.throw("Unexpected end of expression")
		}
	}
}
