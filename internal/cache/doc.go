/*
Package cache provides a generic, size-limited cache for storing key-value pairs efficiently.

The `cache` package offers a simple and performant caching mechanism using Go's generics. It allows
developers to store and retrieve values of any comparable type with automatic eviction of the oldest
entries when the cache reaches its maximum capacity. This ensures optimal memory usage and prevents
the cache from growing indefinitely.
*/
package cache
