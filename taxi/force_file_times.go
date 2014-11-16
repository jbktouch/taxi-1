package taxi

import "fmt"
import "path/filepath"
import "time"
import "os"

func ForceFileTimes(target string) {
	fmt.Printf("Setting time to July 6 1978 on %s\n", ".")
	filepath.Walk(target, ForceFileTimeOnFile)
}

func ForceFileTimeOnFile(path string, info os.FileInfo, err error) error {
	time := time.Date(1978, 7, 6, 5, 4, 3, 2, time.UTC)
	nerr := os.Chtimes(path, time, time)
	check(nerr)
	return nil
}
