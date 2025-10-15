# Hill Cipher in Go

This is a simple command-line program written in Go that demonstrates the Hill Cipher, a classic method of encryption. It allows you to encrypt and decrypt messages using a key matrix of any size (like 2x2 or 3x3 or 4*4).

The program is designed to be educational, showing the step-by-step matrix calculations involved in the process, run the app and provide a key and a message and it will Encrypt/Decrypt the message step by step.

## Features

*   Encrypts alphabetic messages.
*   Decrypts messages back to the original text.
*   Supports custom key matrices (e.g., 2x2, 3x3).
*   Shows detailed matrix calculations for both encryption and decryption.
*   Automatically cleans your message by removing spaces, numbers, and punctuation before processing.

## Getting Started

### Prerequisites

You need to have **Go** (version 1.18 or newer) installed on your computer.

### How to Run

1.  **Download the Code:**
    Clone this repository or download the `main.go` file to your computer.

2.  **Open Your Terminal:**
    Navigate to the directory where you saved the file.

3.  **Run the Program:**
    Execute the following command in your terminal:
    ```sh
    go run main.go
    ```
    The program will automatically download the required `gonum/mat` library the first time you run it.

You will then be prompted to enter a matrix dimension, a key, and the message you want to encrypt. You can validate the the Encryption/Decryption result using online tools.

## How the Hill Cipher Works

The Hill Cipher is a classic cipher that uses matrix multiplication to encrypt text.

1.  **Convert to Numbers:** Each letter in the message and key is turned into a number (A=0, B=1, ..., Z=25).
2.  **Create Vectors:** The message's numbers are grouped into small blocks (vectors) that match the size of the key matrix.
3.  **Encrypt:** Each block is multiplied by the key matrix.
4.  **Apply Modulo 26:** The results are taken modulo 26 to turn them back into numbers from 0 to 25.
5.  **Convert to Letters:** The final numbers are converted back into letters to form the encrypted message.

Decryption is done by following the same process but using the **modular multiplicative inverse** of the key matrix. For further information take a lot at the references.

## References

This implementation was created with the help of the following resources:

*   **Wikipedia - Hill Cipher:** [https://en.wikipedia.org/wiki/Hill_cipher](https://en.wikipedia.org/wiki/Hill_cipher)
*   **Interactive Mathematics - Hill Cipher:** [https://crypto.interactive-maths.com/hill-cipher.html](https://crypto.interactive-maths.com/hill-cipher.html)
   
The crypto.interactive website helped a lot with edge cases, for example filing the the provided key with "X", if it wasnt matching the size of the key.
