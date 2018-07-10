package scynt

/*
 * Now it should be clear that the "IMPORT" keyword does NOT do what it 
 * does in for example BlitzMax and Go, where it imports modules or 
 * libraries.
 * 
 * In Scyndi it imports identifiers DIRECTLY into the Scyndi language.
 * It should be noted that this is not an entirely risk-free procedure 
 * as Scyndi is not able to check if things are set up in the way the 
 * target language expects it, so if you screw up here, the 
 * compiler/parser of the target language will protest and not 
 * Scyndi, and since the target compiler/target cannot know it's 
 * translated code it's dealing with, debugging can form quite a 
 * challenge, so only use this when YOU KNOW WHAT YOU ARE DOING.
 * 
 * Now it works like this, and for this I use a kind of pseudo C++ as
 * target language, if you have the code
 * 
 *     void C_MyPrint(string mystring) {
 *          cout << mystring;
 *     }
 * 
 * You can import it into Scyndi like this:
 *     Import C_MyPrint Procedure MyPrint String
 *
 * And now the Scyndi line: 
 * 		MyPrint "Hello World"
 * Would be translated as:
 * 		C_MyPrint("Hello World");
 *
 * Also note that the IMPORT routine itself does not detect if the 
 * target is right.
 * If you have a Scyndi procedure in Pascal just a standard command, 
 * but not in C, you'll have to use your compiler directives to make 
 * sure things do get set up in order (that is if you want your code
 * to export to both Pascal and C without any modifications).
 * 
 * Overall, I also recommend to only use Imports in modules and not in
 * your main program.
 * 
 * 
 * Also special treatment is not always required.
 * Import MyVar Var MyVar:String
 * Would if you translate to php just result into $MyVar being imported
 * as php required. ;)
 * 
 * Also note the quotations are for the target identifier not requires, 
 * BUT if the target language is case sensitive highly recommended,
 * if you don't want the target parser to throw needless errors :P
 * For case INSENSITVE targets this doesn't matter.
 * 
 */
 
  
// This function will not actually generate a true translation.
// All this does is create an identifer the transalator will directly link
// to its target language counterpart... ;)  
func (s *tsource) performimport(ol *tori){
	imptar:=ol.getword(1).Word;
	imptype:=ol.getword(2).Word;
	isprocedure=imptype=="PROCEDURE" || imptype=="PROC" || imptype=="VOID"
	if (!isprocedure) && imptype!="VAR" && imptype!="FUNCTION" && imptype!="FUNC" && imptype!="DEF" { ol.throw("Syntax error. I do not understand the word "+imptype+" in this import instruction") }
	impid:=ol.getword(3);
	if impid.Wtype!="identifier" { ol.throw("Unexpected "+impid.Wtype)=". I expected an identifier to tie an imported identifier to") }
	dubbelepunt:=ol.getword(4).Word;
	qw:=5;
	if !isprocedure {
		if ol.getword(5).Word!=":" { ol.throw("':' expected") }
		dtpe:=ol.getword(6)
		if dpte.Word=="VARIANT" { ol.throw("VARIANT not allowed for imported functions") }
		if dpte.Word=="INTEGER" || dpte.Word=="STRING" || dpte.Word=="FLOAT" || dpte.Word=="BOOLEAN" || dpte.Wtype=="identifier" {
		} else { ol.throw("Identifier expected for imported "+impype) }
		qw+=2
	}
	if len(ol.sline)>5 && imptype=="VAR" { ol.throw("Variables do not accept parameters") } // Procedure type variables will (for now) not be importable. Perhaps in the future...
	
}
