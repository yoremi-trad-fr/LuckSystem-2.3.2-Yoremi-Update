package siglusluca

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

type Options struct {
	LucaDir      string
	SiglusDir    string
	OutputDir    string
	HDOutput     string
	ReviewOutput string
	TargetCol    int
	MinScore     float64
}

type Summary struct {
	FilesProcessed int
	FilesCopied    int
	FilesSkipped   int
	Imported       int
	HDCandidates   int
	ReviewRows     int
	LowConfidence  int
}

type fileSummary struct {
	file          string
	imported      int
	hdCandidates  int
	reviewRows    int
	lowConfidence int
}

type siglusEntry struct {
	id     string
	source string
	target string
	prep   preparedText
}

type lucaEntry struct {
	seq       int
	fileLine  int
	id        string
	tag       string
	text      string
	quoteIdx  int
	prep      preparedText
	replaced  bool
	review    bool
	reviewWhy string
}

type operation struct {
	kind  string
	sig   int
	luca  int
	score float64
}

type preparedText struct {
	norm   string
	tokens map[string]int
	mag    float64
}

type reportRow struct {
	file       string
	kind       string
	lucaID     string
	line       int
	tag        string
	score      float64
	coverage   float64
	lucaText   string
	siglusID   string
	siglusText string
	siglusFR   string
}

var (
	siglusPairRE = regexp.MustCompile(`^([○●])([0-9]{10})[○●](.*)$`)
	assetRE      = regexp.MustCompile(`(?i)^(?:se|bgm|bg|ev|cg|ef|fg|si|tp|ch|st)[A-Za-z0-9_\-]*(?:\([^)]*\))?$`)
	numberedIDRE = regexp.MustCompile(`^[A-Za-z_]+\d+[A-Za-z0-9_]*$`)
	trophyRE     = regexp.MustCompile(`(?i)^Harmonia_trophy\d+$`)
	wordRE       = regexp.MustCompile(`[a-z0-9]+(?:'[a-z]+)?`)
)

var stopWords = map[string]bool{
	"a": true, "an": true, "the": true, "and": true, "or": true, "but": true, "if": true, "then": true,
	"of": true, "to": true, "in": true, "on": true, "at": true, "for": true, "from": true, "with": true,
	"without": true, "as": true, "by": true, "into": true, "onto": true, "is": true, "are": true,
	"was": true, "were": true, "be": true, "been": true, "being": true, "i": true, "me": true,
	"my": true, "myself": true, "you": true, "your": true, "he": true, "she": true, "it": true,
	"its": true, "they": true, "them": true, "their": true, "we": true, "our": true, "this": true,
	"that": true, "these": true, "those": true, "there": true, "here": true, "not": true, "no": true,
	"yes": true, "do": true, "did": true, "does": true, "done": true, "had": true, "have": true,
	"has": true, "would": true, "could": true, "should": true, "can": true, "may": true, "might": true,
	"must": true, "will": true, "just": true, "even": true, "so": true, "all": true, "only": true,
	"own": true, "more": true, "most": true, "much": true, "many": true, "some": true, "any": true,
	"one": true, "two": true, "after": true, "before": true, "when": true, "while": true, "where": true,
	"what": true, "who": true, "why": true, "how": true, "than": true, "through": true, "out": true,
	"up": true, "down": true, "over": true, "under": true, "again": true, "still": true, "also": true,
	"really": true, "very": true, "perhaps": true, "maybe": true, "probably": true, "because": true,
	"since": true, "though": true, "although": true, "going": true, "went": true, "got": true, "get": true,
	"make": true, "made": true, "take": true, "took": true, "give": true, "gave": true, "find": true,
	"found": true, "look": true, "looked": true, "see": true, "saw": true, "felt": true, "feel": true,
	"think": true, "thought": true, "know": true, "knew": true, "want": true, "wanted": true,
}

