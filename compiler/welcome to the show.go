package scynt

import "trickyunits/ansistring"
import "fmt"
import "os"

var Showdo = false

func doing(a,b string){
	if !Showdo { return }
	fmt.Print(ansistring.SCol(a,ansistring.A_Yellow,0),
	          ansistring.SCol(b,ansistring.A_Cyan,0))
}

func doingln(a,b string){
	doing(a,b)
	fmt.Println()
}


func progress(deel,geheel int){
	if !Showdo { return }
	d:=float64(deel)
	w:=float64(geheel)
	p:=(d/w)*100
	fmt.Print(ansistring.SCol(fmt.Sprintf("%f3.1%",p),ansistring.A_Magenta,0))
	fmt.Print("\010\010\010\010\010\010")
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

func done() {
	if !Showdo { return }
	fmt.Println(ansistring.SCol("Done!     ",ansistring.A_Green,0))
}
