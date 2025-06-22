# Example of filtering out items where the author
# is _not_ 'everettraven'
def authored_by_not_me(item):
  author = item.get("author")
  if author != "everettraven":
    return True
  return False

# Example of increasing priority score by
# a value of '50' if the item contains the
# label 'help wanted'
def priority_help_wanted(item):
  labels = item.get("labels")
  if "help wanted" in labels:
    return 50
  return 0

github(org="kubernetes-sigs", repo="kube-api-linter", filters=[authored_by_not_me], priorities=[priority_help_wanted])
github(org="kubernetes-sigs", repo="crdify", filters=[authored_by_not_me])
