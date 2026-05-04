package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"fen-diagram/builders/fenlib"
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

func extractMoveText(gameText string) string {
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
	
	// Add FEN at the beginning if present (with empty line after)
	if fen != "" {
		movetext = "**FEN:** `" + fen + "`\n\n" + movetext
	}

	lines := strings.Split(movetext, "\n")
	var result []string
	var currentMoveNum string
	var currentMoveLine string
	var pendingComments []string
	i := 0
	fenAdded := fen != ""
	
	reMove := regexp.MustCompile(`^(\d+)(\.{1,3})\s+(.+)$`)
	reComment := regexp.MustCompile(`\{([^}]*)\}`)
	
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
			
			isBlack := dots == "..."
			
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
	
	return strings.Join(result, "\n")
}

func main() {
	src := flag.String("src", "", "Input PGN file")
	out := flag.String("out", "", "Output Markdown file")
	skipChapters := flag.Int("skip", 0, "Number of chapters to skip from the beginning")
	inlineImages := flag.Bool("inline-images", false, "Replace FEN strings with inline base64 images")
	flag.Parse()

	if *src == "" || *out == "" {
		fmt.Fprintf(os.Stderr, "Usage: %s -src <pgn_path> -out <md_path> [-skip <n>] [-inline-images]\n", os.Args[0])
		os.Exit(1)
	}

	raw, err := ioutil.ReadFile(*src)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	games := splitGames(string(raw))

	output := ""
	chapterNum := 0
	for i, g := range games {
		if i < *skipChapters {
			continue
		}
		chapterNum++

		tags := parseGameTags(g)
		tags = translateTags(tags)

		black := tags["Black"]
		if black == "" {
			black = fmt.Sprintf("Партия %d", chapterNum)
		}
		title := buildChapterTitle(tags, chapterNum)

		output += "## " + title + "\n\n"

		moves := extractMoveText(g)
		
		// Replace FEN strings with inline base64 images if flag is set
		if *inlineImages {
			moves = replaceFENWithImages(moves)
		}
		
		output += moves + "\n\n"
	}

	err = ioutil.WriteFile(*out, []byte(output), 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing file: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Done:", *out)
	fmt.Println("Total chapters processed:", chapterNum)
}

// replaceFENWithImages finds FEN strings in text and replaces them with base64-encoded PNG images
func replaceFENWithImages(text string) string {
	// Pattern to match FEN strings: exactly 6 space-separated fields
	// Field 1: piece placement (8 rows separated by /)
	// Fields 2-6: active color, castling, en passant, halfmove, fullmove
	fenPattern := regexp.MustCompile(`[KQRBNPkqrbnp1-8/]+\s+[wb]\s+[KQkq-]+\s+[a-h1-8-]+\s+\d+\s+\d+`)
	
	return fenPattern.ReplaceAllStringFunc(text, func(match string) string {
		// Validate and generate base64 image
		b64, err := fenlib.GenerateDiagramBase64(match)
		if err != nil {
			// Return original FEN if generation fails
			return match
		}
		return fmt.Sprintf("![Position](data:image/png;base64,%s)", b64)
	})
}
