cd src
set GOARCH=amd64
set GOOS=linux

@REM C:\Users\Administrator\go\go1.17.5\bin\go.exe build -i -o C:\Users\Administrator\Desktop\project\github\instagram_project\register makemoney/routine/register
@REM C:\Users\Administrator\go\go1.17.5\bin\go.exe build -i -o C:\Users\Administrator\Desktop\project\github\instagram_project\refresh_cookies makemoney/routine/refresh_cookies
@REM C:\Users\Administrator\go\go1.17.5\bin\go.exe build -i -o C:\Users\Administrator\Desktop\project\github\instagram_project\crawling_tags makemoney/routine/crawling_tags

C:\Users\Administrator\go\go1.17.5\bin\go.exe build -o C:\Users\Administrator\Desktop\project\github\instagram_project\short_link C:\Users\Administrator\Desktop\project\github\instagram_project\src\routine\short_link\short_link.go

