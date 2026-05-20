package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

// Test buildChapterTitle (Rule 8)
func TestBuildChapterTitle_Basic(t *testing.T) {
	tags := GameTags{
		"White": "Nepomniachtchi, Ian",
		"Black": "Alekseenko, Kirill",
		"Site":  "Ekaterinburg",
		"Date":  "2021.05.15",
	}
	title := buildChapterTitle(tags, 1)
	// Since neither White nor Black contains "vs.", use comma separator
	expected := "Chapter 1: Nepomniachtchi, Ian, Alekseenko, Kirill, Ekaterinburg 2021"
	if title != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, title)
	}
}

func TestBuildChapterTitle_WithCommaInSite(t *testing.T) {
	tags := GameTags{
		"White": "Carlsen, Magnus",
		"Black": "Nakamura, Hikaru",
		"Site":  "Moscow, Russia",
		"Date":  "2013.10.05",
	}
	title := buildChapterTitle(tags, 2)
	if !strings.Contains(title, "Moscow") {
		t.Error("City should be extracted from site with comma")
	}
}

func TestBuildChapterTitle_NoDate(t *testing.T) {
	tags := GameTags{
		"White": "Player A",
		"Black": "Player B",
	}
	title := buildChapterTitle(tags, 4)
	if strings.Contains(title, "2021") {
		t.Error("Should not have year without date")
	}
}

func TestExtractMoveText_PairedMovesNoComments(t *testing.T) {
	input := "[Result \"*\"]\n\n1. e4 e5 2. d4 d5\n*"
	output := extractMoveText(input, false, true, false)
	// Check for escaped dots (backslash-dot-space pattern)
	if !strings.Contains(output, `1\. e4`) || !strings.Contains(output, `2\. d4`) {
		t.Errorf("Paired moves not found in output. Got:\n%s", output)
	}
}

func TestConvertNAG_AllCodes(t *testing.T) {
	text := "1. e4 $1 $2 $3 $4 $5 $6 $10 $13 $14 $15 $16 $17 $22 $36 $40 $132 $133"
	result := convertNAG(text)
	expected := []string{"±", "∓", "+-", "-+", "+/=", "=/+", "=", "∞", "⩲", "⩱", "±", "∓", "⨀", "↑", "→", "⌓"}
	for _, exp := range expected {
		if !strings.Contains(result, exp) {
			t.Errorf("NAG code not converted: %s", exp)
		}
	}
}

func TestSplitGames_Multiple(t *testing.T) {
	input := "[Event \"Game 1\"]\n1. e4 e5\n\n[Event \"Game 2\"]\n1. d4 d5"
	games := splitGames(input)
	if len(games) != 2 {
		t.Errorf("Expected 2 games, got %d", len(games))
	}
}

func TestParseGameTags_Simple(t *testing.T) {
	input := `[Event "Test"] [White "Player A"]`
	tags := parseGameTags(input)
	if tags["Event"] != "Test" {
		t.Errorf("Expected 'Test', got '%s'", tags["Event"])
	}
}

func TestExtractCityFromSite_Basic(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Moscow, Russia", "Moscow"},
		{"Riga (Latvia)", "Riga"},
		{"Ekaterinburg", "Ekaterinburg"},
		{"?", ""},
		{"", ""},
	}
	for _, test := range tests {
		result := extractCityFromSite(test.input)
		if result != test.expected {
			t.Errorf("City from '%s': expected '%s', got '%s'", test.input, test.expected, result)
		}
	}
}

func TestExtractYearFromDate_Basic(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"2021.05.15", "2021"},
		{"????.??.??", ""},
		{"", ""},
	}
	for _, test := range tests {
		result := extractYearFromDate(test.input)
		if result != test.expected {
			t.Errorf("Year from '%s': expected '%s', got '%s'", test.input, test.expected, result)
		}
	}
}

