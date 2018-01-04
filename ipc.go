package samchelper

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
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
	nodesContent, err := ioutil.ReadFile(nodesFile)
	if err != nil {
		return "", err
	}
	var nodes map[string]int
	err = json.Unmarshal(nodesContent, &nodes)
	if err != nil {
		return "", err
	}
	if _, ok := nodes[toID]; !ok {
		return "", errors.New("cannot find node number of receiving node")
	}
	data["receiving_node"] = nodes[toID]

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

	hasher := md5.New()
	hasher.Write(buf.Bytes())
	fHash := hex.EncodeToString(hasher.Sum(nil))

	// compute hash id
	iHasher := fnv.New32()
	iHasher.Write(buf.Bytes())
	hashID := int(iHasher.Sum32())
	data["event_id"] = hashID
	buf.WriteString(fmt.Sprintf("event_id=%d\n", hashID))

	path := computePath(filepath.Join("new", ipcPrefix+fHash))
	log.Println("[SAMC-HELPER] writing new ipc", ipcPrefix+fHash, data, "at", path)
	err = ioutil.WriteFile(path, buf.Bytes(), os.ModePerm)
	if err != nil {
		return "", err
	}

	return fHash, err
}

// SendIPC moves new ipc file located at subPath into send folder
func SendIPC(fileName string) error {
	oldPath := computePath(filepath.Join("new", ipcPrefix+fileName))
	newPath := computePath(filepath.Join("send", ipcPrefix+fileName))

	content, err := ioutil.ReadFile(oldPath)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(newPath, content, os.ModePerm)
	if err != nil {
		return err
	}
	return os.RemoveAll(oldPath)
}

// WaitForAckIPC waits until given fileName exists inside ack folder
func WaitForAckIPC(fileName string) error {
	path := computePath(filepath.Join("ack", ipcPrefix+fileName))
	log.Println("[SAMC-HELPER] waiting ack for " + ipcPrefix + fileName + " at " + path)
	for _, err := os.Stat(path); os.IsNotExist(err); _, err = os.Stat(path) {
		time.Sleep(1 * time.Second)
	}
	log.Println("[SAMC-HELPER] got ack for " + ipcPrefix + fileName + " at " + path)
	return nil
}
