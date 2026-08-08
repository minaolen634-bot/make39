# MAKE39 - Bitcoin Seed Phrase Generator

Eesti keeles | In English

---

## Eesti

MAKE39 on Go keeles kirjutatud turvaline tööriist, mis genereerib BIP-39 standardile vastava 24-sõnalise Bitcoin seemnefraasi. Programm kombineerib süsteemse krüptograafilise juhuslikkuse (crypto/rand) ja käsitsi sisestatud füüsilised täringuvisked — luues ühe maailma parima ja turvalisima entroopia seed-fraaside jaoks.

### Peamised omadused

- Kakskeelne liides: Programmi käivitamisel saab valida eesti või inglise keele vahel.
- Tipptasemel entroopia: Süsteemi OS-i juhuslikkus seotakse füüsiliste täringuvisketega (SHA-256 kaudu). See välistab täielikult riistvaralised ja tarkvaralised tagauksed (backdoors).
- Automaksumuse kontroll: Programm kontrollib Linuxi süsteemi saadaolevat entroopiat (/proc/sys/kernel/random/entropy_avail) ja teavitab kasutajat selle seisukorrast.
- Sõltuvusteta ja puhas: Ei kasuta väliseid kolmanda osapoole teeke.

### Kuidas see töötab

1. Loeb BIP-39 ingliskeelse sõnaloendi failist english.txt (täpselt 2048 sõna).
2. Kogub 32 baiti (256 bitti) krüptograafilist juhuslikkust süsteemi crypto/rand kaudu.
3. Küsib kasutajalt füüsilisi täringuviskeid (kohustuslik või valikuline sõltuvalt OS-i entroopiast).
4. Segab süsteemi entroopia ja täringuvisked kokku SHA-256 abil.
5. Arvutab kontrollsumma ja teisendab tulemuse 24-sõnaliseks BIP-39 fraasiks.

### Paigaldamine ja kasutamine

1. Veendu, et english.txt on samas kaustas kus binaar.
2. Kompileeri programm:
   go build -ldflags="-s -w" -o make39 main.go
3. Käivita programm:
   ./make39
4. Vali alguses soovitud keel (1: Eesti, 2: English).
5. Kirjuta saadud 24-sõnaline fraas turvaliselt paberile või terasplaadile.

---

## English

MAKE39 is a secure Bitcoin seed phrase generator written in Go that strictly adheres to the BIP-39 standard. By combining operating system cryptographic randomness (crypto/rand) with manual physical dice rolls, it provides one of the world's most resilient and highest-quality entropy sources for seed phrase generation.

### Key Features

- Bilingual Interface: Built-in interactive menu allowing you to choose between Estonian and English on startup.
- World-Class Entropy: Fuses OS-level randomness with real-world physical dice rolls via SHA-256 mixing, mitigating hardware RNG backdoors and software vulnerabilities.
- Linux Entropy Check: Automatically inspects system entropy level (/proc/sys/kernel/random/entropy_avail) and prompts for physical rolls if system entropy is low.
- Zero External Dependencies: Built strictly using Go standard library for maximum auditability and security.

### How It Works

1. Loads the standard BIP-39 English wordlist from english.txt (exactly 2048 words).
2. Generates 32 bytes (256 bits) of cryptographically secure randomness via crypto/rand.
3. Collects user-provided physical dice rolls (mandatory if OS entropy is low, optional otherwise).
4. Combines system randomness and dice entropy into a unified 256-bit entropy block using SHA-256.
5. Appends the checksum and converts the final entropy bitstream into a 24-word BIP-39 mnemonic.

### Build and Run

1. Ensure english.txt is present in the same directory.
2. Compile the binary:
   go build -ldflags="-s -w" -o make39 main.go
3. Run the executable:
   ./make39
4. Choose your preferred language (1 for Estonian, 2 for English).
5. Safely write down the generated 24-word mnemonic phrase offline.

---

## Turvalisus ja hoiatused / Security & Warnings

- Air-Gapped Execution: For maximum security, always compile and run this software on an offline (air-gapped) machine.
- Physical Dice: Always use real, high-quality physical dice. Do not use software-simulated dice.
- No Digital Copies: Never screenshot, copy to clipboard, or store your generated 24 words on any internet-connected device.

## Litsents / License

Copyright (c) 2026 minaolen634@gmail.com
Litsenseeritud BSD 2-Clause License alusel (Vaba tarkvara / Free Software).
