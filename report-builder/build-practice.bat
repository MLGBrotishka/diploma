cd /d "%~dp0"
powershell.exe -command .\build.ps1 ^
-md ..\report-practice\r-main.md ^
-template ..\report-practice\template-report.docx ^
-docx ..\practice.docx ^
-embedfonts ^
-counters ^
-flags "--resource-path=..\report-practice"