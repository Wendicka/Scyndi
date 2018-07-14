package scynt

import "trickyunits/qstr"

func (s *tsource) declarevar(line []*tword) (string,tidentifier){
	//vr:=self.identifiers
	vr:=tidentifier{}
	vr.private=s.private
	vr.dttype="VARIANT"
	vr.idtype="VAR"
	vr.defaultvalue="NIL"
	if len(line)==0 {return "er:Empty variable declaration",vr } // should normally never happen, but at least this vovr a go panic crash
	n:=line[0]
	name:=n.Word
	vr.translateto="SCYNDI_VAR_"+s.srcname+"_"+name
	vr.tarformed=false
	if contains(keywords,name) { return "er:"+name+" is a keyword and may NOT be used as a variable",vr }
	if n.Wtype!="identifier" { return "er:Unexpected "+n.Wtype+"("+name+"). A name for a variable was expected",vr }
	if len(line)==1 {return name,vr}
	i:=1
	o:=line[i]
	if o.Word==":" {
		if len(line)<3 { return "er:Unexpected end of line. A type for a variable was expected",vr }
		n:=line[i+1]
		if !s.validtype(n) { return "er:Invalid variable type. Either an unknown type or invalud type: "+n.Word,vr }
		vr.dttype=n.Word
		i+=2
		vr.defstring = vr.dttype=="STRING"
	}
	if len(line)>i {
		o=line[i]
		if o.Word=="=" {
			i++
			o:=line[i]
			if len(line)<i+1 { return "er:Unexpected end of line",vr }
			switch vr.dttype {
				case "VARIANT":
					return "er:VARIANTS cannot be defined in a variable block",vr
				case "MAP","ARRAY":
					return "er:MAPS and ARRAYS cannot be defined in a variable block (Hey, psst! They are basically already defined upon declaration, so you don't have to)",vr
				case "STRING":
					if o.Wtype!="string" { return "er:Unexpected "+o.Wtype+". Constant string required",vr }
					vr.defaultvalue = o.Word
					vr.defstring=true
				case "INTEGER":
					if o.Wtype!="integer" { return "er:Unexpected "+o.Wtype+". Constant integer required",vr }
					vr.defaultvalue = o.Word
				case "FLOAT":
					if o.Wtype!="integer" && o.Wtype!="float" { return "er:Unexpected "+o.Wtype+". Constant integer or float required",vr }
					vr.defaultvalue = o.Word
				case "BOOLEAN":
					if o.Wtype!="keyword" {return "er:Unexpected "+o.Wtype+". TRUE or FALSE required",vr }
					if o.Word!="TRUE" && o.Word!="FALSE" { return "er:Unexpected "+o.Word+". TRUE or FALSE required!",vr }
					vr.defaultvalue = o.Word
				default:
					if o.Wtype!="keyword" || o.Word!="NEW" { return "er:Unexpected "+o.Wtype+" ("+o.Word+"). Only the keyword NEW is allowed for "+vr.dttype,vr }
					vr.defaultvalue = "NEW"
			}
		} else { return "er:Syntax error!",vr } // Now it's really beyond me what you were trying to do.... :-/
	} else {
		switch vr.dttype {			
			case "STRING": vr.defaultvalue = ""
			case "INTEGER","FLOAT": vr.defaultvalue="0"
			case "BOOLEAN":  vr.defaultvalue="FALSE"
		}
	}
	return name,vr
}

func (self *tsource) declarevars() string{
	for _,ol:=range self.varblock{
		n,i:=self.declarevar(ol.sline)
		if qstr.Prefixed(n,"er:") { ol.throw(n[3:]) }
		self.identifiers[n]=&i
	}
	t:=TransMod[TARGET]
	ret:=t.TransVars(self)
	return ret
}
