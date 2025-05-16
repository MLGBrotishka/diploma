cd /d "%~dp0"
powershell.exe -command ^
pandoc ..\report\pres.md -o ..\pres.pptx --resource-path=..\report --reference-doc=..\report\template-pres.pptx