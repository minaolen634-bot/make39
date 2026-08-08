/*
 * MAKE39
 * Copyright (c) 2026 minaolen634@gmail.com
 * Litsents: BSD 2-Clause License (Vaba tarkvara)
 */

package main

import (
	"bufio"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
)

const (
	wordlistFile = "english.txt"
	entropySize  = 32  // 256 bits
	minEntropy   = 256 // bits
)

type Lang struct {
	WordlistErr     string
	EntropyErr      string
	AvailEntropy    string
	CryptoRandErr   string
	EntropyLow      string
	EntropyOK       string
	DiceErr         string
	MnemonicErr     string
	ResultHeader    string
	DiceCountPrompt string
	InputEmpty      string
	InputPosNum     string
	RollPrompt      string
	InvalidInput    string
	AddDicePrompt   string
	YesTrigger      []string
}

var langET = Lang{
	WordlistErr:     "Sõnalisti viga:",
	EntropyErr:      "Entroopia kontrolli viga:",
	AvailEntropy:    "Linuxi saadaolev entroopia: %d bitti\n",
	CryptoRandErr:   "crypto/rand viga:",
	EntropyLow:      "Entroopiat on vähe, täringuvisked on kohustuslikud.",
	EntropyOK:       "Entroopiat on piisavalt, täringuvisked on valikulised.",
	DiceErr:         "Täringu sisendi viga:",
	MnemonicErr:     "Mnemoühendi viga:",
	ResultHeader:    "\nSinu 24-sõnaline BIP-39 mnemoühend:",
	DiceCountPrompt: "Mitu täringuviset lisad? ",
	InputEmpty:      "sisend on tühi",
	InputPosNum:     "sisesta positiivne arv",
	RollPrompt:      "Vise %d/%d (1-6): ",
	InvalidInput:    "Vale sisend, proovi uuesti.",
	AddDicePrompt:   "Kas tahad lisada täringuviskeid? (j/e): ",
	YesTrigger:      []string{"j", "jah"},
}

var langEN = Lang{
	WordlistErr:     "Wordlist error:",
	EntropyErr:      "Entropy check error:",
	AvailEntropy:    "Linux entropy available: %d bits\n",
	CryptoRandErr:   "crypto/rand error:",
	EntropyLow:      "Entropy is low, dice rolls are required.",
	EntropyOK:       "Entropy looks sufficient, dice rolls are optional.",
	DiceErr:         "Dice input error:",
	MnemonicErr:     "Mnemonic error:",
	ResultHeader:    "\nYour 24-word BIP-39 mnemonic:",
	DiceCountPrompt: "How many dice rolls to add? ",
	InputEmpty:      "input is empty",
	InputPosNum:     "enter a positive number",
	RollPrompt:      "Roll %d/%d (1-6): ",
	InvalidInput:    "Invalid input, try again.",
	AddDicePrompt:   "Do you want to add dice rolls? (y/n): ",
	YesTrigger:      []string{"y", "yes"},
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	txt := selectLanguage(reader)

	words, err := loadWordlist(wordlistFile)
	if err != nil {
		fmt.Println(txt.WordlistErr, err)
		return
	}

	linuxEntropyOK, availBits, err := linuxEntropySufficient(minEntropy)
	if err != nil {
		fmt.Println(txt.EntropyErr, err)
		return
	}

	fmt.Printf(txt.AvailEntropy, availBits)

	randomEntropy := make([]byte, entropySize)
	if _, err := rand.Read(randomEntropy); err != nil {
		fmt.Println(txt.CryptoRandErr, err)
		return
	}

	var diceEntropy []byte
	if !linuxEntropyOK {
		fmt.Println(txt.EntropyLow)
		diceEntropy, err = collectDice(reader, txt)
		if err != nil {
			fmt.Println(txt.DiceErr, err)
			return
		}
	} else {
		fmt.Println(txt.EntropyOK)
		diceEntropy, err = collectDiceOptional(reader, txt)
		if err != nil {
			fmt.Println(txt.DiceErr, err)
			return
		}
	}

	finalEntropy := mixEntropy(randomEntropy, diceEntropy)

	mnemonic, err := entropyToMnemonic(finalEntropy, words)
	if err != nil {
		fmt.Println(txt.MnemonicErr, err)
		return
	}

	fmt.Println(txt.ResultHeader)
	fmt.Println(mnemonic)
}

