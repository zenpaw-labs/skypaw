<h2 align="center">skypaw: Config</h2>

## ⚙️ Definition of options

| Option | Description |
|---|---|
| `city` | Static location, where we pick weather. Leave empty for autodetect.|
| `location_provider_id` | `1`–`3`. See [here](#-location-providers). |
| `units` | `Metric` or `Imperial` |
| `diagram_hours_before` | Hours to display before current hour (max: 24). |
| `diagram_hours_after` | Hours to display after current hour (max: 48). |
| `windows_local_location_detection` | Runs a PowerShell command to get coordinates from current location. |
| `use_weather_cache` | Enables offline weather cache. |
| `single_line_output` | For automation, scripting, etc. |
| `hide_diagram` | Hides the temperature graph diagram; toggled with `S`. |
| `show_hints` | Shows on-screen hints; toggled with `H`. |
| `colorful_tui` | Enables colors in the terminal. |
| `bad_weather_notification` | Sends a notification when rain, snow or thunderstorm is expected. |
| `bad_weather_lookahead` | Hours ahead to check for bad weather before sending a notification. |
| `always_run_debugger` [DEV] | Always runs the debugger, even without flags. |

## 📍 Location providers
Default provider: **IPInfo** (`location_provider_id: 3`)

| Value | Provider |
|---|---|
| `1` | IpWho |
| `2` | IPApi |
| `3` | IPInfo |
