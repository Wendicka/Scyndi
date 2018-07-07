package main

import Scyndi "Scyndi/compiler"
import "flag"
import ans "trickyunits/ansistring"
import "trickyunits/mkl"
import "os"
import "fmt"
import "trickyunits/gini"
import "trickyunits/dirry"

var sourcefile string
var outputpath string

func init(){
	nonl:=flag.Bool("sco",false,"If set new lines will not count as the end of an instruction and all instructions will have to be ended with a semi-colon")
	trgt:=flag.String("target","Wendicka","Set the target to translate to. (Supported targets: "+Scyndi.TargetsSupported()+")")
	outp:=flag.String("o","","Output file/dir. Please note that the output effect can be different depending on the chosen target")
	flag.Parse()
	// configure
	Scyndi.NLSEP = !*nonl
	Scyndi.TARGET = *trgt
	Scyndi.Showdo = true
	outputpath=*outp
	myargs:=flag.Args()
	if len(myargs)>=1 { sourcefile=myargs[0] }
	mkl.Version("s","18.1.12")
	mkl.Lic("s","To be decided later!")
}

func Copyright(){
	nv:=mkl.Newest()
	fmt.Println(ans.SCol("Scorpion -- Scyndi compiler/translator",ans.A_Yellow,0))
	fmt.Println(ans.SCol("Version: "+nv,ans.A_Cyan,0))
	fmt.Println(ans.SCol("(c) 2018-20"+nv[:2]+" Jeroen Petrus Broks",ans.A_Magenta,0))
	fmt.Println();
}

func InfoScreen(){
	fmt.Println(ans.SCol("Usage: ",ans.A_Red,0),ans.SCol("scorpion ",ans.A_Yellow,0),ans.SCol("[ flags ] ",ans.A_Blue,0),ans.SCol("<<Sourcefile>>",ans.A_Cyan,0))
	fmt.Println();
	flag.PrintDefaults()
	os.Exit(0)
}

func ReadConfig(){
	cdir:=dirry.Dirry("$AppSupport$/ScyndiScorpio/")
	g:=gini.ReadFromFile(cdir+"scorpion.gini")
	Scyndi.SYSTEMDIR=g.C("SystemMods")
}


func LetsGo(){
	ReadConfig()
	//translation,source:=
	Scyndi.CompileFile(sourcefile,"PROGRAM")  // Must be moved up once everything goes, and live above must be unremmed...
}

func main(){
	Copyright()
	if sourcefile=="" { 
		InfoScreen()
	} else {
		LetsGo()
	}
}
