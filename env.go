package env

import (
	"bufio"
	"os"
	"strings"
	"sync"
)

const (
	DefaultFile = ".env"
	EmptyString = ""
)

var mutex = sync.RWMutex{}

// Get an environment variable, returns "" (empty string) if unset.
func Get(key string) (result string) {
	mutex.RLock()
	result = os.Getenv(key)
	mutex.RUnlock()
	return result
}

// Get a key or panic
func MustGet(key string) string {
	v := Get(key)
	if v == "" {
		panic("missing environment variable " + key)
	}

	return v
}

// Get an environment variable if it exists, otherwise return an alternate value.
func GetOr(key string, alt string) (result string) {
	if result = Get(key); result == EmptyString {
		result = alt
	}

	return
}

// Sets an environment variable unconditionally.
func Set(key, value string) (err error) {
	mutex.Lock()
	err = os.Setenv(key, value)
	mutex.Unlock()
	return err
}

// Sets an environment variable only if it is not already set.
func SetDefault(key, value string) (err error) {
	if !IsSet(key) {
		err = Set(key, value)
	}

	return
}

// Performs SetDefault on a map of key/value pairs.
func SetDefaults(vals map[string]string) (err error) {
	for k, v := range vals {
		if err = SetDefault(k, v); err != nil {
			return
		}
	}

	return
}

// Determine if an environment variable is set or not.
func IsSet(key string) bool {
	return (Get(key) != EmptyString)
}

// Unset an environment variable
func Unset(key string) error {
	return Set(key, EmptyString)
}

// Load the environment variables from the ".env" file.
func Load() error {
	return LoadFile(DefaultFile)
}

// Load environment variables from a given filename.
func LoadFile(name string) (err error) {
	var file *os.File
	var scanner *bufio.Scanner

	if file, err = os.Open(name); err != nil {
		return
	}
	defer file.Close()

	scanner = bufio.NewScanner(file)
	for scanner.Scan() {
		parts := strings.SplitN(scanner.Text(), "=", 2)
		if len(parts) == 2 {
			if err = Set(parts[0], parts[1]); err != nil {
				return
			}
		}
	}

	return
}
