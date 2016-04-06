
#
# Copyright (C) Philip Schlump, 2016.
#
# MIT Licnesed. 
#

all:
	go build

test1:
	go build
	mkdir -p ./ref ./out
	./ifit -i note.1 -o ./out/aa.out NameA
	diff ./out/aa.out ./ref/aa.out
	./ifit -i note.1 -o ./out/bb.out NameB
	diff ./out/bb.out ./ref/bb.out
	./ifit -i note.1 -o ./out/cc.out NameC
	diff ./out/cc.out ./ref/cc.out
	./ifit -i note.1 -o ./out/dd.out NameD
	diff ./out/dd.out ./ref/dd.out
	./ifit -i note.1 -o ./out/ee.out NameE
	diff ./out/ee.out ./ref/ee.out
	./ifit -i note.1 -o ./out/ab.out NameA NameB
	diff ./out/ab.out ./ref/ab.out
	
test2:
	go build
	mkdir -p ./ref ./out
	./ifit -i note.2 -o ./out/aa2.out -s sub1.json NameA NameB
	diff ./out/aa2.out ./ref/aa2.out

install: 
	go build
	cp ifit ~/bin

