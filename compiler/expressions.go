/*
	Scyndi
	Expression parser
	
	
	
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
Version: 18.07.21
*/
package scynt

import(
							"fmt"
)

const haveechat = false
var rtt string // ReTurned Type -- Needed for recognition of intruction starting identifiers
var rti *tidentifier // 

func echat(a ...interface{}){
	if !haveechat { return }
	for _,aa:=range a{ fmt.Print(aa," ") }
	fmt.Print("\n")
}


func defaultexpressiontranslation(expect string,source *tsource, c *tchunk, ol *tori,start,level int) (endpos int,ex string){
	echat("Expression:","\n\tExpeciting:",expect,"\n\tStarting at word:",start)
	endpos=start
	ex=""
	einde:=false
	cexpect:=expect
	haakjelevel:=0
	trans:=TransMod[TARGET]
	wantcov:=true // cov = constant or variable
	float2int:="%s"
	int2float:="%s"
	iint2string:="%s"
	iflt2string:="%s"
	cop:=""
	if  trans.float2int!=""  { float2int  =trans. float2int }
	if  trans.int2float!=""  { int2float  =trans. int2float }
	if trans.iint2string!="" { iint2string=trans.iint2string }
	if trans.iflt2string!="" { iflt2string=trans.iflt2string }
	timeout:=int64(2000000000)
	for{
		timeout--
		if timeout<0 { ol.throw("Expression parsing time-out! This could be an internal error") }
		if endpos>=len(ol.sline) { break }
		sexi:=ol.sline[endpos] // Scyndi Expression Index
		if einde { break }
		echat("\tChecking word:",endpos,"\t"+sexi.Wtype+":",sexi.Word)
		echat("\tiddefend:",expect=="identifier" && (sexi.Word!="[") && (endpos!=start))
		if expect=="identifier" && (sexi.Word!="[") && (endpos!=start) { break }
		if wantcov{
			if sexi.Word=="(" { 
				haakjelevel++ 
				ex += "("
			} else {
				if expect=="identifier" && sexi.Wtype!="identifier" { ol.throw("Unexpected "+sexi.Wtype+": "+sexi.Word) }
				switch sexi.Wtype{
					case "identifier":
						id:=source.GetIdentifier(sexi.Word,c,ol)
						out:=id.translateto
						// maybe some type checkups can come here
						if cexpect=="string" && id.dttype=="INTEGER" { out=fmt.Sprintf(iint2string,id.translateto) }
						if cexpect=="string" && id.dttype=="FLOAT"   { out=fmt.Sprintf(iflt2string,id.translateto) }
						// output
						ex += out
						wantcov = false
						rtt=id.dttype
						rti=id
					case "integer","float":
						if cexpect!="string" {
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
						if cexpect=="integer" || expect=="float" { ol.throw("Strings may not be used when nummeric expressions are expected") }
						ex += trans.setstring(sexi.Word)
						wantcov = false
					case "keyword":
						switch sexi.Word{
							case "NOT":
								if expect=="identifier" { ol.throw("That is NOT where you place NOT") }
								if cexpect!="boolean" { ol.throw("Keyword 'NOT' only works in BOOLEAN expressions!") }
								ex += trans.operators["NOT"]
							default:
								ol.throw(fmt.Sprintf("Unexpected keyword '%s' in expression. Identifier expected.",sexi.Word))
						}
					default:
						ol.throw(fmt.Sprintf("Unexpected %s '%s' in expression. Identifier expected.",sexi.Wtype,sexi.Word))
				}
			}
		} else {
			cop=""
			if sexi.Word=="," { break }
			if sexi.Word=="[" {ol.throw("Array indexes and map keys are not YET supported, they will be taken care of as soon as possible (Hey! Rome wasn't built in one day, either, ya know)") 
			} else if sexi.Word=="(" {ol.throw("Function calls are not YET supported in expressions, they will be taken care of as soon as possible (Hey! Rome wasn't built in one day, either, ya know)") 
			} else if sexi.Word==")" {
				if haakjelevel==0 {break}
				haakjelevel--
				ex += ")"
			} else if sexi.Word=="+" && cexpect=="string" {
				cop="concat"
				wantcov=true
			} else if (sexi.Word=="-" || sexi.Word=="+" || sexi.Word=="*" || sexi.Word=="MOD" || sexi.Word=="/") && (cexpect!="integer" && cexpect!="float") {
				ol.throw("Unexpected mathematical operation")
			} else {
				cop=sexi.Word
				wantcov=true
			}
			if cop!="" { 
				if _,ok:=trans.operators[cop];ok{
					ex += " "+trans.operators[cop]+" " 
				} else {
					if haveechat {
						echat("operators:")
						for k,v :=range(trans.operators){ echat("\t",k,"=",v) }
					}
					ol.throw("Unknown operator: "+cop)
				}
			}
		}
		endpos++ // MUST always be last before everything reloops
	}
	if wantcov { echat(start,endpos); ol.throw("Unexpected end of expression") }
	if haakjelevel>1  { ol.throw(fmt.Sprintf("There are %s brackets in this expression that are not properly closed, yet the expression has ended")) }
	if haakjelevel==1 { ol.throw(fmt.Sprintf("There is 1 bracket in this expression that is not properly closed, yet the expression has ended")) }
	return
}


func (s *tsource) translateExpressions(expect string, c *tchunk, ol *tori,start,level int) (endpos int,ex string){
	trans:=TransMod[TARGET]
	et:=defaultexpressiontranslation
	if trans.transexp!=nil { et = trans.transexp }
	endpos,ex = et(expect,s,c,ol,start,level)
	return
}
