#!/bin/sh

set -xe

make build

./shortest-path < testdata/test-01.txt > result.txt
cmp result.txt testdata/result-01.txt

./shortest-path < testdata/test-02.txt > result.txt
cmp result.txt testdata/result-02.txt

./shortest-path < testdata/test-03.txt > result.txt
cmp result.txt testdata/result-03.txt

./shortest-path < testdata/test-04.txt > result.txt
cmp result.txt testdata/result-04.txt

./shortest-path < testdata/test-05.txt > result.txt
cmp result.txt testdata/result-05.txt

./shortest-path < testdata/test-06.txt > result.txt
cmp result.txt testdata/result-06.txt

rm -f shortest-path result.txt

echo "PASS"
