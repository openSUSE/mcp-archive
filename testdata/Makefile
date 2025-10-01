.PHONY: all clean

all: test.cpio

test.cpio: foo/baar.txt foo/bazz
	find foo -print | cpio -o -H newc > test.cpio
	rm -rf foo

foo/baar.txt:
	mkdir -p foo
	echo "baar" > foo/baar.txt

foo/bazz:
	mkdir -p foo
	echo "bazz" > foo/bazz

clean:
	rm -rf foo test.cpio
