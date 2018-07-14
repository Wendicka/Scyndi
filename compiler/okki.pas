program okki;

{ This is just a dirty way to generate a quick array Go needs in Scyndi }

var 
	aantal:integer=1;
	digi:array[1..10] of integer;
	dl:array[1..10] of char;
	i,j:integer;
	
begin
	writeln('package scynt');
	writeln('var okki = map[sring]string{}');
	writeln('func init(){');
	dl[1]:='r';
	dl[2]:='b';
	dl[3]:='n';
	dl[4]:='0';
	for i:=1 to 4 do 
		writeln('	okki["\\',dl[i],'"]=''\\',dl[i],''';');
	for i:=1 to 10 do
		digi[i]:=0;
	for i:=1 to 255 do begin
		digi[1]:=digi[1]+1;
		for j:=1 to 9 do begin
			if digi[j]>7 then begin
				digi[j]:=digi[j]-7;
				digi[j+1]:=digi[j+1]+1;
				if aantal<j+1 then aantal:=j+1
			end
		end;
		write('	okki["\0');
		for j:=1 to aantal do write(digi[j]);
		write('"] = ''\0');
		for j:=1 to aantal do write(digi[j]);
		writeln(''';	{',i,'}')
	end;
	writeln('}')
end.
