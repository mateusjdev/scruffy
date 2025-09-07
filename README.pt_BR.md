# Scruffy

> [!CAUTION]
> Este repositório contém código instável, caso ocorra um `rm -rf /` o problema não é meu, use com cautela.

---

Utilitario geral, mas no momento, é apenas uma versão em `go` do [mateusjdev/rename-files-to-hash](https://github.com/mateusjdev/rename-files-to-hash).

## 🛠️ Compilando

### 📦 Pré-requisitos

Certifique-se de ter o seguinte instalado:

- [Go](https://golang.org/doc/install) (versão ≥ 1.21)
- Para o processo de compilação, [Taskfile](https://taskfile.dev/#/installation) e [Git](https://git-scm.com/downloads) são altamente recomendados.

### 🏗️ Compilando

1. Clone o repositório para sua máquina local:

```shell
git clone https://github.com/mateusjdev/scruffy
cd scruffy
```

2. Compile o projeto:

```shell
go build -o ./build/scruffy
# ou usando Taskfile (recomendado)
task build
```