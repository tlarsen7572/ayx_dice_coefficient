# ayx_dice_coefficient
This repository contains two, independent, Alteryx macros for calculating dice coefficient.

There are two versions of the macros:
* DiceCoefficient.yxmc uses the Python tool to implement the third Python algorithm on [this page](https://en.wikibooks.org/wiki/Algorithm_Implementation/Strings/Dice%27s_coefficient#Python)
* DiceCoefficientAyxOnly.yxmc implements the same algorithm using only standard Alteryx tools.

The Alteryx-only implementation runs faster because it does not have the overhead of the Python runtime. Also, I am unsure how the Python tool will memory-manage large data sources. The Alteryx engine manages memory and disk caching in the standard tools, so the Alteryx-only implementation should be safe. It is possible the Python tool will request as much memory as needed to keep the incoming data in memory and could cause out-of-memory errors.

The only downside to the Alteryx-only implementation is that I find it more difficult to follow the logic. To my eyes, Python more clearly expresses the steps needed to calculate the dice coefficient.
