# cmdproxy

Run any commands on remote host via HTTP/HTTPS.

## Usage

### Server

```sh
./cmdprx -host="0.0.0.0:12345" -secret="abc"
```

### Client

```sh
./cmdprx-cli -url="http://0.0.0.0:12345" -secret="abc" -cmd="echo hello world"
```
```
Error: 
ExitCode: 0
Stdout: hello world

Stderr: 

```

## Future work

### Security

- [ ] Client certificate authentication

### Feature

- [ ] Add stdin (but, no way to support input stream)
- [ ] Encode error (gob+Base64?)

### Usability

- [ ] Enhance logging
