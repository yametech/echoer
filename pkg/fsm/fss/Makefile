CC = gcc
OUT = tcc
TEST_DIR = c_test

flex:
	flex --prefix=yy --header-file=fss.lex.h -o fss.lex.c fss.l

lex:
	lex -c --header-file=fss.lex.h -o fss.lex.c fss.l

go: clean flex
	goyacc -o fss.y.go -p "fss" fss.y

test_yacc: go
	go test .

test_lex: lex
	$(CC) $(TEST_DIR)/scanner.c fss.lex.c -o $(TEST_DIR)/scanner && $(TEST_DIR)/scanner < example/upper_example.fss && rm -f $(TEST_DIR)/scanner

test: clean go
	go test .

clean:
	rm -f fss.y.go fss.lex.h fss.lex.c lex.yy.c