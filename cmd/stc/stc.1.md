% stc(1)
% David Mazi&egrave;res
%

# NAME

stc - Stellar transaction compiler

# SYNOPSIS

stc [-net=_id_] [-sign] [-c] [-l] [-u] [-i | -o FILE] _input-file_ \
stc -edit [-net=ID] _file_ \
stc -post [-net=ID] _input-file_ \
stc -preauth [-net=ID] _input-file_ \
stc -txhash [-net=ID] _input-file_ \
stc -qa [-net=ID] _accountID_ \
stc -qt [-net=ID] _txhash_ \
stc -qta [-net=ID] _accountID_ \
stc -fee-stats \
stc -ledger-header \
stc -create [-net=ID] _accountID_ \
stc -keygen [_name_] \
stc -pub [_name_] \
stc -import-key _name_ \
stc -export-key _name_ \
stc -list-keys \
stc -date YYYY-MM-DDThh:mm:ss[Z] \
stc -builtin-config

# DESCRIPTION

The Stellar transaction compiler, stc, is a command-line tool for
creating, viewing, editing, signing, and posting Stellar network
transactions.  It is intended for use by scripts or for creating test
transactions without the ambiguity of higher-layer wallet
abstractions.  stc is also useful in non-graphical environments, such
as a single-board computer used for cold storage.

The tool runs in one of several modes.  The default mode processes a
transaction in a single shot, optionally updating the sequence numbers
and fees, translating the transaction to/from human-readable form, or
signing it.  In edit mode, stc repeatedly invokes a text editor to
allow somewhat interactive editing of transactions.  In hash mode, stc
hashes a transactions to facilitate creation of pre-signed
transactions or lookup of transaction results.  Key management mode
allows one to maintain a set of signing keys.  Finally, network mode
allows one to post transactions or query the network for account and
fee status.

## Default mode

The default mode parses a transaction (in either textual or
base64-encoded binary), and then outputs it.  The input comes from a
file specified on the command line, or from standard input of the
argument is "`-`".  By default, stc outputs transactions in the
human-readable _txrep_ format, specified by SEP-0011.  With the `-c`
flag, stc outputs base64-encoded binary XDR format.  Various options
modify the transaction as it is being processed, notably `-sign`,
`-key` (which implies `-sign`), and `-u`.

Txrep format is automatically derived from the XDR specification of
`TransactionEnvelope`, with just a few special-cased types.  The
format is a series of lines of the form "`Field-Name: Value Comment`".
The field name is the XDR field name, or one of two pseudo-fields.
Pointers have a boolean pseudofield called `_present` that is true
when the pointer is non-null.  Variable-length arrays have an integer
pseudofield `len` specifying the array length.  There must be no space
between a field name and the colon.  After the colon comes the value
for that field.  Anything after the value is ignored.  stc sometimes
places a comment there, such as when an account ID has been configured
to have a comment (see the FILES section below).

Two field types have specially formatted values:

* Account IDs and Signers are expressed using Stellar's "strkey"
  format, which is a base32-encoded format where public keys start
  with "G", pre-auth transaction hashes start with "T", and hash-X
  signers start with "X".  (Private keys start with "S" in strkey
  format, but never appear in transactions.)

* Assets are formatted as _code_:_issuer_, where codes are formatted
  as printable ASCII bytes and two-byte hex escapes (e.g., `\x1f`),
  with no surrounding quotes.  A literal backslash or colon in an
  asset code must be escaped (e.g., `\\`).

Note that txrep is more likely to change than the base-64 XDR encoding
of transactions.  Hence, if you want to preserve transactions that you
can later read or re-use, compile them with `-c`.  XDR is also
compatible with other tools.  Notably, you can examine the contents of
an XDR transaction with `stellar-core` itself, using the command
"`stellar-core print-xdr --filetype tx --base64 FILE`", or by using
the web-based Stellar XDR viewer at
<https://www.stellar.org/laboratory/#xdr-viewer>.  You can also sign
XDR transactions with `stellar-core`, using "`stellar-core
sign-transaction --base64 --netid "Public Global Stellar Network ;
September 2015" FILE`".

## Edit mode

Edit mode is selected whenever stc is invoked with the `-edit` flag.
In this mode, whether the transaction is originally in base64 binary
or text, it is output in text format to a temporary file, and your
editor is repeatedly invoked to edit the file.  In this way, you can
change union discriminant values or array sizes, quit the editor, and
automatically re-enter the editor with any new fields appropriately
populated.

Note that for enum fields, if you add a question mark ("?") to the end
of the line, stc will populate the line with a comment containing all
possible values.  This is handy if you forget the various options to a
union discriminant such as the operation type.

