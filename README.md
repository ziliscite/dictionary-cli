# Dictionary CLI

A terminal-based Japanese dictionary application that uses the Jisho.org API to look up Japanese words and their definitions.

## Features

- Interactive terminal UI using the Bubble Tea framework
- Search for Japanese words and phrases
- View detailed information about words including:
  - Japanese writing (kanji and kana)
  - Readings
  - English definitions
  - Parts of speech
  - JLPT level information
- Keyboard navigation
- Loading indicators for search operations

## Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/ziliscite/dictionary-cli.git
cd dictionary-cli

# Build the application
go build -o dict-cli ./cmd

# Run the application
./dict-cli
```

### Using Docker

```bash
# Clone the repository
git clone https://github.com/ziliscite/dictionary-cli.git
cd dictionary-cli

# Build the Docker image
docker build -t dict-cli .

# Run the application in a Docker container
docker run -it dict-cli
```

## Usage

1. Start the application:
   ```
   dict-cli
   ```

2. Type a Japanese word or English word to search for
3. Press Enter to search
4. Navigate the results using arrow keys
5. Press Enter to view detailed information about a selected word
6. Press Backspace to return to the results list
7. Press Ctrl+S to start a new search
8. Press Esc or Ctrl+C to quit the application

## Keyboard Shortcuts

- `Enter` - Search or select an item
- `Backspace` - Return to previous view
- `Ctrl+S` - Return to search
- `Esc` or `Ctrl+C` - Quit the application
- Arrow keys - Navigate through search results
