package scynt

import(
//		"fmt"
		"strings"
)

func purecode(s *tsource,c *tchunk,ol *tori) string {
	 /*
	fmt.Println(len(ol.sline) )
	for i,w:=range(ol.sline){
		fmt.Println(i,"\t",w.Wtype,"\t",w.Word)
	}
	// */
	if len(ol.sline)!=4 { ol.throw("Invalid PURCODE instruction") }
	tar:=ol.sline[1]
	comma:=ol.sline[2]
	code:=ol.sline[3]
	if comma.Word!="," { ol.throw("Comma expected") }
	if tar.Word!=TARGET { return "" } // Only do this if the target is correct
	ret:=code.Word
	// replace locals *if* within a function
	if c!=nil {
		for k,i:=range c.locals {
			ret = strings.Replace(ret,"$C{"+k+"}C$",i.translateto,-1)
		}
	}
	// identifier replacement will have to take place here
	for k,i:=range s.identifiers {
		ret = strings.Replace(ret,"$C{"+k+"}C$",i.translateto,-1)
	}
	// Warn if some stuff has been found
	if strings.Index(ret,"$C{")>=0 || strings.Index(ret,"}$C")>=0 {
		ol.warn("Some $C{}$C may not have been properly substituted. Unknown identifiers?")
	}
	
	return ret
}
