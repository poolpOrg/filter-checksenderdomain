# filter-checksenderdomain

## Description
This filter performs a DNS lookup on the domain of the sending e-mail
address to determine if it exists before accepting it.


## Features
The filter currently supports:

- performs a Host lookup on the sender address


## Dependencies
The filter is written in Golang and doesn't have any dependencies beyond standard library.

It requires OpenSMTPD 6.6.0 or higher.


## How to install
Clone the repository, build and install the filter:
```
$ cd filter-checksenderdomain/
$ go build
$ doas install -m 0555 filter-checksenderdomain /usr/local/bin/filter-checksenderdomain
```


## How to configure
The filter itself requires no configuration.

It must be declared in smtpd.conf and attached to a listener:
```
filter "checksenderdomain" proc-exec "/usr/local/bin/filter-checksenderdomain"

listen on all filter "checksenderdomain"
```
