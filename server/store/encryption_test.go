package store

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"math/rand/v2"
	"slices"
	"testing"
)

func TestEncryption_getKeyParts(t *testing.T) {
	tests := []struct {
		ek                  string // base64
		numParts, threshold uint32
	}{
		{
			ek:        "cud9TfePy0XFjt4sPalHJXcHhHY8siSdZOE+84/w5b8=",
			numParts:  5,
			threshold: 2,
		},
		{
			ek:        "ADStQYQ4IgOQb6PuWyEOJ1iX7K81lLokiSkUEU4ugM0=",
			numParts:  2,
			threshold: 2,
		},
		{
			ek:        "kmTIbbj7jcyjNGPoSh1yzWc8FcGZm8z5TkZEhIjBy48=",
			numParts:  20,
			threshold: 10,
		},
	}
	encryption := storeEncryption{}
	for i, test := range tests {
		t.Run(fmt.Sprintf("getKeyParts_%d", i), func(t *testing.T) {
			ek, err := base64.StdEncoding.DecodeString(test.ek)
			if err != nil {
				t.Fatal(err)
			}
			keyParts, err := encryption.getKeyParts(ek, test.numParts, test.threshold)
			if err != nil {
				t.Fatal(err)
			}

			if int(test.numParts) != len(keyParts) {
				t.Fatalf("expected %d parts, got %d", test.numParts, len(keyParts))
			}
		})
	}
}

func TestEncryption_unseal(t *testing.T) {
	tests := []struct {
		threshold   int
		expectedKey string // base64
		keyParts    []string
	}{
		{
			expectedKey: "cud9TfePy0XFjt4sPalHJXcHhHY8siSdZOE+84/w5b8=",
			threshold:   2,
			keyParts: []string{
				"f18d1d4d7c4a652f46b7b8ab59e5f201cca8a8418c52e6b84a36582694f782006e",
				"21fd0b4dd455055f96cb8c47734c0bc66fc717674685d470ebc76cc90ddab99f22",
				"92f9564dda4e685b25a7e627fdbd2cd7c2d8bb81c0fd8a13a35906b3479fa1e1f5",
				"fd1bd94d96f013b94a69bbd4ee015bc7448488be0c74c124deef5b4b4b6adbfb8b",
				"5abed34dbd3ee41ced7cd435a2c79c876aef44421166f6f8310134867b93289368",
			},
		},
		{
			expectedKey: "cud9TfePy0XFjt4sPalHJXcHhHY8siSdZOE+84/w5b8=",
			threshold:   2,
			keyParts: []string{
				"af5e80cf1f7c4b711dcf1215d2244f37b9cd3255b17b12cbcc9129c0c93931046e",
				"461c96cbad779acda2cdee86e3da361090558a984fd17bb89b7958ff4893469369",
				"cd20a89d6d6878dc00e6e858bef7d0732294cdc362feded92a6a1273febc711b8c",
				"6dd37aa1b901ab529ff88b78ba8a41a5aedd7fa41fa2bf2e1ac5fb2d30e0ba07a1",
				"850b5f86149eb208146697b94a08b9708c8302489ddddb0636624d2af39fb5e157",
			},
		},
		{
			expectedKey: "kmTIbbj7jcyjNGPoSh1yzWc8FcGZm8z5TkZEhIjBy48=",
			threshold:   10,
			keyParts: []string{
				"9c054df70f465f69a16c91c027ffe163cce75e7100f3fa7bd70511549f18deb923",
				"52a78d900964e8ff411fd9ed961c38278a63da5edd84456ad71a4061bb2bb5d726",
				"ff8be1c67ed679bac2ef427bded20474c91ba6fbd8522901a1be11b5cf92d20b32",
				"f7e59007cea4467575ede8c08e8615ab1079e24d9452b02a2b5c744bec84f0aff1",
				"2cceee9bd6973a7ca810b3258314c319087411aceb59b8fe10fbd57e411cc7f08d",
				"702b44fe763a9f6090978d4e22df747f8334d69fb503d71277bad0fee2b57fdc65",
				"19065017ba759243e14a1da377c9439b7fb2ca77691cae1ffca10d73192c9b5a12",
				"a4e077ea599cf33b44471dc60e270ffeaddb460a25551ff010d73fff3499ff8b15",
				"e5accded05846fb5ccea55a988d2a92b65a551825a3bc7683677929dc4b80c4d92",
				"f3681c76650f3a4f5b0e545e2d1662c0420a443e6d032131ddbea5677fd6579551",
				"0c132a564da7511d59d57f9fc4f884c8b517a60f4526988adda69676f25e6591c5",
				"54706084452e16b09f597b90cf760ed2bac1b80952de92e2249963ae72a2c1f561",
				"e44bb3bdc49990bbbba96c72bc46ad684c08aac3a9c6f8ee378e52a491f894f1d3",
				"eadb8db7e560d4c07d3220bf017808e46b35f100c0b8df2c360eae010f06bf9c3a",
				"825b014cf524b2d2b85dae47937b8def32b7a839887c666e1afed1983b030160d8",
				"5547e678a4b8a61ec11c2f8b6e0aa4748acbed41b03c1a4b5c8d066c8780367055",
				"4ee2122c607531450ee0f6a9bcbb3fb179c80ffa4e00063a0455aebb2042b02149",
				"f512405636bc2e77c151f94033be7e0963dcbb27c6abb0ffb8b9690a69d76586a1",
				"c74d00e8acb85ad2be568e44473b71a4bfde35f74a1d2cd3bd39a172fecd7d3d05",
				"92643bd73919f8c77c8325f927b46d3d657ac2108aaceab34b4b1f67df4e185417",
			},
		},
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("getKeyParts_%d", i), func(t *testing.T) {
			expectedKey, err := base64.StdEncoding.DecodeString(test.expectedKey)
			if err != nil {
				t.Fatal(err)
			}

			numKeyParts := len(test.keyParts)
			for range 5 * numKeyParts {
				encryption := storeEncryption{}

				used := make([]int, 0, test.threshold)
				count := 0
				for range test.threshold * 5 {
					idx := rand.IntN(numKeyParts)
					if slices.Contains(used, idx) {
						continue
					}
					used = append(used, idx)
					count = count + 1

					key, err := encryption.unseal(test.keyParts[idx], false)
					if count < test.threshold && bytes.Equal(key, expectedKey) && err == nil {
						t.Errorf("Not reached expected threshold %d before unsealed", test.threshold)
						break
					} else if count == test.threshold && err != nil {
						t.Errorf("Reached expected threshold %d for unsealing, but got unexpected err: %v", test.threshold, err)
						break
					} else if count == test.threshold && err == nil {
						if !bytes.Equal(key, expectedKey) {
							t.Errorf("expected key %s but got key %s", test.expectedKey, base64.StdEncoding.EncodeToString(key))
						}
						break
					}
				}
			}
		})
	}
}