var tokenAliases = map[string]string{
	"alright":     "okay",
	"beautiful":   "pretty",
	"cleaned":     "clean",
	"cleaning":    "clean",
	"cooking":     "cook",
	"entertained": "serve",
	"food":        "cook",
	"good":        "nice",
	"helpful":     "help",
	"helping":     "help",
	"helped":      "help",
	"making":      "make",
	"pointed":     "point",
	"productive":  "useful",
	"sounded":     "sound",
	"tasty":       "delicious",
	"tidy":        "clean",
}

var speakerOrControl = map[string]bool{
	"シオナ": true, "レイ": true, "マッド": true, "ティピィ": true, "青年": true, "女": true, "男": true,
	"少女": true, "少年": true, "母親": true, "父親": true, "子供": true, "女の子": true, "男の子": true,
	"店員": true, "ワタライ": true, "京子": true,
	"Madd": true, "Shiona": true, "Rei": true, "Tipi": true, "Watarai": true, "Kyoko": true,
}

func Run(opts Options) (*Summary, error) {
	if opts.LucaDir == "" || opts.SiglusDir == "" || opts.OutputDir == "" {
		return nil, fmt.Errorf("luca, siglus and output directories are required")
	}
	if opts.TargetCol <= 0 {
		opts.TargetCol = 2
	}
	if opts.HDOutput == "" {
		opts.HDOutput = filepath.Join(opts.OutputDir, "hd_candidates.tsv")
	}
	if opts.ReviewOutput == "" {
		opts.ReviewOutput = filepath.Join(opts.OutputDir, "review.tsv")
	}
	if err := os.MkdirAll(opts.OutputDir, 0755); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(filepath.Dir(opts.HDOutput), 0755); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(filepath.Dir(opts.ReviewOutput), 0755); err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(opts.LucaDir)
	if err != nil {
		return nil, err
	}

	var hdRows []reportRow
	var reviewRows []reportRow
	summary := &Summary{}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(strings.ToLower(entry.Name()), ".txt") {
			continue
		}
		if strings.HasSuffix(strings.ToLower(entry.Name()), ".ext.txt") {
			continue
		}

		src := filepath.Join(opts.LucaDir, entry.Name())
		dst := filepath.Join(opts.OutputDir, entry.Name())
		siglusPath := filepath.Join(opts.SiglusDir, strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))+".ss.txt")
		if _, err := os.Stat(siglusPath); err != nil {
			if err := copyFile(src, dst); err != nil {
				return nil, err
			}
			summary.FilesCopied++
			continue
		}

		fs, hd, review, err := processFile(src, siglusPath, dst, opts)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", entry.Name(), err)
		}
		summary.FilesProcessed++
		summary.Imported += fs.imported
		summary.HDCandidates += fs.hdCandidates
		summary.ReviewRows += fs.reviewRows
		summary.LowConfidence += fs.lowConfidence
		hdRows = append(hdRows, hd...)
		reviewRows = append(reviewRows, review...)
	}

	if err := writeReport(opts.HDOutput, hdRows); err != nil {
		return nil, err
	}
	if err := writeReport(opts.ReviewOutput, reviewRows); err != nil {
		return nil, err
	}
	return summary, nil
}

