package samchelper

import (
	"encoding/json"
	"io/ioutil"
)

// ResolveNodeID resolves ethereum peerID into nodeID
func ResolveNodeID(peerID string) int {
	var nodes = map[string]int{}
	nodesContent, err := ioutil.ReadFile(nodesFile)
	if err != nil {
		return -1
	}
	err = json.Unmarshal(nodesContent, &nodes)
	if err != nil {
		return -1
	}

	nodeID, ok := nodes[peerID]
	if ok {
		return nodeID
	}
	return -1
}
