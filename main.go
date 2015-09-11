package hashcash

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"strings"
	"time"
)

// Checker provides a struct for checking that hashcach strings are valid according
// to a set difficulty
type Checker struct {
	difficultySlice []byte
	dateFmt         string
	validityWindow  time.Duration
}

// NewChecker creates a new checker
func NewChecker(difficulty uint8, dateFmt string, validityWindow time.Duration) *Checker {
	return &Checker{
		difficultySlice: make([]byte, difficulty),
		dateFmt:         dateFmt,
		validityWindow:  validityWindow,
	}
}

// Check checks a puzzle and returns nil if it demonstrates the correct amount
// of work or an error
func (c *Checker) Check(puzzle string) error {
	hash := sha1.Sum([]byte(puzzle))
	if valid := bytes.Compare(c.difficultySlice, hash[0:len(c.difficultySlice)]) == 0; !valid {
		return fmt.Errorf("Puzzle doesn't demonstrate the required amount of work")
	}
	fields := strings.Split(puzzle, ":")
	if len(fields) < 3 {
		return fmt.Errorf("Malformed puzzle")
	}
	// Check date
	hashDate, err := time.ParseInLocation(c.dateFmt, fields[2], time.UTC)
	if err != nil {
		return err
	}
	if time.Now().Add(-(c.validityWindow / 2)).After(hashDate) {
		return fmt.Errorf("Puzzle timestamp is too old")
	} else if time.Now().Add(c.validityWindow / 2).Before(hashDate) {
		return fmt.Errorf("Puzzle timestamp is too far in the future")
	}
	// Check something else?
	return nil
}