func TestFixCommentSpacing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"well.14...Bd7", "well. 14...Bd7"},
		{"collapse.20.Qc2", "collapse. 20.Qc2"},
		{"Qh7+Kf8", "Qh7+ Kf8"},
		{"h621.Qh7+", "h6 21.Qh7+"},
		{"well.14...Bd715.Qb3!", "well. 14...Bd7 15.Qb3!"},
		{"20...h621.Qh7+ Kf822.Ne3", "20...h6 21.Qh7+ Kf8 22.Ne3"},
		{"9.O-Ob6!", "9.O-O b6!"},
		{"9.O-O-Ob6!", "9.O-O-O b6!"},
		{"b7-b5", "b7-b5"},
		{"on d2", "on d2"},
		{"0.00", "0.00"},
		{"Rxc1f4!", "Rxc1 f4!"},
		{"9.O- Ob6! Here", "9.O-O b6! Here"},
	}
	for _, tt := range tests {
		result := fixCommentSpacing(tt.input)
		if result != tt.expected {
			t.Errorf("fixCommentSpacing(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestExtractMoveText_ResultConversion(t *testing.T) {
	input := "[Result \"*\"]\n\n1. e4 e5 1-0"
	output := extractMoveText(input, false, true, false)
	if !strings.Contains(output, "1-0 (White wins)") {
		t.Error("Result 1-0 not converted")
	}
}

func TestExtractMoveText_CurlyBracesRemoved(t *testing.T) {
	input := "[Result \"*\"]\n\n1. e4 {Comment} e5\n*"
	output := extractMoveText(input, false, true, false)
	if strings.Contains(output, "{") || strings.Contains(output, "}") {
		t.Error("Curly braces should be removed")
	}
}

func TestBuildChapterTitle_MalformedPGN(t *testing.T) {
	tags := GameTags{
		"White": "Model Games",
		"Black": "Nepomniachtchi, Ian vs. Alekseenko, Kirill, Ekaterinburg 2021",
		"Site": "?",
		"Date": "????.??.??",
	}
	title := buildChapterTitle(tags, 1)
	if !strings.Contains(title, "Nepomniachtchi, Ian vs. Alekseenko, Kirill") {
		t.Errorf("Title should contain players: %s", title)
	}
	if !strings.Contains(title, "Ekaterinburg") {
		t.Errorf("Title should contain city: %s", title)
	}
	if !strings.Contains(title, "2021") {
		t.Errorf("Title should contain year: %s", title)
	}
}

func TestBuildChapterTitle_NoVsInEitherField(t *testing.T) {
	tags := GameTags{
		"White": "Exam Time!",
		"Black": "Exam Time - Mix Exercises",
	}
	title := buildChapterTitle(tags, 257)
	// Should NOT contain "vs." because neither White nor Black has it
	if strings.Contains(title, " vs. ") {
		t.Errorf("Title should NOT contain 'vs.' when neither White nor Black has it: %s", title)
	}
	// Should use comma instead
	expected := "Chapter 257: Exam Time!, Exam Time - Mix Exercises"
	if title != expected {
		t.Errorf("Title = %q, want %q", title, expected)
	}
}

func TestBuildChapterTitle_WithVsInWhite(t *testing.T) {
	tags := GameTags{
		"White": "Player A vs. Player B",
		"Black": "Some game",
	}
	title := buildChapterTitle(tags, 1)
	// Should contain "vs." because White has it
	if !strings.Contains(title, " vs. ") {
		t.Errorf("Title should contain 'vs.' when White has it: %s", title)
	}
}

func TestBuildChapterTitle_WithVsInBlack(t *testing.T) {
	tags := GameTags{
		"White": "Player A",
		"Black": "Player B vs. Player C, City 2021",
		"Site":  "?",
		"Date":  "????.??.??",
	}
	title := buildChapterTitle(tags, 1)
	// Should contain "vs." because Black has it
	if !strings.Contains(title, " vs. ") {
		t.Errorf("Title should contain 'vs.' when Black has it: %s", title)
	}
	// Should contain city from the malformed Black field
	if !strings.Contains(title, "City") {
		t.Errorf("Title should contain 'City' from Black field: %s", title)
	}
}

func TestMain_Integration(t *testing.T) {
	inputFile := "/tmp/test_input.pgn"
	outputFile := "/tmp/test_output.md"
	
	input := `[Event "Test"]
[White "Player A"]
[Black "Player B"]
[Result "*"]

1. e4 e5 *`
	os.WriteFile(inputFile, []byte(input), 0644)
	defer os.Remove(inputFile)
	defer os.Remove(outputFile)
	
	// Save original args
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	
	os.Args = []string{"build_markdown", "-src", inputFile, "-out", outputFile}
	main()
	
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Error("Output file was not created")
	}
}

func TestConvertNAG_UnknownCode(t *testing.T) {
	text := "1. e4 $999"
	result := convertNAG(text)
	if !strings.Contains(result, "(NAG $999: not recognized)") {
		t.Error("Unknown NAG code should be marked as not recognized")
	}
}

func TestConvertNAG_MultipleSameCode(t *testing.T) {
	text := "1. e4 $1 $1 $1"
	result := convertNAG(text)
	count := strings.Count(result, "±")
	if count != 3 {
		t.Errorf("Expected 3 occurrences of ±, got %d", count)
	}
}

func TestSplitGames_EmptyInput(t *testing.T) {
	input := ""
	games := splitGames(input)
	if len(games) != 0 {
		t.Errorf("Expected 0 games for empty input, got %d", len(games))
	}
}

func TestSplitGames_SingleGame(t *testing.T) {
	input := "[Event \"Single\"]\n1. e4 e5"
	games := splitGames(input)
	if len(games) != 1 {
		t.Errorf("Expected 1 game, got %d", len(games))
	}
}

func TestParseGameTags_Empty(t *testing.T) {
	input := ""
	tags := parseGameTags(input)
	if len(tags) != 0 {
		t.Error("Expected empty tags for empty input")
	}
}

func TestParseGameTags_MultipleTags(t *testing.T) {
	input := `[Event "Test"] [Site "Moscow"] [Date "2021.01.01"] [White "A"] [Black "B"] [Result "1-0"]`
	tags := parseGameTags(input)
	expected := map[string]string{
		"Event": "Test",
		"Site":  "Moscow",
		"Date":  "2021.01.01",
		"White": "A",
		"Black": "B",
		"Result": "1-0",
	}
	for key, val := range expected {
		if tags[key] != val {
			t.Errorf("Expected %s=%s, got %s", key, val, tags[key])
		}
	}
}

func TestExtractCityFromSite_WithBrackets(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Moscow (RUS)", "Moscow"},
		{"Riga [LAT]", "Riga"},
		{"St. Petersburg (Russia)", "St. Petersburg"},
	}
	for _, tt := range tests {
		result := extractCityFromSite(tt.input)
		if result != tt.expected {
			t.Errorf("extractCityFromSite(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestExtractYearFromDate_Invalid(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"abc", ""},
		{"2021", "2021"},
		{"2021-05-15", "2021"},
	}
	for _, tt := range tests {
		result := extractYearFromDate(tt.input)
		if result != tt.expected {
			t.Errorf("extractYearFromDate(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestFixCommentSpacing_EdgeCases(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"Qh7+", "Qh7+"},
		{"Kf8", "Kf8"},
		{"e4e5", "e4 e5"},
		{"Nf3Nc6", "Nf3 Nc6"},
		{"Bb5a6", "Bb5 a6"},
		{"a4a5", "a4 a5"},
		{"h3h6", "h3 h6"},
	}
	for _, tt := range tests {
		result := fixCommentSpacing(tt.input)
		if result != tt.expected {
			t.Errorf("fixCommentSpacing(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestExtractMoveText_WithFEN(t *testing.T) {
	input := `[Event "Test"]
[Result "*"]
[FEN "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"]

1. e4 e5 *`
	output := extractMoveText(input, false, true, false)
	if !strings.Contains(output, "**FEN:**") {
		t.Error("FEN should be extracted and formatted")
	}
	if !strings.Contains(output, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR") {
		t.Error("FEN value should be present")
	}
}

func TestExtractMoveText_CommentWithNAG(t *testing.T) {
	input := "[Result \"*\"]\n\n1. e4 {Good move $1} e5\n*"
	output := extractMoveText(input, false, true, false)
	if !strings.Contains(output, `1\. e4`) {
		t.Errorf("Should contain '1\\. e4' with escaped dot, got:\n%s", output)
	}
	if !strings.Contains(output, "±") {
		t.Error("NAG code in comment should be converted")
	}
	if strings.Contains(output, "$1") {
		t.Error("$1 should be converted to ±")
	}
}

func TestBuildChapterTitle_EmptyTags(t *testing.T) {
	tags := GameTags{}
	title := buildChapterTitle(tags, 1)
	if !strings.Contains(title, "Unknown player") {
		t.Error("Should use placeholder for missing players")
	}
}

func TestBuildChapterTitle_EventDate(t *testing.T) {
	tags := GameTags{
		"White": "Player A",
		"Black": "Player B",
		"EventDate": "2020.01.01",
	}
	title := buildChapterTitle(tags, 1)
	if !strings.Contains(title, "2020") {
		t.Error("Should extract year from EventDate if Date is missing")
	}
}

func TestExtractMoveText_MultipleGames(t *testing.T) {
	input := `[Event "Game 1"]
[Result "*"]

1. e4 e5 *

[Event "Game 2"]
[Result "*"]

1. d4 d5 *`
	games := splitGames(input)
	if len(games) != 2 {
		t.Fatalf("Expected 2 games, got %d", len(games))
	}
	for i, game := range games {
		output := extractMoveText(game, false, true, false)
		if !strings.Contains(output, `1\.`) {
			t.Errorf("Game %d should have move text", i+1)
		}
	}
}

func TestConvertNAG_CombinedWithMoves(t *testing.T) {
	text := "12. Bg5 $1 12...Bg4 $2"
	result := convertNAG(text)
	if !strings.Contains(result, "Bg5 ±") {
		t.Error("NAG after move should be converted")
	}
	if !strings.Contains(result, "Bg4 ∓") {
		t.Error("NAG after black move should be converted")
	}
}

func TestExtractMoveText_ResultDraw(t *testing.T) {
	input := "[Result \"*\"]\n\n1. e4 e5 1/2-1/2"
	output := extractMoveText(input, false, true, false)
	if !strings.Contains(output, "1/2-1/2 (Draw)") {
		t.Error("Result 1/2-1/2 should be converted to Draw")
	}
}

func TestExtractMoveText_ResultBlackWins(t *testing.T) {
	input := "[Result \"*\"]\n\n1. e4 e5 0-1"
	output := extractMoveText(input, false, true, false)
	if !strings.Contains(output, "0-1 (Black wins)") {
		t.Error("Result 0-1 should be converted to Black wins")
	}
}

func TestExtractMoveText_ResultOngoing(t *testing.T) {
	input := "[Result \"*\"]\n\n1. e4 e5 *"
	output := extractMoveText(input, false, true, false)
	if strings.HasSuffix(strings.TrimSpace(output), "*") {
		t.Error("Trailing * should be removed")
	}
	if strings.Contains(output, "(Game continues") {
		t.Error("Should not show 'Game continues' text for * result")
	}
}

func TestBuildChapterTitle_SpecialCharacters(t *testing.T) {
	tags := GameTags{
		"White": "René, François",
		"Black": "José, Carlos",
		"Site":  "São Paulo",
		"Date":  "2021.05.15",
	}
	title := buildChapterTitle(tags, 1)
	if !strings.Contains(title, "René, François") {
		t.Error("Special characters in names should be preserved")
	}
}

func TestFixCommentSpacing_Numbers(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"123", "123"},
		{"1.5", "1.5"},
		{"0-0", "0-0"},
		{"1-0", "1-0"},
		{"2-1", "2-1"},
	}
	for _, tt := range tests {
		result := fixCommentSpacing(tt.input)
		if result != tt.expected {
			t.Errorf("fixCommentSpacing(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestParseGameTags_WithNewlines(t *testing.T) {
	input := `[Event "Test"]
[White "Player A"]
[Black "Player B"]`
	tags := parseGameTags(input)
	if tags["Event"] != "Test" {
		t.Errorf("Expected 'Test', got '%s'", tags["Event"])
	}
}

func TestFixCommentSpacing_ChessNotation(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Nf3c5", "Nf3 c5"},
		{"Bb5a6", "Bb5 a6"},
		{"Qh4g5", "Qh4 g5"},
		{"Rae1e2", "Rae1 e2"},
		{"Bxc6b7", "Bxc6 b7"},
		{"Nbxd4c5", "Nbxd4 c5"},
	}
	for _, tt := range tests {
		result := fixCommentSpacing(tt.input)
		if result != tt.expected {
			t.Errorf("fixCommentSpacing(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestExtractMoveText_EmptyMoves(t *testing.T) {
	input := "[Result \"*\"]\n\n*"
	output := extractMoveText(input, false, true, false)
	if strings.Contains(output, "1.") {
		t.Error("Should not have move numbers with no moves")
	}
}

func TestSplitGames_WithEmptyLines(t *testing.T) {
	input := "\n\n[Event \"Game 1\"]\n1. e4 e5\n\n\n\n[Event \"Game 2\"]\n1. d4 d5\n\n"
	games := splitGames(input)
	if len(games) != 2 {
		t.Errorf("Expected 2 games with empty lines, got %d", len(games))
	}
}

func TestBuildChapterTitle_DuplicatePlayerBug(t *testing.T) {
	tags := GameTags{
		"White": "Nepomniachtchi, Ian",
		"Black": "Nepomniachtchi, Ian vs. Alekseenko, Kirill, Ekaterinburg 2021",
		"Site":  "?",
		"Date":  "????.??.??",
	}
	title := buildChapterTitle(tags, 1)
	
	// Should NOT contain duplicate "Nepomniachtchi, Ian"
	count := strings.Count(title, "Nepomniachtchi, Ian")
	if count > 1 {
		t.Errorf("Title should not contain duplicate player name. Count: %d. Title: %s", count, title)
	}
	
	// Should contain both players correctly
	if !strings.Contains(title, "Nepomniachtchi, Ian vs. Alekseenko, Kirill") {
		t.Errorf("Title should contain 'Nepomniachtchi, Ian vs. Alekseenko, Kirill': %s", title)
	}
}

func TestFixCommentSpacing_MoveAtEndOfComment(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"8...Qc7As mentioned", "8...Qc7 As mentioned"},
		{"Qc7As mentioned", "Qc7 As mentioned"},
		{"Nf3As white", "Nf3 As white"},
		{"Bb5As black", "Bb5 As black"},
		{"Rxc1As mentioned", "Rxc1 As mentioned"},
	}
	for _, tt := range tests {
		result := fixCommentSpacing(tt.input)
		if result != tt.expected {
			t.Errorf("fixCommentSpacing(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestExtractMoveText_NoDuplicateFEN(t *testing.T) {
	input := `[Event "Test"]
[Result "*"]
[FEN "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"]

1. e4 e5 *`
	output := extractMoveText(input, false, true, false)
	
	// Count occurrences of FEN
	count := strings.Count(output, "**FEN:**")
	if count > 1 {
		t.Errorf("FEN should appear only once, but found %d times:\n%s", count, output)
	}
	
	if !strings.Contains(output, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR") {
		t.Error("FEN value should be present in output")
	}
}

func TestMain_NoDuplicateFEN(t *testing.T) {
	input := `[Event "Test"]
[White "Player A"]
[Black "Player B"]
[Result "*"]
[FEN "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"]

1. e4 e5 *`
	
	// Test extractMoveText directly (which is where FEN is processed)
	output := extractMoveText(input, false, true, false)
	
	// Count occurrences of FEN
	count := strings.Count(output, "**FEN:**")
	if count > 1 {
		t.Errorf("FEN should appear only once in extractMoveText output, but found %d times", count)
	}
	
	if !strings.Contains(output, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR") {
		t.Error("FEN value should be present in output")
	}
}

func TestFixCommentSpacing_TextBeforeMove(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"after7...Bb7", "after 7...Bb7"},
		{"after7.Bb7", "after 7.Bb7"},
		{"before1.e4", "before 1.e4"},
		{"test3...Nf6", "test 3...Nf6"},
	}
	for _, tt := range tests {
		result := fixCommentSpacing(tt.input)
		if result != tt.expected {
			t.Errorf("fixCommentSpacing(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}


func TestFixCommentSpacing_MoveWithNAGAndText(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"16.Qe3±With", "16.Qe3± With"},
		{"Qe3±With", "Qe3± With"},
	}
	for _, tt := range tests {
		result := fixCommentSpacing(tt.input)
		if result != tt.expected {
			t.Errorf("fixCommentSpacing(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestFixCommentSpacing_MoveWithDotsAndText(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"15...O-OOf course", "15...O-O Of course"},
		{"O-OOf course", "O-O Of course"},
		{"O-O-OOf course", "O-O-O Of course"},
	}
	for _, tt := range tests {
		result := fixCommentSpacing(tt.input)
		if result != tt.expected {
			t.Errorf("fixCommentSpacing(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestFixCommentSpacing_MoveFollowedByText(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"17.Be4I still", "17.Be4 I still"},
		{"Be4I still", "Be4 I still"},
	}
	for _, tt := range tests {
		result := fixCommentSpacing(tt.input)
		if result != tt.expected {
			t.Errorf("fixCommentSpacing(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestFENEmptyLineInOutput(t *testing.T) {
	input := `[Event "Test"]
[Result "*"]
[FEN "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"]

1. e4 e5 *`
	output := extractMoveText(input, false, true, false)
	
	// Check that output has **FEN:** line, then empty line, then moves
	lines := strings.Split(output, "\n")
	fenLine := -1
	emptyLine := -1
	movesLine := -1
	for i, line := range lines {
		if strings.Contains(line, "**FEN:**") {
			fenLine = i
		}
		if fenLine >= 0 && i == fenLine+1 && strings.TrimSpace(line) == "" {
			emptyLine = i
		}
		if emptyLine >= 0 && i == emptyLine+1 && strings.Contains(line, `1\. e4`) {
			movesLine = i
			break
		}
	}
	if fenLine == -1 {
		t.Error("FEN line not found")
	}
	if emptyLine == -1 {
		t.Errorf("Empty line after FEN not found. Output:\n%s", output)
	}
	if movesLine == -1 {
		t.Errorf("Moves line not found after empty line. Output:\n%s", output)
	}
}

func TestFixCommentSpacing_CastlingDouble(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"6.O-OO-O since", "6.O-O O-O since"},
		{"O-OO-O since", "O-O O-O since"},
		{"O-O-OO-O since", "O-O-O O-O since"},
	}
	for _, tt := range tests {
		result := fixCommentSpacing(tt.input)
		if result != tt.expected {
			t.Errorf("fixCommentSpacing(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestFixCommentSpacing_TextWithNumberThenMove(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"worse1 6.dxe5", "worse 16.dxe5"},
		{"Following1 6...Qb8", "Following 16...Qb8"},
	}
	for _, tt := range tests {
		result := fixCommentSpacing(tt.input)
		if result != tt.expected {
			t.Errorf("fixCommentSpacing(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestFixCommentSpacing_NAGThenMove(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"17.Ne4+-16. Nb5", "17.Ne4+- 16. Nb5"},
		{"Ne4+-16. Nb5", "Ne4+- 16. Nb5"},
	}
	for _, tt := range tests {
		result := fixCommentSpacing(tt.input)
		if result != tt.expected {
			t.Errorf("fixCommentSpacing(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestReplaceFENWithImages_WrapInMarkdown(t *testing.T) {
	// Valid FEN string
	fen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	text := fen
	result := replaceFENWithImages(text, true)
	if !strings.Contains(result, "data:image/png;base64,") {
		t.Error("Expected markdown image format when wrapInMarkdown is true")
	}
	if !strings.Contains(result, "![Position]") {
		t.Error("Expected ![Position] prefix when wrapInMarkdown is true")
	}
}

func TestReplaceFENWithImages_NoWrap(t *testing.T) {
	fen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	text := fen
	result := replaceFENWithImages(text, false)
	if !strings.Contains(result, "data:image/png;base64,") {
		t.Error("Should contain data:image/png;base64, when wrapInMarkdown is false")
	}
	if strings.Contains(result, "![Position]") {
		t.Error("Should not contain ![Position] when wrapInMarkdown is false")
	}
	// Result should be a base64 string with the data URI prefix
	if result == "" {
		t.Error("Result should not be empty for valid FEN")
	}
}

func TestReplaceFENWithImages_DataURIPrefix(t *testing.T) {
	fen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	text := fen
	result := replaceFENWithImages(text, false)
	expectedPrefix := "data:image/png;base64,"
	if !strings.HasPrefix(result, expectedPrefix) {
		t.Errorf("Expected result to start with %q, got %q", expectedPrefix, result[:min(len(result), 30)])
	}
}

func TestExtractMoveText_NoMovesWithComment(t *testing.T) {
	// Chapter with only comment, no moves
	input := `[Event "Test"]
[Result "*"]

{This is a comment without any moves}`
	output := extractMoveText(input, false, true, false)
	// Should not contain "1. -- *"
	if strings.Contains(output, "1.") {
		t.Error("Should not contain move numbers when there are no moves")
	}
	// Should contain the comment without curly braces
	if !strings.Contains(output, "This is a comment without any moves") {
		t.Error("Should contain the comment text")
	}
	if strings.Contains(output, "{") || strings.Contains(output, "}") {
		t.Error("Curly braces should be removed from comments")
	}
}

func TestExtractMoveText_NoDoubleDash(t *testing.T) {
	// This pattern was producing "1. -- *"
	input := `[Event "Test"]
[Result "*"]

Some text without moves
1. Doubled Pawns
2. Isolated Pawns`
	output := extractMoveText(input, false, true, false)
	// Should not contain "--"
	if strings.Contains(output, "--") {
		t.Error("Should not contain '--'")
	}
}

func TestExtractMoveText_RemoveOneDashStar(t *testing.T) {
	// Direct test for "1. -- *" pattern
	input := `[Event "Test"]
[Result "*"]

1. -- *
Some text after`
	output := extractMoveText(input, false, true, false)
	// Should not contain the line "1. -- *"
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "1. -- *" {
			t.Error("Should not contain '1. -- *' line")
		}
	}
	// Should contain the text after
	if !strings.Contains(output, "Some text after") {
		t.Error("Should contain 'Some text after'")
		}
}

func TestExtractMoveText_NoMovesOnlyResult(t *testing.T) {
	// Chapter with only result
	input := `[Event "Test"]
[Result "*"]

*`
	output := extractMoveText(input, false, true, false)
	// Should not contain "1. -- *"
	if strings.Contains(output, "1.") {
		t.Error("Should not contain move numbers when there are no moves")
	}
}

func TestExtractMoveText_EscapedDotInMoveNumbers(t *testing.T) {
	// Test that move numbers use escaped dots to avoid Markdown list interpretation
	input := "[Result \"*\"]\n\n1. e4 e5 2. d4 d5 3. Nf3 Nc6\n*"
	output := extractMoveText(input, false, true, false)
	
	// Should contain escaped dots
	if !strings.Contains(output, "1\\. e4") {
		t.Errorf("Expected '1\\. e4' in output, got:\n%s", output)
	}
	if !strings.Contains(output, "2\\. d4") {
		t.Errorf("Expected '2\\. d4' in output, got:\n%s", output)
	}
	if !strings.Contains(output, "3\\. Nf3") {
		t.Errorf("Expected '3\\. Nf3' in output, got:\n%s", output)
	}
	
	// Should NOT contain unescaped move numbers (which would be interpreted as lists)
	// Note: We check that the pattern "N. " (where N is a digit) doesn't appear without backslash
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		// Check if line starts with a move number pattern without escape
		if len(line) > 2 && line[0] >= '0' && line[0] <= '9' {
			if !strings.Contains(line, "\\.") && strings.Contains(line, ". ") {
				t.Errorf("Found unescaped dot in move number: %s", line)
			}
		}
	}
}

func TestExtractMoveText_RemoveBracketPlaceholders(t *testing.T) {
	// Test that @@StartBracket@@ → ( and @@EndBracket@@ → )
	input := "[Result \"*\"]\n\n1. e4 @@StartBracket@@Nf3@@EndBracket@@ e5\n*"
	output := extractMoveText(input, false, true, false)
	if strings.Contains(output, "@@StartBracket@@") {
		t.Errorf("Should not contain '@@StartBracket@@', got:\n%s", output)
	}
	if strings.Contains(output, "@@EndBracket@@") {
		t.Errorf("Should not contain '@@EndBracket@@', got:\n%s", output)
	}
	if !strings.Contains(output, "(Nf3)") {
		t.Errorf("Should contain '(Nf3)', got:\n%s", output)
	}
	// Multiple bracket placeholders
	input2 := "[Result \"*\"]\n\n1. e4 @@StartBracket@@Nf3@@EndBracket@@ 2. d4 @@StartBracket@@Bb5@@EndBracket@@\n*"
	output2 := extractMoveText(input2, false, true, false)
	if strings.Contains(output2, "@@StartBracket@@") || strings.Contains(output2, "@@EndBracket@@") {
		t.Errorf("Should not contain bracket placeholders, got:\n%s", output2)
	}
	if !strings.Contains(output2, "(Nf3)") || !strings.Contains(output2, "(Bb5)") {
		t.Errorf("Should contain '(Nf3)' and '(Bb5)', got:\n%s", output2)
	}
	// Content with @ character inside
	input3 := "[Result \"*\"]\n\n1. e4 @@StartBracket@@Nf3@chessable@@EndBracket@@ e5\n*"
	output3 := extractMoveText(input3, false, true, false)
	if !strings.Contains(output3, "(Nf3@chessable)") {
		t.Errorf("Should handle content with @, got:\n%s", output3)
	}
}

func TestExtractMoveText_NoFalseMoveFromText(t *testing.T) {
	// Text with numbered list (e.g. "1. Doubled Pawns") should not be treated as moves
	input := `[Event "Test"]
[Result "*"]

Chapter 1 - Types of Weak Pawns

1. Doubled Pawns
2. Isolated Pawns
3. Backward Pawns`
	output := extractMoveText(input, false, true, false)
	// Should not contain "1. -- *" or any move-like formatting
	if strings.Contains(output, "1. -- *") {
		t.Error("Should not contain '1. -- *'")
	}
	// The numbered list items should appear as-is (maybe with some formatting)
	if !strings.Contains(output, "Doubled Pawns") {
		t.Error("Should contain 'Doubled Pawns'")
	}
}

func TestExtractMoveText_RealMoveWithComment(t *testing.T) {
	input := `[Event "Test"]
[Result "*"]

1. e4 {Good move} e5
*`
	output := extractMoveText(input, false, true, false)
	// Should contain the move with escaped dot
	if !strings.Contains(output, `1\. e4`) {
		t.Errorf("Should contain '1\\. e4' in output, got:\n%s", output)
	}
	if !strings.Contains(output, "e5") {
		t.Error("Should contain 'e5'")
	}
	// Should not contain "--"
	if strings.Contains(output, "--") {
		t.Error("Should not contain '--'")
	}
}

func TestExtractMoveText_CurlyBracesRemoved_New(t *testing.T) {
	input := `[Event "Test"]
[Result "*"]

1. e4 {Good move} e5 {Black's reply}
*`
	output := extractMoveText(input, false, true, false)
	// Should not contain curly braces
	if strings.Contains(output, "{") || strings.Contains(output, "}") {
		t.Error("Curly braces should be removed from comments")
	}
	// Should contain the move with escaped dot
	if !strings.Contains(output, `1\. e4`) {
		t.Errorf("Should contain '1\\. e4' in output, got:\n%s", output)
	}
	// Should contain the comment text
	if !strings.Contains(output, "Good move") {
		t.Error("Should contain comment text 'Good move'")
	}
	if !strings.Contains(output, "Black's reply") {
		t.Error("Should contain comment text 'Black's reply'")
	}
}

func TestReplaceFENWithImages_InvalidFEN(t *testing.T) {
	// Invalid FEN should be returned as is
	invalidFEN := "invalid_fen_string"
	text := invalidFEN
	result := replaceFENWithImages(text, true)
	if result != invalidFEN {
		t.Errorf("Invalid FEN should be returned as is, got: %s", result)
	}
	result2 := replaceFENWithImages(text, false)
	if result2 != invalidFEN {
		t.Errorf("Invalid FEN should be returned as is (no wrap), got: %s", result2)
	}
}

func TestReplaceFENWithImages_FENInText(t *testing.T) {
	// FEN embedded in text
	text := "Some text before rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1 and after"
	result := replaceFENWithImages(text, false)
	// Should not contain the markdown wrapper
	if strings.Contains(result, "![Position]") {
		t.Error("Should not contain ![Position] in text when wrapInMarkdown is false")
	}
	// The FEN part should be replaced with base64 data URI
	if strings.Contains(result, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR") {
		t.Error("FEN should be replaced with base64 string")
	}
	if !strings.Contains(result, "data:image/png;base64,") {
		t.Error("FEN should be replaced with data:image/png;base64,... format")
	}
}

func TestSplitGames_Count(t *testing.T) {
	input := `[Event "Game 1"]
1. e4 e5 *

[Event "Game 2"]
1. d4 d5 *

[Event "Game 3"]
1. c4 c5 *`
	games := splitGames(input)
	if len(games) != 3 {
		t.Errorf("Expected 3 games, got %d", len(games))
	}
}

func TestMain_ChapterLimit(t *testing.T) {
	inputFile := "/tmp/test_chapters.pgn"
	outputFile := "/tmp/test_chapters_out.md"

	input := `[Event "Game 1"]
[White "A"]
[Black "B"]
[Result "*"]
1. e4 e5 *

[Event "Game 2"]
[White "C"]
[Black "D"]
[Result "*"]
1. d4 d5 *

[Event "Game 3"]
[White "E"]
[Black "F"]
[Result "*"]
1. c4 c5 *`

	os.WriteFile(inputFile, []byte(input), 0644)
	defer os.Remove(inputFile)
	defer os.Remove(outputFile)

	binPath := "../../../bin/build_markdown"
	cmd := exec.Command(binPath, "-src", inputFile, "-out", outputFile, "-chapters", "2")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to run binary: %v", err)
	}

	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatal("Output file not created")
	}
	// Should contain only 2 chapters
	count := strings.Count(string(content), "## Chapter")
	if count != 2 {
		t.Errorf("Expected 2 chapters, got %d", count)
	}
}

func TestMain_ChapterNumbers(t *testing.T) {
	inputFile := "/tmp/test_chapter_numbers.pgn"
	outputFile := "/tmp/test_chapter_numbers_out.md"

	input := `[Event "Game 1"]
[White "A"]
[Black "B"]
[Result "*"]
1. e4 e5 *

[Event "Game 2"]
[White "C"]
[Black "D"]
[Result "*"]
1. d4 d5 *

[Event "Game 3"]
[White "E"]
[Black "F"]
[Result "*"]
1. c4 c5 *`

	os.WriteFile(inputFile, []byte(input), 0644)
	defer os.Remove(inputFile)
	defer os.Remove(outputFile)

	binPath := "../../../bin/build_markdown"
	cmd := exec.Command(binPath, "-src", inputFile, "-out", outputFile, "-chapter-numbers", "1,3")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to run binary: %v", err)
	}

	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatal("Output file not created")
	}
	// Should contain chapters 1 and 3 (skipping 2)
	count := strings.Count(string(content), "## Chapter")
	if count != 2 {
		t.Errorf("Expected 2 chapters, got %d", count)
	}
	// Check that White vs Black is present for chapters 1 and 3 (comma separator since no "vs.")
	if !strings.Contains(string(content), "A, B") {
		t.Error("Should contain 'A, B' for Game 1")
	}
	if strings.Contains(string(content), "C, D") {
		t.Error("Should not contain 'C, D' for Game 2")
	}
	if !strings.Contains(string(content), "E, F") {
		t.Error("Should contain 'E, F' for Game 3")
	}
}

func TestMain_SkipAndChapters(t *testing.T) {
	inputFile := "/tmp/test_skip_chapters.pgn"
	outputFile := "/tmp/test_skip_chapters_out.md"

	input := `[Event "Game 1"]
[White "A"]
[Black "B"]
[Result "*"]
1. e4 e5 *

[Event "Game 2"]
[White "C"]
[Black "D"]
[Result "*"]
1. d4 d5 *

[Event "Game 3"]
[White "E"]
[Black "F"]
[Result "*"]
1. c4 c5 *

[Event "Game 4"]
[White "G"]
[Black "H"]
[Result "*"]
1. Nf3 Nf6 *`

	os.WriteFile(inputFile, []byte(input), 0644)
	defer os.Remove(inputFile)
	defer os.Remove(outputFile)

	binPath := "../../../bin/build_markdown"
	cmd := exec.Command(binPath, "-src", inputFile, "-out", outputFile, "-skip", "1", "-chapters", "2")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to run binary: %v", err)
	}

	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatal("Output file not created")
	}
	// Should skip 1 (Game 1), then take 2: so chapters Game 2 and Game 3
	count := strings.Count(string(content), "## Chapter")
	if count != 2 {
		t.Errorf("Expected 2 chapters, got %d", count)
	}
	// Should contain Game 2 and Game 3 (comma separator since no "vs." in tags)
	if !strings.Contains(string(content), "C, D") {
		t.Error("Should contain 'C, D' for Game 2")
	}
	if !strings.Contains(string(content), "E, F") {
		t.Error("Should contain 'E, F' for Game 3")
	}
	if strings.Contains(string(content), "A vs. B") {
		t.Error("Should not contain 'A vs. B' (Game 1 was skipped)")
	}
	if strings.Contains(string(content), "G vs. H") {
		t.Error("Should not contain 'G vs. H' (Game 4 not in first 2 after skip)")
	}
}

// Test formatLists: numbered lists with 2 spaces
func TestFormatLists_Numbered(t *testing.T) {
	text := "Categories:   1. Doubled Pawns   2. Isolated Pawns   3. Backward Pawns"
	result := extractMoveText("{"+text+"}", false, true, false)
	// After processing, numbered items should start with newline + 2 spaces
	if !strings.Contains(result, "\n  1. Doubled Pawns") {
		t.Error("Numbered list should have newline + 2 spaces before number")
	}
	if !strings.Contains(result, "\n  2. Isolated Pawns") {
		t.Error("Numbered list should have newline + 2 spaces before number")
	}
}

// Test formatLists: dash lists with newline + 2 spaces
func TestFormatLists_Dash(t *testing.T) {
	text := "Scenarios:  - First item  - Second item"
	result := extractMoveText("{"+text+"}", false, true, false)
	// With dash-lists=false (default), dash items should NOT be formatted
	if strings.Contains(result, "\n  - First item") {
		t.Errorf("With dash-lists=false, should not format dash lists, got: %q", result)
	}
	if strings.Contains(result, "\n  - Second item") {
		t.Errorf("With dash-lists=false, should not format dash lists, got: %q", result)
	}
}

// Test formatLists: text after list item moves to next line
func TestFormatLists_TextAfterList(t *testing.T) {
	text := "Categories:   1. Item1   2. Item2   Let's continue"
	result := extractMoveText("{"+text+"}", false, true, false)
	// After processing, text after list item should be on new line
	if strings.Contains(result, "Item2   Let's") {
		t.Errorf("Text after list item should be on new line, got: %q", result)
	}
	if !strings.Contains(result, "\nLet's continue") {
		t.Errorf("Expected '\\nLet's continue' in result, got: %q", result)
	}
}

// Test chapter headers: White as ##, Black as ###
func TestChapterHeaders_WhiteAsH2_BlackAsH3(t *testing.T) {
	input := `[Event "?"]
[White "Introduction to Micro-Plans"]
[Black "Introduction to Micro-Plans: Mastering Weak Pawns"]
[Result "*"]
1. -- *

[Event "?"]
[White "Chapter 1: Recognizing Weak Pawns"]
[Black "Types of vulnerable Pawns"]
[Result "*"]
1. -- *`

	inputFile := "/tmp/test_headers.pgn"
	outputFile := "/tmp/test_headers_out.md"
	os.WriteFile(inputFile, []byte(input), 0644)
	defer os.Remove(inputFile)
	defer os.Remove(outputFile)

	binPath := "../../../bin/build_markdown"
	cmd := exec.Command(binPath, "-src", inputFile, "-out", outputFile)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to run binary: %v", err)
	}

	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatal("Output file not created")
	}
	
	// Should have ## for White
	if !strings.Contains(string(content), "## Introduction to Micro-Plans") {
		t.Error("White should be ## header")
	}
	// Should have ### for Black
	if !strings.Contains(string(content), "### Introduction to Micro-Plans: Mastering Weak Pawns") {
		t.Error("Black should be ### header")
	}
	// Next chapter: White != prevWhite, so should have ##
	if !strings.Contains(string(content), "## Chapter 1: Recognizing Weak Pawns") {
		t.Error("White should be ## header when different from prev")
	}
	// Black as ###
	if !strings.Contains(string(content), "### Types of vulnerable Pawns") {
		t.Error("Black should be ### header")
	}
}

// Test chapter headers: skip ## when White matches previous
func TestChapterHeaders_SkipWhiteWhenSame(t *testing.T) {
	input := `[Event "?"]
[White "Same Book Title"]
[Black "Chapter 1: First Chapter"]
[Result "*"]
1. -- *

[Event "?"]
[White "Same Book Title"]
[Black "Chapter 2: Second Chapter"]
[Result "*"]
1. -- *`

	inputFile := "/tmp/test_skip_white.pgn"
	outputFile := "/tmp/test_skip_white_out.md"
	os.WriteFile(inputFile, []byte(input), 0644)
	defer os.Remove(inputFile)
	defer os.Remove(outputFile)

	binPath := "../../../bin/build_markdown"
	cmd := exec.Command(binPath, "-src", inputFile, "-out", outputFile)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to run binary: %v", err)
	}

	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatal("Output file not created")
	}
	
	// Should have ## only once
	count := strings.Count(string(content), "## Same Book Title")
	if count != 1 {
		t.Errorf("## Same Book Title should appear once, got %d", count)
	}
	// Should have ### for both chapters
	if !strings.Contains(string(content), "### Chapter 1: First Chapter") {
		t.Error("First chapter Black should be ###")
	}
	if !strings.Contains(string(content), "### Chapter 2: Second Chapter") {
		t.Error("Second chapter Black should be ###")
	}
}

// Test chapter title with "vs." in Black field
func TestChapterTitle_VsInBlackField(t *testing.T) {
	input := `[Event "?"]
[Site "?"]
[Date "????.??.??"]
[Round "?"]
[White "Exam Time!"]
[Black "Sethuraman, S.P. vs. Yuffa, D."]
[Result "*"]
[FEN "4r1k1/1b2q1pp/1p1p2r1/pPp1np2/P1P5/4PP1P/1B3QP1/R3RB1K w - a6 0 26"]
[SetUp "1"]

1. e4 e5 *`
	
	inputFile := "/tmp/test_vs_black.pgn"
	outputFile := "/tmp/test_vs_black_out.md"
	os.WriteFile(inputFile, []byte(input), 0644)
	defer os.Remove(inputFile)
	defer os.Remove(outputFile)
	
	binPath := "../../../bin/build_markdown"
	cmd := exec.Command(binPath, "-src", inputFile, "-out", outputFile)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to run binary: %v", err)
	}
	
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatal("Output file not created")
	}
	
	// The Black field contains "vs.", so it should be used as ### header
	if !strings.Contains(string(content), "### Sethuraman, S.P. vs. Yuffa, D.") {
		t.Errorf("Black field with 'vs.' should be used as ### header, got:\n%s", content)
	}
	// White should be ## header
	if !strings.Contains(string(content), "## Exam Time!") {
		t.Errorf("White should be ## header, got:\n%s", content)
	}
}

// Test numbered lists with flag enabled (default)
func TestFormatLists_WithFlagEnabled(t *testing.T) {
	// This tests the actual extractMoveText with numberedLists flag
	// Since we can't easily pass the flag, we test the formatLists function indirectly
	text := "Categories:   1. Doubled Pawns   2. Isolated Pawns"
	result := extractMoveText("{"+text+"}", false, true, false)
	if !strings.Contains(result, "\n  1. Doubled Pawns") {
		t.Errorf("With flag enabled, expected newline+2spaces before number, got: %q", result)
	}
	if !strings.Contains(result, "\n  2. Isolated Pawns") {
		t.Errorf("With flag enabled, expected newline+2spaces before number, got: %q", result)
	}
}

// Test numbered lists with flag disabled
func TestFormatLists_WithFlagDisabled(t *testing.T) {
	text := "Categories:   1. Doubled Pawns   2. Isolated Pawns"
	result := extractMoveText("{"+text+"}", false, false, false)
	// When disabled, should NOT have newline+2spaces formatting
	// But dash lists should still work
	if strings.Contains(result, "\n  1. Doubled Pawns") {
		t.Errorf("With flag disabled, should not format numbered lists, got: %q", result)
	}
}

// Test -numbered-lists flag via command line
func TestMain_NumberedListsEnabled(t *testing.T) {
	inputFile := "/tmp/test_numbered_true.pgn"
	outputFile := "/tmp/test_numbered_true.md"
	
	input := `[Event "Test"]
[White "Player A"]
[Black "Player B"]
[Result "*"]

{There are 3 categories:
1. First item   2. Second item   3. Third item
Some text after.}`
	
	os.WriteFile(inputFile, []byte(input), 0644)
	defer os.Remove(inputFile)
	defer os.Remove(outputFile)
	
	binPath := "../../../bin/build_markdown"
	cmd := exec.Command(binPath, "-src", inputFile, "-out", outputFile, "-numbered-lists=true")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to run binary: %v", err)
	}
	
	output, _ := os.ReadFile(outputFile)
	outputStr := string(output)
	
	// With flag enabled, numbered list should be formatted with newlines and 2 spaces
	if !strings.Contains(outputStr, "\n  1. First item") {
		t.Errorf("With -numbered-lists=true, expected formatted numbered list, got:\n%s", outputStr)
	}
}

func TestMain_NumberedListsDisabled(t *testing.T) {
	inputFile := "/tmp/test_numbered_false.pgn"
	outputFile := "/tmp/test_numbered_false.md"
	
	input := `[Event "Test"]
[White "Player A"]
[Black "Player B"]
[Result "*"]

{There are 3 categories:
1. First item   2. Second item   3. Third item
Some text after.}`
	
	os.WriteFile(inputFile, []byte(input), 0644)
	defer os.Remove(inputFile)
	defer os.Remove(outputFile)
	
	binPath := "../../../bin/build_markdown"
	cmd := exec.Command(binPath, "-src", inputFile, "-out", outputFile, "-numbered-lists=false")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to run binary: %v", err)
	}
	
	output, _ := os.ReadFile(outputFile)
	outputStr := string(output)
	
	// With flag disabled, numbered list should NOT be formatted with newlines
	// But dash lists should still work
	if strings.Contains(outputStr, "\n  1. First item") {
		t.Errorf("With -numbered-lists=false, should not format numbered lists, got:\n%s", outputStr)
	}
}

// Test dash lists with flag enabled
func TestFormatLists_DashEnabled(t *testing.T) {
	text := "Categories:   - First item   - Second item"
	result := extractMoveText("{"+text+"}", false, false, true)
	if !strings.Contains(result, "\n  - First item") {
		t.Errorf("With dash flag enabled, expected formatted dash list, got: %q", result)
	}
}

// Test dash lists with flag disabled
func TestFormatLists_DashDisabled(t *testing.T) {
	text := "Categories:   - First item   - Second item"
	result := extractMoveText("{"+text+"}", false, false, false)
	if strings.Contains(result, "\n  - First item") {
		t.Errorf("With dash flag disabled, should not format dash lists, got: %q", result)
	}
}

// Test -dash-lists flag via command line (enabled)
func TestMain_DashListsEnabled(t *testing.T) {
	inputFile := "/tmp/test_dash_true.pgn"
	outputFile := "/tmp/test_dash_true.md"
	
	input := `[Event "Test"]
[White "Player A"]
[Black "Player B"]
[Result "*"]

{There are categories:
- First item   - Second item
Some text after.}`
	
	os.WriteFile(inputFile, []byte(input), 0644)
	defer os.Remove(inputFile)
	defer os.Remove(outputFile)
	
	binPath := "../../../bin/build_markdown"
	cmd := exec.Command(binPath, "-src", inputFile, "-out", outputFile, "-dash-lists=true")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to run binary: %v", err)
	}
	
	output, _ := os.ReadFile(outputFile)
	outputStr := string(output)
	
	if !strings.Contains(outputStr, "\n  - First item") {
		t.Errorf("With -dash-lists=true, expected formatted dash list, got:\n%s", outputStr)
	}
}

// Test -dash-lists flag via command line (disabled)
func TestMain_DashListsDisabled(t *testing.T) {
	inputFile := "/tmp/test_dash_false.pgn"
	outputFile := "/tmp/test_dash_false.md"
	
	input := `[Event "Test"]
[White "Player A"]
[Black "Player B"]
[Result "*"]

{There are categories:
- First item   - Second item
Some text after.}`
	
	os.WriteFile(inputFile, []byte(input), 0644)
	defer os.Remove(inputFile)
	defer os.Remove(outputFile)
	
	binPath := "../../../bin/build_markdown"
	cmd := exec.Command(binPath, "-src", inputFile, "-out", outputFile, "-dash-lists=false")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to run binary: %v", err)
	}
	
	output, _ := os.ReadFile(outputFile)
	outputStr := string(output)
	
	if strings.Contains(outputStr, "\n  - First item") {
		t.Errorf("With -dash-lists=false, should not format dash lists, got:\n%s", outputStr)
	}
}

// Test exam flag: filter by exact White match
func TestMain_ExamFilterByWhite(t *testing.T) {
	inputFile := "/tmp/test_exam_white.pgn"
	outputFile := "/tmp/test_exam_white_out.md"

	input := `[Event "?"]
[White "Exam Time!"]
[Black "Mix Exercises"]
[Result "*"]
1. -- *

[Event "?"]
[White "Chapter 1: Intro"]
[Black "Something"]
[Result "*"]
1. -- *

[Event "?"]
[White "Exam Time!"]
[Black "Kassis, A. vs. Kuzubov, Y. #1"]
[Result "*"]
[FEN "1rbq1rk1/p1b1n1pp/1p2pp2/2p5/3PBP2/P1P3P1/1P3N1P/R1BQR1K1 b - - 0 18"]
[SetUp "1"]
1. -- *`

	os.WriteFile(inputFile, []byte(input), 0644)
	defer os.Remove(inputFile)
	defer os.Remove(outputFile)

	binPath := "../../../bin/build_markdown"
	cmd := exec.Command(binPath, "-src", inputFile, "-out", outputFile, "-exam", "Exam Time!")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to run binary: %v", err)
	}

	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatal("Output file not created")
	}
	contentStr := string(content)

	// Chapters with White="Exam Time!" should be included
	// First chapter: formatted as "Chapter 1: Exam Time!, Mix Exercises" (no vs.)
	if !strings.Contains(contentStr, "Exam Time!, Mix Exercises") {
		t.Errorf("First chapter with White='Exam Time!' should be included. Got:\n%s", contentStr)
	}
	// Second chapter: "## Exam Time!" because Black contains "vs."
	if !strings.Contains(contentStr, "## Exam Time!") {
		t.Errorf("Second chapter with White='Exam Time!' should have ## header. Got:\n%s", contentStr)
	}

	// Chapter "Chapter 1: Intro" should NOT be included
	if strings.Contains(contentStr, "Chapter 1: Intro") {
		t.Error("Chapter with White='Chapter 1: Intro' should not be included")
	}
}

// Test exam flag: when FEN exists, output only FEN (no formatting)
func TestMain_ExamWithFEN_OutputOnlyFEN(t *testing.T) {
	inputFile := "/tmp/test_exam_fen.pgn"
	outputFile := "/tmp/test_exam_fen_out.md"

	input := `[Event "?"]
[White "Exam Time!"]
[Black "Kassis, A. vs. Kuzubov, Y. #1"]
[Result "*"]
[FEN "1rbq1rk1/p1b1n1pp/1p2pp2/2p5/3PBP2/P1P3P1/1P3N1P/R1BQR1K1 b - - 0 18"]
[SetUp "1"]

{ We have an unusual Pawn structure on the board. Can you find the best move for Black? }
18... cxd4
{ This creates an isolated Queen's Pawn for White which is a potential weakness. }
*`

	os.WriteFile(inputFile, []byte(input), 0644)
	defer os.Remove(inputFile)
	defer os.Remove(outputFile)

	binPath := "../../../bin/build_markdown"
	cmd := exec.Command(binPath, "-src", inputFile, "-out", outputFile, "-exam", "Exam Time!")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to run binary: %v", err)
	}

	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatal("Output file not created")
	}
	contentStr := string(content)

	// FEN should be present
	if !strings.Contains(contentStr, "1rbq1rk1/p1b1n1pp/1p2pp2/2p5/3PBP2/P1P3P1/1P3N1P/R1BQR1K1") {
		t.Error("FEN should be present in output")
	}

	// Comments should NOT be present (no formatting)
	if strings.Contains(contentStr, "unusual Pawn structure") {
		t.Error("Comments should NOT be present when FEN exists in exam mode")
	}
	if strings.Contains(contentStr, "isolated Queen's Pawn") {
		t.Error("Comments should NOT be present when FEN exists in exam mode")
	}

	// Move notation should NOT be present
	if strings.Contains(contentStr, "18...") {
		t.Error("Moves should NOT be present when FEN exists in exam mode")
	}
}

// Test exam flag: when no FEN, normal formatting applies
func TestMain_ExamNoFEN_AppliesFormatting(t *testing.T) {
	inputFile := "/tmp/test_exam_no_fen.pgn"
	outputFile := "/tmp/test_exam_no_fen_out.md"

	input := `[Event "?"]
[Site "?"]
[Date "????.??.??"]
[Round "?"]
[White "Exam Time!"]
[Black "Exam Time - Mix Exercises "]
[Result "*"]

{ On your journey so far in this course, you have come across different types of weak Pawns and worked on the ability to use them in different scenarios. This chapter is a mixture of all the previous chapters. Consider this as your exam. You have to apply all the knowledge and skills that you gained by working on the previous chapters of this course. Ensure that you read the questions carefully. Often you can get hints from the questions themselves. }
1. -- *`

	os.WriteFile(inputFile, []byte(input), 0644)
	defer os.Remove(inputFile)
	defer os.Remove(outputFile)

	binPath := "../../../bin/build_markdown"
	cmd := exec.Command(binPath, "-src", inputFile, "-out", outputFile, "-exam", "Exam Time!")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to run binary: %v", err)
	}

	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatal("Output file not created")
	}
	contentStr := string(content)

	// Title should be present
	if !strings.Contains(contentStr, "##") {
		t.Error("Chapter title should be present")
	}

	// Comment should be present (normal formatting)
	if !strings.Contains(contentStr, "On your journey so far") {
		t.Error("Comments should be present when no FEN in exam mode")
	}
	if !strings.Contains(contentStr, "mixture of all the previous chapters") {
		t.Error("Comments should be present when no FEN in exam mode")
	}
}

