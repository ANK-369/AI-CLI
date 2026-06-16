package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

// ANSI Terminal Colors & Layout Styles
const (
	Reset       = "\033[0m"
	Bold        = "\033[1m"
	Dim         = "\033[2m"
	Red         = "\033[31m"
	Green       = "\033[32m"
	Yellow      = "\033[33m"
	Blue        = "\033[34m"
	Magenta     = "\033[35m"
	Cyan        = "\033[36m"
	White       = "\033[37m"
	BgCyan      = "\033[46m"
	BgMagenta   = "\033[45m"
)

// Current Local Version
const LocalVersion = "4.0"

// Task and Structured Response Schemas
type Task struct {
	Tool      string   `json:"tool"`
	Arguments []string `json:"arguments"`
}

type AIResponse struct {
	Context           string `json:"context"`             // የአሁኑ ሁኔታ መግለጫ
	TacticalPhase     string `json:"tactical_phase"`      // RECON, SCANNING, VULN_ANALYSIS, EXPLOIT, REPORTING
	RecommendedAction Task   `json:"recommended_action"`  // የሚመከር ቀጣይ እርምጃ
	Reasoning         string `json:"reasoning"`           // ምክንያት
	Status            string `json:"status"`              // "CONTINUE", "PROGRESS_PHASE", "COMPLETE"
}

type Content struct {
	Role  string `json:"role"`
	Parts []Part `json:"parts"`
}

type Part struct {
	Text string `json:"text"`
}

type GeminiRequest struct {
	Contents          []Content         `json:"contents"`
	SystemInstruction SystemInstruction `json:"systemInstruction"`
	GenerationConfig  GenerationConfig  `json:"generationConfig"`
}

type SystemInstruction struct {
	Parts []Part `json:"parts"`
}

type GenerationConfig struct {
	ResponseMimeType string       `json:"responseMimeType"`
	ResponseSchema   SchemaObject `json:"responseSchema"`
}

type SchemaObject struct {
	Type       string                 `json:"type"`
	Properties map[string]SchemaField `json:"properties"`
	Required   []string               `json:"required"`
}

type SchemaField struct {
	Type  string                 `json:"type"`
	Items *SchemaField           `json:"items,omitempty"`
	Props map[string]SchemaField `json:"properties,omitempty"`
	Req   []string               `json:"required,omitempty"`
}

var ConversationHistory []Content
var CurrentState string = "INIT"         // የሲስተሙ ቴክኒካል ስቴት (UI Engine)
var ActiveTacticalPhase string = "RECON" // የሳይበር ደህንነት ታክቲካዊ ምዕራፍ (Cyber Security Phase)

// Automated GitHub Self-Updater with Auto-Recompilation
func checkAndApplyUpdate() {
	fmt.Printf("%s[*] Checking GitHub for new updates...%s\r", Dim+Yellow, Reset)
	
	versionURL := "https://raw.githubusercontent.com/ANK-369/AI-CLI/main/version.txt"
	repoURL    := "https://raw.githubusercontent.com/ANK-369/AI-CLI/main/main.go"

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get(versionURL)
	if err != nil {
		fmt.Printf("%s\r[-] Update Check Failed: %v%s\n", Red, err, Reset)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Print("\r\033[K")
		return
	}

	versionBytes, _ := io.ReadAll(resp.Body)
	latestVersion := strings.TrimSpace(string(versionBytes))

	if latestVersion != "" && latestVersion != LocalVersion {
		fmt.Printf("\r\033[K%s[+] New version detected: v%s! Syncing codebase...%s\n", Green, latestVersion, Reset)
		
		respGo, err := client.Get(repoURL)
		if err != nil {
			fmt.Printf("%s[-] Failed to stream updated source code: %v%s\n", Red, err, Reset)
			return
		}
		defer respGo.Body.Close()

		if respGo.StatusCode == http.StatusOK {
			newCode, err := io.ReadAll(respGo.Body)
			if err == nil {
				err = os.WriteFile("main.go", newCode, 0644)
				if err == nil {
					fmt.Printf("%s[*] Optimization Engine: Recompiling aicli binary...%s\r", Dim+Yellow, Reset)
					
					os.Remove("aicli") 
					
					buildCmd := exec.Command("go", "build", "-ldflags=-s -w", "-o", "aicli", "main.go")
					if buildErr := buildCmd.Run(); buildErr != nil {
						fmt.Printf("%s[-] Recompilation failed: %v%s\n", Red, buildErr, Reset)
						return
					}
					
					_ = os.Remove("main.go")
					
					fmt.Printf("%s[+] AI-CLI successfully updated and optimized to v%s! Please restart ./aicli%s\n", Green, latestVersion, Reset)
					os.Exit(0)
				} else {
					fmt.Printf("%s[-] Failed to write update payload: %v%s\n", Red, err, Reset)
				}
			}
		}
	} else {
		fmt.Printf("\r\033[K%s[+] AI-CLI is already up to date (v%s).%s\n", Green, LocalVersion, Reset)
	}
}

