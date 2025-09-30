#!/usr/bin/env bash
set -euo pipefail

# Generate Razorpay webhook signature for a given payload using your webhook secret.
#
# Secret resolution priority:
#  1) --secret <value>
#  2) RAZORPAY_WEBHOOK_SECRET env
#  3) WEBHOOK_SECRET env
#  4) Prompt (hidden input)
#
# Usage examples:
#  scripts/gen_razorpay_signature.sh --payload payload.json
#  scripts/gen_razorpay_signature.sh --data '{"foo":123}'
#  cat payload.json | scripts/gen_razorpay_signature.sh
#  SECRET=your_secret scripts/gen_razorpay_signature.sh --payload payload.json

PAYLOAD_FILE=""
INLINE_DATA=""
SECRET=""

print_help() {
  cat <<'EOF'
Generate Razorpay webhook signature (HMAC-SHA256 over raw payload)

Options:
  --payload <file>     Path to JSON payload file.
  --data <json>        Inline JSON data (use single quotes to avoid shell escaping).
  --secret <value>     Webhook secret (otherwise read from env or prompt).
  -h, --help           Show help.

Notes:
- The signature must be computed over the exact raw bytes you send (use curl --data-binary @file).
- Expected header: X-Razorpay-Signature
EOF
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --payload)
      PAYLOAD_FILE=${2:-}
      shift 2
      ;;
    --data)
      INLINE_DATA=${2:-}
      shift 2
      ;;
    --secret)
      SECRET=${2:-}
      shift 2
      ;;
    -h|--help)
      print_help
      exit 0
      ;;
    *)
      echo "Unknown argument: $1" >&2
      print_help
      exit 1
      ;;
  esac
done

# Resolve secret
if [[ -z "${SECRET}" ]]; then
  SECRET=${RAZORPAY_WEBHOOK_SECRET:-${WEBHOOK_SECRET:-}}
fi
if [[ -z "${SECRET}" ]]; then
  read -r -s -p "Enter webhook secret: " SECRET
  echo "" >&2
fi
if [[ -z "${SECRET}" ]]; then
  echo "Error: webhook secret is empty." >&2
  exit 1
fi

# Prepare payload bytes into a temp file for robust hashing
TMP=$(mktemp)
trap 'rm -f "$TMP"' EXIT

if [[ -n "${PAYLOAD_FILE}" ]]; then
  if [[ ! -f "${PAYLOAD_FILE}" ]]; then
    echo "Error: payload file not found: ${PAYLOAD_FILE}" >&2
    exit 1
  fi
  cat "${PAYLOAD_FILE}" > "$TMP"
elif [[ -n "${INLINE_DATA}" ]]; then
  printf "%s" "${INLINE_DATA}" > "$TMP"
else
  # Read stdin; if stdin is a TTY, guide the user to paste and Ctrl-D
  if [[ -t 0 ]]; then
    echo "Paste JSON payload, then press Ctrl-D:" >&2
  fi
  cat > "$TMP"
fi

# Compute signature as hex
SIG=$(openssl dgst -sha256 -hmac "${SECRET}" -binary "$TMP" | xxd -p -c 256)
printf "%s\n" "$SIG"
