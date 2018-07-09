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
 */
 
  
