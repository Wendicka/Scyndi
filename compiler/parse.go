package scynt

import (
		"trickyunits/qff"
jcr		"trickyunits/jcr6/jcr6main"
		"trickyunits/qstr"
		"strings"
		"fmt"
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
			for i:=0;i<len(word);i++ {
				ok:=word[i]=='_'
				ok =ok || (word[i]>='A' && word[i]<='Z')
				ok =ok || (word[i]>='0' && word[i]<='9')
				lassert(file,line,ok,"Invalid identifier: "+word)
			}
		default:
			lthrow(file,line,"Unknown series of characters: "+word)
	} 
	return t
}


// After loading the code this is step #1
// Basically all this routine does is make sure all instructions are properly separated.
// New lines (0x0A) and ; can do this, although those who prefer it can turn off the new-line as separation.
func Sepsource(src *[] byte,file string) *tsource {
	ret:=&tsource{}
	doing("Pre-parse analysing: ",file)
	ret.target=TARGET
	ret.filename=file
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
			no.sfile=file
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

func contains(a []string,s string) bool{
	for _,v :=range a { if v==s { return true }}
	return false
}

func (s *tsource) validtype(n *tword) bool{
	r:=n.Wtype=="identifier" || n.Word=="VARIANT" || n.Word=="STRING" || n.Word=="INTEGER" || n.Word=="FLOAT" || n.Word=="BOOLEAN" || n.Word=="ARRAY" || n.Word=="MAP"
	return r
}

func (s *tsource) declarevar(line []*tword) (string,tidentifier){
	//vr:=self.identifiers
	vr:=tidentifier{}
	vr.private=s.private
	vr.dttype="VARIANT"
	vr.idtype="VAR"
	vr.defaultvalue="NIL"
	if len(line)==0 {return "er:Empty variable declaration",vr } // should normally never happen, but at least this vovr a go panic crash
	n:=line[0]
	name:=n.Word
	vr.translateto="SCYNDI_VAR_"+s.srcname+"_"+name
	if contains(keywords,name) { return "er:"+name+" is a keyword and may NOT be used as a variable",vr }
	if n.Wtype!="identifier" { return "er:Unexpected "+n.Wtype+"("+name+"). A name for a variable was expected",vr }
	if len(line)==1 {return name,vr}
	i:=1
	o:=line[i]
	if o.Word==":" {
		if len(line)<3 { return "er:Unexpected end of line. A type for a variable was expected",vr }
		n:=line[i+1]
		if !s.validtype(n) { return "er:Invalid variable type. Either an unknown type or invalud type: "+n.Word,vr }
		vr.dttype=n.Word
		i+=2
	}
	if len(line)>i {
		o=line[i]
		if o.Word=="=" {
			i++
			if len(line)<i+1 { return "er:Unexpected end of line",vr }
			switch vr.idtype {
				case "VARIANT":
					return "er:VARIANTS cannot be defined in a variable block",vr
				case "MAP","ARRAY":
					return "er:MAPS and ARRAYS cannot be defined in a variable block (Hey, psst! They are basically already defined upon declaration, so you don't have to)",vr
				case "STRING":
					if o.Wtype!="string" { return "er:Unexpected "+o.Wtype+". Constant string required",vr }
					vr.defaultvalue = o.Word
				case "INTEGER":
					if o.Wtype!="integer" { return "er:Unexpected "+o.Wtype+". Constant integer required",vr }
					vr.defaultvalue = o.Word
				case "FLOAT":
					if o.Wtype!="integer" && o.Wtype!="float" { return "er:Unexpected "+o.Wtype+". Constant integer or float required",vr }
					vr.defaultvalue = o.Word
				case "BOOLEAN":
					if o.Wtype!="keyword" {return "er:Unexpected "+o.Wtype+". TRUE or FALSE required",vr }
					if o.Word!="TRUE" && o.Word!="FALSE" { return "er:Unexpected "+o.Word+". TRUE or FALSE required!",vr }
					vr.defaultvalue = o.Word
				default:
					if o.Wtype!="keyword" || o.Word!="NEW" { return "er:Unexpected "+o.Wtype+". Only the keyword NEW is allowed here",vr }
					vr.defaultvalue = "NEW"
			}
		} else { return "er:Syntax error!",vr } // Now it's really beyond me what you were trying to do.... :-/
	} else {
		switch vr.dttype {			
			case "STRING": vr.defaultvalue = ""
			case "INTEGER","FLOAT": vr.defaultvalue="0"
			case "BOOLEAN":  vr.defaultvalue="FALSE"
		}
	}
	return name,vr
}

// Basically step #2 in compiling.
// Organising the code blocks
func (self *tsource) Organize(){
	agroundkeys:=[]string{"VOID","PROCEDURE","PROC","FUNC","FUNCTION","DEF","VAR","TYPE","USE","XUSE","PRIVATE","PUBLIC"} // On "ground level" only these keywords are allowed.
	ogroundkeys:=[]string{"VOID","PROCEDURE","PROC","FUNC","FUNCTION","DEF","TYPE","USE","XUSE","PRIVATE","PUBLIC"}       // These keywords are ONLY allowed on "ground level".
	ltype:="ground"
	doing("Organising: ",self.filename)
	self.levels=[]tstatementspot{}
	for _,ol:=range self.source {
		if ol.ln==1 {
			lassert(ol.sfile,ol.ln,len(ol.sline)==2,"Illegal source header!  "+fmt.Sprintf("(%d)",len(ol.sline)))
			sl:=ol.sline
			pt:=sl[0]
			id:=sl[1]
			if pt.Word=="UNIT" { pt.Word="MODULE" }
			lassert(ol.sfile,ol.ln,pt.Word=="PROGRAM" || pt.Word=="MODULE" || pt.Word=="SCRIPT","Header must define PROGRAM, MODULE or SCRIPT. Not a "+pt.Word)
			if id.Wtype!="identifier" { lthrow(ol.sfile,ol.ln,"Unexpected "+id.Wtype+"! Expected identifier for "+pt.Word+" in stead!") }
			self.identifiers = map[string]*tidentifier{}
			self.identifiers[id.Word]=&tidentifier{}
			did:=self.identifiers[id.Word]
			did.private=false
			did.idtype="SOURCEGROUP"
			did.translateto="SCYNDI_SOURCEGROUP_"+id.Word
			self.srctype=pt.Word
			self.srcname=id.Word
		} else {
			sl:=ol.sline
			pt:=sl[0]
			if ltype=="ground" {
				if pt.Wtype!="keyword" { ol.throw("Unexpected "+pt.Wtype) }
				if !contains(agroundkeys,pt.Word) { ol.throw("Unexpected "+pt.Word) }
				switch pt.Word {
					case "END": ol.throw("Unexpected END") 
					case "PRIVATE": self.private=true;  if len(sl)>1 { ol.throw("PRIVATE takes no parameters") }
					case "PUBLIC":  self.private=false; if len(sl)>1 { ol.throw("PUBLIC takes no parameters") }
					case "VAR":
						if len(sl)>1 { 
							tv:=sl[1:]
							n,i:=self.declarevar(tv)
							ol.pthrow(n)
							self.identifiers[n]=&i
						} else {
							ltype="var"
							self.levels=append(self.levels,tstatementspot{ol.ln,"Global VAR declaration block"})
						}
					default:
						ol.throw("Unexpected "+pt.Word+"!! (Very likely a bug in the Scyndi compiler! Please report!)")
				}
			} else {
				if pt.Wtype=="keyword" && contains(ogroundkeys,pt.Word) {  ol.throw("Keyword "+pt.Word+" can only be used on the 'lowest' level of the program") }
				if pt.Word=="END" {
					self.levels=self.levels[:len(self.levels)-1]
					if len(self.levels)==0 {ltype="ground"} else {ltype="func"}
				}
			}
			//throw("Unfortunately, the part of the organisor to perform what is set up next has not yet been written")
		}
	}
	
}
