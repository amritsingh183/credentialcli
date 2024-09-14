# About the Credential utility

<!-- TODO: add a section on how to run tests. -->
<!-- You can have three h2 titles: Install, Run, and Test. Then, you can narrow them down
by explaining in a more specific way. -->

This CLI utility exposes the following commands which can be used to generate credentials like password(s)

## How to install

Simply run

```shell
go install ./cmd/credential

```

this will create a binary file named `credential` in your $GOPATH

## How to run this utility

To generate a password run the following command

```shell
~/go/bin/credential password

```

### Available options for generating passwords

<!-- TODO: expand this section and put some run examples. -->

- Length of the password(s) can be specified using `length=8` forexample
- To dump the password(s) to console set `output=0`
- To dump the password(s) to file set `output=1` and set filepath using `file=/your/filepath`
- To control the number of passwords use `count=20` for example
- If you want to include special characters in your passwords(recommended) use the `includeSpecialCharacters true` flag

## How to run/test it without installing

Simply replace `~/go/bin/credential` with `go run cmd/credential/main.go` and you can try everything described above
