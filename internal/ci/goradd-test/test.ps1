# Windows version of test.sh
# This script is an aid to running the browser test as a github action. It is designed to be run from the goradd-test directory.
go build goradd-test
Start-Job -Name "testserver" -ScriptBlock { ./test-helper.ps1 }
./goradd-test.exe