func main() {
	// 💡 አውቶማቲክ ቡትስትራፕ እና ኢንቫይሮንመንት ማጽጃ ሎጂክ
	if !strings.Contains(os.Args[0], "aicli") {
		// 1. 'aicli' ባይነሪ ፋይል በአካባቢው ከሌለ ራሱ በራስ-ሰር ይገነባዋል
		if _, err := os.Stat("aicli"); os.IsNotExist(err) {
			fmt.Printf("%s[*] First-time Initialization: Generating optimized 'aicli' binary...%s\n", Yellow, Reset)
			buildCmd := exec.Command("go", "build", "-ldflags=-s -w", "-o", "aicli", "main.go")
			if err := buildCmd.Run(); err != nil {
				fmt.Printf("%s[-] Auto-build failed: %v%s\n", Red, err, Reset)
				return
			}
			fmt.Printf("%s[+] Success! 'aicli' binary generated.%s\n", Green, Reset)
		}

		// 2. 'main.go' የምንጭ ኮዱን ሙሉ በሙሉ ይሰርዛል (Clean Environment)
		if _, err := os.Stat("main.go"); err == nil {
			_ = os.Remove("main.go")
		}

		// 3. አዲስ የተፈጠረውን `./aicli` በራሱ በራስ-ሰር ተርሚናሉ ላይ ያስነሳል
		cmd := exec.Command("./aicli")
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("%s[-] Execution failed: %v%s\n", Red, err, Reset)
		}
		
		// የ 'go run' ጊዜያዊ ፕሮሰስን እዚህ ላይ ያበቃል
		os.Exit(0)
	}

	// የራስ-አፕዴት ፈንክሽን
	checkAndApplyUpdate()

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		fmt.Printf("%s[-] Error: GEMINI_API_KEY environment variable is not set.%s\n", Red, Reset)
		return
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("%s%s┌─────────────────────────────────────────┐%s\n", Bold, Cyan, Reset)
	fmt.Printf("%s%s│  AI-CLI ANdualem Koriya [ANK369] v%s   │%s\n", Bold, Cyan, LocalVersion, Reset)
	fmt.Printf("%s%s└─────────────────────────────────────────┘%s\n", Bold, Cyan, Reset)
	fmt.Printf("%s[*] Sequential 5-Stage Orchestration Loop: Active%s\n", Dim+White, Reset)

	for {
		fmt.Printf("\n%s [%s ❯ %s] ai-prompt ❯ %s", Bold+Magenta, CurrentState, ActiveTacticalPhase, Reset)
		if !scanner.Scan() {
			break
		}
		userInput := strings.TrimSpace(scanner.Text())

		if userInput == "exit" || userInput == "quit" {
			fmt.Printf("%s[*] Flashing session cache. Safe shutdown complete.%s\n", Yellow, Reset)
			break
		}

		if userInput == "" {
			continue
		}

		ConversationHistory = append(ConversationHistory, Content{
			Role:  "user",
			Parts: []Part{{Text: userInput}},
		})

		var executionSummaryCollector strings.Builder

		// True Feedback Loop (Step-by-Step Strategic Transitions)
		for {
			CurrentState = "PLANNING"
			fmt.Printf("%s[*] Analyzing historic context & mapping operational phases...%s\r", Dim+Yellow, Reset)
			
			aiCmd, rawJsonText, err := askGemini(apiKey, ConversationHistory)
			if err != nil {
				fmt.Printf("%s\r\033[K[-] Orchestration Network Error: %v%s\n", Red, err, Reset)
				break
			}
			fmt.Print("\r\033[K")

			// የታሪክ ማጠራቀሚያ
			ConversationHistory = append(ConversationHistory, Content{
				Role:  "model",
				Parts: []Part{{Text: rawJsonText}},
			})

			// የታክቲካዊ ምዕራፍ ማመሳሰያ
			if aiCmd.TacticalPhase != "" {
				ActiveTacticalPhase = aiCmd.TacticalPhase
			}

			// የመጨረሻው ሪፖርት ማውጫ ቅድመ ሁኔታ
			if aiCmd.Status == "COMPLETE" || aiCmd.TacticalPhase == "REPORTING" || aiCmd.RecommendedAction.Tool == "" {
				CurrentState = "IDLE"
				ActiveTacticalPhase = "REPORTING"
				if executionSummaryCollector.Len() > 0 {
					fmt.Printf("%s[*] Compiling metrics into responsive layout...%s\r", Dim+Yellow, Reset)
					generateAndPrintStructuredSummary(apiKey, executionSummaryCollector.String())
				} else {
					renderFriendlyResponse(apiKey)
				}
				break
			}

			// Recommended Plan Schema ማሳያ ቦክስ
			CurrentState = "PENDING_APPROVAL"
			fmt.Printf("\n%s┌─────────── RECOMMENDED PLAN [%s] ───────────┐%s\n", Blue, ActiveTacticalPhase, Reset)
			fmt.Printf("%s│ Context:%s %s\n", Yellow, White, aiCmd.Context)
			fmt.Printf("%s│ Action :%s %s %s\n", Yellow, Green, aiCmd.RecommendedAction.Tool, strings.Join(aiCmd.RecommendedAction.Arguments, " "))
			fmt.Printf("%s│ Reason :%s %s\n", Yellow, White, aiCmd.Reasoning)
			fmt.Printf("%s└──────────────────────────────────────────────────────┘%s\n", Blue, Reset)

			// User Approval Gatekeeper
			fmt.Printf("%s[?] Approve and execute this step? (y/n/cancel): %s", Bold+White, Reset)
			if !scanner.Scan() {
				break
			}
			approval := strings.ToLower(strings.TrimSpace(scanner.Text()))

			if approval == "n" || approval == "cancel" {
				CurrentState = "INTERRUPTED"
				fmt.Printf("%s[*] Step rejected. Returning to prompt for manual intervention.%s\n", Yellow, Reset)
				break
			}

			if approval != "y" {
				fmt.Printf("%s[-] Invalid input. Action aborted.%s\n", Red, Reset)
				break
			}

			// Execution Engine
			CurrentState = "EXECUTING"
			targetDomain := extractTargetDomain(aiCmd.RecommendedAction.Arguments)
			
			output, success := runToolWithLiveDashboard(targetDomain, aiCmd.RecommendedAction.Tool, ActiveTacticalPhase, 0, aiCmd.RecommendedAction.Arguments)
			
			executionSummaryCollector.WriteString(fmt.Sprintf("[PHASE: %s][TOOL: %s]\nCommand: %s %s\nSuccess: %t\nLogs:\n%s\n[END_DATA]\n", 
				ActiveTacticalPhase, aiCmd.RecommendedAction.Tool, aiCmd.RecommendedAction.Tool, strings.Join(aiCmd.RecommendedAction.Arguments, " "), success, output))

			// Feedback Loop Insertion
			ConversationHistory = append(ConversationHistory, Content{
				Role:  "user",
				Parts: []Part{{Text: fmt.Sprintf("Observation of %s execution in %s phase: Success=%t. Output:\n%s", aiCmd.RecommendedAction.Tool, ActiveTacticalPhase, success, output)}},
			})

			// Dynamic Status Transition handling
			if aiCmd.Status == "PROGRESS_PHASE" {
				fmt.Printf("%s[➔] Tactical phase objective achieved. Progressing to next strategic stage...%s\n", Green, Reset)
			} else if !success {
				CurrentState = "FALLBACK_ROUTING"
				fmt.Printf("%s[⚠️] Action Failed. Feedback routed back to AI Engine for structural alternative plan.%s\n", Red, Reset)
			} else {
				CurrentState = "ANALYZING"
				fmt.Printf("%s[+] Action successful. Moving deeper into current node...%s\n", Green, Reset)
			}
			
			time.Sleep(1 * time.Second)
		}
	}
}

