#!/bin/sh

set -ev

/top-product -f testdata/db.json
/top-product -f testdata/db.csv

