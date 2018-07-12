package scynt

import (
		"trickyunits/qff"
jcr		"trickyunits/jcr6/jcr6main"
		"trickyunits/qstr"

		"strings"
		"fmt"
)

const parchat = false

func pchat(a... string){
	if !parchat { return }
	for _,s :=range a { 
		fmt.Println("= ",s)
	}
}

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
	doingln("Pre-parse analysing: ",file)
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
		if c=='{'  && !mlc && !ins { mlc=true }
		if c=='}'  && mlc { mlc=false; ok=false }
		if c=='"'  && !mlc && (!gbs) { ins=!ins }
		if gbs {
			/*
			nc:=c
			switch c {
				case 'n': nc=10
				case 'r': nc=13
				case 't': nc= 9
				case 'b': nc= 8
			}
			cl=append(cl,nc)
			*/
			cl=append(cl,'\\')
			cl=append(cl,c) 
			ok=false
			gbs=false
		} else if c=='\\' && ins {gbs=true; ok=false }		
		ok=ok && !mlc
		if ((c=='\n' && NLSEP) || c==';') && ok && !ins {
			//psource=append(psource,string(cl))
			psl:=qstr.MyTrim(string(cl))
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
	if gbs { throw("Unexpected end of file. No follow character to \\") }
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
					gbs:=false
					if i>0 { gbs=bline[i-1]=='\\' }
					if b=='"' && !gbs {
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
							if word!="" {
							   nw:=&tword{}
							   nw.Word  = word
							   nw.Wtype = gettype(word, file, so.ln)
							   so.sline = append(so.sline,nw)							
						    }
							nw:=&tword{}
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

func (self *tsource) declarechunk(ol *tori) *tchunk{
	
	ct:=ol.sline[0]
	args:=&targs{}
	locals:= map[string]*tidentifier{}
	
	
	if ct.Word=="BEGIN"{
		if len(ol.sline)>1 { ol.throw("BEGIN does not allow parameters of any sort") }
		if self.srctype=="PROGRAM" {
			ol.sline = []*tword{}
			ol.sline = append(ol.sline,&tword{"PROCEDURE","keyword"})
			ol.sline = append(ol.sline,&tword{"MAIN","identifier"})
		} else {
			ol.sline = []*tword{}
			ol.sline = append(ol.sline,&tword{"PROCEDURE","keyword"})
			ol.sline = append(ol.sline,&tword{"INIT","identifier"})
		}
	}
	wl:=ol.sline
	if len(wl)<2 { ol.throw("Expected identifier") }
	tp:=wl[0]
	id:=wl[1]
	idname:=id.Word
	if id.Wtype=="keyword" { ol.throw("Unexpected keyword ("+idname+")") }
	if _,foundid:=self.identifiers[idname];foundid { ol.throw("Duplicate identifier: "+idname) }	
	myname:="SCYNDI_PRG_"+self.srcname+"_"+tp.Word+"_"+strings.ToUpper(idname)
	mytype:="VARIANT"
	if len(wl)>2 {
		q:=2;
		qt:=wl[q]
		if qt.Word==":" {
			if len(wl)<q+1 { ol.throw("Type expected") }
			if ct.Word=="VOID" || ct.Word=="PROCEDURE" || ct.Word=="PROC" { ol.throw("Procedures have no return type") }
			if qt.Wtype=="keyword" { ol.throw("Unexpected keyword ("+qt.Word+")") }
			if tid,found:=self.identifiers[qt.Word];found {
				if tid.dttype!="TYPE" { ol.throw("Invalid identifier. Expected type but I got "+tid.dttype); }
			} else {
				if qt.Word!="STRING" && qt.Word!="INTEGER" && qt.Word!="BOOLEAN" && qt.Word!="FLOAT" && qt.Word!="VARIANT" {
					ol.throw("Unknown type: "+qt.Word)
				} else {
					mytype=qt.Word
				}
			}
			q+=2
		}
		// parameters for functions/procedures
		stage:=0 // 0 = new + identifier (or var), 1 = dubbele punt, 2 = type or ...,3 = '=' for optional arguments, 4 = default value (constant only)
		constant:=true
		endless:=false
		
		var arg *targ
		var aid *tidentifier
		for q<len(wl) {
			qt=wl[q]
			// TODO HERE: main parameter code (comes later)
			if qt.Word=="," && (stage==0 || stage==2 || stage==4) { ol.throw("Unexpected coma") } else if qt.Word=="," {
				stage=0
				constant=true
				if endless { ol.throw("Endless arguments always come last") }
			} else {
				switch stage {
					case 0: if qt.Word=="VAR" {
								if !constant { ol.throw("Double var") }
								constant=false
							} else if qt.Wtype!="identifier" {
								ol.throw("Unexpected "+qt.Word)
							}
							for _,a:=range args.a {
								if a.argname==qt.Word { ol.throw("Duplicate argument") }
							}
							arg=&targ{}
							args.a=append(args.a,arg)
							arg.argname=qt.Word
							arg.arg=&tidentifier{}; aid=arg.arg
							aid.idtype="VAR"
							aid.dttype="VARIANT"
							aid.translateto="SCYNDI_ARGUMENT_"+qt.Word
							aid.constant=constant
							locals[qt.Word]=aid
							stage=1
					case 1:
							if qt.Word!=":" { ol.throw(": expected") }
							stage=2
					case 2:
							if qt.Word!="STRING" && qt.Word!="INTEGER" && qt.Word!="BOOLEAN" && qt.Word!="FLOAT" && qt.Word!="VARIANT" {
								ol.throw("Unknown type: "+qt.Word)
							}
							aid.dttype=qt.Word
							stage=3
					case 3:
							if qt.Word!=":" { ol.throw("= expected") }
							if !constant { ol.throw("Variable arguments cannot be optional") }
							stage=4
							arg.optional=true
					case 4:
							ok:=false
							if (aid.dttype=="VARIANT" || aid.dttype=="STRING") && qt.Wtype=="string"                         { ok=true }
							if (aid.dttype=="VARIANT" || aid.dttype=="INTEGER" || aid.dttype=="FLOAT") && qt.Wtype=="string" { ok=true }
							if (aid.dttype=="VARIANT" || aid.dttype=="BOOLEAN") && (qt.Word=="TRUE" && qt.Word=="FALSE")     { ok=true }
							if !ok { ol.throw("Unexpected "+qt.Wtype+": "+qt.Word) }
							aid.defaultvalue=qt.Word
					default:
							ol.throw("Internal error. Invalid stage. Please report!")
				}
			}
			q++
		}
		
	}
	// TODO HERE: Create function code chunk
	// TODO HERE: Declare identifier for this function
	self.levels=append(self.levels,&tstatementspot{ol.ln,tp.Word})
	rc:=&tchunk{}
	if ct.Word=="VOID" || ct.Word=="PROCEDURE" || ct.Word=="PROC" || ct.Word=="BEGIN" { rc.pof=0; mytype="VOID" } else { rc.pof=1 }
	rc.instructions = []*tinstruction{}
	rc.locals =map[string]*tidentifier{}
	rc.args=args
	rc.translateto=myname
	rc.locals=locals
	rc.returntype=mytype
	rc.from=ol
	apof:=[]string{"PROCEDURE","FUNCTION"}
	cid:=&tidentifier{}
	cid.idtype=apof[rc.pof]
	cid.dttype=mytype
	cid.translateto=myname
	cid.args=args
	cid.constant=true
	self.chunks = append(self.chunks,rc)
	return rc
}

// Basically step #2 in compiling.
// Organising the code blocks
func (self *tsource) Organize(){
	var mychunk  *tchunk
	agroundkeys:=[]string{"BEGIN","VOID","PROCEDURE","PROC","FUNC","FUNCTION","DEF","VAR","TYPE","USE","XUSE","PRIVATE","PUBLIC","IMPORT"} // On "ground level" only these keywords are allowed.
	ogroundkeys:=[]string{"BEGIN","VOID","PROCEDURE","PROC","FUNC","FUNCTION","DEF","TYPE","USE","XUSE","PRIVATE","PUBLIC","IMPORT"}       // These keywords are ONLY allowed on "ground level".
	ltype:="ground"
	doing("Organising: ",self.filename)
	self.levels=[]*tstatementspot{}
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
				if pt.Wtype!="keyword" { pchat("Codeblock: "+ltype+"; "); ol.throw("Unexpected "+pt.Wtype) }
				if !contains(agroundkeys,pt.Word) { ol.throw("Unexpected "+pt.Word) }
				switch pt.Word {
					case "END": ol.throw("Unexpected END") 
					case "PRIVATE": self.private=true;  if len(sl)>1 { ol.throw("PRIVATE takes no parameters") }
					case "PUBLIC":  self.private=false; if len(sl)>1 { ol.throw("PUBLIC takes no parameters") }
					case "IMPORT":
						id,name:=self.performimport(ol)
						if _,ok:=self.identifiers[name];ok { ol.throw("Duplicate identifier by import request") }
						self.identifiers[name]=id
						pchat("External identifier imported as "+name)
					case "BEGIN","VOID","PROCEDURE","PROC","FUNCTION","FUNC","DEF":
						mychunk = self.declarechunk(ol)
						ltype="func"
					case "VAR":
						if len(sl)>1 { 
							tv:=sl[1:]
							/*
							n,i:=self.declarevar(tv)
							ol.pthrow(n)
							if _,found:=self.identifiers[n];found { ol.throw("Duplicate identifier: "+n) }
							self.identifiers[n]=&i
							*/
							dl:=&tori{}
							dl.sline=tv
							dl.pline=ol.pline
							dl.ln=ol.ln
							self.varblock = append(self.varblock,dl)
							pchat("Instant: VAR-block line added") 
						} else {
							ltype="var"
							self.levels=append(self.levels,&tstatementspot{ol.ln,"Global VAR declaration block"})
						}
					default:
						ol.throw("Unexpected "+pt.Word+"!! (Very likely a bug in the Scyndi compiler! Please report!)")
				}
			} else if ltype=="var" {
				if pt.Word=="END" {
					self.levels=self.levels[:len(self.levels)-1]
					pchat("VAR-block ended on line "+fmt.Sprintf("%5d",ol.ln))
					ltype="ground"
				} else {
					self.varblock = append(self.varblock,ol)
					pchat("VAR-block line added: "+pt.Word)
				}
			} else {
				// -- >> Inside a function or procedure
				if pt.Wtype=="keyword" && contains(ogroundkeys,pt.Word) {  ol.throw("Keyword "+pt.Word+" can only be used on the 'lowest' level of the program") }
				ins:=&tinstruction{}
				ins.ori=ol
				ins.level=len(self.levels)
				ins.state=self.levels[len(self.levels)-1]
				//fmt.Println(mychunk,ltype)
				mychunk.instructions = append(mychunk.instructions,ins)
				pchat("Instruction line added >> "+pt.Word)
				if pt.Word=="IF" || pt.Word=="WHILE" || pt.Word=="DO" || pt.Word=="REPEAT" {
					self.levels=append(self.levels,&tstatementspot{ol.ln,pt.Word+" block"})
				}
				if pt.Word=="LOOP" || pt.Word=="UNTIL" || pt.Word=="FOREVER" {
					bl:=self.levels[len(self.levels)-1]
					if bl.openinstruct!="REPEAT block" && bl.openinstruct!="DO block" {
						ol.throw(pt.Word+" cannot be used to end a "+bl.openinstruct+"! Use and END to end that in stead! "+pt.Word+" can only be used to end a DO/REPEAT block!")
					}
					self.levels=self.levels[:len(self.levels)-1]
					if len(self.levels)==0 {ltype="ground"} else {ltype="func"}
				}
				if pt.Word=="END" {
					self.levels=self.levels[:len(self.levels)-1]
					if len(self.levels)==0 {ltype="ground"} else {ltype="func"}
				}
			}
			//throw("Unfortunately, the part of the organisor to perform what is set up next has not yet been written")
		}
	}
	done()
}

func CompileFile(file string,t string) (string, *tsource) {
	doingln("Processing:",file)
	k:=Grabfromfile(file)
	s:=Sepsource(k,file)
	s.Organize()
	if s.srctype!=t && s.srctype!="SCRIPT" { throw("Source "+file+" may not be used for the purpose that it's been attempted to use\nWanted: "+t+"\nGot:    "+s.srctype) }
	return s.Translate(),s
}



func (self *tsource) Translate() string {
	doingln("Sorting out dependencies for translating: ",self.filename)
	trans:=TransMod[TARGET]
	//trans.NameIdentifiers(self)
	blocks:=map[string]string{}
	if trans.SealBlocks!=nil { trans.SealBlocks(&blocks) }
	blocks["USE"]=""
	useblock(TransMod,self,&blocks)
	doingln("Translating: ",self.filename)
	blocks["VAR"]=self.declarevars()
	blocks["FUN"]=self.translatefunctions()
	return trans.Merge(blocks)
}
