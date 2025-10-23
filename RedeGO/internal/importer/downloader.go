package importer

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
)

// Downloader gerencia o download dos arquivos da Receita Federal
type Downloader struct {
	baseURL string
	zipDir  string
	client  *http.Client
}

// NewDownloader cria um novo downloader
func NewDownloader(baseURL, zipDir string) *Downloader {
	return &Downloader{
		baseURL: baseURL,
		zipDir:  zipDir,
		client: &http.Client{
			Timeout: 30 * time.Minute,
		},
	}
}

// Download baixa todos os arquivos ZIP
func (d *Downloader) Download() error {
	// Cria diretório se não existir
	if err := os.MkdirAll(d.zipDir, 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório %s: %w", d.zipDir, err)
	}

	// Busca última referência (pasta mais recente)
	latestRef, err := d.getLatestReference()
	if err != nil {
		return fmt.Errorf("erro ao buscar última referência: %w", err)
	}

	fmt.Printf("📅 Última referência encontrada: %s\n", latestRef)
	
	// Lista arquivos ZIP disponíveis
	url := d.baseURL + latestRef
	files, err := d.listZipFiles(url)
	if err != nil {
		return fmt.Errorf("erro ao listar arquivos: %w", err)
	}

	fmt.Printf("📋 Encontrados %d arquivos ZIP\n\n", len(files))

	// Baixa arquivos em paralelo (máximo 5 simultâneos)
	return d.downloadParallel(files, 5)
}

// getLatestReference busca a pasta mais recente no site da Receita
func (d *Downloader) getLatestReference() (string, error) {
	resp, err := d.client.Get(d.baseURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return "", err
	}

	var folders []string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" && strings.HasPrefix(attr.Val, "20") && strings.HasSuffix(attr.Val, "/") {
					folders = append(folders, attr.Val)
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	if len(folders) == 0 {
		return "", fmt.Errorf("nenhuma pasta encontrada em %s", d.baseURL)
	}

	// Retorna a última (mais recente)
	latest := folders[len(folders)-1]
	return latest, nil
}

// listZipFiles lista todos os arquivos .zip na URL
func (d *Downloader) listZipFiles(url string) ([]string, error) {
	resp, err := d.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	var files []string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" && strings.HasSuffix(attr.Val, ".zip") {
					fileURL := attr.Val
					if !strings.HasPrefix(fileURL, "http") {
						fileURL = url + fileURL
					}
					files = append(files, fileURL)
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return files, nil
}

// downloadParallel baixa arquivos em paralelo
func (d *Downloader) downloadParallel(urls []string, maxConcurrent int) error {
	sem := make(chan struct{}, maxConcurrent)
	var wg sync.WaitGroup
	errChan := make(chan error, len(urls))

	for i, url := range urls {
		wg.Add(1)
		go func(idx int, fileURL string) {
			defer wg.Done()
			sem <- struct{}{}        // Adquire semáforo
			defer func() { <-sem }() // Libera semáforo

			filename := filepath.Base(fileURL)
			destPath := filepath.Join(d.zipDir, filename)

			// Verifica se já existe
			if _, err := os.Stat(destPath); err == nil {
				fmt.Printf("[%d/%d] ⏭️  %s (já existe)\n", idx+1, len(urls), filename)
				return
			}

			fmt.Printf("[%d/%d] ⬇️  Baixando %s...\n", idx+1, len(urls), filename)
			
			if err := d.downloadFile(fileURL, destPath); err != nil {
				errChan <- fmt.Errorf("erro ao baixar %s: %w", filename, err)
				return
			}

			fmt.Printf("[%d/%d] ✅ %s concluído\n", idx+1, len(urls), filename)
		}(i, url)
	}

	wg.Wait()
	close(errChan)

	// Verifica se houve erros
	if len(errChan) > 0 {
		return <-errChan
	}

	return nil
}

// downloadFile baixa um arquivo individual
func (d *Downloader) downloadFile(url, destPath string) error {
	resp, err := d.client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code: %d", resp.StatusCode)
	}

	out, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