func extractTargetDomain(args []string) string {
	for _, arg := range args {
		if !strings.HasPrefix(arg, "-") && strings.Contains(arg, ".") {
			return arg
		}
	}
	if len(args) > 0 {
		return args[len(args)-1]
	}
	return "Active Target"
}

func printDashboard(target, phase, tool, duration string, retry, percent int, firstPrint bool) {
	if !firstPrint {
		fmt.Print("\033[12A") 
	}
	barWidth := 20
	filled := (percent * barWidth) / 100
	barStr := strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)

	fmt.Printf("\r\033[K%s┌─────────────────────────────────────────┐%s\n", Cyan, Reset)
	fmt.Printf("\r\033[K%s│ %sAI-CLI Pipeline Orchestrator%s            │%s\n", Cyan, Bold+White, Cyan, Reset)
	fmt.Printf("\r\033[K%s├─────────────────────────────────────────┤%s\n", Cyan, Reset)
	fmt.Printf("\r\033[K%s│ %sTarget:%s %-31s %s│%s\n", Cyan, Yellow, White, padStr(target, 31), Cyan, Reset)
	fmt.Printf("\r\033[K%s│ %sPhase :%s %-31s %s│%s\n", Cyan, Yellow, White, padStr(phase, 31), Cyan, Reset)
	fmt.Printf("\r\033[K%s│ %sTool  :%s %-31s %s│%s\n", Cyan, Yellow, White, padStr(tool, 31), Cyan, Reset)
	fmt.Printf("\r\033[K%s│ %sTime  :%s %-31s %s│%s\n", Cyan, Yellow, White, padStr(duration, 31), Cyan, Reset)
	fmt.Printf("\r\033[K%s│ %sRetry :%s %-31d %s│%s\n", Cyan, Yellow, White, retry, Cyan, Reset)
	fmt.Printf("\r\033[K%s├─────────────────────────────────────────┤%s\n", Cyan, Reset)
	fmt.Printf("\r\033[K%s│ %sProgress%s                                │%s\n", Cyan, Bold+White, Cyan, Reset)
	fmt.Printf("\r\033[K%s│ %s%-20s %3d%%               %s│%s\n", Cyan, Green, barStr, percent, Cyan, Reset)
	fmt.Printf("\r\033[K%s└─────────────────────────────────────────┘%s\n", Cyan, Reset)
}

