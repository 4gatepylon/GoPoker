# Construction

## Poker
### Card Sets
Everything is built from the Carset type. 

The cardset type is implemented as a `uint64` bit-vector ordered by `(val, suit)` from right to left. That is to say, `Two of Clubs = 2C = b'...000010000`, `Three of Clubs = 3C = b'...100000000`, `Two of Diamonds = 2D = b'...000100000`. Aces are unique in that they are both the highest and lowest card so in effect an ace is two cards: a little ace (i.e. 1, 4, 8, or 16) and a big ace (bit shift the little ace by 13 to the left). This is purely for practical reasons to calculate ladders more easily. It turned out that it made everything else slightly more convoluted, so I think it may have been a mistake but whatever. Since Aces are two cards per suit and there are another twelve cards per suit, there are a total of `(12 + 2) * 4 < 64` (quick maths). This leaves us with `8` extra cards, which is 2 per suit if we like. Why is this good? It means we can have jokers and other nonsense for example. Basically we can mod the game. 

To those of you who cannot infer from example, this is the structure, split by nice spaces and with parens added to tag bit quadruplets by their values: `b'0000 0000 0000(A) 0000(K) 0000(Q) 0000(J) 0000(T) 0000(9) 0000(8) 0000(7) 0000(6) 0000(5) 0000(4) 0000(3) 0000(2) 0000(A)`. Every quadruplet `0000` turns into `Spades Hearts Diamonds Clubs` so to bitshift one to the left is to go up a suit (i.e. clubs to diamonds) unless you are at spades in which case you go to clubs of the next card (i.e. Two of Spades to Three of Clubs). Note 10 is actually T because it made the code slightly easier.


I was too lazy to create an interface for this. However, Golang is luckily friendly to this sort of anti-pattern I have engaged in. You can create an interface if you like, simply replace the name "CardSet" with something like "CardSetBitVector64" and then change CardSet to be an interface with the relevant methods (check cardset.go) that you need. Then CardSetBitVector64 can implement CardSet and you can take keep the rest of the code (game and hands) the same since they take card sets which use those methods. You may have to read to see which methods are necessary. Good luck!

### Hands and Games
Hands and games give us a state-machine like interface to interact with a game. Games have hands (one per person and one in the middle). (Surprise! hands are card sets, but with additional functionality.) I haven't built this yet, but I intend to make this interfacy so it can be changed.

## Main and elsewhere
### Server and Client
Not yet deciced. Going to use gRPC and play from the commandline. I'll probably create some proto3 protocol buffers ("protobufs" is the usual lingo).

## Building and Running
Tbd. I'll probably change this to use bazel and right a simple skylark rule to run go tests. Also I'd like to change the gopath or goroot so that it will access the bazel repo instead of my host go so that you can run it anywhere easily and without annoying spaghetti.

## Testing
Every file more or less defines a logical grouping of things. For example: a game. Each logical grouping of things is unit tested. If the tests are not pretty or whatever, bad luck bears. I don't care. These tests should have decent coverage of black-box codepaths (that is to say, the input output pairs) such that both sides are covered, even if their product isn't necessarily covered (because ain't nobody got time for that).

Running specific tests:
- https://stackoverflow.com/questions/48465080/how-do-i-run-specific-golang-test-using-go-test-run

# Why
## Did I do this?
To learn golang, gRPC, and Bazel (and here I am not using gRPC nor Bazel: "Sad!" -ex president of the united states). Also to not atrophy.

## Use bitvectors?
Makes me feel like a real programmer. It's also like those vegan people who say they are vegan to help the environment or whatever, except I'm helping my server's memory consumption: it helps me sleep at night. It' also O(1) everything, including random cards/shuffling (just generate a random uint64 you fool). Honestly, I'm inspired to do this in the future. One downside is you can't have multiple copies of each card. We'd have to create some sort of bit-matrix.

## Is it written so shittily?
Notice how the later commits have less shitty code. Notice the trajectory? It was just warmup, not the real game yet.

## Are the tests written so monolothically sometimes?
Sometimes intellectual laziness is actually inferior to temporal laziness.

## ... Just why?
Why not?

# Notes to self from before because I'm a hoarder about information (Whoa will I forget? uh oh ruh roh!)
sudo go get github.com/4gatepylon/GoPoker 

All go commands must be sudoed.

Golang is smart and every file that is Test|string that doesn't start with a lowercase letter|
will be run when you do "go test" and we can use fmt.Errof to throw errors

I'll probably move to bazel later? Not sure... Maybe that'll be for a future project. I don't think I'll even get to SDL on this one, since I'm planning on making it from the lowest level up.