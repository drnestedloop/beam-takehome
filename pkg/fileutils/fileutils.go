package fileutils

import (
	"encoding/base64"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"time"

	"slai.io/takehome/pkg/common"
)

func FileSerializer(path string) (string, error) {
	data, err := os.ReadFile(path)

	if err != nil {
		return "", err
	}

	return base64.RawStdEncoding.EncodeToString(data), nil
}

func snapshot (dir string) (map[string]os.FileInfo, error) {
	// take snapshot of directory
	snap := make(map[string]os.FileInfo)
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			// skip directories
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		relpath, err := filepath.Rel(dir, path)
		
		if err != nil {
			return err
		}
		snap[relpath] = info
		return nil
	})

	if err != nil {
		log.Printf("Encountered error during snapshot: %v", err)
		return nil, err
	}
	return snap, nil
}

// CleanDirs deletes all empty directories in path up to (but not including) the root directory
func CleanDirs(root string, path string) {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		log.Printf("Error obtaining absolute filepath: %v", err)
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		log.Printf("Error obtaining absolute filepath: %v", err)
	}
	for {
		if absPath == absRoot || absPath == "." {
			break
		}
		entries, err := os.ReadDir(absPath)
		if err != nil {
			log.Printf("Error reading directories while cleaning: %v", err)
		}
		if len(entries) != 0 {
			break
		}
		log.Printf("Cleaning up %s", absPath)
		if err := os.Remove(absPath); err !=nil {
			break
		}
		absPath = filepath.Dir(absPath)
	}



	
}
// WatchDir take a directory string, and continuously polls for changes 
// and sends a slice of common.FileOperations over the updates channel
func WatchDir (dir string, freq float64, updatesChan chan <- []common.FileOperation) {
		// DEBGUG	
		log.Println("inside WatchDir")

		interval := time.Duration(freq * float64(time.Second))
		ticker := time.NewTicker(interval)
		prev := make(map[string]os.FileInfo)

		for range ticker.C {
			updates := []common.FileOperation{} // used to store names of updated files
			curr, err := snapshot(dir)
			if err != nil {
				log.Printf("Error %v\n", err)
				continue
			}
			// check to see if file is created or updated
			for name, info := range curr {
				old, exists := prev[name]

				if !exists {
					// file is new, add its name to the updates slice
					updates = append(updates, common.FileOperation{OpType: common.CREATE, FileName: name})
				} else if old.ModTime() != info.ModTime() || old.Size() != info.Size() {
					// file has likely been updated
					updates = append(updates, common.FileOperation{OpType: common.UPDATE, FileName: name})
				}
			}
			// check for deletions
			for name, _ := range prev {
				_, exists := curr[name]

				if !exists {
					// file has been deleted
					updates = append(updates, common.FileOperation{OpType: common.DELETE, FileName: name})
				}
			}
			// if there are updates send them on the channel
			if len(updates) != 0 {
				log.Printf("Sending updates of length %v\n", len(updates))
				updatesChan <- updates
			}
			prev = curr
		}

}