Edit mode terminates when you quit the editor without modifying the
file, at which point stc writes the transaction back to the original
file.

## Hash mode

Stellar hashes transactions to a unique 32-byte value that depends on
the network identification string.  A transaction's hash, in hex
format, can be used to query horizon for the results of the
transaction after it executes.  With the option `-txhash`, stc hashes
transaction and outputs this hex value.

Stellar also allows an account to be configured to allow a
pre-authorized transaction to have a specific signing weight.  These
pre-authorized transactions use the same network-dependent hash values
as computed by `-txhash`.  However, to include such a hash as an
account signer, it must be encoded in strkey format starting with the
letter "T".  Running stc with the `-preauth` flag prints this
strkey-format hash to standard output.

Great care must be taken when creating a pre-authorized transaction,
as any mistake will cause the transaction not to run.  In particular,
make sure you have set the sequence number to one more than it will be
at the time you run the transaction, not one more than it is
currently.  (If the transaction adding the pre-authorized transaction
as a signer uses the same source account, it will consume a sequence
number.)  You should also make sure the transaction fee is high
enough.  You may wish to increase the fee above what is currently
required in case the fee has increased at the time you need to execute
the pre-authorized transaction.

Another potential source of error is that the pre-authorized
transaction hash depends on the network name, so make absolutely sure
the `-net` option is correct when using `-preauth`.

## Key management mode

stc runs in key management mode when one of the following flags is
selected:  `-keygen`, `-pub`, `-import-key`, `-export-key`, and
`-list-keys`.

These options take a key name.  If the key name contains a slash, it
refers to a file in the file system.  If the key name does not contain
a slash, it refers to a file name in the stc configuration directory
(see FILES below).  This allows keys to be stored in the configuration
directory and then accessed from any directory in which stc runs.

The `-keygen` and `-pub` options can be run with no key name, in which
case `-keygen` will output both the secret and public key to standard
output, and `-pub` will read a key from standard input or prompt for
one to be pasted into the terminal.

Keys are generally stored encrypted, but if you supply an empty
passphrase, they will be stored in plaintext.  If you use the
`-nopass` option, stc will never prompt for a passphrase and always
assume you do not encrypt your private keys.

## Network query mode

stc runs in network query mode when one of the `-post`, `-fee-stats`,
`ledger-header`, `-qa`, `-qt`, `-qta`, or `-create` options is
provided.

Post-mode, selected by `-post`, submits a transaction to the Stellar
network.  This is how you actually execute a transaction you have
properly formatted and signed.

`-fee-stats` reports on recent transaction fees.  `-ledger-header`
returns the latest ledger header.  `-qa` reports on the state of a
particular account.  `-qt` reports the result of a transaction that
has been previously submitted.  `-qta` reports transactions on an
account in reverse chronological order (use `-qt` to get more detail
on any transaction ID).  Unfortunately, some of these requests are
parsed from horizon responses in JSON rather than XDR format, and so
are reported in a somewhat incomparable style to txrep format.  For
example, balances are shown as a fixed-point number 10^7 times the
underlying int64.  `-create` creates and funds an account (which only
works when the test network is specified).

## Miscellaneous modes

The `-date` option parses a date and converts it to a Unix time.  This
is convenient for determining the Unix time to place in Timebounds.
The time can have one of several formats:

* `2006-01-02T15:04:05Z` (for parsing in UTC timezone)
* `2006-01-02T15:04:05-07:00` (for parsing in a specific timezone)
* `2006-01-02T15:04:05` (local time)
* `2006-01-02T15:04` (local time)
* `2006-01-02` (local time)

If no `stc.conf` configuration file exists, stc will use a built-in
one.  To see the contents of the built-in file, you can print it with
`-builtin-config`.

# OPTIONS

`-builtin-config`
:	Print the built-in system configuration file that is used if no
`stc.conf` file is found.

`-c`
:	Compile the output to base64 XDR binary.  Otherwise, the default
is to preserve the format (with `-i` and `-edit`) or output in text
mode to standard output or new files.  Only available in default mode.

`-create`
:	Create and fund an account on a network with a "friendbot" that
gives away coins.  Currently the stellar test network has such a bot
available by querying the `/friendbot?addr=ACCOUNT` path on horizon.

`-date`
:	Compute a Unix time from a human-readable time.

`-edit`
:	Select edit mode.

`-export-key`
:	Print a private key in strkey format to standard output.

`-fee-stats`
:	Dump fee stats from network

`-help`
:	Print usage information.

