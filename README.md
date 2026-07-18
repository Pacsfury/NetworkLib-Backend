# NetworkLib - Backend
The core server routing infrastructure for the NetworkLib ecosystem.

<sub>Reference version: A0.1.0 (use this for matching the client)</sub>

---

## Overview
This backend processes incoming client messages, manages state, and routes data back to connected nodes based on the requested commands.

## Purpose
NetworkLib provides an ultra-simple, direct, and lightweight way to build multiplayer games or networked software using a single, unified file right in your directory.

## Current Commands (Alpha)
New commands and functionalities are under active development. The current alpha version supports:

* `SET <name> <value>`: Initializes or updates a variable state.
* `GET <name>`: Retrieves and returns the specified variable's value.
* `CONST <name> <value>`: Declares an immutable, read-only constant.
* `TEMP <name> <value>`: Creates a temporary variable that self-destructs after the first `GET` request.
* `SIGNAL <value>`: Broadcasts an immediate event notification to all connected clients.
* `SUB <name>`: Subscribes a client to receive real-time notifications whenever the target variable changes.

_When sent, and received, are transformed to opcodes for reducing bandwidth (coming soon with std signals)_

## Repository Structure
```
Main/
| README.md
| go.mod
| client.go      => basic test for the server, base for real clients
| client2.go     => another client for testing basic interaction
| defs.go        => global variables definition
| main.go        => main file, calls and initialitzations
| manager.go     => receives and sends values
```

## Current Supported Frontend Languages
_Find this and more at [github.com/Pacsfury/NetworkLib-Clients](github.com/Pacsfury/NetworkLib-Clients)_
- **Lua**: Lua is the first language receiving support

## Roadmap & Future Features

Stability and optimitzations.

### Understanding Signals
A **Signal** is an asynchronous message pushed instantly to the client, triggering a predefined callback or action upon arrival. Signals can optionally carry arguments.

#### Example Scenario:
* **Setup**: Client B listens for the signal `#subscribed_var_changed` with the argument `pos_x`. If received, Client B triggers a `GET` command to sync the player's X position.

**Execution Flow:**
1. **Client B** sends: `SUB pos_x`
2. **Client A** sends: `SET pos_x 10`
3. **Server** broadcasts to Client B: `#subscribed_var_changed pos_x 10`
4. **Client B** automatically updates the game state.
