
# About the password utility

![GitHub CI](https://github.com/amritsingh183/credentialcli/actions/workflows/dev.workflow.yml/badge.svg)
![GitHub CI](https://github.com/amritsingh183/credentialcli/actions/workflows/release.workflow.yml/badge.svg)
![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/amritsingh183/credentialcli?sort=semver)


This CLI utility exposes the following commands which can be used to generate password(s).
It was build with love using <https://github.com/spf13/cobra-cli>

## How to install

<!-- FIXME: why do you care about Git tags? -->
<!-- [x] This tag is needed by Makefile to generate binary, let me remove this and just let makefile fetch the tags-->
<!-- ```shell
git fetch --all --tags
``` -->

### Method 1: Using released version of the binary

Download a released version of the binary from the releases page <https://github.com/amritsingh183/credentialcli/releases> and add the binary to your path. Assuming that you have renamed the downloaded binary as `password` and then added it to your path, you can

Generate a password using following command

```shell
password generate
```

### Method 2: Building from source

Assuming you have `go` installed on your system, clone this repository and then run

```shell
make install
```

this will create a binary file named `password` in your $GOBIN.
You may add $GOBIN to your $PATH to run the `password` utility

To build binaries for distribution, simply run

```shell
make binary
```

This will create binaries for Darwin-amd64, Linux-amd64, Linux-arm64 in `./out/bin` (same as method 1)

## How to run tests

To run the tests from the code, simply run the below command.
Make sure to remove the default password file `passwords.txt` , if it already exists.
For safety, the code will not remove an exising file, matching the default password-file name

```shell
make test
```

## How to run this utility

If you have followed the [Installation instructions](#how-to-install), you can run the utility as:

```shell
password generate
```

this will generate password with default options.

See the next section for available options

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

The next you're asked is to create a hierarchy of commands in this way:

<!-- [x] Done -->
- `password` is the master command
  <!-- [x] Done -->
  - `generate` is the actual command to generate the password
  <!-- [x] Not started yet -->
  - `validate` checks the password against a predefined set of rules and shares those with the user.

Everything else that has not been defined is up to you.