// Test exam flag: filter by Black match (exact match)
func TestMain_ExamFilterByBlack(t *testing.T) {
	inputFile := "/tmp/test_exam_black.pgn"
	outputFile := "/tmp/test_exam_black_out.md"

	input := `[Event "?"]
[White "Chapter 1"]
[Black "Exact Match Title"]
[Result "*"]
1. -- *

[Event "?"]
[White "Chapter 2"]
[Black "Different Title"]
[Result "*"]
1. -- *

[Event "?"]
[White "Chapter 3"]
[Black "Exact Match Title"]
[Result "*"]
[FEN "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"]
1. -- *`

	os.WriteFile(inputFile, []byte(input), 0644)
	defer os.Remove(inputFile)
	defer os.Remove(outputFile)

	binPath := "../../../bin/build_markdown"
	cmd := exec.Command(binPath, "-src", inputFile, "-out", outputFile, "-exam", "Exact Match Title")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to run binary: %v", err)
	}

	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatal("Output file not created")
	}
	contentStr := string(content)

	// Both chapters with Black="Exact Match Title" should be included
	// First chapter: Chapter 1 vs. Exact Match Title
	if !strings.Contains(contentStr, "### Exact Match Title") {
		t.Errorf("First chapter with Black='Exact Match Title' should be included. Got:\n%s", contentStr)
	}
	// Third chapter: Chapter 3 with FEN only
	if !strings.Contains(contentStr, "Chapter 3") {
		t.Errorf("Third chapter with Black='Exact Match Title' should be included. Got:\n%s", contentStr)
	}
	// Should contain FEN from third chapter (exam mode with FEN outputs only FEN)
	if !strings.Contains(contentStr, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR") {
		t.Error("FEN should be present from third chapter")
	}

	// Chapter "Different Title" should NOT be included
	if strings.Contains(contentStr, "Chapter 2") {
		t.Error("Chapter with Black='Different Title' should not be included")
	}
}

