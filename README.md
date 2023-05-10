# multiShellKonfig

Freely inspired by [konf-go](https://github.com/SimonTheLeg/konf-go)

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

You can define some alias for convinience, for example:

```sh
alias kns="msk ns"
alias kctx="msk ctx"
```

## To do

- allow to call 1 specific context
- clean old contexts
- completion
- tests
