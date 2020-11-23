CC = gcc
OUT = tcc
TEST_DIR = c_test

# go get -u golang.org/x/tools/...
# goyacc

go: clean flex
	goyacc -o fss.y.go -p "fss" fss.y

flex:
	flex --prefix=yy --header-file=fss.lex.h -o fss.lex.c fss.l

lex:
	lex -c --header-file=fss.lex.h -o fss.lex.c fss.l


test_yacc: go
	go test .

test_lex1: lex
	$(CC) $(TEST_DIR)/scanner.c fss.lex.c -o $(TEST_DIR)/scanner && $(TEST_DIR)/scanner < example/upper_example.fss && rm -f $(TEST_DIR)/scanner

test_lex2: lex
	$(CC) $(TEST_DIR)/scanner.c fss.lex.c -o $(TEST_DIR)/scanner && $(TEST_DIR)/scanner < example/20201123_string_bug.fss && rm -f $(TEST_DIR)/scanner

test: clean go
	go test .

clean:
	rm -f fss.y.go fss.lex.h fss.lex.c lex.yy.c