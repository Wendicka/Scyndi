/*
	Scyndi
	Main
	
	
	
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
Version: 18.07.27
*/
package main

import Scyndi "Scyndi/compiler"
import "flag"
import ans "trickyunits/ansistring"
import "trickyunits/mkl"
import "os"
import "fmt"
import "strings"
import "trickyunits/gini"
import "trickyunits/dirry"

var sourcefile string
var outputpath string

func init(){
mkl.Version("Scyndi Programming Language - main.go","18.07.27")
mkl.Lic    ("Scyndi Programming Language - main.go","GNU General Public License 3")
	nonl:=flag.Bool("sco",false,"If set new lines will not count as the end of an instruction and all instructions will have to be ended with a semi-colon")
	trgt:=flag.String("target","Wendicka","Set the target to translate to. (Supported targets: "+Scyndi.TargetsSupported()+")")
	outp:=flag.String("o","","Output file/dir. Please note that the output effect can be different depending on the chosen target")
	ver :=flag.Bool("version",false,"Show version information of all used source files to build Scorpion")
	ansi:=flag.String("ansi","","ON = force ANSI to be ON or OFF")
	flag.Parse()
	switch strings.ToUpper(*ansi) {
		case "ON":	ans.ANSI_Use=true
		case "OFF":	ans.ANSI_Use=false
	}
	if *ver {
		Copyright()
		fmt.Println(mkl.ListAll())
		os.Exit(0)
	}
	// configure
	Scyndi.NLSEP = !*nonl
	Scyndi.TARGET = *trgt
	Scyndi.Showdo = true
	outputpath=*outp
	myargs:=flag.Args()
	if len(myargs)>=1 { sourcefile=myargs[0] }
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
	trans,src:=Scyndi.CompileFile(sourcefile,"PROGRAM")  
	src.SaveTranslation(trans,outputpath)
}

func main(){
	Copyright()
	if sourcefile=="" { 
		InfoScreen()
	} else {
		LetsGo()
	}
}