// Test exam flag with inline-images: FEN converted to base64, no **FEN:** prefix
func TestMain_ExamWithFENAndInlineImages(t *testing.T) {
	inputFile := "/tmp/test_exam_inline.pgn"
	outputFile := "/tmp/test_exam_inline_out.md"

	input := `[Event "?"]
[White "Exam Time!"]
[Black "Test"]
[Result "*"]
[FEN "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"]
[SetUp "1"]
1. -- *`

	os.WriteFile(inputFile, []byte(input), 0644)
	defer os.Remove(inputFile)
	defer os.Remove(outputFile)

	binPath := "../../../bin/build_markdown"
	cmd := exec.Command(binPath, "-src", inputFile, "-out", outputFile, "-exam", "Exam Time!", "-inline-images")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to run binary: %v", err)
	}

	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatal("Output file not created")
	}
	contentStr := string(content)

	// Should contain base64 data URI format
	if !strings.Contains(contentStr, "data:image/png;base64,") {
		t.Error("FEN should be converted to base64 image when -inline-images is set")
	}

	// Should NOT contain raw FEN string
	if strings.Contains(contentStr, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR") {
		t.Error("Raw FEN string should not be present when -inline-images is set")
	}

	// Should NOT contain **FEN:** prefix when inline-images is used
	if strings.Contains(contentStr, "**FEN:**") {
		t.Error("**FEN:** prefix should not be present when -inline-images is set in exam mode")
	}
}

