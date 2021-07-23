package lib_test

import "encoding/base64"

var (
	words = []string{
		"keyboard",
		"horse",
		"berry",
		"frog",
		"cloud",
		"mouse",
		"horse", // ensure repeated words support
		"monitor",
		"laptop",
		"glass",
		"grass",
		"čÇÑñ",
		"birdÑÇ 你",
	}

	seed = "bartolo"

	// Encrypted words stored as Base64 values
	encryptedWordsB64 = []string{
		"d2Wl99MsC0d3VKj5ZF3RTg==",
		"TaGXBc0NnmzwtPoU9bBPiQ==",
		"oySmUtWNqrYQgwskrPhdXw==",
		"9c6w9+QTszWThKQByiNiLA==",
		"NJ+aB+c60eu1wZE5Xu9+Tg==",
		"3xNGwidzSN611fDEIa11aQ==",
		"TaGXBc0NnmzwtPoU9bBPiQ==",
		"AFmuz/Rj3I1fHK774waFGA==",
		"3LWjlsmR5w+xWk63JCeD/Q==",
		"fK6433h1aAFU0v0U09FJ3w==",
		"PasndDRJF+3N7Oo16y1Qvg==",
		"rBlTRw/dw4XdZAAQcKbxlA==",
		"OmdtLYQWu7V26kiTcGAxrA==",
	}

	encrypter = &EncrypterMock{
		EncryptWordsFunc: func([]string) ([][]byte, error) {
			var err error
			encrypted := make([][]byte, len(encryptedWordsB64))
			for i, ew := range encryptedWordsB64 {
				encrypted[i], err = base64.StdEncoding.DecodeString(ew)
				if err != nil {
					return nil, err
				}
			}
			return encrypted, nil
		},
	}
	decrypter = &DecrypterMock{
		DecryptWordsFunc: func(byte [][]byte) ([]string, error) {
			return words, nil
		},
	}
)
