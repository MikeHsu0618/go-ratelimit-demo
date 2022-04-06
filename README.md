Rate-limit Demo
---

基於呈現觀念的 Demo 範例，簡略了結構化，著重於實現 IP 限制之功能

## Run Project (docker-compose)

### Run Redis Server
```
docker compose up -d
```

### Run Application
```
cd ./ratelimit-demo

go run main.go
```

### Introduce
使用 Lua 腳本讓 `Incr` 和 `expire` 兩條命令保持原子性，
實際邏輯為以用戶 `IP` 為 key 值, 在 Lua 腳本中使用 `Incr` 指令, 如返回值為初始的 `1`，
則給予 `expire` 過期時間，最後返回次數判斷是否超過需要回傳 `429 too many reqeusts`。

### Demo
 <img src="./ratelimit.gif" width="100%"/>
