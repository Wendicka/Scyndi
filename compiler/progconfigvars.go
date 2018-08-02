/*
	Scyndi
	Program Configuration Variables
	
	
	
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
Version: 18.08.02
*/
package scynt

import "trickyunits/mkl"

type scynt_result struct{
	success, changed bool
}



// This variable will very likely NOT be used when parsing the keywords 
// when used in their proper locations. This array will rather be used
// to make sure none of them are used as identifier names when declaring variables or functions/procedures.
var keywords = []string {"PROGRAM","SCRIPT","MODULE","UNIT", // UNIT WILL AUTOMATICALLY BE REPLACED BY KEYWORD "MODULE"
						"USE","XUSE",
						"PROCEDURE","VOID","PROC",
						"FUNCTION","FUNC","DEF",
						"IF","ELSE","ELSEIF","ELIF",
						"SWITCH","CASE","DEFAULT",
						"WHILE",
						"FOR","FOREACH","FORU", // FORU = For Until It has only a very small difference with FOR, but still one with significant effect ;)
						"REPEAT","UNTIL","DO","LOOP","FOREVER",
						"AND","OR", "NOT",
						"END","BEGIN", // Begin will only be used as an alias for VOID MAIN in a program and VOID INIT in scripts and modules. So its function is NOT the same as in Pascal, although making it into a quick keyword WAS inpired by Pascal :P
						"PRIVATE","PUBLIC",
						"VAR","CONST","ENUM","OPTION","TYPE",
						"INTEGER","FLOAT","STRING","BOOLEAN","MAP","ARRAY","VARIANT",
						"IMPORT","INCLUDE","USE",
						"OPTION","SUPPORT","PURECODE",
						"TRUE","FALSE","NEW",
						"MOD",
						"NIL", 
						"RETURN",
					}
					
var operators = []string { // It's very important here, that the longer ones come first and the smaller ones later, or things might go wrong here.
	                      "...",
	                      "++","--",
	                      ":=",
	                      ":+",
	                      ":-",
	                      "==",
	                      ">=",
	                      "<=",
	                      "<>","!=","~=",
	                      "..",
	                      "=",
	                      ">",
	                      "<",
	                      "+","-","*","/",
	                      "(",")","[","]",
	                      
	                      ":",
	                      ",", // Strictly speaking not an operator, but for the splitting routines it'll count as one.
					  }

var globaldefs = map[string] bool {}

// NLSEP will if set to true (default value) accept a new line as a separator (and then you don't need a semi-colon at the end of each line), turning it off will require such a thing. Please note when putting multiple instructions on one line, the semi-colon will always be required to separate those.
var NLSEP = true

// This variable will contain the target. Default will be "Wendicka"
var TARGET = "Wendicka"

// Compilers will have to define the path here where all the system unit source files can be found for each target.
var SYSTEMDIR = "" 
var USEPATH = []string{}

var TARDIR = ""

func init(){
mkl.Version("Scyndi Programming Language - progconfigvars.go","18.08.02")
mkl.Lic    ("Scyndi Programming Language - progconfigvars.go","GNU General Public License 3")
}
