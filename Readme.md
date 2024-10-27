# Usuários

- ## Criar usuário:
  Nome, email, senha -> enviar email do código de verificação pelo email e salvar em cache por 10min
  -> Se email for verificado, enviar JWT para resposta e logar usuário, senão for expirar código de verificação.
- ## Logar usuário:
  Email e senha -> verificar no banco de dados se existe um usuário com o email e se a senha corresponde -> Enviar JWT de autenticação e logar usuário.
- ## Ler dados:
  Receber requisição com jwt.
- ## Deletar dados:
  Receber requisição com jwt e deletar todos os dados referentes ao id.
- ## Atualizar dados:
  Receber requisição jwt e corpo com os novos dados.

# Plantas

- Procurar plantas por nome.
- Procurar planta por espécie.
- Procurar planta por último acesso.
- Criar planta.
- Deletar planta.
- Atualizar planta.
- Ler dados da planta.
- Gerar relatório da planta.
- Gerar recomendações.
- Criar categoria para planta.

# Hortas

- Todas as atividades da planta

# Tarefas

- Buscar por horta.
- Buscar por prazo.
- Buscar por tópico.
- Criar tópico.
