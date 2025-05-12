# Fluent Bit Filter Generator

This Go program reads a sample ELK-style JSON log (e.g. from `service_log.json`), inspects its nested structure, and generates a Fluent Bit `filters.conf` file that mimics how nested fields would be handled using Fluent Bit filters like `nest`, `lift`, and `rewrite_tag`.

---

## ✨ Features

- Prompts for dynamic values like `Match`, `Rule`, and `Final Nest Key`
- Recursively walks through nested maps and arrays
- Generates `[FILTER]` blocks for all nested levels
- Produces a ready-to-use `filters.conf` file

---

## 📦 Requirements

- Go 1.18+
- A JSON log file in ELK format with a `service_log` structure

---

## 🚀 Running the Script

### 🧑‍💻 Interactive Mode

```bash

go run main.go

```

The program will prompt for:

Enter tag of log: http_ngmi_uat
Enter rule name of filter: syslog_ngmi_uat_access
Enter final nested: NGMI
Browse the service log file: service_log.json

---

Then it generates a `filters.conf` file in the current directory.

---
🔄 Example
Given this sample JSON:

```json
{
  "service_log": {
    "event": "access.api_call",
    "access": {
      "request": {
        "headers": {
          "x_correlation_id": "abc"
        }
      }
    }
  }
}

```

The output `filters.conf` will look like:

```ini
[FILTER]
  Name          rewrite_tag
  Match         http_ngmi_uat
  Rule          $service_log['event'] ^(access).* syslog_ngmi_uat_access true
  Emitter_Name  re_emitted

[FILTER]
  Name nest
  Match syslog_ngmi_uat_access
  Operation lift
  Nested_under service_log
  Add_prefix service_log_

[FILTER]
  Name nest
  Match syslog_ngmi_uat_access
  Operation lift
  Nested_under service_log_access
  Add_prefix service_log_access_

...

[FILTER]
  Name nest
  Match syslog_ngmi_uat_access
  Operation nest
  Wildcard service_log_*
  Nest_under NGMI

```

📁 Output
The result is saved as:
filters.conf

Place this file into your Fluent Bit configuration directory to apply dynamic field flattening and nesting logic.

---

🤝 Contribution
Feel free to submit pull requests or suggest improvements via issues. Designed with extensibility in mind.

---

🧑‍💼 Author
engksa75@gmail.com
