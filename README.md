# Collatz Conjecture Counter Example Miner

If there exists a counter example to the [Collatz conjecture](https://en.wikipedia.org/wiki/Collatz_conjecture), this program can help you find it. Eventually.

## Disclaimer

This project is a semi joke so don't take it too seriously.

## Why

There are already many smarter and better equipped people trying to brute force a counter example to the Collatz conjecture. As reported in [this paper](https://rdcu.be/b5nn1), the Collatz conjecture was already verified up to 2^68 in 2020 using extremely advanced optimizations.

Most people wouldn't be able to remotely replicate the performance let alone beat it, so the motivation for this project is to take the least-effort approach in brute forcing a counter example: by randomly testing some large number that hasn't been verified. To me, this seems like the only way to have a nonzero chance of discovering a counter example if it exists (as opposed to brute forcing numbers that another computer is brute forcing doing faster).

## Installation

Choose one of two main ways to install this program:

1. Download the executable from [Releases](https://github.com/ayoung19/collatz-miner/releases) and run it.

2. Install Go, clone this repository, navigate to the root of the project, and run `go build`. This should create the `collatz-miner` executable.
