package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

type Config struct {
	InputFile    string        `json:"input_file"`
	OutputFile   string        `json:"output_file"`
	Replacements []Replacement `json:"replacements"`
}

type Replacement struct {
	OldHex string `json:"old_hex"`
	NewHex string `json:"new_hex"`
}

func printBanner() {
	fmt.Println("╔══════════════════════════════════════════════╗")
	fmt.Println("║            ██╗  ██╗███████╗██╗  ██╗          ║")
	fmt.Println("║            ██║  ██║██╔════╝╚██╗██╔╝          ║")
	fmt.Println("║            ███████║█████╗   ╚███╔╝           ║")
	fmt.Println("║            ██╔══██║██╔══╝   ██╔██╗           ║")
	fmt.Println("║            ██║  ██║███████╗██╔╝ ██╗          ║")
	fmt.Println("║            ╚═╝  ╚═╝╚══════╝╚═╝  ╚═╝          ║")
	fmt.Println("╟──────────────────────────────────────────────╢")
	fmt.Println("║      HexPatch v1.0 - Binary Hex Patcher      ║")
	fmt.Println("╚══════════════════════════════════════════════╝")
	fmt.Println()
}

func main() {
	printBanner()

	// Read config file
	configData, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatalf("❌ Failed to read config file: %v", err)
	}

	var config Config
	if err := json.Unmarshal(configData, &config); err != nil {
		log.Fatalf("❌ Config parse error: %v", err)
	}

	fmt.Printf("✅ Config loaded: %d replacement rules\n", len(config.Replacements))

	// Read input file
	content, err := ioutil.ReadFile(config.InputFile)
	if err != nil {
		log.Fatalf("❌ Failed to read input file: %v", err)
	}

	fmt.Printf("📥 Input file read: %s (%d bytes)\n", config.InputFile, len(content))

	// Perform all replacements
	modified := content
	for i, r := range config.Replacements {
		if r.OldHex == "" || r.NewHex == "" {
			fmt.Printf("⚠️  Rule #%d skipped: old or new hex is empty\n", i+1)
			continue
		}
		// Clean hex string (remove spaces, etc.)
		oldHex := cleanHex(r.OldHex)
		newHex := cleanHex(r.NewHex)

		// Decode hex string to bytes
		oldBytes, err := hex.DecodeString(oldHex)
		if err != nil {
			fmt.Printf("⚠️  Rule #%d old hex decode failed (skipped): %v\n", i+1, err)
			continue
		}

		newBytes, err := hex.DecodeString(newHex)
		if err != nil {
			fmt.Printf("⚠️  Rule #%d new hex decode failed (skipped): %v\n", i+1, err)
			continue
		}

		// Check if old and new have the same length
		if len(oldBytes) != len(newBytes) {
			fmt.Printf("⚠️  Rule #%d skipped: old and new length mismatch (old:%d, new:%d)\n",
				i+1, len(oldBytes), len(newBytes))
			continue
		}

		// Find and replace all occurrences
		count := 0
		for {
			index := byteIndex(modified, oldBytes)
			if index == -1 {
				break
			}
			copy(modified[index:index+len(oldBytes)], newBytes)
			count++
		}

		fmt.Printf("🔧 Applied rule #%d: %s → %s (%d replacements)\n",
			i+1, r.OldHex, r.NewHex, count)
	}

	// Write output file
	if err := ioutil.WriteFile(config.OutputFile, modified, 0755); err != nil {
		log.Fatalf("❌ Failed to write output file: %v", err)
	}

	fmt.Printf("\n✅ Operation completed!\nInput: %s\nOutput: %s\n", config.InputFile, config.OutputFile)
	fmt.Printf("Bytes: %d → %d\n", len(content), len(modified))
}

// Clean hex string: remove spaces and convert to uppercase
func cleanHex(hexStr string) string {
	return strings.ToUpper(strings.ReplaceAll(hexStr, " ", ""))
}

// Find the index of the first occurrence of sep in s
func byteIndex(s, sep []byte) int {
	n := len(sep)
	if n == 0 {
		return 0
	}

	if n > len(s) {
		return -1
	}

	for i := 0; i <= len(s)-n; i++ {
		if s[i] == sep[0] {
			match := true
			for j := 1; j < n; j++ {
				if s[i+j] != sep[j] {
					match = false
					break
				}
			}
			if match {
				return i
			}
		}
	}
	return -1
}
