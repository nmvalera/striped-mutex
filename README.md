# Striped-Mutex

Implementation of a Striped-Mutex inspired from Guava Java lib Guava: https://github.com/google/guava/wiki/StripedExplained

Striped-Mutex allows to fined grained locking of multiple distinct objects by basing on a indexation key.

```go
package main

import (
    stripedmutex "github.com/nmvalera/striped-mutex"
)

func main() {
    // Create a striped mutex with 20 locks
    smux := stripedmutex.New(20)

    smux.Lock("key")
    defer smux.Lock("key")

    // ... do something in a concurrent manner in scope associated to "key"
}
```

# Implementation

A straightforward implementation would be to create a lock for every key. While this approach leads to minimal lock contention, it results in a linear memory usage wrt the count of keys which can be unsastifying when dealing with a large number of keys.

Striped-Mutex allows to configure a number of locks that are distributed between keys based on their hash code. It allows to select a tradeoff between concurrency and memory consumption, while retaining the property that if `key1 == key2` then lock associated to `key1` and `key2` is the same.
