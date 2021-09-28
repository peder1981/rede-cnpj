# REDE-DADOS-CNPJ - Visualização de dados públicos de CNPJ

### Vídeo no youtube<br>
[![youtube](http://img.youtube.com/vi/nxz9Drhqn_I/0.jpg)](https://youtu.be/nxz9Drhqn_I)

<br>Outros vídeos de utilização:<br>
Opção básicas dos botões: https://youtu.be/-Ug6ToTRnE4 <br>
Criar uma ligação no gráfico: https://youtu.be/8I0oNb4U9Rw <br>
Aumentar tamanho da ligação: https://youtu.be/7hy74LE8e7A <br>
Exportar dados como json: https://youtu.be/WKn02G9yHbQ <br>
Arrastar células do Excel: https://youtu.be/Oxze-d4V7kE <br>
A rotina possibilita visualizar de forma gráfica os relacionamentos entre empresas e sócios, a partir da base de dados públicos de cnpj da Receita Federal. <br>
Foi testada nos navegadores Firefox, Edge e Chrome. NÃO FUNCIONA no Internet Explorer. <br>
A base de dados é o arquivo CNPJ_full.db, banco de dados no formato sqlite. Para exemplificar o funcionamento da rotina, este repositório tem o arquivo com cerca de mil registros com dados fictícios de empresas e de sócios. <br>

## Versão online com base completa de dados públicos de CNPJ:
http://168.138.150.250/rede/ <br>
Leia as informações iniciais, e digite "TESTE", CNPJ, Razão Social Completa, Nome Completo de Sócio ou Radical do CNPJ. Pode-se inserir vários CNPJs de uma só vez, separando-os por (;).
Funciona parcialmente em celular, com menu errático.

## Versão em python:
É preciso ter instalado no computador, um interpretador de linguagem python (versão 3.7 ou posterior) como a distribuída pelo Anaconda ou WinPython.<br> 
Para iniciar esse script, em um console DOS digite<br>
python rede.py<br>
A rotina abrirá o endereço http://127.0.0.1:5000/rede/ no navegador padrão.
Se der algum erro como “module <nome do módulo> not found”, instale o módulo pelo comando pip install <nome do módulo>.<br>
As opções por linha de comando são exibidas fazendo python rede.py -h<br>

## Versão executável:
Para iniciar a versão executável, primeiro descompacte o arquivo [rede-cnpj-exe.7z](https://drive.google.com/file/d/1rjsWmzzudmrjACR61HC2f2HwZ75RueCY/view?usp=sharing). Para executar a rotina, clique duas vezes em rede.exe. Obs: a versão executável foi criada por pyinstaller para funcionar no windows. É possível que falte alguma dll para funcionar corretamente.<br>
A rotina abrirá o endereço http://127.0.0.1:5000/rede/ no navegador padrão e um console do DOS. Para parar a execução, feche o console.<br>
Esta versão antiga executável só irá funcionar com a base de testes que está no arquivo compactado.<br>

## Como utilizar o Banco de dados públicos completo de CNPJs:
A pasta contém um arquivo <b>cnpj_teste.db</b>, que é o banco de dados com poucos dados apenas para testar o funcionamento da rotina. Substitua esse arquivo pela base <b>cnpj.db</b> em sqlite que pode ser obtido no Google Drive https://drive.google.com/drive/folders/1Gkeq27aHv6UgT8m30fc4hZWMPqdhEHWr?usp=sharing (base da SRF de 11/09/2021), alterando o arquivo de configuração rede.ini, mudando o nome do banco na linha para<br>
<b>base_receita = cnpj.db</b><br>
O projeto https://github.com/rictom/cnpj-sqlite faz a conversão dos arquivos zipados do site da Receita para sqlite e gera o arquivo <b>cnpj.db</b> que está no Google Drive.<br>

## OBSERVAÇÃO IMPORTANTE em março de 2021:
Até 2020 este script utilizava o arquivo sqlite gerado a partir do script do Fabio Serpa (https://github.com/fabioserpa/CNPJ-full), mas não funciona mais porque o layout das tabelas mudou.<br>

## Opções:

A roda do mouse expande ou diminui o tamanho da exibição.<br>
Fazendo click duplo em um ícone, a rotina expande as ligações.<br>
Apertando SHIFT, é possível selecionar mais de um ícone. <br>
Pressionando CTRL e arrastando na tela, adiciona a seleção os itens da área.
Clicar no botão do meio do mouse (roda) faz aparecer janela para editar uma Nota, que aparece numa terceira linha abaixo do ícone.

Outras opções da rede estão no menu contextual do mouse (botão direito), sendo configuradas teclas de atalho correspondentes aos comandos:
 

## Tecla – Descrição do comando.
- TECLAS de 1 a 9 - Inserir camadas correspondente ao número sobre o nó selecionado;
- I - Inserir CNPJ, Razão Social completa ou nome completo de sócio. Poderão ser colocados vários CNPJs ao mesmo tempo, separados por ponto e vírgula (;).
- U - Criar item novo (que não seja PF ou PJ) e ligar aos itens selecionados;
- E - Editar dados do item (que não seja PF ou PJ) selecionado;
- CRTL+Z – Desfaz Inserção;

- SubMenu Ligar:
- U - Ligar para novo item;
- L - Ligar itens selecionados, ligação tipo estrela (o primeiro ligado aos demais);
- SHIFT+L - Remover ligação entre itens selecionados;
- Remover Ligacoes - Remove todas as ligações dos itens selecionados;
- K - Ligar itens selecionados, ligação tipo fila (o primeiro ligado ao segundo, o segundo ao terceiro, etc);

- SubMenu Visualização:
- A - Gráfico em Nova Aba - Abre aba com os itens selecionados;
- Q - Quebrar o gráfico em abas - Divide o gráfico em partes menores, mantendo as ligações
- P - Fixar o nó na posição;
- SHIFT+P - Desfixar todos os nós do gráfico;
- CTRL+P - Fixa um nó em cada grupo conexo (para evitar que o gráfico se expanda indefinidamente);
- SubMenu Visualização>Rótulos:
- E - Editar rótulo; 
- N - Rótulo - Exibe apenas o primeiro nome;
- SHIFT+N - Oculta/exibe texto da ligação;


- Alterar Ícone;
- C - Colorir os nós selecionados;
- Escolher Cor;
- D – Abre um popup com dados;
- SHIFT+D – Abre numa nova aba com Dados;
- CTRL+D – Lista ids dos itens selecionados;
- Altera o nome da aba;
- Escala Inicial - Coloca a exibição sem zoom, na escala inicial.
- Barra de Espaço - Parar/reiniciar leiaute (se a tela tiver muitos nós, os comandos funcionam melhor se o leiaute estiver parado);

- F - Localizar - Localizar na Tela Nome, CNPJ ou CPF;
- SHIFT+F - Localizar apenas na seleção;
- CTRL+F - Localiza por campo (como cor do item);
- J – Seleciona itens adjacentes;
- SHIFT+J – Seleciona árvores dos itens selecionados;
- CTRL+J - Itens com mais ligações - Opção para selecionar os itens do gráfico com mais ligações;
- Itens ligados a coloridos;
- Grupos com duas cores;
- CTRL+A - Seleciona todos os itens;
- CTRL+SHIFT+A - Inverte seleção;
- SubMenu - Busca em sites:
- G – Abre o nó numa aba do site Google;
- SHIFT+G – Abre o endereço no Google Maps (só CNPJs);
- Jusbrasil - Busca no site Jusbrasil
- Portal da Transparência - Busca no Portal da Transparência;

- SubMenu Salvar/Abrir:

- Salvar dados em Excel;
- Salvar imagem em formato SVG;
- Salvar Arquivo Json - salva dados do gráfico no formato json;
- Abrir Arquivo Json;
- Exportar/Importar JSON ao Servidor - Exportar ou importar dados do gráfico em formato JSON carregados no servidor;
- Banco de Dados - Exporta dados para banco de dados sqlite (só para usuário local);

- SubMenu Excluir

- DEL – Excluir itens selecionados.
- SHIFT+DEL – Excluir todos os itens.
- Excluir Nó mantendo Link;
- Simplifica Gráfico - Remove itens nas bordas do gráfico que não tenham destaque;
- Excluir itens isolados - Remove itens sem ligação.

Os comandos valem para o último nó selecionado ou nós selecionados, que ficam em destaque com a animação no contorno ods ícones.
Pressionando SHIFT e click, é possível selecionar mais de um ícone para fazer Exclusão ou para Expansão de vínculos.
Pode-se arrastar células com listas de CNPJs do Excel para a janela, ou arrastar arquivos csv ou json.

## Fonte dos dados:

Base de CNPJ. A base de dados públicos de CNPJ da Receita Federal tem informação de Capital Social de empresas. A tabela de sócios contém apenas os sócios ativos de empresas, com CPF descaracterizado e nome completo do sócio.<br>
https://www.gov.br/receitafederal/pt-br/assuntos/orientacao-tributaria/cadastros/consultas/dados-publicos-cnpj ou ou http://200.152.38.155/CNPJ/ em formato csv<br>
Arquivo CNPJ.db completo, referência 11/09/2021, já no formato sqlite, dividido em blocos, foi copiado no Google Drive:<br>
https://drive.google.com/drive/folders/1Gkeq27aHv6UgT8m30fc4hZWMPqdhEHWr?usp=sharing <br>
Para juntar os blocos, abra o primeiro (CNPJ.7z.001) no 7zip. Os arquivos compactados têm o tamanho de 5GB. O arquivo descompactado tem 24GB.<br>

## Outras referências:
Biblioteca em javascript para visualização:<br>
https://github.com/anvaka/VivaGraphJS<br>

Menu Contextual:<br>
https://www.cssscript.com/beautiful-multi-level-context-menu-with-pure-javascript-and-css3/

## Histórico de versões
versão 0.7.3 (setembro/2021)
- correção de erro em ligação em banco de dados local.

versão 0.7.2 (setembro/2021)
- opção para alterar o nome da aba;
- opção para selecionar todos os itens;
- opção para inverter seleção;

versão 0.7.1 (agosto/2021)
- opção para dividir gráficos em outras abas (tecla Q);
- correção de erro quando se apertava tecla CTRL;
- opção para selecionar itens adjacentes aos selecionados (tecla J);
- opção para selecionar árvore que contém o item (tecla SHIFT+J);
- opção para listar itens com mais links (tecla CTRL+J)
- opção para selecionar itens com mais ligações para ícones coloridos;
- opção "Nova Aba" (tecla A) abre nova aba com mais de um item selecionado;
- troca de nomes no menu de exportar json para salvar json;
- opção para exportar dados para banco de dados local (só funciona na máquina local);
- opção para exportar para json apenas itens selecionados;
- opção para ocultar rótulos de ligações (SHIFT+N);
- opção para remover todas as ligações dos itens selecionados;
- opção simplifica gráfico (remove itens que não são coloridos ou com comentário que tem apenas uma ligação).

versão 0.6.3 (julho/2021)
- melhoria para dar clique duplo em ícones;
- correção de erro de ligação para empresa no exterior sem cnpj;
- somente o ícone pode ser clicado;
- mensagem de alerta para utilizar caractere curinga;
- mudança nas tabelas temporárias;
- todas as tabelas de códigos (cnae, natureza jurídica, etc) foram incorporados ao arquivo sqlite;
- OBSERVAÇÃO. A versão 0.6.3 só vai funcionar com a versão mais atualizada do arquivo cnpj.db referência 16/7/2021.

versão 0.5.1 (junho/2021)
- atualização da tabela sqlite cnpj.db com dados públicos de 18/06/2021.

versão 0.5 (abril/2021)
- alteração do código para layout novo das tabelas;
- busca por Radical de CNPJ ou CPF de sócio (busca somente pelo miolo do CPF);

versão 0.4 (janeiro/2021)
- usando lock para evitar erro de consulta em requisições simultâneas;
- opção para fazer busca do termo no Portal da Transparência da CGU;
- correção de link para google search.

versão 0.3.4 (janeiro/2021)
- Possibilita ver o texto do lado direito do ícone;
- diagramas de tabela hierárquica;
- ver diagramas de arquivo com código em python;
- mais opções por linha de comando.

versão 0.3 (janeiro/2021)
- Opção para inserção de novos itens para elaboração de mapas mentais;
- Opções para inserir itens novos como link para sites e arquivos locais.
- Opção para arrastar células do excel, leitura de arquivo csv;
- Opções de leitura de entrada por linha de comando;
- Itens selecionados ficam em destaque com linha animada;
- Alteração no formato do arquivo de configuração rede.ini.

versão 0.2 (dezembro/2020)
- Suporte para busca por parte do nome na base de empresas;
- Exportação/importação de gráfico no formato json para o servidor.

versão 0.1 (setembro/2020)
- Primeira versão