// Test exam flag: no chapters match - empty output
func TestMain_ExamNoMatchingChapters(t *testing.T) {
	inputFile := "/tmp/test_exam_no_match.pgn"
	outputFile := "/tmp/test_exam_no_match_out.md"

	input := `[Event "?"]
[White "Chapter 1"]
[Black "A"]
[Result "*"]
1. -- *

[Event "?"]
[White "Chapter 2"]
[Black "B"]
[Result "*"]
1. -- *`

	os.WriteFile(inputFile, []byte(input), 0644)
	defer os.Remove(inputFile)
	defer os.Remove(outputFile)

	binPath := "../../../bin/build_markdown"
	cmd := exec.Command(binPath, "-src", inputFile, "-out", outputFile, "-exam", "Non-existent Title")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to run binary: %v", err)
	}

	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatal("Output file not created")
	}
	contentStr := strings.TrimSpace(string(content))

	// Output should be empty or contain only whitespace
	if contentStr != "" {
		t.Errorf("Expected empty output when no chapters match, got: %q", contentStr)
	}
}

// Test that exam flag does not affect chapters without matching White or Black
func TestMain_ExamPreservesOtherChapters(t *testing.T) {
	inputFile := "/tmp/test_exam_mixed.pgn"
	outputFile := "/tmp/test_exam_mixed_out.md"

	input := `[Event "?"]
[White "Other Chapter"]
[Black "X"]
[Result "*"]
1. e4 e5 *

[Event "?"]
[White "Exam Time!"]
[Black "Y"]
[Result "*"]
1. d4 d5 *

[Event "?"]
[White "Another Chapter"]
[Black "Z"]
[Result "*"]
1. c4 c5 *`

	os.WriteFile(inputFile, []byte(input), 0644)
	defer os.Remove(inputFile)
	defer os.Remove(outputFile)

	binPath := "../../../bin/build_markdown"
	cmd := exec.Command(binPath, "-src", inputFile, "-out", outputFile, "-exam", "Exam Time!")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to run binary: %v", err)
	}

	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatal("Output file not created")
	}
	contentStr := string(content)

	// Only chapter with White="Exam Time!" should be present
	if !strings.Contains(contentStr, "Exam Time!") {
		t.Error("Chapter with White='Exam Time!' should be included")
	}

	// Other chapters should NOT be present
	if strings.Contains(contentStr, "Other Chapter") {
		t.Error("Chapter with White='Other Chapter' should not be included")
	}
	if strings.Contains(contentStr, "Another Chapter") {
		t.Error("Chapter with White='Another Chapter' should not be included")
	}
}

