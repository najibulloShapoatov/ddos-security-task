package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"strings"
)

var quotes = []string{
	"The best way to predict the future is to invent it. – Alan Kay",
	"Life is 10% what happens to us and 90% how we react to it. – Charles R. Swindoll",
	"The only way to do great work is to love what you do. – Steve Jobs",
	"It does not matter how slowly you go as long as you do not stop. – Confucius",
	"Success is not final, failure is not fatal: It is the courage to continue that counts. – Winston Churchill",
	"The only place where success comes before work is in the dictionary. – Vidal Sassoon",
	"Believe you can and you're halfway there. – Theodore Roosevelt",
	"I can't change the direction of the wind, but I can adjust my sails to always reach my destination. – Jimmy Dean",
	"Whether you think you can or you think you can’t, you’re right. – Henry Ford",
	"I have not failed. I've just found 10,000 ways that won't work. – Thomas A. Edison",
	"The only limit to our realization of tomorrow is our doubts of today. – Franklin D. Roosevelt",
	"To handle yourself, use your head; to handle others, use your heart. – Eleanor Roosevelt",
	"Quality is not an act, it is a habit. – Aristotle",
	"The mind is everything. What you think you become. – Buddha",
	"The best revenge is massive success. – Frank Sinatra",
	"Life shrinks or expands in proportion to one's courage. – Anaïs Nin",
	"Dream big and dare to fail. – Norman Vaughan",
	"Strive not to be a success, but rather to be of value. – Albert Einstein",
	"Do not go where the path may lead, go instead where there is no path and leave a trail. – Ralph Waldo Emerson",
	"Do what you can, with what you have, where you are. – Theodore Roosevelt",
	"You miss 100% of the shots you don’t take. – Wayne Gretzky",
	"I am not a product of my circumstances. I am a product of my decisions. – Stephen Covey",
	"The most difficult thing is the decision to act, the rest is merely tenacity. – Amelia Earhart",
	"The best time to plant a tree was 20 years ago. The second best time is now. – Chinese Proverb",
	"Only a life lived for others is a life worthwhile. – Albert Einstein",
	"An unexamined life is not worth living. – Socrates",
	"Happiness is not something readymade. It comes from your own actions. – Dalai Lama",
	"The only true wisdom is in knowing you know nothing. – Socrates",
	"Believe and act as if it were impossible to fail. – Charles Kettering",
	"The only thing necessary for the triumph of evil is for good men to do nothing. – Edmund Burke",
	"If you want to lift yourself up, lift up someone else. – Booker T. Washington",
	"The best and most beautiful things in the world cannot be seen or even touched - they must be felt with the heart. – Helen Keller",
	"It is during our darkest moments that we must focus to see the light. – Aristotle Onassis",
	"The purpose of our lives is to be happy. – Dalai Lama",
	"In the end, it's not the years in your life that count. It's the life in your years. – Abraham Lincoln",
	"You only live once, but if you do it right, once is enough. – Mae West",
	"Many of life's failures are people who did not realize how close they were to success when they gave up. – Thomas A. Edison",
	"Success is not how high you have climbed, but how you make a positive difference to the world. – Roy T. Bennett",
	"Your time is limited, don’t waste it living someone else’s life. – Steve Jobs",
	"The best dreams happen when you’re awake. – Cherie Gilderbloom",
	"The only thing worse than being blind is having sight but no vision. – Helen Keller",
	"Keep your face always toward the sunshine—and shadows will fall behind you. – Walt Whitman",
	"Success is going from failure to failure without losing your enthusiasm. – Winston Churchill",
	"The difference between ordinary and extraordinary is that little extra. – Jimmy Johnson",
	"The best way to find yourself is to lose yourself in the service of others. – Mahatma Gandhi",
	"What we achieve inwardly will change outer reality. – Plutarch",
	"The only real mistake is the one from which we learn nothing. – Henry Ford",
	"The journey of a thousand miles begins with one step. – Lao Tzu",
	"You must be the change you wish to see in the world. – Mahatma Gandhi",
	"Life is what happens when you’re busy making other plans. – John Lennon",
}

var Quotes []string

func loadQuotes(dataDir string) {
	files, err := os.ReadDir(dataDir)
	if err != nil || len(files) == 0 {
		Quotes = quotes
		log.Println(err)
		return
	}

	for _, file := range files {
		if err := loadQuotesFromFile(dataDir + "/" + file.Name()); err != nil {
			fmt.Println("Error loading file:", file.Name(), "-", err)
		}
	}
	return
}

func loadQuotesFromFile(filename string) error {
	switch {
	case strings.HasSuffix(filename, ".json"):
		return loadQuotesFromJSON(filename)
	case strings.HasSuffix(filename, ".yaml") || strings.HasSuffix(filename, ".yml"):
		return loadQuotesFromYAML(filename)
	case strings.HasSuffix(filename, ".txt"):
		return loadQuotesFromTXT(filename)
	case strings.HasSuffix(filename, ".csv"):
		return loadQuotesFromCSV(filename)
	default:
		fmt.Println("Unsupported file type:", filename)
	}
	return nil
}

func loadQuotesFromJSON(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &Quotes)
}

func loadQuotesFromYAML(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, &Quotes)
}

func loadQuotesFromTXT(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		Quotes = append(Quotes, scanner.Text())
	}
	return scanner.Err()
}

func loadQuotesFromCSV(filename string) error {
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

	for _, record := range records {
		Quotes = append(Quotes, record...)
	}
	return nil
}
