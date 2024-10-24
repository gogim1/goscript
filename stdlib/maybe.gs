letrec (
  nothing = lambda () {
    lambda () { 0 }
  }
  just = lambda (value) {
    lambda () {1}
  }
  isnothing = lambda (maybe) {
    if (not (maybe)) then 1 else 0
  }
  isjust = lambda (maybe) {
    if (isnothing maybe) then 0 else 1
  }
  fromjust = lambda (maybe) {
    if (isjust maybe) then &value maybe else (void)
  }
  show = lambda (maybe) {
    if (isnothing maybe) then (put "Nothing") else (put "Just " &value maybe)
  }

  # functor. fmap :: (a -> b) -> f a -> f b
  fmap = lambda (m) {
    lambda (fa) {
      (just (m (fromjust fa)))
    }
  }

  # applicative. pure :: a -> f a
  pure = just

  # applicative. liftA2 :: (a -> b -> c) -> f a -> f b -> f c
  lift = lambda (m) {
    lambda (fa) {
      lambda (fb) {
        (just (m (fromjust fa) (fromjust fb)))
      }
    }
  }

  # monad. (>>=) :: m a -> (a -> m b) -> m b
  compose = lambda (fa) {
    lambda (m) {
      (m (fromjust fa))
    }
  }
) {[
  (reg "nothing" nothing)
  (reg "just" just)
  (reg "isnothing" isnothing)
  (reg "isjust" isjust)
  (reg "fromjust" fromjust)
]}