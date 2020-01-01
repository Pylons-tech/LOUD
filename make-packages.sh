#! /bin/bash

OSSES="darwin linux windows"
ARCHES="amd64"
ZIPPREFIX=`date +%Y-%m-%d`

for OS in $OSSES
do
	for ARCH in $ARCHES
	do
		echo app-$OS-$ARCH
		OUTDIRNAME=loud-$OS-$ARCH
		mkdir $OUTDIRNAME
		ENABLE_CGO=""

		INFILE=cmd/loud.go
		OUTITEM=$OUTDIRNAME/loud
		if [ "z$OS" = "zlinux" ]
		then
			echo "Linux"
		elif [ "z$OS" = "zwindows" ]
		then
			OUTITEM=$OUTITEM.exe
		elif [ "z$OS" = "zdarwin" ]
		then
			INFILE=cmd/loud-ui.go
		fi

		# brew install mingw-w64
		GOOS=$OS GOARCH=$ARCH go build -o $OUTITEM $INFILE

		cp *.json *.txt $OUTDIRNAME
		if [ "z$OS" = "zwindows" ]
		then
			zip loud-$ZIPPREFIX-$OS-$ARCH.zip $OUTDIRNAME/*
		else
			tar cvzf loud-$ZIPPREFIX-$OS-$ARCH.tgz $OUTDIRNAME
		fi
		rm -rf $OUTDIRNAME
	done
done
