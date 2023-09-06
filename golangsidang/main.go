package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"

	"gorgonia.org/gorgonia"
	"gorgonia.org/tensor"
)

type DataSet struct {
	ContentType string
	Prompt      string
	Completion  string
}

func main() {
	rawDataset, err := readDataFromFile("iriss.csv") // Ganti nama file sesuai kebutuhan
	if err != nil {
		log.Fatalf("Gagal membaca data dari file: %v", err)
	}

	cleanedDataset := cleanDataSet(rawDataset)

	if err := saveDataSetToCSV("cleaned_dataset.csv", cleanedDataset); err != nil {
		log.Fatalf("Gagal menyimpan dataset yang telah dibersihkan ke CSV: %v", err)
	}

	// Example: Adding TensorFlow-like functionality
	data := []float64{1.0, 2.0, 3.0, 4.0, 5.0}
	tensorData := tensor.New(tensor.Of(tensor.Float64), tensor.WithShape(1, 5), tensor.WithBacking(data))

	// Create a Gorgonia expression for the mean
	g := gorgonia.NewGraph()
	x := gorgonia.NewMatrix(g, tensor.Float64, gorgonia.WithShape(1, 5), gorgonia.WithValue(tensorData))
	mean := gorgonia.Must(gorgonia.Mean(x))

	// Create a VM and run the computation graph
	machine := gorgonia.NewTapeMachine(g)
	if err := machine.RunAll(); err != nil {
		log.Fatalf("Error running computation graph: %v", err)
	}

	// Get the mean result
	meanResult := mean.Value().Data().(float64)

	// Calculate the mean percentage
	// Calculate the mean percentage
	meanPercentage := meanResult * 100

	// Batasi nilai maksimum menjadi 100%
	if meanPercentage > 100 {
		meanPercentage = meanResult
	}

	fmt.Printf("Mean Value: %.2f (as a percentage: %.2f%%)\n", meanResult, meanPercentage)

	// Visualize the mean value
	visualizeMean(meanResult)

	fmt.Println("Dataset has been cleaned, saved, and visualized successfully.")
}

// visualizeMean adalah fungsi sederhana untuk menampilkan mean value dalam bentuk tanda plus
func visualizeMean(mean float64) {
	fmt.Println("Visualisasi Mean Value:")
	for i := 0; i < int(mean); i++ {
		fmt.Print("+")
	}
	fmt.Println()
}

// (Sisanya tetap seperti kode sebelumnya)

func readDataFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("Gagal membuka file: %v", err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file) // Inisialisasi scanner di sini
	for scanner.Scan() {
		line := scanner.Text()
		if !shouldSkipLine(line) {
			lines = append(lines, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("Gagal membaca baris dari file: %v", err)
	}

	return lines, nil
}

func shouldSkipLine(line string) bool {
	keywordsToSkip := []string{"_id", "import", "package", "{", "}", "func", "main()", "()"}
	for _, keyword := range keywordsToSkip {
		if strings.Contains(line, keyword) {
			return true
		}
	}
	return false
}

func cleanDataSet(rawData []string) []DataSet {
	var cleanedData []DataSet

	for _, record := range rawData {
		parts := strings.Split(record, ",")
		if len(parts) == 3 {
			contentType := strings.TrimSpace(parts[0])
			prompt := strings.TrimSpace(parts[1])
			completion := strings.TrimSpace(parts[2])

			if contentType != "" && prompt != "" && completion != "" {
				cleanedData = append(cleanedData, DataSet{
					ContentType: contentType,
					Prompt:      prompt,
					Completion:  completion,
				})
			}
		}
	}

	return cleanedData
}

func saveDataSetToCSV(filename string, data []DataSet) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("Gagal membuat file CSV: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Tulis header
	if err := writer.Write([]string{"content_type", "prompts", "completion"}); err != nil {
		return fmt.Errorf("Gagal menulis header CSV: %v", err)
	}

	// Tulis baris data
	for _, entry := range data {
		if err := writer.Write([]string{entry.ContentType, entry.Prompt, entry.Completion}); err != nil {
			return fmt.Errorf("Gagal menulis catatan CSV: %v", err)
		}
	}

	return nil
}
