package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"fen-diagram/src/builders/fenlib"
)

var nagMap = map[string]string{
	"$1": "±", "$2": "∓", "$3": "+-", "$4": "-+",
	"$5": "+/=", "$6": "=/+", "$10": "=", "$13": "∞",
	"$14": "⩲", "$15": "⩱", "$16": "±", "$17": "∓",
	"$22": "⨀", "$36": "↑", "$40": "→", "$44": "⇆",
	"$132": "⌓", "$133": "⌓",
}

var resultMap = map[string]string{
	"1-0": "1-0 (White wins)",
	"0-1": "0-1 (Black wins)",
	"1/2-1/2": "1/2-1/2 (Draw)",
	"*": "*",
}

func convertNAG(input string) string {
	re := regexp.MustCompile(`\$\d+`)
	result := re.ReplaceAllStringFunc(input, func(match string) string {
		if v, ok := nagMap[match]; ok {
			if v != "" {
				return v
			}
			return match
		}
		return match + " (NAG " + match + ": not recognized)"
	})
	return result
}

func formatTime(h, m, s string, withLabel bool) string {
	var res string
	if withLabel {
		res = "Оставшееся время: "
	}
	hi := parseInt(h)
	mi := parseInt(m)
	if hi > 0 {
		res += fmt.Sprintf("%s ч ", h)
	}
	if mi > 0 || hi > 0 {
		res += fmt.Sprintf("%s мин ", m)
	}
	res += fmt.Sprintf("%s сек", s)
	return strings.TrimSpace(res)
}

func parseInt(s string) int {
	var n int
	fmt.Sscanf(s, "%d", &n)
	return n
}

