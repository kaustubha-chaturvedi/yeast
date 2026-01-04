package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"io"
	"os"
	"runtime"
)

func GetDeviceKey() []byte {
	var deviceID string

	switch runtime.GOOS {
	case "linux":
		if id, err := os.ReadFile("/etc/machine-id"); err == nil {
			deviceID = string(id)
		} else if id, err := os.ReadFile("/var/lib/dbus/machine-id"); err == nil {
			deviceID = string(id)
		} else if hostname, err := os.Hostname(); err == nil {
			deviceID = hostname
		} else {
			deviceID = "default-device"
		}
	case "darwin":
		if id, err := os.ReadFile("/etc/hostid"); err == nil {
			deviceID = string(id)
		} else if hostname, err := os.Hostname(); err == nil {
			deviceID = hostname
		} else {
			deviceID = "default-device"
		}
	case "windows":
		if hostname, err := os.Hostname(); err == nil {
			deviceID = hostname
		} else {
			deviceID = "default-device"
		}
	default:
		if hostname, err := os.Hostname(); err == nil {
			deviceID = hostname
		} else {
			deviceID = "default-device"
		}
	}

	key := sha256.Sum256([]byte(deviceID + "yst-index-cache"))
	return key[:32]
}

func Encrypt(data []byte) ([]byte, error) {
	key := GetDeviceKey()
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, HandleError("create cipher", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, HandleError("create GCM", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, HandleError("generate nonce", err)
	}

	return gcm.Seal(nonce, nonce, data, nil), nil
}

func Decrypt(data []byte) ([]byte, error) {
	key := GetDeviceKey()
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, HandleError("create cipher", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, HandleError("create GCM", err)
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, HandleErrorf("decrypt", "ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	result, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, HandleError("decrypt data", err)
	}
	return result, nil
}

