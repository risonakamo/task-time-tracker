set -exu
HERE=$(dirname $(realpath $BASH_SOURCE))
cd $HERE

webdir=$HERE/../../task-time-tracker-web

# generate wails for web
rm -rf wails-js-out
wails generate module
rm -rf $webdir/web/wailsjs
mv wails-js-out/wailsjs $webdir/web/wailsjs
rm -rf wails-js-out

# build web
cd $webdir
pnpm b

# copy in web
cd $HERE
rm -rf web-build
cp -r $webdir/build ./web-build

# build go
go build -tags dev -gcflags "all=-N -l" ttt-desktop.go
./ttt-desktop.exe