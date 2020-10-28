// Build tool for watching and building source code.
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/pflag"
)

var (
	path string
)

func init() {
	pflag.StringVar(&path, "path", ".", "path of Go package to watch and build")
	pflag.Parse()
}

func main() {
	if err := run(); err != nil {
		fmt.Printf("error: %v\n", err)
	}

}

func run() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("initialising file watcher: %w", err)
	}
	defer watcher.Close()
	done := make(chan bool)
	go func() {
		run := func() *exec.Cmd {
			bin, err := build(path)
			if err != nil {
				log.Printf("error: building package: %s\n", err)
				return nil
			}
			handle, err := start(bin)
			if err != nil {
				log.Printf("error: starting bin: %s\n", err)
				return nil
			}
			return handle
		}
		var (
			handle = run()
			events = debounce(watcher.Events, time.Millisecond*250)
		)
		for {
			select {
			case event, ok := <-events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Printf("modified file: %s - restarting\n", event.Name)
					if handle != nil {
						if err := handle.Process.Kill(); err != nil {
							log.Printf("error: killing previous instance: %s\n", err)
						}
						if _, err := handle.Process.Wait(); err != nil {
							log.Printf("error: waiting for process to quit: %s\n", err)
						}
					}
					handle = run()
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Printf("error: %s\n", err)
			}
		}
	}()
	entries, err := directories(path)
	if err != nil {
		return fmt.Errorf("collecting directories: %w", err)
	}
	for _, dir := range entries {
		if err := watcher.Add(dir); err != nil {
			return fmt.Errorf("adding directory to watcher: %w", err)
		}
	}
	<-done
	return nil
}

func directories(root string) ([]string, error) {
	var paths []string
	entries, err := ioutil.ReadDir(root)
	if err != nil {
		return nil, fmt.Errorf("reading dir: %w", err)
	}
	for _, entry := range entries {
		path, err := filepath.Abs(filepath.Join(root, entry.Name()))
		if err != nil {
			return nil, fmt.Errorf("resolving absolute path: %w", err)
		}
		if entry.IsDir() {
			childEntries, err := directories(path)
			if err != nil {
				return nil, err
			}
			paths = append(paths, childEntries...)
		} else {
			paths = append(paths, path)
		}
	}
	return paths, nil
}

func debounce(input chan fsnotify.Event, period time.Duration) chan fsnotify.Event {
	ticker := time.NewTicker(period)
	output := make(chan fsnotify.Event)
	go func() {
		for range ticker.C {
			output <- <-input
			for ii := 0; ii <= len(input); ii++ {
				<-input
			}
		}
	}()
	return output
}

func build(path string) (string, error) {
	dir := os.TempDir()
	if out, err := exec.Command("go", "build", "-o", dir, path).CombinedOutput(); err != nil {
		return "", fmt.Errorf("%s: %s %w\n", path, string(out), err)
	}
	return fmt.Sprintf("%s.exe", filepath.Join(dir, filepath.Base(path))), nil
}

func start(bin string) (*exec.Cmd, error) {
	c := exec.Command(bin)
	c.Stdout = os.Stdout
	c.Stderr = os.Stdout
	if err := c.Start(); err != nil {
		return nil, err
	}
	return c, nil
}
