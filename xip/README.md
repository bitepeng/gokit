# GoKit - xip

IP kits for Golang development.

## Installation

    go get -u github.com/likexian/gokit

## Importing

    import (
        "github.com/likexian/gokit/xip"
    )

## Documentation

Visit the docs on [GoDoc](https://godoc.org/github.com/likexian/gokit/xip)

## Example

### Check string is a valid ip

    ok := xip.IsIP("1.1.1.1")
    fmt.Println("1.1.1.1 is a ip:", ok)

### IPv4 ip2long

    i, err := IPv4ToLong("1.1.1.1")
    if err == nil {
        fmt.Println("1.1.1.1 ip2long is:", i)
    }

### IPv4 long2ip

    ip := LongToIPv4(16843009)
    fmt.Println("16843009 long2ip is:", ip)

## LICENSE

Copyright 2019, Li Kexian

Apache License, Version 2.0

## About

- [Li Kexian](https://www.likexian.com/)

## DONATE

- [Help me make perfect](https://www.likexian.com/donate/)