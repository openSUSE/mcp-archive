.PHONY: all clean

all: test.cpio test.tar.gz test.tar.bz2 test.tar.xz test.zip

test.cpio:
	mkdir -p foo
	echo "das Pferd isst Gurkensalat" > foo/baar.txt
	echo "bazz" > foo/bazz
	find foo -print | cpio -o -H newc > test.cpio
	rm -rf foo

test.tar.gz:
	mkdir -p foo
	echo "das Pferd isst Gurkensalat" > foo/baar.txt
	echo "bazz" > foo/bazz
	tar -czf test.tar.gz foo
	rm -rf foo

test.tar.bz2:
	mkdir -p foo
	echo "das Pferd isst Gurkensalat" > foo/baar.txt
	echo "bazz" > foo/bazz
	tar -cjf test.tar.bz2 foo
	rm -rf foo

test.tar.xz:
	mkdir -p foo
	echo "das Pferd isst Gurkensalat" > foo/baar.txt
	echo "bazz" > foo/bazz
	tar -cJf test.tar.xz foo
	rm -rf foo

test.zip:
	mkdir -p foo
	echo "das Pferd isst Gurkensalat" > foo/baar.txt
	echo "bazz" > foo/bazz
	zip -r test.zip foo
	rm -rf foo

clean:
	rm -rf foo test.cpio test.tar.gz test.tar.bz2 test.tar.xz test.zip