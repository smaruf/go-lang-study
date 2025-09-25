package main

import (
	"bufio"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

// Record represents a data record for processing
type Record struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Category    string    `json:"category"`
	Value       float64   `json:"value"`
	Date        time.Time `json:"date"`
	Active      bool      `json:"active"`
	Tags        []string  `json:"tags"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// ProcessingResult represents the result of data processing
type ProcessingResult struct {
	TotalRecords     int                            `json:"total_records"`
	ProcessedRecords int                            `json:"processed_records"`
	ErrorCount       int                            `json:"error_count"`
	Categories       map[string]int                 `json:"categories"`
	Statistics       map[string]float64             `json:"statistics"`
	TimeTaken        time.Duration                  `json:"time_taken"`
	Errors           []string                       `json:"errors,omitempty"`
}

// DataProcessor handles data processing operations
type DataProcessor struct {
	inputDir    string
	outputDir   string
	concurrency int
	logger      *log.Logger
}

// NewDataProcessor creates a new data processor
func NewDataProcessor(inputDir, outputDir string, concurrency int) *DataProcessor {
	return &DataProcessor{
		inputDir:    inputDir,
		outputDir:   outputDir,
		concurrency: concurrency,
		logger:      log.New(os.Stdout, "[DATA-PROCESSOR] ", log.LstdFlags|log.Lshortfile),
	}
}

// ProcessFiles processes all files in the input directory
func (dp *DataProcessor) ProcessFiles(ctx context.Context) (*ProcessingResult, error) {
	start := time.Now()
	
	// Ensure output directory exists
	if err := os.MkdirAll(dp.outputDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create output directory: %v", err)
	}

	// Find all input files
	files, err := dp.findInputFiles()
	if err != nil {
		return nil, fmt.Errorf("failed to find input files: %v", err)
	}

	dp.logger.Printf("Found %d files to process", len(files))

	// Process files concurrently
	result := &ProcessingResult{
		Categories: make(map[string]int),
		Statistics: make(map[string]float64),
		Errors:     []string{},
	}

	// Channel for file processing jobs
	fileChan := make(chan string, len(files))
	resultChan := make(chan *ProcessingResult, len(files))

	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < dp.concurrency; i++ {
		wg.Add(1)
		go dp.fileWorker(ctx, &wg, fileChan, resultChan)
	}

	// Send files to workers
	for _, file := range files {
		fileChan <- file
	}
	close(fileChan)

	// Collect results
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Aggregate results
	for fileResult := range resultChan {
		result.TotalRecords += fileResult.TotalRecords
		result.ProcessedRecords += fileResult.ProcessedRecords
		result.ErrorCount += fileResult.ErrorCount
		
		// Merge categories
		for category, count := range fileResult.Categories {
			result.Categories[category] += count
		}
		
		// Merge errors
		result.Errors = append(result.Errors, fileResult.Errors...)
	}

	// Calculate statistics
	dp.calculateStatistics(result)
	result.TimeTaken = time.Since(start)

	dp.logger.Printf("Processing completed in %v", result.TimeTaken)
	return result, nil
}

// fileWorker processes individual files
func (dp *DataProcessor) fileWorker(ctx context.Context, wg *sync.WaitGroup, fileChan <-chan string, resultChan chan<- *ProcessingResult) {
	defer wg.Done()

	for file := range fileChan {
		select {
		case <-ctx.Done():
			return
		default:
			result, err := dp.processFile(file)
			if err != nil {
				dp.logger.Printf("Error processing file %s: %v", file, err)
				result = &ProcessingResult{
					ErrorCount: 1,
					Errors:     []string{fmt.Sprintf("File %s: %v", file, err)},
					Categories: make(map[string]int),
					Statistics: make(map[string]float64),
				}
			}
			resultChan <- result
		}
	}
}

// processFile processes a single file
func (dp *DataProcessor) processFile(filePath string) (*ProcessingResult, error) {
	dp.logger.Printf("Processing file: %s", filePath)

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(filePath))
	var records []Record

	switch ext {
	case ".json":
		records, err = dp.processJSONFile(file)
	case ".csv":
		records, err = dp.processCSVFile(file)
	case ".txt":
		records, err = dp.processTextFile(file)
	default:
		return nil, fmt.Errorf("unsupported file type: %s", ext)
	}

	if err != nil {
		return nil, err
	}

	// Process and analyze records
	result := dp.analyzeRecords(records)
	
	// Save processed data
	outputFile := filepath.Join(dp.outputDir, 
		strings.TrimSuffix(filepath.Base(filePath), ext)+"_processed.json")
	
	if err := dp.saveProcessedData(outputFile, records); err != nil {
		dp.logger.Printf("Warning: failed to save processed data: %v", err)
	}

	return result, nil
}

// processJSONFile processes JSON files
func (dp *DataProcessor) processJSONFile(file io.Reader) ([]Record, error) {
	var records []Record
	decoder := json.NewDecoder(file)
	
	// Handle both single objects and arrays
	token, err := decoder.Token()
	if err != nil {
		return nil, err
	}

	if delim, ok := token.(json.Delim); ok && delim == '[' {
		// Array of records
		for decoder.More() {
			var record Record
			if err := decoder.Decode(&record); err != nil {
				dp.logger.Printf("Warning: skipping invalid record: %v", err)
				continue
			}
			records = append(records, record)
		}
	} else {
		// Single record - need to put token back
		var data json.RawMessage
		data = append(data, fmt.Sprintf("%v", token)...)
		remaining, _ := io.ReadAll(decoder.Buffered())
		data = append(data, remaining...)
		
		var record Record
		if err := json.Unmarshal(data, &record); err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	return records, nil
}

// processCSVFile processes CSV files
func (dp *DataProcessor) processCSVFile(file io.Reader) ([]Record, error) {
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1 // Allow variable number of fields

	// Read header
	header, err := reader.Read()
	if err != nil {
		return nil, err
	}

	var records []Record
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			dp.logger.Printf("Warning: skipping invalid CSV row: %v", err)
			continue
		}

		record, err := dp.parseCSVRow(header, row)
		if err != nil {
			dp.logger.Printf("Warning: skipping row: %v", err)
			continue
		}
		records = append(records, record)
	}

	return records, nil
}

// processTextFile processes simple text files
func (dp *DataProcessor) processTextFile(file io.Reader) ([]Record, error) {
	scanner := bufio.NewScanner(file)
	var records []Record
	id := 1

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		// Simple text processing - treat each line as a record
		record := Record{
			ID:       id,
			Name:     line,
			Category: "text",
			Value:    float64(len(line)),
			Date:     time.Now(),
			Active:   true,
			Tags:     strings.Fields(line),
			Metadata: map[string]interface{}{
				"line_length": len(line),
				"word_count":  len(strings.Fields(line)),
			},
		}
		records = append(records, record)
		id++
	}

	return records, scanner.Err()
}

// parseCSVRow converts CSV row to Record
func (dp *DataProcessor) parseCSVRow(header, row []string) (Record, error) {
	if len(row) < len(header) {
		return Record{}, fmt.Errorf("row has fewer fields than header")
	}

	record := Record{
		Metadata: make(map[string]interface{}),
		Tags:     []string{},
	}

	for i, field := range header {
		value := ""
		if i < len(row) {
			value = row[i]
		}

		switch strings.ToLower(field) {
		case "id":
			if id, err := strconv.Atoi(value); err == nil {
				record.ID = id
			}
		case "name":
			record.Name = value
		case "category":
			record.Category = value
		case "value":
			if val, err := strconv.ParseFloat(value, 64); err == nil {
				record.Value = val
			}
		case "date":
			if date, err := time.Parse("2006-01-02", value); err == nil {
				record.Date = date
			} else if date, err := time.Parse(time.RFC3339, value); err == nil {
				record.Date = date
			}
		case "active":
			record.Active = strings.ToLower(value) == "true"
		case "tags":
			if value != "" {
				record.Tags = strings.Split(value, ",")
				for i, tag := range record.Tags {
					record.Tags[i] = strings.TrimSpace(tag)
				}
			}
		default:
			record.Metadata[field] = value
		}
	}

	return record, nil
}

// analyzeRecords analyzes a slice of records
func (dp *DataProcessor) analyzeRecords(records []Record) *ProcessingResult {
	result := &ProcessingResult{
		TotalRecords: len(records),
		ProcessedRecords: len(records),
		Categories: make(map[string]int),
		Statistics: make(map[string]float64),
		Errors: []string{},
	}

	if len(records) == 0 {
		return result
	}

	// Count categories
	var totalValue float64
	var activeCount int
	values := make([]float64, 0, len(records))

	for _, record := range records {
		result.Categories[record.Category]++
		totalValue += record.Value
		values = append(values, record.Value)
		
		if record.Active {
			activeCount++
		}
	}

	// Calculate statistics
	sort.Float64s(values)
	
	result.Statistics["average_value"] = totalValue / float64(len(records))
	result.Statistics["min_value"] = values[0]
	result.Statistics["max_value"] = values[len(values)-1]
	result.Statistics["median_value"] = dp.median(values)
	result.Statistics["active_percentage"] = float64(activeCount) / float64(len(records)) * 100

	return result
}

// calculateStatistics calculates overall statistics
func (dp *DataProcessor) calculateStatistics(result *ProcessingResult) {
	if result.TotalRecords > 0 {
		result.Statistics["success_rate"] = float64(result.ProcessedRecords) / float64(result.TotalRecords) * 100
		result.Statistics["error_rate"] = float64(result.ErrorCount) / float64(result.TotalRecords) * 100
	}

	// Find most common category
	maxCount := 0
	mostCommon := ""
	for category, count := range result.Categories {
		if count > maxCount {
			maxCount = count
			mostCommon = category
		}
	}
	if mostCommon != "" {
		result.Statistics["most_common_category_count"] = float64(maxCount)
	}
}

// median calculates the median of a sorted slice
func (dp *DataProcessor) median(values []float64) float64 {
	n := len(values)
	if n == 0 {
		return 0
	}
	if n%2 == 0 {
		return (values[n/2-1] + values[n/2]) / 2
	}
	return values[n/2]
}

// saveProcessedData saves processed data to JSON file
func (dp *DataProcessor) saveProcessedData(filename string, records []Record) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(records)
}

// findInputFiles finds all processable files in input directory
func (dp *DataProcessor) findInputFiles() ([]string, error) {
	var files []string
	
	err := filepath.Walk(dp.inputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if !info.IsDir() {
			ext := strings.ToLower(filepath.Ext(path))
			if ext == ".json" || ext == ".csv" || ext == ".txt" {
				files = append(files, path)
			}
		}
		
		return nil
	})

	return files, err
}

// GenerateSampleData creates sample data files for testing
func (dp *DataProcessor) GenerateSampleData() error {
	if err := os.MkdirAll(dp.inputDir, 0755); err != nil {
		return err
	}

	// Generate JSON sample
	jsonRecords := []Record{
		{ID: 1, Name: "Product A", Category: "electronics", Value: 99.99, Date: time.Now(), Active: true, Tags: []string{"new", "featured"}},
		{ID: 2, Name: "Product B", Category: "books", Value: 19.99, Date: time.Now().AddDate(0, -1, 0), Active: true, Tags: []string{"bestseller"}},
		{ID: 3, Name: "Product C", Category: "electronics", Value: 199.99, Date: time.Now().AddDate(0, 0, -7), Active: false, Tags: []string{"discontinued"}},
	}

	jsonFile, err := os.Create(filepath.Join(dp.inputDir, "sample.json"))
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	encoder := json.NewEncoder(jsonFile)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(jsonRecords); err != nil {
		return err
	}

	// Generate CSV sample
	csvFile, err := os.Create(filepath.Join(dp.inputDir, "sample.csv"))
	if err != nil {
		return err
	}
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	// Write CSV header and data
	writer.Write([]string{"id", "name", "category", "value", "date", "active", "tags"})
	writer.Write([]string{"4", "Service A", "service", "49.99", "2024-01-15", "true", "premium,support"})
	writer.Write([]string{"5", "Service B", "service", "29.99", "2024-01-20", "true", "basic"})
	writer.Write([]string{"6", "Service C", "service", "99.99", "2024-01-10", "false", "legacy,deprecated"})

	// Generate text sample
	textFile, err := os.Create(filepath.Join(dp.inputDir, "sample.txt"))
	if err != nil {
		return err
	}
	defer textFile.Close()

	textData := []string{
		"Go is an open source programming language",
		"Developed by Google for modern software development",
		"Known for its simplicity and performance",
		"Great for concurrent programming",
		"Used in cloud computing and microservices",
	}

	for _, line := range textData {
		fmt.Fprintln(textFile, line)
	}

	dp.logger.Println("Sample data generated successfully")
	return nil
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found, using defaults")
	}

	// Configuration
	inputDir := getEnv("INPUT_DIR", "./data/input")
	outputDir := getEnv("OUTPUT_DIR", "./data/output")
	concurrency, _ := strconv.Atoi(getEnv("CONCURRENCY", "3"))
	generateSample := getEnv("GENERATE_SAMPLE", "true") == "true"

	// Create processor
	processor := NewDataProcessor(inputDir, outputDir, concurrency)

	// Generate sample data if requested
	if generateSample {
		if err := processor.GenerateSampleData(); err != nil {
			log.Fatalf("Failed to generate sample data: %v", err)
		}
	}

	// Process files
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	result, err := processor.ProcessFiles(ctx)
	if err != nil {
		log.Fatalf("Processing failed: %v", err)
	}

	// Print results
	fmt.Printf("\n=== DATA PROCESSING RESULTS ===\n")
	fmt.Printf("Total Records: %d\n", result.TotalRecords)
	fmt.Printf("Processed Records: %d\n", result.ProcessedRecords)
	fmt.Printf("Error Count: %d\n", result.ErrorCount)
	fmt.Printf("Processing Time: %v\n", result.TimeTaken)
	fmt.Printf("\nCategories:\n")
	for category, count := range result.Categories {
		fmt.Printf("  %s: %d\n", category, count)
	}
	fmt.Printf("\nStatistics:\n")
	for stat, value := range result.Statistics {
		fmt.Printf("  %s: %.2f\n", stat, value)
	}

	if len(result.Errors) > 0 {
		fmt.Printf("\nErrors:\n")
		for _, err := range result.Errors {
			fmt.Printf("  %s\n", err)
		}
	}

	// Save results to file
	resultsFile := filepath.Join(outputDir, "processing_results.json")
	if file, err := os.Create(resultsFile); err == nil {
		defer file.Close()
		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")
		encoder.Encode(result)
		fmt.Printf("\nResults saved to: %s\n", resultsFile)
	}

	fmt.Printf("\nProcessed files are available in: %s\n", outputDir)
}

// getEnv gets environment variable with fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}