func fixCommentSpacing(text string) string {
	// Fix 1: Castling "O-O" followed by piece/pawn without space
	// Must handle cases like "O- Ob6!" -> "O-O b6!" and "9.O- Ob6!" -> "9.O-O b6!"
	// First handle "O- O" -> "O-O O" pattern
	text = regexp.MustCompile(`O-\s+O`).ReplaceAllString(text, "O-O O")
	text = regexp.MustCompile(`O-O-\s+O`).ReplaceAllString(text, "O-O-O O")
	// Handle "O-Ob" -> "O-O b" and "O-O-Ob" -> "O-O-O b"
	text = strings.ReplaceAll(text, "O-Ob", "O-O b")
	text = strings.ReplaceAll(text, "O-O-Ob", "O-O-O b")
	// Handle "9.O- Ob6!" -> "9.O-O b6!" (space between O- and Ob")
	text = regexp.MustCompile(`(\d+\.)?\s*O-\s+Ob`).ReplaceAllString(text, "${1}O-O b")
	text = regexp.MustCompile(`(\d+\.)?\s*O-O-\s+Ob`).ReplaceAllString(text, "${1}O-O-O b")
	// Handle "O-O Ob6!" -> "O-O b6!" (after previous replacements)
	text = regexp.MustCompile(`O-O\s+Ob`).ReplaceAllString(text, "O-O b")
	text = regexp.MustCompile(`O-O-O\s+Ob`).ReplaceAllString(text, "O-O-O b")
	// Handle "O-OO-O" -> "O-O O-O" (two castling moves without space)
	text = strings.ReplaceAll(text, "O-OO-O", "O-O O-O")
	text = strings.ReplaceAll(text, "O-O-OO-O", "O-O-O O-O")
	text = strings.ReplaceAll(text, "O-OO-O-O", "O-O O-O-O")
	
	// Fix 2: Remove space incorrectly added before move number dots
	// "b7- b5" -> "b7-b5"
	re := regexp.MustCompile(`([a-h]\d)\s*-\s*([a-h]\d)`)
	text = re.ReplaceAllString(text, "$1-$2")
	
	// Fix 3: Pawn move ending followed by move number with dots
	// "h621." -> "h6 21."
	re2 := regexp.MustCompile(`([a-h][1-8])(\d{1,3}\.{1,3})`)
	text = re2.ReplaceAllString(text, "$1 $2")
	
	// Fix 4: Sentence ending followed by move number
	// "well.14..." -> "well. 14..."
	re3 := regexp.MustCompile(`([.!?])(\d+\.{1,3})`)
	text = re3.ReplaceAllString(text, "$1 $2")
	
	// Fix 5: Piece move ending without space before next move
	// "Rxc1f4!" -> "Rxc1 f4!"
	re4 := regexp.MustCompile(`([a-h]\d[+#=!?]*)([KQRBNOa-h])`)
	text = re4.ReplaceAllString(text, "$1 $2")
	
	// Fix 6: Move at end of comment followed by text without space
	// "8...Qc7As mentioned" -> "8...Qc7 As mentioned"
	// Match patterns like "Nf3As", "Qc7As", etc. (move followed by capital letter word)
	re5 := regexp.MustCompile(`([KQRBN][a-h][1-8][+#=!?]*|[KQRBN][a-h][a-h][+#=!?]*|[a-h][1-8][+#=!?]*)([A-Z][a-z])`)
	text = re5.ReplaceAllString(text, "$1 $2")
	
	// Also handle cases with dots like "8...Qc7As"
	re6 := regexp.MustCompile(`(\d+\.\.\.[KQRBNa-h][^\s]*)([A-Z])`)
	text = re6.ReplaceAllString(text, "$1 $2")
	
	// Fix 7: Text followed directly by move with dots (e.g., "after7...Bb7")
	// "after7...Bb7" -> "after 7...Bb7"
	re7 := regexp.MustCompile(`([a-z]+)(\d+\.\.\.[KQRBNa-h])`)
	text = re7.ReplaceAllString(text, "$1 $2")
	
	// Also handle "after7.Bb7" pattern
	re8 := regexp.MustCompile(`([a-z]+)(\d+\.[KQRBNa-h])`)
	text = re8.ReplaceAllString(text, "$1 $2")
	
	// Fix 8: Move with NAG symbol followed by text without space
	// "16.Qe3±With" -> "16.Qe3± With"
	// Also handle "Qe3±With" -> "Qe3± With"
	re9 := regexp.MustCompile(`(\d+\.[KQRBNa-h][^\s]*[+#=!?±]*|[KQRBNa-h][1-8][^\s]*[+#=!?±]*)([A-Z])`)
	text = re9.ReplaceAllString(text, "$1 $2")
	
	// Fix 9: Move with dots followed by text without space
	// "15...O-OOf course" -> "15...O-O Of course"
	re10 := regexp.MustCompile(`(\d+\.\.\.O-O)(O[a-z])`)
	text = re10.ReplaceAllString(text, "$1 $2")
	
	// Also handle O-O without number prefix
	re10b := regexp.MustCompile(`(O-O)(O[a-z])`)
	text = re10b.ReplaceAllString(text, "$1 $2")
	
	// Also handle other moves with dots followed by text
	re10c := regexp.MustCompile(`(\d+\.\.\.[KQRBNa-h][^\s]*)([A-Z])`)
	text = re10c.ReplaceAllString(text, "$1 $2")
	
	// Fix 10: Text with number followed by move (e.g., "worse1 6.dxe5" -> "worse 16.dxe5")
	// Handle case where text ends with digit, then space, then digit+move
	re10d := regexp.MustCompile(`([a-z]+)(\d)\s+(\d\.[KQRBNa-h][^\s]*)`)
	text = re10d.ReplaceAllString(text, "$1 $2$3")
	
	// Also handle text directly followed by move (no space): "worse16.dxe5" -> "worse 16.dxe5"
	re10e := regexp.MustCompile(`([a-z]+)(\d+\.[KQRBNa-h][^\s]*)`)
	text = re10e.ReplaceAllString(text, "$1 $2")
	
	// Fix 11: Text with number followed by move with dots (e.g., "Following1 6...Qb8" -> "Following 16...Qb8")
	re11a := regexp.MustCompile(`([a-z]+)(\d)\s+(\d\.\.\.[KQRBNa-h][^\s]*)` )
	text = re11a.ReplaceAllString(text, "$1 $2$3")
	
	// Also handle text directly followed by move with dots: "Following16...Qb8" -> "Following 16...Qb8"
	re11b := regexp.MustCompile(`([a-z]+)(\d+\.\.\.[KQRBNa-h][^\s]*)` )
	text = re11b.ReplaceAllString(text, "$1 $2")
	
	// Fix 12: NAG symbol followed by move without space (e.g., "17.Ne4+-16. Nb5" -> "17.Ne4+- 16. Nb5")
	re12 := regexp.MustCompile(`([+#=!?±⩲⩱\-]+)(\d+\.)`)
	text = re12.ReplaceAllString(text, "$1 $2")
	
	// Fix 13: Move followed by text without space (e.g., "17.Be4I still" -> "17.Be4 I still")
	re13 := regexp.MustCompile(`(\d+\.[KQRBNa-h][^\s]*)([A-Z])` )
	text = re13.ReplaceAllString(text, "$1 $2")
	
	return text
}

