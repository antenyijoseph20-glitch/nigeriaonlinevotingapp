#!/usr/bin/env python3

import json
import os
import sys
from datetime import datetime, timezone


RAW_FILE = "data/raw_lga_source.json"
RECONCILIATION_FILE = "data/lga_reconciliation.json"
OUTPUT_FILE = "data/lgas.json"


def fail(message):
    print(f"FAIL: {message}")
    sys.exit(1)


def load_json(path):
    if not os.path.exists(path):
        fail(f"Missing required file: {path}")

    try:
        with open(path, encoding="utf-8") as f:
            return json.load(f)
    except json.JSONDecodeError as exc:
        fail(f"Invalid JSON in {path}: {exc}")


def main():

    print("=" * 55)
    print(" Nigeria Online Voting System")
    print(" CANONICAL LGA DATASET BUILDER")
    print("=" * 55)

    # =========================================
    # Load source files
    # =========================================

    raw = load_json(RAW_FILE)
    reconciliation = load_json(RECONCILIATION_FILE)

    print()
    print("===== INPUT VALIDATION =====")

    # =========================================
    # Validate raw source count
    # =========================================

    if len(raw) != 774:
        fail(
            f"Expected 774 raw records, "
            f"found {len(raw)}"
        )

    print("Raw records: 774 PASS")

    # =========================================
    # Validate reconciliation count
    # =========================================

    if len(reconciliation) != 8:
        fail(
            f"Expected 8 reconciliation records, "
            f"found {len(reconciliation)}"
        )

    print("Reconciliation records: 8 PASS")

    # =========================================
    # Build reconciliation lookup
    # =========================================

    reconciliation_by_code = {}

    for row in reconciliation:

        code = row.get("source_code")

        if code is None:
            fail(
                "Reconciliation record is missing "
                "source_code"
            )

        if code in reconciliation_by_code:
            fail(
                f"Duplicate reconciliation code: {code}"
            )

        reconciliation_by_code[code] = row

    print("Reconciliation lookup: PASS")

    # =========================================
    # Validate raw source structure
    # =========================================

    required_raw_fields = [
        "code",
        "state_id",
        "state_code",
        "state_name",
        "source_name",
        "type",
        "source"
    ]

    for index, row in enumerate(raw, start=1):

        for field in required_raw_fields:

            if field not in row:
                fail(
                    f"Raw record {index} is missing "
                    f"field: {field}"
                )

        if not str(row["source_name"]).strip():
            fail(
                f"Raw record {index} has empty "
                "source_name"
            )

    print("Raw source structure: PASS")

    # =========================================
    # Validate raw source codes
    # =========================================

    raw_codes = [
        row["code"]
        for row in raw
    ]

    if len(raw_codes) != len(set(raw_codes)):
        fail(
            "Duplicate raw source codes detected"
        )

    print("Unique raw source codes: PASS")

    # =========================================
    # Prepare timestamp
    # =========================================

    now = datetime.now(
        timezone.utc
    ).replace(
        microsecond=0
    )

    timestamp = now.isoformat().replace(
        "+00:00",
        "Z"
    )

    # =========================================
    # Build canonical dataset
    # =========================================

    canonical = []

    for internal_id, raw_row in enumerate(
        raw,
        start=1
    ):

        source_code = raw_row["code"]

        # -------------------------------------
        # Default: preserve NBS source name
        # -------------------------------------

        canonical_name = raw_row[
            "source_name"
        ]

        verification_status = "source_only"

        verified_against = raw_row[
            "source"
        ]

        # -------------------------------------
        # Apply verified reconciliation
        # -------------------------------------

        if source_code in reconciliation_by_code:

            correction = reconciliation_by_code[
                source_code
            ]

            canonical_name = correction[
                "canonical_name"
            ]

            verification_status = correction[
                "verification_status"
            ]

            verified_against = correction[
                "verified_against"
            ]

        # -------------------------------------
        # Build final record
        # -------------------------------------

        record = {
            "id": internal_id,

            "code": raw_row["code"],

            "state_id": raw_row["state_id"],

            "state_code": raw_row["state_code"],

            "state_name": raw_row["state_name"],

            "source_name": raw_row["source_name"],

            "name": canonical_name,

            "type": raw_row["type"],

            "is_active": True,

            "verification_status": verification_status,

            "verified_against": verified_against,

            "source": raw_row["source"],

            "created_at": timestamp,

            "updated_at": timestamp
        }

        canonical.append(record)

    print(
        "Canonical records prepared:",
        len(canonical)
    )

    # =========================================
    # Validate canonical count
    # =========================================

    if len(canonical) != 774:
        fail(
            f"Canonical dataset contains "
            f"{len(canonical)} records"
        )

    # =========================================
    # Validate internal IDs
    # =========================================

    ids = [
        row["id"]
        for row in canonical
    ]

    expected_ids = list(
        range(1, 775)
    )

    if ids != expected_ids:
        fail(
            "Internal IDs are not sequential "
            "from 1 to 774"
        )

    print("Internal IDs 1-774: PASS")

    # =========================================
    # Validate source codes
    # =========================================

    codes = [
        row["code"]
        for row in canonical
    ]

    if len(codes) != len(set(codes)):
        fail(
            "Duplicate canonical source codes"
        )

    print("Unique source codes: PASS")

    # =========================================
    # Validate state mapping
    # =========================================

    for row in canonical:

        state_id = row["state_id"]

        if not isinstance(
            state_id,
            int
        ):
            fail(
                f'Invalid state_id type for '
                f'code {row["code"]}'
            )

        if state_id < 1 or state_id > 37:
            fail(
                f'Invalid state_id {state_id} '
                f'for code {row["code"]}'
            )

        if not row["state_code"].strip():
            fail(
                f'Empty state_code for '
                f'code {row["code"]}'
            )

        if not row["state_name"].strip():
            fail(
                f'Empty state_name for '
                f'code {row["code"]}'
            )

        if not row["name"].strip():
            fail(
                f'Empty canonical name for '
                f'code {row["code"]}'
            )

    print("State mapping: PASS")
    print("Canonical names: PASS")

    # =========================================
    # Validate LGA / FCT types
    # =========================================

    for row in canonical:

        if row["state_id"] == 37:

            if row["type"] != "AREA_COUNCIL":
                fail(
                    f'FCT record {row["code"]} '
                    f'must be AREA_COUNCIL'
                )

        else:

            if row["type"] != "LGA":
                fail(
                    f'Non-FCT record {row["code"]} '
                    f'must be LGA'
                )

    print("LGA/FCT type validation: PASS")

    # =========================================
    # Validate FCT
    # =========================================

    fct_records = [
        row
        for row in canonical
        if row["state_id"] == 37
    ]

    if len(fct_records) != 6:
        fail(
            f"Expected 6 FCT Area Councils, "
            f"found {len(fct_records)}"
        )

    print("FCT Area Councils: 6 PASS")

    # =========================================
    # Validate reconciliations
    # =========================================

    verified_records = [
        row
        for row in canonical
        if row["verification_status"]
        == "verified"
    ]

    if len(verified_records) != 8:
        fail(
            f"Expected 8 verified records, "
            f"found {len(verified_records)}"
        )

    print("Verified reconciliations: 8 PASS")

    # =========================================
    # Validate source-only records
    # =========================================

    source_only_records = [
        row
        for row in canonical
        if row["verification_status"]
        == "source_only"
    ]

    if len(source_only_records) != 766:
        fail(
            f"Expected 766 source-only records, "
            f"found {len(source_only_records)}"
        )

    print("Source-only records: 766 PASS")

    # =========================================
    # Confirm total verification categories
    # =========================================

    if (
        len(verified_records)
        + len(source_only_records)
        != 774
    ):
        fail(
            "Verification categories do not "
            "account for all 774 records"
        )

    print("Verification categories: PASS")

    # =========================================
    # Write canonical dataset
    # =========================================

    with open(
        OUTPUT_FILE,
        "w",
        encoding="utf-8"
    ) as f:

        json.dump(
            canonical,
            f,
            indent=4,
            ensure_ascii=False
        )

        f.write("\n")

    # =========================================
    # Final output
    # =========================================

    print()
    print("===== BUILD COMPLETE =====")
    print()
    print(
        "Canonical dataset:"
    )
    print(
        f"  {OUTPUT_FILE}"
    )
    print()
    print(
        "Records written:        774"
    )
    print(
        "Verified reconciliations: 8"
    )
    print(
        "Source-only records:    766"
    )
    print(
        "FCT Area Councils:        6"
    )
    print()
    print(
        "BUILD STATUS: PASS"
    )
    print()
    print(
        "The canonical dataset was generated"
    )
    print(
        "programmatically from the validated"
    )
    print(
        "NBS source and reconciliation layer."
    )
    print()


if __name__ == "__main__":
    main()
