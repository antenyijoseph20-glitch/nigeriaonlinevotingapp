#!/usr/bin/env python3

import json
import sys
from pathlib import Path


# ============================================================
# Configuration
# ============================================================

PROJECT_ROOT = Path(__file__).resolve().parent.parent

STATES_FILE = PROJECT_ROOT / "data" / "states.json"
LGAS_FILE = PROJECT_ROOT / "data" / "lgas.json"

EXPECTED_STATE_COUNT = 37
EXPECTED_NATIONAL_LGA_COUNT = 774

VALID_LGA_TYPE = "LGA"
VALID_AREA_COUNCIL_TYPE = "AREA_COUNCIL"


# ============================================================
# Helpers
# ============================================================

def fail(message):
    print(f"FAIL: {message}")
    sys.exit(1)


def load_json(path):
    if not path.exists():
        fail(f"File does not exist: {path}")

    try:
        with path.open("r", encoding="utf-8") as file:
            return json.load(file)

    except json.JSONDecodeError as error:
        fail(
            f"Invalid JSON in {path}: "
            f"line {error.lineno}, column {error.colno}"
        )

    except OSError as error:
        fail(f"Unable to read {path}: {error}")


def normalize(value):
    return " ".join(str(value).strip().lower().split())


# ============================================================
# Validate states
# ============================================================

def validate_states(states):
    print("===== VALIDATING STATES =====")

    if not isinstance(states, list):
        fail("states.json must contain a JSON array")

    if len(states) != EXPECTED_STATE_COUNT:
        fail(
            f"Expected {EXPECTED_STATE_COUNT} state/FCT records, "
            f"found {len(states)}"
        )

    state_ids = set()
    state_codes = set()
    state_names = set()

    fct_count = 0
    normal_state_count = 0

    for index, state in enumerate(states, start=1):

        if not isinstance(state, dict):
            fail(f"State record #{index} is not an object")

        required_fields = [
            "id",
            "name",
            "code",
            "is_fct",
            "is_active",
        ]

        for field in required_fields:
            if field not in state:
                fail(
                    f"State record #{index} is missing field '{field}'"
                )

        state_id = state["id"]
        state_name = state["name"]
        state_code = state["code"]
        is_fct = state["is_fct"]

        if not isinstance(state_id, int):
            fail(
                f"State '{state_name}' has a non-integer ID"
            )

        if state_id <= 0:
            fail(
                f"State '{state_name}' has invalid ID {state_id}"
            )

        if state_id in state_ids:
            fail(
                f"Duplicate state ID: {state_id}"
            )

        state_ids.add(state_id)

        if not isinstance(state_name, str):
            fail(
                f"State ID {state_id} has an invalid name"
            )

        if normalize(state_name) == "":
            fail(
                f"State ID {state_id} has an empty name"
            )

        normalized_name = normalize(state_name)

        if normalized_name in state_names:
            fail(
                f"Duplicate state name: {state_name}"
            )

        state_names.add(normalized_name)

        if not isinstance(state_code, str):
            fail(
                f"State '{state_name}' has an invalid code"
            )

        normalized_code = normalize(state_code)

        if normalized_code == "":
            fail(
                f"State '{state_name}' has an empty code"
            )

        if normalized_code in state_codes:
            fail(
                f"Duplicate state code: {state_code}"
            )

        state_codes.add(normalized_code)

        if not isinstance(is_fct, bool):
            fail(
                f"State '{state_name}' has invalid is_fct value"
            )

        if is_fct:
            fct_count += 1
        else:
            normal_state_count += 1

    if normal_state_count != 36:
        fail(
            f"Expected 36 states, found {normal_state_count}"
        )

    if fct_count != 1:
        fail(
            f"Expected exactly 1 FCT record, found {fct_count}"
        )

    print("State/FCT count: PASS")
    print("36 states: PASS")
    print("1 FCT: PASS")
    print("Unique state IDs: PASS")
    print("Unique state names: PASS")
    print("Unique state codes: PASS")

    return {
        "ids": state_ids,
        "codes": state_codes,
        "records": states,
    }


# ============================================================
# Validate LGA structure
# ============================================================

