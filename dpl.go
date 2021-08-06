package main

import (
	"DiamondProtectorLink/config"
	"bytes"
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"unicode"
)

var redisClient *redis.Client
var ctx = context.Background()

func stripSpaces(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			// if the character is a space, drop it
			return -1
		}
		// else keep it in the string
		return r
	}, str)
}

func main() {
	out, errOut, err := executeCommand("whoami")
	if err != nil {
		log.Println(errOut)
		log.Fatal("Failed to execute whoami command!")
	}

	if stripSpaces(out) != "root" {
		log.Fatal("Failed to start DiamondProtectorLink cause the application must have sudo perms!")
	}

	// Load the config
	var loadedConfig = config.Get()

	log.Println("Connecting to Redis with " + (loadedConfig.RedisHostname + ":" + strconv.Itoa(loadedConfig.RedisPort)))
	redisClient = redis.NewClient(&redis.Options{
		Addr:     loadedConfig.RedisHostname + ":" + strconv.Itoa(loadedConfig.RedisPort),
		Password: loadedConfig.RedisPassword,
		DB:       0, // use default DB
	})

	ping := redisClient.Ping(context.Background())
	if strings.Contains(ping.String(), "refused") {
		log.Fatal("Error while connecting to redis, invalid credentials?!")
	}

	log.Println("Successfully connected to redis")

	pubsub := redisClient.Subscribe(ctx, "dpl-channel")
	channel := pubsub.Channel()

	go func() {
		for true {
			for cd := range channel {
				var command = cd.Payload
				_, errOut, err := executeCommand(command)
				if err != nil {
					return
				}

				if len(errOut) != 0 {
					log.Println("Error while executing command " + command)
					log.Println(errOut)
					log.Println(err)
				}
			}
		}
	}()

	select {}
}

func executeCommand(command string) (string, string, error) {
	cmd := exec.Command("sh", "-c", command)

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}