func padStr(s string, max int) string {
	if len(s) > max {
		return s[:max-3] + "..."
	}
	return s
}

func runToolWithLiveDashboard(target, tool, phase string, retry int, arguments []string) (string, bool) {
	if _, err := exec.LookPath(tool); err != nil {
		return fmt.Sprintf("Execution Error: %s missing on localized PATH.", tool), false
	}

	stopChan := make(chan bool)
	doneChan := make(chan bool)
	startTime := time.Now()

	go func() {
		first := true
		percent := 5
		ticker := time.NewTicker(400 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-stopChan:
				durationStr := fmt.Sprintf("%02dm %02ds", int(time.Since(startTime).Minutes()), int(time.Since(startTime).Seconds())%60)
				printDashboard(target, "Completed", tool, durationStr, retry, 100, first)
				doneChan <- true
				return
			case <-ticker.C:
				elapsed := time.Since(startTime)
				durationStr := fmt.Sprintf("%02dm %02ds", int(elapsed.Minutes()), int(elapsed.Seconds())%60)
				if percent < 92 {
					percent += 4
				}
				printDashboard(target, phase, tool, durationStr, retry, percent, first)
				first = false
			}
		}
	}()

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd := exec.Command(tool, arguments...)
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	close(stopChan)
	<-doneChan 

	if err != nil {
		return stdoutBuf.String() + "\n" + stderrBuf.String(), false
	}
	return stdoutBuf.String(), true
}

