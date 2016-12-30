# Lambda functions
def plusN (n): return lambda m: m + n
plus5 = plusN(5)
plus10 = plusN(10)

print(plus5(3))
# >> 8

print(plus10(3))
# >> 13

print(list(filter(lambda x: x % 2 == 0, range(1, 10))))
# >> 2, 4, 6, 8

phrase = 'The quick brown fox jumps'
words = phrase.split()
print(words)
# >> ['The', 'quick', 'brown', 'fox', 'jumps']
wordLengths = map(lambda w: len(w), words)
print(list(wordLengths))
# >> [3, 5, 5, 3, 5]

