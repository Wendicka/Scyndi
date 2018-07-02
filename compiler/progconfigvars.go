package scynt

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
						"FOR","FOREACH",
						"REPEAT","UNTIL","DO","LOOP","FOREVER",
						"AND","OR",// "NOT"
						"END","BEGIN", // Begin will only be used as an alias for VOID MAIN in a program and VOID INIT in scripts and modules. So its function is NOT the same as in Pascal, although making it into a quick keyword WAS inpired by Pascal :P
						"PRIVATE","PUBLIC",
						"VAR","CONST","OPTION","TYPE",
						"INTEGER","FLOAT","STRING","BOOLEAN","MAP","ARRAY","VARIANT",
						"IMPORT","INCLUDE","USE",
						"OPTION","SUPPORT","PURECODE",
						"TRUE","FALSE","NEW",
						"NIL", 
					}
					
var operators = []string { // It's very important here, that the longer ones come first and the smaller ones later, or things might go wrong here.
	                      ":=",
	                      ":+",
	                      ":-",
	                      "==",
	                      ">=",
	                      "<=",
	                      "..",
	                      "=",
	                      ">",
	                      "<",
	                      "+","-","*","/",
	                      "(",")",
	                      
	                      ":",
	                      ",", // Strictly speaking not an operator, but for the splitting routines it'll count as one.
					  }

// NLSEP will if set to true (default value) accept a new line as a separator (and then you don't need a semi-colon at the end of each line), turning it off will require such a thing. Please note when putting multiple instructions on one line, the semi-colon will always be required to separate those.
var NLSEP = true
