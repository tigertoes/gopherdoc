# gopherdoc
I'm really surprised I haven't found anyone else attempt this yet. Gopher is an
old protocol that is an incredibly simple TCP service that serves primarily
ASCII. It features linking and a simple document structure and that's about it.

# WARNING! ACHTUNG!
This is more than likely to be insecure code, given it's an un-authenticated,
insecure protocol TCP server that is making system calls.

## Using it

`go run gopherdoc.go`

Then find yourself a nice gopher client (I prefer lynx myself) and open
[gopher://localhost:7000/1/buf/ScanLines](gopher://localhost:7000/1/buf/ScanLines)
and browse around!

# License
MIT.
