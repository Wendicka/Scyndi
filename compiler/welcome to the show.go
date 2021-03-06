/*
	Scyndi
	Show stuff on screen during translating/compiling/whatever
	
	
	
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
Version: 18.08.04
*/
package scynt

import "trickyunits/ansistring"
import "trickyunits/mkl"
import "fmt"
import "os"

var Showdo = false

func doing(a,b string){
	if !Showdo { return }
	fmt.Print(ansistring.SCol(a,ansistring.A_Yellow,0),
	          ansistring.SCol(b,ansistring.A_Cyan,0))
	fmt.Print(" ")
}

func doingln(a,b string){
	doing(a,b)
	
	fmt.Println()
}


func progress(deel,geheel int){
	if !Showdo { return }
	d:=float64(deel)
	w:=float64(geheel)
	p:=float64((d/w)*100)
	pi:=int(p)
	fmt.Print(ansistring.SCol(fmt.Sprintf("%3d",pi)+"%",ansistring.A_Magenta,0))
	fmt.Print("\010\010\010\010")
}

func throw(e string){
	fmt.Println(ansistring.SCol("ERROR",ansistring.A_Red,ansistring.A_Blink))
	fmt.Println(ansistring.SCol(e,ansistring.A_Yellow,0))
	os.Exit(1)
}

func ethrow(e error){
	throw(e.Error())
}

func lthrow(f string, l int,e string){
	fmt.Println(ansistring.SCol("ERROR",ansistring.A_Red,ansistring.A_Blink))
	fmt.Println(ansistring.SCol(f+":",ansistring.A_Cyan,0)+" "+ansistring.SCol(fmt.Sprintf("%d",l),ansistring.A_Magenta,0)+"\t\t"+ansistring.SCol(e,ansistring.A_Yellow,0))
	os.Exit(1)
}

func lwarn(f string, l int,e string){
	fmt.Println(ansistring.SCol("WARNING",ansistring.A_White,ansistring.A_Blink))
	fmt.Println(ansistring.SCol(f+":",ansistring.A_Cyan,0)+" "+ansistring.SCol(fmt.Sprintf("%d",l),ansistring.A_Magenta,0)+"\t\t"+ansistring.SCol(e,ansistring.A_Yellow,0))
}

func warn(e string){
	fmt.Println(ansistring.SCol("WARNING",ansistring.A_White,ansistring.A_Blink))
	fmt.Println(ansistring.SCol(e,ansistring.A_Yellow,0))
}

func sumdot(i int){
	s:=fmt.Sprintf("%5d. ",i)
	fmt.Print(ansistring.SCol(s,ansistring.A_Green,0))
}

func lassert(f string,l int, check bool, e string){
	if !check { lthrow(f,l,e) }
}

func (s *tori) throw(e string) {
	lthrow ( s.sfile,s.ln,e )
}

func (s *tori) ethrow(e error){
	s.throw(e.Error())
}


func (s *tori) warn(e string) {
	lwarn ( s.sfile,s.ln,e )
}

func (s *tori) pthrow(e string) {
	if len(e)<4 { return }
	ce:=e[:3]
	pe:=e[3:]
	if ce=="er:" {throw(pe)}
}

func done() {
	if !Showdo { return }
	fmt.Println(ansistring.SCol("Done!     ",ansistring.A_Green,0))
}

func hey() {
	if !Showdo { return }
	fmt.Println(ansistring.SCol("HEY!     ",ansistring.A_Magenta,0))
}

func init(){
mkl.Lic    ("Scyndi Programming Language - welcome to the show.go","GNU General Public License 3")
mkl.Version("Scyndi Programming Language - welcome to the show.go","18.08.04")
}
