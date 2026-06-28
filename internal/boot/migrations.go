package boot

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/os/gres"
)

func RunMigrations() error {
	db := getDB()
	if db == nil {
		return fmt.Errorf("database not initialized")
	}
	ctx := context.Background()

	if _, err := db.Exec(ctx, `CREATE TABLE IF NOT EXISTS ydsterm_schema_migrations (
		version INTEGER PRIMARY KEY,
		applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`); err != nil {
		return fmt.Errorf("create migration table: %w", err)
	}

	applied, err := getAppliedVersions(ctx, db)
	if err != nil {
		return err
	}

	files, err := readMigrationFiles()
	if err != nil {
		return err
	}
	sort.Strings(files)

	appliedCount := 0
	for _, f := range files {
		if !strings.HasSuffix(f, ".up.sql") {
			continue
		}
		ver := parseVersion(filepath.Base(f))
		if ver < 0 {
			continue
		}
		if applied[ver] {
			continue
		}
		sql, err := readFileContent(f)
		if err != nil {
			return fmt.Errorf("read %s: %w", f, err)
		}
		if _, err := db.Exec(ctx, sql); err != nil {
			return fmt.Errorf("apply migration %s: %w", f, err)
		}
		if _, err := db.Exec(ctx, "INSERT INTO ydsterm_schema_migrations (version) VALUES (?)", ver); err != nil {
			return fmt.Errorf("record migration %s: %w", f, err)
		}
		appliedCount++
		log.Printf("migration applied: %s", filepath.Base(f))
	}

	if appliedCount == 0 {
		log.Println("no pending migrations")
	} else {
		log.Printf("applied %d migration(s)", appliedCount)
	}
	return nil
}

func DefaultDBPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "./data/ydsterm.db"
	}
	switch runtime.GOOS {
	case "windows":
		appData := os.Getenv("APPDATA")
		if appData != "" {
			return filepath.Join(appData, "ydsterm", "data", "ydsterm.db")
		}
		return filepath.Join(home, "AppData", "Roaming", "ydsterm", "data", "ydsterm.db")
	case "darwin":
		return filepath.Join(home, "Library", "Application Support", "ydsterm", "data", "ydsterm.db")
	default:
		xdg := os.Getenv("XDG_DATA_HOME")
		if xdg != "" {
			return filepath.Join(xdg, "ydsterm", "data", "ydsterm.db")
		}
		return filepath.Join(home, ".local", "share", "ydsterm", "data", "ydsterm.db")
	}
}

func EnsureDataDir(dbPath string) error {
	return os.MkdirAll(filepath.Dir(dbPath), 0755)
}

func getAppliedVersions(ctx context.Context, db gdb.DB) (map[int]bool, error) {
	rows, err := db.Query(ctx, "SELECT version FROM ydsterm_schema_migrations ORDER BY version")
	if err != nil {
		return nil, fmt.Errorf("query migrations: %w", err)
	}
	applied := map[int]bool{}
	for _, r := range rows {
		applied[r["version"].Int()] = true
	}
	return applied, nil
}

func readMigrationFiles() ([]string, error) {
	if !gres.Contains("migrations") {
		return nil, fmt.Errorf("migrations not found in gres, run 'make pack' first")
	}
	files := gres.ScanDirFile("migrations", "*.sql", false)
	var names []string
	for _, f := range files {
		names = append(names, f.Name())
	}
	return names, nil
}

func readFileContent(path string) (string, error) {
	b := gres.GetContent(path)
	if b == nil {
		return "", fmt.Errorf("gres: %s not found", path)
	}
	return string(b), nil
}

func getDB() gdb.DB {
	db, err := gdb.Instance("default")
	if err != nil {
		return nil
	}
	return db
}

func parseVersion(name string) int {
	parts := strings.SplitN(name, "_", 2)
	var ver int
	if _, err := fmt.Sscanf(parts[0], "%d", &ver); err != nil {
		return -1
	}
	return ver
}
