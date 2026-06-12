# Criptografia de Snapshots — Quick Start

Para versionar snapshots com segurança em repositórios git, use criptografia AES-256-GCM.

## TL;DR — 30 segundos

```bash
# 1. Defina uma frase-senha FORTE
export ANALYZER_PASSPHRASE="Sua_Frase_Secreta_Muito_Forte_123!"

# 2. Salve o snapshot (criptografado automaticamente)
./investment-analyzer analyze --file carteira.csv --save

# 3. Commite com segurança
git add ~/.investment-analyzer/snapshots/*.json.enc
git commit -m "snapshot jun/2026"
git push

# 4. Em outra máquina, use a mesma frase-senha
export ANALYZER_PASSPHRASE="Sua_Frase_Secreta_Muito_Forte_123!"
./investment-analyzer history --months 12  # funciona automaticamente
```

## Detalhes Importantes

### ✅ Faça isso

```bash
# Use variável de ambiente (não fica em shell history)
export ANALYZER_PASSPHRASE="frase-secreta-forte"
./investment-analyzer analyze --file carteira.csv --save

# Snapshots são criados como .json.enc (criptografados)
ls ~/.investment-analyzer/snapshots/
# 2026-06.json.enc ← seguro para git

# Carregar é transparente
./investment-analyzer history --months 12
```

### ❌ Evite isso

```bash
# Não use --passphrase na linha de comando (fica em shell history)
./investment-analyzer analyze --file carteira.csv --passphrase "senha" --save

# Não commite snapshots plaintext
git add ~/.investment-analyzer/snapshots/2026-06.json

# Não escriba a passphrase em arquivo
echo "ANALYZER_PASSPHRASE=minha-senha" > .env
git add .env
```

## Primitivas de Segurança

| Elemento | Detalhe |
|---|---|
| Cifra | AES-256-GCM (padrão de indústria) |
| Derivação | Argon2id com 64MB RAM (resistente a GPU) |
| Integridade | GCM detecta qualquer bit alterado |
| Salt | 32 bytes aleatórios por arquivo |
| Nonce | 12 bytes aleatórios por arquivo |

## Força da Frase-Senha

Recomendação: **mínimo 16 caracteres** com mistura de:
- Maiúsculas (`A-Z`)
- Minúsculas (`a-z`)
- Números (`0-9`)
- Símbolos (`!@#$%^&*`)

```bash
# ✅ BOM
export ANALYZER_PASSPHRASE="MyPortfolio_2026_AES256!Secure$"

# ✅ BOM
export ANALYZER_PASSPHRASE="Investimentos_Carteira_Segura_2026_#123"

# ❌ FRACO
export ANALYZER_PASSPHRASE="password"
export ANALYZER_PASSPHRASE="123456"
export ANALYZER_PASSPHRASE="carteira"
```

## Perda de Frase-Senha

Se perder a frase-senha:
- Os arquivos `.json.enc` continuam intactos
- Mas são **irrecuperáveis sem a frase-senha original**
- Não há "modo recovery"

**Use um gerenciador de senhas** (1Password, Bitwarden, KeePass) para guardar a frase-senha com segurança.

## Ver Arquivo Encriptado

```bash
# Verificar que é binário (não texto)
file ~/.investment-analyzer/snapshots/2026-06.json.enc
# output: data

# Ver bytes aleatórios (salt + ciphertext)
xxd ~/.investment-analyzer/snapshots/2026-06.json.enc | head -3

# Tentar ler como JSON (não funciona)
cat ~/.investment-analyzer/snapshots/2026-06.json.enc
# output: (caracteres aleatórios ilegíveis)
```

## Mudar Frase-Senha

Não é possível "reencriptar" um arquivo. Para mudar:

1. Descriptografar com a frase ANTIGA
2. Deletar arquivo criptografado
3. Salvar com frase NOVA

```bash
# Frase antiga
export ANALYZER_PASSPHRASE="frase-antiga"
./investment-analyzer analyze --file carteira.csv --month 2026-06 --save

# Remover arquivo
rm ~/.investment-analyzer/snapshots/2026-06.json.enc

# Frase nova
export ANALYZER_PASSPHRASE="frase-nova"
./investment-analyzer analyze --file carteira.csv --month 2026-06 --save

# Commitar
git add ~/.investment-analyzer/snapshots/2026-06.json.enc
git commit -m "update snapshot with new encryption"
```

## Compatibilidade com Git

```bash
# Snapshots criptografados são seguros para git
git add ~/.investment-analyzer/snapshots/*.json.enc
git commit -m "encrypted snapshots"
git push

# .gitignore bloqueia .json plaintext, permite .json.enc
*.json       ← ignorado (plaintext)
!*.json.enc  ← permitido (criptografado)
```

## Tampering Detection

Se alguém tentar modificar um arquivo `.json.enc`:

```bash
# Modificar 1 bit do arquivo
dd if=/dev/urandom of=~/.investment-analyzer/snapshots/2026-06.json.enc bs=1 count=1 conv=notrunc seek=100

# Tentar carregar
ANALYZER_PASSPHRASE=minha-senha ./investment-analyzer history
# Error: cipher: message authentication failed ← detectado!
```

GCM garante que qualquer alteração é detectada.

---

**Para detalhes completos, ver [SECURITY.md](SECURITY.md)**