// Test exam mode without inline-images: no **FEN:** prefix
func TestMain_ExamWithFEN_NoInlineImages_NoPrefix(t *testing.T) {
	inputFile := "/tmp/test_exam_no_inline.pgn"
	outputFile := "/tmp/test_exam_no_inline_out.md"

	input := `[Event "?"]
[White "Exam Time!"]
[Black "Test"]
[Result "*"]
[FEN "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"]
[SetUp "1"]
1. -- *`

	os.WriteFile(inputFile, []byte(input), 0644)
	defer os.Remove(inputFile)
	defer os.Remove(outputFile)

	binPath := "../../../bin/build_markdown"
	cmd := exec.Command(binPath, "-src", inputFile, "-out", outputFile, "-exam", "Exam Time!")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to run binary: %v", err)
	}

	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatal("Output file not created")
	}
	contentStr := string(content)

	// In exam mode without inline-images, FEN should be present but WITHOUT **FEN:** prefix
	if !strings.Contains(contentStr, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR") {
		t.Error("FEN should be present in output")
	}

	// Should NOT contain **FEN:** prefix in exam mode
	if strings.Contains(contentStr, "**FEN:**") {
		t.Error("**FEN:** prefix should not be present in exam mode (FEN should be output directly)")
	}
}

// Test inline FEN markers: basic handling
func TestExtractMoveText_InlineFENMarkers(t *testing.T) {
	input := "[Result \"*\"]\n\n1. e4 {Some text @@StartFEN@@4k3/8/8/8/8/8/8/4K3 w - - 0 1@@EndFEN@@ more text} e5\n*"
	output := extractMoveText(input, false, true, false)

	// Markers should be removed
	if strings.Contains(output, "@@StartFEN@@") {
		t.Error("@@StartFEN@@ marker should be removed")
	}
	if strings.Contains(output, "@@EndFEN@@") {
		t.Error("@@EndFEN@@ marker should be removed")
	}

	// FEN should be formatted with **FEN:** label
	if !strings.Contains(output, "**FEN:**") {
		t.Error("FEN should be formatted with **FEN:** label")
	}
	if !strings.Contains(output, "4k3/8/8/8/8/8/8/4K3") {
		t.Error("FEN value should be present")
	}

	// Comment text should still be present
	if !strings.Contains(output, "Some text") {
		t.Error("Comment text before FEN should be present")
	}
	if !strings.Contains(output, "more text") {
		t.Error("Comment text after FEN should be present")
	}
}

