# Scyndi

![Scyndi](http://tricky1975.github.io/63/Icons/scyndi.png)

Experimental stuff.
Not much to see here yet.
If you really want to follow go ahead, but expect nothing unless further notice is given in da future!


Scyndi is an experiment programming language.
Its "compiler" is rather a translator which translates Scyndi code to other languages.
The aim is to make any language possible as long either I or other people can add modules to Scyndi allowing to make this happen.

A few aims for Scyndi as a language will be
- Symplistic Syntax
- Case INSENSITIVITY 
- Should cover most purposes if not all. 
  - Therefore Scyndi could cover compiler based languages (BlitzMax is my test case here)
  - My own "VM" (For that the Wendicka language was developed)
  - Scripting languages for addons (Lua will be my test case here)
  - Web scripting (php will be my testcase here.
- By means of "PURECODE" statements code from the target language can be directly imported or "in-lined"... This could basically be the point where Scyndi is no longer a simple language, but for building modules for the rest of us, this could help a lot. ;)
- I do have plans to make OOP possible

Deliberate limitations and the reasons why:
- Small group of base types. This may not be most efficient on the RAM, but makes things easier on the Scyndi language
- "if a==1" and "if a=1" will have the same effect. I hear C-junkies grunt in anger, but I will not allow variable definitions inside a boolean-check, as this is one of the pieces of rope C provides to hang yourself.
- Any type not standard type will be pointers only. With that I follow the line of BlitzMax, Lua and Python. This to prevent loads of syntax conflicts (which I already experienced in Go, thank you). I also want to reserve the "\*" symbol (as C and Go use for pointer calls) for multiplications only and the "^" symbol (which Pascal uses) only for empower math... And when it comes to OOP, pointers are the better choince anyway.


And I hear the question, will the generated code be as powerful as when written directly in the target language. I've always said to people using older translation based languages translating to C.
"If you want the full power of C, use C".
My prime concern is getting stuff to work, and then trying to make it as optimal as possible will certainly come into play, and for now it's too early to say if Scyndi will be efficient enough to fully replace their target languages. It's a personal experiment anyway, so future tests will have to tell. But always keep in mind, no matter how efficient a translation based language is, nothing beats using the target language for real, no exceptions!

Now the prototype for Scyndi is being written in Go, which was for me just because of Go having a compiler that easily compiles to all three great platforms (Mac, Linux and Windows). I have no short-term plans for writing Scyndi in Scyndi allowing it to "compile itself", but maybe on the longer run, it would be great if it could, right? If that will ever happen a new repository will very likely be started and this will then be an 'archive repository'. This is all long term planning, though.
