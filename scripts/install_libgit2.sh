#! /bin/sh

git clone https://github.com/libgit2/libgit2.git ~/libgit2
cd ~/libgit2
git checkout v1.5.0
mkdir build && cd build
cmake .. -DBUILD_SHARED_LIBS=0 -DCMAKE_C_FLAGS="-fPIC"
cmake --build . --target install
