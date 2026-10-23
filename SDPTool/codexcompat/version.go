// Package codexcompat defines admission to the bounded PG Codex wire profile.
// Admission is not a claim that every future host has been workflow-tested.
package codexcompat

import (
	"strconv"
	"strings"
)

const MinimumVersion = "0.164.0"
const Profile = "sdp-codex-wire/1"

// Supported admits stable numeric versions at or above the measured floor.
// Preview/build labels must receive an explicit compatibility disposition.
func Supported(version string) bool {
	parts := strings.Split(version, ".")
	if len(parts) != 3 {
		return false
	}
	var numbers [3]uint64
	for i, part := range parts {
		if part == "" || len(part) > 9 || len(part) > 1 && part[0] == '0' {
			return false
		}
		for _, r := range part {
			if r < '0' || r > '9' {
				return false
			}
		}
		n, err := strconv.ParseUint(part, 10, 32)
		if err != nil {
			return false
		}
		numbers[i] = n
	}
	return numbers[0] > 0 || numbers[1] >= 164
}