func convertInlineMarkers(text string) string {
	t := text
	t = regexp.MustCompile(`\[%emt\s+([0-9:]+)\]`).ReplaceAllStringFunc(t, func(m string) string {
		parts := regexp.MustCompile(`\[%emt\s+([0-9:]+)\]`).FindStringSubmatch(m)
		if len(parts) == 2 {
			timeParts := strings.Split(parts[1], ":")
			if len(timeParts) == 3 {
				return fmt.Sprintf("(%s)", formatTime(timeParts[0], timeParts[1], timeParts[2], false))
			}
		}
		return m
	})
	t = regexp.MustCompile(`\[%clk\s+([0-9:]+)\]`).ReplaceAllStringFunc(t, func(m string) string {
		parts := regexp.MustCompile(`\[%clk\s+([0-9:]+)\]`).FindStringSubmatch(m)
		if len(parts) == 2 {
			timeParts := strings.Split(parts[1], ":")
			if len(timeParts) == 3 {
				return fmt.Sprintf("(Оставшееся время: %s)", formatTime(timeParts[0], timeParts[1], timeParts[2], true))
			}
		}
		return m
	})
	t = regexp.MustCompile(`\[%eval\s+([^\]]+)\]`).ReplaceAllStringFunc(t, func(m string) string {
		parts := regexp.MustCompile(`\[%eval\s+([^\]]+)\]`).FindStringSubmatch(m)
		if len(parts) == 2 {
			if strings.HasPrefix(parts[1], "#") {
				return fmt.Sprintf("(Мат в %s хода)", parts[1][1:])
			}
			return fmt.Sprintf("(%s)", parts[1])
		}
		return m
	})
	t = regexp.MustCompile(`\[%M\s*(\d+)\]`).ReplaceAllStringFunc(t, func(m string) string {
		parts := regexp.MustCompile(`\[%M\s*(\d+)\]`).FindStringSubmatch(m)
		if len(parts) == 2 {
			return fmt.Sprintf("(Мат в %s ходов)", parts[1])
		}
		return m
	})
	return t
}

func splitGames(content string) []string {
	content = strings.TrimSpace(content)
	if content == "" {
		return []string{}
	}
	re := regexp.MustCompile(`(?m)^\[Event\s+"[^"]*"\]\s*\n`)
	parts := re.Split(content, -1)
	var games []string
	for _, p := range parts {
		if strings.TrimSpace(p) != "" {
			games = append(games, p)
		}
	}
	if len(games) == 0 {
		games = append(games, content)
	}
	return games
}

func parseGameTags(gameText string) GameTags {
	re := regexp.MustCompile(`\[(\w+)\s+"([^"]*)"`)
	matches := re.FindAllStringSubmatch(gameText, -1)
	tags := make(GameTags)
	for _, m := range matches {
		if len(m) >= 3 {
			tags[m[1]] = m[2]
		}
	}
	return tags
}

