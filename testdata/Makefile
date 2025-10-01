..PHONY: all clean

all: test.cpio test.tar.gz

test.cpio:
	mkdir -p foo
	echo "baar" > foo/baar.txt
	echo "bazz" > foo/bazz
	find foo -print | cpio -o -H newc > test.cpio
	rm -rf foo

test.tar.gz:
	mkdir -p foo
	echo "baar" > foo/baar.txt
	echo "bazz" > foo/bazz
	tar -czf test.tar.gz foo
	rm -rf foo

clean:
	rm -rf foo test.cpio test.tar.gz