// Test inline FEN markers with inline-images flag
func TestExtractMoveText_InlineFENMarkersWithImages(t *testing.T) {
	input := "[Result \"*\"]\n\n1. e4 {Text @@StartFEN@@4k3/8/8/8/8/8/8/4K3 w - - 0 1@@EndFEN@@ end} e5\n*"
	output := extractMoveText(input, true, true, false)

	// Markers should be removed
	if strings.Contains(output, "@@StartFEN@@") {
		t.Error("@@StartFEN@@ marker should be removed")
	}
	if strings.Contains(output, "@@EndFEN@@") {
		t.Error("@@EndFEN@@ marker should be removed")
	}

	// FEN should be present WITHOUT **FEN:** label (inlineImages mode)
	if strings.Contains(output, "**FEN:**") {
		t.Error("FEN should NOT have **FEN:** label in inlineImages mode")
	}
	if !strings.Contains(output, "4k3/8/8/8/8/8/8/4K3") {
		t.Error("FEN value should be present")
	}
}

// Test inline FEN markers: starting position FEN should be skipped
func TestExtractMoveText_InlineFEN_StartingPositionSkipped(t *testing.T) {
	input := "[Result \"*\"]\n\n1. e4 {Start @@StartFEN@@rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1@@EndFEN@@ end} e5\n*"
	output := extractMoveText(input, false, true, false)

	// Starting position FEN should NOT appear
	if strings.Contains(output, "**FEN:**") {
		t.Error("Starting position FEN should not be output")
	}
	if strings.Contains(output, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR") {
		t.Error("Starting position FEN string should not appear")
	}

	// Markers should be removed
	if strings.Contains(output, "@@StartFEN@@") {
		t.Error("@@StartFEN@@ marker should be removed")
	}
	if strings.Contains(output, "@@EndFEN@@") {
		t.Error("@@EndFEN@@ marker should be removed")
	}

	// Comment text should still be present
	if !strings.Contains(output, "Start") {
		t.Error("Comment text before FEN should be present")
	}
	if !strings.Contains(output, "end") {
		t.Error("Comment text after FEN should be present")
	}
}

// Test inline FEN markers: starting position with inline-images
func TestExtractMoveText_InlineFEN_StartingPositionWithImages(t *testing.T) {
	input := "[Result \"*\"]\n\n1. e4 {Start @@StartFEN@@rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1@@EndFEN@@ end} e5\n*"
	output := extractMoveText(input, true, true, false)

	// Starting position FEN should NOT appear (even in inline-images mode)
	if strings.Contains(output, "**FEN:**") {
		t.Error("Starting position FEN should not have **FEN:** label")
	}
	if strings.Contains(output, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR") {
		t.Error("Starting position FEN string should not appear in inline-images mode")
	}
}

