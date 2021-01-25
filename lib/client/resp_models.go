package client

// Solc compiler
// SolcBuild and SolcVersion parse solc version build response
/* raw data
{
	"path": "soljson-v0.1.1+commit.6ff4cd6.js",
	"version": "0.1.1",
	"build": "commit.6ff4cd6",
	"longVersion": "0.1.1+commit.6ff4cd6",
	"keccak256": "0xd8b8c64f4e9de41e6604e6ac30274eff5b80f831f8534f0ad85ec0aff466bb25",
	"urls": [
		"bzzr://8f3c028825a1b72645f46920b67dca9432a87fc37a8940a2b2ce1dd6ddc2e29b",
		"dweb:/ipfs/QmPPGxsMtQSEUt9hn9VAtrtBTa1u9S5KF1myw78epNNFkx"
	]
}
*/
type SolcBuild struct {
	Path 		string 		`json:"path"`
	Version		string 		`json:"version"`
	Build		string 		`json:"build"`
	LongVersion string 		`json:"long_version"`
	Keccak256	string 		`json:"keccak_256"`
	Urls		[]string 	`json:"urls"`
}

type SolcVersion struct{
	Builds 			[]SolcBuild 		`json:"builds"`
	Releases		map[string]string	`json:"releases"`
	LatestRelease 	string 				`json:"latest_release"`
}


