#!/bin/sh
echo "Building server..."
cd server
go build
mv server ../bin/authserv
cd ..
cd tools/netkey
echo "Building netkey..."
go build
mv netkey ../../bin/netkey
cd ..
cd adminconsole
echo "Building adminconsole..."
go build
mv adminconsole ../../bin/adminconsole
cd ../..
echo "Done."

