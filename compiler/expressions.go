package scynt

func defaultexpressiontranslation(expect string,source *tsource, ol *tori,start,level int) (endpos int,ex string){
	endpos=start
	ex="Temp shit"
	return
}


func (s *tsource) translateExpressions(expect string,ol *tori,start,level int) (endpos int,ex string){
	trans:=TransMod[TRANS]
	et:=defaultexpressiontranslation
	if trans.transexp { et = transexp }
	endpos,ex = et(expect,s,ol,start,level)
	return
}
