#!/usr/bin/env python3

import json
import re
import sys
from datetime import datetime, timezone
from pathlib import Path
from urllib.request import Request, urlopen


# ============================================================
# CONFIGURATION
# ============================================================

SOURCE_URL = (
    "https://microdata.nigerianstat.gov.ng/"
    "index.php/catalog/82/variable/F80/V3316?name=lga"
)

OUTPUT_FILE = Path("data/raw_lga_source.json")

EXPECTED_TOTAL = 774

# NBS state-code prefixes.
#
# 01 = Abia
# 02 = Adamawa
# ...
# 36 = Zamfara
# 37 = FCT
#
# We derive the prefix from the first two digits of the
# NBS LGA code, then verify it against states.json.

EXPECTED_STATE_IDS = set(range(1, 38))


# ============================================================
# HELPERS
# ============================================================

def fetch_source(url: str) -> str:
    """
    Download the NBS source page.

    We deliberately use the public NBS catalogue page rather
    than a third-party LGA website.
    """

    request = Request(
        url,
        headers={
            "User-Agent": (
                "NigeriaOnlineVotingSystem/"
                "1.0 LGA Data Importer"
            )
        },
    )

    with urlopen(request, timeout=30) as response:
        charset = response.headers.get_content_charset() or "utf-8"
        return response.read().decode(charset, errors="replace")


def html_to_text(html: str) -> str:
    """
    Convert the small amount of HTML we need into plain text.

    We do not use BeautifulSoup so this script depends only on
    Python's standard library.
    """

    text = re.sub(
        r"<script\b[^>]*>.*?</script>",
        " ",
        html,
        flags=re.IGNORECASE | re.DOTALL,
    )

    text = re.sub(
        r"<style\b[^>]*>.*?</style>",
        " ",
        text,
        flags=re.IGNORECASE | re.DOTALL,
    )

    text = re.sub(
        r"<[^>]+>",
        " ",
        text,
    )

    # Decode common HTML entities without requiring bs4.
    text = text.replace("&nbsp;", " ")
    text = text.replace("&amp;", "&")
    text = text.replace("&quot;", '"')
    text = text.replace("&#39;", "'")
    text = text.replace("&lt;", "<")
    text = text.replace("&gt;", ">")

    return text


def parse_lga_records(html: str) -> list[dict]:
    """
    Extract NBS LGA code/name pairs from the catalogue page.

    Expected source representation:

        101. ABA NORTH
        102. ABA SOUTH
        ...
        3706. ABUJA MUNICIPAL

    We intentionally preserve the source spelling.
    """

    text = html_to_text(html)

    records = []

    # The NBS catalogue displays the category as:
    #
    # 101. ABA NORTH
    #
    # The code is three or four digits.
    #
    # We only accept codes in the known NBS LGA range.
    pattern = re.compile(
        r"\b(\d{3,4})\.\s+([A-Z][A-Z0-9 &'()\-./]+?)(?=\s+\d{3,4}\.\s+|\s*$)",
        re.MULTILINE,
    )

    for match in pattern.finditer(text):

        code = int(match.group(1))
        name = " ".join(match.group(2).split())

        if code < 101 or code > 3706:
            continue

        if not name:
            continue

        records.append(
            {
                "code": code,
                "name": name,
            }
        )

    return records


def load_states() -> list[dict]:
    """
    Load the project's existing states.json.

    We do NOT recreate states here.
    """

    states_file = Path("data/states.json")

    if not states_file.exists():
        raise RuntimeError(
            "data/states.json does not exist."
        )

    with states_file.open(
        "r",
        encoding="utf-8",
    ) as file:
        states = json.load(file)

    if not isinstance(states, list):
        raise RuntimeError(
            "data/states.json must contain a JSON array."
        )

    return states


def validate_states(states: list[dict]) -> dict:
    """
    Validate the parent state dataset before accepting
    any LGA records.
    """

    if len(states) != 37:
        raise RuntimeError(
            f"Expected 37 state/FCT records, got {len(states)}."
        )

    state_map = {}

    for state in states:

        state_id = state.get("id")
        name = state.get("name")
        code = state.get("code")

        if not isinstance(state_id, int):
            raise RuntimeError(
                f"Invalid state ID: {state_id!r}"
            )

        if state_id in state_map:
            raise RuntimeError(
                f"Duplicate state ID: {state_id}"
            )

        if not name:
            raise RuntimeError(
                f"State {state_id} has no name."
            )

        if not code:
            raise RuntimeError(
                f"State {state_id} has no code."
            )

        state_map[state_id] = {
            "name": name,
            "code": code,
            "is_fct": bool(state.get("is_fct", False)),
        }

    if set(state_map.keys()) != EXPECTED_STATE_IDS:
        raise RuntimeError(
            "State IDs must be exactly 1 through 37."
        )

    fct_count = sum(
        1
        for state in state_map.values()
        if state["is_fct"]
    )

    if fct_count != 1:
        raise RuntimeError(
            f"Expected exactly 1 FCT record, got {fct_count}."
        )

    return state_map


