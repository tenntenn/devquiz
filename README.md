# Devquiz

devquiz validator for golang.tokyo.

## Instal

```sh
$ go get github.com/tenntenn/devquiz
```

## Usage

```sh
$ nkf --utf8 event_xxxxx_participants.csv > golangtokyoxxx.csv
$ devquiz -urlrow 10 golangtokyoxxx.csv > golangtokyoxxx.txt
$ cat golangtokyoxxx.txt | grep -v skip | grep -v ok
```
