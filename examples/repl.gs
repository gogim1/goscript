letrec (
  prompt = [
    (put "  ___  _____  ___   ___  ____  ____  ____  ____  " "\n")
    (put " / __)(  _  )/ __) / __)(  _ \\(_  _)(  _ \\(_  _) " "\n")
    (put "( (_-. )(_)( \\__ \\( (__  )   / _)(_  )___/  )(    " "\n")
    (put " \\___/(_____)(___/ \\___)(_)\\_)(____)(__)   (__)" " Hopes you have a `SUGOI` day!\n")
  ]
  loop = lambda () {
    letrec (
      line = [(put "> ") (getline)]
    ) {
      [(put (eval line) "\n") (loop)]
    }
  }
) {
  (loop)
}