def determine_state_id(code: int) -> int:
    """
    Convert the NBS code to our parent State ID.

    Examples:

        101  -> 1
        117  -> 1
        201  -> 2
        3301 -> 33
        3601 -> 36
        3701 -> 37
    """

    return code // 100


def determine_type(state_id: int) -> str:
    """
    FCT records are Area Councils.
    All other records are LGAs.
    """

    if state_id == 37:
        return "AREA_COUNCIL"

    return "LGA"


def validate_source_records(
    records: list[dict],
    state_map: dict,
) -> None:
    """
    Perform structural validation before writing anything.

    This function deliberately rejects bad source data rather
    than trying to repair it automatically.
    """

    print()
    print("===== SOURCE VALIDATION =====")

    # --------------------------------------------------------
    # Count
    # --------------------------------------------------------

    print(
        f"Source records found: {len(records)}"
    )

    if len(records) != EXPECTED_TOTAL:
        raise RuntimeError(
            "SOURCE COUNT FAILED: "
            f"expected {EXPECTED_TOTAL}, "
            f"got {len(records)}."
        )

    print("774 records: PASS")

    # --------------------------------------------------------
    # Duplicate codes
    # --------------------------------------------------------

    codes = [record["code"] for record in records]

    if len(codes) != len(set(codes)):
        duplicates = sorted(
            code
            for code in set(codes)
            if codes.count(code) > 1
        )

        raise RuntimeError(
            "DUPLICATE SOURCE CODES: "
            f"{duplicates}"
        )

    print("Unique source codes: PASS")

    # --------------------------------------------------------
    # Duplicate names
    # --------------------------------------------------------

    normalized_names = [
        record["name"].strip().casefold()
        for record in records
    ]

    if len(normalized_names) != len(
        set(normalized_names)
    ):
        duplicates = sorted(
            name
            for name in set(normalized_names)
            if normalized_names.count(name) > 1
        )

        raise RuntimeError(
            "DUPLICATE SOURCE NAMES: "
            f"{duplicates}"
        )

    print("Unique source names: PASS")

    # --------------------------------------------------------
    # Parent state validation
    # --------------------------------------------------------

    invalid_parent_codes = []

    for record in records:

        code = record["code"]
        state_id = determine_state_id(code)

        if state_id not in state_map:
            invalid_parent_codes.append(code)

    if invalid_parent_codes:
        raise RuntimeError(
            "INVALID PARENT STATE FOR CODES: "
            f"{invalid_parent_codes}"
        )

    print("Parent state mapping: PASS")

    # --------------------------------------------------------
    # Code structure
    # --------------------------------------------------------

    invalid_codes = []

    for record in records:

        code = record["code"]

        state_id = determine_state_id(code)

        # Codes must be 3 digits for the first state block
        # and 4 digits for subsequent state blocks.
        #
        # More importantly, the first two digits must map
        # to state IDs 1-37.

        if state_id < 1 or state_id > 37:
            invalid_codes.append(code)

    if invalid_codes:
        raise RuntimeError(
            "INVALID NBS CODES: "
            f"{invalid_codes}"
        )

    print("Source code structure: PASS")

    # --------------------------------------------------------
    # Empty names
    # --------------------------------------------------------

    empty_names = [
        record["code"]
        for record in records
        if not record["name"].strip()
    ]

    if empty_names:
        raise RuntimeError(
            "EMPTY LGA NAMES: "
            f"{empty_names}"
        )

    print("Non-empty names: PASS")

    # --------------------------------------------------------
    # FCT validation
    # --------------------------------------------------------

    fct_records = [
        record
        for record in records
        if determine_state_id(record["code"]) == 37
    ]

    if len(fct_records) != 6:
        raise RuntimeError(
            "FCT AREA COUNCIL COUNT FAILED: "
            f"expected 6, got {len(fct_records)}."
        )

    print("FCT Area Councils: PASS")

    # --------------------------------------------------------
    # State distribution
    # --------------------------------------------------------

    distribution = {}

    for record in records:

        state_id = determine_state_id(
            record["code"]
        )

        distribution[state_id] = (
            distribution.get(state_id, 0) + 1
        )

    print()
    print("LGA distribution by parent state:")

    for state_id in sorted(distribution):

        state_name = state_map[state_id]["name"]

        print(
            f"  {state_id:02d} "
            f"{state_name:<30} "
            f"{distribution[state_id]}"
        )