type GameTags map[string]string

func translateTags(tags GameTags) GameTags {
	return tags
}

func translateCommentsInline(movetext string) string {
	return movetext
}

func buildChapterTitle(tags GameTags, chapterNum int) string {
	white := tags["White"]
	black := tags["Black"]
	site := tags["Site"]
	date := tags["Date"]
	
	// Handle malformed PGN where Black field contains entire title
	if strings.Contains(black, " vs. ") {
		parts := strings.Split(black, " vs. ")
		if len(parts) == 2 {
			whiteFromBlack := parts[0]
			blackFull := parts[1]
			
			// Always use white from black field if white is "Model Games" or empty
			// Otherwise, don't overwrite it
			if white == "Model Games" || white == "" {
				white = whiteFromBlack
			}
			
			// Always extract black player name from blackFull
			locationParts := strings.Split(blackFull, ", ")
			if len(locationParts) >= 2 {
				lastPart := locationParts[len(locationParts)-1]
				if idx := strings.Index(lastPart, " "); idx >= 0 {
					cityFromBlack := lastPart[:idx]
					yearFromBlack := lastPart[idx+1:]
					
					// Always update black to just the player name
					black = strings.Join(locationParts[:len(locationParts)-1], ", ")
					
					if site == "?" || site == "" {
						site = cityFromBlack
					}
					if date == "????.??.??" || date == "" {
						date = yearFromBlack + ".??.??"
					}
				} else {
					black = blackFull
					if site == "?" || site == "" {
						site = lastPart
					}
				}
			} else {
				// No comma in blackFull, use it as-is
				black = blackFull
			}
		}
	}
	
	if white == "" || white == "Model Games" {
		white = "Unknown player"
	}
	if black == "" {
		black = "Unknown player"
	}
	
	city := extractCityFromSite(site)
	year := extractYearFromDate(date)
	
	if year == "" {
		if eventDate := tags["EventDate"]; eventDate != "" {
			year = extractYearFromDate(eventDate)
		}
	}
	
	title := fmt.Sprintf("Chapter %d: %s vs. %s", chapterNum, white, black)
	if city != "" {
		title += ", " + city
	}
	if year != "" {
		title += " " + year
	}
	return title
}

func extractCityFromSite(site string) string {
	if idx := strings.Index(site, "("); idx >= 0 {
		site = site[:idx]
	}
	if idx := strings.Index(site, "["); idx >= 0 {
		site = site[:idx]
	}
	if strings.TrimSpace(site) == "?" || strings.TrimSpace(site) == "" {
		return ""
	}
	parts := strings.Split(site, ",")
	return strings.TrimSpace(parts[0])
}

func extractYearFromDate(date string) string {
	if len(date) >= 4 {
		year := date[:4]
		for _, c := range year {
			if c < '0' || c > '9' {
				return ""
			}
		}
		return year
	}
	return ""
}

// escapeMoveDots adds a backslash before dots that follow move numbers
// to prevent Markdown from interpreting them as numbered lists
func escapeMoveDots(s string) string {
	result := []byte(s)
	// Find all positions where we have digit(s) followed by "." and space
	// Then add a backslash before the dot if not already present
	i := 0
	for i < len(result) {
		// Check if current char is a digit
		if result[i] >= '0' && result[i] <= '9' {
			// Find the end of the digit sequence
			j := i
			for j < len(result) && result[j] >= '0' && result[j] <= '9' {
				j++
			}
			// Check if followed by "." and space
			if j < len(result) && result[j] == '.' {
				// Check if the dot is already escaped
				if i > 0 && result[i-1] == '\\' {
					// Already escaped, skip
					i = j + 1
					continue
				}
				// Need to insert a backslash before the dot
				newResult := make([]byte, len(result)+1)
				copy(newResult[:j], result[:j])
				newResult[j] = '\\'
				copy(newResult[j+1:], result[j:])
				result = newResult
				i = j + 2 // Skip past the inserted backslash and the dot
			} else {
				i = j
			}
		} else {
			i++
		}
	}
	return string(result)
}

