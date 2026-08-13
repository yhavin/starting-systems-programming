#!/usr/bin/env bash
echo "to the last I grappe with thee; from hell's heart I stab at thee; for hate's sake I spit my last breath at thee" > moby.txt

shexdump moby.txt > moby.hex
unhexdump moby.hex > moby2.txt
diff -s moby.txt moby2.txt