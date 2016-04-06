
#
# Copyright (C) Philip Schlump, 2016.
#
# MIT Licnesed. 
#

all:
	go build

# basics
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
	echo PASS

# variable substitution	
test2:
	go build
	mkdir -p ./ref ./out
	./ifit -i note.2 -o ./out/aa2.out -s sub1.json NameA NameB
	diff ./out/aa2.out ./ref/aa2.out
	echo PASS

# markers not in col(1)
test3:
	go build
	mkdir -p ./ref ./out
	./ifit -i note.3 -o ./out/aa3.out -s sub1.json NameA NameB
	diff ./out/aa3.out ./ref/aa3.out
	echo PASS

# nested ifs
test4:
	go build
	mkdir -p ./ref ./out
	./ifit -i note.4 -o ./out/aa4_1.out -s sub1.json NameA NameB
	diff ./out/aa4_1.out ./ref/aa4_1.out
	./ifit -i note.4 -o ./out/aa4_2.out -s sub1.json NameA 
	diff ./out/aa4_2.out ./ref/aa4_2.out
	./ifit -i note.4 -o ./out/aa4_3.out -s sub1.json NameB
	diff ./out/aa4_3.out ./ref/aa4_3.out
	./ifit -i note.4 -o ./out/aa4_4.out -s sub1.json NameC
	diff ./out/aa4_4.out ./ref/aa4_4.out
	./ifit -i note.4 -o ./out/aa4_5.out -s sub1.json NameA NameC
	diff ./out/aa4_5.out ./ref/aa4_5.out
	./ifit -i note.4 -o ./out/aa4_6.out -s sub1.json NameC NameA
	diff ./out/aa4_6.out ./ref/aa4_6.out
	echo PASS

# test "else"
test5:
	echo FAIL

install: 
	go build
	cp ifit ~/bin

