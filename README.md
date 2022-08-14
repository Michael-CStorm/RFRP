# Remotely controlled frp

## Sample usage

### Step 1 : Host a temporary file server at socket 9999.

```
python3 -m http.server 9999
```

### Step 2 : Start the docker.

```
~/rfrp $ docker compose build && docker compose up
```

### Step 3 : Start the rfrp client.

Download the frp source code at
https://github.com/Darwin-Che/libfrp

```
~/libfrp $ make
```

Create a file of `libfrp/bin/test.ini` with the following content:
```
[common]
server_addr = localhost
server_port = 9000

[web]
type = http
local_port = 9999
subdomain = test
```

Start the client
```
~/libfrp/bin $ ./frpc -c test.ini
```

### Step 4 : Toggle the router for `test.localhost`.

Send the following POST request to `localhost:9010/router`

```
{
  "operation": "enable",
  "subdomain": "test.localhost"
}
```
to enabled the subdomain.

```
{
  "operation": "enable",
  "subdomain": "test.localhost"
}
```
to disable the subdomain.

The result can be accessed at `test.localhost:9090` from the browser.



