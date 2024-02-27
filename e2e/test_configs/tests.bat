@echo off
setlocal enabledelayedexpansion

cd database
call test.bat
if %ERRORLEVEL% == 1 (
    echo Error: test.bat in database directory failed with exit code %ERRORLEVEL%
    exit /b %ERRORLEVEL%
)
cd ..

cd no_database
call test.bat
if %ERRORLEVEL% == 1 (
    echo Error: test.bat in no_database directory failed with exit code %ERRORLEVEL%
    exit /b %ERRORLEVEL%
)
cd ..

echo E2E Tests Passed! Yippee
