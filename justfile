default:
    air

gen:
    go run ./cmd/internal gen

render-css:
    tailwindcss -i apis/render/css/style.css -o apis/render/static/style.css

meilisearch port="7700" key="dev-key":
    mkdir -p work/meilisearch
    meilisearch --config-file-path /dev/null --env development --db-path work/meilisearch --http-addr 127.0.0.1:{{port}} --master-key {{key}}

clean:
    rm -rf work
    rm -rf tmp
    rm -f result

test-build:
    nix build --no-link .#backend
    nix build --no-link .#web
