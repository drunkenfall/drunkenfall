function isGoZeroDateOrFalsy (m) { // moment date
  if (m) {
    return m.isSame("0001-01-01T00:00:00Z")
  } else {
    return true
  }
}

export { isGoZeroDateOrFalsy }
