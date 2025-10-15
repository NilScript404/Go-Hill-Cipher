package main

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"gonum.org/v1/gonum/mat"
)

type HillCipher struct {
	dimension        int
	keyMatrix        *mat.Dense
	inverseKeyMatrix *mat.Dense
}

var alphabet map[string]int
var reverseAlphabet map[int]string

func main() {
	initializeAlphabetMaps()

	fmt.Println("--- Hill Cipher Implementation ---")

	keyString, dimension, err := promptForKeyDetails()
	if err != nil {
		fmt.Printf("Error during key setup: %v\n", err)
		return
	}

	message, err := promptForMessage()
	if err != nil {
		fmt.Printf("Error reading message: %v\n", err)
		return
	}

	fmt.Println("\n--- Key Setup ---")
	cipher, err := NewHillCipher(keyString, dimension)
	if err != nil {
		fmt.Printf("Error initializing cipher: %v\n", err)
		return
	}
	printMatrix(cipher.keyMatrix, "1. Key Matrix (Numerical):")
	printMatrix(cipher.inverseKeyMatrix, "2. Inverse Key Matrix (K⁻¹):")

	fmt.Println("\n--- Encryption Process ---")
	ciphertext, err := cipher.Encrypt(message)
	if err != nil {
		fmt.Printf("Error during encryption: %v\n", err)
		return
	}
	fmt.Printf("\n4. Final Encrypted Message: %s\n", ciphertext)

	fmt.Println("\n--- Decryption Process ---")
	decryptedText, err := cipher.Decrypt(ciphertext)
	if err != nil {
		fmt.Printf("Error during decryption: %v\n", err)
		return
	}
	fmt.Printf("\n4. Final Decrypted Message: %s\n", decryptedText)
}

func NewHillCipher(key string, dimension int) (*HillCipher, error) {
	keyMatrix, err := stringToMatrix(key, dimension)
	if err != nil {
		return nil, fmt.Errorf("could not create key matrix: %w", err)
	}

	inverseKeyMatrix, err := calculateInverseKeyMatrix(keyMatrix)
	if err != nil {
		return nil, fmt.Errorf("could not calculate inverse key: %w", err)
	}

	return &HillCipher{
		dimension:        dimension,
		keyMatrix:        keyMatrix,
		inverseKeyMatrix: inverseKeyMatrix,
	}, nil
}

func (hc *HillCipher) Encrypt(plaintext string) (string, error) {
	if plaintext == "" {
		return "", nil
	}
	fmt.Printf("Plaintext: %s\n", plaintext)
	paddedText := padMessage(plaintext, hc.dimension)
	if len(paddedText) != len(plaintext) {
		fmt.Printf("Message padded to '%s' for vectorization.\n", paddedText)
	}

	vectors, err := textToVectors(paddedText, hc.dimension)
	if err != nil {
		return "", fmt.Errorf("failed to convert plaintext to vectors: %w", err)
	}

	fmt.Println("\n1. Plaintext to Vectors:")
	for i, v := range vectors {
		printMatrix(v, fmt.Sprintf("   Vector P%d:", i+1))
	}

	encryptedVectors := make([]*mat.Dense, len(vectors))
	for i, p := range vectors {
		c := mat.NewDense(hc.dimension, 1, nil)
		c.Mul(hc.keyMatrix, p)
		encryptedVectors[i] = c
	}

	fmt.Println("\n2. Multiplying Key Matrix with Plaintext Vectors (K * P):")
	for i, v := range encryptedVectors {
		printMatrix(v, fmt.Sprintf("   Result Vector %d:", i+1))
	}

	applyMod26ToMatrices(encryptedVectors)

	fmt.Println("\n3. Applying Modulo 26 to Result Vectors:")
	for i, v := range encryptedVectors {
		printMatrix(v, fmt.Sprintf("   Final Vector C%d:", i+1))
	}

	ciphertext, err := vectorsToText(encryptedVectors)
	if err != nil {
		return "", fmt.Errorf("failed to convert encrypted vectors to text: %w", err)
	}

	return ciphertext, nil
}

