# Building Capybot

## Package

1. Build the APK
    - A signing pair for rsa keys is needed. Assuming you have `melange.rsa` & `melange.rsa.pub`, run

    ```bash
    podman run --rm -v "${PWD}":/work cgr.dev/chainguard/melange keygen
    ```

    - Build APK for all available architectures. Run

    ```bash
    podman run --rm --privileged -v "${PWD}":/work \
        cgr.dev/chainguard/melange build melange.yaml \
        --arch amd64,aarch64,armv7 \
        --signing-key melange.rsa
    ```

    - Make sure the package exists under `./packages`

## Containerise

### Building Locally

```bash
GITHUB_USERNAME="myuser"
REF="ghcr.io/${GITHUB_USERNAME}/capybot"

podman run --rm -v "${PWD}":/work \
    cgr.dev/chainguard/apko build --debug apko.yaml \
    "${REF}" output.tar -k melange.rsa.pub \
    --arch amd64,aarch64,armv7
```

To use or test this image, run

```bash
ARCH_REF="$(podman load < output.tar | grep "Loaded image" | sed 's/^Loaded image: //' | head -1)"

podman run -v ./config.json:/config.json -v ./images.txt:/images.txt "${ARCH_REF}" -config-file /config.json
```

### Uploading to Registry

```bash
GITHUB_USERNAME="myuser"
REF="ghcr.io/${GITHUB_USERNAME}/capybot"

# A personal access token with the "write:packages" scope
GITHUB_TOKEN="*****"

podman run --rm -v "${PWD}":/work \
    -e REF="${REF}" \
    -e GITHUB_USERNAME="${GITHUB_USERNAME}" \
    -e GITHUB_TOKEN="${GITHUB_TOKEN}" \
    --entrypoint sh \
    cgr.dev/chainguard/apko -c \
        'echo "${GITHUB_TOKEN}" | \
            apko login ghcr.io -u "${GITHUB_USERNAME}" --password-stdin && \
            apko publish --debug apko.yaml \
                "${REF}" -k melange.rsa.pub \
                --arch amd64,aarch64,armv7'
```