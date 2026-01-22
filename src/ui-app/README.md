# Go UI Application Examples

This directory contains examples of building graphical user interfaces (GUIs) in Go using various UI frameworks.

## Frameworks Covered

### Fyne
- **Cross-platform**: Works on Windows, macOS, Linux, iOS, Android
- **Modern design**: Material Design-inspired widgets
- **Easy to use**: Simple, declarative API

### Examples

1. **Hello World** (`hello_fyne.go`)
   - Basic window with text and button
   - Demonstrates window creation and widget usage

2. **Form Application** (`form_app.go`)
   - Complete form with input fields
   - Form validation and submission
   - Layout management

3. **Todo List** (`todo_app.go`)
   - Full-featured todo list application
   - List management, add/remove items
   - Data persistence

4. **Calculator** (`calculator.go`)
   - Simple calculator application
   - Button grid layout
   - Basic arithmetic operations

## Prerequisites

### For Fyne
Install Fyne dependencies:

**Linux (Ubuntu/Debian):**
```bash
sudo apt-get install gcc libgl1-mesa-dev xorg-dev
```

**macOS:**
```bash
# Xcode command line tools required
xcode-select --install
```

**Windows:**
- Install GCC (via MinGW-w64 or TDM-GCC)

### Install Fyne
```bash
go get fyne.io/fyne/v2
```

## Running Examples

### Hello World
```bash
go run hello_fyne.go
```

### Form Application
```bash
go run form_app.go
```

### Todo List
```bash
go run todo_app.go
```

### Calculator
```bash
go run calculator.go
```

## Building Standalone Applications

### Build for current platform
```bash
# Build hello_fyne
go build -o hello_fyne hello_fyne.go

# Build with Fyne bundler for better packaging
go install fyne.io/fyne/v2/cmd/fyne@latest
fyne package -executable hello_fyne
```

### Cross-compile for different platforms
```bash
# For Windows from Linux/macOS
fyne package -os windows -icon icon.png

# For macOS from Linux/Windows
fyne package -os darwin -icon icon.png

# For Linux from Windows/macOS
fyne package -os linux -icon icon.png
```

## Alternative UI Frameworks

### Gio
- Pure Go, no CGO required
- Immediate mode GUI
- Very lightweight

```bash
go get gioui.org
```

### Walk (Windows only)
- Native Windows GUI
- Windows-specific features

```bash
go get github.com/lxn/walk
```

### Wails (Web-based)
- Build desktop apps with web technologies
- Go backend, HTML/CSS/JS frontend

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

## Key Concepts

### Widgets
- Buttons, Labels, Entries (text inputs)
- Lists, Tables, Trees
- Containers and layouts

### Layouts
- **VBox**: Vertical box layout
- **HBox**: Horizontal box layout
- **Grid**: Grid layout
- **Border**: Border layout with top, bottom, left, right, center
- **Form**: Form layout for labels and inputs

### Containers
- **Split**: Resizable split containers
- **Scroll**: Scrollable containers
- **Tab**: Tabbed containers

### Themes
- Light and Dark themes
- Custom theme support
- Icon theming

## Best Practices

1. **Keep UI logic separate from business logic**
2. **Use data binding for reactive UIs**
3. **Handle long-running operations in goroutines**
4. **Test business logic separately from UI**
5. **Use proper error handling and user feedback**

## Resources

- [Fyne Documentation](https://developer.fyne.io/)
- [Fyne Examples](https://github.com/fyne-io/examples)
- [Gio Documentation](https://gioui.org/)
- [Go GUI Projects](https://github.com/topics/gui?l=go)

## Troubleshooting

### Build issues on Linux
```bash
# Install required libraries
sudo apt-get install gcc libgl1-mesa-dev xorg-dev

# For some systems, also install:
sudo apt-get install libxcursor-dev libxrandr-dev libxinerama-dev libxi-dev libxxf86vm-dev
```

### Build issues on macOS
```bash
# Ensure Xcode command line tools are installed
xcode-select --install
```

### Build issues on Windows
- Ensure GCC is in your PATH
- Use MinGW-w64 or TDM-GCC
- Try running from Git Bash or PowerShell as administrator