func (hc *HillCipher) Decrypt(ciphertext string) (string, error) {
	if ciphertext == "" {
		return "", nil
	}
	fmt.Printf("Ciphertext: %s\n", ciphertext)

	vectors, err := textToVectors(ciphertext, hc.dimension)
	if err != nil {
		return "", fmt.Errorf("failed to convert ciphertext to vectors: %w", err)
	}

	fmt.Println("\n1. Ciphertext to Vectors:")
	for i, v := range vectors {
		printMatrix(v, fmt.Sprintf("   Vector C%d:", i+1))
	}

	decryptedVectors := make([]*mat.Dense, len(vectors))
	for i, c := range vectors {
		p := mat.NewDense(hc.dimension, 1, nil)
		p.Mul(hc.inverseKeyMatrix, c)
		decryptedVectors[i] = p
	}

	fmt.Println("\n2. Multiplying Inverse Key Matrix with Ciphertext Vectors (K⁻¹ * C):")
	for i, v := range decryptedVectors {
		printMatrix(v, fmt.Sprintf("   Result Vector %d:", i+1))
	}

	applyMod26ToMatrices(decryptedVectors)

	fmt.Println("\n3. Applying Modulo 26 to Result Vectors:")
	for i, v := range decryptedVectors {
		printMatrix(v, fmt.Sprintf("   Final Vector P%d:", i+1))
	}

	plaintext, err := vectorsToText(decryptedVectors)
	if err != nil {
		return "", fmt.Errorf("failed to convert decrypted vectors to text: %w", err)
	}

	return plaintext, nil
}

func initializeAlphabetMaps() {
	alphabet = make(map[string]int)
	reverseAlphabet = make(map[int]string)
	for i := range 26 {
		letter := string('A' + i)
		alphabet[letter] = i
		reverseAlphabet[i] = letter
	}
}

func promptForKeyDetails() (key string, dimension int, err error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter the key matrix dimension (e.g., 2 for a 2x2 matrix): ")
	dimInput, err := reader.ReadString('\n')
	if err != nil {
		return "", 0, fmt.Errorf("failed to read dimension: %w", err)
	}
	dimension, err = strconv.Atoi(strings.TrimSpace(dimInput))
	if err != nil {
		return "", 0, errors.New("dimension must be a valid integer")
	}

	if dimension <= 0 {
		return "", 0, errors.New("matrix dimension must be positive")
	}

	fmt.Print("Enter the key value: ")
	keyInput, err := reader.ReadString('\n')
	if err != nil {
		return "", 0, fmt.Errorf("failed to read key: %w", err)
	}

	processedKey := strings.TrimSpace(keyInput)
	for _, char := range processedKey {
		if (char < 'a' || char > 'z') && (char < 'A' || char > 'Z') {
			return "", 0, fmt.Errorf("key must only contain alphabetic characters, but found '%c'", char)
		}
	}

	key = strings.ToUpper(processedKey)

	if len(key) != dimension*dimension {
		err = fmt.Errorf("key value length (%d) does not match matrix dimensions (%d*%d=%d)", len(key), dimension, dimension, dimension*dimension)
		return "", 0, err
	}

	return key, dimension, nil
}

func promptForMessage() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter the message to encrypt: ")
	message, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read message: %w", err)
	}

	originalMessage := strings.TrimSpace(message)
	var builder strings.Builder
	for _, char := range originalMessage {
		if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') {
			builder.WriteRune(char)
		}
	}
	sanitizedMessage := strings.ToUpper(builder.String())

	if len(originalMessage) > 0 && len(sanitizedMessage) == 0 {
		return "", errors.New("message contains no valid alphabetic characters to encrypt")
	}

	if len(sanitizedMessage) > 0 {
		fmt.Printf("Sanitized message for processing: %s\n", sanitizedMessage)
	}

	return sanitizedMessage, nil
}

func stringToMatrix(s string, dim int) (*mat.Dense, error) {
	values := make([]float64, len(s))
	for i, char := range s {
		val, ok := alphabet[string(char)]
		if !ok {
			return nil, fmt.Errorf("invalid character in key: %c", char)
		}
		values[i] = float64(val)
	}
	return mat.NewDense(dim, dim, values), nil
}