func extractMoveText(gameText string, inlineImages bool, numberedLists bool, dashLists bool) string {
	// formatLists formats numbered and dash lists based on flags
	formatLists := func(text string, numbered bool, dash bool) string {
		if numbered {
			// Handle numbered lists: replace "  i. " or "i. " at start with "\n  i. "
			for i := 1; i <= 10; i++ {
				// Case: "  i. " (already has spaces)
				old := fmt.Sprintf("  %d. ", i)
				newStr := fmt.Sprintf("\n  %d. ", i)
				text = strings.Replace(text, old, newStr, -1)
				// Case: "i. " at start or after newline
				old2 := fmt.Sprintf("\n%d. ", i)
				text = strings.Replace(text, old2, newStr, -1)
				// Case: " i. " with single space
				old3 := fmt.Sprintf(" %d. ", i)
				text = strings.Replace(text, old3, newStr, -1)
			}
			// Then, if text follows a list item with multiple spaces, add newline
			re := regexp.MustCompile(`(\n  \d+\. [^\n]*?) {2,}(\w)`)
			text = re.ReplaceAllString(text, "$1\n$2")
		}
		// Handle dash lists if flag is set
		if dash {
			// Handle dash lists: replace any whitespace + "- " with "\n  - "
			reDash := regexp.MustCompile(`\s+- `)
			text = reDash.ReplaceAllString(text, "\n  - ")
		}
		// Clean up: remove extra newlines but preserve single empty lines
		for strings.Contains(text, "\n\n\n") {
			text = strings.Replace(text, "\n\n\n", "\n\n", -1)
		}
		// Remove leading newline if any
		text = strings.TrimPrefix(text, "\n")
		return text
	}

	// Extract FEN if present
	fenRe := regexp.MustCompile(`\[FEN\s+"([^"]*)"\]`)
	fenMatch := fenRe.FindStringSubmatch(gameText)
	var fen string
	if len(fenMatch) > 1 {
		fen = fenMatch[1]
	}
	
	// Remove all PGN tags
	re := regexp.MustCompile(`\[[^\]]+\]\s*`)
	movetext := re.ReplaceAllString(gameText, "")
	movetext = strings.TrimSpace(movetext)
	
	// Convert results
	for k, v := range resultMap {
		movetext = strings.ReplaceAll(movetext, k, v)
	}
	
	// Convert NAG codes
	movetext = convertNAG(movetext)

	// Remove curly braces from comments (keep content inside)
	reComment := regexp.MustCompile(`\{([^}]*)\}`)
	movetext = reComment.ReplaceAllStringFunc(movetext, func(match string) string {
		// Extract content inside braces, remove braces
		content := strings.TrimSpace(match[1:len(match)-1])
		// Clean up extra newlines
		for strings.Contains(content, "\n\n") {
			content = strings.Replace(content, "\n\n", "\n", -1)
		}
		// Remove leading newline if any
		content = strings.TrimPrefix(content, "\n")
		// Apply formatLists to handle dash lists and optionally numbered lists
		content = formatLists(content, numberedLists, dashLists)
		return content
	})
  
	// FEN will be added at the end, after formatLists is applied
	// This prevents formatLists from breaking the FEN string

	lines := strings.Split(movetext, "\n")
	var result []string
	var currentMoveNum string
	var currentMoveLine string
	var pendingComments []string
	i := 0
	fenAdded := fen != "" && !inlineImages // Don't set fenAdded if using inline images
	
reMove := regexp.MustCompile(`^(\d+)(\.{1,3})\s+(.+)$`)
	// reComment is already defined above

	// Helper to check if a line looks like a real chess move
	isRealMove := func(line string) bool {
		return reMove.MatchString(line)
	}
	
	// Check if there are any real moves in the text
	hasMoves := false
	for _, line := range lines {
		if isRealMove(line) {
			hasMoves = true
			break
		}
	}
	
	// If no moves found, return only the text without move numbers
	if !hasMoves {
		// Return the text as-is (without move formatting)
		return strings.Join(lines, "\n")
	}
	
	for i < len(lines) {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			// Don't skip the empty line immediately after FEN
			if fenAdded && i > 0 && strings.Contains(lines[i-1], "**FEN:**") {
				fenAdded = false
				// Keep this empty line
			} else {
				i++
				continue
			}
		}
		
		m := reMove.FindStringSubmatch(line)
		if m != nil {
			num := m[1]
			dots := m[2]
			moveText := m[3]
			
			moveText = fixCommentSpacing(moveText)
			
			var comments []string
			commentMatches := reComment.FindAllStringSubmatch(moveText, -1)
			for _, cm := range commentMatches {
				if len(cm) > 1 {
					comment := fixCommentSpacing(cm[1])
					comments = append(comments, comment)
				}
			}
			
		moveTextClean := reComment.ReplaceAllString(moveText, "")
		moveTextClean = strings.TrimSpace(moveTextClean)
		
		// Check if there's actual move text (remove result symbols for check)
		moveTextCheck := moveTextClean
		for k := range resultMap {
			moveTextCheck = strings.ReplaceAll(moveTextCheck, k, "")
		}
		moveTextCheck = strings.TrimSpace(moveTextCheck)
		
		isBlack := dots == "..."
		
		// If no actual move text (only result or empty), skip this line
		if moveTextCheck == "" {
			i++
			continue
		}
		
		if !isBlack && currentMoveNum == num && currentMoveLine != "" {
			isBlack = true
		}
		
		if !isBlack {
			if currentMoveNum == num && currentMoveLine != "" {
				currentMoveLine += " " + moveTextClean
				pendingComments = append(pendingComments, comments...)
			} else {
				if currentMoveLine != "" {
					result = append(result, currentMoveLine)
					if len(pendingComments) > 0 {
						result = append(result, "")
						for _, c := range pendingComments {
							result = append(result, "  "+c)
						}
						result = append(result, "")
					}
				}
				currentMoveNum = num
				currentMoveLine = num + ". " + moveTextClean
				pendingComments = comments
			}
		} else {
			// Skip if no actual move text (only result)
			if moveTextCheck == "" {
				i++
				continue
			}
			if currentMoveNum == num && currentMoveLine != "" {
				currentMoveLine += " " + moveTextClean
				pendingComments = append(pendingComments, comments...)
			} else {
					if currentMoveLine != "" {
						result = append(result, currentMoveLine)
					}
					currentMoveNum = num
			currentMoveLine = num + ". " + moveTextClean
					pendingComments = comments
				}
			}
		} else {
			if currentMoveLine != "" {
				result = append(result, currentMoveLine)
				if len(pendingComments) > 0 {
					result = append(result, "")
					for _, c := range pendingComments {
						result = append(result, "  "+c)
					}
					result = append(result, "")
				}
				currentMoveLine = ""
				pendingComments = nil
			}
			result = append(result, line)
		}
		i++
	}
	
	if currentMoveLine != "" {
		result = append(result, currentMoveLine)
		if len(pendingComments) > 0 {
			result = append(result, "")
			for _, c := range pendingComments {
				result = append(result, "  "+c)
			}
			result = append(result, "")
		}
	}
	
	// Remove lines that match "1. -- *" or "1. *" patterns (no real moves)
	// Also match "1. -- *" with optional spaces
	reNoMoves := regexp.MustCompile(`^\s*\d+\.\s*--\s*\*?\s*$`)
	reNoMoves2 := regexp.MustCompile(`^\s*\d+\.\s*\*?\s*$`)
	filtered := make([]string, 0, len(result))
	for _, line := range result {
		if reNoMoves.MatchString(line) || reNoMoves2.MatchString(line) {
			continue
		}
		filtered = append(filtered, line)
	}
	
	// Apply formatLists to the final result if needed
	resultStr := strings.Join(filtered, "\n")
	resultStr = formatLists(resultStr, numberedLists, dashLists)
	
	// Escape dots in move numbers to prevent Markdown list interpretation
	// Process the string to add backslash before dots that follow move numbers
	resultStr = escapeMoveDots(resultStr)

	// Remove bracket placeholders like @@StartBracket@@__TEXT__@@EndBracket@@
	// Handle optional space after "End" to catch cases where fixCommentSpacing added a space
	reBracket := regexp.MustCompile(`@@StartBracket@@([^@]+)@@End\s*Bracket@@`)
	resultStr = reBracket.ReplaceAllString(resultStr, `$1`)

	// Add FEN at the beginning if present (with empty line after)
	if fen != "" {
		if inlineImages {
			// When using inline images, just add FEN without the label
			// It will be replaced by replaceFENWithImages later
			resultStr = fen + "\n\n" + resultStr
		} else {
			resultStr = "**FEN:** `" + fen + "`\n\n" + resultStr
		}
	}
	
	return resultStr
}

