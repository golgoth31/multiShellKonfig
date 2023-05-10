# multiShellKonfig

## Install

Download the needed binary from the [release page](https://github.com/golgoth31/multiShellKonfig/releases) and save it as **msk-bin** somewhere in your **$PATH**.

```sh
curl -fsL -o msk-bin https://github.com/golgoth31/multiShellKonfig/releases/download/v0.0.4/msk-bin_v0.0.4_darwin_arm64
chmod +x msk-bin
mv msk-bin <bin path\>
```

## How to use

Source the wapper for the needed shell:

```sh
source <(msk-bin shellwrapper zsh)
```

To enter a new context:

```sh
msk context (or ctx)
```

To enter a new namespace:

```sh
msk namespace (or ns)
```

## To do

- context in alphabetical order
- clean old contexts
