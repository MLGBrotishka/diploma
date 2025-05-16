cd /d "%~dp0"
powershell.exe -command .\build.ps1 ^
-md ..\report\r-beginning.md, ..\report\r-main.md, ..\report\r-end.md ^
-template ..\report\template-report.docx ^
-docx ..\report.docx ^
-embedfonts ^
-counters ^
-flags "--resource-path=..\report"
echo "Не забудь шрифт оглавления и разрывы таблиц"