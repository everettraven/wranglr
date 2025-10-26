# Interactive Output

The `interactive` output format (`--output interactive`, which is the default) renders items in a TUI.

The general structure of the interactive output is:
- Vertical tabs to represent groupings of items (i.e "Feature X", "API Review", "SIG Auth", etc.).
  - If you do not specify groups, or only have a singular group, no vertical tabs will be present.
- Horizontal tabs to represent the current status of items (i.e "Todo", "Needs Review", etc.).
  - If you do not specify statuses, or only have a singular status, no horizontal tabs will be present.
- A paginated set of "pages" for each item in the currently selected group and status. A single page is displayed at a time.

## Keybindings

### Navigating groups

- `J` - go to next vertical tab
- `K` - go to previous vertical tab

### Navigating statuses

- `H` - go to previous horizontal tab
- `L` - go to next horizontal tab

### Navigating items

- `h` - go to previous item
- `l` - go to next item
- `j` - scroll down
- `k` - scroll up

### Actions

- `o` - open item in your browser

### Quitting

- `q` - quits the interactive view
