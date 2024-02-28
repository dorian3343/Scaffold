@echo off
setlocal enabledelayedexpansion

start /b go run main.go > NUL 2>&1
rem Sleep without input redirection error
ping -n 1 127.0.0.1 > nul
rem GET http://localhost:8080/Greeting
for /f "delims=" %%a in ('curl -s http://localhost:8080/Greeting') do set "response=%%a"
if NOT "!response!" == "{"message":"Hello"}" (
    echo Wrong Response in test: no_database. Case : Greeting
    echo Received:" "!response!
    taskkill /f /im main.exe
    exit /b 1
)
rem GET http://localhost:8080/status
for /f "delims=" %%a in ('curl -s http://localhost:8080/status') do set "response=%%a"
if NOT "!response!" == "{"status":"OK"}" (
    echo Wrong Response in test: no_database. Case : status
    echo Received: "!response!"
    taskkill /f /im main.exe
    exit /b 1
)
rem GET http://localhost:8080/int
for /f "delims=" %%a in ('curl -s http://localhost:8080/int') do set "response=%%a"
if NOT "!response!" == "69" (
    echo  Wrong Response in test: no_database. Case : int
     echo Received:" "!response!
    taskkill /f /im main.exe
    exit /b 1
)
taskkill /f /im main.exe > NUL 2>&1
exit /b 0