# Guia de Desenvolvimento de Bibliotecas e Evolucao de Repositorios

Este documento estabelece as diretrizes completas de arquitetura, estilo de código (derivadas de `code-guide.md`), regras de ignorados (`.gitignore`), estratégia de tagging, versionamento de backlog (`CHANGELOG.md`), convenção de nomenclatura baseada em **The Silmarillion** (e O Senhor dos Anéis) e a estratégia de evolução cronológica retroativa de commits para bibliotecas do ecossistema.

---

## 1. Diretrizes de Codigo e Arquitetura (`code-guide.md`)

### Estilo de Codigo e Estrutura
- **Tamanho de Funcoes**: Entre 4 e 20 linhas de codigo. Se exceder 20 linhas, refatorar em subfuncoes com responsabilidade unica (SRP).
- **Tamanho de Arquivos**: Abaixo de 500 linhas. Dividir por responsabilidade clara.
- **Nomenclatura**: Nomes especificos e unicos. Evitar termos genericos como `data`, `handler`, `Manager`. Dar preferencia a nomes que retornem poucas ocorrencias no busca do codigo.
- **Tipagem Explicita**: Tipagem forte e explicita. Proibido o uso de `any`, `Dict` genericos ou funcoes sem tipagem de retorno.
- **Sem Duplicacao**: Extrair logica compartilhada em funcoes ou modulos dedicados.
- **Early Returns**: Preferir retornos antecipados em vez de `if`s aninhados. Nivel maximo de indentacao: 2 niveis.
- **Mensagens de Erro/Excecao**: Devem obrigatoriamente incluir o valor invalido recebido e o formato/shape esperado (ex: `config: invalid port -1, expected positive integer between 1 and 65535`).

### Comentarios e Documentacao
- **Manter Comentarios**: Nao remover comentarios existentes em refatoracoes; eles carregam a intencao e a proveniencia do codigo.
- **Escrever o PORQUÊ, nao o O QUÊ**: Evitar comentarios obvios como `// incrementa contador` sobre `i++`.
- **Docstrings em Funcoes Publicas**: Toda funcao/metodo publico deve ter docstring com a intencao + exemplo de uso pratico (ex: `Usage example: ...`).
- **Referencias**: Citar numeros de issues ou SHAs de commits quando uma linha existir por restricao especifica.

### Testes (F.I.R.S.T.)
- **Comando Único**: Todos os testes devem rodar com um unico comando (`go test ./...` ou `pytest`).
- **Cobertura**: Toda nova funcao recebe um teste unitario. Toda correcao de bug recebe um teste de regressao.
- **Mock de I/O Externo**: Simular chamadas externas (API, DB, Filesystem) com fakes nomeados, evitando stubs inline.
- **Principios F.I.R.S.T.**: Testes devem ser **F**ast (rapidos), **I**ndependent (independentes), **R**epeatable (repetiveis), **S**elf-validating (auto-validaveis) e **T**imely (escritos no momento certo).

### Injecao de Dependencias
- **Injecao Explicita**: Injetar dependencias via construtor ou parametros de funcao, nunca via importacao/global ou container oculta.
- **Encapsulamento de Bibliotecas de Terceiros**: Envolver bibliotecas externas sob uma interface fina pertencente ao proprio projeto.
- **Composicao Transparente**: Compor dependencias explicitamente na funcao `main()`.

### Logging e Formatação
- **Logging Estruturado**: Utilizar formato JSON estruturado ao gravar logs para observabilidade/debugging.
- **Texto Puro apenas no CLI**: Logs em texto puro sao restritos a saidas CLI para usuarios humanos.
- **Formatador Oficial**: Utilizar o formatador padrao da linguagem (`gofmt`, `prettier`, `black`).

---

## 2. Boas Praticas Especificas para Go (Golang)

1. **Alinhamento de Memoria em Structs**:
   - Ordenar os campos das structs do maior tamanho em bytes para o menor (`pointers` -> `int64`/`uint64` -> `int32`/`float32` -> `bool`) para reduzir padding de memoria e otimizar cache locality.

2. **Interfaces Minimais no Consumidor**:
   - Definir interfaces pequenas e focadas no lado do consumidor, e nao no provedor.

3. **Receutores de Valor vs Ponteiro**:
   - Usar *value receivers* para tipos pequenos e imutaveis; usar *pointer receivers* para tipos com mutacao de estado ou structs grandes.

4. **Variacoes de Injecao**:
   - Dependencia obrigatoria: injecao via construtor (`NewService(...)`).
   - Dependencia opcional: injecao via setter ou functional options (`WithTimeout(...)`).
   - Comportamento único: injecao via funcao de primeira classe (`func(req *http.Request) bool`).

---

## 3. Regras de Estilizacao e Documentacao Markdown

1. **PROIBIDO USO DE ICONES E EMOJIS**:
   - Nao utilizar icones, emojis ou caracteres decorativos em arquivos Markdown, titulos, README.md, documentacao, docstrings ou mensagens de commit no Git.

2. **PROIBIDO CAMELCASE EM TITULOS**:
   - Titulos e secoes (headers `#`, `##`, `###`) devem utilizar texto claro separado por espacos (Title Case ou Sentence Case).
   - Nao utilizar CamelCase em titulos de secao (exemplo: use `## Basic routing engine` em vez de `## BasicRoutingEngine`).

---

## 4. Convencao de Nomenclatura baseada em "The Silmarillion"

Para nomear bibliotecas, proxies, motores de observabilidade e microsservicos, utilize como inspiracao a mitologia de J.R.R. Tolkien em **The Silmarillion** e **O Senhor dos Anéis**:

