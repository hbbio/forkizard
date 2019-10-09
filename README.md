# forkizard

A wizard for comparing GitHub forks.

Will become useless once GitHub displays the number of ahead/behind commits on the Forks page!

## Problem

[StackOverflow: How to determine which forks on GitHub are ahead?](https://stackoverflow.com/questions/54868988/how-to-determine-which-forks-on-github-are-ahead)

## Solution

Since there does not seem to be an API for that, we scrape the info message of each GitHub fork.

## Usage

Build with `go build` and run `forkizard org/repo`:

```
$ ./forkizard k0kubun/pp
2019/06/20 16:27:35 44 forks
 44 / 44 [===============================================================================================================================] 100.00% 32s
done
/wtertius/pp       +6 -3
/hbbio/pp          +5 -3
/quanticko/pp      +2 -12
/federico-bollo/pp +1 -12
/juntaki/pp        +1 -12
/sniperkit/pp      +1 -12
/yudai/pp          +4 -39
/wangkechun/pp     +2 -57
```

## Building with Docker

To build without installing golang, run:

```
bin/build
```

## Bugs

Current version might not compare subforks well.

## About

Written in #Go by @hbbio and released under the MIT license.
