@echo off
setlocal enabledelayedexpansion

start /b go run main.go > NUL 2>&1
rem Sleep without input redirection error
ping -n 6 127.0.0.1 > nul

rem GET http://localhost:8080/get_user test 1
for /f "delims=" %%a in ('curl -s http://localhost:8080/get_user') do set "response=%%a"
if NOT "!response!" == "null" (
    echo Wrong Response in test: database. Case : get_user 1
    echo Received: !response!
    taskkill /f /im main.exe
    del main.db
    exit /b 1
)

set "json_data={\"Name\":\"John Doe\",\"Age\":30}"
curl -s -X POST ^
     -H "Content-Type: application/json" ^
     -d "%json_data%" ^
     http://localhost:8080/post_user > NUL 2>&1
if %errorlevel% neq 0 (
    echo Request failed in test: database. Case : post_user
    del main.db
    taskkill /f /im main.exe
    exit /b 1
)
rem GET http://localhost:8080/get_user test 2
rem Rewrite this
set json="[{age:30,id:1,name:John Doe}]"
for /f "delims=" %%a in ('curl -s http://localhost:8080/get_user') do set "response=%%a"
for /f "tokens=* delims=" %%b in ("%response%") do set "response=%%b"
set "response=%response:"=%"
if not "!response!"==%json% (
    echo Wrong Response in test: database. Case : get_user 2
    echo Received: !response!
    echo %json%
    taskkill /f /im main.exe
    del main.db
    exit /b 1
)
taskkill /f /im main.exe > NUL 2>&1
del main.db
exit /b 0