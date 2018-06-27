package scynt

import (
		"trickyunits/qff"
)

// Basically all this routine does is make sure all instructions are properly separated.
// New lines (0x0A) and ; can do this, although those who prefer it can turn off the new-line as separation.
func pre_parse(sourcefile string) *tsource{
	doing("Analysing: ",sourcefile)
	bank,e:=qff.EGetFile(sourcefile)
	if e!=nil { ethrow(e) }
	if bank[0]=='"' { throw("Not a single file in Scyndi may start with a \"!") }
	ret:= &tsource{} 
	return ret
}
