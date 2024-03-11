grep id a.out | sed 's/id:/id$/' | cut -d '$' -f 2 | cut -d ' ' -f 1 | sort -u | wc -l
