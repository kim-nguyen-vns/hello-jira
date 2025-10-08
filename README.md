# Hello Jira

Hello Jira is a tool that automatically pulls events from `.ics` calendar url and logs work entries to Jira.

## Features

- Parse `.ics` url for calendar events
- Map events to Jira issues
- Log work entries to Jira automatically

## Installation

1. **Install with Go:**
    ```bash
    go install github.com/kim-nguyen-vns/hello-jira@latest
    ```

## Usage

1. **Prepare your `.env` file:** refer to [`.env.example`](./.env.example).

3. **Sample command:**
    ```bash
    hello-jira --env /path/to/.env --date yyyy-mm-dd  
    ```

