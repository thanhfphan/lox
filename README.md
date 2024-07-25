# Lox

Go implementation of Lox interpreter (based on the book "Crafting Interpreters")

## Getting started
```
fun makeCounter() {
  var i = 0;
  fun count() {
    i = i + 1;
    print i;
  }

  return count;
}

var counter = makeCounter();
counter(); // "1".
counter(); // "2".

print clock(); // built-in function
```
