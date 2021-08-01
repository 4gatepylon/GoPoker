# Docs

## Card Sets
Everything is built from the Carset type. 

The cardset type is implemented as a `uint64` bit-vector ordered by `(val, suit)` from right to left. That is to say, `Two of Clubs = 2C = b'...000010000`, `Three of Clubs = 3C = b'...100000000`, `Two of Diamonds = 2D = b'...000100000`. Aces are unique in that they are both the highest and lowest card so in effect an ace is two cards: a little ace (i.e. 1, 4, 8, or 16) and a big ace (bit shift the little ace by 13 to the left). This is purely for practical reasons to calculate ladders more easily. It turned out that it made everything else slightly more convoluted, so I think it may have been a mistake but whatever. Since Aces are two cards per suit and there are another twelve cards per suit, there are a total of `(12 + 2) * 4 < 64` (quick maths). This leaves us with `8` extra cards, which is 2 per suit if we like. Why is this good? It means we can have jokers and other nonsense for example. Basically we can mod the game. 

To those of you who cannot infer from example, this is the structure, split by nice spaces and with parens added to tag bit quadruplets by their values: `b'0000 0000 0000(A) 0000(K) 0000(Q) 0000(J) 0000(T) 0000(9) 0000(8) 0000(7) 0000(6) 0000(5) 0000(4) 0000(3) 0000(2) 0000(A)`. Every quadruplet `0000` turns into `Spades Hearts Diamonds Clubs` so to bitshift one to the left is to go up a suit (i.e. clubs to diamonds) unless you are at spades in which case you go to clubs of the next card (i.e. Two of Spades to Three of Clubs). Note 10 is actually T because it made the code slightly easier.


I was too lazy to create an interface for this. However, Golang is luckily friendly to this sort of anti-pattern I have engaged in. You can create an interface if you like, simply replace the name "CardSet" with something like "CardSetBitVector64" and then change CardSet to be an interface with the relevant methods (check cardset.go) that you need. Then CardSetBitVector64 can implement CardSet and you can take keep the rest of the code (game and hands) the same since they take card sets which use those methods. You may have to read to see which methods are necessary. Good luck!

## Games
A state-machine like interface. Will keep different rooms in files using temp directories. Json is the commonly used format in files.