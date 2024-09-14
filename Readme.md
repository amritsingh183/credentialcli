# About the Credential utility

This CLI utility exposes the following commands which can be used to generate credentials like password(s)

## How to install

Simply run

```shell
go install ./cmd/credential
```

this will create a binary file named `credential` in your $GOPATH.
You may add $GOPATH to your $PATH to run the `credential` utility

## How to run this utility

To generate a password run the following command

```shell
~/go/bin/credential password
```

### Available options for generating passwords

- Length of the password(s) can be specified using `length=8` forexample
- To dump the password(s) to console set `output=0`
- To dump the password(s) to file set `output=1` and set filepath using `file=/your/filepath`
- To control the number of passwords use `count=20` for example
- If you want to include special characters in your passwords(recommended) use the `includeSpecialCharacters true` flag

some examples to run the utility( Assuming you have put `credential` in your $PATH)

- To generate a password with default options
    ```shell
    credential password
    ```

    where the default options are:
    - length=7
    - output=terminal
    - includeSpecialCharacters=true
    - urlSafe=False

- To generate a password with specified length and output to default file `./passwords.txt`
    ```shell
    credential password --length=40 --output=1 
    ```

- To generate 10 passwords with specified length and output to specific file `./myPassword.txt`
    ```shell
    credential password --length=40 --output=1 --count=10 --file=./myPassword.txt
    ```

- To generate 10 passwords with specified length and output to console
    ```shell
    credential password --length=40 --count=10
    ```

- To generate 10 URL safe passwords with specified length and output to console
    ```shell
    credential password --length=40 --count=10 --urlSafe=true
    ```
    
## How to run/test it without installing

Simply replace `~/go/bin/credential` with `go run cmd/credential/main.go` and you can try everything described above

