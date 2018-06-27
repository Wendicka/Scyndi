# Keywords

This is just an overview of the planned keywords

# First line declaration

Every Scyndi file (except for include files) MUST have either one of these instructions on the first line

- PROGRAM
  - Program will require a Main procedure (keyword "BEGIN" will then act as a "PROCEDURE MAIN" instruction)
  - Program is only meant for sources being to be compiled in actual .exe files (like when translating to compiler based languages) or interpreters being set up for working as programs (such as Python).


- SCRIPT
  - Script does not require anything specific, and are specifically set up for setting up add-ons for programs or when translating to languages that are 100% script (such as Lua).
  - The keyword "BEGIN" will act as "PROCEDURE INIT" which will basically be run when the translated code is loaded, before any other action takes place.

- MODULE
  - Set up for being imported with USE.
  - Modules are not meant for direct running, mostly only to be imported by their main programs/scripts
  - The keyword "BEGIN" will act as "PROCEDURE INIT" which will basically be run when the translated code is loaded, before any other action takes place.



# Other keywords
- USE
  - Imports module into the main script or program.
  - Of course modules may call for other modules.
  - Depending on the support of the target language, the translation of this code is merged with the main code.

- XUSE
  - Basically the same as use, but depending on the support of the target language the translated code will remain in a separate file, acting as a dependency.



- PROCEDURE / PROC / VOID
  - Official keyword is "PROCEDURE", but "PROC" and "VOID" will do the same thing.
  - These are used for functions that won't return a value.
  - It's up to you if you want to make use of these, but depending on the target language this can be pretty important. Lua and BlitzMax won't care if you use a PROCEDURE or a FUNCTION as the generated code will be the same, but Pascal makes a difference between the two, and GO won't even compile if a return statement is not present in a normal function. All in all, if you want your Scyndi code to work on whatever translation module may ever arrive, use procedures when no values have to be returned.

- FUNCTION / FUNC / DEF
  - Official keyword if function, but FUNC and DEF have the same function
  - For functions that will return values.

- IF
  - Start if statement.
  - No keyword then, just as soon as a separator is found, the code within the statement starts

- ELSE
  - Else statement

- ELSEIF / ELIF
  - ELSEIF statement
  - ELIF is just an alias

- SWITCH, CASE, DEFAULT
  - Commands for the Switch statement

- WHILE
  - Start while statement

- FOR
  - FOR loop
- FORUNTIL
  - Same as FOR loop, but where a normal FOR stops after the loop value has equalled its ending value (or tops it), FORUNTIL will stop immediately when this is reached.

- REPEAT / DO
  - Start repeat statement
  - DO is the alias

- UNTIL
  - Ends REPEAT if condition is met

- LOOP / FOREVER
  - Ends REPEAT and loops forever

- AND, OR, NOT
  - Regular boolean keywords

- END
  - Ends any statement
  - Alternatively when using VAR/CONST/OPTION/USE without keyword, but a separator in stead, you have to end those with END as well

- BEGIN
  - In programs it will be a "PROCEDURE MAIN"
  - In scripts and modules it will be a "PROCEDURE INIT"


- PRIVATE / PUBLIC
  - In modules they can be used to determine if the programs, scripts, parent modules, calling them may or may not use the identifiers declared on the next lines

- VAR
  - Variable declaration

- CONST
  - Constant definition

- OPTION
  - Always comes first after PROGRAM, SCRIPT and MODULE
  - Sets config for the file to be compiled.

- "INTEGER","FLOAT","STRING","BOOLEAN","MAP","ARRAY","VARIANT"
  - Basic data types

- IMPORT
  - Does NOT do what you expect, I'm afraid :P
  - With IMPORT you can declare an identifier from the target language directly as an identifier in Scyndi.
  - Must be accompanied with a declaration for Scyndi to understand

- INCLUDE
  - Includes a file and merges the code with the parent completely


- SUPPORT
  - Tells the compiler which target languages are supported for this file. If left out all target languages that Scorpion can translate to will be supported.
  - Prefixing with a ! will allow all languages except the one you put in this way.
