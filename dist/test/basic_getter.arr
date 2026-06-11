def fn() {
  def a = 69
  defer free a
  ret 696969 + a
}

def tttt = fn() + 1
