letrec (
  leaf = lambda () {
    lambda () { 0 }
  }
  node = lambda (value left right) {
    lambda () { 1 }
  }
  dfs = lambda (tree) {
    if (not (tree)) then (void)
    else [
      (dfs &left tree)
      (put &value tree "\n")
      (dfs &right tree)
    ]
  }
) {
  (dfs
    (node 4
      (node 2
        (node 1 (leaf) (leaf))
        (node 3 (leaf) (leaf)))
      (node 5 (leaf) (leaf))))
}
