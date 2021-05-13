# Witchery

## Synopsis

Witchery is an experimental, embeddable programming language that employs concepts from Joy, Forth, Lisp, and Erlang, and also has some of its own unique features such as keyword restriction, stack traversal, and embedding of stacks.

## Specification

Witchery's code is encoded via UTF-8, always read from left to right, top to bottom, and only has three rules:

* Literals of data are always represented as strings, which are enclosed in square bracks: `[This is a string.]`. When strings terminate, the compiler embeds the string into an executable, which can then be executed once compilation is completed.
* Words can be written, which are composed entirely of glyphs that are not whitespace: `this-is-a-word This_is_also_a_WORD!!! $%^@@#$#$@#$%This too, is a word!@#@#$@#$^%&%^*`. Words are either keywords, which are built into the language, or definitions, which programmers can define. The convention for naming is to only use lowercase letters, represent spaces as dashes, prepend a singular dash to primary definitions of the program (`-this-is-an-example`), and represent definitions from a library as `library-name/name-of-definition`. Words can have effects at both the compilation and runtime stages.
* Whitespace separates words from strings. Interestingly, you can use an odd variant of postfix notation by doing something like `swap[foo][bar]`.

## Distribution

Witchery is released as a package written in [Go](https://golang.org) under the MIT License. Witchery is currently at version 0.1.0, and will likely undergo many various changes due to it being an experimental language.

## Examples

### Initialization

```
package main
import "yourmodule/witchery"
func main() {
	witchery.Initialize("[Hello, world] echo-line")
}
```

### Adding keywords

```
package main
import "yourmodule/witchery"
func main() {
	witchery.Keywords["foo"] = func(thread *witchery.Thread){
		second := witchery.DropStack(thread.Stack).(string)
		witchery.PushStack(witchery.DropStack(thread.Stack).(string) + second, thread.Stack)
	}
}
```

### Hello, world

```
[Hello, world!] echo-line
```

### Definitions

```
[-hello]
[[Hello, world!] echo-line] compile-now
define-now
-hello
```

### Arithmetic

```
[40] numberize [2] numberize add echo-line
[44] numberize [2] numberize subtract echo-line
[21] numberize [2] numberize multiply echo-line
[84] numberize [2] numberize divide echo-line
```

### Navigation

```
[foo] [bar] [baz] report
visit-previous report
visit-previous report
visit-previous report
[Uh oh, the ordering is messed up!] report
visit-next echo-line
visit-last echo-line
visit-first echo-line
```

### Conditionals

```
[t]
[[This should appear.] echo-line] compile-now
[[This should NOT appear.] echo-line] compile-now
if-else

[]
[[This should NOT appear.] echo-line] compile-now
[[This should appear.] echo-line] compile-now
if-else

[] not
[[This should appear.] echo-line] compile-now
if
```

### Loops

#### Eternal

```
[[t]] compile-now
[[The ride never ends.] echo-line] compile-now
while
```

#### Conditional

```
[-counter]
[
	[visit-previous duplicate visit-next duplicate visit-previous lesser?] compile-now
	[duplicate echo-line [1] numberize add visit-next] compile-now
	while visit-next drop drop
] compile-now define-now
[1] numberize [43] numberize -counter
```

### Vectors

```
[2] numberize create-vector
[0] numberize [foo] set-vector
[1] numberize [bar] set-vector

[3] numberize create-vector
concatenate-vectors

[0] numberize get-vector echo-line
[2] numberize get-vector echo-line
measure-vector echo-line
```

### Maps

```
create-map
[foo] [[Give me a foo!] echo-line] compile-now set-map
[bar] [Give me a bar!] set-map

[foo] get-map execute
[bar] get-map echo-line
measure-map echo-line
```

### Stacks

```
[foo] [bar] [baz] create-stack report
descend report ascend report

measure-stack echo-line
descend measure-stack echo-line
```

### Ceilings

#### Error

```
create-stack descend set-ceiling ascend
```

#### Success

```
create-stack descend set-ceiling reset-ceiling ascend
```

#### Error

```
[-ascend]
[reset-ceiling ascend] compile-now
define-now

create-stack descend set-ceiling -ascend
```

### Restrictions

#### Error

```
[echo-line] restrict [Hello, world!] echo-line
```

#### Success

```
[-echo-line]
[echo-line] compile-now
define-now

[echo-line] restrict [Hello, world!] -echo-line
```

#### Success

```
[echo-line] restrict reset-restrictions
[Hello, world!] echo-line
```

#### Error

```
[echo-line] restrict

[-echo-line]
[reset-restrictions echo-line] compile-now
define-now
```

### Executables

```
[[40] numberize] compile-now
[[2] numberize] compile-now
[add echo-line] compile-now
concatenate-executables concatenate-executables execute
```

### Files

Hint: run this script twice in a file called `script.witchery` via [Cauldron](cauldron.html).

```
[script.witchery] read-file echo-line

[script.witchery] [[Hello, world!] echo-line] write-file
```

### Input

```
[Enter your name: ] echo-raw
read-line
visit-previous [Your name is ] visit-next concatenate-strings [.] concatenate-strings echo-line
```

### REPL

```
[-quit?]
[[]] compile-now define-now

[-repl]
[
	create-stack descend set-ceiling
	[read-file] restrict
	[write-file] restrict
	[-quit? not] compile-now
	[
		[RUN [-quit?] [[t]] compile-now define-now TO QUIT.] echo-line
		[> ] echo-raw read-line compile-later execute
	] compile-now while
] compile-now define-now

-repl
```
