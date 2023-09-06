package main

import (
	"bytes"
	"context"
	"encoding/csv"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/fogleman/gg" // Import modul "gg"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/websocket/v2"
	"github.com/jdkato/prose/v2"
	"github.com/sashabaranov/go-openai"
)

type DataSet struct {
	Prompt string
	Output string
}

var dataSet []DataSet

func main() {
	if err := readDataSetFromCSV("iriss.csv"); err != nil {
		log.Fatalf("Failed to read dataset from CSV: %v", err)
	}
	app := fiber.New()
	app.Use(logger.New())

	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				log.Println("Websocket read error:", err.Error())
				break
			}

			// Convert the WebSocket message (byte slice) to a string prompt
			prompt := string(msg)

			go func() {
				cli := openai.NewClient("ENV")

				var output string
				for _, entry := range dataSet {
					if entry.Prompt == prompt {
						output = entry.Output
						break
					}
				}

				if output == "" {
					stream, err := cli.CreateCompletionStream(
						context.Background(),
						openai.CompletionRequest{
							Model:       openai.GPT3TextDavinci003,
							Prompt:      prompt,
							MaxTokens:   2048,
							Temperature: 0.5,
							// ... (other parameters)
						},
					)
					if err != nil {
						log.Println("Error:", err.Error())
						return
					}

					var response strings.Builder
					for {
						resp, err := stream.Recv()
						if err != nil {
							log.Println("Error:", err.Error())
							return
						}

						if len(resp.Choices) > 0 {
							chunk := resp.Choices[0].Text
							time.Sleep(500 * time.Millisecond)
							err = c.WriteMessage(websocket.TextMessage, []byte(chunk))
							if err != nil {
								log.Println("Websocket write error:", err.Error())
								return
							}
							response.WriteString(chunk)
						}

						if resp.Choices[0].FinishReason == "stop" {
							break
						}
					}

					output = response.String()
				}

				// Memproses prompt dengan fungsi NLP
				words := strings.Fields(prompt)
				if err != nil {
					log.Println("Tokenization error:", err.Error())
					return
				}
				// ...

				// Memproses prompt dengan fungsi NLP
				// ...

				// Memproses prompt dengan fungsi NLP
				var tokens []prose.Token
				for _, word := range words {
					token := prose.Token{
						Text: word,
						Tag:  "", // You can set the tag as needed
					}
					tokens = append(tokens, token)
				}

				processedTokens := YourNLPAlgorithm(tokens)

				// ...

				processedText := JoinTokens(processedTokens)

				// Menggabungkan teks hasil NLP dengan output dari model AI
				output = processedText + " " + output
				output = addImageHTML(output)
				// Mengirimkan hasil ke WebSocket
				err = c.WriteMessage(websocket.TextMessage, []byte(output))
				if err != nil {
					log.Println("Websocket write error:", err.Error())
					return
				}

				// Menggabungkan teks hasil NLP dengan output dari model AI

				// Membuat gambar dengan modul "gg"
				const (
					Width  = 800
					Height = 400
				)
				dc := gg.NewContext(Width, Height)

				// Set warna latar belakang
				dc.SetRGB(1, 1, 1)
				dc.Clear()

				// Set warna teks
				dc.SetRGB(0, 0, 0)

				// Gambar teks hasil dari respon chatan AI
				dc.DrawStringAnchored(output, Width/2, Height/2, 0.5, 0.5)

				// Simpan gambar ke file
				if err := dc.SavePNG("hasil_nlp.png"); err != nil {
					log.Println("Error saving image:", err.Error())
				}

			}()
		}
	}))

	// app.Get("/gambar", func(c *fiber.Ctx) error {
	// 	// Membaca gambar yang telah dihasilkan
	// 	imageBytes, err := os.ReadFile("hasil_nlp.png")
	// 	if err != nil {
	// 		return c.Status(fiber.StatusInternalServerError).SendString("Error reading image")
	// 	}

	// 	// Mengirim gambar sebagai respons dengan tipe konten "image/png"
	// 	return c.Type("image/png").Send(imageBytes)
	// })
	log.Println("Server running on port 8080")
	log.Fatal(app.Listen(":8080"))
}

func readDataSetFromCSV(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	// Skip header row if present
	if len(records) > 0 {
		records = records[1:]
	}

	dataSet = make([]DataSet, len(records))
	for i, record := range records {
		dataSet[i] = DataSet{
			Prompt: record[0],
			Output: record[1],
		}
	}
	return nil
}

// Fungsi NLP Anda
func YourNLPAlgorithm(tokens []prose.Token) []prose.Token {
	var processedTokens []prose.Token
	wordPattern := regexp.MustCompile(`^[a-zA-Z]+$`) // Pola Regex untuk cocokkan kata-kata berhuruf
	for _, token := range tokens {
		if wordPattern.MatchString(token.Text) {
			// Memproses token lebih lanjut (terapkan logika Anda sendiri)
			// Contoh: mengonversi ke huruf kecil
			processedToken := prose.Token{
				Text: strings.ToLower(token.Text),
				Tag:  token.Tag,
			}
			processedTokens = append(processedTokens, processedToken)
		}
	}
	return processedTokens
}

// Gabungkan teks hasil NLP
func JoinTokens(tokens []prose.Token) string {
	var buffer bytes.Buffer
	for i, token := range tokens {
		if i > 0 {
			buffer.WriteString(" ")
		}
		buffer.WriteString(token.Text)
	}
	return buffer.String()
}
func addImageHTML(text string) string {
	// Menggunakan ekspresi reguler untuk mencari tautan gambar dalam format ![Gorilla](https://example.com/gorilla.jpg)
	linkPattern := regexp.MustCompile(`!\[.*\]\((.*\.(jpg|png))\)`)
	matches := linkPattern.FindStringSubmatch(text)
	if len(matches) > 1 {
		log.Println("Image URL:", matches[1])
		// Jika tautan gambar ditemukan, ekstrak URL gambar
		imageURL := matches[1]

		// Membuat elemen HTML <img> dengan atribut src yang mengarah ke URL gambar
		imageHTML := "<img src='" + imageURL + "'/>"
		log.Println("Image HTML:", imageHTML)
		// Gantikan tautan gambar dengan elemen HTML dalam teks
		textWithImage := linkPattern.ReplaceAllString(text, imageHTML)

		return textWithImage
	}
	// Jika tidak ada tautan gambar yang ditemukan, kembalikan teks asli
	log.Println(text, "text")
	return text
}
