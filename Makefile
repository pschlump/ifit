
#
# Copyright (C) Philip Schlump, 2016.
#
# MIT Licnesed. 
#

all:
	go build

.PHONY: test pre_test test1 test2 test3 test4 test5 test6 test7 test8 test9 test10

test: test1 test2 test3 test4 test5 test6 test7 test8 test9 test10

pre_test:
	go build
	mkdir -p ./ref ./out

# basics
test1: pre_test
	./ifit -m test -i note.1 -o ./out/aa.out NameA
	diff ./out/aa.out ./ref/aa.out
	./ifit -m test -i note.1 -o ./out/bb.out NameB
	diff ./out/bb.out ./ref/bb.out
	./ifit -m test -i note.1 -o ./out/cc.out NameC
	diff ./out/cc.out ./ref/cc.out
	./ifit -m test -i note.1 -o ./out/dd.out NameD
	diff ./out/dd.out ./ref/dd.out
	./ifit -m test -i note.1 -o ./out/ee.out NameE
	diff ./out/ee.out ./ref/ee.out
	./ifit -m test -i note.1 -o ./out/ab.out NameA NameB
	diff ./out/ab.out ./ref/ab.out
	echo PASS

# variable substitution	
test2: pre_test
	./ifit -m test -i note.2 -o ./out/aa2.out -s sub1.json NameA NameB
	diff ./out/aa2.out ./ref/aa2.out
	echo PASS

# more variable substitution	
test6: pre_test
	./ifit -m prod -i note.2 -o ./out/aa2_test6.out -s sub1.json NameA NameB
	diff ./out/aa2_test6.out ./ref/aa2_test6.out
	echo PASS


# markers not in col(1)
test3: pre_test
	./ifit -m test -i note.3 -o ./out/aa3.out -s sub1.json NameA NameB
	diff ./out/aa3.out ./ref/aa3.out
	echo PASS

# nested ifs
test4: pre_test
	./ifit -m test -i note.4 -o ./out/aa4_1.out -s sub1.json NameA NameB
	diff ./out/aa4_1.out ./ref/aa4_1.out
	./ifit -m test -i note.4 -o ./out/aa4_2.out -s sub1.json NameA 
	diff ./out/aa4_2.out ./ref/aa4_2.out
	./ifit -m test -i note.4 -o ./out/aa4_3.out -s sub1.json NameB
	diff ./out/aa4_3.out ./ref/aa4_3.out
	./ifit -m test -i note.4 -o ./out/aa4_4.out -s sub1.json NameC
	diff ./out/aa4_4.out ./ref/aa4_4.out
	./ifit -m test -i note.4 -o ./out/aa4_5.out -s sub1.json NameA NameC
	diff ./out/aa4_5.out ./ref/aa4_5.out
	./ifit -m test -i note.4 -o ./out/aa4_6.out -s sub1.json NameC NameA
	diff ./out/aa4_6.out ./ref/aa4_6.out
	echo PASS

# test "else"
test5: pre_test
	./ifit -m test -i note.5 -o ./out/aa5_1.out -s sub1.json NameA NameB
	diff ./out/aa5_1.out ./ref/aa5_1.out
	./ifit -m test -i note.5 -o ./out/aa5_2.out -s sub1.json NameA 
	diff ./out/aa5_2.out ./ref/aa5_2.out
	./ifit -m test -i note.5 -o ./out/aa5_3.out -s sub1.json NameB
	diff ./out/aa5_3.out ./ref/aa5_3.out
	./ifit -m test -i note.5 -o ./out/aa5_4.out -s sub1.json NameC
	diff ./out/aa5_4.out ./ref/aa5_4.out
	./ifit -m test -i note.5 -o ./out/aa5_5.out -s sub1.json NameA NameC
	diff ./out/aa5_5.out ./ref/aa5_5.out
	./ifit -m test -i note.5 -o ./out/aa5_6.out -s sub1.json NameC NameA
	diff ./out/aa5_6.out ./ref/aa5_6.out
	echo PASS

# test include and set_path
test7: pre_test
	./ifit -m test -i inc.1 -o ./out/inc1.1.out -s sub1.json NameC NameA
	diff ./out/inc1.1.out ./ref/inc1.1.out
	./ifit -m test -i inc.2 -o ./out/inc2.1.out -s sub1.json NameC NameA
	diff ./out/inc2.1.out ./ref/inc2.1.out
	echo PASS

# Verify works with command line args A=BBB and that the command line args override the in file ones.
test8: pre_test
	./ifit -m test -i note.2 -o ./out/test8.out -s sub1.json NameA aa=AaAaAaA
	diff ./out/test8.out ./ref/test8.out
	echo PASS

# test of set_path and include
test9: pre_test
	./ifit -m test -i path.1 -o ./out/test9.out -s sub1.json NameA aa=AaAaAaA
	diff ./out/test9.out ./ref/test9.out
	echo PASS

# test define / undef
test10: pre_test
	./ifit -m test -i def-undef.1 -o ./out/test10.out -s sub1.json NameA aa=AaAaAaA
	diff ./out/test10.out ./ref/test10.out
	echo PASS

install: 
	go build
	cp ifit ~/bin

