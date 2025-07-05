# WhatsApp Template Variables - Constraints & Limitations

## Variable Content Restrictions

### Newlines and Line Breaks
- **Variables CANNOT contain newlines** (`\n`), tab characters (`\t`), or carriage returns (`\r`)
- Error message: "Param text cannot have new-line/tab characters"
- If you need multi-line content, you must structure it in the template itself with separate variables

### Spaces
- Variables cannot contain more than 4 consecutive spaces
- Regular spaces within variables are allowed

## Character Limits for Variables

### Individual Variable Limits
- No specific character limit per individual variable ({{1}}, {{2}}, etc.)
- The limitation comes from the overall template character limits:
  - **Template Body**: 1,024 characters total (including all variables)
  - **Header**: 60 characters total (maximum 1 variable allowed)
  - **Footer**: 60 characters total (NO variables allowed)

### Variable Count Limits
- Up to **100 variables** per template maximum
- Body section: Up to **15 variables** maximum
- Header section: Maximum **1 variable** only
- Footer section: **No variables** allowed

## Special Characters and Formatting

### Allowed Characters
- Letters, numbers, and most special characters are allowed
- Emojis are supported in variables
- Standard punctuation marks are allowed

### Forbidden Characters
- Variable parameters cannot contain special characters such as `#`, `$`, or `%`
- No newline characters (`\n`, `\r`)
- No tab characters (`\t`)

## Variable Structure Requirements

### Proper Format
- Must use double curly braces: `{{1}}`, `{{2}}`, `{{3}}`
- Variables must be sequential (no skipping numbers)
- ✅ Correct: `{{1}}` then `{{2}}` then `{{3}}`
- ❌ Incorrect: `{{1}}` then `{{3}}` (skipping {{2}})

### Spacing Rules
- Variables cannot be adjacent without separating text
- There must be characters (separated by spaces) between variables
- For every 'x' variables, there must be 2x+1 non-variable words

## Component-Specific Differences

### Header Variables
- Limited to 60 characters total
- Maximum 1 variable allowed
- Cannot contain emojis, asterisks, formatting markup, or newlines in TEXT headers
- Media headers don't support text variables

### Body Variables
- 1,024 characters total limit
- Up to 15 variables allowed
- Supports markdown formatting (`*bold*`, `_italic_`, `~strikethrough~`, ``` `monospace` ```)
- Emojis and special characters allowed

### Footer Variables
- 60 characters total limit
- **NO variables allowed at all**
- Plain text only, no formatting

## Template vs Interactive Messages

| **Template Messages** | **Interactive Messages** |
|----------------------|--------------------------|
| **Approval**: Required from WhatsApp | **Approval**: None required |
| **Usage**: Can initiate conversations outside 24-hour window | **Usage**: Only within 24-hour customer service window |
| **Cost**: Higher cost per message | **Cost**: Lower cost (session messages) |
| **Flexibility**: Fixed, pre-approved structure | **Flexibility**: Dynamic content, can be generated on-the-fly |
| **Variables**: Up to 15 with strict formatting rules | **Variables**: No variable restrictions |
| **Content**: Static structure with variable substitution | **Content**: Fully dynamic lists, buttons, and interactive elements |
| **Character Limit**: 1024 characters | **Character Limit**: More flexible limits |
| **Loops/Logic**: Not supported | **Loops/Logic**: Can be implemented in code |

## Tide Notification Template Options

### Option 1: Pre-formatted in template
```
🌊 Daily Tide Report for {{1}}

⬆️ High Tide: {{2}} ({{3}}m)
⬇️ Low Tide: {{4}} ({{5}}m)
⬆️ High Tide: {{6}} ({{7}}m)

📍 Risco del Paso, Fuerteventura
```

### Option 2: Single formatted variable
```
🌊 Daily Tide Report

{{1}}

📍 Risco del Paso, Fuerteventura
```

Where `{{1}}` contains: `"⬆️ High: 6:30 AM (8.2m) ⬇️ Low: 12:45 PM (1.1m) ⬆️ High: 7:15 PM (7.9m)"` (no newlines)

## Best Practices

### Authentication Templates
- Verification codes: maximum 25 characters
- No URLs, media, or emojis allowed in variables

### Content Organization
- Use separate variables for different pieces of information
- Avoid cramming multiple data points into single variables
- Plan for localization if supporting multiple languages

### Error Prevention
- Always validate variable content before sending
- Remove newlines and excessive spaces from dynamic content
- Test templates thoroughly before approval submission

### Template Approval Considerations
- Must provide sample content for all variables during template submission
- Samples help WhatsApp understand the intended use case
- Variables without samples may cause template rejection
- Variable content should be meaningful and relevant
- Avoid placeholder text like "Lorem ipsum"
- Ensure proper grammar and spelling in samples

## Current Project Implementation

### Interactive Template
- Template SID: `HX6f156e3466407a835bef6505f85cf9b1`
- Question: "What would you like to do?"
- Buttons: 
  - ID: `tides`, Text: `🌊 Get Current Tides`
  - ID: `start`, Text: `🔔 Enable Notifications`

### Webhook Button Handling
- Extract `ButtonPayload` for button ID
- Extract `ButtonText` for display text
- Check `MessageType=button` to detect button responses
- Use button ID for command routing instead of display text