`-i`
:	Edit in place---overwrite the input file with the stc's output.
The original file is saved with a `~` appended to the name.  Only
available in default mode.

`-import-key`
:	Read a private key from the terminal (or standard input) and write
it (optionally encrypted) into a file (if the name has a slash) or
into the configuration directory.

`-key` _name_
:	Specifies the name of a key to sign with.  Implies the `-sign`
option.  Only available in default mode.

`-keygen` [_file_]
:	Creates a new public keypair.  With no argument, prints first the
secret then the public key to standard output.  When given an
argument, writes the public key to standard output and the private key
to a file, asking for a passphrase if you don't supply `-nopass`.
Note that if file contains a '/' character, the file is taken relative
to the current working directory or root directory.  If it does not,
the file is stored in stc's configuration directory.

`-l`
:	Learn all signers associated with an account.  Queries horizon and
stores the signers under the network's configuration directory, so
that it can verify signatures from all keys associated with the
account.  Only available in default mode.

`-list-keys`
:	List all private keys stored under the configuration directory.

`-net` _name_
:	Specify which network to use for hashing, signing, and posting
transactions, as well as for querying signers with the `-l` option.
Two pre-defined names are "main" and "test", but you can configure
other networks in `stc.conf` or by creating per-network configuration
files as discussed in the FILES section below.

`-nopass`
:	Never prompt for a passphrase, so assume an empty passphrase
anytime one is required.

`-o` _file_
:	Specify a file in which to write the output.  The default is to
send the transaction to standard output unless `-i` has been
supplied.  `-i` and `-o` are mutually exclusive, and can only be used
in default mode.

`-post`
:	Submit the transaction to the network.

`-preauth`
:	Hash a transaction to strkey for use as a pre-auth transaction
signer.  Beware that `-net` must be set correctly or the hash will be
incorrect, since the input to the hash function includes the network
ID as well as the transaction.

`-pub`
:	Print the public key corresponding to a particular private key.

`-qa`
:	Query the network for the state of a particular account.

`-qt`
:	Query the network for the results and effects of a particular
transaction.  The transaction must be specified in the hex format
output by `-txhash`.

`-qta`
:	Query the network for all transactions that have affected a
particular account, in reverse chronological order.  Also shows the
effects those transactions had on the target account.  To see effects
on all accounts, you can look up a particular transaction using `-qt`.

`-sign`
:	Sign the transaction.  If no `-key` option is specified, it will
prompt for the private key on the terminal (or read it from standard
input if standard input is not a terminal).

`-txhash`
:	Like `-preauth`, but outputs the hash in hex format.  Like
`-preauth`, also gives incorrect results if `-net` is not properly
specified.

`-u`
:	Query the network to update the fee and sequence number.  The fee
depends on the number of operations, so be sure to re-run this if you
change the number of transactions.  Only available in default mode.

`-v`
:	Produce more verbose output for the query options.

# EXAMPLES

`stc trans`
:	Reads a transaction from a file called `trans` and prints it to
standard output in human-readable form.

`stc -edit trans`
:	Run the editor on the text format of the transaction in file
`trans` (which can be either text or base64 XDR, or not exist yet in
which case it will be created in XDR format).  Keep editing the file
until the editor quits without making any changes.

`stc -c -i -key mykey trans`
:	Reads a transaction in file `trans`, signs it using key `mykey`,
then overwrite the `trans` file with the signed transaction in base64
format.  The original unsigned transaction is backed up in `trans~`.

`stc -post trans`
:	Posts a transaction in file `trans` to the network.  The
transaction must previously have been signed.

`stc -keygen`
:	Generate a new private/public key pair and print them both to
standard output, one per line (private key first).

`stc -keygen mykey`
:	Generate a new private/public key pair.  Prompt for a passphrase.
Print the public key to standard output.  Write the private key to
`$HOME/.config/stc/keys/mykey` encrypted with the passphrase.

# ENVIRONMENT

EDITOR
:	Name of editor to invoke with the `-edit` argument (default: `vi`)

STCDIR
:	Directory containing all the configuration files (default:
`$XDG_CONFIG_HOME/stc` or `$HOME/.config/stc`)

STCNET
:	Name of network to use by default if not overridden by `-net`
argument (default: `default`)

# FILES

Configuration files use the INI file format specified in the
git-config(1) manual page.  This also means you can use a command such
as `git config -f ~/.config/stc/stc.conf net.main.horizon https://...`
to edit the configuration files.  An example of this syntax is:

    [net]
    name = main
    network-id = "Public Global Stellar Network ; September 2015"
    horizon = https://horizon.stellar.org/
    native-asset = XLM

