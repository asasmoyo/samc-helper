package samchelper

import (
	"log"
	"os"
)

var ipcDir, ipcPrefix, currNodeID, nodesFile string

// Init inits the samchelper package
// It sets:
//  - nodeID from NODE_ID env
//  - ipcDir from IPC_DIR env
//  - ipcPrefix from IPC_PREFIX env
//  - nodesFile from NODES_FILE
func Init() {
	ipcDir = os.Getenv("IPC_DIR")
	ipcPrefix = os.Getenv("IPC_PREFIX")
	currNodeID = os.Getenv("NODE_ID")
	nodesFile = os.Getenv("NODES_FILE")

	log.Printf("[SAMC-HELPER] Setting ipcDir: [%s] ipcPrefix: [%s] nodeID: [%s] nodesFile: [%s]\n", ipcDir, ipcPrefix, currNodeID, nodesFile)
}
