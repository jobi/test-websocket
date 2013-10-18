# Experiments with websockets

This is a little experiment with using websockets as an IPC mechanism.

The little server (in Go) opens a websocket on an available port, then publishes
the service with the mDns responder.

The little client (in Cocoa) browses for these services, and connects to the websocket.
