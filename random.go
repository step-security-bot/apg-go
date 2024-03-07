package apg

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"
	"strings"
)

const (
	// 7 bits to represent a letter index
	letterIdxBits = 7
	// All 1-bits, as many as letterIdxBits
	letterIdxMask = 1<<letterIdxBits - 1
	// # of letter indices fitting in 63 bits)
	letterIdxMax = 63 / letterIdxBits
)

// maxInt32 is the maximum positive value for a int32 number type
const maxInt32 = 2147483647

var (
	// ErrInvalidLength is returned if the provided maximum number is equal or less than zero
	ErrInvalidLength = errors.New("provided length value cannot be zero or less")
	// ErrLengthMismatch is returned if the number of generated bytes does not match the expected length
	ErrLengthMismatch = errors.New("number of generated random bytes does not match the expected length")
	// ErrInvalidCharRange is returned if the given range of characters is not valid
	ErrInvalidCharRange = errors.New("provided character range is not valid or empty")
)

// CoinFlip performs a simple coinflip based on the rand library and returns 1 or 0
func (g *Generator) CoinFlip() int64 {
	cf, _ := g.RandNum(2)
	return cf
}

// CoinFlipBool performs a simple coinflip based on the rand library and returns true or false
func (g *Generator) CoinFlipBool() bool {
	return g.CoinFlip() == 1
}

// Generate generates a password based on all the different config flags and returns
// it as string type. If the generation fails, an error will be thrown
func (g *Generator) Generate() (string, error) {
	switch g.config.Algorithm {
	case AlgoCoinFlip:
		return g.generateCoinFlip()
	case AlgoRandom:
		return g.generateRandom()
	case AlgoUnsupported:
		return "", fmt.Errorf("unsupported algorithm")
	}
	return "", nil
}

// GetPasswordLength returns the password length based on the given config
// parameters
func (g *Generator) GetPasswordLength() (int64, error) {
	if g.config.FixedLength > 0 {
		return g.config.FixedLength, nil
	}
	minLength := g.config.MinLength
	maxLength := g.config.MaxLength
	if minLength > maxLength {
		maxLength = minLength
	}
	diff := maxLength - minLength + 1
	randNum, err := g.RandNum(diff)
	if err != nil {
		return 0, err
	}
	length := minLength + randNum
	if length <= 0 {
		return 1, nil
	}
	return length, nil
}

