package analytics

import (
	"sort"

	"github.com/peder1981/rede-cnpj/RedeGO/internal/models"
)

// GraphStats estatísticas do grafo
type GraphStats struct {
	TotalNodes       int                `json:"totalNodes"`
	TotalEdges       int                `json:"totalEdges"`
	Empresas         int                `json:"empresas"`
	Pessoas          int                `json:"pessoas"`
	PessoasExternas  int                `json:"pessoasExternas"`
	Densidade        float64            `json:"densidade"`
	GrauMedio        float64            `json:"grauMedio"`
	NosMaisConectados []NodeDegree      `json:"nosMaisConectados"`
	TiposRelacao     map[string]int     `json:"tiposRelacao"`
	ComponentesConexos int              `json:"componentesConexos"`
}

// NodeDegree grau de um nó
type NodeDegree struct {
	ID     string  `json:"id"`
	Label  string  `json:"label"`
	Degree int     `json:"degree"`
}

// Analyzer analisador de grafos
type Analyzer struct{}

// NewAnalyzer cria um novo analisador
func NewAnalyzer() *Analyzer {
	return &Analyzer{}
}

// AnalyzeGraph analisa um grafo e retorna estatísticas
func (a *Analyzer) AnalyzeGraph(graph *models.Graph) *GraphStats {
	stats := &GraphStats{
		TotalNodes:     len(graph.Nodes),
		TotalEdges:     len(graph.Edges),
		TiposRelacao:   make(map[string]int),
	}

	// Conta tipos de nós
	for _, node := range graph.Nodes {
		if len(node.ID) >= 3 {
			prefix := node.ID[:3]
			switch prefix {
			case "PJ_":
				stats.Empresas++
			case "PF_":
				stats.Pessoas++
			case "PE_":
				stats.PessoasExternas++
			}
		}
	}

	// Conta tipos de relação
	for _, edge := range graph.Edges {
		stats.TiposRelacao[edge.Label]++
	}

	// Calcula densidade
	if stats.TotalNodes > 1 {
		maxEdges := stats.TotalNodes * (stats.TotalNodes - 1)
		stats.Densidade = float64(stats.TotalEdges) / float64(maxEdges)
	}

	// Calcula grau médio
	degrees := a.calculateDegrees(graph)
	totalDegree := 0
	for _, degree := range degrees {
		totalDegree += degree.Degree
	}
	if stats.TotalNodes > 0 {
		stats.GrauMedio = float64(totalDegree) / float64(stats.TotalNodes)
	}

	// Nós mais conectados (top 10)
	sort.Slice(degrees, func(i, j int) bool {
		return degrees[i].Degree > degrees[j].Degree
	})
	if len(degrees) > 10 {
		stats.NosMaisConectados = degrees[:10]
	} else {
		stats.NosMaisConectados = degrees
	}

	// Componentes conexos
	stats.ComponentesConexos = a.countConnectedComponents(graph)

	return stats
}

// calculateDegrees calcula o grau de cada nó
func (a *Analyzer) calculateDegrees(graph *models.Graph) []NodeDegree {
	degreeMap := make(map[string]int)
	labelMap := make(map[string]string)

	// Inicializa com todos os nós
	for _, node := range graph.Nodes {
		degreeMap[node.ID] = 0
		labelMap[node.ID] = node.Label
	}

	// Conta graus
	for _, edge := range graph.Edges {
		degreeMap[edge.From]++
		degreeMap[edge.To]++
	}

	// Converte para slice
	degrees := make([]NodeDegree, 0, len(degreeMap))
	for id, degree := range degreeMap {
		degrees = append(degrees, NodeDegree{
			ID:     id,
			Label:  labelMap[id],
			Degree: degree,
		})
	}

	return degrees
}

// countConnectedComponents conta componentes conexos usando DFS
func (a *Analyzer) countConnectedComponents(graph *models.Graph) int {
	visited := make(map[string]bool)
	adjList := make(map[string][]string)

	// Constrói lista de adjacência
	for _, edge := range graph.Edges {
		adjList[edge.From] = append(adjList[edge.From], edge.To)
		adjList[edge.To] = append(adjList[edge.To], edge.From)
	}

	count := 0
	for _, node := range graph.Nodes {
		if !visited[node.ID] {
			a.dfs(node.ID, visited, adjList)
			count++
		}
	}

	return count
}

// dfs busca em profundidade
func (a *Analyzer) dfs(nodeID string, visited map[string]bool, adjList map[string][]string) {
	visited[nodeID] = true
	
	for _, neighbor := range adjList[nodeID] {
		if !visited[neighbor] {
			a.dfs(neighbor, visited, adjList)
		}
	}
}

// DetectCentralNodes detecta nós centrais (betweenness centrality simplificado)
func (a *Analyzer) DetectCentralNodes(graph *models.Graph, top int) []NodeDegree {
	// Por simplicidade, usa grau como proxy para centralidade
	degrees := a.calculateDegrees(graph)
	
	sort.Slice(degrees, func(i, j int) bool {
		return degrees[i].Degree > degrees[j].Degree
	})

	if len(degrees) > top {
		return degrees[:top]
	}
	return degrees
}

// DetectCommunities detecta comunidades (algoritmo simples baseado em conexões)
func (a *Analyzer) DetectCommunities(graph *models.Graph) map[int][]string {
	visited := make(map[string]bool)
	adjList := make(map[string][]string)

	// Constrói lista de adjacência
	for _, edge := range graph.Edges {
		adjList[edge.From] = append(adjList[edge.From], edge.To)
		adjList[edge.To] = append(adjList[edge.To], edge.From)
	}

	communities := make(map[int][]string)
	communityID := 0

	for _, node := range graph.Nodes {
		if !visited[node.ID] {
			community := []string{}
			a.dfsCollect(node.ID, visited, adjList, &community)
			communities[communityID] = community
			communityID++
		}
	}

	return communities
}

// dfsCollect coleta nós de uma comunidade
func (a *Analyzer) dfsCollect(nodeID string, visited map[string]bool, adjList map[string][]string, community *[]string) {
	visited[nodeID] = true
	*community = append(*community, nodeID)
	
	for _, neighbor := range adjList[nodeID] {
		if !visited[neighbor] {
			a.dfsCollect(neighbor, visited, adjList, community)
		}
	}
}

// CalculateShortestPath calcula caminho mais curto entre dois nós (BFS)
func (a *Analyzer) CalculateShortestPath(graph *models.Graph, from, to string) []string {
	adjList := make(map[string][]string)
	
	// Constrói lista de adjacência
	for _, edge := range graph.Edges {
		adjList[edge.From] = append(adjList[edge.From], edge.To)
		adjList[edge.To] = append(adjList[edge.To], edge.From)
	}

	// BFS
	queue := []string{from}
	visited := make(map[string]bool)
	parent := make(map[string]string)
	visited[from] = true

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if current == to {
			// Reconstrói caminho
			path := []string{}
			for node := to; node != ""; node = parent[node] {
				path = append([]string{node}, path...)
				if node == from {
					break
				}
			}
			return path
		}

		for _, neighbor := range adjList[current] {
			if !visited[neighbor] {
				visited[neighbor] = true
				parent[neighbor] = current
				queue = append(queue, neighbor)
			}
		}
	}

	return nil // Sem caminho
}