def validate_lgas(lgas, state_info):
    print()
    print("===== VALIDATING LGAs =====")

    if not isinstance(lgas, list):
        fail("lgas.json must contain a JSON array")

    if len(lgas) != EXPECTED_NATIONAL_LGA_COUNT:
        fail(
            f"Expected {EXPECTED_NATIONAL_LGA_COUNT} LGA/Area Council "
            f"records, found {len(lgas)}"
        )

    state_ids = state_info["ids"]
    state_codes = state_info["codes"]

    lga_ids = set()
    source_codes = set()
    names = set()
    state_name_pairs = set()

    lga_count = 0
    area_council_count = 0

    records_by_state = {}

    for index, lga in enumerate(lgas, start=1):

        if not isinstance(lga, dict):
            fail(
                f"LGA record #{index} is not an object"
            )

        required_fields = [
            "id",
            "code",
            "state_id",
            "state_code",
            "name",
            "type",
            "is_active",
        ]

        for field in required_fields:
            if field not in lga:
                fail(
                    f"LGA record #{index} is missing field '{field}'"
                )

        # ----------------------------------------------------
        # Internal ID
        # ----------------------------------------------------

        lga_id = lga["id"]

        if not isinstance(lga_id, int):
            fail(
                f"LGA record #{index} has a non-integer ID"
            )

        if lga_id <= 0:
            fail(
                f"LGA record #{index} has invalid ID {lga_id}"
            )

        if lga_id in lga_ids:
            fail(
                f"Duplicate LGA internal ID: {lga_id}"
            )

        lga_ids.add(lga_id)

        # ----------------------------------------------------
        # Source code
        # ----------------------------------------------------

        source_code = lga["code"]

        if not isinstance(source_code, int):
            fail(
                f"LGA '{lga.get('name')}' has "
                f"a non-integer source code"
            )

        if source_code <= 0:
            fail(
                f"LGA '{lga.get('name')}' has "
                f"invalid source code {source_code}"
            )

        if source_code in source_codes:
            fail(
                f"Duplicate LGA source code: {source_code}"
            )

        source_codes.add(source_code)

        # ----------------------------------------------------
        # State ID
        # ----------------------------------------------------

        state_id = lga["state_id"]

        if not isinstance(state_id, int):
            fail(
                f"LGA '{lga.get('name')}' has "
                f"a non-integer state_id"
            )

        if state_id not in state_ids:
            fail(
                f"LGA '{lga.get('name')}' references "
                f"unknown state_id {state_id}"
            )

        # ----------------------------------------------------
        # State code
        # ----------------------------------------------------

        state_code = lga["state_code"]

        if not isinstance(state_code, str):
            fail(
                f"LGA '{lga.get('name')}' has invalid state_code"
            )

        normalized_state_code = normalize(state_code)

        if normalized_state_code not in state_codes:
            fail(
                f"LGA '{lga.get('name')}' references "
                f"unknown state code '{state_code}'"
            )

        # ----------------------------------------------------
        # Name
        # ----------------------------------------------------

        name = lga["name"]

        if not isinstance(name, str):
            fail(
                f"LGA ID {lga_id} has invalid name"
            )

        normalized_name = normalize(name)

        if normalized_name == "":
            fail(
                f"LGA ID {lga_id} has an empty name"
            )

        if normalized_name in names:
            fail(
                f"Duplicate LGA name: {name}"
            )

        names.add(normalized_name)

        # ----------------------------------------------------
        # State + LGA combination
        # ----------------------------------------------------

        state_name_key = (
            state_id,
            normalized_name,
        )

        if state_name_key in state_name_pairs:
            fail(
                f"Duplicate state/LGA combination: "
                f"state_id={state_id}, name='{name}'"
            )

        state_name_pairs.add(state_name_key)

        # ----------------------------------------------------
        # Type
        # ----------------------------------------------------

        lga_type = lga["type"]

        if not isinstance(lga_type, str):
            fail(
                f"LGA '{name}' has invalid type"
            )

        lga_type = normalize(lga_type).upper()

        if lga_type not in {
            VALID_LGA_TYPE,
            VALID_AREA_COUNCIL_TYPE,
        }:
            fail(
                f"LGA '{name}' has invalid type "
                f"'{lga['type']}'"
            )

        if lga_type == VALID_LGA_TYPE:
            lga_count += 1
        else:
            area_council_count += 1

        # ----------------------------------------------------
        # Active status
        # ----------------------------------------------------

        if not isinstance(lga["is_active"], bool):
            fail(
                f"LGA '{name}' has invalid is_active value"
            )

        # ----------------------------------------------------
        # Group records by state
        # ----------------------------------------------------

        records_by_state.setdefault(state_id, 0)
        records_by_state[state_id] += 1

    # ========================================================
    # National totals
    # ========================================================

    if len(lga_ids) != EXPECTED_NATIONAL_LGA_COUNT:
        fail(
            "Internal LGA IDs are not unique"
        )

    if len(source_codes) != EXPECTED_NATIONAL_LGA_COUNT:
        fail(
            "Source LGA codes are not unique"
        )

    if len(names) != EXPECTED_NATIONAL_LGA_COUNT:
        fail(
            "LGA names are not unique"
        )

    print(
        f"Total records: {len(lgas)} / "
        f"{EXPECTED_NATIONAL_LGA_COUNT}: PASS"
    )

    print(
        f"Unique internal IDs: {len(lga_ids)}: PASS"
    )

    print(
        f"Unique source codes: {len(source_codes)}: PASS"
    )

    print(
        f"Unique names: {len(names)}: PASS"
    )

    print(
        f"LGAs: {lga_count}"
    )

    print(
        f"FCT Area Councils: {area_council_count}"
    )

    return records_by_state


# ============================================================
# Main
# ============================================================

def main():

    print("==============================================")
    print(" Nigeria Online Voting System")
    print(" LGA DATA VALIDATOR")
    print("==============================================")
    print()

    states = load_json(STATES_FILE)

    state_info = validate_states(states)

    if not LGAS_FILE.exists():
        print()
        print("LGA dataset has not been created yet.")
        print()
        print(
            "Validation infrastructure is ready."
        )
        print(
            f"Expected file: {LGAS_FILE}"
        )
        print()
        print(
            "NEXT STEP:"
        )
        print(
            "Place the verified 774-record dataset "
            "in data/lgas.json."
        )
        print()
        print(
            "No LGA data has been accepted yet."
        )
        return

    lgas = load_json(LGAS_FILE)

    records_by_state = validate_lgas(
        lgas,
        state_info,
    )

    print()
    print("==============================================")
    print(" VALIDATION COMPLETE")
    print("==============================================")
    print()
    print("RESULT: PASS")
    print()
    print(
        "The LGA dataset passed all structural checks."
    )


if __name__ == "__main__":
    main()