#!/usr/bin/env python3

import json
import sys


RAW_FILE = "data/raw_lga_source.json"
RECONCILIATION_FILE = "data/lga_reconciliation.json"


def fail(message):
    print(f"FAIL: {message}")
    sys.exit(1)


def main():
    print("=" * 55)
    print(" Nigeria Online Voting System")
    print(" LGA RECONCILIATION VALIDATOR")
    print("=" * 55)

    # -----------------------------------------
    # Load raw source
    # -----------------------------------------

    try:
        with open(RAW_FILE, encoding="utf-8") as f:
            raw = json.load(f)
    except FileNotFoundError:
        fail(f"Missing {RAW_FILE}")

    # -----------------------------------------
    # Load reconciliation data
    # -----------------------------------------

    try:
        with open(RECONCILIATION_FILE, encoding="utf-8") as f:
            reconciliation = json.load(f)
    except FileNotFoundError:
        fail(f"Missing {RECONCILIATION_FILE}")

    print()
    print("===== RAW SOURCE =====")

    if len(raw) != 774:
        fail(f"Expected 774 raw records, found {len(raw)}")

    print("774 raw records: PASS")

    # -----------------------------------------
    # Validate raw source codes
    # -----------------------------------------

    raw_codes = [row["code"] for row in raw]

    if len(raw_codes) != len(set(raw_codes)):
        fail("Duplicate raw source codes detected")

    print("Unique raw source codes: PASS")

    # -----------------------------------------
    # Validate reconciliation count
    # -----------------------------------------

    print()
    print("===== RECONCILIATION =====")

    if len(reconciliation) != 8:
        fail(
            f"Expected 8 reconciliation records, "
            f"found {len(reconciliation)}"
        )

    print("8 reconciliation records: PASS")

    # -----------------------------------------
    # Validate reconciliation codes
    # -----------------------------------------

    reconciliation_codes = [
        row["source_code"]
        for row in reconciliation
    ]

    if len(reconciliation_codes) != len(
        set(reconciliation_codes)
    ):
        fail("Duplicate reconciliation source codes")

    print("Unique reconciliation codes: PASS")

    # -----------------------------------------
    # Confirm every reconciliation record
    # matches a raw NBS record
    # -----------------------------------------

    raw_by_code = {
        row["code"]: row
        for row in raw
    }

    for row in reconciliation:

        code = row["source_code"]

        if code not in raw_by_code:
            fail(
                f"Reconciliation code {code} "
                f"does not exist in raw source"
            )

        source = raw_by_code[code]

        if source["source_name"] != row["source_name"]:
            fail(
                f"Source name mismatch for code {code}: "
                f'{source["source_name"]} != '
                f'{row["source_name"]}'
            )

        if source["state_code"] != row["state_code"]:
            fail(
                f"State code mismatch for code {code}"
            )

        if source["type"] != row["type"]:
            fail(
                f"Type mismatch for code {code}"
            )

        if row["verification_status"] != "verified":
            fail(
                f"Record {code} is not marked verified"
            )

        if not row["canonical_name"].strip():
            fail(
                f"Record {code} has no canonical name"
            )

        if not row["verified_against"].strip():
            fail(
                f"Record {code} has no verification source"
            )

    print("All reconciliation records match raw source: PASS")
    print("Source names preserved: PASS")
    print("State mappings preserved: PASS")
    print("Types preserved: PASS")
    print("Canonical names present: PASS")
    print("Verification status present: PASS")

    # -----------------------------------------
    # Confirm FCT structure
    # -----------------------------------------

    print()
    print("===== FCT CHECK =====")

    fct_records = [
        row for row in raw
        if row["state_id"] == 37
    ]

    if len(fct_records) != 6:
        fail(
            f"Expected 6 FCT Area Councils, "
            f"found {len(fct_records)}"
        )

    for row in fct_records:
        if row["type"] != "AREA_COUNCIL":
            fail(
                f'FCT record {row["code"]} '
                f'is not AREA_COUNCIL'
            )

    print("6 FCT Area Councils: PASS")
    print("FCT type validation: PASS")

    # -----------------------------------------
    # Final count
    # -----------------------------------------

    print()
    print("===== FINAL RECONCILIATION STATUS =====")
    print()
    print("Raw records:            774")
    print("Reconciled records:       8")
    print("Unchanged records:       766")
    print("FCT Area Councils:         6")
    print()
    print("RECONCILIATION VALIDATION: PASS")
    print()
    print("The raw NBS dataset is safe to transform")
    print("into the canonical dataset.")
    print()
    print("NEXT STEP:")
    print("Build data/lgas.json from the validated")
    print("raw source + reconciliation layer.")
    print()


if __name__ == "__main__":
    main()
