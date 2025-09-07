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

# Example of setting the status of an item based
# on fields of the item.
# In this case we want:
# - IF item is an Issue AND does not have label triage/accepted THEN status = "Needs Triage"
# - IF item is a PullRequest AND does not have any assignees AND is missing the LGTM label THEN status = "Needs Review"
def set_status(item):
  type = item.get("type")
  labels = item.get("labels")
  assignees = item.get("assignees")

  if type == "Issue":
    if "triage/accepted" not in labels:
      return "Needs Triage"
  if type == "PullRequest":
    if len(assignees) == 0 and "lgtm" not in labels:
      return "Needs Review"
  return ""

github(repo="kubernetes-sigs/kube-api-linter", filters=[authored_by_not_me], priorities=[priority_help_wanted], status=set_status)
github(repo="kubernetes-sigs/crdify", filters=[authored_by_not_me])
github(repo="kubernetes-sigs/crdify", mentioned="everettraven")

