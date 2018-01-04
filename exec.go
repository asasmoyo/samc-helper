package samchelper

import "log"

// ExecAt executes fn if current node is equals to nodeID
func ExecAt(nodeID, fnName string, fn func() error) {
	if currNodeID != nodeID {
		log.Printf("[SAMC-HELPER] skip executing [%s]. currNode: [%s] wantNode: [%s]\n", fnName, currNodeID, nodeID)
		return
	}

	err := fn()
	if err != nil {
		log.Printf("[SAMC-HELPER] got an error executing [%s]. err: [%s]\n", fnName, err.Error())
	} else {
		log.Println("[SAMC-HELPER] successfully executed", fnName)
	}
}
