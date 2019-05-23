# usurp - Utterly Stupid Unprotected Reverse Proxy

I needed an ultra-lightweight reverse proxy for dumping some HTTP traffic. I'm putting it into a Git repo so I can keep track of it, 
but also so others can borrow ideas from it. It only does HTTP (not HTTPS) and it can only dump everything on a specified
port to a file. If that's all you need, this may be for you!

## Usage

```usurp -port <port to listen to> -target <host:port> -file <dump file>```
