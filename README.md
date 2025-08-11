# Japanese Learning CLI

Not just a dictionary anymore, lol

A terminal-based Japanese dictionary, translation, and Japanese analysis application that uses the Jisho.org API for word lookups, the DeepL API for text translation, and DeepSeek AI for Japanese sentence analysis.

## Features

- Interactive terminal UI using the Bubble Tea framework
- Dictionary functionality:
  - Search for Japanese words and phrases
  - View detailed information about words including:
    - Japanese writing (kanji and kana)
    - Readings
    - English definitions
    - Parts of speech
    - JLPT level information
- Translation functionality:
  - Translate text between multiple languages (Japanese, English, Indonesian)
  - Powered by DeepL API
- Japanese sentence explainer:
  - Analyze Japanese sentences for in-depth understanding
  - Get kana reading, romaji, and both literal and natural translations
  - Word-by-word breakdown with parts of speech and meanings
  - Grammar point explanations with similar examples
  - Nuance and register information
  - Common errors and alternative expressions
  - Practice exercises with answers
  - Powered by DeepSeek AI
- Keyboard navigation
- Loading indicators for search and translation operations

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

2. From the main menu, select either "Search" (dictionary lookup), "Translate" (text translation), or "Explainer" (Japanese sentence analysis) using arrow keys and press Enter

### Dictionary Mode
1. Type a Japanese word or English word to search for
2. Press Enter to search
3. Navigate the results using arrow keys
4. Press Enter to view detailed information about a selected word
5. Press Ctrl+Q to return to the results' list
6. Press Ctrl+S to start a new search
7. Press Ctrl+Q to return to the main menu
8. Press Esc or Ctrl+C to quit the application

### Translation Mode
1. Type the text you want to translate
2. Use Shift+Left and Shift+Right to cycle between target languages (Japanese, English, Indonesian)
3. Press Ctrl+T to translate the text
4. Press Ctrl+Q to return to the main menu
5. Press Esc or Ctrl+C to quit the application

### Explainer Mode
1. Type a Japanese sentence you want to analyze
2. Press Enter to get the explanation
3. Use arrow keys or j/k to scroll through the detailed explanation
4. Press Ctrl+Q to return to the input screen
5. Press Ctrl+Q again to return to the main menu
6. Press Esc or Ctrl+C to quit the application

## Keyboard Shortcuts

### General
- `Esc` or `Ctrl+C` - Quit the application
- `Ctrl+Q` - Return to previous view

### Main Menu
- Arrow keys - Navigate between options
- `Enter` - Select an option

### Dictionary Mode
- `Enter` - Search or select an item
- `Ctrl+Q` - Return to previous view
- `Ctrl+S` - Return to search
- Arrow keys - Navigate through search results

### Translation Mode
- `Shift+Left` / `Shift+Right` - Cycle between target languages
- `Ctrl+T` - Translate the entered text

### Explainer Mode
- `Enter` - Submit Japanese sentence for analysis
- `↑/k` / `↓/j` - Scroll through explanation
- `Ctrl+Q` - Return to input screen or main menu

## Warnings

Your terminal font may not support Japanese characters or it's too small.
1) Set your terminal font to a CJK-capable font (e.g. Noto Sans Mono CJK, Source Han Code JP, Sarasa Gothic).
2) Increase font size in your terminal preferences.
3) If you're on iTerm2, create a profile with a large font and name it "LargeFont".
   Then run: echo -e "\033]50;SetProfile=LargeFont\a"
