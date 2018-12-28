package naming

import "testing"

func TestLegacyNonSteamGameId(t *testing.T) {
	name := LegacyNonSteamGameId("Pikmin", `"D:\Program Files\Dolphin\Dolphin.exe"`)

	expected := "11271507026838028288"
	if name != expected {
		t.Fatal("Did not get expected value of '" + expected + "' - got '" + name + "'")
	}
}