func textToVectors(text string, dim int) ([]*mat.Dense, error) {
	numVectors := len(text) / dim
	vectors := make([]*mat.Dense, numVectors)

	for i := range numVectors {
		chunk := text[i*dim : (i+1)*dim]
		values := make([]float64, dim)
		for j, char := range chunk {
			val, ok := alphabet[string(char)]
			if !ok {
				return nil, fmt.Errorf("invalid character in text: %c", char)
			}
			values[j] = float64(val)
		}
		vectors[i] = mat.NewDense(dim, 1, values)
	}
	return vectors, nil
}

func vectorsToText(vectors []*mat.Dense) (string, error) {
	var builder strings.Builder
	for _, v := range vectors {
		r, c := v.Dims()
		if c != 1 {
			return "", fmt.Errorf("expected column vectors (nx1), but got %dx%d", r, c)
		}
		for i := range r {
			val := int(math.Round(v.At(i, 0)))
			char, ok := reverseAlphabet[val]
			if !ok {
				return "", fmt.Errorf("invalid numerical value in vector: %d", val)
			}
			builder.WriteString(char)
		}
	}
	return builder.String(), nil
}

func padMessage(message string, blockSize int) string {
	if len(message)%blockSize == 0 {
		return message
	}
	paddingNeeded := blockSize - (len(message) % blockSize)
	return message + strings.Repeat("X", paddingNeeded)
}

func applyMod26ToMatrices(matrices []*mat.Dense) {
	for _, m := range matrices {
		r, c := m.Dims()
		for i := range r {
			for j := range c {
				val := m.At(i, j)
				mod := int(math.Round(val)) % 26
				if mod < 0 {
					mod += 26
				}
				m.Set(i, j, float64(mod))
			}
		}
	}
}

func calculateInverseKeyMatrix(m *mat.Dense) (*mat.Dense, error) {
	det := mat.Det(m)
	modDet := int(math.Round(det)) % 26
	if modDet < 0 {
		modDet += 26
	}
	fmt.Printf("Determinant of key matrix: %.2f | mod 26 -> %d\n", det, modDet)

	if modDet == 0 {
		return nil, errors.New("matrix is not invertible as determinant is 0 mod 26")
	}

	detInverse, err := calculateModInverse(modDet, 26)
	if err != nil {
		return nil, fmt.Errorf("cannot find modular inverse of determinant: %w", err)
	}
	fmt.Printf("Multiplicative inverse of determinant (%d) is: %d\n", modDet, detInverse)

	adjugate, err := calculateAdjugate(m)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate adjugate matrix: %w", err)
	}

	r, c := m.Dims()
	inverseMatrix := mat.NewDense(r, c, nil)
	inverseMatrix.Scale(float64(detInverse), adjugate)

	applyMod26ToMatrices([]*mat.Dense{inverseMatrix})

	return inverseMatrix, nil
}

func calculateAdjugate(m *mat.Dense) (*mat.Dense, error) {
	r, c := m.Dims()
	if r != c {
		return nil, errors.New("matrix must be square to find its adjugate")
	}

	adjugate := mat.NewDense(r, c, nil)
	if r == 1 {
		adjugate.Set(0, 0, 1)
		return adjugate, nil
	}

	for i := range r {
		for j := range c {
			subMatrixData := make([]float64, (r-1)*(c-1))
			subI := 0
			for origI := range r {
				if origI == i {
					continue
				}
				for origJ := range c {
					if origJ == j {
						continue
					}
					subMatrixData[subI] = m.At(origI, origJ)
					subI++
				}
			}
			subMatrix := mat.NewDense(r-1, c-1, subMatrixData)
			minor := mat.Det(subMatrix)
			cofactor := math.Pow(-1, float64(i+j)) * minor
			adjugate.Set(j, i, cofactor)
		}
	}
	return adjugate, nil
}

func calculateModInverse(n, modulus int) (int, error) {
	for i := 1; i < modulus; i++ {
		if (n*i)%modulus == 1 {
			return i, nil
		}
	}
	return 0, fmt.Errorf("no modular inverse for %d mod %d", n, modulus)
}

func printMatrix(m *mat.Dense, label string) {
	if m == nil {
		fmt.Printf("%s\n    (nil matrix)\n", label)
		return
	}
	fmt.Println(label)
	matrixString := fmt.Sprintf("%v", mat.Formatted(m, mat.Squeeze()))
	lines := strings.Split(strings.TrimRight(matrixString, "\n"), "\n")
	for _, line := range lines {
		fmt.Printf("    %s\n", line)
	}
}
