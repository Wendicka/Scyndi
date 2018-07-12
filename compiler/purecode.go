package scynt

//import "fmt"

func purecode(s *tsource,c *tchunk,ol *tori) string {
	 /*
	fmt.Println(len(ol.sline) )
	for i,w:=range(ol.sline){
		fmt.Println(i,"\t",w.Wtype,"\t",w.Word)
	}
	// */
	if len(ol.sline)!=4 { ol.throw("Invalid PURCODE instruction") }
	tar:=ol.sline[1]
	comma:=ol.sline[2]
	code:=ol.sline[3]
	if comma.Word!="," { ol.throw("Comma expected") }
	if tar.Word!=TARGET { return "" } // Only do this if the target is correct
	
	// identifier replacement will have to take place here
	ret:=code.Word
	
	return ret
}