// RandomBytes returns a byte slice of random bytes with given length that got generated by
// the crypto/rand generator
func (g *Generator) RandomBytes(length int64) ([]byte, error) {
	if length < 1 {
		return nil, ErrInvalidLength
	}
	bytes := make([]byte, length)
	numBytes, err := rand.Read(bytes)
	if int64(numBytes) != length {
		return nil, ErrLengthMismatch
	}
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

// RandNum generates a random, non-negative number with given maximum value
func (g *Generator) RandNum(max int64) (int64, error) {
	if max < 1 {
		return 0, ErrInvalidLength
	}
	max64 := big.NewInt(max)
	randNum, err := rand.Int(rand.Reader, max64)
	if err != nil {
		return 0, fmt.Errorf("random number generation failed: %w", err)
	}
	return randNum.Int64(), nil
}

// RandomStringFromCharRange returns a random string of length l based of the range of characters given.
// The method makes use of the crypto/random package and therfore is
// cryptographically secure
func (g *Generator) RandomStringFromCharRange(length int64, charRange string) (string, error) {
	if length < 1 {
		return "", ErrInvalidLength
	}
	if len(charRange) < 1 {
		return "", ErrInvalidCharRange
	}
	randString := strings.Builder{}

	// As long as the length is smaller than the max. int32 value let's grow
	// the string builder to the actual size, so we need less allocations
	if length <= maxInt32 {
		randString.Grow(int(length))
	}

	charRangeLength := len(charRange)

	randPool := make([]byte, 8)
	_, err := rand.Read(randPool)
	if err != nil {
		return randString.String(), err
	}
	for idx, char, rest := length-1, binary.BigEndian.Uint64(randPool), letterIdxMax; idx >= 0; {
		if rest == 0 {
			_, err = rand.Read(randPool)
			if err != nil {
				return randString.String(), err
			}
			char, rest = binary.BigEndian.Uint64(randPool), letterIdxMax
		}
		if i := int(char & letterIdxMask); i < charRangeLength {
			randString.WriteByte(charRange[i])
			idx--
		}
		char >>= letterIdxBits
		rest--
	}

	return randString.String(), nil
}

// GetCharRangeFromConfig checks the Mode from the Config and returns a
// list of all possible characters that are supported by these Mode
func (g *Generator) GetCharRangeFromConfig() string {
	charRange := strings.Builder{}
	if MaskHasMode(g.config.Mode, ModeLowerCase) {
		switch MaskHasMode(g.config.Mode, ModeHumanReadable) {
		case true:
			charRange.WriteString(CharRangeAlphaLowerHuman)
		default:
			charRange.WriteString(CharRangeAlphaLower)
		}
	}
	if MaskHasMode(g.config.Mode, ModeNumeric) {
		switch MaskHasMode(g.config.Mode, ModeHumanReadable) {
		case true:
			charRange.WriteString(CharRangeNumericHuman)
		default:
			charRange.WriteString(CharRangeNumeric)
		}
	}
	if MaskHasMode(g.config.Mode, ModeSpecial) {
		switch MaskHasMode(g.config.Mode, ModeHumanReadable) {
		case true:
			charRange.WriteString(CharRangeSpecialHuman)
		default:
			charRange.WriteString(CharRangeSpecial)
		}
	}
	if MaskHasMode(g.config.Mode, ModeUpperCase) {
		switch MaskHasMode(g.config.Mode, ModeHumanReadable) {
		case true:
			charRange.WriteString(CharRangeAlphaUpperHuman)
		default:
			charRange.WriteString(CharRangeAlphaUpper)
		}
	}
	return charRange.String()
}

func (g *Generator) checkMinimumRequirements(password string) bool {
	ok := true
	if g.config.MinLowerCase > 0 {
		var charRange string
		switch MaskHasMode(g.config.Mode, ModeHumanReadable) {
		case true:
			charRange = CharRangeAlphaLowerHuman
		default:
			charRange = CharRangeAlphaLower
		}

		count := 0
		for _, char := range charRange {
			count += strings.Count(password, string(char))
		}
		if int64(count) < g.config.MinLowerCase {
			ok = false
		}
	}
	if g.config.MinNumeric > 0 {
		var charRange string
		switch MaskHasMode(g.config.Mode, ModeHumanReadable) {
		case true:
			charRange = CharRangeNumericHuman
		default:
			charRange = CharRangeNumeric
		}

		count := 0
		for _, char := range charRange {
			count += strings.Count(password, string(char))
		}
		if int64(count) < g.config.MinNumeric {
			ok = false
		}
	}
	if g.config.MinSpecial > 0 {
		var charRange string
		switch MaskHasMode(g.config.Mode, ModeHumanReadable) {
		case true:
			charRange = CharRangeSpecialHuman
		default:
			charRange = CharRangeSpecial
		}

		count := 0
		for _, char := range charRange {
			count += strings.Count(password, string(char))
		}
		if int64(count) < g.config.MinSpecial {
			ok = false
		}
	}
	if g.config.MinUpperCase > 0 {
		var charRange string
		switch MaskHasMode(g.config.Mode, ModeHumanReadable) {
		case true:
			charRange = CharRangeAlphaUpperHuman
		default:
			charRange = CharRangeAlphaUpper
		}

		count := 0
		for _, char := range charRange {
			count += strings.Count(password, string(char))
		}
		if int64(count) < g.config.MinUpperCase {
			ok = false
		}
	}
	return ok
}

// generateCoinFlip is executed when Generate() is called with Algorithm set
// to AlgoCoinFlip
func (g *Generator) generateCoinFlip() (string, error) {
	if g.CoinFlipBool() {
		return "Heads", nil
	}
	return "Tails", nil
}

// generateRandom is executed when Generate() is called with Algorithm set
// to AlgoRandmom
func (g *Generator) generateRandom() (string, error) {
	length, err := g.GetPasswordLength()
	if err != nil {
		return "", fmt.Errorf("failed to calculate password length: %w", err)
	}
	charRange := g.GetCharRangeFromConfig()
	var password string
	var ok bool
	for !ok {
		password, err = g.RandomStringFromCharRange(length, charRange)
		if err != nil {
			return "", err
		}
		ok = g.checkMinimumRequirements(password)
	}

	return password, nil
}
