# usurp - Utterly Stupid Unprotected Reverse Proxy

I needed an ultra-lightweight reverse proxy for dumping some HTTP traffic. I'm putting it into a Git repo so I can keep track of it, 
but also so others can borrow ideas from it. It only does HTTP (not HTTPS) and it can only dump everything on a specified
port to a file. If that's all you need, this may be for you!

## Usage
```
Usage of ./zreplay:
  -file string
    	Input file captured with usurp
  -forever
    	Run realtime with delays according to original timing
  -gz
    	Input is gzipped
  -log.format value
    	Set the log target and format. Example: "logger:syslog?appname=bob&local=7" or "logger:stdout?json=true" (default "logger:stderr")
  -log.level value
    	Only log messages with the given severity or above. Valid levels: [debug, info, warn, error, fatal] (default "info")
  -namespace string
    	Namespace substitution (default "default")
  -offset duration
    	Time offset in minutes
  -realtime
    	Run realtime with delays according to original timing
  -target string
    	The target URL
```

## Example
```usurp -port <port to listen to> -target <host:port> -file <dump file>```


