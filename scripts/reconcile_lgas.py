#!/usr/bin/env python3

import json
import os
import re
import sys


RAW_FILE = "data/raw_lga_source.json"


def normalize_name(name):
    """
    Normalize a name only for comparison.

    This does NOT change the stored source name.
    """
    name = name.strip().upper()
    name = re.sub(r"\s+", " ", name)
    return name


def main():
    print("=" * 55)
    print(" Nigeria Online Voting System")
    print(" LGA SOURCE RECONCILIATION REPORT")
    print("=" * 55)

    if not os.path.exists(RAW_FILE):
        print()
        print("ERROR: Raw NBS LGA source does not exist.")
        print(f"Expected: {RAW_FILE}")
        sys.exit(1)

    with open(RAW_FILE, "r", encoding="utf-8") as f:
        records = json.load(f)

    print()
    print("===== RAW SOURCE =====")
    print(f"Records loaded: {len(records)}")

    if len(records) != 774:
        print("FAIL: Expected exactly 774 records.")
        sys.exit(1)

    print("774 records: PASS")

    print()
    print("===== POTENTIAL NAME RECONCILIATIONS =====")

    # These are NOT corrections.
    #
    # They are records that require verification against
    # an authoritative/current administrative source.
    #
    # The NBS source_name remains unchanged.
    review_codes = {
        713,
        801,
        905,
        1023,
        1024,
        1025,
        1807,
        2415,
    }

    found = 0

    for record in records:
        code = record["code"]

        if code in review_codes:
            found += 1

            print(
                f'{code} | '
                f'{record["state_name"]} | '
                f'{record["source_name"]} | '
                f'{record["type"]}'
            )

    print()
    print(f"Records requiring review: {found}")

    if found != len(review_codes):
        print("FAIL: Review-code count mismatch.")
        sys.exit(1)

    print("Review list integrity: PASS")

    print()
    print("===== SOURCE NAME PRESERVATION =====")

    missing_names = []

    for record in records:
        source_name = record.get("source_name")

        if not isinstance(source_name, str):
            missing_names.append(record["code"])
            continue

        if not source_name.strip():
            missing_names.append(record["code"])

    if missing_names:
        print("FAIL: Missing source names:")
        print(missing_names)
        sys.exit(1)

    print("All source names present: PASS")

    print()
    print("===== IMPORTANT =====")
    print()
    print("No LGA name has been automatically changed.")
    print("The NBS source names remain the source of record.")
    print()
    print("Next step:")
    print("Verify the eight flagged records against")
    print("an independent authoritative/current source.")
    print()
    print("Final lgas.json has NOT been created.")

    print()
    print("RECONCILIATION REPORT COMPLETE.")


if __name__ == "__main__":
    main()