func generateAndPrintStructuredSummary(apiKey string, allOutputs string) {
	apiURL := "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash:generateContent?key=" + apiKey

	summaryPrompt := "You are an elite automated security auditing parser. Analyze the given tool raw outputs and map out a structured matrix report. " +
		"CRITICAL STYLING MANIFESTO:\n" +
		"1. NEVER use standard markdown characters (such as *, #, `, _).\n" +
		"2. Keep EVERY bullet list line ultra-short and precise (UNDER 50 characters per line) so it remains perfectly responsive and never wraps on mobile terminal windows.\n" +
		"3. You must group and isolate findings for EACH tool separately using our tags.\n\n" +
		"Strict Template Formats:\n" +
		"[TARGET_HEADER]\nTarget specifications and target metadata goes here.\n[END_SECTION]\n\n" +
		"[TOOL_REPORT: <Tool Name Here>]\n• Grouped brief clean responsive bullet points for this specific tool.\n[END_SECTION]\n\n" +
		"[CRITICAL_ALERTS]\n• Alarming exposures discovered across vectors.\n[END_SECTION]\n\n" +
		"[NEXT_STEPS_PROMPT]\nAsk the user an interactive question outlining the exact next step recommendation for verification.\n[END_SECTION]\n\n" +
		"Data Stream Logs:\n" + allOutputs

	reqPayload := map[string]interface{}{
		"contents": []map[string]interface{}{
			{"role": "user", "parts": []map[string]interface{}{{"text": summaryPrompt}}},
		},
	}

	jsonData, _ := json.Marshal(reqPayload)
	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("[-] Matrix reporting engine unreachable.")
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var rawResponse map[string]interface{}
	json.Unmarshal(body, &rawResponse)

	candidates, ok := rawResponse["candidates"].([]interface{})
	if !ok || len(candidates) == 0 {
		fmt.Println("[-] Failed to safely stream processed data matrices.")
		return
	}
	candidate := candidates[0].(map[string]interface{})
	content := candidate["content"].(map[string]interface{})
	parts := content["parts"].([]interface{})
	part := parts[0].(map[string]interface{})
	summaryText := part["text"].(string)

	ConversationHistory = append(ConversationHistory, Content{
		Role:  "model",
		Parts: []Part{{Text: summaryText}},
	})

	lines := strings.Split(summaryText, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || trimmed == "[END_SECTION]" {
			continue
		}

		if trimmed == "[TARGET_HEADER]" {
			fmt.Printf("\n%s%s┌─────────────────────────────────────────┐%s\n", Bold, Blue, Reset)
			fmt.Printf("%s%s│       TARGET RISK REPORT MATRIX         │%s\n", Bold, Blue, Reset)
			fmt.Printf("%s%s└─────────────────────────────────────────┘%s\n", Bold, Blue, Reset)
			continue
		}

		if strings.HasPrefix(trimmed, "[TOOL_REPORT:") {
			toolName := strings.TrimSuffix(strings.TrimPrefix(trimmed, "[TOOL_REPORT:"), "]")
			fmt.Printf("\n%s%s%s ENGINE AUDIT SUMMARY: %s %s\n", Bold, White, BgCyan, strings.ToUpper(toolName), Reset)
			fmt.Printf("%s%s───────────────────────────────────────────%s\n", Dim, Cyan, Reset)
			continue
		}

		if trimmed == "[CRITICAL_ALERTS]" {
			fmt.Printf("\n%s%s%s ISOLATED VULNERABILITY VECTORS %s\n", Bold, White, BgMagenta, Reset)
			fmt.Printf("%s%s───────────────────────────────────────────%s\n", Dim, Magenta, Reset)
			continue
		}

		if trimmed == "[NEXT_STEPS_PROMPT]" {
			fmt.Printf("\n%s%s RECOMMENDED ACTION PLAN & INTERACTIVE PROMPT:%s\n", Bold, Yellow, Reset)
			fmt.Printf("%s%s───────────────────────────────────────────%s\n", Dim, Yellow, Reset)
			continue
		}

		if strings.HasPrefix(trimmed, "•") || strings.HasPrefix(trimmed, "-") {
			fmt.Printf("   %s%s%s\n", Green, trimmed, Reset)
		} else if strings.Contains(trimmed, "Critical") || strings.Contains(trimmed, "High") {
			fmt.Printf("   %s%s%s\n", Red, trimmed, Reset)
		} else {
			fmt.Printf("   %s%s\n", White, trimmed)
		}
	}
	fmt.Printf("\n%s%s───────────────────────────────────────────%s\n", Dim, Blue, Reset)
}

func renderFriendlyResponse(apiKey string) {
	apiURL := "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash:generateContent?key=" + apiKey

	reqPayload := map[string]interface{}{"contents": ConversationHistory}
	jsonData, _ := json.Marshal(reqPayload)
	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var rawResponse map[string]interface{}
	json.Unmarshal(body, &rawResponse)

	candidates := rawResponse["candidates"].([]interface{})
	candidate := candidates[0].(map[string]interface{})
	content := candidate["content"].(map[string]interface{})
	parts := content["parts"].([]interface{})
	part := parts[0].(map[string]interface{})
	responseText := part["text"].(string)

	ConversationHistory = append(ConversationHistory, Content{
		Role:  "model",
		Parts: []Part{{Text: responseText}},
	})

	responseText = strings.ReplaceAll(responseText, "**", "")
	responseText = strings.ReplaceAll(responseText, "##", "")
	responseText = strings.ReplaceAll(responseText, "`", "")

	fmt.Printf("\n%s%s AI Response Engine:%s\n  %s\n", Bold, Cyan, Reset, responseText)
}