func selectLanguage(reader *bufio.Reader) Lang {
	for {
		fmt.Println("Select language / Vali keel:")
		fmt.Println("1) Eesti (ET)")
		fmt.Println("2) English (EN)")
		fmt.Print("Choice / Valik (1-2): ")

		line, err := reader.ReadString('\n')
		if err != nil {
			return langET
		}

		choice := strings.TrimSpace(line)
		switch choice {
		case "1", "et", "ET", "eesti", "Eesti":
			fmt.Println()
			return langET
		case "2", "en", "EN", "english", "English":
			fmt.Println()
			return langEN
		default:
			fmt.Println("Invalid choice, please select 1 or 2.\n")
		}
	}
}

func linuxEntropySufficient(minBits int) (bool, int, error) {
	data, err := os.ReadFile("/proc/sys/kernel/random/entropy_avail")
	if err != nil {
		return false, 0, err
	}

	s := strings.TrimSpace(string(data))
	avail, err := strconv.Atoi(s)
	if err != nil {
		return false, 0, err
	}

	return avail >= minBits, avail, nil
}

func loadWordlist(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) != 2048 {
		return nil, fmt.Errorf("wordlist must contain exactly 2048 words, got %d", len(lines))
	}

	for i := range lines {
		lines[i] = strings.TrimSpace(lines[i])
	}

	return lines, nil
}

func collectDice(reader *bufio.Reader, txt Lang) ([]byte, error) {
	fmt.Print(txt.DiceCountPrompt)
	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	line = strings.TrimSpace(line)
	if line == "" {
		return nil, fmt.Errorf("%s", txt.InputEmpty)
	}

	n, err := strconv.Atoi(line)
	if err != nil || n <= 0 {
		return nil, fmt.Errorf("%s", txt.InputPosNum)
	}

	rolls := make([]byte, n)

	for i := 0; i < n; i++ {
		for {
			fmt.Printf(txt.RollPrompt, i+1, n)
			line, err := reader.ReadString('\n')
			if err != nil {
				return nil, err
			}

			line = strings.TrimSpace(line)
			roll, err := strconv.Atoi(line)
			if err != nil || roll < 1 || roll > 6 {
				fmt.Println(txt.InvalidInput)
				continue
			}

			rolls[i] = byte(roll)
			break
		}
	}

	return rolls, nil
}

func collectDiceOptional(reader *bufio.Reader, txt Lang) ([]byte, error) {
	fmt.Print(txt.AddDicePrompt)
	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	line = strings.TrimSpace(strings.ToLower(line))
	isYes := false
	for _, trigger := range txt.YesTrigger {
		if line == trigger {
			isYes = true
			break
		}
	}

	if !isYes {
		return []byte{}, nil
	}

	return collectDice(reader, txt)
}

func mixEntropy(randomEntropy, diceEntropy []byte) []byte {
	h := sha256.New()
	h.Write(randomEntropy)
	h.Write(diceEntropy)
	return h.Sum(nil)
}

func entropyToMnemonic(entropy []byte, words []string) (string, error) {
	if len(entropy) != 32 {
		return "", fmt.Errorf("entropy must be 32 bytes")
	}

	checksumBits := len(entropy) * 8 / 32
	checksum := sha256.Sum256(entropy)

	entropyBits := bytesToBits(entropy)
	checksumBitString := bitsFromBytes(checksum[:], checksumBits)
	bitString := entropyBits + checksumBitString

	if len(bitString)%11 != 0 {
		return "", fmt.Errorf("invalid bit length")
	}

	var mnemonic []string
	for i := 0; i < len(bitString); i += 11 {
		index := bitsToInt(bitString[i : i+11])
		if index >= len(words) {
			return "", fmt.Errorf("word index out of range")
		}
		mnemonic = append(mnemonic, words[index])
	}

	return strings.Join(mnemonic, " "), nil
}

func bytesToBits(b []byte) string {
	var sb strings.Builder
	for _, v := range b {
		sb.WriteString(fmt.Sprintf("%08b", v))
	}
	return sb.String()
}

func bitsFromBytes(b []byte, count int) string {
	bits := bytesToBits(b)
	if count > len(bits) {
		count = len(bits)
	}
	return bits[:count]
}

func bitsToInt(bitStr string) int {
	n := big.NewInt(0)
	base := big.NewInt(2)
	for _, c := range bitStr {
		n.Mul(n, base)
		if c == '1' {
			n.Add(n, big.NewInt(1))
		}
	}
	return int(n.Int64())
}
