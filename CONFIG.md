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
| `template` | Template for your fullscreen output. [See here](#-customize). |
| `single_line_template` | Same as `template`, but for single line. Also [read here](#-customize). |
| `colors` | Colors used in the TUI. |
| `always_run_debugger` [DEV] | Always runs the debugger, even without flags. |

## 📍 Location providers
Default provider: **IPInfo** (`location_provider_id: 3`)

| Value | Provider |
|---|---|
| `1` | IpWho |
| `2` | IPApi |
| `3` | IPInfo |

## 🎨 Customize
You can customize output of skypaw by using `text/template`. Read docs [here](https://pkg.go.dev/text/template).

### 🧩 Functions
The following template functions are available:

| Function | Example | Description |
| --- | --- | --- |
| `theme` | {{.Temp \| theme "temp"}} | Color from the [colors](#-colors) section. |
| `color` | {{.City \| color "#89b4fa"}} |Arbitrary hex color. |
| `bold` | {{.Temp \| bold}} | Makes text bold. |
| `upper`, `lower` | {{.City \| upper}} | Changes case. |
| `round` | {{round .TempValue}} | Numbers without decimals.

- Put `upper` and `lower` before `theme`, `color` or `bold`, because color codes will be corrupted.
- With `colorful_tui: false` the following variables will return plain text: `theme`, `color`, `bold`.
- Go functions also available. For more details, see [text/template](https://pkg.go.dev/text/template).

### 📋 Fields
| Field | Description | Example |
| --- | --- | --- |
| `.City` | City | `London` |
| `.Region` | Region | `England` |
| `.Location` | Region, City | `England, London` |
| `.Temp` | Temperature with unit | `19.1°C` |
| `.TempValue` | Temperature as a number | `19.1` |
| `.Condition` | Weather description | `Partly Cloudy` |
| `.Time` | Time | `15:04` |
| `.TimeSec` | Time with seconds | `15:04:05` |
| `.Date` | Date | `18.11.2026` |
| `.DateLong` | Weekday, month, day | `Wednesday, November 18` |
| `.Weekday` | Weekday | `Wednesday` |
| `.Offline` | `true`, if data was loaded from cache. | |
| `.OfflineNotice` | Offline message with last update time (empty when online). | |
| `.Art` | ASCII art (fullscreen only). | |
| `.Hints` | Key hints (fullscreen only, empty if `show_hints` is off). | |

### 🌈 Colors

You can customize the colors used in the TUI. Default values are provided below:

| Key | Value |
|---|---|
| `location` | `#40E0D0` ![](https://placehold.co/16x16/40E0D0/40E0D0.png) |
| `temp` | `#FFD700` ![](https://placehold.co/16x16/FFD700/FFD700.png) |
| `weather_name` | `#9bcbff` ![](https://placehold.co/16x16/9bcbff/9bcbff.png) |
| `time` | `#5FFFD7` ![](https://placehold.co/16x16/5FFFD7/5FFFD7.png) |
| `date` | `#98FB98` ![](https://placehold.co/16x16/98FB98/98FB98.png) |
| `hints` | `#4D4D4D` ![](https://placehold.co/16x16/4D4D4D/4D4D4D.png) |
| `error` | `#ff0000` ![](https://placehold.co/16x16/ff0000/ff0000.png) |
| `loading` | `#ff7ad7` ![](https://placehold.co/16x16/ff7ad7/ff7ad7.png) |
| `offline_or_timeout` | `#3f3f3f` ![](https://placehold.co/16x16/3f3f3f/3f3f3f.png) |
| `weather_clear` | `#FFD700` ![](https://placehold.co/16x16/FFD700/FFD700.png) |
| `weather_partly_cloudy` | `#B0C4DE` ![](https://placehold.co/16x16/B0C4DE/B0C4DE.png) |
| `weather_overcast` | `#778899` ![](https://placehold.co/16x16/778899/778899.png) |
| `weather_fog` | `#9DB4C0` ![](https://placehold.co/16x16/9DB4C0/9DB4C0.png) |
| `weather_drizzle` | `#87CEEB` ![](https://placehold.co/16x16/87CEEB/87CEEB.png) |
| `weather_freezing_drizzle` | `#B0E0E6` ![](https://placehold.co/16x16/B0E0E6/B0E0E6.png) |
| `weather_rain` | `#4169E1` ![](https://placehold.co/16x16/4169E1/4169E1.png) |
| `weather_freezing_rain` | `#00CED1` ![](https://placehold.co/16x16/00CED1/00CED1.png) |
| `weather_snowfall` | `#E0FFFF` ![](https://placehold.co/16x16/E0FFFF/E0FFFF.png) |
| `weather_snow_grains` | `#CAE9FF` ![](https://placehold.co/16x16/CAE9FF/CAE9FF.png) |
| `weather_rain_showers` | `#1E90FF` ![](https://placehold.co/16x16/1E90FF/1E90FF.png) |
| `weather_snow_showers` | `#B0E0E6` ![](https://placehold.co/16x16/B0E0E6/B0E0E6.png) |
| `weather_thunderstorm` | `#9400D3` ![](https://placehold.co/16x16/9400D3/9400D3.png) |
| `weather_thunderstorm_hail` | `#FF6347` ![](https://placehold.co/16x16/FF6347/FF6347.png) |
| `weather_unknown` | `#DA70D6` ![](https://placehold.co/16x16/DA70D6/DA70D6.png) |
