package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	pflag "github.com/spf13/pflag"
)

const version = "1.0.0"

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	Items []Item `xml:"item"`
}

type Item struct {
	Title       string       `xml:"title"`
	Link        string       `xml:"link"`
	Key         Key          `xml:"key"`
	Summary     string       `xml:"summary"`
	Type        TypeField    `xml:"type"`
	Priority    Priority     `xml:"priority"`
	Status      Status       `xml:"status"`
	Resolution  Resolution   `xml:"resolution"`
	Assignee    string       `xml:"assignee"`
	Reporter    string       `xml:"reporter"`
	Labels      Labels       `xml:"labels"`
	Description string       `xml:"description"`
	Created     string       `xml:"created"`
	Updated     string       `xml:"updated"`
	Due         string       `xml:"due"`
	Comments    Comments     `xml:"comments"`
	CustomFields CustomFields `xml:"customfields"`
}

type Key struct {
	ID    string `xml:"id,attr"`
	Value string `xml:",chardata"`
}

type TypeField struct {
	ID      string `xml:"id,attr"`
	IconURL string `xml:"iconUrl,attr"`
	Value   string `xml:",chardata"`
}

type Priority struct {
	ID      string `xml:"id,attr"`
	IconURL string `xml:"iconUrl,attr"`
	Value   string `xml:",chardata"`
}

type Status struct {
	ID      string `xml:"id,attr"`
	IconURL string `xml:"iconUrl,attr"`
	Value   string `xml:",chardata"`
}

type Resolution struct {
	ID    string `xml:"id,attr"`
	Value string `xml:",chardata"`
}

type Labels struct {
	Label []string `xml:"label"`
}

type Comments struct {
	Comment []Comment `xml:"comment"`
}

type Comment struct {
	ID      string `xml:"id,attr"`
	Author  string `xml:"author,attr"`
	Created string `xml:"created,attr"`
	Value   string `xml:",chardata"`
}

type CustomFields struct {
	CustomField []CustomField `xml:"customfield"`
}

type CustomField struct {
	ID              string                 `xml:"id,attr"`
	Key             string                 `xml:"key,attr"`
	CustomFieldName string                 `xml:"customfieldname"`
	CustomFieldValues CustomFieldValues    `xml:"customfieldvalues"`
}

type CustomFieldValues struct {
	CustomFieldValue []CustomFieldValue `xml:"customfieldvalue"`
}

type CustomFieldValue struct {
	Key   string `xml:"key,attr"`
	Value string `xml:",chardata"`
}

type Config struct {
	inputFiles []string
	output     string
	details    bool
	verbose    bool
	force      bool
	showVersion bool
}

func main() {
	config := parseFlags()

	if config.showVersion {
		fmt.Printf("converttomd-jira version %s\n", version)
		os.Exit(0)
	}

	if len(config.inputFiles) == 0 {
		fmt.Fprintln(os.Stderr, "Error: no input files specified")
		pflag.Usage()
		os.Exit(1)
	}

	for _, inputFile := range config.inputFiles {
		if err := processFile(inputFile, config); err != nil {
			fmt.Fprintf(os.Stderr, "Error processing %s: %v\n", inputFile, err)
			os.Exit(1)
		}
	}
}

func parseFlags() Config {
	config := Config{}

	var detailsStr string
	pflag.StringVarP(&config.output, "output", "o", "", "Output file path (defaults to *.details.md or *.md)")
	pflag.StringVarP(&detailsStr, "details", "d", "enabled", "Include custom fields details (on|off|enabled|disabled|1|0)")
	pflag.BoolVarP(&config.verbose, "verbose", "v", false, "Verbose output")
	pflag.BoolVarP(&config.force, "force", "f", false, "Force overwrite existing files")
	pflag.BoolVar(&config.showVersion, "version", false, "Show version")

	pflag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS] FILE [FILE...]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Convert JIRA XML exports to Markdown format.\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		pflag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s AI-538.xml\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s --output output.md AI-538.xml\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s --details off AI-538.xml\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s *.xml\n", os.Args[0])
	}

	pflag.Parse()

	config.inputFiles = pflag.Args()
	config.details = parseDetailsFlag(detailsStr)

	return config
}

func parseDetailsFlag(s string) bool {
	switch strings.ToLower(s) {
	case "on", "enabled", "1", "true", "yes":
		return true
	case "off", "disabled", "0", "false", "no":
		return false
	default:
		return true // default to enabled
	}
}

