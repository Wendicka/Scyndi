package scynt

import(
		"strings"
		"fmt"
	)
	

func legalidentifier(s string) (bool,string) {
	us:=strings.ToUpper(s)
	bs:=[]byte(us)
	ok:=true
	for i:=0;i<len(bs);i++{
		c:=bs[i]
		ok = ok && ( (c>='0' && c<='9' && i>0) || (c>='A' && c<='Z') || c=='_')
		if !ok {return false,fmt.Sprintf("Illegal character in identifier name: %c",rune(c))}
	}
	for _,kw:=range keywords{
		if kw==us { return false,kw+" is a keyword and can therefore NOT be used as an identifier name" }
	}
	return true,""
}
