/*
	Scyndi
	Import processing and parsing
	
	
	
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
Version: 18.07.21
*/
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

//import "fmt"

func (s *tsource) performimport(ol *tori) (*tidentifier,string){
	tar:=&tidentifier{}
	tar.imported=true
	tar.args=&targs{}
	tar.private=s.private
	imptar:=ol.getword(1).Word;
	imptype:=ol.getword(2).Word;
	isprocedure:=imptype=="PROCEDURE" || imptype=="PROC" || imptype=="VOID"
	if (!isprocedure) && imptype!="VAR" && imptype!="FUNCTION" && imptype!="FUNC" && imptype!="DEF" { ol.throw("Syntax error. I do not understand the word "+imptype+" in this import instruction") }
	impid:=ol.getword(3);
	if impid.Wtype!="identifier" { ol.throw("Unexpected "+impid.Wtype+". I expected an identifier to tie an imported identifier to") }
	dubbelepunt:=ol.getword(4).Word;
	qw:=5;
	if !isprocedure {
		if dubbelepunt!=":" { ol.throw("':' expected") }
		dpte:=ol.getword(5)
		if dpte.Word=="VARIANT" { ol.throw("VARIANT not allowed for imported identifiers") }
		if dpte.Word=="INTEGER" || dpte.Word=="STRING" || dpte.Word=="FLOAT" || dpte.Word=="BOOLEAN" || dpte.Wtype=="identifier" {
			tar.dttype=dpte.Word
		} else { ol.throw("Identifier expected for imported "+imptype) }
		qw+=2
	} else { qw=4 }
	//fmt.Println("Start",qw,len(ol.sline))
	if len(ol.sline)>4 {
		if imptype=="VAR" { 
			ol.throw("Variables do not accept parameters")  // Procedure type variables will (for now) not be importable. Perhaps in the future...
		} else {
			wantcomma:=false
			endless:=false
			for i:=qw;i<len(ol.sline);i++{				
				//fmt.Println("Walk",qw,i,len(ol.sline))
				if wantcomma {
					if ol.getword(i).Word!="," { ol.throw("Comma expected") }
					wantcomma=false
				} else if ol.getword(i).Word=="..." {
					//fmt.Print(i,len(ol.sline))
					if i+1>=len(ol.sline) { ol.throw("Type wanted for infinite argument") }
					if i+3<=len(ol.sline) { ol.throw("Endless arguments always come last") }
					dpte:=ol.getword(i+1)
					if dpte.Word=="INTEGER" || dpte.Word=="STRING" || dpte.Word=="FLOAT" || dpte.Word=="BOOLEAN" || dpte.Wtype=="identifier" {
						arn:=dpte.Word
						ara:=&targ{}
						ara.argtype=arn
						tar.args.endless=ara
						endless=true
					} else {
						{ ol.throw("Type identifier expected for endless argument for imported "+imptype) }
					}
				} else if !endless {
					dpte:=ol.getword(i)
					if dpte.Word=="INTEGER" || dpte.Word=="STRING" || dpte.Word=="FLOAT" || dpte.Word=="BOOLEAN" || dpte.Wtype=="identifier" {
						arn:=dpte.Word
						ara:=&targ{}
						ara.argtype=arn
						ara.arg=&tidentifier{}
						ara.arg.dttype=dpte.Word
						tar.args.a = append(tar.args.a,ara)
						wantcomma=true
					} else { ol.throw("Type identifier expected for argument for imported "+imptype) }

				}
			}
		}
	}
	switch imptype{
		case "PROCEDURE","PROC","VOID":
			tar.idtype="PROCEDURE"
		case "FUNCTION","FUNC","DEF":
			tar.idtype="FUNCTION"
		case "VAR":
			tar.idtype="VAR"
		default:
			ol.throw("Cannot import '"+tar.idtype+"'")
	}
	tar.translateto = imptar
	tar.tarformed = true
	
	
	return tar,impid.Word // Allow the translator to do what is right here. 
}
