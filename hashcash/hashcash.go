package hashcash

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var (
	ErrParse             = errors.New("error parsing hashcash")
	ErrInvalidDifficulty = errors.New("invalid difficulty")
	ErrExpired           = errors.New("expired hashcash")
	ErrInvalidAlgorithm  = errors.New("invalid algorithm")
	ErrInvalidSolution   = errors.New("invalid solution")
	ErrDifficultyTooHigh = errors.New("difficulty too high")
	ErrInvalidDate       = errors.New("invalid date")
	ErrInvalidTag        = errors.New("expected version to be '1'")
)

// MaxDifficulty is the upper limit for difficulty
const (
	MaxDifficulty          = 26
	defaultHashcashVersion = 1
	defaultHashcashAlg     = "SHA-256"
)

// Hashcash represents a parsed Hashcash string
type Hashcash struct {
	Version    int       `json:"version"`    // Hashcash version (always "1")
	Difficulty int       `json:"difficulty"` // Number of leading zero bits required in hash
	ExpiresAt  time.Time `json:"exp"`        // Expiration time
	Subject    string    `json:"sub"`        // Subject or resource being protected
	Nonce      string    `json:"nonce"`      // Random nonce
	Alg        string    `json:"alg"`        // Hashing algorithm (always SHA-256)
	Solution   string    `json:"solution"`   // Solution to the hashcash challenge
}

// New creates a new Hashcash with default values
func New(subject string, difficulty int) *Hashcash {
	nonce := make([]byte, 16)
	if _, err := rand.Read(nonce); err != nil {
		panic(fmt.Errorf("error generating nonce: %w", err))
	}

	return &Hashcash{
		Version:    defaultHashcashVersion,
		Difficulty: difficulty,
		ExpiresAt:  time.Now().UTC().Add(5 * time.Minute).Truncate(time.Second),
		Subject:    subject,
		Nonce:      base64.RawURLEncoding.EncodeToString(nonce),
		Alg:        defaultHashcashAlg,
	}
}

// String returns the hashcash as a formatted string
func (h *Hashcash) String() string {
	expAt := strconv.FormatInt(h.ExpiresAt.UTC().Unix(), 10)
	parts := []string{strconv.Itoa(h.Version), strconv.Itoa(h.Difficulty), expAt, h.Subject, h.Nonce, h.Alg}
	if h.Solution != "" {
		parts = append(parts, h.Solution)
	}
	return strings.Join(parts, ":")
}

// Verify checks if the hashcash is valid
func (h *Hashcash) Verify() error {
	if h.Difficulty < 0 || h.Difficulty > MaxDifficulty {
		return ErrInvalidDifficulty
	}
	if h.Alg != defaultHashcashAlg {
		return ErrInvalidAlgorithm
	}
	if h.ExpiresAt.Before(time.Now().UTC()) {
		return ErrExpired
	}
	return h.verifySolution()
}

func (h *Hashcash) verifySolution() error {
	hash := sha256.Sum256([]byte(h.String()))
	bitCount := h.Difficulty
	for i := 0; i < len(hash); i++ {
		for j := 7; j >= 0; j-- {
			if bitCount == 0 {
				return nil
			}
			if (hash[i]>>j)&1 == 1 {
				return ErrInvalidSolution
			}
			bitCount--
		}
	}
	return nil
}

// Solve attempts to find a solution for the hashcash
func (h *Hashcash) Solve(ctx context.Context) error {
	if h.Difficulty > MaxDifficulty {
		return ErrDifficultyTooHigh
	}
	if h.Alg != defaultHashcashAlg {
		return ErrInvalidAlgorithm
	}

	var solution uint32
	sb := make([]byte, 4)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			binary.LittleEndian.PutUint32(sb, solution)
			h.Solution = base64.RawURLEncoding.EncodeToString(sb)
			if err := h.verifySolution(); err == nil {
				return nil
			}
			solution++
		}
	}
}

// Parse parses a hashcash string into a Hashcash struct
func Parse(hcStr string) (*Hashcash, error) {
	parts := strings.Split(hcStr, ":")
	if len(parts) < 6 || len(parts) > 7 {
		return nil, ErrParse
	}

	versionStr, difficultyStr, expStr, subject, nonce, alg := parts[0], parts[1], parts[2], parts[3], parts[4], parts[5]
	version, err := strconv.Atoi(versionStr)
	if err != nil || version != defaultHashcashVersion {
		return nil, ErrInvalidTag
	}

	difficulty, err := strconv.Atoi(difficultyStr)
	if err != nil || difficulty < 0 || difficulty > MaxDifficulty {
		return nil, ErrInvalidDifficulty
	}

	var exp time.Time
	if expStr != "" {
		expAt, err := strconv.ParseInt(expStr, 10, 64)
		if err != nil || expAt < 0 {
			return nil, ErrInvalidDate
		}
		exp = time.Unix(expAt, 0).UTC()
	}

	var solution string
	if len(parts) == 7 {
		solution = parts[6]
	}

	return &Hashcash{
		Version:    version,
		Difficulty: difficulty,
		ExpiresAt:  exp,
		Subject:    subject,
		Nonce:      nonce,
		Alg:        alg,
		Solution:   solution,
	}, nil
}
