package knx

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"sync"
	"time"

	"github.com/pakerfeldt/knx-mqtt/internal/dpt"
	"github.com/pakerfeldt/knx-mqtt/internal/models"
	"github.com/pakerfeldt/knx-mqtt/internal/msg"
	"github.com/pakerfeldt/knx-mqtt/internal/utils"
	"github.com/rs/zerolog/log"
	knxgo "github.com/vapourismo/knx-go/knx"
)

// KNXLogger handles logging of KNX messages to a file with rotation capabilities
type KNXLogger struct {
	config     models.KNXLogConfig
	file       *os.File
	mu         sync.Mutex
	fileSize   int64
	lastRotate time.Time
	knxItems   *models.KNX
}

// KNXLogEntry represents a log entry for a KNX message
type KNXLogEntry struct {
	Timestamp    time.Time   `json:"timestamp"`
	Direction    string      `json:"direction"`
	Destination  string      `json:"destination"`
	Command      string      `json:"command"`
	Data         string      `json:"data,omitempty"`
	Name         string      `json:"name,omitempty"`
	DecodedValue interface{} `json:"decoded_value,omitempty"`
	Unit         string      `json:"unit,omitempty"`
}

// NewKNXLogger creates a new KNX logger
func NewKNXLogger(config models.KNXLogConfig, knxItems *models.KNX) (*KNXLogger, error) {
	if !config.Enabled {
		return nil, nil
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(config.File)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	// Open log file
	file, err := os.OpenFile(config.File, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	// Get file info for size
	fileInfo, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	return &KNXLogger{
		config:     config,
		file:       file,
		fileSize:   fileInfo.Size(),
		lastRotate: time.Now(),
		knxItems:   knxItems,
	}, nil
}

// LogIncoming logs an incoming KNX message
func (l *KNXLogger) LogIncoming(message *msg.KNXMessage) error {
	if l == nil {
		return nil
	}

	entry := KNXLogEntry{
		Timestamp:   time.Now(),
		Direction:   "incoming",
		Destination: message.Destination(),
		Command:     message.Command(),
		Data:        fmt.Sprintf("%X", message.Data()),
	}

	if message.IsResolved() {
		entry.Name = message.Name()

		// Add decoded value and unit
		if l.knxItems != nil {
			// Get the datapoint type from the message
			dpType := message.Datapoint()

			// Extract the decoded value
			if dp, ok := dpt.Produce(dpType); ok {
				if err := dp.Unpack(message.Data()); err == nil {
					entry.DecodedValue = utils.ExtractDatapointValue(dp, dpType)
					entry.Unit = dp.Unit()
				}
			}
		}
	}

	return l.writeLogEntry(entry)
}

// LogOutgoing logs an outgoing KNX message
func (l *KNXLogger) LogOutgoing(event knxgo.GroupEvent) error {
	if l == nil {
		return nil
	}

	commandStr := utils.KNXCommandToString(event.Command)

	entry := KNXLogEntry{
		Timestamp:   time.Now(),
		Direction:   "outgoing",
		Destination: event.Destination.String(),
		Command:     commandStr,
		Data:        fmt.Sprintf("%X", event.Data),
	}

	// Try to resolve the group address and add decoded value and unit
	if l.knxItems != nil {
		flatAddr := models.FlatGroupAddress(event.Destination)
		if index, exists := l.knxItems.GadToIndex[flatAddr]; exists {
			groupAddress := l.knxItems.GroupAddresses[index]
			entry.Name = groupAddress.Name

			// Get the datapoint type and decode the value
			dpType := groupAddress.Datapoint
			if dp, ok := dpt.Produce(dpType); ok {
				if err := dp.Unpack(event.Data); err == nil {
					entry.DecodedValue = utils.ExtractDatapointValue(dp, dpType)
					entry.Unit = dp.Unit()
				}
			}
		}
	}

	return l.writeLogEntry(entry)
}

// writeLogEntry writes a log entry to the file
func (l *KNXLogger) writeLogEntry(entry KNXLogEntry) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Check if we need to rotate the log
	if err := l.checkRotation(); err != nil {
		return err
	}

	var data []byte
	var err error

	if l.config.Format == "json" {
		data, err = json.Marshal(entry)
		if err != nil {
			return fmt.Errorf("failed to marshal log entry: %w", err)
		}
		data = append(data, '\n')
	} else {
		// Text format
		data = fmt.Appendf(nil, "[%s] %s %s %s %s %s %v %s\n",
			entry.Timestamp.Format(time.RFC3339),
			entry.Direction,
			entry.Destination,
			entry.Command,
			entry.Data,
			entry.Name,
			entry.DecodedValue,
			entry.Unit,
		)
	}

	n, err := l.file.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write to log file: %w", err)
	}

	l.fileSize += int64(n)
	return nil
}

// checkRotation checks if the log file needs to be rotated
func (l *KNXLogger) checkRotation() error {
	// Check if we need to rotate based on size
	if l.config.MaxSize > 0 && l.fileSize >= int64(l.config.MaxSize*1024) {
		log.Info().Msg("Log file exceeded max size: rotating")
		if err := l.rotate(); err != nil {
			return err
		}
	}

	// Check if we need to rotate based on age
	if l.config.MaxAge > 0 && time.Since(l.lastRotate) > l.config.MaxAge {
		log.Info().Msg("Log file exceeded max age: rotating")
		if err := l.rotate(); err != nil {
			return err
		}
	}

	return nil
}

