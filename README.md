# converttomd-jira

A Go command-line tool to convert JIRA XML exports to Markdown format.

## Installation

Build the binary:
```bash
go build -o converttomd-jira
```

Install system-wide:
```bash
make install
```

### Optional: Create a Shorter Alias

For convenience, you can alias `jira-md` to `converttomd-jira`:

**For zsh:**
```bash
echo 'alias jira-md="converttomd-jira"' >> ~/.zshrc
source ~/.zshrc
```

**For bash:**
```bash
echo 'alias jira-md="converttomd-jira"' >> ~/.bashrc
source ~/.bashrc
```

Then use it as:
```bash
jira-md AI-538.xml
```

## Usage

```bash
converttomd-jira [OPTIONS] FILE [FILE...]
```

### Options

- `-o, --output <file>` - Output file path (defaults to `*.details.md` if details is on or `*.md` if details is off)
- `-d, --details <value>` - Include custom fields details (on|off|enabled|disabled|1|0) - defaults to enabled
- `-v, --verbose` - Verbose output
- `-f, --force` - Force overwrite existing files
- `--version` - Show version

### Examples

Convert a single file with default settings (includes custom fields):
```bash
converttomd-jira AI-538.xml
# Creates: AI-538.details.md
```

Convert without custom fields:
```bash
converttomd-jira -d off AI-538.xml
# Creates: AI-538.md
```

Specify custom output file:
```bash
converttomd-jira -o output.md AI-538.xml
```

Convert multiple files:
```bash
converttomd-jira *.xml
```

Force overwrite with verbose output:
```bash
converttomd-jira -f -v AI-538.xml
```

## Features

- Converts JIRA XML RSS format to clean Markdown
- Supports custom fields (can be toggled)
- HTML entity decoding
- Converts HTML tags to Markdown equivalents
- Handles comments, dates, labels, and attachments
- Multiple file processing
- Configurable output paths

## Output Format

The generated Markdown includes:

- Issue title and link
- Overview (type, priority, status, assignee, reporter, labels)
- Dates (created, updated, and custom date fields if details enabled)
- Full description with formatted HTML converted to Markdown
- Comments
- Custom fields (when details mode is enabled)

## License

MIT
