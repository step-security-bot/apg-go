package main

import (
	"math/big"
	"testing"
)

// Make sure the flags are initalized
var _ = func() bool {
	testing.Init()
	return true
}()

// Test getRandNum with max 1000
func TestGetRandNum(t *testing.T) {
	t.Run("maxNum_is_1000", func(t *testing.T) {
		randNum, err := getRandNum(1000)
		if err != nil {
			t.Errorf("Random number generation failed: %v", err.Error())
		}
		if randNum > 1000 {
			t.Errorf("Generated random number between 0 and 1000 is too big: %v", randNum)
		}
		if randNum < 0 {
			t.Errorf("Generated random number between 0 and 1000 is too small: %v", randNum)
		}
	})
	t.Run("maxNum_is_1", func(t *testing.T) {
		randNum, err := getRandNum(1)
		if err != nil {
			t.Errorf("Random number generation failed: %v", err.Error())
		}
		if randNum > 1 {
			t.Errorf("Generated random number between 0 and 1000 is too big: %v", randNum)
		}
		if randNum < 0 {
			t.Errorf("Generated random number between 0 and 1000 is too small: %v", randNum)
		}
	})
	t.Run("maxNum_is_0", func(t *testing.T) {
		randNum, err := getRandNum(0)
		if err == nil {
			t.Errorf("Random number expected to fail, but provided a value instead: %v", randNum)
		}
	})
}

// Test getRandChar
func TestGetRandChar(t *testing.T) {
	t.Run("return_value_is_A_B_or_C", func(t *testing.T) {
		charRange := "ABC"
		randChar, err := getRandChar(&charRange, 1)
		if err != nil {
			t.Errorf("Random character generation failed => %v", err.Error())
		}
		if randChar != "A" && randChar != "B" && randChar != "C" {
			t.Errorf("Random character generation failed. Expected A, B or C but got: %v", randChar)
		}
	})

	t.Run("return_value_has_specific_length", func(t *testing.T) {
		charRange := "ABC"
		randChar, err := getRandChar(&charRange, 1000)
		if err != nil {
			t.Errorf("Random character generation failed => %v", err.Error())
		}
		if len(randChar) != 1000 {
			t.Errorf("Generated random characters with 1000 chars returned wrong amount of chars: %v",
				len(randChar))
		}
	})

	t.Run("fail", func(t *testing.T) {
		charRange := "ABC"
		randChar, err := getRandChar(&charRange, -2000)
		if err == nil {
			t.Errorf("Generated random characters expected to fail, but returned a value => %v",
				randChar)
		}
	})
}