def build_raw_records(
    records: list[dict],
    state_map: dict,
) -> list[dict]:
    """
    Build the raw source representation.

    IMPORTANT:

    The NBS name is stored exactly as sourced.
    We do not correct spelling here.
    """

    retrieved_at = datetime.now(
        timezone.utc
    ).replace(microsecond=0).isoformat()

    result = []

    for record in sorted(
        records,
        key=lambda item: item["code"],
    ):

        code = record["code"]
        name = record["name"]

        state_id = determine_state_id(code)

        state = state_map[state_id]

        result.append(
            {
                "source": "National Bureau of Statistics",
                "source_reference": (
                    "NGA_2023_GHSP-W5_v01_M"
                ),
                "source_url": SOURCE_URL,
                "retrieved_at": retrieved_at,

                "code": code,

                # Preserve the NBS spelling exactly.
                "source_name": name,

                "state_id": state_id,

                "state_code": state["code"],

                "state_name": state["name"],

                "type": determine_type(state_id),

                "is_active": True,
            }
        )

    return result


def save_raw_records(records: list[dict]) -> None:
    """
    Write the raw source dataset.

    This is NOT yet the production lgas.json file.
    """

    OUTPUT_FILE.parent.mkdir(
        parents=True,
        exist_ok=True,
    )

    with OUTPUT_FILE.open(
        "w",
        encoding="utf-8",
    ) as file:

        json.dump(
            records,
            file,
            indent=4,
            ensure_ascii=False,
        )

        file.write("\n")


# ============================================================
# MAIN
# ============================================================

def main() -> int:

    print("=" * 55)
    print(" Nigeria Online Voting System")
    print(" NBS RAW LGA SOURCE IMPORTER")
    print("=" * 55)

    try:

        # ----------------------------------------------------
        # Load and validate existing states
        # ----------------------------------------------------

        print()
        print("===== VALIDATING EXISTING STATES =====")

        states = load_states()

        state_map = validate_states(states)

        print("37 state/FCT records: PASS")
        print("State IDs 1-37: PASS")

        # ----------------------------------------------------
        # Download NBS source
        # ----------------------------------------------------

        print()
        print("===== DOWNLOADING NBS SOURCE =====")

        print(
            "Source:"
        )

        print(SOURCE_URL)

        html = fetch_source(
            SOURCE_URL
        )

        if not html:
            raise RuntimeError(
                "NBS source returned empty content."
            )

        print("NBS source download: PASS")

        # ----------------------------------------------------
        # Parse source
        # ----------------------------------------------------

        print()
        print("===== PARSING NBS SOURCE =====")

        records = parse_lga_records(html)

        if not records:
            raise RuntimeError(
                "No LGA records could be extracted "
                "from the NBS source."
            )

        print(
            f"Records extracted: {len(records)}"
        )

        # ----------------------------------------------------
        # Validate source
        # ----------------------------------------------------

        validate_source_records(
            records,
            state_map,
        )

        # ----------------------------------------------------
        # Build raw dataset
        # ----------------------------------------------------

        print()
        print("===== BUILDING RAW DATASET =====")

        raw_records = build_raw_records(
            records,
            state_map,
        )

        print(
            f"Prepared records: {len(raw_records)}"
        )

        # ----------------------------------------------------
        # Save
        # ----------------------------------------------------

        save_raw_records(
            raw_records
        )

        print()
        print("===== IMPORT COMPLETE =====")

        print(
            f"Raw source file created:"
        )

        print(
            f"  {OUTPUT_FILE}"
        )

        print()
        print(
            "IMPORTANT:"
        )

        print(
            "This file contains the NBS source names "
            "without correction."
        )

        print(
            "It is NOT yet the final verified lgas.json."
        )

        print()
        print(
            "NEXT STEP:"
        )

        print(
            "Run the reconciliation validator."
        )

        return 0

    except Exception as exc:

        print()
        print("===== IMPORT FAILED =====")
        print()
        print(str(exc))
        print()

        return 1


if __name__ == "__main__":
    sys.exit(main())