func processFile(inputFile string, config Config) error {
	if config.verbose {
		fmt.Printf("Processing %s...\n", inputFile)
	}

	// Read and parse XML
	data, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	var rss RSS
	if err := xml.Unmarshal(data, &rss); err != nil {
		return fmt.Errorf("failed to parse XML: %w", err)
	}

	if len(rss.Channel.Items) == 0 {
		return fmt.Errorf("no items found in XML")
	}

	item := rss.Channel.Items[0]

	// Determine output file
	outputFile := config.output
	if outputFile == "" {
		base := strings.TrimSuffix(inputFile, filepath.Ext(inputFile))
		if config.details {
			outputFile = base + ".details.md"
		} else {
			outputFile = base + ".md"
		}
	}

	// Check if file exists
	if !config.force {
		if _, err := os.Stat(outputFile); err == nil {
			return fmt.Errorf("output file %s already exists (use -f to overwrite)", outputFile)
		}
	}

	// Generate markdown
	md := generateMarkdown(item, config.details)

	// Write output
	if err := os.WriteFile(outputFile, []byte(md), 0644); err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}

	if config.verbose {
		fmt.Printf("Created %s\n", outputFile)
	}

	return nil
}

func generateMarkdown(item Item, includeDetails bool) string {
	var sb strings.Builder

	// Title
	fmt.Fprintf(&sb, "# %s: %s\n\n", item.Key.Value, item.Summary)
	fmt.Fprintf(&sb, "**Link:** [%s](%s)\n\n", item.Link, item.Link)

	// Overview
	sb.WriteString("## Overview\n\n")
	fmt.Fprintf(&sb, "- **Type:** %s\n", item.Type.Value)
	fmt.Fprintf(&sb, "- **Priority:** %s\n", item.Priority.Value)
	fmt.Fprintf(&sb, "- **Status:** %s\n", item.Status.Value)
	fmt.Fprintf(&sb, "- **Resolution:** %s\n", item.Resolution.Value)
	fmt.Fprintf(&sb, "- **Assignee:** %s\n", item.Assignee)
	fmt.Fprintf(&sb, "- **Reporter:** %s\n", item.Reporter)
	if len(item.Labels.Label) > 0 {
		fmt.Fprintf(&sb, "- **Labels:** %s\n", strings.Join(item.Labels.Label, ", "))
	}
	sb.WriteString("\n")

	// Dates
	sb.WriteString("## Dates\n\n")
	fmt.Fprintf(&sb, "- **Created:** %s\n", item.Created)
	fmt.Fprintf(&sb, "- **Updated:** %s\n", item.Updated)
	
	// Add custom date fields if details enabled
	if includeDetails {
		for _, cf := range item.CustomFields.CustomField {
			if strings.Contains(strings.ToLower(cf.CustomFieldName), "date") && len(cf.CustomFieldValues.CustomFieldValue) > 0 {
				val := cf.CustomFieldValues.CustomFieldValue[0].Value
				if val != "" {
					fmt.Fprintf(&sb, "- **%s:** %s\n", cf.CustomFieldName, val)
				}
			}
		}
	}
	sb.WriteString("\n")

	// Description/Details
	sb.WriteString("## Details\n\n")
	sb.WriteString(decodeHTML(item.Description))
	sb.WriteString("\n\n")

	// Comments
	if len(item.Comments.Comment) > 0 {
		sb.WriteString("## Comments\n\n")
		for _, comment := range item.Comments.Comment {
			fmt.Fprintf(&sb, "### %s\n\n", comment.Created)
			sb.WriteString(decodeHTML(comment.Value))
			sb.WriteString("\n\n")
		}
	}

	// Custom Fields (if details enabled)
	if includeDetails && len(item.CustomFields.CustomField) > 0 {
		sb.WriteString("## Custom Fields\n\n")
		for _, cf := range item.CustomFields.CustomField {
			// Skip date fields (already included above)
			if strings.Contains(strings.ToLower(cf.CustomFieldName), "date") {
				continue
			}
			
			// Skip empty fields
			if len(cf.CustomFieldValues.CustomFieldValue) == 0 {
				continue
			}

			hasContent := false
			for _, val := range cf.CustomFieldValues.CustomFieldValue {
				if val.Value != "" {
					hasContent = true
					break
				}
			}
			
			if !hasContent {
				continue
			}

			// Multi-value fields
			if len(cf.CustomFieldValues.CustomFieldValue) > 1 {
				fmt.Fprintf(&sb, "- **%s:** ", cf.CustomFieldName)
				var values []string
				for _, val := range cf.CustomFieldValues.CustomFieldValue {
					if val.Value != "" {
						values = append(values, val.Value)
					}
				}
				sb.WriteString(strings.Join(values, ", "))
				sb.WriteString("\n")
			} else {
				// Single value fields
				val := cf.CustomFieldValues.CustomFieldValue[0].Value
				if val != "" {
					fmt.Fprintf(&sb, "- **%s:** %s\n", cf.CustomFieldName, val)
				}
			}
		}
		
		// Add audit description if present
		for _, cf := range item.CustomFields.CustomField {
			if cf.CustomFieldName == "Audit Description" && len(cf.CustomFieldValues.CustomFieldValue) > 0 {
				val := cf.CustomFieldValues.CustomFieldValue[0].Value
				if val != "" {
					sb.WriteString("\n## Audit Description\n\n")
					sb.WriteString(decodeHTML(val))
					sb.WriteString("\n")
				}
			}
		}
	}

	return sb.String()
}

