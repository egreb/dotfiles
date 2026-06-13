function review --description "Save flagged JSON from clipboard under <org>/<wall>/<log>.json and launch claude review"
    if test (count $argv) -ne 1
        echo "usage: review <wall_id>" >&2
        return 2
    end
    set -l wall_id $argv[1]

    set -l json (pbpaste)
    if test -z "$json"
        echo "error: clipboard is empty" >&2
        return 1
    end

    set -l parsed (printf '%s' "$json" | python3 -c '
import json, sys
raw = sys.stdin.read().lstrip()
try:
    d, end = json.JSONDecoder().raw_decode(raw)
    trailing = raw[end:].strip()
    l = d["request"]["labels"]
    print(l["cybrOrgId"])
    print(l["siegeId"])
    print("TRAILING" if trailing else "CLEAN")
    sys.stdout.write(json.dumps(d))
except Exception as e:
    sys.stderr.write(f"parse error: {e}\n")
    sys.exit(1)
')
    or return 1

    set -l org_id $parsed[1]
    set -l log_id $parsed[2]
    set -l status $parsed[3]
    set -l clean_json $parsed[4..-1]
    set -l dir "$org_id/$wall_id"
    set -l file "$dir/$log_id.json"

    if test -e $file
        echo "error: $file already exists, refusing to overwrite" >&2
        return 1
    end

    if test "$status" = TRAILING
        echo "note: clipboard had extra data after the first JSON object — saving only the first object" >&2
    end

    mkdir -p $dir
    printf '%s' "$clean_json" > $file
    echo "saved: $file  (org=$org_id wall=$wall_id log=$log_id)"

    claude "Review $file"
end
