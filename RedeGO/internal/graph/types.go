package graph

import (
	"database/sql"

	"github.com/peder1981/rede-cnpj/RedeGO/internal/models"
	_ "github.com/mattn/go-sqlite3"
)

// GraphType tipo de grafo
type GraphType string

const (
	GraphTypeRede           GraphType = "rede"            // Rede completa
	GraphTypeCaminhos       GraphType = "caminhos"        // Caminhos entre entidades
	GraphTypeCaminhosDireto GraphType = "caminhos-direto" // Apenas caminhos diretos
	GraphTypeCaminhosComum  GraphType = "caminhos-comum"  // Entidades em comum
)

// PathFinder encontra caminhos entre entidades
type PathFinder struct {
	dbPath string
}

// NewPathFinder cria um novo buscador de caminhos
func NewPathFinder(dbPath string) *PathFinder {
	return &PathFinder{dbPath: dbPath}
}

// FindPaths encontra caminhos entre duas entidades
func (p *PathFinder) FindPaths(from, to string, maxDepth int) (*models.Graph, error) {
	db, err := sql.Open("sqlite3", p.dbPath)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Usa BFS para encontrar caminhos
	visited := make(map[string]bool)
	paths := [][]string{}
	
	// Busca caminhos
	p.bfs(db, from, to, maxDepth, []string{from}, visited, &paths)

	// Converte caminhos para grafo
	return p.pathsToGraph(db, paths)
}

// bfs busca em largura para encontrar caminhos
func (p *PathFinder) bfs(db *sql.DB, current, target string, depth int, path []string, visited map[string]bool, paths *[][]string) {
	if depth == 0 {
		return
	}

	if current == target && len(path) > 1 {
		// Encontrou caminho
		pathCopy := make([]string, len(path))
		copy(pathCopy, path)
		*paths = append(*paths, pathCopy)
		return
	}

	visited[current] = true

	// Busca vizinhos
	query := `
		SELECT id2 FROM ligacao WHERE id1 = ?
		UNION
		SELECT id1 FROM ligacao WHERE id2 = ?
	`

	rows, err := db.Query(query, current, current)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var neighbor string
		if err := rows.Scan(&neighbor); err != nil {
			continue
		}

		if !visited[neighbor] {
			newPath := append(path, neighbor)
			p.bfs(db, neighbor, target, depth-1, newPath, visited, paths)
		}
	}

	visited[current] = false
}

// pathsToGraph converte caminhos para grafo
func (p *PathFinder) pathsToGraph(db *sql.DB, paths [][]string) (*models.Graph, error) {
	graph := &models.Graph{
		Nodes: []models.Node{},
		Edges: []models.Edge{},
	}

	nodeSet := make(map[string]bool)
	edgeSet := make(map[string]bool)

	for _, path := range paths {
		for i, nodeID := range path {
			// Adiciona nó
			if !nodeSet[nodeID] {
				node, err := p.getNodeInfo(db, nodeID)
				if err == nil {
					graph.Nodes = append(graph.Nodes, node)
					nodeSet[nodeID] = true
				}
			}

			// Adiciona aresta
			if i < len(path)-1 {
				nextID := path[i+1]
				edgeKey := nodeID + "->" + nextID
				
				if !edgeSet[edgeKey] {
					edge, err := p.getEdgeInfo(db, nodeID, nextID)
					if err == nil {
						graph.Edges = append(graph.Edges, edge)
						edgeSet[edgeKey] = true
					}
				}
			}
		}
	}

	return graph, nil
}

// getNodeInfo obtém informações do nó
func (p *PathFinder) getNodeInfo(db *sql.DB, nodeID string) (models.Node, error) {
	// Busca label do nó
	var label string
	
	// Tenta buscar em id_search
	err := db.QueryRow("SELECT id_descricao FROM id_search WHERE id_descricao LIKE ? LIMIT 1", nodeID+"%").Scan(&label)
	if err != nil {
		// Se não encontrar, usa o próprio ID
		label = nodeID
	}

	return models.Node{
		ID:    nodeID,
		Label: label,
	}, nil
}