func decodeHTML(s string) string {
	// Basic HTML entity decoding and tag removal
	s = strings.ReplaceAll(s, "&lt;", "<")
	s = strings.ReplaceAll(s, "&gt;", ">")
	s = strings.ReplaceAll(s, "&quot;", "\"")
	s = strings.ReplaceAll(s, "&amp;", "&")
	s = strings.ReplaceAll(s, "&#8217;", "'")
	
	// Convert HTML tags to markdown
	s = strings.ReplaceAll(s, "<p>", "")
	s = strings.ReplaceAll(s, "</p>", "\n\n")
	s = strings.ReplaceAll(s, "<br/>", "\n")
	s = strings.ReplaceAll(s, "<br />", "\n")
	s = strings.ReplaceAll(s, "<b>", "**")
	s = strings.ReplaceAll(s, "</b>", "**")
	s = strings.ReplaceAll(s, "<ul>", "")
	s = strings.ReplaceAll(s, "</ul>", "")
	s = strings.ReplaceAll(s, "<li>", "- ")
	s = strings.ReplaceAll(s, "</li>", "\n")
	
	// Convert links
	s = convertHTMLLinks(s)
	
	// Convert images
	s = convertHTMLImages(s)
	
	// Clean up extra whitespace
	s = strings.TrimSpace(s)
	
	return s
}

func convertHTMLLinks(s string) string {
	// Simple regex-like replacement for <a href="url">text</a>
	for {
		start := strings.Index(s, "<a href=\"")
		if start == -1 {
			break
		}
		
		urlStart := start + 9
		urlEnd := strings.Index(s[urlStart:], "\"")
		if urlEnd == -1 {
			break
		}
		urlEnd += urlStart
		
		url := s[urlStart:urlEnd]
		
		textStart := strings.Index(s[urlEnd:], ">")
		if textStart == -1 {
			break
		}
		textStart += urlEnd + 1
		
		textEnd := strings.Index(s[textStart:], "</a>")
		if textEnd == -1 {
			break
		}
		textEnd += textStart
		
		text := s[textStart:textEnd]
		
		// Replace with markdown link
		markdown := fmt.Sprintf("[%s](%s)", text, url)
		s = s[:start] + markdown + s[textEnd+4:]
	}
	
	return s
}

func convertHTMLImages(s string) string {
	// Convert <img src="url" ... /> to ![Image](url)
	for {
		start := strings.Index(s, "<img src=\"")
		if start == -1 {
			// Also check for <span class="image-wrap">
			start = strings.Index(s, "<span class=\"image-wrap\"")
			if start == -1 {
				break
			}
			// Find the img tag inside
			imgStart := strings.Index(s[start:], "<img src=\"")
			if imgStart == -1 {
				break
			}
			start = start + imgStart
		}
		
		urlStart := start + 10
		urlEnd := strings.Index(s[urlStart:], "\"")
		if urlEnd == -1 {
			break
		}
		urlEnd += urlStart
		
		url := s[urlStart:urlEnd]
		
		// Find end of img tag
		tagEnd := strings.Index(s[urlEnd:], "/>")
		if tagEnd == -1 {
			tagEnd = strings.Index(s[urlEnd:], ">")
		}
		if tagEnd == -1 {
			break
		}
		tagEnd += urlEnd + 2
		
		// Check if there's a closing </span>
		endTag := s[tagEnd:]
		if strings.HasPrefix(strings.TrimSpace(endTag), "</span>") {
			tagEnd += strings.Index(s[tagEnd:], "</span>") + 7
		}
		
		// Replace with markdown image
		markdown := fmt.Sprintf("![Image](%s)", url)
		s = s[:start] + markdown + s[tagEnd:]
	}
	
	return s
}
