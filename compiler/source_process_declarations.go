package scynt

type targ struct {
	argtype string // Just the types... Scyndi will use this to translate function calls.. 
	
	// These two fields may play a role when translating functions
	argname string 
	arg *tidentifier
	optional bool
	
}
type targs struct {
	a []*targ
	endless *targ
}

type tidentifier struct {
	imported bool
	private bool
	idtype string // function, procedure, type, constant, variable, sourcegroup(either program, module, script)
	dttype string // data type string, int etc. (for variables and functions returning values)
	defaultvalue string
	defstring bool // When a string, it must be made sure, as strings often have quotations and other required stuff.
	translateto string // as some legal names in Scyndi can be keywords in the target language, Scyndi will use different names in its translations.
	tarformed bool // Will be true if the translation module already reformed this variable, in order to rule out ANY POSSIBILITY AT ALL it will happen twice.
	args *targs
	constant bool // When set true the code cannot change this identifier after being defined.
}

type texpression struct{
	pure string
}

type tword struct{
	Word string
	Wtype string
}

type tori struct{
	pline string
	sline []*tword
	sfile string
	ln int
}

func (o *tori) getword(i int) *tword{
	if i>=len(o.sline) { o.throw("Syntax error! Not all data needed provided!") }
	return o.sline[i]
}

func (s *tori) LORI() (int,string,[]*tword) { return s.ln,s.pline,s.sline } // debug

type tstatementspot struct {
	openline int
	openinstruct string
	startfor int
}


type tinstruction struct {
	ori *tori
	parsed bool // All instructions only come in with the ori value set. Only keywords starting extra sub-chunks will have the extra level set, but other than that, no. This is to make sure everything can be declared in random order in the source code Scyndi processes.
	instruct byte // 0 = call function/procedure
				  // 1 = define variable
				  // 2 = start "if"
				  // 3 = start "elseif"
				  // 4 = start "else"
				  // 6 = start "while"
				  // 7 = start "four"
				  // 8 = start "repeat"
				  // 8 = end
				  // 9 = until
	expressions [] texpression
	id string // used for "end" en "until" to be properly tied to their respective start instructions. 
			  // for exporting to languages using { and } it may not matter too much (except for repeat/until statements) same goes for languages just using "end"
			  // But if exporting to languages based on BASIC and COMAL for example, it can hurt pretty bad if the translator doesn't know all this :P
			  // And besides when parse checing Scyndi code, knowing all this would be better anyway :P

	level int
		// This basically determines the statement level. Comments like "if" and such make it go up one level, and "end" makes it go down a level.
		// The info stored in the array is only for the error/warning handler in order to refer to line numbers were statements began.
	state *tstatementspot
}


type tchunk struct {
	isimported bool
	ismethod string
	translateto string
	pof byte	// 0 = procedure, 1 = function 
				// for translating to Wendicka or a scripting language such as php or even lua, this may not matter, but when translating to languages 
				// like Pascal, C or Go, this information can be crucial (especially in Go where the compiler is very very strict on these matters).
	instructions [] *tinstruction
	locals map[string]*tidentifier
	args *targs
	returntype string
	from *tori
	forid map[string]*tidentifier
	fors map[int] bool

}

type smap struct {
	m map[string]*tsource
}


type tsource struct {
	srctype string // May contain either "PROGRAM", "SCRIPT" or "MODULE"
	srcname string
	inputname string
	filename string
	chunks [] *tchunk
	currentchunk *tchunk
	identifiers map[string]*tidentifier
	source []*tori
	varblock []*tori
	// orilinerem will place the original line in the translation as a comment or remark
	// write traceback will instruct the parser of the translated code to process the traceback data, providing the target language has any way to support such a thing.
	// nlsep will if turned on (default value) accept a new line as a separator (and then you don't need a semi-colon at the end of each line), turning it off will require such a thing. Please note when putting multiple instructions on one line, the semi-colon will always be required to separate those.
	orilinerem,writetraceback,nlsep bool
	private bool
	levels []*tstatementspot
	target string
	//spackage *TPackage
	used []*tsource
	usedmap *smap
	userequested []string
	usepurecode string
	allid map[string]*tidentifier // All own identifiers plus the imported ones.
}

func (s *tsource) GetIdentifier(name string,c *tchunk, o *tori) *tidentifier {
	var ret *tidentifier
	if c!=nil {
		loc:=c.locals
		if v,ok:=loc[name]; ok { return v }
	}
	if v,ok:=s.identifiers[name];ok { return v }
	if v,ok:=s.allid[name];ok { return v }
	// for k,_ := range(s.allid) { doing("ID: ",k) } // debug only
	if ret==nil {
		if o==nil { 
			throw("Unknown identifier: "+name)
		} else {
			o.throw("Unknown identifier: "+name)
		}
	}
	return ret
}

func (s *tsource) Lsource() []*tori { return s.source }



/*
type TPackage struct {
	sources [] *tsource
	mainsource *tsource
	outputf string
	translateto string
}
*/