// getEdgeInfo obtém informações da aresta
func (p *PathFinder) getEdgeInfo(db *sql.DB, from, to string) (models.Edge, error) {
	var label string
	
	err := db.QueryRow("SELECT descricao FROM ligacao WHERE id1 = ? AND id2 = ? LIMIT 1", from, to).Scan(&label)
	if err != nil {
		label = "relacionamento"
	}

	return models.Edge{
		From:  from,
		To:    to,
		Label: label,
	}, nil
}

// FindCommonEntities encontra entidades em comum entre duas entidades
func (p *PathFinder) FindCommonEntities(id1, id2 string) (*models.Graph, error) {
	db, err := sql.Open("sqlite3", p.dbPath)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Busca entidades em comum
	query := `
		SELECT DISTINCT l1.id2 as comum
		FROM ligacao l1
		JOIN ligacao l2 ON l1.id2 = l2.id2
		WHERE l1.id1 = ? AND l2.id1 = ?
		UNION
		SELECT DISTINCT l1.id1 as comum
		FROM ligacao l1
		JOIN ligacao l2 ON l1.id1 = l2.id1
		WHERE l1.id2 = ? AND l2.id2 = ?
	`

	rows, err := db.Query(query, id1, id2, id1, id2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	graph := &models.Graph{
		Nodes: []models.Node{},
		Edges: []models.Edge{},
	}

	// Adiciona nós principais
	node1, _ := p.getNodeInfo(db, id1)
	node2, _ := p.getNodeInfo(db, id2)
	graph.Nodes = append(graph.Nodes, node1, node2)

	// Adiciona entidades em comum
	for rows.Next() {
		var commonID string
		if err := rows.Scan(&commonID); err != nil {
			continue
		}

		commonNode, err := p.getNodeInfo(db, commonID)
		if err != nil {
			continue
		}
		graph.Nodes = append(graph.Nodes, commonNode)

		// Adiciona arestas
		edge1, _ := p.getEdgeInfo(db, id1, commonID)
		edge2, _ := p.getEdgeInfo(db, id2, commonID)
		graph.Edges = append(graph.Edges, edge1, edge2)
	}

	return graph, nil
}

// FilterGraph filtra grafo por critérios
type FilterCriteria struct {
	MinConnections int      // Mínimo de conexões
	MaxConnections int      // Máximo de conexões
	NodeTypes      []string // Tipos de nó (PJ_, PF_, PE_)
	EdgeTypes      []string // Tipos de aresta
}

// FilterGraph filtra um grafo existente
func FilterGraph(graph *models.Graph, criteria FilterCriteria) *models.Graph {
	filtered := &models.Graph{
		Nodes: []models.Node{},
		Edges: []models.Edge{},
	}

	// Conta conexões por nó
	connections := make(map[string]int)
	for _, edge := range graph.Edges {
		connections[edge.From]++
		connections[edge.To]++
	}

	// Filtra nós
	validNodes := make(map[string]bool)
	for _, node := range graph.Nodes {
		// Verifica tipo de nó
		if len(criteria.NodeTypes) > 0 {
			valid := false
			for _, t := range criteria.NodeTypes {
				if len(node.ID) >= len(t) && node.ID[:len(t)] == t {
					valid = true
					break
				}
			}
			if !valid {
				continue
			}
		}

		// Verifica conexões
		conn := connections[node.ID]
		if criteria.MinConnections > 0 && conn < criteria.MinConnections {
			continue
		}
		if criteria.MaxConnections > 0 && conn > criteria.MaxConnections {
			continue
		}

		filtered.Nodes = append(filtered.Nodes, node)
		validNodes[node.ID] = true
	}

	// Filtra arestas
	for _, edge := range graph.Edges {
		// Verifica se ambos os nós são válidos
		if !validNodes[edge.From] || !validNodes[edge.To] {
			continue
		}

		// Verifica tipo de aresta
		if len(criteria.EdgeTypes) > 0 {
			valid := false
			for _, t := range criteria.EdgeTypes {
				if edge.Label == t {
					valid = true
					break
				}
			}
			if !valid {
				continue
			}
		}

		filtered.Edges = append(filtered.Edges, edge)
	}

	return filtered
}