func askGemini(apiKey string, history []Content) (*AIResponse, string, error) {
	apiURL := "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash:generateContent?key=" + apiKey

	systemRule := "You are an elite automated security orchestration engine built for a Kali Linux PRoot environment on Android. " +
		"You operate inside a strict 5-stage tactical security framework. You MUST execute them in sequential order based on current findings:\n" +
		"1. 'RECON' (Subdomain discovery, OSINT)\n" +
		"2. 'SCANNING' (Port scanning, service enumeration)\n" +
		"3. 'VULN_ANALYSIS' (Vulnerability analysis, script scanning)\n" +
		"4. 'EXPLOIT' (Exploit simulation, payload verification)\n" +
		"5. 'REPORTING' (Final matrices aggregation)\n\n" +
		"CRITICAL OPERATIONAL RULES:\n" +
		"- Evaluate the historic tool observations. Do NOT loops tools or re-scan targets endlessly in the same phase once sufficient findings are available.\n" +
		"- When a phase objective is fulfilled, you MUST progress the lifecycle by setting 'tactical_phase' to the next stage and 'status' to 'PROGRESS_PHASE'.\n" +
		"- If a tool fails, dynamically route to an alternative fallback action ('FALLBACK_ROUTING').\n" +
		"- Any 'nmap' execution MUST include '--unprivileged' and '-sT'.\n" +
		"- Return strict JSON conforming to the structural Schema object."

	reqPayload := GeminiRequest{
		Contents: history,
		SystemInstruction: SystemInstruction{
			Parts: []Part{{Text: systemRule}},
		},
		GenerationConfig: GenerationConfig{
			ResponseMimeType: "application/json",
			ResponseSchema: SchemaObject{
				Type: "OBJECT",
				Properties: map[string]SchemaField{
					"context":        {Type: "STRING"},
					"tactical_phase": {Type: "STRING"}, 
					"reasoning":      {Type: "STRING"},
					"status":         {Type: "STRING"}, 
					"recommended_action": {
						Type: "OBJECT",
						Props: map[string]SchemaField{
							"tool":      {Type: "STRING"},
							"arguments": {Type: "ARRAY", Items: &SchemaField{Type: "STRING"}},
						},
						Req: []string{"tool", "arguments"},
					},
				},
				Required: []string{"context", "tactical_phase", "recommended_action", "reasoning", "status"},
			},
		},
	}

	jsonData, err := json.Marshal(reqPayload)
	if err != nil {
		return nil, "", err
	}

	for attempt := 1; attempt <= 2; attempt++ {
		client := &http.Client{Timeout: 30 * time.Second}
		resp, err := client.Post(apiURL, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			return nil, "", err
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close() 
		if err != nil {
			return nil, "", err
		}

		if resp.StatusCode == 429 {
			fmt.Printf("\r\033[K%s[⚠️] API Rate Limit Hit! Cooling down for 25s before auto-retry...%s\n", Yellow, Reset)
			for i := 25; i > 0; i-- {
				fmt.Printf("\r%s[⏳] Retrying in %ds...%s", Dim+White, i, Reset)
				time.Sleep(1 * time.Second)
			}
			fmt.Print("\r\033[K") 
			continue 
		}

		if resp.StatusCode != http.StatusOK {
			return nil, "", fmt.Errorf("API Error (Status %d): %s", resp.StatusCode, string(body))
		}

		var rawResponse map[string]interface{}
		if err := json.Unmarshal(body, &rawResponse); err != nil {
			return nil, "", err
		}

		candidates, ok := rawResponse["candidates"].([]interface{})
		if !ok || len(candidates) == 0 {
			return nil, "", fmt.Errorf("empty core framework payload. Raw Response: %s", string(body))
		}

		candidate := candidates[0].(map[string]interface{})
		content, ok := candidate["content"].(map[string]interface{})
		if !ok {
			return nil, "", fmt.Errorf("candidate content missing")
		}

		parts, ok := content["parts"].([]interface{})
		if !ok || len(parts) == 0 {
			return nil, "", fmt.Errorf("parts missing")
		}

		part := parts[0].(map[string]interface{})
		cleanJSONText := part["text"].(string)

		var finalResponse AIResponse
		if err := json.Unmarshal([]byte(cleanJSONText), &finalResponse); err != nil {
			return nil, "", err
		}

		return &finalResponse, cleanJSONText, nil
	}

	return nil, "", fmt.Errorf("failed after maximum API rate-limit retries")
}