// rotate rotates the log file
func (l *KNXLogger) rotate() error {
	// Close current file
	if err := l.file.Close(); err != nil {
		return fmt.Errorf("failed to close log file: %w", err)
	}

	// Rename current file with timestamp
	// If this format ever gets changed here, don't forget to also change it in the regexp in cleanupOldFiles
	timestamp := time.Now().Format("20060102-150405")
	rotatedName := fmt.Sprintf("%s.%s", l.config.File, timestamp)
	if err := os.Rename(l.config.File, rotatedName); err != nil {
		return fmt.Errorf("failed to rename log file: %w", err)
	}

	// Open new file
	file, err := os.OpenFile(l.config.File, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open new log file: %w", err)
	}
	l.file = file
	l.fileSize = 0
	l.lastRotate = time.Now()

	// Compress old file if needed
	if l.config.Compress {
		if err := compressFile(rotatedName); err != nil {
			log.Error().Err(err).Str("file", rotatedName).Msg("Failed to compress rotated log file")
		}
	}

	// Clean up old files if needed
	if l.config.MaxFiles > 0 {
		if err := cleanupOldFiles(l.config.File, l.config.MaxFiles); err != nil {
			log.Error().Err(err).Str("file", l.config.File).Msg("Failed to clean up old log files")
		}
	}

	return nil
}

// Close closes the logger
func (l *KNXLogger) Close() error {
	if l.file == nil {
		return nil
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	return l.file.Close()
}

// compressFile compresses a file using gzip
func compressFile(filename string) error {
	// Open the original file for reading
	srcFile, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open source file for compression: %w", err)
	}
	defer srcFile.Close()

	// Create the destination file with .gz extension
	destFilename := filename + ".gz"
	destFile, err := os.Create(destFilename)
	if err != nil {
		return fmt.Errorf("failed to create compressed file: %w", err)
	}

	// Create a gzip writer
	gzipWriter := gzip.NewWriter(destFile)
	defer func() {
		// Close the gzip writer properly
		if err := gzipWriter.Close(); err != nil {
			log.Error().Err(err).Str("file", destFilename).Msg("Failed to close gzip writer")
		}

		// Close the destination file
		if err := destFile.Close(); err != nil {
			log.Error().Err(err).Str("file", destFilename).Msg("Failed to close compressed file")
		}
	}()

	// Copy the contents from source to gzip writer
	if _, err := io.Copy(gzipWriter, srcFile); err != nil {
		return fmt.Errorf("failed to compress file: %w", err)
	}

	// Flush any pending compressed data
	if err := gzipWriter.Flush(); err != nil {
		return fmt.Errorf("failed to flush compressed data: %w", err)
	}

	// Remove the original file after successful compression
	if err := os.Remove(filename); err != nil {
		return fmt.Errorf("failed to remove original file after compression: %w", err)
	}

	return nil
}

// cleanupOldFiles removes old rotated log files
func cleanupOldFiles(baseFilename string, maxFiles int) error {
	dir := filepath.Dir(baseFilename)
	base := filepath.Base(baseFilename)

	// List all files in the directory
	files, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	// Find rotated log files
	var logFiles []string
	for _, file := range files {
		// Check if this is a rotated log file:
		// 1. Not a directory
		// 2. Not the base file itself
		// 3. Either a compressed version or a rotated version with timestamp
		fileName := filepath.Base(file.Name())
		if !file.IsDir() && fileName != base {
			// Check if it's a compressed version (base.gz)
			if fileName == base+".gz" {
				logFiles = append(logFiles, filepath.Join(dir, file.Name()))
				continue
			}

			// Check if it's a rotated version (base.timestamp or base.timestamp.gz)
			// Timestamp format is YYYYMMDD-HHMMSS as defined in rotate() function
			if matched, _ := regexp.MatchString(fmt.Sprintf(`^%s\.\d{8}-\d{6}(\.gz)?$`, regexp.QuoteMeta(base)), fileName); matched {
				logFiles = append(logFiles, filepath.Join(dir, file.Name()))
			}
		}
	}

	if len(logFiles) <= maxFiles {
		return nil
	}

	// Sort files by modification time (oldest first)
	type fileInfo struct {
		path    string
		modTime time.Time
	}

	fileInfos := make([]fileInfo, 0, len(logFiles))
	for _, file := range logFiles {
		info, err := os.Stat(file)
		if err != nil {
			return fmt.Errorf("failed to get file info: %w", err)
		}
		fileInfos = append(fileInfos, fileInfo{path: file, modTime: info.ModTime()})
	}

	// Sort by modification time (oldest first)
	sort.Slice(fileInfos, func(i, j int) bool {
		return fileInfos[i].modTime.Before(fileInfos[j].modTime)
	})

	// Remove excess files
	for i := range len(fileInfos) - maxFiles {
		if err := os.Remove(fileInfos[i].path); err != nil {
			return fmt.Errorf("failed to remove old log file: %w", err)
		}
		log.Info().
			Str("file", fileInfos[i].path).
			Int("max_files", maxFiles).
			Msg("Removed old log file as part of rotation cleanup")
	}

	// Log information about the files we're keeping
	if len(fileInfos) > maxFiles {
		keptFiles := fileInfos[len(fileInfos)-maxFiles:]
		for _, file := range keptFiles {
			log.Debug().
				Str("file", file.path).
				Time("modified", file.modTime).
				Int("max_files", maxFiles).
				Msg("Keeping log file after rotation cleanup")
		}
	}

	return nil
}