When using a network _NetName_, as specified by `$STCNET` or the
`-net` command-line argument, three configuration files are parsed in
order:

1. $STCDIR/_NetName_.net (or the default value of $STCDIR as described
   in the ENVIRONMENT section if $STCDIR is unset)

1. `$STCDIR/global.conf`

1. The system configuration, which comes from the first to exist of
   the following files:  `$STCDIR/stc.conf`, `/etc/stc.conf`, or
   `../share/stc.conf` relative to the executable.  If none of these
   files exist, stc uses the built-in version returned by the
   `-builtin-config` option.

A key is set to the first value encountered.  This means definitions
in the $STCDIR/_NetName_.net file take precedence over ones in the
`global.conf` file, which in turn has precedence over the global
configuration file.  However, it is possible to undefine a key by
including it without an equals sign, in which case it can be
redefined.  For example, the following would override any previously
set network-id:

    [net "main"]
    network-id
    network-id = "Public Global Stellar Network ; September 2015"

Subsections are only considered when the subsection string matches the
network name.  Hence, the section `[signers]` applies to all networks,
while `[signers "main"]` only applies to network main.  Generally the
$STCDIR/_NetName_.net file will include a `[net]` section, since it is
for only one network, while the global and system defaults will
include sections `[net "main"]` and `[net "test"]` for per-network
defaults.

The recognized keys are as follows:

`net.name`
:	Specifies the name of the network, which affects which subsections
will be parsed as described above.  This parameter can only be set in
a `[net]` section in the $STCDIR/_NetName_.net file, as it does not
make sense to set this globally.  Note that the value only changes
subsequently parsed sections; if the network name is changed,
previously parsed sections with the new name have already been ignored
and will not be reconsidered.

`net.network-id`
:	The network ID that permutes signatures and pre-signed-transaction
hashes (which prevents signatures from being valid on more than one
instantiation of the Stellar network).  If this is not specified, stc
automatically fetches and stores the network ID the first time it is
used.

`net.horizon`
:	The base URL of the horizon instance to use for this network.  You
may wish to change this URL to use your own local validator if you are
running one, or else that of an exchange that you trust.  Note that
the URL _must_ end with a `/` (slash) character.

`net.native-asset`
:	Shows how to render the native asset---e.g., `XLM` for the stellar
main network, and `TestXLM` for the stellar test network.  If not
specified, it defaults to the string `NATIVE`.  Note that this only
controls how the asset is rendered not parsed.  When parsing, any
string not ending ":IssuerAccountID" is considered the native asset.

accounts._AccountID_
:	Specifies a human-readable comment for _AccountID_ (which must be in
strkey format)

signers._SignerKey_
:	Specifies a human-readable comment for _SigherKey_ (in strkey
format)

# SEE ALSO

stellar-core(1), gpg(1), git-config(1)

The Stellar web site:  <https://www.stellar.org/>

Stellar's web-based XDR viewer:\
<https://www.stellar.org/laboratory/#xdr-viewer>

SEP-0011, the specification for txrep format:\
<https://github.com/stellar/stellar-protocol/blob/master/ecosystem/sep-0011.md>

RFC4506, the specification for XDR:\
<https://tools.ietf.org/html/rfc4506>

The XDR definition of a `TransactionEnvelope`:\
<https://github.com/stellar/stellar-core/blob/master/src/xdr/Stellar-transaction.x>

# BUGS

stc accepts and generates any `TransactionEnvelope` that is valid
according to the XDR specification.  However, a `TransactionEnvelope`
that is syntactically valid XDR may not be a valid Stellar
transaction.  stellar-core imposes additional restrictions on
transactions, such as prohibiting non-ASCII characters in certain
string fields.  This fact is important to keep in mind when using stc
to examine pre-signed transactions:  what looks like a valid, signed
transaction may not actually be valid.

stc uses a potentially imperfect heuristic to decide whether a file
contains a base64-encoded binary transaction or a textual one.

stc can only encrypt secret keys with symmetric encryption.  However,
the `-sign` option will read a key from standard input, so you can
always run `gpg -d keyfile.pgp | stc -sign -i txfile` to sign the
transaction in `txfile` with a public-key-encrypted signature key in
`keyfile.pgp`.

The options that interact with Horizon and parse JSON (such as `-qa`)
report things in a different style from the options that manipulate
XDR.

The txrep format has periodically been updated, and stc does not
attempt to maintain backwards compatibility with old files.  Binary
XDR, however, has been standard since 1995, so stc should be able to
parse any binary transaction since the launch of the Stellar network.
