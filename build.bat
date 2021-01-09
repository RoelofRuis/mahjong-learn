@echo off

setlocal
  cd server
  set GOOS=windows
  set GOARCH=amd64
  go build -o ../bin/win-server.exe
endlocal