# Fetch all open issues and pull requests
# in the kubernetes-sigs/crdify repository
crdify_open_items = github.search(query="repo:kubernetes-sigs/crdify is:open")

# Tell wranglr to include them in the rendered output
wranglr.render(crdify_open_items)
