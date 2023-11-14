# mutual-exclusion

Handin 4 — Distributed Systems 2023

# Guide — How to run the program

The program uses ports that start and end on the same port to create a "ring".
No ports can be reused.
We suggest the ports: 8430, 8431, 8432

You must use port 8430 for exactly one of the clients/peers

From the peer folder, run:
```bash
go run peer.go -port 8430 -next 8431
```
From another terminal, run:
```bash
go run peer.go -port 8431 -next 8432
```
From a third terminal, run:
```bash
go run peer.go -port 8432 -next 8430
```
Each individual terminal will display a number representing logical time at which the critical section is accessed by the peer.

# Assumptions

- token ring, assume no crash failures
- take port of self and next in ring
