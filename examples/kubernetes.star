def has_sig_api_machinery(item):
  labels = item.get("labels")
  if "sig/api-machinery" in labels:
    return True
  return False

github(org="kubernetes", repo="kubernetes", filters=[has_sig_api_machinery])
