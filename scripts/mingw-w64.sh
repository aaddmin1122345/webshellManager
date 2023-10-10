#/bin/bash
# Arch build
echo "请全部安装!"
paru -S mingw-w64

export CC=x86_64-w64-mingw32-gcc
export CXX=x86_64-w64-mingw32-g++

if [ -f "go.mod" ]; then
  GOARCH=amd64 CGO_ENABLED=1 GOOS=windows go build -o webshellManager.exe
else
  cd ../;GOARCH=amd64 CGO_ENABLED=1 GOOS=windows go build -o webshellManager.exe
fi


echo "编译完成!"
unset CC
unset CXX
