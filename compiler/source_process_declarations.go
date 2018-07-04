package scynt


type tidentifier struct {
	private bool
	idtype string // function, procedure, type, constant, variable, sourcegroup(either program, module, script)
	dttype string // data type string, int etc. (for variables and functions returning values)
	defaultvalue string
	translateto string // as some legal names in Scyndi can be keywords in the target language, Scyndi will use different names in its translations.
	tarformed bool // Will be true if the translation module already reformed this variable, in order to rule out ANY POSSIBILITY AT ALL it will happen twice.
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

func (s *tori) LORI() (int,string,[]*tword) { return s.ln,s.pline,s.sline } // debug

type tstatementspot struct {
	openline int
	openinstruct string
}


type tinstruction struct {
	ori *tori
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
}


type tchunk struct {
	ismethod string
	pof byte	// 0 = procedure, 1 = function 
				// for translating to Wendicka or a scripting language such as php or even lua, this may not matter, but when translating to languages 
				// like Pascal, C or Go, this information can be crucial (especially in Go where the compiler is very very strict on these matters).
	instructions [] tinstruction
	locals []*tidentifier

}



type tsource struct {
	srctype string // May contain either "PROGRAM", "SCRIPT" or "MODULE"
	srcname string
	inputname string
	filename string
	chunks [] tchunk
	identifiers map[string]*tidentifier
	source []*tori
	// orilinerem will place the original line in the translation as a comment or remark
	// write traceback will instruct the parser of the translated code to process the traceback data, providing the target language has any way to support such a thing.
	// nlsep will if turned on (default value) accept a new line as a separator (and then you don't need a semi-colon at the end of each line), turning it off will require such a thing. Please note when putting multiple instructions on one line, the semi-colon will always be required to separate those.
	orilinerem,writetraceback,nlsep bool
	private bool
	levels []tstatementspot
	target string
	spackage *TPackage
}

func (s *tsource) Lsource() []*tori { return s.source }




type TPackage struct {
	sources [] *tsource
	mainsource *tsource
	outputf string
	translateto string
}

