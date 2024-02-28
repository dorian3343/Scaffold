@echo off
setlocal enabledelayedexpansion

cd database
echo Test no_database start %time%
call test.bat
if %ERRORLEVEL% == 1 (
    echo Error: test.bat in database directory failed with exit code %ERRORLEVEL%
    exit /b %ERRORLEVEL%
)
echo Test no_database end %time%
cd ..

cd no_database
echo Test database start %time%
call test.bat
if %ERRORLEVEL% == 1 (
    echo Error: test.bat in no_database directory failed with exit code %ERRORLEVEL%
    exit /b %ERRORLEVEL%
)
echo Test database end %time%
cd ..

echo E2E Tests Passed! Yippee
