package aceonline

import "DiscordBotGO/pkg/emoji"

var (
	ANImaps          = map[string]string{"ev": "Edmont Valley", "doa": "Desert of Ardor", "cc": "Crystal Cave", "pdm": "Plain of Doleful Melody", "hrs": "Herremeze Relic Site", "ab": "Atus Beach", "gr": "Gjert Road", "sp": "Slope Port", "pmc": "Portsmouth Canyon"}
	BCUmaps          = map[string]string{"bmc": "Bach Mountain Chain", "bs": "Blackburn Site", "zb": "Zaylope Beach", "sv": "Starlite Valley", "rl": "Redline", "kb": "Kahlua Beach", "nc": "Nubarke Cave", "op": "Orina Peninsula", "dr": "Daisy Riverhead"}
	ANImapsEmoji     = map[string]string{"ev": emoji.EdmontValley, "doa": emoji.DesertofArdor, "cc": emoji.CrystalCave, "pdm": emoji.PlainofDolefulMelody, "hrs": emoji.HerremezeRelicSite, "ab": emoji.AtusBeach, "gr": emoji.GjertRoad, "sp": emoji.SlopePort, "pmc": emoji.PortsmouthCanyon}
	BCUmapsEmoji     = map[string]string{"bmc": emoji.BachMountainChain, "bs": emoji.BlackburnSite, "zb": emoji.ZaylopeBeach, "sv": emoji.StarliteValley, "rl": emoji.Redline, "kb": emoji.KahluaBeach, "nc": emoji.NubarkeCave, "op": emoji.OrinaPeninsula, "dr": emoji.DaisyRiverhead}
	SortedANImapKeys = []string{"ev", "doa", "cc", "pdm", "hrs", "ab", "gr", "sp", "pmc"}
	SortedBCUmapKeys = []string{"bmc", "bs", "zb", "sv", "rl", "kb", "nc", "op", "dr"}
)
