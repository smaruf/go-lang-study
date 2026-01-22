package main

import (
	"fmt"
	"sync"
)

// Singleton is a thread-safe singleton instance
type Singleton struct {
	data string
}

var (
	instance *Singleton
	once     sync.Once
)

// GetInstance returns the singleton instance
// Using sync.Once ensures thread-safe lazy initialization
func GetInstance() *Singleton {
	once.Do(func() {
		fmt.Println("Creating singleton instance...")
		instance = &Singleton{
			data: "I am a singleton!",
		}
	})
	return instance
}

// GetData returns the singleton's data
func (s *Singleton) GetData() string {
	return s.data
}

// SetData sets the singleton's data
func (s *Singleton) SetData(data string) {
	s.data = data
}

// DatabaseConnection demonstrates a practical singleton use case
type DatabaseConnection struct {
	connectionString string
	connected        bool
}

var (
	dbInstance *DatabaseConnection
	dbOnce     sync.Once
)

// GetDatabaseConnection returns the singleton database connection
func GetDatabaseConnection() *DatabaseConnection {
	dbOnce.Do(func() {
		fmt.Println("Establishing database connection...")
		dbInstance = &DatabaseConnection{
			connectionString: "localhost:5432",
			connected:        true,
		}
	})
	return dbInstance
}

// Query simulates a database query
func (db *DatabaseConnection) Query(sql string) string {
	if !db.connected {
		return "Not connected"
	}
	return fmt.Sprintf("Executing: %s on %s", sql, db.connectionString)
}

// Close closes the database connection
func (db *DatabaseConnection) Close() {
	db.connected = false
	fmt.Println("Database connection closed")
}

// ConfigManager demonstrates another singleton pattern use case
type ConfigManager struct {
	settings map[string]string
	mu       sync.RWMutex
}

var (
	configInstance *ConfigManager
	configOnce     sync.Once
)

// GetConfigManager returns the singleton config manager
func GetConfigManager() *ConfigManager {
	configOnce.Do(func() {
		fmt.Println("Initializing configuration manager...")
		configInstance = &ConfigManager{
			settings: make(map[string]string),
		}
		// Load default settings
		configInstance.settings["app_name"] = "MyApp"
		configInstance.settings["version"] = "1.0.0"
		configInstance.settings["debug"] = "false"
	})
	return configInstance
}

// Get retrieves a configuration value
func (cm *ConfigManager) Get(key string) (string, bool) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	val, exists := cm.settings[key]
	return val, exists
}

// Set sets a configuration value
func (cm *ConfigManager) Set(key, value string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.settings[key] = value
}

// GetAll returns all configuration settings
func (cm *ConfigManager) GetAll() map[string]string {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	
	// Return a copy to prevent external modification
	result := make(map[string]string)
	for k, v := range cm.settings {
		result[k] = v
	}
	return result
}

func main() {
	fmt.Println("=== Basic Singleton Pattern ===")
	
	// Get singleton instance multiple times
	s1 := GetInstance()
	s2 := GetInstance()
	s3 := GetInstance()

	fmt.Printf("s1 data: %s\n", s1.GetData())
	fmt.Printf("s1 == s2: %v\n", s1 == s2)
	fmt.Printf("s2 == s3: %v\n", s2 == s3)

	// Modify data through one reference
	s1.SetData("Modified data")
	fmt.Printf("s2 data after s1 modification: %s\n", s2.GetData())

	fmt.Println("\n=== Database Connection Singleton ===")
	
	// Get database connection
	db1 := GetDatabaseConnection()
	db2 := GetDatabaseConnection()

	fmt.Printf("db1 == db2: %v\n", db1 == db2)
	
	result := db1.Query("SELECT * FROM users")
	fmt.Println(result)

	result = db2.Query("SELECT * FROM orders")
	fmt.Println(result)

	fmt.Println("\n=== Configuration Manager Singleton ===")
	
	config1 := GetConfigManager()
	config2 := GetConfigManager()

	fmt.Printf("config1 == config2: %v\n", config1 == config2)

	// Read configuration
	if appName, exists := config1.Get("app_name"); exists {
		fmt.Printf("App Name: %s\n", appName)
	}

	// Modify configuration
	config1.Set("debug", "true")
	config2.Set("max_connections", "100")

	// Verify changes are shared
	if debug, exists := config2.Get("debug"); exists {
		fmt.Printf("Debug mode (from config2): %s\n", debug)
	}

	if maxConn, exists := config1.Get("max_connections"); exists {
		fmt.Printf("Max Connections (from config1): %s\n", maxConn)
	}

	fmt.Println("\nAll settings:")
	for key, value := range config1.GetAll() {
		fmt.Printf("  %s: %s\n", key, value)
	}

	fmt.Println("\n=== Thread Safety Test ===")
	
	// Test concurrent access
	var wg sync.WaitGroup
	numGoroutines := 10

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			instance := GetInstance()
			fmt.Printf("Goroutine %d got instance: %p\n", id, instance)
		}(i)
	}

	wg.Wait()
	fmt.Println("\nAll goroutines completed - all instances should have the same address")
}
