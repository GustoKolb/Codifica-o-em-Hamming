
# Codificador e Decodificador de Hamming em Go

Este projeto implementa um codificador e decodificador Hamming (31, 26) em Go. O programa codifica um arquivo original em um arquivo maior que pode corrigir erros de um único bit, sendo ideal para testar capacidades de correção de erros.

## Funcionalidades

* **Codificação**: Converte uma sequência de 26 bits do arquivo original em uma sequência de 32 bytes, onde 31 bytes representam os dados e o último byte serve como separador.
* **Decodificação**: Restaura o arquivo original a partir do arquivo codificado com Hamming.
* **Teste de Correção de Erros**: Permite simular erros de bits no arquivo codificado para testar a funcionalidade de correção de erros.

## Uso

### Codificar um Arquivo

Para codificar um arquivo, execute o seguinte comando:

```bash
go run main.go -c <ARQUIVO>
```

Isso criará um arquivo `<ARQUIVO>.hamming`, que será significativamente maior do que o arquivo original.

### Decodificar um Arquivo

Para decodificar um arquivo codificado com Hamming, use o comando:

```bash
go run main.go -d <ARQUIVO>.hamming
```

Isso gerará um arquivo `<ARQUIVO>.dec` que será idêntico ao arquivo original.

## Notas Importantes

* Modificar o arquivo `.hamming` diretamente em um editor de código pode introduzir bytes adicionais, já que alguns editores podem escrever um caractere de 2 bytes em vez de 1.
* Este programa foi projetado para codificação Unix, portanto codificar um arquivo com codificação do Windows pode gerar pequenas alterações.

## Testando a Correção de Erros

Para testar as capacidades de correção de erros, siga os passos abaixo:

1. Acesse o diretório `functions`:

   ```bash
   cd functions/
   ```

2. Abra o arquivo `functions_test.go` e localize a função `MakeMistakes`. Modifique o arquivo de origem e defina o número de alterações desejadas na variável `x`.

3. Execute o teste com o seguinte comando:

   ```bash
   go test -run TestMakeMistakes
   ```

Isso selecionará aleatoriamente `x` sequências do arquivo `.hamming` e alterará um bit aleatório em cada uma delas.

## Conclusão

Este Codificador e Decodificador de Hamming oferece uma maneira simples e eficaz de codificar e decodificar arquivos, além de permitir testes de correção de erros. Sinta-se à vontade para modificar e aprimorar o código conforme necessário!

