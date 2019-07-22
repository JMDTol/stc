# Stellar transaction compiler

stc is a library and command-line tool for creating, editing, and
manipulating transactions for the [Stellar](https://www.stellar.org/)
blockchain.  It supports translating back and forth between
human-readable [txrep] format and Stellar's native binary XDR
representation, as well submitting transactions to the network,
querying account status, and more.  The library makes it easy to build
and submit transactions programmatically from go applications.

# Installing `stc` for non-developers

To install or upgrade this software if you don't plan to hack on it,
run the following two commands (assuming your GOPATH is in the default
location of `~/go`):

    rm -rf ~/go/src/github.com/xdrpp/stc
    go get github.com/xdrpp/stc/...

The `rm` command is necessary when upgrading because [some `go get`
limitation](https://github.com/golang/go/issues/27526) leaves your
tree in a detached state, so that `go get -u` cannot pull from the
remote `go1` branch.

Once this command completes, put the `~/go/bin` directory on your path
and you should be able to run `stc`.

To install the software from within a go module (a directory with a
`go.mod` file), you will need to specify the `go1` branch to get
autogenerated files.  Do so by running:

    go get github.com/xdrpp/stc/...@go1

# Using `stc`

See the [stc(1)][stc.1] man page for the command-line tool.

See the [godoc documentation][gh-pages] for the library.  Because
`stc` contains auto-generated source files that are not on the master
branch in git, it is best to view the documentation locally, rather
than through on-line godoc viewers that will not show you the correct
branch.  To view godoc, after installing `stc` with `go get` as
described above, start a local godoc server and open it in your
browser as follows:

    godoc -index -http localhost:6060 &
    xdg-open http://localhost:6060/pkg/github.com/xdrpp/stc

On MacOS computers, run `open` instead of `xdg-open`, or just paste
the URL into your browser.

# Building `stc` for developers

Because `stc` requires autogenerated files, the `master` branch is not
meant to be compiled under `$GOPATH`, but rather in a standalone
directory with `make`.

Furthermore, to build `stc` from the master branch, you also need to
have the [`goxdr`](https://github.com/xdrpp/goxdr) compiler.  Because
`stc` is codeveloped with goxdr, you may want to use a development
version of `goxdr`, which you can do by placing a the `goxdr` source
tree (or a symbolic link to it) in `cmd/goxdr`.

Once you have `goxdr`, you can build `stc` by running:

    make

To install `stc`, you will also need [pandoc](https://pandoc.org/) to
format the man page.

# Disclaimer

There is no warranty for the program, to the extent permitted by
applicable law.  Except when otherwise stated in writing the copyright
holders and/or other parties provide the program "as is" without
warranty of any kind, either expressed or implied, including, but not
limited to, the implied warranties of merchantability and fitness for
a particular purpose.  The entire risk as to the quality and
performance of the program is with you.  Should the program prove
defective, you assume the cost of all necessary servicing, repair or
correction.

In no event unless required by applicable law or agreed to in writing
will any copyright holder, or any other party who modifies and/or
conveys the program as permitted above, be liable to you for damages,
including any general, special, incidental or consequential damages
arising out of the use or inability to use the program (including but
not limited to loss of data or data being rendered inaccurate or
losses sustained by you or third parties or a failure of the program
to operate with any other programs), even if such holder or other
party has been advised of the possibility of such damages.

[gh-pages]: https://xdrpp.github.io/stc/pkg/github.com/xdrpp/stc/
[stc.1]: https://xdrpp.github.io/stc/pkg/github.com/xdrpp/stc/cmd/stc/stc.1.html
[txrep]: https://github.com/stellar/stellar-protocol/blob/master/ecosystem/sep-0011.md