func processFile(lucaPath, siglusPath, outputPath string, opts Options) (fileSummary, []reportRow, []reportRow, error) {
	var result fileSummary
	result.file = filepath.Base(lucaPath)

	lucaData, err := os.ReadFile(lucaPath)
	if err != nil {
		return result, nil, nil, err
	}
	lines := strings.Split(string(lucaData), "\n")
	lucaEntries := extractLucaEntries(lines, opts.TargetCol)
	siglusEntries, err := loadSiglusEntries(siglusPath)
	if err != nil {
		return result, nil, nil, err
	}
	if len(lucaEntries) == 0 || len(siglusEntries) == 0 {
		if err := os.WriteFile(outputPath, lucaData, 0644); err != nil {
			return result, nil, nil, err
		}
		return result, nil, nil, nil
	}

	ops := align(siglusEntries, lucaEntries)
	reviewOps := map[int]string{}
	var hdRows []reportRow
	var reviewRows []reportRow

	markMergedMatchesForReview(ops, siglusEntries, lucaEntries, reviewOps)

	for _, group := range groupedOps(ops, "skipL") {
		lucaIdxs := lucaIndexesForOps(ops, group)
		if len(lucaIdxs) == 0 {
			continue
		}
		var contentIdxs []int
		for _, idx := range lucaIdxs {
			if isMinorLucaOnlyText(lucaEntries[idx].text) {
				row := makeReportRow(result.file, "luca_only_short_or_punctuation", lucaEntries[idx], 0, 0, siglusEntry{})
				reviewRows = append(reviewRows, row)
				result.reviewRows++
				continue
			}
			contentIdxs = append(contentIdxs, idx)
		}
		if len(contentIdxs) == 0 {
			continue
		}

		if skippedSig, skippedCoverage := bestNearbySkippedSiglusCoverage(ops, group, siglusEntries, lucaEntries); skippedCoverage >= 0.20 {
			for _, idx := range contentIdxs {
				row := makeReportRow(result.file, "siglus_luca_split_or_rewrite", lucaEntries[idx], 0, skippedCoverage, siglusEntryAt(siglusEntries, skippedSig))
				reviewRows = append(reviewRows, row)
				result.reviewRows++
			}
			continue
		}

		bestSig, coverage := bestNearbySiglusCoverage(ops, group, siglusEntries, lucaEntries)
		if coverage >= 0.20 && bestSig >= 0 {
			for _, opIndex := range group {
				reviewOps[opIndex] = "siglus_merged_or_luca_split"
			}
			if matchedOpIndex := findMatchedOpBySiglus(ops, bestSig, group[0], group[len(group)-1]); matchedOpIndex >= 0 {
				reviewOps[matchedOpIndex] = "siglus_merged_or_luca_split"
			}
			for _, idx := range contentIdxs {
				row := makeReportRow(result.file, "siglus_merged_or_luca_split", lucaEntries[idx], 0, coverage, siglusEntryAt(siglusEntries, bestSig))
				reviewRows = append(reviewRows, row)
				result.reviewRows++
			}
			continue
		}
		for _, idx := range contentIdxs {
			row := makeReportRow(result.file, "hd_candidate", lucaEntries[idx], 0, coverage, siglusEntryAt(siglusEntries, bestSig))
			hdRows = append(hdRows, row)
			result.hdCandidates++
		}
	}

	for opIndex, op := range ops {
		switch op.kind {
		case "match":
			luca := &lucaEntries[op.luca]
			siglus := siglusEntries[op.sig]
			if reason, ok := reviewOps[opIndex]; ok {
				luca.review = true
				luca.reviewWhy = reason
				reviewRows = append(reviewRows, makeReportRow(result.file, reason, *luca, op.score, 1, siglus))
				result.reviewRows++
				continue
			}
			if op.score < opts.MinScore {
				reviewRows = append(reviewRows, makeReportRow(result.file, "below_min_score_not_imported", *luca, op.score, 0, siglus))
				result.reviewRows++
				result.lowConfidence++
				continue
			}
			if op.score < 0.25 {
				reviewRows = append(reviewRows, makeReportRow(result.file, "low_confidence_imported", *luca, op.score, 0, siglus))
				result.reviewRows++
				result.lowConfidence++
			}
			lines[luca.fileLine] = replaceNthQuotedString(lines[luca.fileLine], luca.quoteIdx, siglus.target)
			luca.replaced = true
			result.imported++
		case "skipS":
			siglus := siglusEntries[op.sig]
			reviewRows = append(reviewRows, reportRow{
				file:       result.file,
				kind:       "siglus_only_or_removed_in_hd",
				score:      0,
				siglusID:   siglus.id,
				siglusText: siglus.source,
				siglusFR:   siglus.target,
			})
			result.reviewRows++
		}
	}

	output := strings.Join(lines, "\n")
	if err := os.WriteFile(outputPath, []byte(output), 0644); err != nil {
		return result, nil, nil, err
	}
	return result, hdRows, reviewRows, nil
}

