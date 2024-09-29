# About the password utility

This CLI utility exposes the following commands which can be used to generate password(s).
It was build with love using <https://github.com/spf13/cobra-cli>

## How to install or create binaries

<!-- FIXME: why do you care about Git tags? -->
<!-- [x] This tag is needed by Makefile to generate binary, let me remove this and just let makefile fetch the tags-->
<!-- ```shell
git fetch --all --tags
``` -->

To install, simply run

```shell
make install
```

this will create a binary file named `password` in your $GOPATH.
You may add $GOPATH to your $PATH to run the `password` utility

To build binaries for distribution, simply run

```shell
make binary
```

This will create binaries for Darwin-linux, Linux-amd64, Linux-arm64 in `./out/bin`

## How to run tests

To run the tests from the code, simply run the below command.
Make sure to remove the default password file `passwords.txt` , if it already exists.
For safety, the code will not remove an exising file, matching the default password-file name

```shell
make test
```

## How to run this utility

To generate a password run the following command

```shell
~/go/bin/password generate
```

or if you have the `password` in your $PATH

```shell
password generate
```

### Available options for generating passwords

<!-- [x] Done adding examples -->

- Length of the password(s) can be specified using `length=8` forexample
- To dump the password(s) to console set `output=0`
- To dump the password(s) to file set `output=1` and set filepath using `file=/your/filepath`
- To control the number of passwords use `count=20` for example
- If you want to include special characters in your passwords(recommended) use the `includeSpecialCharacters true` flag

Some examples to run the utility( Assuming you have put `password` in your $PATH)

- To generate a password with default options

    ```shell
    password generate
    ```

  where the default options are:
  - length=7
  - output=terminal
  - includeSpecialCharacters=true
  - count=1

- To generate a password with specified length and output to default file `./passwords.txt`

    ```shell
    password generate --length=40 --output=1 
    ```

- To generate 10 passwords with specified length and output to specific file `./myPassword.txt`

    ```shell
    password generate --length=40 --output=1 --count=10 --file=./myPassword.txt
    ```

- To generate 10 passwords with specified length and output to console

    ```shell
    password generate --length=40 --count=10
    ```

### How to run/test it without installing
<!-- FIXME: this is outdated. -->
<!-- [x] updated it -->
Simply replace `~/go/bin/password` with `go run main.go` and you can try everything described above.

## New Requirements
<!-- [x] Not started yet -->
The next you're asked is to create a hierarchy of commands in this way:

- `password` is the master command
  - `generate` is the actual command to generate the password
  - `validate` checks the password against a predefined set of rules and shares those with the user.

Everything else that has not been defined is up to you.
