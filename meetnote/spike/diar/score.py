#!/usr/bin/env python3
"""Scores a diarization against the reference turns.

    python3 score.py out/dialogue.ref.tsv out/hyp.tsv

Hypothesis speaker ids are mapped to reference names by the mapping that agrees most (brute force).
- time error: share of reference speech time labelled with the wrong speaker or with nobody
- turn accuracy: share of reference turns whose dominant hypothesis speaker is the right one,
  which is what labelling each transcript line would get
"""
import itertools
import sys

STEP = 0.01


def load(path):
    return [(float(a), float(b), c) for a, b, c in (l.rstrip("\n").split("\t") for l in open(path) if l.strip())]


def frames(segments, end):
    labels = [None] * int(end / STEP + 1)
    for s, e, who in segments:
        for i in range(int(s / STEP), min(len(labels), int(e / STEP))):
            labels[i] = who
    return labels


def main():
    ref, hyp = load(sys.argv[1]), load(sys.argv[2])
    end = max(e for _, e, _ in ref + hyp)
    r, h = frames(ref, end), frames(hyp, end)
    names = sorted({c for *_, c in ref})
    ids = sorted({c for *_, c in hyp})

    def agreement(mapping):
        return sum(1 for a, b in zip(r, h) if a is not None and b is not None and mapping.get(b) == a)

    best = {}
    best_score = -1
    # Map each hypothesis id to a distinct reference name (extra ids stay unmapped).
    for perm in itertools.permutations(names + [None] * max(0, len(ids) - len(names)), len(ids)):
        mapping = dict(zip(ids, perm))
        score = agreement(mapping)
        if score > best_score:
            best, best_score = mapping, score

    speech = sum(1 for a in r if a is not None)
    time_error = 1 - best_score / speech

    right = 0
    details = []
    for s, e, who in ref:
        votes = {}
        for i in range(int(s / STEP), int(e / STEP)):
            if h[i] is not None:
                votes[h[i]] = votes.get(h[i], 0) + 1
        guess = best.get(max(votes, key=votes.get)) if votes else None
        right += guess == who
        details.append(f"  {s:6.2f} {who:8} -> {guess or '-':8} {'' if guess == who else '✗'}")

    print(f"speakers found: {len(ids)} (actual {len(names)}), mapping {best}")
    print(f"time error {time_error:.1%}, turn accuracy {right}/{len(ref)} = {right / len(ref):.0%}")
    if "-v" in sys.argv:
        print("\n".join(details))


if __name__ == "__main__":
    main()
