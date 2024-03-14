package main

import (
	"net"
)

// GetLocalIP secured local ip address
func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

// Filter filters a slice with a predicate
func Filter[V any](collection []V, predicate func(item V, index int) bool) []V {
	result := make([]V, 0, len(collection))

	for i, item := range collection {
		if predicate(item, i) {
			result = append(result, item)
		}
	}

	return result
}

// EmptyOr returns fallback if v is considered empty, otherwise returns v.
func EmptyOr[T comparable](v T, fallback T) T {
	var zero T
	if zero == v {
		return fallback
	}
	return v
}

// Try0 has the same behavior as Try, but callback returns no variable.
// Try0 is used to capture panics that may occur in goroutine to ，
// ensure that the entire program does not crash due to panic.
// Try0 用于捕获 goroutine 中可能发生的 panic，确保整个程序不因 panic 而崩溃。
func Try0(callback func()) bool {
	return Try(func() error {
		callback()
		return nil
	})
}

// Try calls the function and return false in case of error.
func Try(callback func() error) (ok bool) {
	ok = true

	defer func() {
		if r := recover(); r != nil {
			ok = false
		}
	}()

	err := callback()
	if err != nil {
		ok = false
	}

	return
}
