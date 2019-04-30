@echo off
set /p VERSION=please input today version:
echo today version:%VERSION%

set GOARCH=amd64
set GOOS=linux
set CGO_ENABLED=0

rem 导出的根路径
rem %EXPORT_PATH%/out目录下是需要上传到ftp的资源文件目录
rem %EXPORT_PATH%/server目录下最新的资源
set EXPORT_PATH=D:\test
rem 配置文件路径
set CONF_PATH=G:\Server\GbServer\bin
rem 压缩工具路径
set ZIP_TOOL_PATH=E:\Haozip\HaoZipC.exe
rem go path路径
set GOPATH=G:\Library;G:\Server\GbServer

set ZIP_NAME=GbxcServer_%date:~0,4%%date:~5,2%%date:~8,2%%VERSION%
set OUTPUT_ZIP_FILE=%EXPORT_PATH%\out\%ZIP_NAME%.zip
set OUTPUT_MD5_FILE=%OUTPUT_ZIP_FILE%.txt

go build -o %EXPORT_PATH%\server\gb_server main.go
xcopy /S /Y %CONF_PATH% %EXPORT_PATH%\server\
%ZIP_TOOL_PATH% a -tzip %OUTPUT_ZIP_FILE% %EXPORT_PATH%\server

echo.

call certutil -hashfile %OUTPUT_ZIP_FILE% MD5>%OUTPUT_MD5_FILE%
setlocal enabledelayedexpansion
for /f "usebackq delims=" %%a in (%OUTPUT_MD5_FILE%) do (
   set /a n+=1
   if !n! geq 2 (
      if !n! leq 2 echo %%a>%OUTPUT_MD5_FILE%
   )
)

echo output zip: %OUTPUT_ZIP_FILE%
echo output md5: %OUTPUT_MD5_FILE%
echo build success!

::@echo Off
::echo open 120.25.206.203 21 >ftp.up
::echo gbxc>>ftp.up
::echo lBBR2wd7k5rQheSHM4TA>>ftp.up

::echo Cd .\ >>ftp.up
::echo binary>>ftp.up
::echo put %OUTPUT_ZIP_FILE%>>ftp.up
::echo put %OUTPUT_MD5_FILE%>>ftp.up
::echo bye>>ftp.up
::FTP -s:ftp.up
::del ftp.up /q
