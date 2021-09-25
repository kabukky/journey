package watcher

import (
	"github.com/Landria/journey/helpers"
	"gopkg.in/fsnotify.v1"
	"log"
	"os"
	"path/filepath"
)

var watcher *fsnotify.Watcher
var watchedDirectories []string

func Watch(paths []string, extensionsFunctions map[string]func() error) error {
	// Prepare watcher to generate the theme on changes to the files
	if watcher == nil {
		var err error
		watcher, err = createWatcher(extensionsFunctions)
		if err != nil {
			return err
		}
	} else {
		// Remove all current directories from watcher
		for _, dir := range watchedDirectories {
			err := watcher.Remove(dir)
			if err != nil {
				return err
			}
		}
	}
	watchedDirectories = make([]string, 0)
	// Watch all subdirectories in the given paths
	for _, path := range paths {
		err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
			if info.IsDir() {
				err := watcher.Add(filePath)
				if err != nil {
					return err
				}
				watchedDirectories = append(watchedDirectories, filePath)
			}
			return nil
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func createWatcher(extensionsFunctions map[string]func() error) (*fsnotify.Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write {
					for key, value := range extensionsFunctions {
						if !helpers.IsDirectory(event.Name) && filepath.Ext(event.Name) == key {
							// Call the function associated with this file extension
							err := value()
							if err != nil {
								log.Panic("Error while reloading theme or plugins:", err)
							}
						}
					}
				}
			case err := <-watcher.Errors:
				log.Println("Error while watching directory:", err)
			}
		}
	}()
	return watcher, nil
}