// Test getCharRange() with different config settings
func TestGetCharRange(t *testing.T) {

	t.Run("lower_case_only", func(t *testing.T) {
		// Lower case only
		allowedBytes := []int{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r',
			's', 't', 'u', 'v', 'w', 'x', 'y', 'z'}
		config.useLowerCase = true
		charRange := getCharRange()
		for _, curChar := range charRange {
			searchAllowedBytes := containsByte(allowedBytes, int(curChar))
			if !searchAllowedBytes {
				t.Errorf("Character range for lower-case only returned invalid value: %v",
					string(curChar))
			}
		}
	})

	// Lower case only (human readable)
	t.Run("lower_case_only_human_readable", func(t *testing.T) {
		allowedBytes := []int{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'j', 'k', 'm', 'n', 'p', 'q', 'r',
			's', 't', 'u', 'v', 'w', 'x', 'y', 'z'}
		config.humanReadable = true
		charRange := getCharRange()
		for _, curChar := range charRange {
			searchAllowedBytes := containsByte(allowedBytes, int(curChar))
			if !searchAllowedBytes {
				t.Errorf("Character range for lower-case only (human readable) returned invalid value: %v",
					string(curChar))
			}
		}
	})

	// Upper case only
	t.Run("upper_case_only", func(t *testing.T) {
		allowedBytes := []int{'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R',
			'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'}
		config.humanReadable = false
		config.useLowerCase = false
		config.useUpperCase = true
		charRange := getCharRange()
		for _, curChar := range charRange {
			searchAllowedBytes := containsByte(allowedBytes, int(curChar))
			if !searchAllowedBytes {
				t.Errorf("Character range for upper-case only returned invalid value: %v",
					string(curChar))
			}
		}
	})

	// Upper case only (human readable)
	t.Run("upper_case_only_human_readable", func(t *testing.T) {
		allowedBytes := []int{'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'J', 'K', 'M', 'N', 'P', 'Q',
			'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'}
		config.humanReadable = true
		charRange := getCharRange()
		for _, curChar := range charRange {
			searchAllowedBytes := containsByte(allowedBytes, int(curChar))
			if !searchAllowedBytes {
				t.Errorf("Character range for upper-case only (human readable) returned invalid value: %v",
					string(curChar))
			}
		}
	})

	// Numbers only
	t.Run("numbers_only", func(t *testing.T) {
		allowedBytes := []int{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}
		config.humanReadable = false
		config.useUpperCase = false
		config.useNumber = true
		charRange := getCharRange()
		for _, curChar := range charRange {
			searchAllowedBytes := containsByte(allowedBytes, int(curChar))
			if !searchAllowedBytes {
				t.Errorf("Character range for numbers only returned invalid value: %v",
					string(curChar))
			}
		}
	})

	// Numbers only (human readable)
	t.Run("numbers_only_human_readable", func(t *testing.T) {
		allowedBytes := []int{'2', '3', '4', '5', '6', '7', '8', '9'}
		config.humanReadable = true
		charRange := getCharRange()
		for _, curChar := range charRange {
			searchAllowedBytes := containsByte(allowedBytes, int(curChar))
			if !searchAllowedBytes {
				t.Errorf("Character range for numbers (human readable) only returned invalid value: %v",
					string(curChar))
			}
		}
	})

	// Special characters only
	t.Run("special_chars_only", func(t *testing.T) {
		allowedBytes := []int{'!', '"', '#', '$', '%', '&', '\'', '(', ')', '*', '+', ',', '-', '.', '/', ':',
			';', '<', '=', '>', '?', '@', '[', '\\', ']', '^', '_', '`', '{', '|', '}', '~'}
		config.humanReadable = false
		config.useNumber = false
		config.useSpecial = true
		charRange := getCharRange()
		for _, curChar := range charRange {
			searchAllowedBytes := containsByte(allowedBytes, int(curChar))
			if !searchAllowedBytes {
				t.Errorf("Character range for special characters only returned invalid value: %v",
					string(curChar))
			}
		}
	})

	// Special characters only (human readable)
	t.Run("special_chars_only_human_readable", func(t *testing.T) {
		allowedBytes := []int{'"', '#', '%', '*', '+', '-', '/', ':', ';', '=', '\\', '_', '|', '~'}
		config.humanReadable = true
		charRange := getCharRange()
		for _, curChar := range charRange {
			searchAllowedBytes := containsByte(allowedBytes, int(curChar))
			if !searchAllowedBytes {
				t.Errorf("Character range for special characters only returned invalid value: %v",
					string(curChar))
			}
		}
	})
}

// Forced failures
func TestForceFailures(t *testing.T) {
	t.Run("too_big_big.NewInt_value", func(t *testing.T) {
		maxNum := 9223372036854775807
		maxNumBigInt := big.NewInt(int64(maxNum) + 1)
		if maxNumBigInt.IsUint64() {
			t.Errorf("Calling big.NewInt() with too large number expected to fail: %v", maxNumBigInt)
		}
	})

	t.Run("negative value for big.NewInt()", func(t *testing.T) {
		randNum, err := getRandNum(-20000)
		if err == nil {
			t.Errorf("Calling getRandNum() with negative value is expected to fail, but returned value: %v",
				randNum)
		}
	})
}

// Benchmark: Random number generation
func BenchmarkGetRandNum(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = getRandNum(100000)
	}
}

// Benchmark: Random char generation
func BenchmarkGetRandChar(b *testing.B) {
	charRange := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890\"#/!\\$%&+-*.,?=()[]{}:;~^|"
	for i := 0; i < b.N; i++ {
		_, _ = getRandChar(&charRange, 20)
	}
}

// Contains function to search a given slice for values
func containsByte(allowedBytes []int, currentChar int) bool {
	for _, charInt := range allowedBytes {
		if charInt == currentChar {
			return true
		}
	}
	return false
}
