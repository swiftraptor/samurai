# Samurai

Simple chat app

Client connects
Server sends message REQ_PUB_KEY
Client responds with
PUB_KEY
public rsa key encoded in base64
Server sends ACK_PUB_KEY

Session established
Client sends encrypted commands
Commands
JOIN <room name>
MESSAGE <room name> <message>
LIST