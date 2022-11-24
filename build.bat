@ECHO off

SET WORKDIR=%~dp0

SET OPERATOR_ENGINE=%WORKDIR%operator
SET OPERATOR_UI=%WORKDIR%ui

ECHO.
ECHO --- WORKDIR: %WORKDIR%
ECHO --- OPERATOR_ENGINE: %OPERATOR_ENGINE%
ECHO --- OPERATOR_UI: %OPERATOR_UI%
ECHO.

REM ### Build Operator ENGINE
ECHO "Build Operator ENGINE"
ECHO.

cd %OPERATOR_ENGINE%
docker build -t jnnkrdb/configurj-engine:latest .
docker push jnnkrdb/configurj-engine:latest

set /p enginetag_release=Set the ReleaseTag of the Engine-Container:

if "%enginetag_release%" == "" goto END

docker tag jnnkrdb/configurj-engine:latest jnnkrdb/configurj-engine:%enginetag_release%
docker push jnnkrdb/configurj-engine:%enginetag_release%

goto END

REM ### Build Operator UI
ECHO "Build Operator UI"
ECHO.

cd %OPERATOR_UI%
docker build -t jnnkrdb/configurj-ui:latest .
docker push jnnkrdb/configurj-ui:latest

set /p uitag_release=Set the ReleaseTag of the UI-Container:

docker tag jnnkrdb/configurj-ui:latest jnnkrdb/configurj-ui:%uitag_release%
docker push jnnkrdb/configurj-ui:%uitag_release%

REM ### back to original workdir
:END
cd %WORKDIR%