package library

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

/* this is the entry like for example folder, or files*/
type Entry struct {
	Name  string
	Path  string // tthis will store the full path
	IsDir bool
	Size  int64 // 0 for directories
}

var audioExts = map[string]bool{
	".mp3":  true,
	".aac":  true,
	".ac3":  true,
	".aiff": true,
	".aif":  true,
	".alac": true,
	".amr":  true,
	".ape":  true,
	".au":   true,
	".dts":  true,
	".flac": true,
	".m4a":  true,
	".m4b":  true,
	".m4p":  true,
	".mka":  true,
	".mp2":  true,
	".mpa":  true,
	".oga":  true,
	".ogg":  true,
	".opus": true,
	".ra":   true,
	".ram":  true,
	".tta":  true,
	".voc":  true,
	".wav":  true,
	".weba": true,
	".wma":  true,
	".wv":   true,
	".3ga":  true,
}

/*checks if the file is audio file or not*/
func isAudioFile(name string) bool {
	return audioExts[strings.ToLower(filepath.Ext(name))]
}

func readDir(dir string) ([]Entry, error) {
	os_entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	// creates slice of entry struct with 0 size and capacity of entries in the current dir
	entries := make([]Entry, 0, len(os_entries))

	for _, entry := range os_entries {
		name := entry.Name()

		// ignore the hidden files
		if strings.HasPrefix(name, ".") {
			continue
		}

		full := filepath.Join(dir, name)

		/*get the info about the file pointed by full and if theres error just skip that shit
		*info.Name()        // "song.mp3"
		*info.IsDir()       // false
		*info.Size()        // file size
		*info.ModTime()     // modification time
		* */
		info, err := os.Stat(full)
		if err != nil {
			continue // broken symlink or perms error
		}

		if info.IsDir() {
			entries = append(entries, Entry{
				Name:  name,
				Path:  full,
				IsDir: true,
				Size:  0,
			})
		}

		if !isAudioFile(name) {
			continue
		}

		entries = append(entries, Entry{
			Name:  name,
			Path:  full,
			IsDir: false,
			Size:  info.Size(),
		})
	}

	// sort time
	sort.SliceStable(entries, func(i, j int) bool {
		a, b := entries[i], entries[j]

		if a.IsDir != b.IsDir {
			return a.IsDir
		}

		return strings.ToLower(a.Name) < strings.ToLower(b.Name)
	})

	return entries, nil
}
