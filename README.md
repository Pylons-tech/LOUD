## Building

You need a correctly set up Golang environment over 1.11; this proejct uses `go mod`.

Run make.

    make

Then run `bin/loud` from this folder.

# Connecting to Play

## Overview

This LOUD (Legend of Undead Dragon) is a terminal-based SSH server. You need an ssh client installed and a private key generated.

## Connecting with macOS/Linux

```
    ssh localhost -p 2222
```

Or

```
    ssh username@localhost -p 2222
```

assuming you're running the loud server locally.

## Connecting with Windows

You'll need Putty and PuttyGen. [Follow the instructions here](https://system.cs.kuleuven.be//cs/system/security/ssh/setupkeys/putty-with-key.html) for how to make a key to connect.