| Nome | Origem em The Silmarillion / LOTR | Metaphora para Engenharia de Software |
| :--- | :--- | :--- |
| **Palantír / PalantirProxy** | As Pedras de Visao de Gondor e Arnor | Reverse Proxy, Load Balancer, Observabilidade em tempo real e Roteamento |
| **Silmaril / SilmarilProxy** | As três gemas que contêm a luz imaculada de Valinor | Criptografia TLS/HTTPS pura, Seguranca e Alta Performance |
| **Vingilot** | O navio de Eärendil que navega pelos ceus | Message Queue, Event Streaming, Transportador de Mensagens e Rotas |
| **Illuin / Ormal** | As Grandes Lampadas que iluminavam Arda | Sistema de Monitoramento, Tracing distribuido e Telemetria |
| **Utumno** | A antiga fortaleza subterrânea de Melkor | Engine de Armazenamento de Baixo Nível, Data Lake ou Log Store |
| **Ainur / Valar** | Os espíritos primordiais que moldaram Arda | Planejador de Orquestracao, Control Plane de Contêineres / Kubernetes |
| **Eru / Ilúvatar** | O Criador Supremo | Process Supervisor, Root Orchestrator ou Init System |
| **Dagor Bragollach** | A Batalha da Chama Repentina | Circuit Breakers, Testes de Carga Extrema e Resiliencia a picos de tráfego |
| **Noldor** | O povo elfo dos maiores artífices e sabios | Pacotes de Algoritmos de Alta Eficiência e Estruturas de Dados Core |
| **Ancalagon** | O maior dragao de Morgoth | Benchmark Suite, Suite de Stress Test e Carga Massiva |

---

## 5. Arquivos de Suporte e Versionamento (`.gitignore`, `VERSION`, `CHANGELOG.md`)

Todo repositorio deve incluir os seguintes arquivos padronizados:

1. **`.gitignore`**: Excluir binarios (`*.exe`, `*.dll`), pastas de build (`dist/`, `bin/`), arquivos de IDE (`.idea/`, `.vscode/`, `.gemini/`), logs temporarios (`*.log`, `*.tmp`) e artefatos de cobertura (`coverage.out`).
2. **`VERSION`**: Arquivo com a versao SemVer atual (exemplo: `1.0.0`).
3. **`CHANGELOG.md`**: Seguindo [Keep a Changelog](https://keepachangelog.com/en/1.0.0/) e [Semantic Versioning](https://semver.org/spec/v2.0.0.html).
4. **`CONTRIBUTING.md`** e **`CODE_OF_CONDUCT.md`**: Diretrizes de contribuicao e conduta da comunidade.

---

## 6. Estrategia de Evolucao Cronologica de Commits no Git

### Diretrizes de Histórico Git

1. **Idioma dos Commits**: Sempre em **Ingles**.
2. **Formato da Mensagem**:
   - **Titulo**: Curto, direto e no imperativo (`feat: ...`, `fix: ...`, `docs: ...`, `refactor: ...`, `perf: ...`).
   - **Corpo**: Descricao multilinha explicando *o que* foi alterado e *por que*.
   - **Zero Emojis**: Sem caracteres decorativos ou simbolos na mensagem.

3. **Variacao de Frequência de Commits**:
   - Distribuição de commits ao longo de meses e anos a partir de uma data inicial (exemplo: a partir de Fevereiro de 2018).
   - Em determinados dias ativos, realizar **mais de 1 commit no mesmo dia** (ex: 2 ou 3 commits em horários distintos do dia com aleatoriedade).

4. **Horários e Segundos Aleatórios**:
   - **NUNCA utilizar segundos zerados (`00`)** nas datas retroativas.
   - Sempre incluir horas, minutos e segundos variaveis e realistas (exemplo: `11:15:37 -0300`, `14:42:19 -0300`, `18:05:48 -0300`).

5. **Criacao de Tags Git (Release Versioning)**:
   - Criar tags Git anotadas mapeando as versoes do `CHANGELOG.md` e `VERSION` (ex: `git tag -a v1.0.0 -m "Release v1.0.0"`).

---

## 7. Workflow de Execucao PowerShell para Commits e Tags Retroativos

```powershell
# 1. Adicionar arquivos alterados
git add .

# 2. Criar commit inicial com mensagem padronizada
git commit -m "feat: short imperative title in english

Detailed description of what was implemented and the rationale behind
the engineering decisions."

# 3. Definir variaveis de data no ambiente PowerShell com segundos especificos
$env:GIT_AUTHOR_DATE="YYYY-MM-DD HH:MM:SS -0300"
$env:GIT_COMMITTER_DATE="YYYY-MM-DD HH:MM:SS -0300"

# 4. Alterar o commit para aplicar as datas retroativas
git commit --amend --no-edit --date "YYYY-MM-DD HH:MM:SS -0300"

# 5. Criar tag Git anotada para versoes de release
$env:GIT_COMMITTER_DATE="YYYY-MM-DD HH:MM:SS -0300"
git tag -a v1.0.0 -m "Release v1.0.0"
```

---

## 8. Checklist de Validacao Antes do Push Remote

Before pushing commits to remote host (`git push origin main`), verify:

1. [ ] Function lengths are between 4 and 20 lines of code.
2. [ ] All Go struct fields are ordered largest to smallest byte size.
3. [ ] Exception/error messages include offending value and expected shape.
4. [ ] Docstrings on public functions include usage examples.
5. [ ] Tests are FIRST and run with a single command (`go test ./...`).
6. [ ] `.gitignore` excludes binaries, IDE folders, and temporary logs.
7. [ ] `VERSION` and `CHANGELOG.md` reflect current release version.
8. [ ] No emojis or icons in Markdown files or git commit messages.
9. [ ] No CamelCase headers in Markdown documentation.
10. [ ] Git timeline has retroactively stamped dates with non-zero seconds (`HH:MM:SS`).
