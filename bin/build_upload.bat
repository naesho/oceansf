set GOARCH=amd64
set GOOS=linux
go build github.com/naesho/oceansf/

@echo off
echo user ftpuser> ftpcmd.dat
echo ftpxptmxm@1234>> ftpcmd.dat
echo bin>> ftpcmd.dat
echo cd /app/bin>>ftpcmd.dat
echo put oceansf>> ftpcmd.dat
echo quit>> ftpcmd.dat
ftp -n -s:ftpcmd.dat smon-dev
del ftpcmd.dat
C:\Program Files (x86)\WinSCP\WinSCP.exe