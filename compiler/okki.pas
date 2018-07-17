program okki;

{ This is just a dirty way to generate a quick array Go needs in Scyndi }

var 
	aantal:integer=1;
	digi:array[1..10] of integer;
	dl:array[1..10] of char;
	i,j:integer;
	
begin
	writeln('package scynt');
	writeln('var okki = map[string]byte{}');
	writeln('func init(){');
	dl[1]:='r';
	dl[2]:='b';
	dl[3]:='n';
	dl[4]:='t';
	{dl[4]:='"';}
	{dl[5]:='0';}
	writeln('	okki["\\\""] = ''"''');
	for i:=1 to 4 do 
		writeln('	okki["\\',dl[i],'"]=''\',dl[i],''';');
	for i:=1 to 10 do
		digi[i]:=0;
	digi[1]:=-1;
	for i:=0 to 255 do begin
		digi[1]:=digi[1]+1;
		for j:=1 to 9 do begin
			if digi[j]>7 then begin
				digi[j]:=digi[j]-8;
				digi[j+1]:=digi[j+1]+1;
				if aantal<j+1 then aantal:=j+1
			end
		end;
		write('	okki["\\');
		for j:=3 downto 1 do write(digi[j]);
		write('"] = ''\');
		for j:=3 downto 1 do write(digi[j]);
		writeln(''';	// ',i:3,' //')
	end;
	writeln('}')
end.
