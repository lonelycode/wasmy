# Wasmy
A utility wrapper to make WASM plugins less painful

## Context
Wasm is great, but passing values between host <-> guest is uneccessarily complicated and not well documented across the various runtimes. 

The reason it is hard is because WASM only allows for integer and float types to be shared between guest and host. That means no strings, no structured objects or arrays.

WASM does provide a shared linear memory block though, which is what should be used to move data into and out of WASM modules. Whgile using it for trivial examples is fine, using it in anger starts getting very messy and in no-way general purpose.

There is a standard WASI - Web Assembly Shared Interfaces - that might one day make this a great deal easier by enabling socket-based i/o between guest and host, but it's very early days. If you want to do anything complex with WASM (e.g. making an HTTP request), you need to import that capability from your host langauge. 

## What does Wasmy do?
Wasmy provides wrapper functions and a shared mempry manager that handles function i/o for you, giving you clean interface{} types on either end of the toolchain. When using Wasmy, you do not need to worry about how to get data into and out of WASM, because it's handled for you. 

The goal is to make WASM more accessible for us average folks that just want to JFDI.

## Usage

See the `wasm-tests/managedv2.go` file for an example of how to write WASM functions that can be exported in go. to compile the wasm file you'll need TinyGo:

```
tinygo build -o wasm-tests/managedv2.wasm -wasm-abi=generic  -target=wasi wasm-tests/managedv2.go
```

To see how to call a function in this wasm file, you'll also need to compile a runner, see `example/example.go` for a sample application. 

Compiling it is as simple as runnign `go build`.

When you ruin the example app, you should see output like the below:

```
vmuser@codeserv:~/wasmy/example$ ./example 
inside module: hello martin
From Host: Hello Mr. anderson
host function output from inside module: From Host: Hello Mr. anderson
function output (from runner): hello martin 
```

## Warnings and Caveats

- This is an experimental library and has not been used in anger
- The i/o buffers for guest and host functions are currently set to 1024 bytes, this can be increased by modiying `FUNCBUFFER_SIZE`. WASM requires fixed array sizes so a byte slice isn't suitable.
- There's no overflow checks in place yet, so use at your peril
