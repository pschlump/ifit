
#
# Copyright (C) Philip Schlump, 2016.
#
# MIT Licnesed. 
#

DIFF=diff


all:
	go build



install: 
	go build
	( cd ~/bin ; rm -f ifit ; ln -s ../go/src/github.com/pschlump/ifit/ifit . )






.PHONY: test pre_test test0 test1 test2 test3 test4 test5 test6 test7 test8 test9 test10 test11 test12 test14 test15

test: test0 test1 test2 test3 test4 test5 test6 test7 test8 test9 test10 test11 test12 test14 test15

pre_test:
	go build
	mkdir -p ./ref ./out

# Test library
test0:
	( cd ./ifitlib ; go test )
	( cd ./stk ; go test )
	( cd ./fstk ; go test )

# basics
test1: pre_test
	echo test1
	./ifit -m test -i test/note.1 -o ./out/aa.out NameA
	$(DIFF) ./out/aa.out ./ref/aa.out
	./ifit -m test -i test/note.1 -o ./out/bb.out NameB
	$(DIFF) ./out/bb.out ./ref/bb.out
	./ifit -m test -i test/note.1 -o ./out/cc.out NameC
	$(DIFF) ./out/cc.out ./ref/cc.out
	./ifit -m test -i test/note.1 -o ./out/dd.out NameD
	$(DIFF) ./out/dd.out ./ref/dd.out
	./ifit -m test -i test/note.1 -o ./out/ee.out NameE
	$(DIFF) ./out/ee.out ./ref/ee.out
	./ifit -m test -i test/note.1 -o ./out/ab.out NameA NameB
	$(DIFF) ./out/ab.out ./ref/ab.out
	echo PASS

# variable substitution	
test2: pre_test
	echo test2
	./ifit -m test -i test/note.2 -o ./out/aa2.out -s sub1.json NameA NameB
	$(DIFF) ./out/aa2.out ./ref/aa2.out
	echo PASS

# more variable substitution	
test6: pre_test
	echo test6
	./ifit -m prod -i test/note.2 -o ./out/aa2_test6.out -s sub1.json NameA NameB
	$(DIFF) ./out/aa2_test6.out ./ref/aa2_test6.out
	echo PASS


# markers not in col(1)
test3: pre_test
	echo test3
	./ifit -m test -i test/note.3 -o ./out/aa3.out -s sub1.json NameA NameB
	$(DIFF) ./out/aa3.out ./ref/aa3.out
	echo PASS

# nested ifs
test4: pre_test
	echo test4
	./ifit -m test -i test/note.4 -o ./out/aa4_1.out -s sub1.json NameA NameB
	$(DIFF) ./out/aa4_1.out ./ref/aa4_1.out
	./ifit -m test -i test/note.4 -o ./out/aa4_2.out -s sub1.json NameA 
	$(DIFF) ./out/aa4_2.out ./ref/aa4_2.out
	./ifit -m test -i test/note.4 -o ./out/aa4_3.out -s sub1.json NameB
	$(DIFF) ./out/aa4_3.out ./ref/aa4_3.out
	./ifit -m test -i test/note.4 -o ./out/aa4_4.out -s sub1.json NameC
	$(DIFF) ./out/aa4_4.out ./ref/aa4_4.out
	./ifit -m test -i test/note.4 -o ./out/aa4_5.out -s sub1.json NameA NameC
	$(DIFF) ./out/aa4_5.out ./ref/aa4_5.out
	./ifit -m test -i test/note.4 -o ./out/aa4_6.out -s sub1.json NameC NameA
	$(DIFF) ./out/aa4_6.out ./ref/aa4_6.out
	echo PASS

# test "else"
test5: pre_test
	echo test5
	./ifit -m test -i test/note.5 -o ./out/aa5_1.out -s sub1.json NameA NameB
	$(DIFF) ./out/aa5_1.out ./ref/aa5_1.out
	./ifit -m test -i test/note.5 -o ./out/aa5_2.out -s sub1.json NameA 
	$(DIFF) ./out/aa5_2.out ./ref/aa5_2.out
	./ifit -m test -i test/note.5 -o ./out/aa5_3.out -s sub1.json NameB
	$(DIFF) ./out/aa5_3.out ./ref/aa5_3.out
	./ifit -m test -i test/note.5 -o ./out/aa5_4.out -s sub1.json NameC
	$(DIFF) ./out/aa5_4.out ./ref/aa5_4.out
	./ifit -m test -i test/note.5 -o ./out/aa5_5.out -s sub1.json NameA NameC
	$(DIFF) ./out/aa5_5.out ./ref/aa5_5.out
	./ifit -m test -i test/note.5 -o ./out/aa5_6.out -s sub1.json NameC NameA
	$(DIFF) ./out/aa5_6.out ./ref/aa5_6.out
	echo PASS

# test include and set_path
test7: pre_test
	echo test7
	./ifit -m test -i test/inc.1 -o ./out/inc1.1.out -s sub1.json NameC NameA
	$(DIFF) ./out/inc1.1.out ./ref/inc1.1.out
	./ifit -m test -i test/inc.2 -o ./out/inc2.1.out -s sub1.json NameC NameA
	$(DIFF) ./out/inc2.1.out ./ref/inc2.1.out
	echo PASS

# Verify works with command line args A=BBB and that the command line args override the in file ones.
test8: pre_test
	echo test8
	./ifit -m test -i test/note.2 -o ./out/test8.out -s sub1.json NameA aa=AaAaAaA
	$(DIFF) ./out/test8.out ./ref/test8.out
	echo PASS

# test of set_path and include
test9: pre_test
	echo test9
	./ifit -m test -i test/path.1 -o ./out/test9.out -s sub1.json NameA aa=AaAaAaA
	$(DIFF) ./out/test9.out ./ref/test9.out
	echo PASS

# test define / undef
test10: pre_test
	echo test10
	./ifit -m test -i test/def-undef.1 -o ./out/test10.out -s sub1.json NameA aa=AaAaAaA
	$(DIFF) ./out/test10.out ./ref/test10.out
	echo PASS

# Simple file include test
test11: pre_test
	echo test11
	./ifit -m test -i test/inc.3 -o ./out/test11.out -s sub1.json NameA aa=AaAaAaA
	$(DIFF) ./out/test11.out ./ref/test11.out
	echo PASS

test12:
	echo test12
	./ifit -m test -i test/inc.4 -o ./out/test12.out -s sub1.json NameA aa=AaAaAaA
	$(DIFF) ./out/test12.out ./ref/test12.out
	echo PASS

# new test for $base$
test14:
	echo test14 New test for base
	./ifit -m dev -i test/inc.14 -o ./out/test14.out -s sub14.json 
	$(DIFF) ./out/test14.out ./ref/test14.out
	./ifit -m test -i test/inc.14 -o ./out/test14a.out -s sub14.json 
	$(DIFF) ./out/test14a.out ./ref/test14a.out

test15:
	echo PASS

