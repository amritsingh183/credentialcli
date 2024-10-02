
# About the password utility

[![main-branch](https://github.com/amritsingh183/credentialcli/actions/workflows/main.workflow.yml/badge.svg)](https://github.com/amritsingh183/credentialcli/actions/workflows/main.workflow.yml)
[![dev-branch](https://github.com/amritsingh183/credentialcli/actions/workflows/dev.workflow.yml/badge.svg)](https://github.com/amritsingh183/credentialcli/actions/workflows/dev.workflow.yml)
[![release-branch](https://github.com/amritsingh183/credentialcli/actions/workflows/release.workflow.yml/badge.svg)](https://github.com/amritsingh183/credentialcli/actions/workflows/release.workflow.yml)
[![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/amritsingh183/credentialcli?sort=semver)](https://img.shields.io/github/v/release/amritsingh183/credentialcli)

This CLI utility is built for generating password(s).
It was made with love using <https://github.com/spf13/cobra-cli>

## How to install

<!-- FIXME: why do you care about Git tags? -->
<!-- [x] This tag is needed by Makefile to generate binary, let me remove this and just let makefile fetch the tags-->
<!-- ```shell
git fetch --all --tags
``` -->
There are two ways to install the `password` utility

### Method 1: Download a pre-built binary

You can download the pre-built exectable binary for your operating system from the releases page <https://github.com/amritsingh183/credentialcli/releases> and add the binary to your programs path `/usr/local/bin`.

```shell
wget https://github.com/amritsingh183/credentialcli/releases/download/v0.0.2/password-0.0.2-linux-amd64

mv password-0.0.2-linux-amd64 /usr/local/bin/password

```

```shell
password -v # should show you the current version
```

Now, you can generate a password using following command

```shell
password generate
```

this will generate password with default options. All options, inluding the default ones are defined in the section [Default options](#available-options-for-controlling-the-generated-password)

### Method 2: Building from source

Assuming you have `go` installed on your system, clone this repository and then run

```shell
make install
```

this will create a binary file named `password` in your $GOBIN.
You may add $GOBIN to your $PATH to run the `password` utility

Another option is to build binaries for distribution by running

```shell
make binary
```

This will create binaries for Darwin-amd64, Linux-amd64, Linux-arm64 in `./out/bin`.
Rest of the steps are same as [Method 1](#method-1-download-a-pre-built-binary)

## Available options for controlling the generated password

The default options are:

- length=7
- output=terminal
- includeSpecialCharacters=true
- count=1

The purpose of each option is as follows:

- Length of the password(s) can be specified using `length=8` forexample
- To dump the password(s) to console set `output=0`
- To dump the password(s) to file set `output=1` and set filepath using `file=/your/filepath`
- To control the number of passwords use `count=20` for example
- If you want to include special characters in your passwords(recommended) use the `includeSpecialCharacters true` flag

## How to use this utility

If you have followed the [Installation instructions](#how-to-install), you can use the following examples. You can check [All available options](#available-options-for-controlling-the-generated-password).

<!-- [x] Done adding examples -->

- To generate a password with default options

    ```shell
    password generate
    ```

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

## How to run tests for the source code

To run the tests from the code, simply run the below command.
Make sure to remove the default password file `passwords.txt` , if it already exists.
For safety, the code will not remove an exising file, matching the default password-file name

```shell
make test
```

<!-- ### How to test it without installing -->
<!-- FIXME: this is outdated. -->
<!-- [x] removing this section-->
<!-- Simply replace `~/go/bin/password` with `go run main.go` and you can try everything described above. -->

## New Requirements

The next you're asked is to create a hierarchy of commands in this way:

<!-- [x] Done -->
- `password` is the master command
  <!-- [x] Done -->
  - `generate` is the actual command to generate the password
  <!-- [x] Not started yet -->
  - `validate` checks the password against a predefined set of rules and shares those with the user.

Everything else that has not been defined is up to you.