func main() {
	src := flag.String("src", "", "Input PGN file")
	out := flag.String("out", "", "Output Markdown file")
	skipChapters := flag.Int("skip", 0, "Number of chapters to skip from the beginning")
	chaptersLimit := flag.Int("chapters", 0, "Number of chapters to process (first N chapters after skip)")
	chapterNumbersStr := flag.String("chapter-numbers", "", "Comma-separated list of chapter numbers to process (1-based)")
	inlineImages := flag.Bool("inline-images", false, "Replace FEN strings with inline base64 images")
	numberedLists := flag.Bool("numbered-lists", true, "Format numbered lists with 2-space indent (default: true)")
	dashLists := flag.Bool("dash-lists", false, "Format dash lists with 2-space indent (default: false)")
	exam := flag.String("exam", "", "Process only chapters where White field exactly matches this value. If chapter contains FEN, output only FEN. Otherwise, format normally.")
	flag.Parse()

	if *src == "" || *out == "" {
		fmt.Fprintf(os.Stderr, "Usage: %s -src <pgn_path> -out <md_path> [-skip <n>] [-chapters <n>] [-chapter-numbers <list>] [-inline-images]\n", os.Args[0])
		os.Exit(1)
	}

	// Parse chapter-numbers if provided
	chapterNumbers := make(map[int]bool)
	if *chapterNumbersStr != "" {
		parts := strings.Split(*chapterNumbersStr, ",")
		for _, p := range parts {
			p = strings.TrimSpace(p)
			num, err := strconv.Atoi(p)
			if err == nil && num > 0 {
				chapterNumbers[num] = true
			}
		}
	}

	raw, err := os.ReadFile(*src)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	games := splitGames(string(raw))

	output := ""
	chapterNum := 0
	var prevWhite string
	for i, g := range games {
		// If chapter-numbers is provided, filter by index (1-based)
		if len(chapterNumbers) > 0 {
			if !chapterNumbers[i+1] {
				continue
			}
			chapterNum++
		} else {
			// Apply skip
			if i < *skipChapters {
				continue
			}
			chapterNum++
			// Apply chapters limit
			if *chaptersLimit > 0 && chapterNum > *chaptersLimit {
				break
			}
		}

		tags := parseGameTags(g)
		
		// If exam filter is set, check if White or Black matches exactly
		if *exam != "" {
			rawWhite := strings.TrimSpace(tags["White"])
			rawBlack := strings.TrimSpace(tags["Black"])
			if rawWhite != *exam && rawBlack != *exam {
				continue
			}
		}
		
		// Detect if this is a book chapter format BEFORE translation
		rawWhite := tags["White"]
		rawBlack := tags["Black"]
		isBookChapter := strings.HasPrefix(strings.TrimSpace(rawBlack), "Chapter ") || 
		                  strings.Contains(rawWhite, "Introduction") || 
		                  strings.Contains(rawWhite, "Part") ||
		                  strings.HasPrefix(strings.TrimSpace(rawWhite), "Chapter ") ||
		                  strings.Contains(rawWhite, " vs. ") ||
		                  strings.Contains(rawBlack, " vs. ")
		
		tags = translateTags(tags)
		white := tags["White"]
		black := tags["Black"]

// Check if chapter has FEN (for exam mode)
		hasFEN := false
		if *exam != "" {
			fenRe := regexp.MustCompile(`\[FEN\s+"([^"]*)"\]`)
			fenMatch := fenRe.FindStringSubmatch(g)
			hasFEN = len(fenMatch) > 1
		}

		if isBookChapter {
			// Book chapter format: ## White, ### Black
			// If "vs." is found in White or Black, use the full text from that field
			whiteTitle := white
			blackTitle := black
			
			if strings.Contains(white, " vs. ") {
				whiteTitle = white
				// If vs. is in White, use it as ## and don't add separate ###
				if white != prevWhite {
					output += "## " + whiteTitle + "\n\n"
					prevWhite = white
				}
			} else if strings.Contains(black, " vs. ") {
				// If vs. is in Black, use it as ###
				blackTitle = black
				if white != prevWhite {
					output += "## " + whiteTitle + "\n\n"
					prevWhite = white
				}
				output += "### " + blackTitle + "\n\n"
			} else {
				// Normal case: no "vs." in either field
				if white != prevWhite {
					output += "## " + whiteTitle + "\n\n"
					prevWhite = white
				}
				if blackTitle != "" {
					output += "### " + blackTitle + "\n\n"
				}
			}
		} else {
			// Normal game format: ## White vs. Black, City Year
			title := buildChapterTitle(tags, chapterNum)
			output += "## " + title + "\n\n"
		}

		// Exam mode: if chapter has FEN, output only FEN (no prefix)
		if *exam != "" && hasFEN {
			fenRe := regexp.MustCompile(`\[FEN\s+"([^"]*)"\]`)
			fenMatch := fenRe.FindStringSubmatch(g)
			if len(fenMatch) > 1 {
				fen := fenMatch[1]
				if *inlineImages {
					b64, err := fenlib.GenerateDiagramBase64(fen)
					if err == nil {
						fen = fmt.Sprintf("data:image/png;base64,%s", b64)
					}
				}
				output += fen + "\n\n"
			}
		} else {
			moves := extractMoveText(g, *inlineImages, *numberedLists, *dashLists)

			// Replace FEN strings with inline base64 images if flag is set
			if *inlineImages {
				moves = replaceFENWithImages(moves, false)
			}

			output += moves + "\n\n"
		}
	}

	err = os.WriteFile(*out, []byte(output), 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing file: %v\n", err)
		os.Exit(1)
	}
}

// replaceFENWithImages finds FEN strings in text and replaces them with base64-encoded PNG images
// If wrapInMarkdown is true, wraps base64 string in ![Position](data:image/png;base64,...)
// If false, returns data:image/png;base64,... format without markdown wrapper
func replaceFENWithImages(text string, wrapInMarkdown bool) string {
	// Pattern to match FEN strings: exactly 6 space-separated fields
	// Field 1: piece placement (8 rows separated by /)
	// Fields 2-6: active color, castling, en passant, halfmove, fullmove
	fenPattern := regexp.MustCompile(`([KQRBNPkqrbnp1-8/]+)\s+([wb])\s+([KQkq-]+)\s+([a-h1-8-]+)\s+(\d+)\s+(\d+)`)

	return fenPattern.ReplaceAllStringFunc(text, func(match string) string {
		// Validate and generate base64 image
		b64, err := fenlib.GenerateDiagramBase64(match)
		if err != nil {
			// Return original FEN if generation fails
			return match
		}
		if wrapInMarkdown {
			return fmt.Sprintf("![Position](data:image/png;base64,%s)", b64)
		}
		return fmt.Sprintf("data:image/png;base64,%s", b64)
	})
}
