# Data Processing Pipeline

A comprehensive data processing pipeline demonstrating file I/O, concurrent processing, data transformation, and statistical analysis in Go.

## Features

- **Multi-format Support**: Process JSON, CSV, and text files
- **Concurrent Processing**: Configurable worker pool for parallel file processing
- **Data Analysis**: Statistical analysis and categorization
- **Error Handling**: Robust error handling with detailed reporting
- **Flexible Output**: JSON output with processed data and analytics
- **Sample Data Generation**: Automatic generation of test data
- **Context-based Cancellation**: Timeout and cancellation support
- **Environment Configuration**: Configurable via environment variables

## Quick Start

1. **Install dependencies:**
   ```bash
   go mod download
   ```

2. **Setup environment:**
   ```bash
   cp .env.example .env
   # Edit .env if needed
   ```

3. **Run with sample data generation:**
   ```bash
   go run main.go
   ```

4. **Check results:**
   ```bash
   ls -la data/output/
   cat data/output/processing_results.json
   ```

## Configuration

Environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `INPUT_DIR` | `./data/input` | Directory containing input files |
| `OUTPUT_DIR` | `./data/output` | Directory for processed output |
| `CONCURRENCY` | `3` | Number of concurrent workers |
| `GENERATE_SAMPLE` | `true` | Generate sample data files |

## Supported File Formats

### JSON Files (.json)
Supports both single objects and arrays of objects:

```json
[
  {
    "id": 1,
    "name": "Product A",
    "category": "electronics",
    "value": 99.99,
    "date": "2024-01-15T10:30:00Z",
    "active": true,
    "tags": ["new", "featured"],
    "metadata": {"key": "value"}
  }
]
```

### CSV Files (.csv)
Standard CSV with header row:

```csv
id,name,category,value,date,active,tags
1,Product A,electronics,99.99,2024-01-15,true,"new,featured"
```

### Text Files (.txt)
Simple text files where each line becomes a record:

```text
Go is an open source programming language
Developed by Google for modern software development
```

## Data Processing Pipeline

### 1. File Discovery
- Recursively scan input directory
- Filter by supported file extensions
- Queue files for processing

### 2. Concurrent Processing
- Configurable worker pool
- Context-based cancellation
- Error isolation per file

### 3. Data Transformation
- Format-specific parsing
- Data validation and cleaning
- Structure normalization

### 4. Analysis & Statistics
- Category counting
- Value statistics (min, max, average, median)
- Active/inactive ratios
- Error rates and success metrics

### 5. Output Generation
- Processed data in JSON format
- Comprehensive processing report
- Error logging and reporting

## Usage Examples

### Basic Processing
```bash
# Process files with default settings
go run main.go
```

### Custom Configuration
```bash
# Set custom directories and worker count
export INPUT_DIR="/path/to/input"
export OUTPUT_DIR="/path/to/output"
export CONCURRENCY=5
go run main.go
```

### Processing Existing Data
```bash
# Disable sample generation for existing data
export GENERATE_SAMPLE=false
go run main.go
```

## Output Structure

### Processed Data Files
Each input file generates a corresponding `*_processed.json` file:
```
data/output/
├── sample_processed.json
├── sample_processed.json
├── sample_processed.json
└── processing_results.json
```

### Processing Results
The `processing_results.json` contains comprehensive analytics:

```json
{
  "total_records": 11,
  "processed_records": 11,
  "error_count": 0,
  "categories": {
    "electronics": 2,
    "books": 1,
    "service": 3,
    "text": 5
  },
  "statistics": {
    "average_value": 67.32,
    "min_value": 5.0,
    "max_value": 199.99,
    "median_value": 49.99,
    "active_percentage": 72.73,
    "success_rate": 100.00,
    "error_rate": 0.00
  },
  "time_taken": "15.2ms"
}
```

## Architecture

### Core Components

1. **DataProcessor**: Main processing coordinator
2. **Record**: Unified data structure for all formats
3. **ProcessingResult**: Analytics and statistics container
4. **File Workers**: Concurrent processing workers

### Processing Flow

```
Input Files → File Discovery → Worker Pool → Format Parsing → 
Data Analysis → Statistics → Output Generation
```

### Concurrency Model

- **Producer-Consumer**: File queue with multiple workers
- **Worker Pool**: Configurable number of processing goroutines
- **Result Aggregation**: Channel-based result collection
- **Context Cancellation**: Timeout and cancellation support

## Key Concepts Demonstrated

1. **File I/O**: Reading various file formats (JSON, CSV, text)
2. **Concurrent Processing**: Worker pool pattern with goroutines
3. **Data Structures**: Complex data modeling and transformation
4. **Error Handling**: Graceful error handling and reporting
5. **Statistics**: Mathematical analysis and data summarization
6. **Context Usage**: Timeout and cancellation patterns
7. **Environment Config**: Configuration management
8. **JSON Processing**: Marshaling and unmarshaling complex data
9. **CSV Processing**: Reading and parsing CSV files
10. **Text Processing**: Line-by-line text file processing

## Performance Characteristics

- **Concurrent**: Parallel processing of multiple files
- **Memory Efficient**: Streaming processing where possible
- **Scalable**: Configurable worker pool size
- **Robust**: Continues processing despite individual file errors

## Error Handling

The pipeline implements comprehensive error handling:

- **File Level**: Individual file errors don't stop overall processing
- **Record Level**: Invalid records are logged but processing continues
- **System Level**: OS-level errors are properly handled
- **Timeout**: Context-based timeout prevents hanging

## Production Considerations

For production use, consider:

- **Database Integration**: Store results in persistent storage
- **Monitoring**: Add metrics and health checks
- **Validation**: More comprehensive data validation
- **Scheduling**: Integrate with cron or job schedulers
- **Notifications**: Alert on processing failures
- **Cleanup**: Automatic cleanup of old processed files
- **Compression**: Handle compressed input files
- **Streaming**: Process very large files in streaming mode

## Testing

```bash
# Run with sample data
go run main.go

# Check output
ls -la data/output/
jq . data/output/processing_results.json

# Test with custom data
mkdir -p custom_input
echo '{"id": 1, "name": "Test", "value": 100}' > custom_input/test.json
INPUT_DIR=custom_input OUTPUT_DIR=custom_output GENERATE_SAMPLE=false go run main.go
```

## Learning Objectives

This example teaches:

- **File Processing**: Handling multiple file formats
- **Concurrent Programming**: Worker pool patterns
- **Data Analysis**: Statistical computation
- **Error Management**: Robust error handling strategies
- **JSON/CSV Processing**: Data format handling
- **Context Usage**: Timeout and cancellation
- **Performance Optimization**: Concurrent processing benefits
- **Configuration Management**: Environment-based settings

## Related Examples

- [Concurrency](../concurrency/) - Worker pool patterns
- [CLI Tool](../cli-tool/) - Command-line applications
- [Web Server](../web-server/) - HTTP API development