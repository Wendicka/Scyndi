{ --- START LICENSE BLOCK ---
	Scyndi
	Generate okki.go
	
	
	
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
Version: 18.07.21
  --- END LICENSE BLOCK --- } 
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