func loadSiglusEntries(path string) ([]siglusEntry, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	type pair struct {
		source string
		target string
	}
	order := []string{}
	pairs := map[string]*pair{}
	scanner := bufio.NewScanner(file)
	buf := make([]byte, 0, 1024*1024)
	scanner.Buffer(buf, 16*1024*1024)
	for scanner.Scan() {
		line := scanner.Text()
		m := siglusPairRE.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		mark, id, text := m[1], m[2], m[3]
		if pairs[id] == nil {
			pairs[id] = &pair{}
			order = append(order, id)
		}
		if mark == "○" {
			pairs[id].source = text
		} else {
			pairs[id].target = text
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	entries := []siglusEntry{}
	for _, id := range order {
		p := pairs[id]
		if p == nil || p.source == "" || p.target == "" {
			continue
		}
		if isControlText(p.source) || isControlText(p.target) {
			continue
		}
		if !containsLatinLetter(p.target) {
			continue
		}
		if !containsTextLetter(p.source) && !containsTextLetter(p.target) {
			continue
		}
		entry := siglusEntry{
			id:     id,
			source: p.source,
			target: p.target,
			prep:   prepareText(p.source),
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

func extractLucaEntries(lines []string, targetCol int) []lucaEntry {
	entries := []lucaEntry{}
	quoteIdx := targetCol - 1
	seq := 0
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		tag, ok := dialogueOpcode(trimmed)
		if !ok {
			continue
		}
		quoted := extractQuotedStrings(trimmed)
		if quoteIdx < 0 || quoteIdx >= len(quoted) {
			continue
		}
		text := quoted[quoteIdx]
		if normalizeText(text) == "" || isControlText(text) {
			continue
		}
		seq++
		entries = append(entries, lucaEntry{
			seq:      seq,
			fileLine: i,
			id:       strconv.Itoa(seq),
			tag:      tag,
			text:     text,
			quoteIdx: quoteIdx,
			prep:     prepareText(text),
		})
	}
	return entries
}

func align(siglus []siglusEntry, luca []lucaEntry) []operation {
	n, m := len(siglus), len(luca)
	gapPenalty := -0.18
	matchBase := 0.20
	dp := make([][]float64, n+1)
	bt := make([][]operation, n+1)
	for i := 0; i <= n; i++ {
		dp[i] = make([]float64, m+1)
		bt[i] = make([]operation, m+1)
		for j := 0; j <= m; j++ {
			dp[i][j] = math.Inf(-1)
		}
	}
	dp[0][0] = 0
	for i := 1; i <= n; i++ {
		dp[i][0] = dp[i-1][0] + gapPenalty
		bt[i][0] = operation{kind: "skipS", sig: i - 1, luca: -1}
	}
	for j := 1; j <= m; j++ {
		dp[0][j] = dp[0][j-1] + gapPenalty
		bt[0][j] = operation{kind: "skipL", sig: -1, luca: j - 1}
	}
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			best := dp[i-1][j] + gapPenalty
			op := operation{kind: "skipS", sig: i - 1, luca: -1}

			if score := dp[i][j-1] + gapPenalty; score > best {
				best = score
				op = operation{kind: "skipL", sig: -1, luca: j - 1}
			}

			sim := textSimilarity(siglus[i-1].prep, luca[j-1].prep)
			if score := dp[i-1][j-1] + sim - matchBase; score > best {
				best = score
				op = operation{kind: "match", sig: i - 1, luca: j - 1, score: sim}
			}

			dp[i][j] = best
			bt[i][j] = op
		}
	}

	ops := []operation{}
	for i, j := n, m; i > 0 || j > 0; {
		op := bt[i][j]
		ops = append(ops, op)
		switch op.kind {
		case "match":
			i--
			j--
		case "skipS":
			i--
		case "skipL":
			j--
		default:
			i, j = 0, 0
		}
	}
	for i, j := 0, len(ops)-1; i < j; i, j = i+1, j-1 {
		ops[i], ops[j] = ops[j], ops[i]
	}
	return ops
}

func groupedOps(ops []operation, kind string) [][]int {
	var groups [][]int
	var current []int
	for i, op := range ops {
		if op.kind == kind {
			current = append(current, i)
			continue
		}
		if len(current) > 0 {
			groups = append(groups, current)
			current = nil
		}
	}
	if len(current) > 0 {
		groups = append(groups, current)
	}
	return groups
}

func lucaIndexesForOps(ops []operation, group []int) []int {
	idxs := []int{}
	for _, opIndex := range group {
		if ops[opIndex].kind == "skipL" && ops[opIndex].luca >= 0 {
			idxs = append(idxs, ops[opIndex].luca)
		}
	}
	return idxs
}

func bestNearbySiglusCoverage(ops []operation, group []int, siglus []siglusEntry, luca []lucaEntry) (int, float64) {
	return bestSiglusCoverage(ops, group, siglus, luca, false)
}

func bestNearbySkippedSiglusCoverage(ops []operation, group []int, siglus []siglusEntry, luca []lucaEntry) (int, float64) {
	return bestSiglusCoverage(ops, group, siglus, luca, true)
}

func bestSiglusCoverage(ops []operation, group []int, siglus []siglusEntry, luca []lucaEntry, skippedOnly bool) (int, float64) {
	lucaIdxs := lucaIndexesForOps(ops, group)
	if len(lucaIdxs) == 0 {
		return -1, 0
	}
	groupTokens := map[string]bool{}
	for _, idx := range lucaIdxs {
		for tok := range luca[idx].prep.tokens {
			groupTokens[tok] = true
		}
	}
	if len(groupTokens) == 0 {
		return -1, 0
	}

	candidates := map[int]bool{}
	start := group[0] - 6
	end := group[len(group)-1] + 6
	if start < 0 {
		start = 0
	}
	if end >= len(ops) {
		end = len(ops) - 1
	}
	for i := start; i <= end; i++ {
		if skippedOnly && ops[i].kind != "skipS" {
			continue
		}
		if ops[i].sig >= 0 {
			candidates[ops[i].sig] = true
		}
	}

	bestIdx := -1
	bestCoverage := 0.0
	for sigIdx := range candidates {
		covered := 0
		for tok := range groupTokens {
			if siglus[sigIdx].prep.tokens[tok] > 0 {
				covered++
			}
		}
		coverage := float64(covered) / float64(len(groupTokens))
		if coverage > bestCoverage {
			bestCoverage = coverage
			bestIdx = sigIdx
		}
	}
	return bestIdx, bestCoverage
}

func findMatchedOpBySiglus(ops []operation, sigIdx, start, end int) int {
	best := -1
	bestDistance := math.MaxInt32
	for i, op := range ops {
		if op.kind != "match" || op.sig != sigIdx {
			continue
		}
		distance := 0
		if i < start {
			distance = start - i
		} else if i > end {
			distance = i - end
		}
		if distance < bestDistance {
			bestDistance = distance
			best = i
		}
	}
	return best
}

func markMergedMatchesForReview(ops []operation, siglus []siglusEntry, luca []lucaEntry, reviewOps map[int]string) {
	for opIndex, op := range ops {
		if op.kind != "match" || op.sig < 0 || op.luca < 0 {
			continue
		}
		sigTokens := len(siglus[op.sig].prep.tokens)
		lucaTokens := len(luca[op.luca].prep.tokens)
		if sigTokens < 14 || lucaTokens == 0 || float64(sigTokens)/float64(lucaTokens) < 1.8 {
			continue
		}

		covered := []int{opIndex}
		start := opIndex - 4
		end := opIndex + 4
		if start < 0 {
			start = 0
		}
		if end >= len(ops) {
			end = len(ops) - 1
		}
		for i := start; i <= end; i++ {
			if i == opIndex || ops[i].luca < 0 {
				continue
			}
			coverage := tokenCoverageBySiglus(siglus[op.sig], luca[ops[i].luca])
			if coverage >= 0.25 {
				covered = append(covered, i)
			}
		}
		if len(covered) < 2 {
			continue
		}
		for _, idx := range covered {
			reviewOps[idx] = "siglus_merged_or_luca_split"
		}
	}
}

func tokenCoverageBySiglus(sig siglusEntry, luc lucaEntry) float64 {
	if len(luc.prep.tokens) == 0 {
		return 0
	}
	covered := 0
	for tok := range luc.prep.tokens {
		if sig.prep.tokens[tok] > 0 {
			covered++
		}
	}
	return float64(covered) / float64(len(luc.prep.tokens))
}

func siglusEntryAt(entries []siglusEntry, idx int) siglusEntry {
	if idx < 0 || idx >= len(entries) {
		return siglusEntry{}
	}
	return entries[idx]
}

func makeReportRow(file, kind string, luca lucaEntry, score, coverage float64, siglus siglusEntry) reportRow {
	return reportRow{
		file:       file,
		kind:       kind,
		lucaID:     luca.id,
		line:       luca.fileLine + 1,
		tag:        luca.tag,
		score:      score,
		coverage:   coverage,
		lucaText:   luca.text,
		siglusID:   siglus.id,
		siglusText: siglus.source,
		siglusFR:   siglus.target,
	}
}

func writeReport(path string, rows []reportRow) error {
	sort.SliceStable(rows, func(i, j int) bool {
		if rows[i].file != rows[j].file {
			return rows[i].file < rows[j].file
		}
		return rows[i].line < rows[j].line
	})
	var sb strings.Builder
	sb.WriteString("File\tKind\tLucaID\tScriptLine\tTag\tScore\tCoverage\tLucaText\tSiglusID\tSiglusSource\tSiglusFR\n")
	for _, row := range rows {
		fields := []string{
			row.file,
			row.kind,
			row.lucaID,
			strconv.Itoa(row.line),
			row.tag,
			fmt.Sprintf("%.3f", row.score),
			fmt.Sprintf("%.3f", row.coverage),
			escapeTSV(row.lucaText),
			row.siglusID,
			escapeTSV(row.siglusText),
			escapeTSV(row.siglusFR),
		}
		sb.WriteString(strings.Join(fields, "\t"))
		sb.WriteString("\n")
	}
	return os.WriteFile(path, []byte(sb.String()), 0644)
}

func escapeTSV(s string) string {
	s = strings.ReplaceAll(s, "\r", "")
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\t", "\\t")
	return s
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0644)
}

func stripLabelPrefix(trimmed string) string {
	rest := trimmed
	for {
		next, ok := stripOneLabel(rest, "label")
		if !ok {
			next, ok = stripOneLabel(rest, "global")
		}
		if !ok {
			return rest
		}
		rest = next
	}
}

func stripOneLabel(trimmed, prefix string) (string, bool) {
	if !strings.HasPrefix(trimmed, prefix) {
		return trimmed, false
	}
	rest := trimmed[len(prefix):]
	i := 0
	for i < len(rest) && rest[i] >= '0' && rest[i] <= '9' {
		i++
	}
	if i == 0 || i >= len(rest) || rest[i] != ':' {
		return trimmed, false
	}
	return strings.TrimLeft(rest[i+1:], " \t"), true
}

func dialogueOpcode(trimmed string) (string, bool) {
	line := stripLabelPrefix(trimmed)
	for _, opcode := range []string{"MESSAGE", "LOG_BEGIN", "SELECT"} {
		if hasOpcodePrefix(line, opcode) {
			return opcode, true
		}
	}
	return "", false
}

func hasOpcodePrefix(line, opcode string) bool {
	if !strings.HasPrefix(line, opcode) {
		return false
	}
	if len(line) == len(opcode) {
		return true
	}
	switch line[len(opcode)] {
	case ' ', '\t', '(':
		return true
	default:
		return false
	}
}

func extractQuotedStrings(line string) []string {
	var out []string
	runes := []rune(line)
	inQuote := false
	escaped := false
	var current strings.Builder
	for _, ch := range runes {
		if !inQuote {
			if ch == '"' {
				inQuote = true
				escaped = false
				current.Reset()
			}
			continue
		}
		if escaped {
			switch ch {
			case 'n':
				current.WriteString("\\n")
			case 't':
				current.WriteString("\\t")
			case '"':
				current.WriteRune('"')
			case '\\':
				current.WriteRune('\\')
			default:
				current.WriteRune('\\')
				current.WriteRune(ch)
			}
			escaped = false
			continue
		}
		if ch == '\\' {
			escaped = true
			continue
		}
		if ch == '"' {
			out = append(out, current.String())
			inQuote = false
			continue
		}
		current.WriteRune(ch)
	}
	return out
}

func replaceNthQuotedString(line string, n int, newText string) string {
	quoted := extractQuotedStrings(line)
	if n >= 0 && n < len(quoted) {
		newText = preserveOriginalLineBreakSuffix(quoted[n], newText)
	}
	escaped := strings.ReplaceAll(newText, "\\", "\\\\")
	escaped = strings.ReplaceAll(escaped, "\"", "\\\"")
	escaped = strings.ReplaceAll(escaped, "\r", "")
	escaped = strings.ReplaceAll(escaped, "\n", "\\n")
	escaped = strings.ReplaceAll(escaped, "\t", "\\t")

	runes := []rune(line)
	quoteCount := 0
	inQuote := false
	skipTarget := false
	var result strings.Builder
	for i := 0; i < len(runes); i++ {
		ch := runes[i]
		if !inQuote {
			if ch == '"' {
				inQuote = true
				if quoteCount == n {
					skipTarget = true
					result.WriteRune('"')
					result.WriteString(escaped)
					continue
				}
			}
			result.WriteRune(ch)
			continue
		}
		if skipTarget {
			if ch == '\\' && i+1 < len(runes) {
				i++
				continue
			}
			if ch == '"' {
				result.WriteRune('"')
				inQuote = false
				skipTarget = false
				quoteCount++
			}
			continue
		}
		result.WriteRune(ch)
		if ch == '\\' && i+1 < len(runes) {
			i++
			result.WriteRune(runes[i])
			continue
		}
		if ch == '"' {
			inQuote = false
			quoteCount++
		}
	}
	return result.String()
}

func preserveOriginalLineBreakSuffix(original, replacement string) string {
	if strings.HasSuffix(replacement, "\\n") {
		replacement = strings.TrimSuffix(replacement, "\\n") + "\n"
	}
	if (strings.HasSuffix(original, "\\n") || strings.HasSuffix(original, "\n")) && !strings.HasSuffix(replacement, "\n") {
		return replacement + "\n"
	}
	return replacement
}

func prepareText(s string) preparedText {
	norm := strings.ToLower(normalizeText(s))
	norm = strings.ReplaceAll(norm, "can't", "cannot")
	norm = strings.ReplaceAll(norm, "won't", "will not")
	norm = strings.ReplaceAll(norm, "n't", " not")
	tokens := map[string]int{}
	for _, tok := range wordRE.FindAllString(norm, -1) {
		tok = normalizeToken(tok)
		if len(tok) <= 1 || stopWords[tok] {
			continue
		}
		tokens[tok]++
	}
	var mag float64
	for _, count := range tokens {
		mag += float64(count * count)
	}
	return preparedText{norm: norm, tokens: tokens, mag: math.Sqrt(mag)}
}

func normalizeToken(tok string) string {
	if alias, ok := tokenAliases[tok]; ok {
		return alias
	}
	switch {
	case len(tok) > 5 && strings.HasSuffix(tok, "ies"):
		tok = strings.TrimSuffix(tok, "ies") + "y"
	case len(tok) > 5 && strings.HasSuffix(tok, "ing"):
		tok = strings.TrimSuffix(tok, "ing")
	case len(tok) > 4 && strings.HasSuffix(tok, "ed"):
		tok = strings.TrimSuffix(tok, "ed")
	case len(tok) > 4 && strings.HasSuffix(tok, "ly"):
		tok = strings.TrimSuffix(tok, "ly")
	case len(tok) > 4 && strings.HasSuffix(tok, "s"):
		tok = strings.TrimSuffix(tok, "s")
	}
	if alias, ok := tokenAliases[tok]; ok {
		return alias
	}
	return tok
}

func textSimilarity(a, b preparedText) float64 {
	if a.norm != "" && a.norm == b.norm {
		return 1
	}
	if len(a.tokens) == 0 || len(b.tokens) == 0 {
		return 0
	}
	intersection := 0
	dot := 0
	for tok, countA := range a.tokens {
		if countB, ok := b.tokens[tok]; ok {
			intersection++
			dot += countA * countB
		}
	}
	if intersection == 0 {
		return 0
	}
	union := len(a.tokens) + len(b.tokens) - intersection
	minLen := len(a.tokens)
	if len(b.tokens) < minLen {
		minLen = len(b.tokens)
	}
	cosine := float64(dot) / (a.mag * b.mag)
	jaccard := float64(intersection) / float64(union)
	containment := float64(intersection) / float64(minLen)
	return 0.58*cosine + 0.24*jaccard + 0.18*containment
}

func normalizeText(s string) string {
	replacer := strings.NewReplacer(
		"\\n", " ", "\r", " ", "\n", " ",
		"❝", "\"", "❞", "\"", "“", "\"", "”", "\"", "„", "\"",
		"❛", "'", "❜", "'", "‘", "'", "’", "'", "＇", "'",
		"﹣", "-", "－", "-", "–", "-", "—", "-", "―", "-",
		"…", "...", "　", " ", "\ufeff", "",
	)
	s = replacer.Replace(s)
	fields := strings.Fields(s)
	return strings.Join(fields, " ")
}

func isControlText(s string) bool {
	x := normalizeText(strings.TrimSpace(s))
	if x == "" {
		return true
	}
	if speakerOrControl[x] {
		return true
	}
	switch x {
	case "_", "0", "1", "2", "01", "02", "03", "04", "ja", "en", "jp", "cn", "fr", "L", "R", "M", "S", "CG", "BG", "SE", "BGM", "_EN":
		return true
	}
	lower := strings.ToLower(x)
	if lower == "dummy" || lower == "attack" || lower == "intro1" || lower == "intro2" || lower == "intro3" {
		return true
	}
	if strings.HasPrefix(x, "$") && !strings.ContainsAny(x, " \t") {
		return true
	}
	if assetRE.MatchString(x) || trophyRE.MatchString(x) {
		return true
	}
	if numberedIDRE.MatchString(x) && !strings.ContainsAny(x, " \t") {
		return true
	}
	if strings.Contains(x, "_") && !strings.ContainsAny(x, " \t") {
		return true
	}
	if isShortJapaneseLabel(x) {
		return true
	}
	return false
}

func isMinorLucaOnlyText(s string) bool {
	x := normalizeText(strings.TrimSpace(s))
	if x == "" {
		return true
	}
	if !containsLatinLetter(x) && !containsJapaneseRune(x) {
		return true
	}
	stripped := strings.Trim(x, "\"'“”‘’❝❞.?!,;:…-—–﹣－― ")
	if stripped == "" {
		return true
	}
	words := wordRE.FindAllString(strings.ToLower(x), -1)
	quotedDialogue := strings.HasPrefix(x, "❝") || strings.HasPrefix(x, "\"")
	if quotedDialogue && len(words) <= 2 && len([]rune(x)) <= 28 {
		return true
	}
	return false
}

func containsTextLetter(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) {
			return true
		}
	}
	return false
}

func containsLatinLetter(s string) bool {
	for _, r := range s {
		if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || (r >= 0x00C0 && r <= 0x024F) {
			return true
		}
	}
	return false
}

func isShortJapaneseLabel(s string) bool {
	runes := []rune(s)
	if len(runes) == 0 || len(runes) > 4 {
		return false
	}
	for _, r := range runes {
		if !isJapaneseRune(r) {
			return false
		}
	}
	return true
}

func containsJapaneseRune(s string) bool {
	for _, r := range s {
		if isJapaneseRune(r) {
			return true
		}
	}
	return false
}

func isJapaneseRune(r rune) bool {
	return (r >= 0x3040 && r <= 0x30ff) || (r >= 0x3400 && r <= 0x9fff)
}