// Test inline FEN markers: empty FEN content
func TestExtractMoveText_InlineFENMarkers_Empty(t *testing.T) {
	input := "[Result \"*\"]\n\n1. e4 {Text @@StartFEN@@@@EndFEN@@ end} e5\n*"
	output := extractMoveText(input, false, true, false)

	if strings.Contains(output, "@@StartFEN@@") {
		t.Error("Empty @StartFEN@@...@@EndFEN@@ should be removed")
	}
}

// Test inline FEN markers with full example from real PGN
func TestExtractMoveText_InlineFENMarkers_FullExample(t *testing.T) {
	input := `[Event "?"]
[Site "?"]
[Date "????.??.??"]
[Round "?"]
[White "1. Noteboom 4.Nc3 dxc4 − 5th Moves"]
[Black "Noteboom 4.Nc3 dxc4 5.Bg5 #1"]
[Result "*"]

1. d4 d5 2. c4 e6 3. Nf3 c6 4. Nc3 dxc4
{ The Noteboom Variation text. }
5. Bg5
{ Visual description. }
5... Be7
{ Materialistic note. }
6. Bxe7 Nxe7
{ Simplest. White is probably better off playing the more modest @@StartFEN@@rnb1k1nr/pp2qppp/2p1p3/8/2pP4/2N2N2/PP2PPPP/R2QKB1R w KQkq - 0 7@@EndFEN@@  7.e3 b5 8.a4 , which does seem to give good compensation. }
7. e4 { Gambit-style play. } 7... b5 8. a4 Bb7
{ Black holds a pleasant advantage. } 9. axb5 cxb5 10. Nxb5 Bxe4 11. Bxc4 O-O { White is left with an ugly isolated pawn. } *`

	output := extractMoveText(input, false, true, false)

	// Markers should be removed
	if strings.Contains(output, "@@StartFEN@@") {
		t.Error("@@StartFEN@@ marker should be removed")
	}
	if strings.Contains(output, "@@EndFEN@@") {
		t.Error("@@EndFEN@@ marker should be removed")
	}

	// FEN should be formatted with **FEN:** label
	if !strings.Contains(output, "**FEN:**") {
		t.Errorf("FEN should be formatted with **FEN:** label.\nOutput:\n%s", output)
	}
	if !strings.Contains(output, "rnb1k1nr/pp2qppp/2p1p3/8/2pP4/2N2N2/PP2PPPP/R2QKB1R") {
		t.Error("FEN value should be present")
	}

	// FEN should be on its own line
	lines := strings.Split(output, "\n")
	fenLineFound := false
	for _, line := range lines {
		if strings.Contains(line, "**FEN:**") {
			fenLineFound = true
		}
	}
	if !fenLineFound {
		t.Error("**FEN:** line not found in output")
	}

	// Comment text should be preserved
	if !strings.Contains(output, "Simplest.") {
		t.Error("Comment text before FEN should be preserved")
	}
}

// Test inline FEN markers: no conflict with regular FEN tag
func TestExtractMoveText_InlineFENWithRegularFEN(t *testing.T) {
	input := `[Event "Test"]
[Result "*"]
[FEN "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"]

1. e4 {Comment with @@StartFEN@@4k3/8/8/8/8/8/8/4K3 w - - 0 1@@EndFEN@@ inline} e5 *`

	output := extractMoveText(input, false, true, false)

	// Both FENs should be present
	// The regular FEN (from tag) should appear as **FEN:** at the start
	if !strings.Contains(output, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR") {
		t.Error("Regular FEN from tag should be present")
	}

	// The inline FEN should also be formatted
	if !strings.Contains(output, "4k3/8/8/8/8/8/8/4K3") {
		t.Error("Inline FEN should be present")
	}

	// Both should have **FEN:** prefix
	count := strings.Count(output, "**FEN:**")
	if count != 2 {
		t.Errorf("Expected 2 **FEN:** occurrences (one for tag, one for inline), got %d", count)
	}
}

// Test inline-images with inline FEN: image replacement works
func TestMain_InlineImagesWithInlineFEN(t *testing.T) {
	inputFile := "/tmp/test_inline_fen_img.pgn"
	outputFile := "/tmp/test_inline_fen_img_out.md"

	input := `[Event "?"]
[White "Player A"]
[Black "Player B"]
[Result "*"]

1. e4 {Start @@StartFEN@@4k3/8/8/8/8/8/8/4K3 w - - 0 1@@EndFEN@@ end} e5 *`

	os.WriteFile(inputFile, []byte(input), 0644)
	defer os.Remove(inputFile)
	defer os.Remove(outputFile)

	binPath := "../../../bin/build_markdown"
	cmd := exec.Command(binPath, "-src", inputFile, "-out", outputFile, "-inline-images")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to run binary: %v", err)
	}

	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatal("Output file not created")
	}
	contentStr := string(content)

	// Markers should be removed
	if strings.Contains(contentStr, "@@StartFEN@@") {
		t.Error("@@StartFEN@@ marker should be removed")
	}
	if strings.Contains(contentStr, "@@EndFEN@@") {
		t.Error("@@EndFEN@@ marker should be removed")
	}

	// Should contain base64 data URI format (from replaceFENWithImages)
	if !strings.Contains(contentStr, "data:image/png;base64,") {
		t.Error("FEN should be converted to base64 image when -inline-images is set")
	}

	// Should NOT contain raw FEN string
	if strings.Contains(contentStr, "4k3/8/8/8/8/8/8/4K3") {
		t.Error("Raw FEN string should not be present when -inline-images is set")
	}

	// Should NOT contain **FEN:** prefix
	if strings.Contains(contentStr, "**FEN:**") {
		t.Error("**FEN:** prefix should not be present with -inline-images")
	}
}

// Test inline-images without exam mode: no **FEN:** prefix
func TestMain_InlineImagesNoExam_NoPrefix(t *testing.T) {
	inputFile := "/tmp/test_inline_no_exam.pgn"
	outputFile := "/tmp/test_inline_no_exam_out.md"

	input := `[Event "?"]
[White "Player A"]
[Black "Player B"]
[Result "*"]
[FEN "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"]
[SetUp "1"]
1. e4 e5 *`

	os.WriteFile(inputFile, []byte(input), 0644)
	defer os.Remove(inputFile)
	defer os.Remove(outputFile)

	binPath := "../../../bin/build_markdown"
	cmd := exec.Command(binPath, "-src", inputFile, "-out", outputFile, "-inline-images")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to run binary: %v", err)
	}

	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatal("Output file not created")
	}
	contentStr := string(content)

	// Should contain base64 data URI format
	if !strings.Contains(contentStr, "data:image/png;base64,") {
		t.Error("FEN should be converted to base64 image when -inline-images is set")
	}

	// Should NOT contain raw FEN string
	if strings.Contains(contentStr, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR") {
		t.Error("Raw FEN string should not be present when -inline-images is set")
	}

	// Should NOT contain **FEN:** prefix when inline-images is used (even without exam)
	if strings.Contains(contentStr, "**FEN:**") {
		t.Error("**FEN:** prefix should not be present when -inline-images is set")
	}
}

// Test trailing * result marker removed from end of moves
func TestExtractMoveText_TrailingStarRemoved(t *testing.T) {
	input := "[Result \"*\"]\n\n1. e4 e5 *"
	output := extractMoveText(input, false, true, false)
	if strings.HasSuffix(strings.TrimSpace(output), "*") {
		t.Errorf("Trailing * should be removed, got:\n%s", output)
	}
}

// Test trailing * on its own line removed
func TestExtractMoveText_TrailingStarOwnLine(t *testing.T) {
	input := "[Result \"*\"]\n\n1. e4 e5\n*"
	output := extractMoveText(input, false, true, false)
	if strings.HasSuffix(strings.TrimSpace(output), "*") {
		t.Errorf("Trailing * on own line should be removed, got:\n%s", output)
	}
}

// Test trailing * with whitespace removed
func TestExtractMoveText_TrailingStarWithWhitespace(t *testing.T) {
	input := "[Result \"*\"]\n\n1. e4 e5 *   "
	output := extractMoveText(input, false, true, false)
	if strings.HasSuffix(strings.TrimSpace(output), "*") {
		t.Errorf("Trailing * with whitespace should be removed, got:\n%s", output)
	}
}

// Test * in middle of text preserved
func TestExtractMoveText_StarInMiddlePreserved(t *testing.T) {
	input := "[Result \"*\"]\n\n1. e4 {Some * text here} e5 *"
	output := extractMoveText(input, false, true, false)
	if !strings.Contains(output, "*") {
		t.Error("Star in middle of comment should be preserved")
	}
}

// Test no change when there's no trailing *
func TestExtractMoveText_NoTrailingStar(t *testing.T) {
	input := "[Result \"1-0\"]\n\n1. e4 e5 1-0"
	output := extractMoveText(input, false, true, false)
	if !strings.Contains(output, "1-0 (White wins)") {
		t.Error("Result 1-0 should be preserved")
	}
}

// Test trailing * removed with comments present
func TestExtractMoveText_TrailingStarWithComments(t *testing.T) {
	input := "[Result \"*\"]\n\n1. e4 {Good} e5 {Reply} *"
	output := extractMoveText(input, false, true, false)
	if strings.HasSuffix(strings.TrimSpace(output), "*") {
		t.Errorf("Trailing * with comments should be removed, got:\n%s", output)
	}
	if !strings.Contains(output, "Good") {
		t.Error("Comment 'Good' should be preserved")
	}
	if !strings.Contains(output, "Reply") {
		t.Error("Comment 'Reply' should be preserved")
	}
}
