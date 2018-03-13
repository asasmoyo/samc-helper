package samchelper

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

func computePath(subPath string) string {
	return filepath.Join(ipcDir, subPath)
}

func preparePath(path string) {
	d := filepath.Dir(path)
	if info, _ := os.Stat(d); info.IsDir() {
		return
	}
	os.MkdirAll(d, os.ModePerm)
}

// ReadIPC reads ipc located at subPath
func ReadIPC(subPath string) (map[string]string, error) {
	path := computePath(subPath)
	f, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	data := map[string]string{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "=")
		if len(parts) != 2 {
			return nil, errors.New("invalid ipc file format")
		}

		data[parts[0]] = parts[1]
	}

	return data, nil
}

// NewIPC writes data into subPath
// The data will be written into ipcDir/subPath/md5(data)
// It uses NODES environment variable to determine node number from toID
func NewIPC(data map[string]int, toID string) (string, error) {
	nodeID, _ := strconv.ParseInt(os.Getenv("NODE_ID"), 10, 64)
	data["sender_node"] = int(nodeID)
	if toID == "" {
		data["receiving_node"] = int(nodeID)
	} else {
		data["receiving_node"] = ResolveNodeID(toID)
	}

	keys := []string{}
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	buf := bytes.NewBufferString("")
	for _, key := range keys {
		val := data[key]
		buf.WriteString(fmt.Sprintf("%s=%d\n", key, val))
	}

	// compute hash id
	iHasher := fnv.New32()
	iHasher.Write(buf.Bytes())
	hashID := int(iHasher.Sum32())
	data["event_id"] = hashID
	buf.WriteString(fmt.Sprintf("event_id=%d\n", hashID))

	fileName := fmt.Sprintf("%d-%d", hashID, time.Now().UnixNano())
	path := computePath(filepath.Join("new", ipcPrefix+fileName))
	log.Println("[SAMC-HELPER] writing new ipc from", nodeID, "at", path, "data", data)
	err := ioutil.WriteFile(path, buf.Bytes(), os.ModePerm)
	if err != nil {
		return "", err
	}

	return fileName, err
}

// SendIPC moves new ipc file located at subPath into send folder
func SendIPC(fileName string) error {
	oldPath := computePath(filepath.Join("new", ipcPrefix+fileName))
	newPath := computePath(filepath.Join("send", ipcPrefix+fileName))
	cmd := exec.Command("mv", oldPath, newPath)
	return cmd.Run()
}

// WaitForAckIPC waits until given fileName exists inside ack folder
func WaitForAckIPC(fileName string) {
	path := computePath(filepath.Join("ack", ipcPrefix+fileName))
	log.Println("[SAMC-HELPER] waiting ack for " + ipcPrefix + fileName + " at " + path)
	for _, err := os.Stat(path); os.IsNotExist(err); _, err = os.Stat(path) {
		time.Sleep(1 * time.Second)
	}
	log.Println("[SAMC-HELPER] got ack for " + ipcPrefix + fileName + " at " + path)
}
