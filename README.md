## Blockchain in Golang

This is an up-to-date version of Tensor-Programming's Blockchain-Golang found [here](https://github.com/tensor-programming/golang-blockchain).

Since the original was written using an outdated version of Go and just doesn't work using the latest version, I decided to have it rewritten (and cleaned up some) using `Go 1.24.1`.

The project is split into different branches so anyone can follow along his walkthrough which can be found [here](https://www.youtube.com/playlist?list=PLJbE2Yu2zumC5QE39TQHBLYJDB2gfFE5Q).


At around part5, as an alternative, I decided to use `Protobuf` as the serialization method for creating wallets. Branches with `_b` suffix follows this branch from here on while the original implementation using `gob` has been left intact in branches without the suffix.

I also refactored some functions to clean it up a bit. The refactored functions can be found inside the [`utils` folder](./utils/bcutils.go)

Enjoy.

