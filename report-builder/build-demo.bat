cd /d "%~dp0"
powershell.exe -command .\build.ps1 ^
-md ..\report-demo\r-beginning.md, ..\report-demo\r-main.md, ..\report-demo\r-end.md ^
-template ..\report-demo\template-report.docx ^
-docx ..\demo.docx ^
-embedfonts ^
-counters ^
-flags "--resource-path=..\report-demo"
echo "Не забудь шрифт оглавления и разрывы таблиц"