package junos

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
)

func (clt *Client) appendFakeCreateSetFile(lines []string) error {
	dirSetFile := path.Dir(clt.fakeCreateSetFile)
	if _, err := os.Stat(dirSetFile); err != nil {
		if err := os.MkdirAll(dirSetFile, os.FileMode(directoryPermission)); err != nil {
			return fmt.Errorf("creating parent directory of '%s': %w", clt.fakeCreateSetFile, err)
		}
	}
	f, err := os.OpenFile(clt.fakeCreateSetFile,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.FileMode(clt.filePermission))
	if err != nil {
		return fmt.Errorf("opening file '%s': %w", clt.fakeCreateSetFile, err)
	}
	defer f.Close()
	for _, v := range lines {
		if _, err := f.WriteString(v + "\n"); err != nil {
			return fmt.Errorf("writing in file '%s': %w", clt.fakeCreateSetFile, err)
		}
	}

	return nil
}

// log message in junosLogFile.
func (clt *Client) logFile(message string) {
	if clt.logFileDst != "" {
		f, err := os.OpenFile(clt.logFileDst,
			os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.FileMode(clt.filePermission))
		if err != nil {
			var perr *os.PathError
			if !errors.As(err, &perr) {
				log.Fatal(err)
			}
			log.Printf("[WARN] appending debug_netconf_log file: %s", perr.Error())
		}
		defer f.Close()
		logger := log.New(f, "", log.LstdFlags)

		logger.Printf("%s", message)
	}
}

func (clt *Client) FilePermission() int64 {
	return clt.filePermission
}
