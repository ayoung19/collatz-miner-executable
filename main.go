package main

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"strconv"

	"github.com/google/uuid"
	"github.com/tcnksm/go-input"
)

var (
	bigZero = big.NewInt(0)
	bigOne  = big.NewInt(1)
	bigTwo  = big.NewInt(2)
	bigFour = big.NewInt(4)
)

type client struct {
	uuid            string
	name            string
	max             *big.Int
	goroutinesCount uint64

	logger *log.Logger
}

func newClient(logger *log.Logger) (*client, error) {
	ui := &input.UI{
		Writer: os.Stdout,
		Reader: os.Stdin,
	}

	clientUUID := uuid.New().String()
	clientName, err := ui.Ask("collatz-miner instance name", &input.Options{
		Default: clientUUID,
	})
	if err != nil {
		return nil, err
	}

	nthPowerOfTwoString, err := ui.Ask("random integer upper bound (nth power of 2)", &input.Options{
		Default: "10000",
	})
	if err != nil {
		return nil, err
	}
	nthPowerOfTwo, err := strconv.ParseUint(nthPowerOfTwoString, 10, 64)
	if err != nil {
		return nil, err
	}
	clientMax := big.NewInt(0).Lsh(big.NewInt(2), uint(nthPowerOfTwo))

	goroutinesCountString, err := ui.Ask("number of goroutines", &input.Options{
		Default: "1",
	})
	if err != nil {
		return nil, err
	}
	clientGoroutinesCount, err := strconv.ParseUint(goroutinesCountString, 10, 64)
	if err != nil {
		return nil, err
	}

	logFilePath, err := ui.Ask("log file path", &input.Options{
		Default: filepath.Join(".", "collatz-miner.logs"),
	})
	if err != nil {
		return nil, err
	}
	f, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	logger.SetOutput(io.MultiWriter(os.Stdout, f))

	return &client{
		uuid:            clientUUID,
		name:            clientName,
		max:             clientMax,
		goroutinesCount: clientGoroutinesCount,

		logger: logger,
	}, nil
}

func mine(i uint64, c *client) {
	for {
		initial, err := rand.Int(rand.Reader, c.max)
		if err != nil {
			c.logger.Printf("an unexpected error occured while generating a random big int: %v", err)
			os.Exit(1)
		}

		current := big.NewInt(0).Set(initial)
		seen := map[string]bool{}

		for {
			currentB64 := base64.StdEncoding.EncodeToString(current.Bytes())
			currentSeen := seen[currentB64]

			if currentSeen {
				break
			}

			seen[currentB64] = true

			if current.Bit(0) == 0 {
				current.Rsh(current, 1)
			} else {
				current.Mul(current, big.NewInt(3))
				current.Add(current, big.NewInt(1))
			}
		}

		if current.Cmp(bigZero) == 0 || current.Cmp(bigOne) == 0 || current.Cmp(bigTwo) == 0 || current.Cmp(bigFour) == 0 {
			continue
		}

		c.logger.Printf("counter example found. initial value: %s, cycle start: %s, path length: %d", initial.String(), current.String(), len(seen))
		os.Exit(0)
	}
}

func main() {
	logger := log.Default()

	c, err := newClient(logger)
	if err != nil {
		logger.Printf("an unexpected error occured while initializing the client: %v", err)
		os.Exit(1)
	}

	c.logger.Printf("successfully initialized the client (%s). mining...", c.uuid)

	for i := uint64(1); i < c.goroutinesCount; i++ {
		go mine(i, c)
	}

	mine(0, c)
}
