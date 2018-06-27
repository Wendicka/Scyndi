package scynt

import (
		"trickyunits/qff"
jcr		"trickyunits/jcr6/jcr6main"
)

func grabfromfile(sourcefile string) *[]byte{
	doing("Reading file: ",sourcefile)
	bank,e:=qff.EGetFile(sourcefile)
	if e!=nil { ethrow(e) }
	if bank[0]=='"' { throw("Not a single file in Scyndi may start with a \"!") }
	done()
	return &bank
}

func grabfromjcr(j jcr.TJCRDir,entry string) *[]byte{
	doing("Reading JCR entry: ",sourcefile)
	bank:=jcr.JCR_B(j,entry)
	if jcr.JCRERROR!="" { throw(jcr.JCRERROR) }
	if bank[0]=='"' { throw("Not a single file in Scyndi may start with a \"!") }
	done()
	return &bank
}


// Basically all this routine does is make sure all instructions are properly separated.
// New lines (0x0A) and ; can do this, although those who prefer it can turn off the new-line as separation.
func sepsource(*[] byte,file string) *tsource {
	ret:=&tsource{}
	doing("Pre-parse analysing: ",file)
	ret.nlsep = NLSEP
	return ret
}
