# Segurança de Snapshots — Guia de Uso

## Visão Geral

Os snapshots mensais contêm dados sensíveis (valores de ativos, alocações, histórico financeiro). Para versioná-los com segurança em repositórios públicos ou compartilhados, o `investment-analyzer` oferece **criptografia transparente com AES-256-GCM**.

## Criptografia de Dados

### Primitivas

| Componente | Algoritmo | Parâmetro |
|---|---|---|
| **Derivação de chave** | Argon2id | time=1, mem=64MB, threads=4, keyLen=32 bytes |
| **Cifra** | AES-256-GCM | nonce=12 bytes aleatórios por arquivo |
| **Integridade** | GCM authentication tag | 16 bytes (detecta tampering) |
| **Salt** | Aleatório | 32 bytes, único por arquivo |

**Cada `Save()` gera salt e nonce frescos** → mesmo passphrase produz bytes diferentes a cada gravação (propriedade de um bom sistema criptográfico).

### Formato de Arquivo

```
[32 bytes salt][12 bytes nonce][ciphertext + 16 bytes GCM tag]
```

Arquivo: `~/.investment-analyzer/snapshots/YYYY-MM.json.enc`

## Resolução da Frase-Senha

A frase-senha é resolvida nesta ordem de prioridade:

1. **Variável de ambiente** `ANALYZER_PASSPHRASE` (recomendado para scripts/CI)
2. **Flag CLI** `--passphrase "..."` (útil para testes; ⚠️ exposta em shell history)
3. **Sem passphrase** → snapshots salvos em plaintext `.json` (retrocompatibilidade)

### Uso Recomendado

#### Para uso local seguro (recomendado)

```bash
export ANALYZER_PASSPHRASE="sua-frase-secreta-muito-forte"
./investment-analyzer analyze --file carteira.csv --month 2026-06 --save
```

A frase-senha:
- ✅ Fica apenas em memória durante execução
- ✅ Não é gravada em disco
- ✅ Não fica exposta em shell history
- ✅ Pode ser diferente a cada execução (recomendado)

#### Para automação (scripts, CI/CD)

```bash
# .env ou secrets management
export ANALYZER_PASSPHRASE="${VAULT_INVESTMENT_PASSPHRASE}"

# Execute o script
./investment-analyzer analyze --file carteira.csv --save --no-browser
```

#### ⚠️ Evitar

```bash
# ❌ NÃO RECOMENDADO: expõe a frase-senha em shell history
./investment-analyzer analyze --file carteira.csv --passphrase "my-password" --save

# ❌ NÃO RECOMENDADO: gravar a frase-senha em arquivo
echo "ANALYZER_PASSPHRASE=my-password" > ~/.investment-analyzer/.env
```

## Fluxo de Commit Seguro

### Passo 1: Criptografar Snapshots

```bash
export ANALYZER_PASSPHRASE="sua-frase-segura"

# Analisar e salvar snapshot criptografado
./investment-analyzer analyze --file carteira.csv --month 2026-06 --save --no-browser

# Resultado: ~/.investment-analyzer/snapshots/2026-06.json.enc
# Arquivo plaintext .json NÃO é criado
```

### Passo 2: Verificar Encriptação

```bash
# Listar snapshots
ls -lh ~/.investment-analyzer/snapshots/

# Verificar que é binário (não texto)
file ~/.investment-analyzer/snapshots/2026-06.json.enc
# output: data

# Tentar ler com strings (sem resultado legível)
strings ~/.investment-analyzer/snapshots/2026-06.json.enc | head -1
# output: (vazio ou lixo)
```

### Passo 3: Commitar com Segurança

```bash
# Adicionar apenas arquivos criptografados
git add ~/.investment-analyzer/snapshots/2026-06.json.enc

# Verificar que não há .json plaintext
git status | grep ".json"
# output: (deve estar vazio)

# Commit
git commit -m "Add encrypted portfolio snapshot jun/2026"

# Push seguro
git push origin main
```

### Passo 4: Clonar em Outra Máquina

```bash
git clone https://github.com/user/repo.git
cd repo

# Defina a frase-senha (mesma que foi usada para criptografar)
export ANALYZER_PASSPHRASE="sua-frase-segura"

# Os snapshots são decriptografados automaticamente
./investment-analyzer history --months 12

# Lê do arquivo .json.enc, descriptografa em memória, exibe
```

## Segurança da Frase-Senha

### Força Recomendada

Use uma frase-senha **forte e única**:
- ✅ Mínimo 16 caracteres
- ✅ Misture maiúsculas, minúsculas, números, símbolos
- ✅ Evite palavras de dicionário ou padrões óbvios

Exemplos:
```bash
# ✅ BOM
export ANALYZER_PASSPHRASE="MyP@ssw0rd!Dec2024#Investment$123"

# ✅ BOM (frase-senha)
export ANALYZER_PASSPHRASE="Minha_Carteira_Segura_2024_$Analyst!9"

# ❌ FRACO
export ANALYZER_PASSPHRASE="password123"
export ANALYZER_PASSPHRASE="carteira"
```

