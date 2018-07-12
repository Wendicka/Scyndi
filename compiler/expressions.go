package scynt

type texi struct{
	word string
	wtype string
	subex *tex
}

type tex struct{
	oriline *tori
	orisrc *tsource
	exi []*texi
	cptype string
	tottype string 
}

func (s *tsource) translateExpressions(expect string){
}
