package scynt

import(
							"fmt"
)

func defaultexpressiontranslation(expect string,source *tsource, c *tchunk, ol *tori,start,level int) (endpos int,ex string){
	endpos=start
	ex=""
	einde:=false
	haakjelevel:=0
	trans:=TransMod[TARGET]
	wantcov:=true // cov = constant or variable
	float2int:="%s"
	int2float:="%s"
	if trans.float2int!="" { float2int=trans.float2int }
	if trans.int2float!="" { int2float=trans.int2float }
	timeout:=int64(2000000000)
	for{
		timeout--
		if timeout<0 { ol.throw("Expression parsing time-out! This could be an internal error") }
		if endpos>=len(ol.sline) { break }
		sexi:=ol.sline[endpos] // Scyndi Expression Index
		if einde { break }
		if expect=="identifier" && (sexi.Word!="[") { break }
		if wantcov{
			if sexi.Word=="(" { 
				haakjelevel++ 
			} else {
				switch sexi.Wtype{
					case "identifer":
						id:=source.GetIdentifier(sexi.Word,c,ol)
						// maybe some type checkups can come here
						// output
						ex += id.translateto
						wantcov = false
					case "integer","float":
						if expect!="string" {
							switch sexi.Wtype { // yes again, and for good reasons, trust me 
								case "integer":
									ex += fmt.Sprintf(float2int,sexi.Word)
								case "float":
									ex += fmt.Sprintf(int2float,sexi.Word)
							}
							wantcov = false
							break;
						}
						sexi.Wtype="string"
						fallthrough
					case "string":
						if expect=="integer" || expect=="float" { ol.throw("Strings may not be used when nummeric expressions are expected") }
						ex += trans.setstring(sexi.Word)
						wantcov = false
					case "keyword":
						switch sexi.Word{
							case "NOT":
								if expect!="boolean" { ol.throw("Keyword 'NOT' only works in BOOLEAN expressions!") }
								ex += trans.operators["NOT"]
							default:
								ol.throw(fmt.Sprintf("Unexpected keyword '%s' in expression. Identifier expected.",sexi.Word))
						}
					default:
						ol.throw(fmt.Sprintf("Unexpected %s '%s' in expression. Identifier expected.",sexi.Wtype,sexi.Word))
				}
			}
		} else {
			if sexi.Word=="[" {ol.throw("Array indexes and map keys are not YET supported, they will be taken care of as soon as possible (Hey! Rome wasn't built in one day, either, ya know)") }
			if sexi.Word=="(" {ol.throw("Function calls are not YET supported in expressions, they will be taken care of as soon as possible (Hey! Rome wasn't built in one day, either, ya know)") }
			if sexi.Word==")" {
				if haakjelevel==0 {break}
				haakjelevel--
			}
		}
		endpos++ // MUST always be last before everything reloops
	}
	if wantcov { ol.throw("Unexpected end of expression") }
	if haakjelevel>1  { ol.throw(fmt.Sprintf("There are %s brackets in this expression that are not properly closed, yet the expression has ended")) }
	if haakjelevel==1 { ol.throw(fmt.Sprintf("There is 1 bracket in this expression that are not properly closed, yet the expression has ended")) }
	return
}


func (s *tsource) translateExpressions(expect string, c *tchunk, ol *tori,start,level int) (endpos int,ex string){
	trans:=TransMod[TARGET]
	et:=defaultexpressiontranslation
	if trans.transexp!=nil { et = trans.transexp }
	endpos,ex = et(expect,s,c,ol,start,level)
	return
}
