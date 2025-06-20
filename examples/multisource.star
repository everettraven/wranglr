def authored_by_me(item):
  author = item.get("author")
  if author != "everettraven":
    return False
  return True

github(org="kubernetes-sigs", repo="kube-api-linter", filters=[authored_by_me])
github(org="kubernetes-sigs", repo="crdify", filters=[authored_by_me])
