set -exu
HERE=$(dirname $(realpath $BASH_SOURCE))
cd $HERE

go build -o ttt.exe ttt.go
./ttt.exe
