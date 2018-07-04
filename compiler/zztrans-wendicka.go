package scynt


func init(){
	
	
	
	
	TransMod["Wendicka"] = &T_TransMod {}
	tmw:=TransMod["Wendicka"]
	
	tmw.NameIdentifiers = func(p *TPackage) {
		for _,src := range p.sources {
			for _,id := range src.identifiers {
				if !id.tarformed {
					if id.idtype=="VAR" { id.translateto = "$"+id.translateto }
				}
			}
		}
	}
	
	
	
	
	
	
}
