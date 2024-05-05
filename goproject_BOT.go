package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	//"telegram-bot-api"
	//	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	//	tgbotapi "C:\GO_projects\GO\pkg\mod\github.com\go-telegram-bot-api\telegram-bot-api\v5@v5.5.1"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	//  tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	//	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5@latest"
)

const (
	n_nomer = 3
	text_1  = "Выберите одну з 2-хкоманд:\n1) введите 3-х значный номер  для получения отзыва о нем\n2) введите 3-х значный номер и затем свой отзыв внутри звездочек (*...*), пример: 111*не бери*"
	text_2  = "команда введена не корректно"
	text_3  = "отзыв о даном номере отсутствует"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("6731538180:AAFn14QP4Xg7hpdhZSZc3RON6ovYvxgr2sQ")
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	// Setup long-polling request
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // Есть новое сообщение
			text := update.Message.Text      // Текст сообщения
			chatID := update.Message.Chat.ID //  ID чата
			userID := update.Message.From.ID // ID пользователя
			var replyMsg string

			log.Printf("[%s](%d) %s", update.Message.From.UserName, userID, text)

			// Анализируем текст сообщения и записываем ответ в переменную

			replyMsg = otvet(text)
			// Отправляем ответ
			msg := tgbotapi.NewMessage(chatID, replyMsg)    // Создаем новое сообщение
			msg.ReplyToMessageID = update.Message.MessageID // Указываем сообщение, на которое нужно ответить

			bot.Send(msg)
		}
	}
}
func otvet(text_vvod string) string {
	var textout string
	fmt.Println("введи текст")
	//	text_vvod = "111*gtr*"
	file := "tmp/Telephone.csv"
	if len(text_vvod) <= n_nomer {
		textout = output(file, text_vvod)
		//	fmt.Println(textout)
	} else {
		textout = input(file, text_vvod, nil)
	}
	return textout
}

func readFile(filePath string) ([][]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	return records, nil
}

func writeFile(filePath string, records [][]string) error {
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	csvWriter := csv.NewWriter(f)
	err = csvWriter.WriteAll(records)
	if err != nil {
		return err
	}

	return nil
}

func output(file string, text_vvod string) string {

	if len(text_vvod) < n_nomer {
		return text_2 + "\n" + text_1
	}
	//+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	var textout string
	var text_search string
	record, err := readFile(file)
	if err != nil {
		log.Fatal(err) // может return nil, err
	}
	if len(record) == 0 {
		return "записи отсутствуют"
	}
	fmt.Println(record)

	for i, v := range text_vvod {
		if i < n_nomer {
			text_search = text_search + string(v)
		} else {
			break
		}
	}
	fmt.Println(text_search)

	A_out := make([][]string, 0, 10)
	for i := range record {
		if record[i][0] == text_search {
			A_out = append(A_out, record[i])
		}
	}

	if len(A_out) == 0 {
		return text_3
	}

	for _, v := range A_out {
		for _, h := range v {
			textout = textout + " " + h
		}
		textout = textout + "\n"
	}

	return textout
}

func input(file string, text_vvod string, records [][]string) string {
	var x_1 string
	var x_2 string
	for i, v := range text_vvod {
		if i == n_nomer {
			x_1 = string(v)
		}
		if i == len(text_vvod)-1 {
			x_2 = string(v)
		}
	}
	if x_1 != "*" || x_2 != "*" {
		return text_2 + "\n" + text_1
	}

	//+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	record_write := make([]string, 2)
	for i, v := range text_vvod {
		if i < n_nomer {
			record_write[0] = record_write[0] + string(v)
		}
		if i > n_nomer && i < len(text_vvod)-1 {
			record_write[1] = record_write[1] + string(v)
		}
	}

	records = append(records, record_write)
	err := writeFile(file, records)
	if err != nil {
		//	log.Fatal(err) // может return nil, err
		return "номер не добавлен"
	}
	return "номер успешно добавлен"
}
