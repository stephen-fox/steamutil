package naming

import (
	"hash/crc32"
	"strconv"
)

// LegacyNonSteamGameId returns the legacy ID for a non-Steam game.
//
// Based on work by Lucas Boppre Niehues:
// https://github.com/boppreh/steamgrid/blob/master/games.go#L117
func LegacyNonSteamGameId(gameName string, executablePath string) string {
	uniqueName := executablePath + gameName
	// Does IEEE CRC32 of target concatenated with gameName, then convert
	// to 64bit Steam ID. No idea why Steam chose this operation.
	top := uint64(crc32.ChecksumIEEE([]byte(uniqueName)) | 0x80000000)
	return strconv.FormatUint(top<<32|0x02000000, 10)
}
