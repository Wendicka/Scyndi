/*
	Scyndi
	General Parsing
	
	
	
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
Version: 18.08.11
*/
package scynt

import (
		"trickyunits/qff"
jcr		"trickyunits/jcr6/jcr6main"
		"trickyunits/qstr"
		"trickyunits/mkl"
fpath	"path/filepath"
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
	bank=append(bank,';')
	bank=append(bank,10)
	done()
	return &bank
}

func Grabfromjcr(j jcr.TJCR6Dir,entry string) *[]byte{
	doing("Reading JCR entry: ",entry)
	bank:=jcr.JCR_B(j,entry)
	if jcr.JCR6Error!="" { throw(jcr.JCR6Error) }
	if bank[0]=='"' { throw("Not a single file in Scyndi may start with a \"!") }
	bank=append(bank,';')
	bank=append(bank,10)
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
				ok =ok || (word[i]=='.')
				lassert(file,line,ok,"Invalid identifier: "+word)
			}
		case '#':
			t = "Preprocessor tag"
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
	ret.allid=map[string]*tidentifier{}
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
	// make sure escape codes are properly seen
	for _,so:=range ret.source {
		for _,fw:=range so.sline{
			if fw.Wtype=="string" {
				for k,v:=range okki {
					//doingln("okki:",k)//debug
					fw.Word=strings.Replace(fw.Word,k,string([]byte{v}),-1)
				}
				fw.Word=strings.Replace(fw.Word,"\\\\","\\",-1) // Must always be done last.
			}
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
		//doingln(self.srctype+" ",self.filename) // debug
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
			qt:=wl[q+1]
			if len(wl)<q+1 { ol.throw("Type expected") }
			if ct.Word=="VOID" || ct.Word=="PROCEDURE" || ct.Word=="PROC" { ol.throw("Procedures have no return type") }
			if tid,found:=self.identifiers[qt.Word];found {
				if tid.dttype!="TYPE" { ol.throw("Invalid identifier. Expected type but I got "+tid.dttype); }
			} else {
				if qt.Word!="STRING" && qt.Word!="INTEGER" && qt.Word!="BOOLEAN" && qt.Word!="FLOAT" && qt.Word!="VARIANT" {
					if qt.Wtype=="keyword" { ol.throw("Unexpected keyword ("+qt.Word+")") }
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
			if qt.Word=="," && (stage==0 || stage==2 || stage==4) { ol.throw(fmt.Sprintf("Unexpected comma   (code: %d)",stage)) } else if qt.Word=="," {
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
							if qt.Word=="..." {
								endless=true
								//aid.endless=true
								args.a=args.a[:len(args.a)-1]
								args.endless=arg
							} else if qt.Word!="STRING" && qt.Word!="INTEGER" && qt.Word!="BOOLEAN" && qt.Word!="FLOAT" && qt.Word!="VARIANT" {
								ol.throw("Unknown type: "+qt.Word)
								aid.dttype=qt.Word
								stage=3
							} else {
								aid.dttype=qt.Word
								stage=3
							}
					case 3:
							if qt.Word!="=" { ol.throw("= expected") }
							if !constant { ol.throw("Variable arguments cannot be optional") }
							if endless { ol.throw("Default values not allowed in endless argument declaration") }
							stage=4
							arg.optional=true
					case 4:
							ok:=false
							if (aid.dttype=="VARIANT" || aid.dttype=="STRING") && qt.Wtype=="string"                                                 { ok=true }
							if (aid.dttype=="VARIANT" || aid.dttype=="INTEGER" || aid.dttype=="FLOAT") && (qt.Wtype=="integer" || qt.Wtype=="float") { ok=true }
							if (aid.dttype=="VARIANT" || aid.dttype=="BOOLEAN") && (qt.Word=="TRUE" || qt.Word=="FALSE")                             { ok=true }
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
	self.levels=append(self.levels,&tstatementspot{ol.ln,tp.Word,0})
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
	self.identifiers[id.Word]=cid
	self.chunks = append(self.chunks,rc)
	return rc
}

func (self *tsource) doinclude(ol *tori, file string,isprogram bool) (dump []*tori,dumpid map[string]*tidentifier){
	workfile:=""
	if (qstr.Prefixed(file,"/") || file[1]==':'){
		if qff.IsFile(file) { workfile=file }
	} else {
		//doingln("f:",self.filename)
		paths:=[]string{ fpath.Dir(self.filename) }
		for _,spath:=range paths {
			path:=strings.Replace(spath,"\\","/",-1)
			if !qstr.Suffixed(path,"/") { path+="/" }
			//doingln("#INCLUDE: Trying: ",path+file)
			if qff.IsFile(path+file) { workfile=path+file; break }
		}
	}
	if workfile=="" { ol.throw("I could not find \""+file+"\" for #INCLUDE request") }
	hey()
	doingln("Including: ",workfile)
	tempbank:=Grabfromfile(workfile)
	tempsrc:=Sepsource(tempbank,workfile)
	tempsrc.srctype="SCRIPT"
	if isprogram {tempsrc.srctype="PROGRAM"}
	tempsrc.noheader=true
	tempsrc.identifiers=map[string]*tidentifier{}
	tempsrc.Organize()
	dump=tempsrc.source
	doing("Resuming organisation of: ",self.filename)
	return
}

// Basically step #2 in compiling.
// Organising the code blocks
func (self *tsource) Organize(){
	var mychunk  *tchunk
	agroundkeys:=[]string{"BEGIN","VOID","PROCEDURE","PROC","FUNC","FUNCTION","DEF","VAR","TYPE","USE","XUSE","PRIVATE","PUBLIC","IMPORT","CONST","ENUM"} // On "ground level" only these keywords are allowed.
	ogroundkeys:=[]string{"BEGIN","VOID","PROCEDURE","PROC","FUNC","FUNCTION","DEF","TYPE","USE","XUSE","PRIVATE","PUBLIC","IMPORT","CONST","ENUM"}       // These keywords are ONLY allowed on "ground level".
	ltype:="ground"
	doing("Organising: ",self.filename)
	self.levels=[]*tstatementspot{}
	headerset:=self.noheader
	ppb:=false // Pre-Process block
	localdefs:=map[string] bool {}
	// Presearch for includes
	subsource:=[]*tori{}
	isprogram:=false
	for _,ol:=range self.source {
		if ol.sline[0].Word=="PROGRAM" {isprogram=true}
		if ol.sline[0].Word=="#INCLUDE" {
			if len(ol.sline)!=2 { ol.throw("Invalid #INCLUDE request") }
			if ol.sline[1].Wtype!="string" { ol.throw("Constant string expected for #INCLUDE request") }
			incori,incid:=self.doinclude(ol,ol.sline[1].Word,isprogram)
			for _,iol:=range incori {
				subsource=append(subsource,iol)
			}
			for ikey,iid:=range incid{
				self.identifiers[ikey]=iid
			}
		} else {
			subsource=append(subsource,ol)
		}
	}
	self.source=subsource
	// Default defines
	for k,_:=range TransMod { globaldefs["$TARGET_"+strings.ToUpper(k)]=k==TARGET }
	// Let's go
	for _,ol:=range self.source {
		//doingln("let's parse: ",ol.sline[0].Word)
		if qstr.Prefixed(ol.sline[0].Word,"#") {
			switch ol.sline[0].Word {
				case "#DEFINE","#UNDEF":
					if ppb { break }
					for i:=1;i<len(ol.sline);i++{
						mydef:=ol.sline[i].Word
						switch mydef[0]{
							case '$':	ol.throw("$ prefix is reserved for system based definitions and may not be used in a #DEFINE or #UNDEF statement so "+mydef+" is invalid!")
							case '#':	localdefs[mydef]=ol.sline[0].Word=="DEFINE"
							default:	globaldefs[mydef]=ol.sline[0].Word=="DEFINE"
						}
					}
				case "#IF","#IFNOT":
					ppb=true
					//doing("#","IF")
					for i:=1;i<len(ol.sline);i++{
						mydef:=ol.sline[i].Word						
						switch mydef[0]{
							case '#':	if _,ok:=localdefs[mydef];!ok{localdefs[mydef]=false}
										ppb = ppb && localdefs[mydef]
							default:	if _,ok:=globaldefs[mydef];!ok{globaldefs[mydef]=false}
										ppb = ppb && globaldefs[mydef]
						}
						//fmt.Println("#IF",i,ol.sline[i].Word,ppb)
					}
					
					if ol.sline[0].Word=="#IF" {ppb=!ppb}
				case "#ENDIF","#FI":
					ppb=false
				case "#ELSE":
					ppb=!ppb
				case "#ERROR":
					if len(ol.sline)<=1 { ol.throw("#ERROR requires a message to throw!") }
					if !ppb { 
						cerr:=""
						for i:=1;i<len(ol.sline);i++ { cerr += ol.sline[i].Word + " "}
						ol.throw("Custom Error: "+cerr)
					}
				default: ol.throw("Unknown preprocessor definition: "+ol.sline[0].Word)
			}
		} else if (!headerset) && (!ppb) {
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
			headerset=true
		} else if !ppb {
			sl:=ol.sline
			pt:=sl[0]
			if ltype=="ground" {
				if pt.Wtype!="keyword" { pchat("Codeblock: "+ltype+"; "); ol.throw("Unexpected "+pt.Wtype) }
				if !contains(agroundkeys,pt.Word) { ol.throw("Unexpected "+pt.Word) }
				switch pt.Word {
					case "END": ol.throw("Unexpected END") 
					case "PRIVATE": self.private=true;  if len(sl)>1 { ol.throw("PRIVATE takes no parameters") }
					case "PUBLIC":  self.private=false; if len(sl)>1 { ol.throw("PUBLIC takes no parameters") }
					case "USE":
						if len(sl)>=2 { 
							for i:=1;i<len(sl);i++ {
								switch sl[i].Wtype{
									case "string": self.userequested=append(self.userequested,sl[i].Word) 
									case "identifier": self.userequested=append(self.userequested,strings.ToLower(sl[i].Word) )
								}
								//doing("REQUESTED FOR USE:",sl[i].Word)
							}
						} else if len(sl)==1 { 
							ltype="use"
						} else {
							ol.throw("Misunderstood USE request")
						}
					case "CONST":
						if len(sl)!=1 { self.defconst(sl[1:]) } else { 
							ltype="const" 
							self.levels=append(self.levels,&tstatementspot{ol.ln,"Constant definition block",0})
						}
					case "ENUM":
						self.enumstart(ol)
						ltype="enum"
						self.levels=append(self.levels,&tstatementspot{ol.ln,"ENUM block",0})
					case "IMPORT":
						id,name:=self.performimport(ol)
						if _,ok:=self.identifiers[name];ok { ol.throw("Duplicate identifier by import request") }
						self.identifiers[name]=id
						pchat("External identifier imported as "+name)
					case "BEGIN","VOID","PROCEDURE","PROC","FUNCTION","FUNC","DEF":
						mychunk = self.declarechunk(ol)
						ltype="func"
					case "TYPE":
						if self.starttype(ol) {
							ltype="type"
						}
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
							self.levels=append(self.levels,&tstatementspot{ol.ln,"Global VAR declaration block",0})
						}
					default:
						ol.throw("Unexpected "+pt.Word+"!!") // (Very likely a bug in the Scyndi compiler! Please report!)")
				}				
			} else if ltype=="use" {
				for _,w:=range sl {
					if w.Word=="END" { 
						ltype="ground"
					} else if w.Wtype=="string" || w.Wtype=="identifier" {
						self.userequested=append(self.userequested,w.Word) 
					}
				}
			} else if ltype=="type" {
				if pt.Word=="END" {
					ltype="ground"
				} else {
					self.totype(ol)
				}
			} else if ltype=="enum" {
				if pt.Word=="END" {
					self.levels=self.levels[:len(self.levels)-1]
					ltype="ground"
				} else {
					self.enumid(ol)
				}
			} else if ltype=="const" {
				if pt.Word=="END" {
					self.levels=self.levels[:len(self.levels)-1]
					pchat("CONST-block ended on line "+fmt.Sprintf("%5d",ol.ln))
					ltype="ground"
				} else {
					self.defconst(sl)
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
				if pt.Word=="IF" || pt.Word=="WHILE" || pt.Word=="DO" || pt.Word=="REPEAT" || pt.Word=="FOR" || pt.Word=="FORU" || pt.Word=="FOREACH" || pt.Word=="SELECT" || pt.Word=="SWITCH" {
					self.levels=append(self.levels,&tstatementspot{ol.ln,pt.Word+" block",0})
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
	s.noheader=false
	s.Organize()
	if s.srctype!=t && s.srctype!="SCRIPT" { throw("Source "+file+" may not be used for the purpose that it's been attempted to use\nWanted: "+t+"\nGot:    "+s.srctype) }
	return s.Translate(),s
}



func (self *tsource) Translate() string {
	doingln("Sorting out dependencies for translating: ",self.filename)
	trans,present:=TransMod[TARGET]
	if !present { throw("Target "+TARGET+" is not supported") }
	//trans.NameIdentifiers(self)
	blocks:=map[string]string{}
	if trans.SealBlocks!=nil { 
		trans.SealBlocks(&blocks) 
	}
	blocks["USE"]=""
	useblock(TransMod,self,&blocks)
	blocks["TYPES"]=trans.TransTypes(self)
	doingln("Translating: ",self.filename)
	blocks["VAR"]=self.declarevars()
	blocks["FUN"]=self.translatefunctions()
	return trans.Merge(blocks)
}

func (self *tsource) SaveTranslation(strans,outputpath string) {
	doingln("Saving: ","Translation")
	trans:=TransMod[TARGET]
	trans.savetrans(self,strans,outputpath)
}

func init(){
mkl.Lic    ("Scyndi Programming Language - parse.go","GNU General Public License 3")
mkl.Version("Scyndi Programming Language - parse.go","18.08.11")
}