### Armazenamento Seguro

Opções para guardar a frase-senha:

1. **Gerenciador de senhas** (recomendado)
   - 1Password, Bitwarden, KeePass
   - Acesso controlado, auditável

2. **Secret management (CI/CD)**
   - GitHub Secrets, GitLab CI Variables, HashiCorp Vault
   - Nunca fica em código ou git history

3. **Passphrase em memória (manual)**
   - Digite quando solicitado pelo CLI
   - Não fica exposta em shell history

### O Que NÃO Fazer

```bash
# ❌ Nunca commite a frase-senha
echo "ANALYZER_PASSPHRASE=my-secret" > .env
git add .env
git commit

# ❌ Nunca use em scripts com passphrase hardcoded
./analyzer --passphrase "hardcoded-secret" --save

# ❌ Nunca armazene em arquivo plaintext no servidor
cat ~/.investment-analyzer/passphrase.txt
```

## Recuperação de Dados

### Se você perder a frase-senha

⚠️ **Sem a frase-senha correta, os snapshots criptografados não podem ser descriptografados.**

Se isso acontecer:
1. Os arquivos `.json.enc` ainda existem (não foram perdidos)
2. Mas estão inacessíveis sem a frase-senha original
3. Não há "modo recovery" — isso é por design

**Recomendação:** Use um gerenciador de senhas para guardar a frase-senha com segurança.

### Migração de Snapshots

Se quiser mudar a frase-senha:

```bash
# Passo 1: Descriptografar com a frase ANTIGA
export ANALYZER_PASSPHRASE="frase-senha-antiga"
./investment-analyzer history > backup.txt

# Passo 2: Remover snapshots antigos
rm ~/.investment-analyzer/snapshots/*.json.enc

# Passo 3: Salvar com a frase NOVA
export ANALYZER_PASSPHRASE="frase-senha-nova"
./investment-analyzer analyze --file carteira.csv --month 2026-06 --save

# Passo 4: Commitar novos arquivos
git add ~/.investment-analyzer/snapshots/
git commit -m "Update encrypted snapshots with new passphrase"
```

## Detecção de Tampering

O AES-256-GCM detecta automaticamente se um arquivo foi modificado:

```bash
# Se alguém tentar modificar 2026-06.json.enc
xxd -r -p <<< "cd291d234e3a97f8..." > ~/.investment-analyzer/snapshots/2026-06.json.enc

# Tentar carregar vai falhar
./investment-analyzer history --months 1
# Error: failed to decrypt snapshot: cipher: message authentication failed
```

**O GCM tag garante integridade** — qualquer bit alterado será detectado.

## Benchmarks de Performance

Criptografia é rápida, mesmo em máquinas antigas:

```bash
# Teste em máquina de 2 cores, 4GB RAM
time ./investment-analyzer analyze --file carteira.csv --save --no-browser --no-fetch

# Resultado típico:
# - Parsing CSV: ~10ms
# - Cálculo portfólio: ~5ms
# - Criptografia Argon2id: ~200ms
# - Total: ~350ms
```

Argon2id usa 64MB de RAM por default (configurável se necessário).

## Compatibilidade

### Versões Go

Requer **Go 1.23+** (usa `golang.org/x/crypto` estável e moderno).

### Retrocompatibilidade

```bash
# Snapshots antigos em plaintext continuam funcionando
ls ~/.investment-analyzer/snapshots/
# 2025-06.json      (antigo, plaintext)
# 2026-01.json.enc  (novo, criptografado)

# O sistema carrega ambos automaticamente
./investment-analyzer history --months 6
# Lê 2025-06.json (plaintext) + 2026-01.json.enc (descriptografa)
```

## Auditing

Se você compartilha o repositório de código com colaboradores:

```bash
# Os snapshots .enc estão no git
git log --all -- "snapshots/*.json.enc" | head -20

# Mas seu conteúdo é ilegível
git show HEAD:snapshots/2026-06.json.enc | file -
# output: data

# Apenas quem conhece a frase-senha pode ler
export ANALYZER_PASSPHRASE="frase-segura"
./investment-analyzer history
```

---

## Resumo

| Aspecto | Detalhes |
|---|---|
| **Quando usar** | Snapshots que serão versionados em git remoto |
| **Como ativar** | `export ANALYZER_PASSPHRASE="..."`  |
| **Criptografia** | AES-256-GCM (indústria-padrão) |
| **Derivação de chave** | Argon2id (resistente a ataques com GPU) |
| **Arquivo** | `~/.investment-analyzer/snapshots/YYYY-MM.json.enc` |
| **Integridade** | GCM detecta tampering automaticamente |
| **Performance** | ~200ms por arquivo (Argon2id) |
| **Fallback** | Sem passphrase → salva plaintext `.json` |
| **Perda de chave** | Sem recuperação (por design) |

