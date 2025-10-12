# Telugu SSML Builder

A comprehensive SSML (Speech Synthesis Markup Language) builder specifically designed for Telugu text, providing emotion detection, proper pauses, and stress handling for Telugu vowel signs (vothulu).

## Features

### 1. Emotion Detection Based on Punctuation
- **Exclamation marks (!)**: Adds excitement with higher pitch and medium rate
- **Double exclamation (!!)**: Very excited with fast rate and high pitch
- **Question marks (?)**: Questioning tone with rising pitch
- **Ellipsis (...)**: Suspense with slow rate and low pitch

### 2. Punctuation-Based Pauses
- **Commas (,)**: 300ms pause for natural speech flow
- **Semicolons (;)**: 500ms pause for sentence separation
- **Periods (.)**: 700ms pause between sentences
- **Colons (:)**: 400ms pause for emphasis

### 3. Vothulu (Vowel Signs) Stress Detection
Automatically detects and adds moderate emphasis to words containing Telugu vowel signs:
- ా, ి, ీ, ు, ూ, ె, ే, ొ, ో, ౌ
- అ, ఆ, ఇ, ఈ, ఉ, ఊ, ఎ, ఏ, ఒ, ఓ, ఔ
- ం, ః (anusvara and visarga)

### 4. Telugu Emotional Keywords
Pre-configured emotional keywords with appropriate emphasis levels:

#### Positive Emotions
- ఆనందం, సంతోషం, ఆశ్చర్యం, ఉత్సాహం
- ప్రేమ, ఆశ, విజయం, గర్వం
- ఆశీర్వాదం, ఆనందించు, ఆశ్చర్యపడు

#### Negative Emotions
- భయం, దుఃఖం, కోపం, అసహనం
- భయపడు, ఏడ్చు, కోపగించు, అసహనపడు
- భయంకర, భయానక, భయపెట్టు

#### Action/Intensity Words
- వేగంగా, నెమ్మదిగా, ధైర్యంగా, భయంకరంగా
- ఆశ్చర్యకరంగా, అద్భుతంగా, అత్యంత

#### Story Elements
- రాజు, రాణి, యువరాజు, యువరాణి
- మంత్రి, సేనాధిపతి, యోధుడు, యోధురాలు
- మంత్రుడు, సాధువు, సాధ్వి

## Usage

### Basic Usage

```go
package main

import (
    "fmt"
    "github.com/your-org/rio-go-model/internal/util"
)

func main() {
    teluguStory := "ఒక రోజు, ఒక బాలుడు అడవిలో నడుస్తున్నాడు. అతను ఒక భయంకరమైన భూతాన్ని చూశాడు!"
    
    // Using convenience function
    ssml := util.BuildTeluguSSMLFromStory(teluguStory)
    fmt.Println(ssml)
}
```

### Advanced Usage with Custom Builder

```go
package main

import (
    "fmt"
    "github.com/your-org/rio-go-model/internal/util"
)

func main() {
    // Create custom builder
    builder := util.NewTeluguSSMLBuilder()
    
    // Add custom emotional keywords
    builder.AddCustomEmotionalKeyword("అద్భుతం", "strong")
    builder.AddCustomEmotionalKeyword("అతిశయం", "moderate")
    
    // Add custom vothulu if needed
    builder.AddCustomVothulu('ృ') // Add custom vowel sign
    
    // Generate SSML
    ssml := builder.BuildTeluguSSML(teluguStory)
    fmt.Println(ssml)
}
```

## Output Example

For the input:
```
ఒక రోజు, ఒక బాలుడు అడవిలో నడుస్తున్నాడు. అతను ఒక భయంకరమైన భూతాన్ని చూశాడు!
```

The generated SSML will be:
```xml
<speak>
<p>
ఒక రోజు,<break time="300ms"/> ఒక <emphasis level="moderate">బాలుడు</emphasis> <emphasis level="moderate">అడవిలో</emphasis> <emphasis level="moderate">నడుస్తున్నాడు</emphasis>.<break time="700ms"/> 
అతను ఒక <emphasis level="strong">భయంకరమైన</emphasis> <emphasis level="moderate">భూతాన్ని</emphasis> <emphasis level="moderate">చూశాడు</emphasis><prosody rate="medium" pitch="high">!</prosody>
</p>
</speak>
```

## API Reference

### TeluguSSMLBuilder

#### Methods

- `NewTeluguSSMLBuilder() *TeluguSSMLBuilder`: Creates a new builder instance
- `BuildTeluguSSML(story string) string`: Converts Telugu text to SSML
- `AddCustomEmotionalKeyword(keyword, emphasisLevel string)`: Adds custom emotional keyword
- `AddCustomVothulu(vothulu rune)`: Adds custom vowel sign for stress detection
- `GetSupportedEmotions() []string`: Returns supported emphasis levels
- `GetSupportedPunctuationEmotions() []string`: Returns supported punctuation emotions

#### Emphasis Levels
- `"strong"`: High emphasis for very emotional words
- `"moderate"`: Medium emphasis for emotional words
- `"reduced"`: Low emphasis for subtle emotional words

#### Punctuation Emotions
- `"excited"`: Single exclamation mark
- `"very_excited"`: Double exclamation marks
- `"questioning"`: Question marks
- `"suspense"`: Ellipsis
- `"surprised"`: Question-exclamation combination

## Technical Details

### Vothulu Detection
The builder uses Unicode ranges to detect Telugu vowel signs and applies appropriate stress to words containing them.

### XML Escaping
All special XML characters are properly escaped to ensure valid SSML output.

### Whitespace Normalization
Text is normalized to handle different line ending formats and collapse multiple spaces.

### Word Boundary Detection
Emotional keyword matching uses word boundaries to avoid partial matches.

## Contributing

To add new emotional keywords or modify existing ones, you can:

1. Modify the `emotionalKeywords` map in the `NewTeluguSSMLBuilder()` function
2. Use the `AddCustomEmotionalKeyword()` method at runtime
3. Extend the `vothulu` map for additional vowel signs

## License

This code is part of the rio-go-model project and follows the same licensing terms.
