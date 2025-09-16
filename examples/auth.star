# Fetch all open issues and pull requests
# in everettraven/masquerade, a private repository.
masquerade_open_items = github.search(query="repo:everettraven/masquerade is:open")

# Tell wranglr to include the items in the output
wranglr.render(masquerade_open_items)
