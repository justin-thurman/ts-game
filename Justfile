deploy:
  fly deploy --ha=false

run:
  go run .

local:
  telnet 127.0.0.1 8080

prod:
  telnet ts-game.fly.dev 8080
