package metrics

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/rahulkumarparida/roxkv/internal/store"
	"github.com/rahulkumarparida/roxkv/internal/utils"
)

type KeyNSize struct {
	Key  string
	Size int64
}

type SnapshotInfo struct {
	Name       string
	Path       string
	Size       int64
	ModifiedAt time.Time
}
type PersistenceHealth struct {
	DbFolderExists       bool
	SnapshotFolderExists bool
	SnapshotCount        int
	TotalSnapshotSize    int64
	LatestSnapshot       SnapshotInfo
}

type KeyMetadata struct {
	Key            string
	TTL            time.Time
	UpdatedAt      time.Time
	CreatedAt      time.Time
	LastAcessedBy  *utils.NewClient
	KeyAccessCount int64
	Size           int64
	Namespace      string
}

func liveItems(stre *store.MemoryAlloc) []store.Item {
	stre.Mu.RLock()
	defer stre.Mu.RUnlock()

	now := time.Now()
	items := make([]store.Item, 0, len(stre.Data))

	for _, item := range stre.Data {
		if !item.Meta.TTL.IsZero() && item.Meta.TTL.Before(now) {
			continue
		}
		items = append(items, item)
	}

	return items
}

func topSizedKeys(items []store.Item, topN int, largest bool) []KeyNSize {
	if len(items) == 0 || topN <= 0 {
		return []KeyNSize{}
	}

	sort.Slice(items, func(i, j int) bool {
		if largest {
			return items[i].Meta.Size > items[j].Meta.Size
		}
		return items[i].Meta.Size < items[j].Meta.Size
	})

	if topN > len(items) {
		topN = len(items)
	}

	keys := make([]KeyNSize, 0, topN)
	for idx := 0; idx < topN; idx++ {
		keys = append(keys, KeyNSize{
			Key:  items[idx].Key,
			Size: items[idx].Meta.Size,
		})
	}

	return keys
}

func snapshotFiles() []SnapshotInfo {
	dirPath := utils.SnapshotFolder()
	if dirPath == "" {
		return []SnapshotInfo{}
	}

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return []SnapshotInfo{}
	}

	snaps := make([]SnapshotInfo, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		snaps = append(snaps, SnapshotInfo{
			Name:       entry.Name(),
			Path:       filepath.Join(dirPath, entry.Name()),
			Size:       info.Size(),
			ModifiedAt: info.ModTime(),
		})
	}

	sort.Slice(snaps, func(i, j int) bool {
		return snaps[i].ModifiedAt.After(snaps[j].ModifiedAt)
	})

	return snaps
}

func toKeyMetadata(item store.Item) KeyMetadata {
	return KeyMetadata{
		Key:            item.Key,
		TTL:            item.Meta.TTL,
		UpdatedAt:      item.Meta.UpdatedAt,
		CreatedAt:      item.Meta.CreatedAt,
		LastAcessedBy:  item.Meta.LastAcessedBy,
		KeyAccessCount: item.Meta.KeyAccessCount,
		Size:           item.Meta.Size,
		Namespace:      item.Meta.Namespace,
	}
}

// Key statistcs
func GetTotalKeys(stre *store.MemoryAlloc, namespace *store.NameSpace) int {
	_ = namespace
	return len(liveItems(stre))
}

func GetLargestKeys(stre *store.MemoryAlloc, namespace *store.NameSpace, topN int) []KeyNSize {
	_ = namespace
	return topSizedKeys(liveItems(stre), topN, true)
}

func GetSmallestKeys(stre *store.MemoryAlloc, namespace *store.NameSpace, topN int) []KeyNSize {
	_ = namespace
	return topSizedKeys(liveItems(stre), topN, false)
}

func GetAverageValueSize(stre *store.MemoryAlloc, namespace *store.NameSpace) int {
	_ = namespace
	items := liveItems(stre)
	if len(items) == 0 {
		return 0
	}

	var totalSize int64
	for _, item := range items {
		totalSize += item.Meta.Size
	}

	return int(totalSize) / len(items)
}

func GetNamespaces(stre *store.MemoryAlloc, namespace *store.NameSpace) []KeyNSize {
	_ = namespace
	items := liveItems(stre)
	if len(items) == 0 {
		return []KeyNSize{}
	}

	namespaceCounts := make(map[string]int64)
	for _, item := range items {
		ns := item.Meta.Namespace
		if ns == "" {
			if strings.Contains(item.Key, ":") {
				ns = strings.Split(item.Key, ":")[0]
			} else {
				ns = item.Key
			}
		}
		namespaceCounts[ns] += 1
	}

	spaces := make([]KeyNSize, 0, len(namespaceCounts))
	for key, count := range namespaceCounts {
		spaces = append(spaces, KeyNSize{Key: key, Size: count})
	}

	sort.Slice(spaces, func(i, j int) bool {
		if spaces[i].Size == spaces[j].Size {
			return spaces[i].Key < spaces[j].Key
		}
		return spaces[i].Size > spaces[j].Size
	})

	return spaces
}

