package crypto

// HashCode generate hash code
func HashCode(key string) int {
	if len(key) == 0 {
		return 0
	}

	hash := 0
	chars := []byte(key)
	for _, char := range chars {
		// Better decentralized hash
		// s[0]*31^(n-1) + s[1]*31^(n-2) + ... + s[n-1]
		hash = 31*hash + int(char)
	}

	return hash
}
