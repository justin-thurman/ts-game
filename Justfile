deploy:
  fly deploy --ha=false

run:
  go run .

local:
  telnet 127.0.0.1 8080

prod:
  telnet ts-game.fly.dev 8080

roomid:
  #!/bin/bash
  # outputs the current highest room ID value
  command -v yq >/dev/null 2>&1 || { echo >&2 "yq is required but it's not installed. Aborting."; exit 1; }

  max_id=0

  for file in ./room/roomdata/*.yaml; do
    file_max_id=$(yq eval '.rooms[].id' "$file" | sort -nr | head -n 1)
    if (( file_max_id > max_id )); then
      max_id=$file_max_id
    fi
  done

  echo "Largest ID value: $max_id"

test:
  go test -v ./engine