// TTL Analysis
func GetTTLMetrics() store.TTLMetrics {
	store.TTLMetricsContainer.Mu.RLock()
	defer store.TTLMetricsContainer.Mu.RUnlock()

	nextExpiring := append([]store.TTLInfo(nil), store.TTLMetricsContainer.NextExpiringKeys...)
	sort.Slice(nextExpiring, func(i, j int) bool {
		return nextExpiring[i].ExpiresIn < nextExpiring[j].ExpiresIn
	})

	return store.TTLMetrics{
		ActiveTTLKeys:    store.TTLMetricsContainer.ActiveTTLKeys,
		PermanentKeys:    store.TTLMetricsContainer.PermanentKeys,
		TotalExpiredKeys: store.TTLMetricsContainer.TotalExpiredKeys,
		ExpiresToday:     store.TTLMetricsContainer.ExpiresToday,
		NextExpiringKeys: nextExpiring,
	}
}


// Remove
func GetExpiredKeys(stre *store.MemoryAlloc) []string {
	stre.Mu.RLock()
	defer stre.Mu.RUnlock()

	now := time.Now()
	expired := make([]string, 0)
	for key, item := range stre.Data {
		if !item.Meta.TTL.IsZero() && item.Meta.TTL.Before(now) {
			expired = append(expired, key)
		}
	}

	sort.Strings(expired)
	return expired
}


func GetUpcomingExpirations(stre *store.MemoryAlloc, topN int) []store.TTLInfo {
	items := liveItems(stre)
	upcoming := make([]store.TTLInfo, 0)

	for _, item := range items {
		if item.Meta.TTL.IsZero() {
			continue
		}
		upcoming = append(upcoming, store.TTLInfo{
			Key:       item.Key,
			ExpiresIn: time.Until(item.Meta.TTL),
		})
	}

	sort.Slice(upcoming, func(i, j int) bool {
		return upcoming[i].ExpiresIn < upcoming[j].ExpiresIn
	})

	if topN <= 0 || topN >= len(upcoming) {
		return upcoming
	}

	return upcoming[:topN]
}

func GetKeysWithoutTTL(stre *store.MemoryAlloc) []string {
	items := liveItems(stre)
	keys := make([]string, 0)

	for _, item := range items {
		if item.Meta.TTL.IsZero() {
			keys = append(keys, item.Key)
		}
	}

	sort.Strings(keys)
	return keys
}

// Persistence
func GetSnapshotCount() int {
	return len(snapshotFiles())
}

func GetLatestSnapshot() SnapshotInfo {
	snaps := snapshotFiles()
	if len(snaps) == 0 {
		return SnapshotInfo{}
	}

	return snaps[0]
}

func GetSnapshotSize() int64 {
	snaps := snapshotFiles()
	var totalSize int64
	for _, snap := range snaps {
		totalSize += snap.Size
	}

	return totalSize
}

func GetPersistenceHealth() PersistenceHealth {
	dbFolder := utils.DbFolder()
	snapshotFolder := utils.SnapshotFolder()
	_, dbErr := os.Stat(dbFolder)
	_, snapshotErr := os.Stat(snapshotFolder)

	return PersistenceHealth{
		DbFolderExists:       dbErr == nil,
		SnapshotFolderExists: snapshotErr == nil,
		SnapshotCount:        GetSnapshotCount(),
		TotalSnapshotSize:    GetSnapshotSize(),
		LatestSnapshot:       GetLatestSnapshot(),
	}
}

// Metadata
func GetOldestKey(stre *store.MemoryAlloc) KeyMetadata {
	items := liveItems(stre)
	if len(items) == 0 {
		return KeyMetadata{}
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Meta.CreatedAt.Before(items[j].Meta.CreatedAt)
	})

	return toKeyMetadata(items[0])
}

func GetNewestKey(stre *store.MemoryAlloc) KeyMetadata {
	items := liveItems(stre)
	if len(items) == 0 {
		return KeyMetadata{}
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Meta.CreatedAt.After(items[j].Meta.CreatedAt)
	})

	return toKeyMetadata(items[0])
}

func GetMostAccessedKey(stre *store.MemoryAlloc) KeyMetadata {
	items := liveItems(stre)
	if len(items) == 0 {
		return KeyMetadata{}
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Meta.KeyAccessCount > items[j].Meta.KeyAccessCount
	})

	return toKeyMetadata(items[0])
}

func GetLeastAccessedKey(stre *store.MemoryAlloc) KeyMetadata {
	items := liveItems(stre)
	if len(items) == 0 {
		return KeyMetadata{}
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Meta.KeyAccessCount < items[j].Meta.KeyAccessCount
	})

	return toKeyMetadata(items[0])
}

func GetRecentlyModifiedKeys(stre *store.MemoryAlloc, topN int) []KeyMetadata {
	items := liveItems(stre)
	if len(items) == 0 || topN <= 0 {
		return []KeyMetadata{}
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Meta.UpdatedAt.After(items[j].Meta.UpdatedAt)
	})

	if topN > len(items) {
		topN = len(items)
	}

	metadata := make([]KeyMetadata, 0, topN)
	for idx := 0; idx < topN; idx++ {
		metadata = append(metadata, toKeyMetadata(items[idx]))
	}

	return metadata
}
