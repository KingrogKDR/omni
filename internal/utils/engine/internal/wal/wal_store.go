package wal

import (
	"fmt"
	"os"
)

type WALOptions struct {
	DirPath     string
	MaxSegments uint32
}

func DefaultWALOptions() WALOptions {
	return WALOptions{
		DirPath:     ".wal",
		MaxSegments: 10,
	}
}

func mergeDefaults(opts WALOptions) WALOptions {
	defaults := DefaultWALOptions()
	if opts.DirPath != "" {
		defaults.DirPath = opts.DirPath
	}
	if opts.MaxSegments != 0 {
		defaults.MaxSegments = opts.MaxSegments
	}
	return defaults
}

func Open(opts WALOptions) (*WAL, error) {
	opts = mergeDefaults(opts)
	if err := os.MkdirAll(opts.DirPath, 0o750); err != nil {
		return nil, fmt.Errorf("creating WAL directory: %w", err)
	}

	wal := &WAL{
		MaxSegments: opts.MaxSegments,
	}
	return wal, nil
}
