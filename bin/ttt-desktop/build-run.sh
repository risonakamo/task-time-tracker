set -exu
HERE=$(dirname $(realpath $BASH_SOURCE))
cd $HERE

cp -r ../../task-time-tracker-web/build ./web-build

wails generate module

go build -tags dev -gcflags "all=-N -l" ttt-desktop.go
# ./ttt-desktop.exe