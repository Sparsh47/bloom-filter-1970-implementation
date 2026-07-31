# Bloom Filter 1970 Implementation - Validated against its theoretical false-positive formula

### Motivation
Normal hash set stores data completely and so consumes more memory. Bloom filter stores only a few bits per element thereby reducing memory usage by allowing a few false positives.

### How it works

### Theoretical Formula
__False Positive Probability__

```math
p = \left(1 - e^{-\frac{kn}{m}}\right)^k
```

* p = probability of false positive
* k = number of hash functions
* m = size of the bit array
* n = number of elements in filter

### Implementation Details
* Language: GO
* Hash Library Used: [MurmurHash](https://github.com/spaolacci/murmur3)
* OS: Unix/Linux

### Experimental Setup
* I made my words test set using the existing words directory present in Unix/Linux systems (/usr/share/dict/words)
* I tested the implementation using 4 variations: 

| Config | m (bits) | n (words inserted) | k (hash functions) | m/n ratio |
|--------|----------|---------------------|----------------------|-----------|
| A      | 300,000  | 50,000              | 3                    | 6         |
| B      | 300,000  | 50,000              | 5                    | 6         |
| C      | 300,000  | 100,000             | 5                    | 3         |
| D      | 500,000  | 100,000             | 3                    | 5         |

### Results
| Config | m (bits) | n (words) | k | m/n ratio | Predicted rate | Measured rate | Gap |
|--------|---------:|----------:|--:|----------:|----------------:|----------------:|-----:|
| A      | 300,000  | 50,000    | 3 | 6         | 6.09%           | 5.78%           | 0.31 pts |
| B      | 300,000  | 50,000    | 5 | 6         | 5.78%           | 5.16%           | 0.62 pts |
| C      | 300,000  | 100,000   | 5 | 3         | 35.13%          | 34.71%          | 0.42 pts |
| D      | 500,000  | 100,000   | 3 | 5         | 9.18%           | 9.04%           | 0.14 pts |

![False positive rate vs. number of hash functions, theoretical curves with measured points overlaid](https://res.cloudinary.com/deuj8faqq/image/upload/v1785502823/Screenshot_2026-07-31_at_4.50.09_PM_newnbr.png)

### Discussion

The measured false-positive rates track the theoretical formula closely across
all four configurations, with gaps ranging from 0.14 to 0.62 percentage points.
These small deviations are expected — MurmurHash approximates but doesn't
perfectly achieve the independence assumption the formula relies on, and any
finite test sample carries some statistical noise.

Two trends emerge from the data:

**1. False-positive rate depends on the m/n ratio, not raw scale.**
Config D (500,000 bits / 100,000 words) has a higher predicted and measured
rate than Config A (300,000 bits / 50,000 words), even though its arrays are
larger in absolute terms. What matters is bits allocated per inserted element
(m/n) — D has 5 bits/word, A has 6 — and the formula and measurements both
reflect that lower ratio as a worse rate.

**2. There is an optimal k for a given m/n ratio, and overshooting it hurts.**
At m/n=6 (Configs A and B), increasing k from 3 to 5 *improved* the rate
(5.78% → 5.16%), consistent with the optimal k for this ratio
(k ≈ ln(2) × m/n ≈ 4.16) sitting closer to 5 than 3.
At m/n=3 (Config C), the same k=5 produces a much worse rate (34.71%) because
the optimal k for this ratio is only ≈2.08 — five hash functions oversaturate
a bit array that small, setting far more bits than necessary and driving up
collisions rather than reducing them.

This confirms the filter isn't just numerically close to theory on a single
config, but reproduces the actual shape of the tradeoff curve Bloom's analysis
predicts.

### How to Run

```bash
go run main.go
```

The program loads the word list, builds the filter, runs the false-positive
test against the dictionary word set, and prints the bit array size, false
positive count, and measured rate.

### References

* Bloom, B. H. (1970). [Space/Time Trade-offs in Hash Coding with Allowable Errors](https://dl.acm.org/doi/10.1145/362686.362692). *Communications of the ACM*, 13(7), 422–426.
* [MurmurHash / go-murmur3](https://github.com/spaolacci/murmur3)
