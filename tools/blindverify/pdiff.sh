#!/bin/bash
# pdiff.sh OLDVER NEWVER [maxlines] — per-file protocol diff between two cached gophertunnel versions
GT="C:/Users/Пк/go/pkg/mod/github.com/sandertv/gophertunnel"
A="$GT@$1/minecraft/protocol"; B="$GT@$2/minecraft/protocol"; N=${3:-60}
cd "$A" || exit 1
files=$( (find . -type f -name '*.go'; cd "$B" && find . -type f -name '*.go') | sort -u | sed 's|^\./||')
for f in $files; do
  [ "$f" = "info.go" ] && continue
  if [ ! -f "$A/$f" ]; then echo "######## $f — NEW in $2 ########"; grep -E "^func \(pk \*|^type |^\t*ID" "$B/$f" | head -5; continue; fi
  if [ ! -f "$B/$f" ]; then echo "######## $f — REMOVED in $2 ########"; continue; fi
  d=$(diff "$A/$f" "$B/$f" | grep -vE '^[<>]\s*//|^---|^[0-9]' )
  [ -z "$d" ] && continue
  echo "######## $f ########"; echo "$d" | head -$N
done
