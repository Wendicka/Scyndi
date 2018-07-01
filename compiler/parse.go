package scynt

import (
		"trickyunits/qff"
jcr		"trickyunits/jcr6/jcr6main"
		"trickyunits/qstr"
		"strings"
)

func Grabfromfile(sourcefile string) *[]byte{
	doing("Reading file: ",sourcefile)
	bank,e:=qff.EGetFile(sourcefile)
	if e!=nil { ethrow(e) }
	if bank[0]=='"' { throw("Not a single file in Scyndi may start with a \"!") }
	done()
	return &bank
}

func Grabfromjcr(j jcr.TJCR6Dir,entry string) *[]byte{
	doing("Reading JCR entry: ",entry)
	bank:=jcr.JCR_B(j,entry)
	if jcr.JCR6Error!="" { throw(jcr.JCR6Error) }
	if bank[0]=='"' { throw("Not a single file in Scyndi may start with a \"!") }
	done()
	return &bank
}


func gettype(word string,file string,line int) string{
	t:="???"
	if len(word)==0{ return "" }
	switch word[0]{
		case '0','1','2','3','4','5','6','7','8','9','%','$':
			t = "integer"
			if strings.IndexByte(word,'.')>=0 { t="float" }
		case 'A','B','C','D','E','F','G','H','I','J','K','L','M','N','O','P','Q','R','S','T','U','V','W','X','Y','Z','_':
			t = "identifier"
			for _,w:=range keywords {
				if w==word { t = "keyword" }
			}
		default:
			lthrow(file,line,"Unknown series of characters: "+word)
	} 
	return t
}


// Basically all this routine does is make sure all instructions are properly separated.
// New lines (0x0A) and ; can do this, although those who prefer it can turn off the new-line as separation.
func Sepsource(src *[] byte,file string) *tsource {
	ret:=&tsource{}
	doing("Pre-parse analysing: ",file)
	ret.nlsep = NLSEP
	mlc:=false // Multi-line comment
	nwp:=true  // No whitespace
	ins:=false // in string
	gbs:=false // Got backslash?
	lnb:=1
	pc:=0
	ret.source=[]*tori{}
	cl:=[]byte{}
	for i:=0;i<len(*src);i++{
		pc++
		if pc>100 { progress(i+1,len(*src)) }
		c:=(*src)[i]
		ok:=true
		// get rid of useless whitespaces
		if (c==' ' || c=='\r' || c=='\t' || (c=='\n' && !NLSEP)) {
			if nwp {
				ok=false
			} else {
				if (!mlc) && (!ins) { nwp=true }
			}
		} else { nwp = false }
		// I have no need for the text in multi-line comments. Get rid of them
		if c=='{' && !mlc && !ins { mlc=true }
		if c=='}' && mlc { mlc=false; ok=false }
		if c=='"' && !mlc && (!gbs) { ins=!ins }
		ok=ok && !mlc
		if ((c=='\n' && NLSEP) || c==';') && ok && !ins {
			//psource=append(psource,string(cl))
			psl:=
			qstr.MyTrim(string(cl))
			cl=[]byte{}
			ok=false
			no:=&tori{}
			no.pline=psl
			no.ln=lnb
			if psl!="" {
				ret.source = append(ret.source,no)
			}
			nwp=true
			//ci:=len(ret.instructions)-1
			//ret.instructions[ci].ori = no
			if c=='\n' {lnb++}
		} else if c=='\n' { lnb++ }
		if ok {
			if (!ins) && c>='a' && c<='z' { c-=32 }
			cl=append(cl,c)
		}
	}
	if ins { throw("Unexpected end of file. String is not finished") }
	if mlc { throw("Unexpected end of file. Comment is not finished") }
	//chkstuff:=[][]string{keywords,operators}
	for _,so:=range ret.source {
		ig:=0
		so.sline = []*tword{}
		//for _,sc:=range ret.source{
			bline:=[]byte(so.pline)
			forcenw:=false
			//t:="?"
			word:=""
			instring:=false
			for i:=0;i<len(so.pline);i++{
				b:=bline[i]
				if ig>0 {
					ig--
				} else if instring {
					if b=='"' {
						instring=false
						nw:=&tword{}
						nw.Word  = word
						nw.Wtype = "string"
						so.sline = append(so.sline,nw)
						word=""
					} else {
						word = word + qstr.Mid(so.pline,i+1,1)
					}
				} else {
					if b==' ' || b=='\t' || b=='\r' || b=='\n' { forcenw=true }
					if b=='"' { forcenw=true; instring=true }
					for _,o:=range operators{
						if i+len(o)<=len(so.pline) && qstr.Mid(so.pline,i+1,len(o))==o && !forcenw {
							nw:=&tword{}
							nw.Word  = word
							nw.Wtype = gettype(word, file, so.ln)
							so.sline = append(so.sline,nw)
							nw=&tword{}
							nw.Word=o
							nw.Wtype="operator"
							so.sline = append(so.sline,nw)
							ig=len(o)-1
							word=""
							forcenw=true
						}
					}
					if forcenw {
						if word!="" {
							nw:=&tword{}
							nw.Word  = word
							nw.Wtype = gettype(word, file, so.ln)
							so.sline = append(so.sline,nw)
						}
						word=""
						forcenw=false
					} else {
						word=word+qstr.Mid(so.pline,i+1,1)
					}
				}
			}
		//}
		if word!="" {
			nw:=&tword{}
			nw.Word  = word
			nw.Wtype = gettype(word, file, so.ln)
			so.sline = append(so.sline,nw)
		}
			
	}
	return ret
}
