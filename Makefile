main: output.asm
	C:\masm32\bin\ml /c /Zd /coff output.asm && C:\masm32\bin\Link /SUBSYSTEM:CONSOLE output.obj && .\output.exe
