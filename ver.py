#!/usr/bin/env python3

import argparse
from datetime import datetime

def v_today():
    now = datetime.now()
    return now.strftime("v%Y%m%d")

parser = argparse.ArgumentParser()
parser.add_argument("-i", "--inc", help="updates VERSION.txt", action="store_true")
parser.add_argument("-n", "--dry_run", help="prints what would happen", action="store_true")
args = parser.parse_args()

with open("VERSION.txt") as f:
    try:
        ver, date = f.readline().strip().split("_")
        if args.inc:
            # increment ver
            major, minor, patch = map(int, ver.split("."))
            patch += 1
            ver = ".".join(map(str, [major, minor, patch]))
            # increment date
            date = v_today()
    except ValueError:
        ver, date = ("0.0.1", v_today())

verstr = f"{ver}_{date}\n"
if args.inc and not args.dry_run:
    with open("VERSION.txt", "w") as f:
        f.write(verstr)
        f.flush()
    with open("VERSION.txt") as f:
        for line in f.readlines():
            print(line, end="")
else:
    print(verstr, end="")
            
