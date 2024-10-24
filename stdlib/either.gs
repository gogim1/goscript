letrec (
  left = lambda (value) {
    lambda () { 0 }
  }
  right = lambda (value) {
    lambda () {1}
  }
  isleft = lambda (either) {
    if (not (either)) then 1 else 0
  }
  isright = lambda (either) {
    if (isleft either) then 0 else 1
  }
  fromleft = lambda (either) {
    if (isleft either) then &value either else (void)
  }
  fromright = lambda (either) {
    if (isright either) then &value either else (void)
  }
) {[
  (reg "left" left)
  (reg "right" right)
  (reg "isleft" isleft)
  (reg "isright" isright)
  (reg "fromleft" fromleft)
  (reg "fromright" fromright)
]}