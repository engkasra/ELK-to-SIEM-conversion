package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"sort"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	// 	for example:
	// Enter tag of log: http_ngmi_uat
	// Enter rule name of filter: syslog_ngmi_uat_access
	// Enter final nested: NGMI
	// Browse the service log file: service_log.json

	fmt.Print("Enter tag of log: ")
	matchTag, _ := reader.ReadString('\n')
	matchTag = strings.TrimSpace(matchTag)

	fmt.Print("Enter rule name of filter: ")
	newTag, _ := reader.ReadString('\n')
	newTag = strings.TrimSpace(newTag)

	fmt.Print("Enter final nested: ")
	finalNestKey, _ := reader.ReadString('\n')
	finalNestKey = strings.TrimSpace(finalNestKey)

	fmt.Print("Browse the service log file: ")
	filePath, _ := reader.ReadString('\n')
	filePath = strings.TrimSpace(filePath)

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read file: %v\n", err)
		os.Exit(1)
	}

	var input map[string]interface{}
	if err := json.Unmarshal(data, &input); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse JSON: %v\n", err)
		os.Exit(1)
	}

	output := []string{
		"[FILTER]",
		"  Name          rewrite_tag",
		fmt.Sprintf("  Match         %s", matchTag),
		fmt.Sprintf("  Rule          $service_log['event'] ^(access).* %s true", newTag),
		"  Emitter_Name  re_emitted",
		"",
	}

	flatKeys := map[string]bool{}
	extractNestedKeys(input, "", flatKeys)

	keys := make([]string, 0, len(flatKeys))
	for k := range flatKeys {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		output = append(output, "[FILTER]")
		output = append(output, "  Name nest")
		output = append(output, fmt.Sprintf("  Match %s", newTag))
		output = append(output, "  Operation lift")
		output = append(output, fmt.Sprintf("  Nested_under %s", key))
		output = append(output, fmt.Sprintf("  Add_prefix %s_", key))
		output = append(output, "")
	}

	output = append(output, "[FILTER]")
	output = append(output, "  Name nest")
	output = append(output, fmt.Sprintf("  Match %s", newTag))
	output = append(output, "  Operation nest")
	output = append(output, "  Wildcard service_log_*")
	output = append(output, fmt.Sprintf("  Nest_under %s", finalNestKey))

	if err := ioutil.WriteFile("filters.conf", []byte(strings.Join(output, "\n")), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write filters.conf: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("filters.conf generated successfully.")
}

func extractNestedKeys(data interface{}, prefix string, result map[string]bool) {
	m, ok := data.(map[string]interface{})
	if !ok {
		return
	}
	for k, v := range m {
		fullKey := k
		if prefix != "" {
			fullKey = prefix + "_" + k
		}
		valKind := reflect.TypeOf(v).Kind()
		if valKind == reflect.Map {
			result[fullKey] = true
			extractNestedKeys(v, fullKey, result)
		} else if valKind == reflect.Slice {
			slice := v.([]interface{})
			for _, item := range slice {
				extractNestedKeys(item, fullKey, result)
			}
		}
	}
}
